package ericsson.core.nrf.dbproxy.functionservice;

import com.ericsson.geode.functionservice.ClearCode;
import com.ericsson.geode.functionservice.RemoteCacheDistributeLock;
import ericsson.core.nrf.dbproxy.DBProxyServer;
import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.config.GeodeConfig;
import java.net.InetAddress;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Timer;
import java.util.TimerTask;
import java.util.UUID;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicLong;
import org.apache.geode.cache.execute.Execution;
import org.apache.geode.cache.execute.FunctionService;
import org.apache.geode.cache.execute.ResultCollector;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class RemoteCacheMonitorThread {

  private static final Logger LOGGER = LogManager.getLogger(RemoteCacheMonitorThread.class);
  private static RemoteCacheMonitorThread instance;

  static {
    instance = null;
  }

  private AtomicLong putOperationCount = new AtomicLong();
  private String hostname = "";
  private int capacity = 100;
  private long sleepTime = 100;
  private Map<String, String> lockArgs = new HashMap<String, String>();
  private Map<String, String> unlockArgs = new HashMap<String, String>();

  public static synchronized RemoteCacheMonitorThread getInstance() {
    if (null == instance) {
      instance = new RemoteCacheMonitorThread();
    }

    return instance;
  }

  public int getCapacity() {
    return this.capacity;
  }

  public void setCapacity(int capacity) {
    LOGGER.debug("capacity : " + capacity);
    this.capacity = capacity;
  }

  public void incCacheOperationCount() {
    this.putOperationCount.incrementAndGet();
  }

  public void resetCacheOperationCount() {
    this.putOperationCount.getAndSet(0);
  }

  public long getCacheOperationCount() {
    return this.putOperationCount.longValue();
  }

  public void setHostname() {
    try {
      this.hostname = InetAddress.getLocalHost().getHostName();
    } catch (Exception e) {
      this.hostname = UUID.randomUUID().toString() + System.currentTimeMillis();
      LOGGER
          .error("Get HostName fail :" + e.toString() + "Generator UUID as hostname : " + hostname);
    }
  }

  public int lockDistributedlock() {
    int ret = ClearCode.GetLockFail;
    try {
      Execution execution = FunctionService
          .onRegion(ClientCacheService.getInstance().getRegion(Code.DISTRIBUTEDLOCK_INDICE))
          .setArguments(lockArgs);
      ResultCollector<Object, List> results = execution
          .execute(new RemoteCacheDistributeLock());
      ret = (int) results.getResult().get(0);
    } catch (Exception e) {
      LOGGER.error(e.toString());
      ret = ClearCode.GetLockFail;
    }

    return ret;
  }

  public int unlockDistributedlock() {
    int unLockRet = ClearCode.GetLockFail;
    try {
      Execution execution2 = FunctionService
          .onRegion(ClientCacheService.getInstance().getRegion(Code.DISTRIBUTEDLOCK_INDICE))
          .setArguments(unlockArgs);
      ResultCollector<Object, List> results2 = execution2
          .execute(new RemoteCacheDistributeLock());
      unLockRet = (int) results2.getResult().get(0);
    } catch (Exception e) {
      LOGGER.error(e.toString());
      unLockRet = ClearCode.GetLockFail;
    }

    return unLockRet;
  }

  public String getHostname() {
    return this.hostname;
  }

  public void start() {
    setHostname();
    LOGGER.debug("Hostname: " + hostname);
    if (!hostname.contains("discovery")) {
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

    lockArgs.put("lock", hostname);
    unlockArgs.put("unlock", hostname);

    AtomicBoolean b = new AtomicBoolean(false);

    Timer timer = new Timer();
    timer.scheduleAtFixedRate(new TimerTask() {
      @Override
      public void run() {
        b.set(true);
      }
    }, GeodeConfig.getRemoteCacheClearInterval(), GeodeConfig.getRemoteCacheClearInterval());

    Thread monitor = new Thread(() -> {
      while (DBProxyServer.getInstance().isRunning()) {
        while (getCacheOperationCount() >= GeodeConfig.getRemoteCachePutCount()
            || b.get()) {

          if (b.get()) {
            b.set(false);
          } else {
            resetCacheOperationCount();
          }

          try {
            int ret = lockDistributedlock();
            if (ret == ClearCode.GetLockSucc) {
              if (!RemoteCacheClearThread.getInstance().getClearFlag()) {
                RemoteCacheClearThread.getInstance().setClearFlag(true);
              }
              while (RemoteCacheClearThread.getInstance().getClearFlag()) {
                LOGGER.debug("Start Clear cache");
                Thread.sleep((long) ClearCode.UpdataDistributeInterval * 1000);
                LOGGER.debug(
                    "Clear cachenfprofiles not finished, still occupy distributedlock  cache capacity: "
                        + capacity);
                int relock = lockDistributedlock();
                if (relock != ClearCode.GetLockSucc) {
                  LOGGER.error("Clear cachenfprofiles not finished, but lost distributedlock");
                }
              }

              int unLockRet = unlockDistributedlock();
              if (unLockRet == ClearCode.GetLockFail) {
                LOGGER.error("After Capacity Clear, to UnLock Distributedlock fail");
              }
              LOGGER.debug("End Clear cache");
            } else {
              break;
            }
          } catch (Exception e) {
            LOGGER.error(e.toString());
          }
        }

        try {
          Thread.sleep(sleepTime);
        } catch (Exception e) {
          LOGGER.error(e.toString());
        }
      }
    });

    monitor.start();
  }
}
