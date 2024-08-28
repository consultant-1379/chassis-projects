package ericsson.core.nrf.dbproxy.helper.subscription;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PutResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionPutRequestProto.SubscriptionPutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionPutResponseProto.SubscriptionPutResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class SubscriptionPutHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(SubscriptionPutHelper.class);

  private static SubscriptionPutHelper instance;

  private SubscriptionPutHelper() {
  }

  public static synchronized SubscriptionPutHelper getInstance() {
    if (null == instance) {
      instance = new SubscriptionPutHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    SubscriptionPutRequest request = message.getRequest().getPutRequest()
        .getSubscriptionPutRequest();
    String subscriptionId = request.getSubscriptionId();
    if (subscriptionId.isEmpty()) {
      LOGGER.error("subscriptionId field is empty in SubscriptionPutRequest");
      return Code.EMPTY_SUBSCRIPTION_ID;
    }

    if (subscriptionId.length() > Code.KEY_MAX_LENGTH) {
      LOGGER.error("subscriptionId length {} is too large, max length is {}",
          subscriptionId.length(), Code.KEY_MAX_LENGTH);
      return Code.SUBSCRIPTION_ID_LENGTH_EXCEED_MAX;
    }

    ByteString subscriptionData = request.getSubscriptionData();
    if (subscriptionData == ByteString.EMPTY) {
      LOGGER.error("subscription_data is empty in SubscriptionPutRequest");
      return Code.EMPTY_SUBSCRIPTION_DATA;
    }

    return Code.VALID;
  }

  public NFMessage createResponse(int code) {

    SubscriptionPutResponse subscriptionPutResponse = SubscriptionPutResponse.newBuilder()
        .setCode(code).build();
    PutResponse putResponse = PutResponse.newBuilder()
        .setSubscriptionPutResponse(subscriptionPutResponse).build();
    return createNFMessage(putResponse);
  }
}
