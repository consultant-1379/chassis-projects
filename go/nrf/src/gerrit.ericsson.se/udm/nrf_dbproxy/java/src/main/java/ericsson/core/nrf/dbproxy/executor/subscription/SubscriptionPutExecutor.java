package ericsson.core.nrf.dbproxy.executor.subscription;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.clientcache.schema.KeyAggregation;
import ericsson.core.nrf.dbproxy.clientcache.schema.SubscriptionData;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionPutRequestProto.SubscriptionPutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionIndexProto.SubKeyStruct;
import ericsson.core.nrf.dbproxy.helper.subscription.SubscriptionPutHelper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import com.google.protobuf.ProtocolStringList;

public class SubscriptionPutExecutor extends Executor
{
    private static final Logger logger = LogManager.getLogger(SubscriptionPutExecutor.class);

    private static SubscriptionPutExecutor instance = null;

    private SubscriptionPutExecutor()
    {
        super(SubscriptionPutHelper.getInstance());
    }

    public static synchronized SubscriptionPutExecutor getInstance()
    {
        if(null == instance) {
            instance = new SubscriptionPutExecutor();
        }

        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        SubscriptionPutRequest put_request = request.getRequest().getPutRequest().getSubscriptionPutRequest();
        String subscription_id = put_request.getSubscriptionId();
        SubscriptionData subscription_data = createSubscriptionData(put_request);
        int code = ClientCacheService.getInstance().put(Code.SUBSCRIPTION_INDICE, subscription_id, subscription_data);
        return new ExecutionResult(code);
    }

    private SubscriptionData createSubscriptionData(SubscriptionPutRequest request)
    {
        SubscriptionData subscription_data = new SubscriptionData();
        subscription_data.setSubscriptionID(request.getSubscriptionId());
        subscription_data.setData(request.getSubscriptionData());
        subscription_data.setNoCond(request.getIndex().getNoCond());
        subscription_data.setNfStatusNotificationUri(request.getIndex().getNfStatusNotificationUri());
        subscription_data.setNfInstanceId(request.getIndex().getNfInstanceId());
        subscription_data.setNfType(request.getIndex().getNfType());
        subscription_data.setServiceName(request.getIndex().getServiceName());

        int id = 0;
        SubKeyStruct amfCond = request.getIndex().getAmfCond();
        if (amfCond != null) {

            KeyAggregation ka = new KeyAggregation();
            ka.setSubKey1(amfCond.getSubKey1());
            ka.setSubKey2(amfCond.getSubKey2());

            subscription_data.addAmfCond(id, ka);
        }

        id = 0;
        for(SubKeyStruct ks : request.getIndex().getGuamiListList()) {

            KeyAggregation ka = new KeyAggregation();
            ka.setSubKey1(ks.getSubKey1());
            ka.setSubKey2(ks.getSubKey2());
            ka.setSubKey3(ks.getSubKey3());
            ka.setSubKey4(ks.getSubKey4());
            ka.setSubKey5(ks.getSubKey5());

            subscription_data.addGuamiList(id, ka);
            id++;
        }

        id = 0;
        for(SubKeyStruct ks : request.getIndex().getSnssaiListList()) {

            KeyAggregation ka = new KeyAggregation();
            ka.setSubKey1(ks.getSubKey1());
            ka.setSubKey2(ks.getSubKey2());
            ka.setSubKey3(ks.getSubKey3());
            ka.setSubKey4(ks.getSubKey4());
            ka.setSubKey5(ks.getSubKey5());

            subscription_data.addSnssaiList(id, ka);
            id++;
        }

        for(String nsi : request.getIndex().getNsiListList()) {
            subscription_data.addNsiList(nsi);
        } 

        id = 0;
        SubKeyStruct nfGroupCond = request.getIndex().getNfGroupCond();
        if (nfGroupCond != null) {

            KeyAggregation ka = new KeyAggregation();
            ka.setSubKey1(nfGroupCond.getSubKey1());
            ka.setSubKey2(nfGroupCond.getSubKey2());

            subscription_data.addNfGroupCond(id, ka);
        }

        subscription_data.setValidityTime(request.getIndex().getValidityTime());

        logger.trace("Subscription Data : {} ", subscription_data.toString());

        return subscription_data;
    }

}
