package ericsson.core.nrf.dbproxy.functionservice;

import com.ericsson.geode.functionservice.ClearCode;
import com.ericsson.geode.functionservice.RemoteCacheClearService;
import ericsson.core.nrf.dbproxy.DBProxyServer;
import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import java.util.List;
import org.apache.geode.cache.execute.Execution;
import org.apache.geode.cache.execute.FunctionService;
import org.apache.geode.cache.execute.ResultCollector;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class RemoteCacheClearThread {

  private static final Logger LOGGER = LogManager.getLogger(RemoteCacheClearThread.class);
  private static RemoteCacheClearThread instance;

  static {
    instance = null;
  }

  private boolean running;
  private boolean clearFlag;

  public static synchronized RemoteCacheClearThread getInstance() {
    if (null == instance) {
      instance = new RemoteCacheClearThread();
    }
    return instance;
  }

  public boolean getClearFlag() {
    return this.clearFlag;
  }

  public void setClearFlag(boolean b) {
    this.clearFlag = b;
  }

  public boolean isRunning() {
    return this.running;
  }

  public void start() {
    if (!RemoteCacheMonitorThread.getInstance().getHostname().contains("discovery")) {
      return;
    }

    while (!ClientCacheService.getInstance().isAvailable()) {
      LOGGER.debug("ClientCache not available wait 1 seconds");
      try {
        Thread.sleep(1000);
      } catch (Exception e) {
        LOGGER.error(e.toString());
      }
    }

    Thread clear = new Thread(() -> {
      running = true;
      while (DBProxyServer.getInstance().isRunning()) {
        while (clearFlag) {
          try {
            Execution execution = FunctionService
                .onRegion(ClientCacheService.getInstance().getRegion(Code.CACHENFPROFILE_INDICE))
                .setArguments(RemoteCacheMonitorThread.getInstance().getCapacity());
            ResultCollector<Object, List> results = execution
                .execute(new RemoteCacheClearService());
            if ((int) results.getResult().get(0) != ClearCode.ClearSuccess) {
              LOGGER.error("Cache Region Capacity Clear Fail");
            } else {
              LOGGER.debug("Cache Region Capacity Clear Success");
            }
          } catch (Exception e) {
            LOGGER.error(e.toString());
          }
          clearFlag = false;
        }
        try {
          Thread.sleep(500);
        } catch (Exception e) {
          LOGGER.error(e.toString());
        }
      }
    });

    clear.start();
  }
}
