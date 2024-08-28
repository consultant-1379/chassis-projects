package ericsson.core.nrf.dbproxy.clientcache.region;

import java.util.List;

import ericsson.core.nrf.dbproxy.common.*;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.geode.cache.Region;
import org.apache.geode.cache.client.*;
import org.apache.geode.cache.RegionExistsException;
import org.apache.geode.cache.CacheClosedException;
import org.apache.geode.distributed.LeaseExpiredException;
import org.apache.geode.cache.TimeoutException;
import org.apache.geode.cache.CacheWriterException;
import org.apache.geode.cache.CacheLoaderException;
import org.apache.geode.cache.query.SelectResults;
import org.apache.geode.cache.query.FunctionDomainException;
import org.apache.geode.cache.query.TypeMismatchException;
import org.apache.geode.cache.query.NameResolutionException;
import org.apache.geode.cache.query.QueryInvocationTargetException;
import org.apache.geode.cache.LowMemoryException;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;

public class ClientRegion
{
    private static final Logger logger = LogManager.getLogger(ClientRegion.class);

    private static final String LOG_1 = "Fail to put key = ";
    private static final String LOG_2 = "in region ";
    private static final String LOG_3 = "Fail to delete the region item by key = ";
    protected String region_name;
    protected Region region;

    public ClientRegion(String region_name)
    {
        this.region_name = region_name;
    }

    public boolean initialize(ClientCache client_cache, String pool_name)
    {
        try {
            int core_number = Runtime.getRuntime().availableProcessors();
            region = client_cache.createClientRegionFactory(ClientRegionShortcut.PROXY).setConcurrencyLevel(core_number).setPoolName(pool_name).create(region_name);
        } catch (RegionExistsException | CacheClosedException e) {
            logger.error(e.toString());
            logger.error("Fail to create region = " + region_name);
            return false;
        } catch (Exception e) {
            logger.error(e.toString());
            logger.error("Fail to create region = " + region_name);
            return false;
        }

        logger.debug("Create Region = " + region_name + " successfully");

        return true;

    }

    public String getRegionName()
    {
        return region_name;
    }

    public Region getRegion() {
	return region;
    }

    public int put(Object region_key, Object region_item)
    {

        int code = Code.CREATED;
        try {
            region.put(region_key, region_item);
            logger.debug("Key = " + region_key + " is put successfully in region " + region_name);
        } catch (NullPointerException | ClassCastException e) {
            logger.error(e.toString());
            logger.error(LOG_1 + region_key + LOG_2 + region_name);
            code = Code.INVALID_NF_INSTANCE_ID_NOT_SERIALIZABLE;
        } catch (LeaseExpiredException | TimeoutException |
                     CacheWriterException | LowMemoryException e) {
            logger.error(e.toString());
            logger.error(LOG_1 + region_key + LOG_2 + region_name);
            code = Code.INTERNAL_ERROR;
        } catch (Exception e) {
            logger.error(e.toString());
            logger.error(LOG_1 + region_key + LOG_2 + region_name);
            code = Code.INTERNAL_ERROR;
        }

        return code;

    }

    public int delete(Object region_key, boolean withTransaction)
    {
        int code = Code.SUCCESS;

        try {
            if (withTransaction == true) {
                ClientCacheService.getInstance().getCacheTransactionManager().begin();

                if(region.containsValueForKey(region_key) == false) {
                    code = Code.DATA_NOT_EXIST;
                    logger.debug("Region Item is not found by key = " +  region_key + LOG_2 + region_name);
                } else {
                    region.remove(region_key);
                    logger.debug("Region Item is deleted successfully by key = " +  region_key + LOG_2 + region_name);
                }
                ClientCacheService.getInstance().getCacheTransactionManager().commit();
            } else {
                if(region.containsValueForKey(region_key) == false) {
                    code = Code.DATA_NOT_EXIST;
                    logger.debug("Region Item is not found by key = " +  region_key + LOG_2 + region_name);
                } else {
                    region.remove(region_key);
                    logger.debug("Region Item is deleted successfully by key = " +  region_key + LOG_2 + region_name);
                }
            }
        } catch (NullPointerException | IllegalArgumentException e) {
            logger.error(e.toString());
            logger.debug(LOG_3 +  region_key + LOG_2 + region_name);
            code = Code.INVALID_NF_INSTANCE_ID_NOT_SERIALIZABLE;
        } catch (LeaseExpiredException | TimeoutException | CacheWriterException e) {
            logger.error(e.toString());
            logger.debug(LOG_3 +  region_key + LOG_2 + region_name);
            code = Code.INTERNAL_ERROR;
        } catch (Exception e) {
            logger.error(e.toString());
            logger.debug(LOG_3 +  region_key + LOG_2 + region_name);
            code = Code.INTERNAL_ERROR;
        }

        return code;
    }

