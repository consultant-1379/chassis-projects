package ericsson.core.nrf.dbproxy.clientcache;

import java.util.Map;
import java.util.HashMap;
import java.util.List;
import java.util.ArrayList;
import org.apache.geode.pdx.ReflectionBasedAutoSerializer;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.geode.cache.CacheXmlException;
import org.apache.geode.cache.TimeoutException;
import org.apache.geode.cache.CacheWriterException;
import org.apache.geode.cache.RegionExistsException;
import org.apache.geode.security.AuthenticationFailedException;
import org.apache.geode.security.AuthenticationRequiredException;
import org.apache.geode.cache.CacheClosedException;
import org.apache.geode.cache.client.*;
import org.apache.geode.cache.CacheTransactionManager;
import org.apache.geode.cache.Region;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.FragmentSessionManagement;
import ericsson.core.nrf.dbproxy.config.GeodeConfig;
import ericsson.core.nrf.dbproxy.clientcache.region.*;
import ericsson.core.nrf.dbproxy.clientcache.state.ClientCacheServiceInitState;
import ericsson.core.nrf.dbproxy.clientcache.state.ClientCacheServiceState;

public class ClientCacheService
{
    private static final Logger logger = LogManager.getLogger(ClientCacheService.class);

    private static ClientCacheService instance = null;

    private ClientCache client_cache =  null;

    private Map<String,ClientRegion> client_regions;

    private ClientCacheServiceState state;
    private boolean available;
    private GeodeConfig config;

    private ClientCacheService()
    {
        this.client_regions = new HashMap<String,ClientRegion>();
        reset();
    }

    public static synchronized ClientCacheService getInstance()
    {
        if(null == instance) {
            instance = new ClientCacheService();
        }
        return instance;
    }

    public void setState(ClientCacheServiceState service_state)
    {
        state = service_state;
    }

    public void apply(GeodeConfig geode_config)
    {
        logger.debug("Apply geode configuration on client cache service");
        state.apply(geode_config);
    }

    public int put(String region_name, Object region_key, Object region_item)
    {
        ClientRegion client_region = client_regions.get(region_name);
        return client_region.put(region_key, region_item);
    }

    public int delete(String region_name, Object region_key)
    {
        ClientRegion client_region = client_regions.get(region_name);
        return client_region.delete(region_key,true);
    }

    public int delete(String region_name, Object region_key, boolean withTransaction)
    {
        ClientRegion client_region = client_regions.get(region_name);
        return client_region.delete(region_key,withTransaction);
    }

    public ExecutionResult getByID(String region_name, Object region_key)
    {
        ClientRegion client_region = client_regions.get(region_name);
        return client_region.getByID(region_key);
    }
    public ExecutionResult getAllByID(String region_name, List<Long> id_list)
    {
        ClientRegion client_region = client_regions.get(region_name);
        return client_region.getAllByID(id_list);
    }
    public ExecutionResult getByFilter(String region_name, String query_string)
    {
        ClientRegion client_region = client_regions.get(region_name);
        return client_region.getByFilter(region_name, query_string);
    }
    
    public ExecutionResult getByCountFilter(String region_name, String query_string)
    {
        ClientRegion client_region = client_regions.get(region_name);
        return client_region.getByCountFilter(query_string);
    }

    public ExecutionResult getByFragSessionId(String region_name, String fragment_session_id)
    {
        return FragmentSessionManagement.getInstance().get(region_name, fragment_session_id);
    }

    public boolean initialize(GeodeConfig geode_config)
    {
        if(createClientCache(geode_config) == false)
            return false;

        String[] regions = geode_config.getRegionList();
        for (int i=0; i<regions.length; i++) {
            ClientRegion region = new ClientRegion(regions[i]);
            if(region.initialize(client_cache, geode_config.getPoolName()) == false)
                return false;
            client_regions.put(regions[i], region);
        }

        config = geode_config;
        available = true;
        return true;
    }

