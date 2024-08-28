package ericsson.core.nrf.dbproxy.config;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class EnvironmentConfig {

  private static final Logger LOGGER = LogManager.getLogger(EnvironmentConfig.class);

  private static EnvironmentConfig instance;

  static {
    instance = null;
  }

  private String db_info_conf;
  private String internal_conf;
  private String attributes_conf;
  private int grpc_port;

  private EnvironmentConfig() {
    reset();
  }

  public static synchronized EnvironmentConfig getInstance() {
    if (instance == null) {
      instance = new EnvironmentConfig();
    }
    return instance;
  }

  public boolean initialize() {
    String envDbInfoConf = "";
    try {
      envDbInfoConf = System.getenv("DB_INFO_CONF");
      if (null == envDbInfoConf || envDbInfoConf.isEmpty()) {
        LOGGER.error("Env DB_INFO_CONF is not configured");
        return false;
      }
      db_info_conf = envDbInfoConf;
    } catch (NullPointerException | SecurityException e) {
      LOGGER.error("Exception occurs, DB_INFO_CONF = {} is configured", envDbInfoConf);
      return false;
    }

    String envInternalConf = "";
    try {
      envInternalConf = System.getenv("DBPROXY_INTERNAL_CONF");
      if (null == envInternalConf || envInternalConf.isEmpty()) {
        LOGGER.error("Env DBPROXY_INTERNAL_CONF is not configured");
        return false;
      }
      internal_conf = envInternalConf;
    } catch (NullPointerException | SecurityException e) {
      LOGGER.error("Exception occurs, DBPROXY_INTERNAL_CONF = {} is configured", envInternalConf);
      return false;
    }

    String envAttributesConf = "";
    try {
      envAttributesConf = System.getenv("DBPROXY_ATTRIBUTES_CONF");
      if (null == envAttributesConf || envAttributesConf.isEmpty()) {
        LOGGER.error("Env DBPROXY_ATTRIBUTES_CONF is not configured");
        return false;
      }
      attributes_conf = envAttributesConf;
    } catch (NullPointerException | SecurityException e) {
      LOGGER.error("Exception occurs, DBPROXY_ATTRIBUTES_CONF = {} is configured",
          envAttributesConf);
      return false;
    }

    String envGrpcPort = "";
    try {
      envGrpcPort = System.getenv("DBPROXY_GRPC_PORT");
      if (null == envGrpcPort || envGrpcPort.isEmpty()) {
        LOGGER.error("Env DBPROXY_GRPC_PORT is not configured");
        return false;
      }
      grpc_port = Integer.parseInt(envGrpcPort);
    } catch (NumberFormatException | NullPointerException | SecurityException e) {
      LOGGER.error("Exception occurs, DBPROXY_GRPC_PORT = {} is configured", envGrpcPort);
      return false;
    }

    LOGGER.debug(toString());

    return true;
  }

  public void reset() {
    internal_conf = attributes_conf = "";
    grpc_port = 0;
  }

  public String getDBInfoConf() {
    return db_info_conf;
  }

  public String getInternalConf() {
    return internal_conf;
  }

  public String getAttributesConf() {
    return attributes_conf;
  }

  public int getGRPCPort() {
    return grpc_port;
  }

  public String toString() {

    StringBuilder sb = new StringBuilder("");
    sb.append("DB_INFO_CONF = " + db_info_conf + ",");
    sb.append("DBPROXY_INTERNAL_CONF = " + internal_conf + ",");
    sb.append("DBPROXY_ATTRIBUTES_CONF = " + attributes_conf + ",");
    sb.append("DBPROXY_GRPC_PORT = " + Integer.toString(grpc_port));
    return sb.toString();

  }
}
