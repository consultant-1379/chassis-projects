package ericsson.core.nrf.dbproxy.executor.subscription;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelRequestProto.SubscriptionDelRequest;
import ericsson.core.nrf.dbproxy.helper.subscription.SubscriptionDelHelper;

public class SubscriptionDeleteExecutor extends Executor {

  private static SubscriptionDeleteExecutor instance;

  static {
    instance = null;
  }

  private SubscriptionDeleteExecutor() {
    super(SubscriptionDelHelper.getInstance());
  }

  public static synchronized SubscriptionDeleteExecutor getInstance() {
    if (null == instance) {
      instance = new SubscriptionDeleteExecutor();
    }

    return instance;
  }

  protected ExecutionResult execute(NFMessage request) {
    SubscriptionDelRequest delRequest = request.getRequest().getDelRequest()
        .getSubscriptionDelRequest();
    String subscriptionId = delRequest.getSubscriptionId();
    int code = ClientCacheService.getInstance().delete(Code.SUBSCRIPTION_INDICE, subscriptionId);
    return new ExecutionResult(code);
  }
}
