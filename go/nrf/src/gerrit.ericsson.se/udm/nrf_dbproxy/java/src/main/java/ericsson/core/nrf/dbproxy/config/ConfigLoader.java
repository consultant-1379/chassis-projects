package ericsson.core.nrf.dbproxy.config;

import java.util.concurrent.Semaphore;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import ericsson.core.nrf.dbproxy.DBProxyServer;
import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;

public class ConfigLoader
{
    private static final Logger logger = LogManager.getLogger(ConfigLoader.class);
    private static ConfigLoader instance = null;
    private static final long CONF_CHECKING_INTERVAL = 15000;
    private final Semaphore available = new Semaphore(1, true);
    private long configuration_checking_interval;

    private ConfigLoader()
    {
        this.configuration_checking_interval = CONF_CHECKING_INTERVAL;
    }

    public static synchronized ConfigLoader getInstance()
    {
        if(null == instance) {
            instance = new ConfigLoader();
        }
        return instance;
    }

    public void start()
    {
        while(DBProxyServer.getInstance().isRunning()) {
            if(EnvironmentConfig.getInstance().initialize() == true && AttributeConfig.getInstance().load() == true) {
                break;
            }

            EnvironmentConfig.getInstance().reset();

            try {
                logger.debug("Fail to read environment parameters, sleep " + Long.toString(ConfigLoader.this.configuration_checking_interval/1000) + " seconds and try again");
                Thread.sleep(ConfigLoader.this.configuration_checking_interval);
            } catch(Exception e) {
                logger.error(e.toString());
            }
        }

        Thread timer = new Thread(() -> {

            while(DBProxyServer.getInstance().isRunning()) {
                GeodeConfig geode_config = new GeodeConfig();
                if(geode_config.initialize() == true) {
                    ConfigLoader.this.acquire();
                    ClientCacheService.getInstance().apply(geode_config);
                    ConfigLoader.this.release();
                }
                try {
                    logger.debug("Sleep " + Long.toString(ConfigLoader.this.configuration_checking_interval/1000) + " seconds and load geode configuration parameters again");
                    Thread.sleep(ConfigLoader.this.configuration_checking_interval);
                } catch (Exception e) {
                    logger.error(e.toString());
                }
            }
        });

        timer.start();
    }

    private void acquire()
    {
        try {
            this.available.acquire();
        } catch (InterruptedException e) {
            logger.error(e.toString());
            Thread.currentThread().interrupt();
        }
    }

    private void release()
    {
        this.available.release();
    }
}
