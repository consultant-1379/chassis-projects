package ericsson.core.nrf.dbproxy.helper.subscription;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.DelResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelRequestProto.SubscriptionDelRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionDelResponseProto.SubscriptionDelResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class SubscriptionDelHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(SubscriptionDelHelper.class);

  private static SubscriptionDelHelper instance;

  private SubscriptionDelHelper() {
  }

  public static synchronized SubscriptionDelHelper getInstance() {
    if (null == instance) {
      instance = new SubscriptionDelHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    SubscriptionDelRequest request = message.getRequest().getDelRequest()
        .getSubscriptionDelRequest();
    String subscriptionId = request.getSubscriptionId();
    if (subscriptionId.isEmpty()) {
      LOGGER.error("subscription_id field is empty in SubscriptionDelRequest");
      return Code.EMPTY_SUBSCRIPTION_ID;
    }

    return Code.VALID;
  }

  public NFMessage createResponse(int code) {

    SubscriptionDelResponse subscriptionDelResponse = SubscriptionDelResponse.newBuilder()
        .setCode(code).build();
    DelResponse delResponse = DelResponse.newBuilder()
        .setSubscriptionDelResponse(subscriptionDelResponse).build();
    return createNFMessage(delResponse);
  }
}
