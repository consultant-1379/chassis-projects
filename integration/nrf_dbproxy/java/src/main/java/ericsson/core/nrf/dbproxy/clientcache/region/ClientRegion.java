package ericsson.core.nrf.dbproxy.clientcache.region;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.FragmentResult;
import ericsson.core.nrf.dbproxy.common.FragmentUtil;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import java.util.List;
import java.util.Map;
import org.apache.geode.cache.CacheClosedException;
import org.apache.geode.cache.CacheLoaderException;
import org.apache.geode.cache.CacheTransactionManager;
import org.apache.geode.cache.CacheWriterException;
import org.apache.geode.cache.CommitConflictException;
import org.apache.geode.cache.InterestResultPolicy;
import org.apache.geode.cache.LowMemoryException;
import org.apache.geode.cache.Region;
import org.apache.geode.cache.RegionExistsException;
import org.apache.geode.cache.TimeoutException;
import org.apache.geode.cache.client.ClientCache;
import org.apache.geode.cache.client.ClientRegionShortcut;
import org.apache.geode.cache.query.FunctionDomainException;
import org.apache.geode.cache.query.NameResolutionException;
import org.apache.geode.cache.query.QueryInvocationTargetException;
import org.apache.geode.cache.query.SelectResults;
import org.apache.geode.cache.query.TypeMismatchException;
import org.apache.geode.distributed.LeaseExpiredException;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class ClientRegion {

  private static final Logger LOGGER = LogManager.getLogger(ClientRegion.class);

  private static final String LOG_1 = "Fail to put key = ";
  private static final String LOG_2 = " in region ";
  protected String region_name;
  protected Region region;

  public ClientRegion(String regionName) {
    this.region_name = regionName;
  }

  public boolean initialize(ClientCache clientCache, String poolName,
      boolean isSubscriptionEnabled) {
    try {
      int coreNumber = Runtime.getRuntime().availableProcessors();
      if (region_name.equals(Code.NFPROFILE_INDICE) && isSubscriptionEnabled) {
        region = clientCache.createClientRegionFactory(ClientRegionShortcut.CACHING_PROXY)
            .setConcurrencyLevel(coreNumber).setPoolName(poolName).create(region_name);
        region.registerInterestForAllKeys(InterestResultPolicy.KEYS_VALUES);
      } else {
        region = clientCache.createClientRegionFactory(ClientRegionShortcut.PROXY)
            .setConcurrencyLevel(coreNumber).setPoolName(poolName).create(region_name);
      }
    } catch (RegionExistsException | CacheClosedException e) {
      LOGGER.error(e.toString());
      LOGGER.error("Fail to create region = " + region_name);
      return false;
    } catch (Exception e) {
      LOGGER.error(e.toString());
      LOGGER.error("Fail to create region = " + region_name);
      return false;
    }

    LOGGER.debug("Create Region = " + region_name + " successfully");

    return true;

  }

  public String getRegionName() {
    return region_name;
  }

  public Region getRegion() {
    return region;
  }

  public int put(Object regionKey, Object regionItem) {

    int code = Code.CREATED;
    try {
      region.put(regionKey, regionItem);
      LOGGER.debug("Key = " + regionKey + " is put successfully in region " + region_name);
    } catch (NullPointerException | ClassCastException e) {
      LOGGER.error(e.toString());
      LOGGER.error(LOG_1 + regionKey + LOG_2 + region_name);
      code = Code.INVALID_NF_INSTANCE_ID_NOT_SERIALIZABLE;
    } catch (LeaseExpiredException | TimeoutException |
        CacheWriterException | LowMemoryException e) {
      LOGGER.error(e.toString());
      LOGGER.error(LOG_1 + regionKey + LOG_2 + region_name);
      code = Code.INTERNAL_ERROR;
    } catch (Exception e) {
      LOGGER.error(e.toString());
      LOGGER.error(LOG_1 + regionKey + LOG_2 + region_name);
      code = Code.INTERNAL_ERROR;
    }

    return code;

  }

  //putAll function is used transaction to put some key/value into region
  public int putAll(Map<Object, Object> kvMap) {
    int code = Code.INTERNAL_ERROR;
    CacheTransactionManager txManager = ClientCacheService.getInstance()
        .getCacheTransactionManager();
    boolean retryTransaction = false;
    int retryTime = 0;
    do {
      try {
        retryTime++;
        if (retryTime > 3) {
          break;
        }
        txManager.begin();
        region.putAll(kvMap);
        txManager.commit();
        retryTransaction = false;
        code = Code.CREATED;
      } catch (CommitConflictException conflictException) {
        LOGGER.error(conflictException.toString());
        LOGGER.error(
            "Transaction commit failed, put keys " + kvMap.keySet() + " to region " + region_name
                + " failed");
        retryTransaction = true;
      } catch (Exception e) {
        LOGGER.error(e.toString());
      } finally {
        if (txManager.exists()) {
          txManager.rollback();
        }
      }
    } while (retryTransaction);
    return code;
  }

  public int delete(Object regionKey, boolean withTransaction) {
    int code = Code.INTERNAL_ERROR;
    if (withTransaction) {
      CacheTransactionManager txManager = ClientCacheService.getInstance()
          .getCacheTransactionManager();
      boolean retryTransaction = false;
      int retryTime = 0;
      do {
        try {
          retryTime++;
          if (retryTime > 3) {
            break;
          }
          if (!region.containsKeyOnServer(regionKey)) {
            code = Code.DATA_NOT_EXIST;
            LOGGER.debug("Region Item is not found by key = " + regionKey + LOG_2 + region_name);
            break;
          }
          txManager.begin();
          region.remove(regionKey);
          LOGGER.debug(
              "Region Item is deleted successfully by key = " + regionKey + LOG_2 + region_name);
          txManager.commit();
          retryTransaction = false;
          code = Code.SUCCESS;
        } catch (CommitConflictException e) {
          LOGGER.error(e.toString());
          LOGGER.error(
              "Transaction commit failed, delete key " + regionKey + " to region " + region_name
                  + " failed");
          retryTransaction = true;
        } catch (Exception e) {
          LOGGER.error(e.toString());
        } finally {
          if (txManager.exists()) {
            txManager.rollback();
          }
        }
      } while (retryTransaction);


    } else {
      try {
        if (!region.containsKeyOnServer(regionKey)) {
          code = Code.DATA_NOT_EXIST;
          LOGGER.debug("Region Item is not found by key = " + regionKey + LOG_2 + region_name);
        } else {
          region.remove(regionKey);
          LOGGER.debug(
              "Region Item is deleted successfully by key = " + regionKey + LOG_2 + region_name);
          code = Code.SUCCESS;
        }
      } catch (Exception e) {
        LOGGER.error(e.toString());
      }
    }
    return code;
  }


  public int deleteAll(List<String> keyList) {
    int code = Code.INTERNAL_ERROR;
    CacheTransactionManager txManager = ClientCacheService.getInstance()
        .getCacheTransactionManager();
    boolean retryTransaction = false;
    int retryTime = 0;
    do {
      try {
        retryTime++;
        if (retryTime > 3) {
          break;
        }
        txManager.begin();
        region.removeAll(keyList);
        LOGGER.debug(
            "Region Item is deleted successfully by keys = " + keyList + LOG_2 + region_name);
        txManager.commit();
        retryTransaction = false;
        code = Code.SUCCESS;
      } catch (CommitConflictException e) {
        LOGGER.error(e.toString());
        LOGGER.error(
            "Transaction commit failed, delete keys " + keyList + " to region " + region_name
                + " failed");
        retryTransaction = true;
      } catch (Exception e) {
        LOGGER.error(e.toString());
      } finally {
        if (txManager.exists()) {
          txManager.rollback();
        }
      }
    } while (retryTransaction);
    return code;
  }

  public ExecutionResult getByID(Object regionKey) {
    int code = Code.SUCCESS;
    try {
      Object object = region.get(regionKey);
      if (null == object) {
        code = Code.DATA_NOT_EXIST;
      } else {
        SearchResult searchResult = new SearchResult(false);
        searchResult.add(object);
        return searchResult;
      }
    } catch (NullPointerException | IllegalArgumentException e) {
      LOGGER.error(e.toString());
      code = Code.INVALID_NF_INSTANCE_ID_NOT_SERIALIZABLE;
    } catch (LeaseExpiredException | TimeoutException | CacheLoaderException e) {
      LOGGER.error(e.toString());
      code = Code.INTERNAL_ERROR;
    } catch (Exception e) {
      LOGGER.error(e.toString());
      code = Code.INTERNAL_ERROR;
    }

    return new ExecutionResult(code);
  }

  public ExecutionResult getAllByID(List<Long> idList) {
    int code = Code.SUCCESS;
    try {
      Object object = region.getAll(idList);
      if (null == object) {
        code = Code.DATA_NOT_EXIST;
      } else {
        SearchResult searchResult = new SearchResult(false);
        searchResult.add(object);
        return searchResult;
      }
    } catch (NullPointerException | IllegalArgumentException e) {
      LOGGER.error(e.toString());
      code = Code.INVALID_NF_INSTANCE_ID_NOT_SERIALIZABLE;
    } catch (LeaseExpiredException | TimeoutException | CacheLoaderException e) {
      LOGGER.error(e.toString());
      code = Code.INTERNAL_ERROR;
    } catch (Exception e) {
      LOGGER.error(e.toString());
      code = Code.INTERNAL_ERROR;
    }

    return new ExecutionResult(code);
  }

  public ExecutionResult getByFilter(String regionName, String queryString) {
    LOGGER.debug("OQL = " + queryString);
    int code = Code.SUCCESS;
    try {
      SelectResults<Object> results = (SelectResults<Object>) region.getRegionService()
          .getQueryService().newQuery(queryString).execute();
      if (results.size() <= 0) {
        code = Code.DATA_NOT_EXIST;
      } else if (FragmentUtil.isNeedFragment(regionName, results)) {
        LOGGER.debug("result byte bigger than 3MB, need to transmit fragment");
        FragmentResult fragmentResult = new FragmentResult();
        for (Object object : results) {
          fragmentResult.add(object);
        }
        return fragmentResult;
      } else {
        LOGGER.debug("result byte smaller than 3MB, can transmit one time");
        SearchResult searchResult = new SearchResult(false);
        for (Object object : results) {
          searchResult.add(object);
        }
        return searchResult;
      }
    } catch (FunctionDomainException | TypeMismatchException | NameResolutionException | QueryInvocationTargetException e) {
      LOGGER.error(e.toString());
      code = Code.INTERNAL_ERROR;
    } catch (Exception e) {
      LOGGER.error(e.toString());
      code = Code.INTERNAL_ERROR;
    }

    return new ExecutionResult(code);
  }

  public ExecutionResult getByCountFilter(String queryString) {
    LOGGER.debug("OQL = " + queryString);
    int code = Code.SUCCESS;
    try {
      SelectResults resultsTest = (SelectResults) region.getRegionService().getQueryService()
          .newQuery(queryString).execute();
      List<Integer> results = resultsTest.asList();
      if (results.size() <= 0) {
        code = Code.DATA_NOT_EXIST;
      } else {
        LOGGER.debug("The count value is " + results.get(0));
        SearchResult searchResult = new SearchResult(false);
        for (Object object : results) {
          searchResult.add(object);
        }
        return searchResult;
      }
    } catch (FunctionDomainException | TypeMismatchException | NameResolutionException | QueryInvocationTargetException e) {
      LOGGER.error(e.toString());
      code = Code.INTERNAL_ERROR;
    } catch (Exception e) {
      LOGGER.error(e.toString());
      code = Code.INTERNAL_ERROR;
    }

    return new ExecutionResult(code);
  }
}
