package ericsson.core.nrf.dbproxy.helper.subscription;

import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import com.google.protobuf.ByteString;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionPutRequestProto.SubscriptionPutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionPutResponseProto.SubscriptionPutResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;

public class SubscriptionPutHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(SubscriptionPutHelper.class);

    private static SubscriptionPutHelper instance;

    private SubscriptionPutHelper() { }

    public static synchronized SubscriptionPutHelper getInstance()
    {
        if(null == instance) {
            instance = new SubscriptionPutHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        SubscriptionPutRequest request = message.getRequest().getPutRequest().getSubscriptionPutRequest();
        String subscription_id = request.getSubscriptionId();
        if(subscription_id.isEmpty() == true) {
            logger.error("subscription_id field is empty in SubscriptionPutRequest");
            return Code.EMPTY_SUBSCRIPTION_ID;
        }

        if(subscription_id.length() > Code.KEY_MAX_LENGTH) {
            logger.error("subscription_id length {} is too large, max length is {}",
                         subscription_id.length(), Code.KEY_MAX_LENGTH);
            return Code.SUBSCRIPTION_ID_LENGTH_EXCEED_MAX;
        }

        ByteString subscription_data = request.getSubscriptionData();
        if(subscription_data == ByteString.EMPTY) {
            logger.error("subscription_data is empty in SubscriptionPutRequest");
            return Code.EMPTY_SUBSCRIPTION_DATA;
        }

        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {

        SubscriptionPutResponse subscription_put_response = SubscriptionPutResponse.newBuilder().setCode(code).build();
        PutResponse put_response = PutResponse.newBuilder().setSubscriptionPutResponse(subscription_put_response).build();
        return createNFMessage(put_response);
    }
}
