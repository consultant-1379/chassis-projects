package ericsson.core.nrf.dbproxy.config;

import ericsson.core.nrf.dbproxy.common.Code;
import java.io.File;
import java.io.IOException;
import java.net.InetAddress;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import org.apache.commons.io.FileUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.json.JSONException;
import org.json.JSONObject;

public class GeodeConfig {

  private static final Logger LOGGER = LogManager.getLogger(GeodeConfig.class);

  private static final String LOG_1 = ", new name = ";
  private static int maxTransmitFragmentSize;
  private static int remoteCachePutCount;
  private static int remoteCacheClearInterval;
  private String kvdb_locator_name;
  private List<String> kvdb_locator_ip_list;
  private int kvdb_locator_port;
  private String[] region_list;
  private String pool_name;
  private int free_connection_timeout;
  private int idle_timeout;
  private int load_conditioning_interval;
  private int max_connections;
  private int min_connections;
  private int ping_interval;
  private boolean pr_single_hop_enabled;
  private int read_timeout;
  private int retry_attempts;
  private int socket_buffer_size;
  private int socket_connect_timeout;
  private boolean subscription_enabled;
  private int subscription_redundancy;
  private boolean thread_local_connections;

  public GeodeConfig() {
    kvdb_locator_ip_list = new ArrayList<>();
    kvdb_locator_port = 0;
    kvdb_locator_name = "";
    resetInternalValues();
  }

  public static int getRemoteCachePutCount() {
    return remoteCachePutCount;
  }

  public static int getRemoteCacheClearInterval() {
    return remoteCacheClearInterval;
  }

  public static int getMaxTransmitFragmentSize() {
    return maxTransmitFragmentSize;
  }

  private void resetInternalValues() {
    pool_name = "connection-pool";
    free_connection_timeout = 3000;
    idle_timeout = 5000;
    load_conditioning_interval = 120000;
    max_connections = 5000;
    min_connections = 10;
    ping_interval = 10000;
    pr_single_hop_enabled = false;
    read_timeout = 3000;
    retry_attempts = -1;
    socket_buffer_size = 32768;
    socket_connect_timeout = 3000;
    subscription_enabled = true;
    subscription_redundancy = -1;
    thread_local_connections = false;
    maxTransmitFragmentSize = 3145728;
    remoteCacheClearInterval = 1000;
    remoteCachePutCount = 100;
  }

  public boolean initialize() {
    try {
      String dbInfoConf = EnvironmentConfig.getInstance().getDBInfoConf();
      String dbContent = FileUtils.readFileToString(new File(dbInfoConf), "UTF-8");
      JSONObject dbJsonObject = new JSONObject(dbContent);

      kvdb_locator_name = dbJsonObject.getString("locator-server-name");
      String locatorPort = dbJsonObject.getString("locator-server-port");

      String regionNames = dbJsonObject.getString("region-names");
      region_list = regionNames.split(";");

      InetAddress[] address = InetAddress.getAllByName(kvdb_locator_name);
      for (int i = 0; i < address.length; i++) {
        kvdb_locator_ip_list.add(address[i].getCanonicalHostName());
      }

      kvdb_locator_port = Integer.parseInt(locatorPort);

      String internalConf = EnvironmentConfig.getInstance().getInternalConf();
      String content = FileUtils.readFileToString(new File(internalConf), "UTF-8");
      JSONObject jsonObject = new JSONObject(content);

      pool_name = jsonObject.getString("name");
      free_connection_timeout = jsonObject.getInt("free-connection-timeout");
      idle_timeout = jsonObject.getInt("idle-timeout");
      load_conditioning_interval = jsonObject.getInt("load-conditioning-interval");
      max_connections = jsonObject.getInt("max-connections");
      min_connections = jsonObject.getInt("min-connections");
      ping_interval = jsonObject.getInt("ping-interval");
      if (jsonObject.getString("pr-single-hop-enabled").equals("true")) {
        pr_single_hop_enabled = true;
      } else {
        pr_single_hop_enabled = false;
      }
      read_timeout = jsonObject.getInt("read-timeout");
      retry_attempts = jsonObject.getInt("retry-attempts");
      socket_buffer_size = jsonObject.getInt("socket-buffer-size");
      socket_connect_timeout = jsonObject.getInt("socket-connect-timeout");
      if (jsonObject.getString("subscription-enabled").equals("true")) {
        subscription_enabled = true;
      } else {
        subscription_enabled = false;
      }
      subscription_redundancy = jsonObject.getInt("subscription-redundancy");
      if (jsonObject.getString("thread-local-connections")
          .equals("true")) {
        thread_local_connections = true;
      } else {
        thread_local_connections = false;
      }
      maxTransmitFragmentSize = jsonObject.getInt("max-transmit-fragment-size");
      remoteCachePutCount = jsonObject.getInt("remote-cache-put-threshold");
      remoteCacheClearInterval = jsonObject.getInt("remote-cache-clear-interval");

    } catch (IOException | JSONException | NullPointerException | SecurityException | NumberFormatException e) {
      LOGGER.error(e.toString());
      LOGGER.warn("Fail to parse " + EnvironmentConfig.getInstance().getInternalConf());
      resetInternalValues();
    } catch (Exception e) {
      LOGGER.error(e.toString());
      LOGGER.warn("Fail to parse " + EnvironmentConfig.getInstance().getInternalConf());
      resetInternalValues();
    }

    LOGGER.debug(toString());

    return validate();
  }

