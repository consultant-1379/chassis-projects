package ericsson.core.nrf.dbproxy.log;

import ericsson.core.nrf.dbproxy.fsnotify.FileMonitor;
import ericsson.core.nrf.dbproxy.fsnotify.EventHandler;
import java.io.File;
import java.io.IOException;
import org.apache.commons.io.FileUtils;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;
import org.everit.json.schema.Schema;
import org.everit.json.schema.loader.SchemaLoader;
import org.everit.json.schema.ValidationException;
import java.util.Map;
import java.util.HashMap;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;
import org.apache.logging.log4j.core.config.Configurator;
import org.apache.logging.log4j.Level;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class LogLevelControl implements EventHandler{

  private static LogLevelControl instance;
  private static final Logger LOGGER = LogManager.getLogger(LogLevelControl.class);
  private static final Level DEFAULT_LEVEL = Level.INFO;
  private static final String OTHER_MODULE_NAME = "OTHERS";
  private static final String LOG_JSON_SCHEMA = "{\n"
      + "  \"$id\": \"file://schema/LogControl#\",\n"
      + "  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n"
      + "  \"title\": \"Log level control\",\n"
      + "  \"description\": \"Definition of a log level control file INT.LOG.CTRL at /etc/adp/logcontrol.json\",\n"
      + "  \"eric-adp-version\": \"0.1.0\",\n"
      + "  \"type\": \"array\",\n"
      + "  \"uniqueItems\": true,\n"
      + "  \"items\": {\n"
      + "    \"type\": \"object\",\n"
      + "    \"properties\": {\n"
      + "      \"container\": {\n"
      + "        \"type\": \"string\",\n"
      + "        \"description\": \"Name of the container producing the log event.\"\n"
      + "      },\n"
      + "      \"severity\": {\n"
      + "        \"type\": \"string\",\n"
      + "        \"enum\": [\n"
      + "          \"debug\",\n"
      + "          \"info\"\n"
      + "        ],\n"
      + "        \"default\": \"info\",\n"
      + "        \"description\": \"Log event severity level.\"\n"
      + "      },\n"
      + "      \"customFilters\": {\n"
      + "        \"type\": \"array\",\n"
      + "        \"items\": {\n"
      + "          \"$ref\": \"#/components/schemas/CustomFilter\"\n"
      + "        },\n"
      + "        \"description\": \"Optional list of log events filters.\"\n"
      + "      }\n"
      + "    },\n"
      + "    \"additionalProperties\": false,\n"
      + "    \"required\": [\n"
      + "      \"container\",\n"
      + "      \"severity\"\n"
      + "    ]\n"
      + "  },\n"
      + "  \"components\": {\n"
      + "    \"schemas\": {\n"
      + "      \"CustomFilter\": {\n"
      + "        \"type\": \"object\",\n"
      + "        \"additionalProperties\": false,\n"
      + "        \"properties\": {\n"
      + "          \"pod\": {\n"
      + "            \"type\": \"string\",\n"
      + "            \"description\": \"Name of the pod producing the log event.\"\n"
      + "          },\n"
      + "          \"module\": {\n"
      + "            \"type\": \"string\",\n"
      + "            \"description\": \"Name of the module producing the log event.\"\n"
      + "          }\n"
      + "        }\n"
      + "      }\n"
      + "    }\n"
      + "  }\n"
      + "}";
  private String podName;
  private String containerName;
  private String configFile;
  private String configFileName;
  private Map<String, String> moduleLevel;
  private String rootLevel;
  private FileMonitor fileMonitor;
  private Schema jsonSchema;
  private int fileUpdateCount;
  private Lock lock;

  static {
    instance = null;
  }

  private LogLevelControl() {
    podName = "";
    containerName = "";
    configFile = "";
    configFileName = "";
    moduleLevel = new HashMap<String, String>();
    rootLevel = "INFO";
    fileMonitor = null;
    jsonSchema = null;
    fileUpdateCount = 0;
    lock = new ReentrantLock();
  }

  public static synchronized LogLevelControl getInstance() {
    if (null == instance) {
      instance = new LogLevelControl();
    }

    return instance;
  }

  public boolean initialize() {
    String podNameTmp = "";
    try {
      podNameTmp = System.getenv("POD_NAME");
      if (null == podNameTmp || podNameTmp.isEmpty()) {
        LOGGER.error("Env POD_NAME is not configured");
        return false;
      }
      this.podName = podNameTmp;
    } catch (NullPointerException | SecurityException e) {
      LOGGER.error("Exception occurs, POD_NAME = {} is configured",
          podNameTmp);
      return false;
    }

    String containerNameTmp = "";
    try {
      containerNameTmp = System.getenv("CONTAINER_NAME");
      if (null == containerNameTmp || containerNameTmp.isEmpty()) {
        LOGGER.error("Env CONTAINER_NAME is not configured");
        return false;
      }
      this.containerName = containerNameTmp;
    } catch (NullPointerException | SecurityException e) {
      LOGGER.error("Exception occurs, CONTAINER_NAME = {} is configured",
          containerNameTmp);
      return false;
    }

    String configFileTmp = "";
    try {
      configFileTmp = System.getenv("LOG_CONFIG_FILE");
      if (null == configFileTmp || configFileTmp.isEmpty()) {
        LOGGER.error("Env LOG_CONFIG_FILE is not configured");
        return false;
      }
      this.configFile = configFileTmp;
    } catch (NullPointerException | SecurityException e) {
      LOGGER.error("Exception occurs, LOG_CONFIG_FILE = {} is configured",
          configFileTmp);
      return false;
    }

    int lastIndex = configFile.lastIndexOf("/");
    if (-1 == lastIndex) {
      LOGGER.error("Fail to get the absolute path of log config file directory");
      return false;
    }
    configFileName = configFile.substring(lastIndex + 1);


    if (!registerModule()) {
      LOGGER.error("Fail to register module");
      return false;
    }

    if (!loadJsonSchema()) {
      LOGGER.error("Fail to load json schema");
      return false;
    }

    if (!setLevel()) {
      LOGGER.error("Fail to set log level");
      return false;
    }

    if (!monitorConf()) {
      LOGGER.error("Fail to monitor log level config file");
      return false;
    }

    return true;
  }

  public boolean registerModule() {
    if ("eric-nrf-accesstoken-dbproxy".equals(this.containerName)) {
      moduleLevel.put("ericsson.core.nrf.dbproxy.DBProxyServer", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceImpl", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.Executor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.ExecutorManager", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nfprofile.NFProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.clientcache.region.ClientRegion", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nfprofile.NFProfileGetHelper", DEFAULT_LEVEL.name());

    } else if ("eric-nrf-disc-dbproxy".equals(this.containerName)) {
      moduleLevel.put("ericsson.core.nrf.dbproxy.DBProxyServer", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceImpl", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.Executor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.ExecutorManager", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.imsiprefixprofile.ImsiprefixProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.gpsiprefixprofile.GpsiprefixProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nfprofile.NFProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.cachenfprofile.CacheNFProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.cachenfprofile.CacheNFProfilePutExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nrfaddress.NRFAddressGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfileProcesser", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.clientcache.region.ClientRegion", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.imsiprefixprofile.ImsiprefixProfileGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.gpsiprefixprofile.GpsiprefixProfileGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nfprofile.NFProfileGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.cachenfprofile.CacheNFProfilePutHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.cachenfprofile.CacheNFProfileGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nrfaddress.NRFAddressGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nrfprofile.NRFProfileGetHelper", DEFAULT_LEVEL.name());

    } else if ("eric-nrf-nfm-dbproxy".equals(this.containerName)) {
      moduleLevel.put("ericsson.core.nrf.dbproxy.DBProxyServer", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceImpl", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.Executor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.ExecutorManager", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.groupprofile.GroupProfileDeleteExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.groupprofile.GroupProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.groupprofile.GroupProfilePutExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.gpsiprofile.GpsiProfileDeleteExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.gpsiprofile.GpsiProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.gpsiprofile.GpsiProfilePutExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nfprofile.NFProfileDeleteExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nfprofile.NFProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nfprofile.NFProfilePutExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.common.ImsiPrefixProfilesUtil", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.common.GpsiPrefixProfilesUtil", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nrfaddress.NRFAddressGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.subscription.SubscriptionDeleteExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.subscription.SubscriptionGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.subscription.SubscriptionPutExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfileDeleteExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfilePutExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfileProcesser", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.clientcache.region.ClientRegion", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.groupprofile.GroupProfileDelHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.groupprofile.GroupProfilePutHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.groupprofile.GroupProfileGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.gpsiprofile.GpsiProfileDelHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.gpsiprofile.GpsiProfilePutHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.gpsiprofile.GpsiProfileGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nfprofile.NFProfileDelHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nfprofile.NFProfilePutHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nfprofile.NFProfileGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nrfaddress.NRFAddressGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.subscription.SubscriptionDelHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.subscription.SubscriptionPutHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.subscription.SubscriptionGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nrfprofile.NRFProfileDelHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nrfprofile.NRFProfilePutHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nrfprofile.NRFProfileGetHelper", DEFAULT_LEVEL.name());

    } else if ("eric-nrf-notify-dbproxy".equals(this.containerName)) {
      moduleLevel.put("ericsson.core.nrf.dbproxy.DBProxyServer", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceImpl", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.Executor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.ExecutorManager", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.groupprofile.GroupProfileDeleteExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.gpsiprofile.GpsiProfileDeleteExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nfprofile.NFProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nfprofile.NFProfilePutExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.common.ImsiPrefixProfilesUtil", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.common.GpsiPrefixProfilesUtil", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.subscription.SubscriptionGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfilePutExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfileProcesser", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.clientcache.region.ClientRegion", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.groupprofile.GroupProfileDelHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.gpsiprofile.GpsiProfileDelHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nfprofile.NFProfilePutHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nfprofile.NFProfileGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.subscription.SubscriptionGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nrfprofile.NRFProfilePutHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nrfprofile.NRFProfileGetHelper", DEFAULT_LEVEL.name());

    } else if ("eric-nrf-prov-dbproxy".equals(this.containerName)) {
      moduleLevel.put("ericsson.core.nrf.dbproxy.DBProxyServer", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceImpl", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.Executor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.ExecutorManager", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.groupprofile.GroupProfileDeleteExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.groupprofile.GroupProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.groupprofile.GroupProfilePutExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.gpsiprofile.GpsiProfileDeleteExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.gpsiprofile.GpsiProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.gpsiprofile.GpsiProfilePutExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.common.ImsiPrefixProfilesUtil", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.imsiprefixprofile.ImsiprefixProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.common.GpsiPrefixProfilesUtil", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.gpsiprefixprofile.GpsiprefixProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nfprofile.NFProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nrfaddress.NRFAddressDeleteExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nrfaddress.NRFAddressGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nrfaddress.NRFAddressPutExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.subscription.SubscriptionGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfileGetExecutor", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.clientcache.region.ClientRegion", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.groupprofile.GroupProfileDelHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.groupprofile.GroupProfilePutHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.groupprofile.GroupProfileGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.gpsiprofile.GpsiProfileDelHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.gpsiprofile.GpsiProfilePutHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.gpsiprofile.GpsiProfileGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nfprofile.NFProfileGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nrfaddress.NRFAddressDelHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nrfaddress.NRFAddressPutHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nrfaddress.NRFAddressGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.subscription.SubscriptionGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.nrfprofile.NRFProfileGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfileProcesser", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.imsiprefixprofile.ImsiprefixProfileGetHelper", DEFAULT_LEVEL.name());
      moduleLevel.put("ericsson.core.nrf.dbproxy.helper.gpsiprefixprofile.GpsiprefixProfileGetHelper", DEFAULT_LEVEL.name());

    } else {
      LOGGER.error("unknown container name");
      return false;
    }

    return true;
  }

  public boolean loadJsonSchema() {
    try {
      JSONObject schemaObject = new JSONObject(this.LOG_JSON_SCHEMA);
      jsonSchema = SchemaLoader.load(schemaObject);
    } catch (JSONException e) {
      LOGGER.error(e.toString());
      return false;
    } catch (Exception e) {
      LOGGER.error(e.toString());
      return false;
    }

    if (null == jsonSchema) {
      LOGGER.error("jsonSchema is null");
      return false;
    }

    return true;
  }

  public boolean monitorConf() {
    //should monitor the directory
    String configFilePath = "";
    int lastIndex = configFile.lastIndexOf("/");
    if (-1 == lastIndex) {
      LOGGER.error("Fail to get the absolute path of log config file directory");
      return false;
    }
    configFilePath = configFile.substring(0, lastIndex);
    fileMonitor = new FileMonitor(configFilePath, this, 1000);
    try {
      fileMonitor.start();
    } catch (Exception e) {
      LOGGER.error(e.toString());
      return false;
    }

    return true;
  }

  public boolean setLevel() {
    JSONArray jsonArray = loadConf();
    if (null == jsonArray) {
      LOGGER.error("Fail to load log level config file {}", this.configFile);
      return false;
    }

    if (!validateConf(jsonArray)) {
      LOGGER.error("Fail to validate log level config file {}", this.configFile);
      return false;
    }

    if (!parseConf(jsonArray)) {
      LOGGER.error("Fail to parse log level config file {}", this.configFile);
      return false;
    }

    return true;
  }

  public JSONArray loadConf() {
    JSONArray jsonArray = null;
    try {
      String content = FileUtils.readFileToString(new File(this.configFile), "UTF-8");
      jsonArray = new JSONArray(content);
    } catch (IOException | JSONException e) {
      LOGGER.error(e.toString());
      jsonArray = null;
    } catch (Exception e) {
      LOGGER.error(e.toString());
      jsonArray = null;
    } finally {
      return jsonArray;
    }
  }

  public boolean validateConf(JSONArray jsonArray) {
    if (null == jsonArray) {
      LOGGER.warn("The file content is empty");
      return false;
    }

    if (0 == jsonArray.length()) {
      LOGGER.warn("Zero items are configured");
      return false;
    }

    try {
      jsonSchema.validate(jsonArray);
    } catch (ValidationException e) {
      LOGGER.error(e.toString());
      return false;
    } catch (Exception e) {
      LOGGER.error(e.toString());
      return false;
    }

    return true;
  }

  public boolean parseConf(JSONArray jsonArray) {
    if (null == jsonArray) {
      LOGGER.warn("The file content is empty");
      return false;
    }

    if (0 == jsonArray.length()) {
      LOGGER.warn("Zero items are configured");
      return false;
    }

    boolean allModuleMatched = false;
    Map<String, String> matchedModuleList = new HashMap<String, String>();
    String configedSeverity = DEFAULT_LEVEL.name();

    //filter matchedModuleList
    try {
      for (int i = 0; i < jsonArray.length(); i++) {
        JSONObject obj = (JSONObject) jsonArray.get(i);
        String container = obj.getString("container");
        if (container != null && container.equals(this.containerName)) {
          configedSeverity = obj.getString("severity");
          if (!obj.has("customFilters") || obj.isNull("customFilters")) {
            allModuleMatched = true;
            break;
          }

          JSONArray filters = obj.getJSONArray("customFilters");
          if (0 == filters.length()) {
            allModuleMatched = true;
            break;
          }

          for (int j = 0; j < filters.length(); j++) {
            JSONObject filter = (JSONObject) filters.get(j);
            boolean podMatched = false;
            if (!filter.has("pod") || filter.isNull("pod") || filter.getString("pod").equals(this.podName)) {
              podMatched = true;
            }
            if (podMatched) {
              if (!filter.has("module") || filter.isNull("module")) {
                allModuleMatched = true;
                break;
              }
              matchedModuleList.put(filter.getString("module"), configedSeverity);
            }
          }

          break;
        }
      }
    } catch (JSONException | NullPointerException e) {
      LOGGER.error(e.toString());
      return false;
    } catch (Exception e) {
      LOGGER.error(e.toString());
      return false;
    }

    //set level
    Map<String, String> moduleLevelTmp = new HashMap<String, String>();
    if (allModuleMatched) {
      for (Map.Entry<String, String> entry : moduleLevel.entrySet()) {
        if (entry.getKey().isEmpty()) {
          continue;
        }
        setLevelToLogger(entry.getKey(), toLevel(configedSeverity));
        moduleLevelTmp.put(entry.getKey(), toLevel(configedSeverity).name());
      }
      setRootLevel(toLevel(configedSeverity));
      this.rootLevel = toLevel(configedSeverity).name();
    } else {
      for (Map.Entry<String, String> entry : moduleLevel.entrySet()) {
        if (matchedModuleList.containsKey(entry.getKey())) {
          setLevelToLogger(entry.getKey(), toLevel(configedSeverity));
          moduleLevelTmp.put(entry.getKey(), toLevel(configedSeverity).name());
        } else {
          setLevelToLogger(entry.getKey(), this.DEFAULT_LEVEL);
          moduleLevelTmp.put(entry.getKey(), this.DEFAULT_LEVEL.name());
        }
      }
      if (matchedModuleList.containsKey(this.OTHER_MODULE_NAME)) {
        setRootLevel(toLevel(configedSeverity));
        this.rootLevel = toLevel(configedSeverity).name();
      } else {
        setRootLevel(this.DEFAULT_LEVEL);
        this.rootLevel = this.DEFAULT_LEVEL.name();
      }
    }

    for (Map.Entry<String, String> entry : moduleLevelTmp.entrySet()) {
      moduleLevel.put(entry.getKey(), entry.getValue());
    }
    moduleLevelTmp.clear();

    return true;
  }

  public boolean resetLevel() {
    return setLevel();
  }

  public void handle(String eventName, int op) {
    if (eventName.equals(configFileName) && op == EventHandler.FILE_CHANGE) {
      this.lock.lock();
      try {
        this.fileUpdateCount++;
        if (2 == this.fileUpdateCount) {
          this.fileUpdateCount = 0;
          LOGGER.info("log level config file {} is changed, now to reset log level", this.configFile);
          if (!resetLevel()) {
            LOGGER.error("Fail to reset log level. Keep the current level {}", toString());
          } else {
            LOGGER.info("Succeed to reset log level. It's {}", toString());
          }
        }
      } catch (Exception e) {
        LOGGER.error(e.toString());
      } finally {
        this.lock.unlock();
      }
    }
  }

  public Level toLevel(String level) {
    if ("debug".equals(level)) {
      return Level.DEBUG;
    } else if ("info".equals(level)) {
      return Level.INFO;
    } else {
      return Level.INFO;
    }
  }

  public void setLevelToLogger(String loggerName, final Level level) {
    Configurator.setLevel(loggerName, level);
  }

  public void setRootLevel(final Level level) {
    Configurator.setRootLevel(level);
  }

  public String toString() {
    String str = "";
    for (Map.Entry<String, String> entry : moduleLevel.entrySet()) {
      if (str.isEmpty()) {
        str = entry.getKey() + ":" + entry.getValue();
      } else {
        str = str + "; " + entry.getKey() + ":" + entry.getValue();
      }
    }
    str = str + "; " + this.OTHER_MODULE_NAME + ":" + this.rootLevel;

    return "<" + str + ">";
  }
}
