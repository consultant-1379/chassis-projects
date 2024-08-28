package ericsson.core.nrf.dbproxy.config;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class EnvironmentConfig
{
    private static final Logger logger = LogManager.getLogger(EnvironmentConfig.class);

    private static EnvironmentConfig instance = null;

    private String db_info_conf;
    private String internal_conf;
    private String attributes_conf;
    private int grpc_port;

    private EnvironmentConfig()
    {
        reset();
    }

    public static synchronized EnvironmentConfig getInstance()
    {
        if(instance == null) {
            instance = new EnvironmentConfig();
        }
        return instance;
    }

    public boolean initialize()
    {
        String env_db_info_conf = "";
        try {
            env_db_info_conf = System.getenv("DB_INFO_CONF");
            if(null == env_db_info_conf || env_db_info_conf.isEmpty()) {
                logger.error("Env DB_INFO_CONF is not configured");
                return false;
            }
            db_info_conf = env_db_info_conf;
        } catch (NullPointerException | SecurityException e ) {
            logger.error("Exception occurs, DB_INFO_CONF = {} is configured", env_db_info_conf);
            return false;
        }

        String env_internal_conf = "";
        try {
            env_internal_conf = System.getenv("DBPROXY_INTERNAL_CONF");
            if(null == env_internal_conf || env_internal_conf.isEmpty()) {
                logger.error("Env DBPROXY_INTERNAL_CONF is not configured");
                return false;
            }
            internal_conf = env_internal_conf;
        } catch (NullPointerException | SecurityException e ) {
            logger.error("Exception occurs, DBPROXY_INTERNAL_CONF = {} is configured", env_internal_conf);
            return false;
        }

	String env_attributes_conf = "";
	try {
            env_attributes_conf = System.getenv("DBPROXY_ATTRIBUTES_CONF");
            if(null == env_attributes_conf || env_attributes_conf.isEmpty()) {
                logger.error("Env DBPROXY_ATTRIBUTES_CONF is not configured");
                return false;
            }
            attributes_conf = env_attributes_conf;
	}
	catch(NullPointerException | SecurityException e ) {
            logger.error("Exception occurs, DBPROXY_ATTRIBUTES_CONF = {} is configured", env_attributes_conf);
            return false;
        }

        String env_grpc_port = "";
        try {
            env_grpc_port = System.getenv("DBPROXY_GRPC_PORT");
            if(null == env_grpc_port || env_grpc_port.isEmpty()) {
                logger.error("Env DBPROXY_GRPC_PORT is not configured");
                return false;
            }
            grpc_port = Integer.parseInt(env_grpc_port);
        } catch (NumberFormatException | NullPointerException | SecurityException e) {
            logger.error("Exception occurs, DBPROXY_GRPC_PORT = {} is configured", env_grpc_port);
            return false;
        }

        logger.debug(toString());

        return true;
    }

    public void reset()
    {
        internal_conf = attributes_conf = "";
        grpc_port = 0;
    }

    public String getDBInfoConf()
    {
        return db_info_conf;
    }

    public String getInternalConf()
    {
        return internal_conf;
    }

    public String getAttributesConf()
    {
	return attributes_conf;
    }

    public int getGRPCPort()
    {
        return grpc_port;
    }

    public String toString()
    {

        StringBuilder sb = new StringBuilder("");
        sb.append("DB_INFO_CONF = " + db_info_conf + ",");
        sb.append("DBPROXY_INTERNAL_CONF = " + internal_conf + ",");
        sb.append("DBPROXY_ATTRIBUTES_CONF = " + attributes_conf + ",");
        sb.append("DBPROXY_GRPC_PORT = " + Integer.toString(grpc_port));
        return sb.toString();

    }
}
