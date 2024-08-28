package ericsson.core.nrf.dbproxy.clientcache.state;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.config.GeodeConfig;

public class ClientCacheServiceCreatedState implements ClientCacheServiceState
{
    private static ClientCacheServiceCreatedState instance = null;

    private ClientCacheServiceCreatedState() {}

    public static synchronized ClientCacheServiceCreatedState getInstance()
    {
        if(instance == null) {
            instance = new ClientCacheServiceCreatedState();
        }
        return instance;
    }

    public void apply(GeodeConfig geode_config)
    {
        if(ClientCacheService.getInstance().compare(geode_config) != Code.NOT_CHANGED) {
            ClientCacheService.getInstance().clean();
            ClientCacheService.getInstance().apply(geode_config);
        }
    }
}
