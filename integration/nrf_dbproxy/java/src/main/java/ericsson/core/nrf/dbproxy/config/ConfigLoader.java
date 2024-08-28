package ericsson.core.nrf.dbproxy.config;

import ericsson.core.nrf.dbproxy.DBProxyServer;
import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import java.util.concurrent.Semaphore;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class ConfigLoader {

  private static final Logger LOGGER = LogManager.getLogger(ConfigLoader.class);
  private static final long CONF_CHECKING_INTERVAL = 15000;
  private static ConfigLoader instance;

  static {
    instance = null;
  }

  private final Semaphore available = new Semaphore(1, true);
  private long configuration_checking_interval;

  private ConfigLoader() {
    this.configuration_checking_interval = CONF_CHECKING_INTERVAL;
  }

  public static synchronized ConfigLoader getInstance() {
    if (null == instance) {
      instance = new ConfigLoader();
    }
    return instance;
  }

  public void start() {
    while (DBProxyServer.getInstance().isRunning()) {
      if (EnvironmentConfig.getInstance().initialize()
          && AttributeConfig.getInstance().load()) {
        break;
      }

      EnvironmentConfig.getInstance().reset();

      try {
        LOGGER.debug("Fail to read environment parameters, sleep " + Long
            .toString(ConfigLoader.this.configuration_checking_interval / 1000)
            + " seconds and try again");
        Thread.sleep(ConfigLoader.this.configuration_checking_interval);
      } catch (Exception e) {
        LOGGER.error(e.toString());
      }
    }

    Thread timer = new Thread(() -> {

      while (DBProxyServer.getInstance().isRunning()) {
        GeodeConfig geodeConfig = new GeodeConfig();
        if (geodeConfig.initialize()) {
          ConfigLoader.this.acquire();
          ClientCacheService.getInstance().apply(geodeConfig);
          ConfigLoader.this.release();
        }
        try {
          LOGGER.debug(
              "Sleep " + Long.toString(ConfigLoader.this.configuration_checking_interval / 1000)
                  + " seconds and load geode configuration parameters again");
          LOGGER.debug(
              "KVDB locator ip is got, current ip list = " + geodeConfig.getLocatorIPList()
                  .toString());
          Thread.sleep(ConfigLoader.this.configuration_checking_interval);
        } catch (Exception e) {
          LOGGER.error(e.toString());
        }
      }
    });

    timer.start();
  }

  private void acquire() {
    try {
      this.available.acquire();
    } catch (InterruptedException e) {
      LOGGER.error(e.toString());
      Thread.currentThread().interrupt();
    }
  }

  private void release() {
    this.available.release();
  }
}