    public ExecutionResult getByID(Object region_key)
    {
        int code = Code.SUCCESS;
        try {
            Object object = region.get(region_key);
            if(null == object) {
                code = Code.DATA_NOT_EXIST;
            } else {
                SearchResult search_result = new SearchResult(false);
                search_result.add(object);
                return search_result;
            }
        } catch (NullPointerException | IllegalArgumentException e) {
            logger.error(e.toString());
            code = Code.INVALID_NF_INSTANCE_ID_NOT_SERIALIZABLE;
        } catch (LeaseExpiredException | TimeoutException | CacheLoaderException e) {
            logger.error(e.toString());
            code = Code.INTERNAL_ERROR;
        } catch (Exception e) {
            logger.error(e.toString());
            code = Code.INTERNAL_ERROR;
        }

        return new ExecutionResult(code);
    }

    public ExecutionResult getAllByID(List<Long> id_list)
    {
        int code = Code.SUCCESS;
        try {
            Object object = region.getAll(id_list);
            if(null == object) {
                code = Code.DATA_NOT_EXIST;
            } else {
                SearchResult search_result = new SearchResult(false);
                search_result.add(object);
                return search_result;
            }
        } catch (NullPointerException | IllegalArgumentException e) {
            logger.error(e.toString());
            code = Code.INVALID_NF_INSTANCE_ID_NOT_SERIALIZABLE;
        } catch (LeaseExpiredException | TimeoutException | CacheLoaderException e) {
            logger.error(e.toString());
            code = Code.INTERNAL_ERROR;
        } catch (Exception e) {
            logger.error(e.toString());
            code = Code.INTERNAL_ERROR;
        }

        return new ExecutionResult(code);
    }

    public ExecutionResult getByFilter(String region_name, String query_string)
    {
	logger.trace("OQL = " + query_string);
        int code = Code.SUCCESS;
        try {
            SelectResults<Object> results = (SelectResults<Object>) region.getRegionService().getQueryService().newQuery(query_string).execute();
            if (results.size() <= 0) {
                code = Code.DATA_NOT_EXIST;
            } else if(FragmentUtil.isNeedFragment(region_name, results)){
                logger.trace("result byte bigger than 3MB, need to transmit fragment");
                FragmentResult fragmentResult = new FragmentResult();
                for(Object object : results) {
                    fragmentResult.add(object);
                }
                return fragmentResult;
            } else {
                logger.trace("result byte smaller than 3MB, can transmit one time");
                SearchResult searchResult = new SearchResult(false);
                for(Object object : results) {
                    searchResult.add(object);
                }
                return searchResult;
            }
        } catch (FunctionDomainException | TypeMismatchException | NameResolutionException | QueryInvocationTargetException e) {
            logger.error(e.toString());
            code = Code.INTERNAL_ERROR;
        } catch (Exception e) {
            logger.error(e.toString());
            code = Code.INTERNAL_ERROR;
        }

        return new ExecutionResult(code);
    }
    
    public ExecutionResult getByCountFilter(String query_string)
    {
	logger.trace("OQL = " + query_string);
        int code = Code.SUCCESS;
        try {
        	SelectResults results_test = (SelectResults)region.getRegionService().getQueryService().newQuery(query_string).execute();
        	List<Integer> results = results_test.asList();
            if (results.size() <= 0) {
                code = Code.DATA_NOT_EXIST;
            } else {
                logger.trace("The count value is " + results.get(0));
                SearchResult searchResult = new SearchResult(false);
                for(Object object : results) {
                    searchResult.add(object);
                }
                return searchResult;
            }
        } catch (FunctionDomainException | TypeMismatchException | NameResolutionException | QueryInvocationTargetException e) {
            logger.error(e.toString());
            code = Code.INTERNAL_ERROR;
        } catch (Exception e) {
            logger.error(e.toString());
            code = Code.INTERNAL_ERROR;
        }

        return new ExecutionResult(code);
    }
}
