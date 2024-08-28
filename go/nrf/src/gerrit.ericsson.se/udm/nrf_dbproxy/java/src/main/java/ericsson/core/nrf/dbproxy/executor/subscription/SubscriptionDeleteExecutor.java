package ericsson.core.nrf.dbproxy.executor.subscription;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelRequestProto.SubscriptionDelRequest;
import ericsson.core.nrf.dbproxy.helper.subscription.SubscriptionDelHelper;

public class SubscriptionDeleteExecutor extends Executor
{
    private static SubscriptionDeleteExecutor instance = null;

    private SubscriptionDeleteExecutor()
    {
        super(SubscriptionDelHelper.getInstance());
    }

    public static synchronized SubscriptionDeleteExecutor getInstance()
    {
        if(null == instance) {
            instance = new SubscriptionDeleteExecutor();
        }

        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        SubscriptionDelRequest del_request = request.getRequest().getDelRequest().getSubscriptionDelRequest();
        String subscription_id = del_request.getSubscriptionId();
        int code = ClientCacheService.getInstance().delete(Code.SUBSCRIPTION_INDICE, subscription_id);
        return new ExecutionResult(code);
    }
}
