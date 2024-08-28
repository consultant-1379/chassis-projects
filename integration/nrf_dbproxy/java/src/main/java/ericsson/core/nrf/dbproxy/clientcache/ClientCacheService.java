package ericsson.core.nrf.dbproxy.clientcache;

import ericsson.core.nrf.dbproxy.clientcache.region.ClientRegion;
import ericsson.core.nrf.dbproxy.clientcache.state.ClientCacheServiceInitState;
import ericsson.core.nrf.dbproxy.clientcache.state.ClientCacheServiceState;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.FragmentSessionManagement;
import ericsson.core.nrf.dbproxy.config.GeodeConfig;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import org.apache.geode.cache.CacheClosedException;
import org.apache.geode.cache.CacheTransactionManager;
import org.apache.geode.cache.CacheWriterException;
import org.apache.geode.cache.CacheXmlException;
import org.apache.geode.cache.Region;
import org.apache.geode.cache.RegionExistsException;
import org.apache.geode.cache.TimeoutException;
import org.apache.geode.cache.client.ClientCache;
import org.apache.geode.cache.client.ClientCacheFactory;
import org.apache.geode.cache.client.Pool;
import org.apache.geode.cache.client.PoolFactory;
import org.apache.geode.cache.client.PoolManager;
import org.apache.geode.pdx.ReflectionBasedAutoSerializer;
import org.apache.geode.security.AuthenticationFailedException;
import org.apache.geode.security.AuthenticationRequiredException;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class ClientCacheService {

  private static final Logger LOGGER = LogManager.getLogger(ClientCacheService.class);

  private static ClientCacheService instance;

  static {
    instance = null;
  }

  private static ClientCache clientCache;

  static {
    clientCache = null;
  }

  private Map<String, ClientRegion> clientRegions;

  private ClientCacheServiceState state;
  private boolean available;
  private GeodeConfig config;

  private ClientCacheService() {
    this.clientRegions = new HashMap<String, ClientRegion>();
    reset();
  }

  public static synchronized ClientCacheService getInstance() {
    if (null == instance) {
      instance = new ClientCacheService();
    }
    return instance;
  }

  public void setState(ClientCacheServiceState serviceState) {
    state = serviceState;
  }

  public void apply(GeodeConfig geodeConfig) {
    LOGGER.debug("Apply geode configuration on client cache service");
    state.apply(geodeConfig);
  }

  public int put(String regionName, Object regionKey, Object regionItem) {
    ClientRegion clientRegion = clientRegions.get(regionName);
    return clientRegion.put(regionKey, regionItem);
  }

  public int putAll(String regionName, Map<Object, Object> kvMap) {
    ClientRegion clientRegion = clientRegions.get(regionName);
    return clientRegion.putAll(kvMap);
  }

  public int delete(String regionName, Object regionKey) {
    ClientRegion clientRegion = clientRegions.get(regionName);
    return clientRegion.delete(regionKey, true);
  }

  public int deleteAll(String regionName, List<String> keyList) {
    ClientRegion clientRegion = clientRegions.get(regionName);
    return clientRegion.deleteAll(keyList);
  }

  public int delete(String regionName, Object regionKey, boolean withTransaction) {
    ClientRegion clientRegion = clientRegions.get(regionName);
    return clientRegion.delete(regionKey, withTransaction);
  }

  public ExecutionResult getByID(String regionName, Object regionKey) {
    ClientRegion clientRegion = clientRegions.get(regionName);
    return clientRegion.getByID(regionKey);
  }

  public ExecutionResult getAllByID(String regionName, List<Long> idList) {
    ClientRegion clientRegion = clientRegions.get(regionName);
    return clientRegion.getAllByID(idList);
  }

  public ExecutionResult getByFilter(String regionName, String queryString) {
    ClientRegion clientRegion = clientRegions.get(regionName);
    return clientRegion.getByFilter(regionName, queryString);
  }

  public ExecutionResult getByCountFilter(String regionName, String queryString) {
    ClientRegion clientRegion = clientRegions.get(regionName);
    return clientRegion.getByCountFilter(queryString);
  }

  public ExecutionResult getByFragSessionId(String regionName, String fragmentSessionId) {
    return FragmentSessionManagement.getInstance().get(regionName, fragmentSessionId);
  }

  public boolean initialize(GeodeConfig geodeConfig) {
    if (!createClientCache(geodeConfig)) {
      return false;
    }

    String[] regions = geodeConfig.getRegionList();
    for (int i = 0; i < regions.length; i++) {
      ClientRegion region = new ClientRegion(regions[i]);
      if (!region.initialize(clientCache, geodeConfig.getPoolName(),
          geodeConfig.isSubscriptionEnabled())) {
        return false;
      }
      clientRegions.put(regions[i], region);
    }

    config = geodeConfig;
    available = true;
    return true;
  }

  private boolean createClientCache(GeodeConfig geodeConfig) {
    ClientCacheFactory factory = new ClientCacheFactory();
    try {
      clientCache = factory.setPdxSerializer(
          new ReflectionBasedAutoSerializer("ericsson.core.nrf.dbproxy.clientcache.schema.*"))
          .create();
    } catch (CacheXmlException | TimeoutException | RegionExistsException | CacheWriterException | IllegalStateException | AuthenticationFailedException | AuthenticationRequiredException e) {
      LOGGER.error(e.toString());
      return false;
    }

    PoolFactory poolFactory = PoolManager.createFactory();
    List<String> kvdbLocatorIpList = geodeConfig.getLocatorIPList();
    int kvdbLocatorPort = geodeConfig.getLocatorPort();
    for (String locatorIp : kvdbLocatorIpList) {
      try {
        poolFactory = poolFactory.addLocator(locatorIp, kvdbLocatorPort);
      } catch (IllegalArgumentException | IllegalStateException e) {
        LOGGER.error(e.toString());
        LOGGER.error(
            "Fail to add locator = " + locatorIp + ":" + Integer.toString(kvdbLocatorPort));
        return false;
      } catch (Exception e) {
        LOGGER.error(e.toString());
        LOGGER.error(
            "Fail to add locator = " + locatorIp + ":" + Integer.toString(kvdbLocatorPort));
        return false;
      }
    }

    try {
      poolFactory.setFreeConnectionTimeout(geodeConfig.getFreeConnectionTimeout())
          .setIdleTimeout(geodeConfig.getIdleTimeout())
          .setLoadConditioningInterval(geodeConfig.getLoadConditioningInterval())
          .setMaxConnections(geodeConfig.getMaxConnections())
          .setMinConnections(geodeConfig.getMinConnections())
          .setPingInterval(geodeConfig.getPingInterval())
          .setPRSingleHopEnabled(geodeConfig.isPrSingleHopEnabled())
          .setReadTimeout(geodeConfig.getReadTimeout())
          .setRetryAttempts(geodeConfig.getRetryAttempts())
          .setSocketBufferSize(geodeConfig.getSocketBufferSize())
          .setSocketConnectTimeout(geodeConfig.getSocketConnectTimeout())
          .setSubscriptionEnabled(geodeConfig.isSubscriptionEnabled())
          .setThreadLocalConnections(geodeConfig.isThreadLocalConnections());

      if (geodeConfig.isSubscriptionEnabled()) {
        poolFactory.setSubscriptionRedundancy(geodeConfig.getSubscriptionRedundancy());
      }

      poolFactory.create(geodeConfig.getPoolName());
    } catch (IllegalStateException e) {
      LOGGER.error(e.toString());
      LOGGER.error("Fail to create pool = " + geodeConfig.getPoolName());
      return false;
    } catch (Exception e) {
      LOGGER.error(e.toString());
      LOGGER.error("Fail to create pool = " + geodeConfig.getPoolName());
      return false;
    }

    LOGGER.debug("Create Client Cache successfully");

    return true;
  }

  public void clean() {
    available = false;

    try {
      if (null != clientCache) {
        clientCache.close(true);
      }
    } catch (CacheClosedException e) {
      LOGGER.error(e.toString());
    } catch (Exception e) {
      LOGGER.error(e.toString());
    }

    Map<String, Pool> pools = PoolManager.getAll();
    for (Pool pool : pools.values()) {
      String name = pool.getName();
      while (!pool.isDestroyed()) {
        try {
          pool.destroy();
        } catch (Exception e) {
          LOGGER.error(e.toString());
          LOGGER.error("Fail to destroy pool = " + name);
        }

        if (pool.isDestroyed()) {
          break;
        }

        try {
          LOGGER.warn("Sleep one second and try to destroy again");
          Thread.sleep(1000);
        } catch (Exception e) {
          LOGGER.error(e.toString());
        }
      }
    }

    reset();
  }

  private void reset() {
    available = false;
    config = null;
    clientCache = null;
    clientRegions.clear();
    state = ClientCacheServiceInitState.getInstance();
  }

  public int compare(GeodeConfig geodeConfig) {
    return config.compare(geodeConfig);
  }

  public boolean isAvailable() {
    return available;
  }

  public CacheTransactionManager getCacheTransactionManager() {
    return clientCache.getCacheTransactionManager();
  }

  public Region getRegion(String regionName) {
    return clientRegions.get(regionName).getRegion();
  }
}
