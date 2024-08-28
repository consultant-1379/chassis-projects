package ericsson.core.nrf.dbproxy.clientcache.state;

import ericsson.core.nrf.dbproxy.config.GeodeConfig;

public interface ClientCacheServiceState
{

    public void apply(GeodeConfig geode_config);
}
