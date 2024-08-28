package ericsson.core.nrf.dbproxy.helper.subscription;

import java.util.List;
import java.util.ArrayList;

import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import com.google.protobuf.ByteString;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.clientcache.schema.SubscriptionData;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionGetRequestProto.SubscriptionGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionFilterProto.SubscriptionFilter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionGetResponseProto.SubscriptionGetResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionGetResponseProto.SubscriptionDataList;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionGetResponseProto.SubscriptionIDList;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;

public class SubscriptionGetHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(SubscriptionGetHelper.class);

    private static SubscriptionGetHelper instance;

    private SubscriptionGetHelper() { }

    public static synchronized SubscriptionGetHelper getInstance()
    {
        if(null == instance) {
            instance = new SubscriptionGetHelper();
        }
        return instance;
    }

    public boolean validateFilter(SubscriptionFilter filter)
    {
        if(!filter.getIndex().getNoCond().isEmpty()) {
            return true;
        }

        if(!filter.getIndex().getNfStatusNotificationUri().isEmpty()) {
            return true;
        }

        if(!filter.getIndex().getNfInstanceId().isEmpty()) {
            return true;
        }

        if(!filter.getIndex().getNfType().isEmpty()) {
            return true;
        }

        if(filter.getIndex().getServiceNamesList().size() > 0) {
            return true;
        }

        if(filter.getIndex().getAmfCondsList().size() > 0) {
            return true;
        }

        if(filter.getIndex().getGuamiListList().size() > 0) {
            return true;
        }

        if(filter.getIndex().getSnssaiListList().size() > 0) {
            return true;
        }

        if(filter.getIndex().getStartValidityTime() < filter.getIndex().getEndValidityTime()) {
            return true;
        }

        return false;
    }

    public int validate(NFMessage message)
    {

        SubscriptionGetRequest request = message.getRequest().getGetRequest().getSubscriptionGetRequest();
        switch(request.getDataCase()) {
        case SUBSCRIPTION_ID:
            if(request.getSubscriptionId().isEmpty()) {
                logger.error("Subscription ID is empty in SubscriptionGetRequest");
                return Code.EMPTY_SUBSCRIPTION_ID;
            }
            break;
        case FILTER:
            if(!validateFilter(request.getFilter())) {
                logger.error("Subscription Filter is empty in SubscriptionGetRequest");
                return Code.EMPTY_SUBSCRIPTION_FILTER;
            }
            break;
        default:
            logger.error("Empty SubscriptionGetRequest is received");
            return Code.NFMESSAGE_PROTOCOL_ERROR;
        }

        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {

        SubscriptionGetResponse subscription_get_response = SubscriptionGetResponse.newBuilder().setCode(code).build();
        GetResponse get_response = GetResponse.newBuilder().setSubscriptionGetResponse(subscription_get_response).build();
        return createNFMessage(get_response);
    }

    public NFMessage createResponse(ExecutionResult execution_result)
    {
        if(execution_result.getCode() == Code.SUCCESS) {
            SearchResult search_result = (SearchResult)execution_result;

            List<ByteString> subscription_data_list = new ArrayList<>();
            for(Object obj : search_result.getItems()) {
                SubscriptionData item = (SubscriptionData)obj;
                subscription_data_list.add(item.getData());
            }
            SubscriptionDataList data = SubscriptionDataList.newBuilder().addAllSubscriptionData(subscription_data_list).build();
            SubscriptionGetResponse subscription_get_response = SubscriptionGetResponse.newBuilder().setCode(search_result.getCode()).setSubscriptionDataList(data).build();
            GetResponse get_response = GetResponse.newBuilder().setSubscriptionGetResponse(subscription_get_response).build();
            return createNFMessage(get_response);
        } else if (execution_result.getCode() == Code.SUBSCRIPTION_MONITOR_SUCCESS) {
            SearchResult search_result = (SearchResult)execution_result;

            List<String> subscription_id_list = new ArrayList<>();
            for(Object obj : search_result.getItems()) {
                SubscriptionData item = (SubscriptionData)obj;
                subscription_id_list.add(item.getSubscriptionID());
            }
            SubscriptionIDList data = SubscriptionIDList.newBuilder().addAllSubscriptionId(subscription_id_list).build();
            SubscriptionGetResponse subscription_get_response = SubscriptionGetResponse.newBuilder().setCode(Code.SUCCESS).setSubscriptionIdList(data).build();
            GetResponse get_response = GetResponse.newBuilder().setSubscriptionGetResponse(subscription_get_response).build();
            return createNFMessage(get_response);
        } else {
            return createResponse(execution_result.getCode());
        }
    }
}