  private boolean validate() {

    if (kvdb_locator_name.isEmpty()) {
      LOGGER.error("kvdb_locator_name is empty");
      return false;
    }

    if (kvdb_locator_ip_list.size() == 0) {
      LOGGER.error("No available kvdb locator ip");
      return false;
    }

    if (kvdb_locator_port <= 0) {
      LOGGER.error("kvdb_locator_port " + Integer.toString(kvdb_locator_port) + " is invalid");
      return false;
    }

    if (region_list.length == 0) {
      LOGGER.error("region_list is empty");
      return false;
    }

    LOGGER.debug("Geode Configuration Validation is Successful");
    return true;
  }

  public String getLocatorName() {
    return kvdb_locator_name;
  }

  public void addLocatorIP(String ip) {
    kvdb_locator_ip_list.add(ip);
  }

  public List<String> getLocatorIPList() {
    return kvdb_locator_ip_list;
  }

  public int getLocatorPort() {
    return kvdb_locator_port;
  }

  public void setLocatorPort(int port) {
    kvdb_locator_port = port;
  }

  public String[] getRegionList() {
    return region_list;
  }

  public String getPoolName() {
    return pool_name;
  }

  public int getFreeConnectionTimeout() {
    return free_connection_timeout;
  }

  public int getIdleTimeout() {
    return idle_timeout;
  }

  public int getLoadConditioningInterval() {
    return load_conditioning_interval;
  }

  public int getMaxConnections() {
    return max_connections;
  }

  public int getMinConnections() {
    return min_connections;
  }

  public int getPingInterval() {
    return ping_interval;
  }

  public boolean isPrSingleHopEnabled() {
    return pr_single_hop_enabled;
  }

  public int getReadTimeout() {
    return read_timeout;
  }

  public int getRetryAttempts() {
    return retry_attempts;
  }

  public int getSocketBufferSize() {
    return socket_buffer_size;
  }

  public int getSocketConnectTimeout() {
    return socket_connect_timeout;
  }

  public boolean isSubscriptionEnabled() {
    return subscription_enabled;
  }

  public int getSubscriptionRedundancy() {
    return subscription_redundancy;
  }

  public boolean isThreadLocalConnections() {
    return thread_local_connections;
  }

  public String toString() {
    StringBuilder sb = new StringBuilder("");
    sb.append("kvdb_locator_name = " + kvdb_locator_name + kvdb_locator_ip_list.toString() + ",");
    sb.append("kvdb_locator_port = " + Integer.toString(kvdb_locator_port) + ",");
    sb.append("region_name = " + Arrays.toString(region_list) + ",");
    sb.append("pool_name = " + pool_name + ",");
    sb.append("free_connection_timeout = " + Integer.toString(free_connection_timeout) + ",");
    sb.append("idle_timeout = " + Integer.toString(idle_timeout) + ",");
    sb.append("load_conditioning_interval = " + Integer.toString(load_conditioning_interval) + ",");
    sb.append("max_connections = " + Integer.toString(max_connections) + ",");
    sb.append("min_connections = " + Integer.toString(min_connections) + ",");
    sb.append("ping_interval = " + Integer.toString(ping_interval) + ",");
    sb.append("pr_single_hop_enabled = " + pr_single_hop_enabled + ",");
    sb.append("read_timeout = " + Integer.toString(read_timeout) + ",");
    sb.append("retry_attempts = " + Integer.toString(retry_attempts) + ",");
    sb.append("socket_buffer_size = " + Integer.toString(socket_buffer_size) + ",");
    sb.append("socket_connect_timeout = " + Integer.toString(socket_connect_timeout) + ",");
    sb.append("subscription_enabled = " + subscription_enabled + ",");
    sb.append("subscription_redundancy = " + subscription_redundancy + ",");
    sb.append("thread_local_connections = " + thread_local_connections + ",");
    sb.append("maxTransmitFragmentSize = " + maxTransmitFragmentSize + ",");
    sb.append("remote_cache_put_threshold = " + remoteCachePutCount + ",");
    sb.append("remoteCacheClearInterval = " + remoteCacheClearInterval);

    return sb.toString();
  }

  public int compare(GeodeConfig geodeConfig) {

    if (kvdb_locator_name.compareTo(geodeConfig.getLocatorName()) != 0) {
      LOGGER.warn(
          "KVDB locator name is changed, old name = " + kvdb_locator_name + LOG_1 + geodeConfig
              .getLocatorName());
      return Code.KVDB_LOCATOR_NAME_CHANGED;
    }

    if (kvdb_locator_port != geodeConfig.getLocatorPort()) {
      LOGGER.warn("KVDB locator port is changed, old port = " + Integer.toString(kvdb_locator_port)
          + ", new port = " + Integer.toString(geodeConfig.getLocatorPort()));
      return Code.KVDB_LOCATOR_PORT_CHANGED;
    }

    if (!Arrays.equals(region_list, geodeConfig.getRegionList())) {
      LOGGER.warn(
          "region list is changed, old = " + region_list + LOG_1 + geodeConfig.getRegionList());
      return Code.KVDB_REGION_NAME_CHANGED;
    }

    boolean found = false;
    for (String locatorIp : kvdb_locator_ip_list) {
      if (geodeConfig.getLocatorIPList().contains(locatorIp)) {
        found = true;
        break;
      }
    }

    if (!found) {
      LOGGER.warn("KVDB locator ip is changed, old ip list = " + kvdb_locator_ip_list.toString()
          + ", new ip list = " + geodeConfig.getLocatorIPList().toString());
      return Code.KVDB_LOCATOR_IP_CHANGED;
    }

    LOGGER.debug("Geode configuration is same as before");
    return Code.NOT_CHANGED;
  }
}