    private boolean createClientCache(GeodeConfig geode_config)
    {
        ClientCacheFactory factory = new ClientCacheFactory();
        try {
            client_cache = factory.setPdxSerializer(new ReflectionBasedAutoSerializer("ericsson.core.nrf.dbproxy.clientcache.schema.*")).create();
        } catch (CacheXmlException | TimeoutException | RegionExistsException | CacheWriterException | IllegalStateException | AuthenticationFailedException | AuthenticationRequiredException e ) {
            logger.error(e.toString());
            return false;
        }

        PoolFactory pool_factory  = PoolManager.createFactory();
        List<String> kvdb_locator_ip_list = geode_config.getLocatorIPList();
        int kvdb_locator_port = geode_config.getLocatorPort();
        for(String locator_ip : kvdb_locator_ip_list) {
            try {
                pool_factory = pool_factory.addLocator(locator_ip, kvdb_locator_port);
            } catch (IllegalArgumentException | IllegalStateException e) {
                logger.error(e.toString());
                logger.error("Fail to add locator = " + locator_ip + ":" + Integer.toString(kvdb_locator_port));
                return false;
            } catch (Exception e) {
                logger.error(e.toString());
                logger.error("Fail to add locator = " + locator_ip + ":" + Integer.toString(kvdb_locator_port));
                return false;
            }
        }

        try {
            pool_factory.setFreeConnectionTimeout(geode_config.getFreeConnectionTimeout())
            .setIdleTimeout(geode_config.getIdleTimeout())
            .setLoadConditioningInterval(geode_config.getLoadConditioningInterval())
            .setMaxConnections(geode_config.getMaxConnections())
            .setMinConnections(geode_config.getMinConnections())
            .setPingInterval(geode_config.getPingInterval())
            .setPRSingleHopEnabled(geode_config.isPrSingleHopEnabled())
            .setReadTimeout(geode_config.getReadTimeout())
            .setRetryAttempts(geode_config.getRetryAttempts())
            .setSocketBufferSize(geode_config.getSocketBufferSize())
            .setSocketConnectTimeout(geode_config.getSocketConnectTimeout())
            .setSubscriptionEnabled(geode_config.isSubscriptionEnabled())
            .setThreadLocalConnections(geode_config.isThreadLocalConnections())
            .create(geode_config.getPoolName());
        } catch(IllegalStateException e) {
            logger.error(e.toString());
            logger.error("Fail to create pool = " + geode_config.getPoolName());
            return false;
        } catch(Exception e) {
            logger.error(e.toString());
            logger.error("Fail to create pool = " + geode_config.getPoolName());
            return false;
        }

        logger.debug("Create Client Cache successfully");

        return true;
    }

    public void clean()
    {
        available = false;

        try {
            if(null != client_cache)
                client_cache.close(true);
        } catch(CacheClosedException e) {
            logger.error(e.toString());
        } catch(Exception e) {
            logger.error(e.toString());
        }

        Map<String,Pool> pools = PoolManager.getAll();
        for(Pool pool : pools.values()) {
            String name = pool.getName();
            while(pool.isDestroyed() == false) {
                try {
                    pool.destroy();
                } catch(Exception e) {
                    logger.error(e.toString());
                    logger.error("Fail to destroy pool = " + name);
                }

                if(pool.isDestroyed()) break;

                try {
                    logger.warn("Sleep one second and try to destroy again");
                    Thread.sleep(1000);
                } catch(Exception e) {
                    logger.error(e.toString());
                }
            }
        }

        reset();
    }

    private void reset()
    {
        available = false;
        config = null;
        client_cache = null;
        client_regions.clear();
        state = ClientCacheServiceInitState.getInstance();
    }

    public int compare(GeodeConfig geode_config)
    {
        return config.compare(geode_config);
    }

    public boolean isAvailable()
    {
        return available;
    }

    public CacheTransactionManager getCacheTransactionManager()
    {
        return client_cache.getCacheTransactionManager();
    }

    public Region getRegion(String region_name)
    {
        return client_regions.get(region_name).getRegion();
    }
}
