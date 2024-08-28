package ericsson.core.nrf.dbproxy.helper.subscription;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.clientcache.schema.SubscriptionData;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionFilterProto.SubscriptionFilter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionGetRequestProto.SubscriptionGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionGetResponseProto.SubscriptionDataList;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionGetResponseProto.SubscriptionGetResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionGetResponseProto.SubscriptionIDList;
import ericsson.core.nrf.dbproxy.helper.Helper;
import java.util.ArrayList;
import java.util.List;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class SubscriptionGetHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(SubscriptionGetHelper.class);

  private static SubscriptionGetHelper instance;

  private SubscriptionGetHelper() {
  }

  public static synchronized SubscriptionGetHelper getInstance() {
    if (null == instance) {
      instance = new SubscriptionGetHelper();
    }
    return instance;
  }

  public boolean validateFilter(SubscriptionFilter filter) {
    if (!filter.getIndex().getNoCond().isEmpty()) {
      return true;
    }

    if (!filter.getIndex().getNfStatusNotificationUri().isEmpty()) {
      return true;
    }

    if (!filter.getIndex().getNfInstanceId().isEmpty()) {
      return true;
    }

    if (!filter.getIndex().getNfType().isEmpty()) {
      return true;
    }

    if (filter.getIndex().getServiceNamesList().size() > 0) {
      return true;
    }

    if (filter.getIndex().getAmfCondsList().size() > 0) {
      return true;
    }

    if (filter.getIndex().getGuamiListList().size() > 0) {
      return true;
    }

    if (filter.getIndex().getSnssaiListList().size() > 0) {
      return true;
    }

    if (filter.getIndex().getStartValidityTime() < filter.getIndex().getEndValidityTime()) {
      return true;
    }

    return false;
  }

  public int validate(NFMessage message) {

    SubscriptionGetRequest request = message.getRequest().getGetRequest()
        .getSubscriptionGetRequest();
    switch (request.getDataCase()) {
      case SUBSCRIPTION_ID:
        if (request.getSubscriptionId().isEmpty()) {
          LOGGER.error("Subscription ID is empty in SubscriptionGetRequest");
          return Code.EMPTY_SUBSCRIPTION_ID;
        }
        break;
      case FILTER:
        if (!validateFilter(request.getFilter())) {
          LOGGER.error("Subscription Filter is empty in SubscriptionGetRequest");
          return Code.EMPTY_SUBSCRIPTION_FILTER;
        }
        break;
      default:
        LOGGER.error("Empty SubscriptionGetRequest is received");
        return Code.NFMESSAGE_PROTOCOL_ERROR;
    }

    return Code.VALID;
  }

  public NFMessage createResponse(int code) {

    SubscriptionGetResponse subscriptionGetResponse = SubscriptionGetResponse.newBuilder()
        .setCode(code).build();
    GetResponse getResponse = GetResponse.newBuilder()
        .setSubscriptionGetResponse(subscriptionGetResponse).build();
    return createNFMessage(getResponse);
  }

  @Override
  public NFMessage createResponse(ExecutionResult executionResult) {
    if (executionResult.getCode() == Code.SUCCESS) {
      SearchResult searchResult = (SearchResult) executionResult;

      List<ByteString> subscriptionDataList = new ArrayList<>();
      for (Object obj : searchResult.getItems()) {
        SubscriptionData item = (SubscriptionData) obj;
        subscriptionDataList.add(item.getData());
      }
      SubscriptionDataList data = SubscriptionDataList.newBuilder()
          .addAllSubscriptionData(subscriptionDataList).build();
      SubscriptionGetResponse subscriptionGetResponse = SubscriptionGetResponse.newBuilder()
          .setCode(searchResult.getCode()).setSubscriptionDataList(data).build();
      GetResponse getResponse = GetResponse.newBuilder()
          .setSubscriptionGetResponse(subscriptionGetResponse).build();
      return createNFMessage(getResponse);
    } else if (executionResult.getCode() == Code.SUBSCRIPTION_MONITOR_SUCCESS) {
      SearchResult searchResult = (SearchResult) executionResult;

      List<String> subscriptionIdList = new ArrayList<>();
      for (Object obj : searchResult.getItems()) {
        SubscriptionData item = (SubscriptionData) obj;
        subscriptionIdList.add(item.getSubscriptionID());
      }
      SubscriptionIDList data = SubscriptionIDList.newBuilder()
          .addAllSubscriptionId(subscriptionIdList).build();
      SubscriptionGetResponse subscriptionGetResponse = SubscriptionGetResponse.newBuilder()
          .setCode(Code.SUCCESS).setSubscriptionIdList(data).build();
      GetResponse getResponse = GetResponse.newBuilder()
          .setSubscriptionGetResponse(subscriptionGetResponse).build();
      return createNFMessage(getResponse);
    } else {
      return createResponse(executionResult.getCode());
    }
  }
}
