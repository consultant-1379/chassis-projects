package ericsson.core.nrf.dbproxy.clientcache.state;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.config.GeodeConfig;

public class ClientCacheServiceInitState implements ClientCacheServiceState {

  private static ClientCacheServiceInitState instance;

  static {
    instance = null;
  }

  private ClientCacheServiceInitState() {
  }

  public static synchronized ClientCacheServiceInitState getInstance() {
    if (instance == null) {
      instance = new ClientCacheServiceInitState();
    }
    return instance;
  }

  public void apply(GeodeConfig geodeConfig) {
    if (ClientCacheService.getInstance().initialize(geodeConfig)) {
      ClientCacheService.getInstance().setState(ClientCacheServiceCreatedState.getInstance());
    } else {
      ClientCacheService.getInstance().clean();
    }
  }
}
