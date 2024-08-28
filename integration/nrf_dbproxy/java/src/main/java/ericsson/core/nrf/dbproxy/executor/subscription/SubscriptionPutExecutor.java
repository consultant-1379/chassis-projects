package ericsson.core.nrf.dbproxy.executor.subscription;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.clientcache.schema.KeyAggregation;
import ericsson.core.nrf.dbproxy.clientcache.schema.SubscriptionData;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionIndexProto.SubKeyStruct;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionPutRequestProto.SubscriptionPutRequest;
import ericsson.core.nrf.dbproxy.helper.subscription.SubscriptionPutHelper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class SubscriptionPutExecutor extends Executor {

  private static final Logger LOGGER = LogManager.getLogger(SubscriptionPutExecutor.class);

  private static SubscriptionPutExecutor instance;

  static {
    instance = null;
  }

  private SubscriptionPutExecutor() {
    super(SubscriptionPutHelper.getInstance());
  }

  public static synchronized SubscriptionPutExecutor getInstance() {
    if (null == instance) {
      instance = new SubscriptionPutExecutor();
    }

    return instance;
  }

  protected ExecutionResult execute(NFMessage request) {
    SubscriptionPutRequest putRequest = request.getRequest().getPutRequest()
        .getSubscriptionPutRequest();
    String subscriptionId = putRequest.getSubscriptionId();
    SubscriptionData subscriptionData = createSubscriptionData(putRequest);
    int code = ClientCacheService.getInstance()
        .put(Code.SUBSCRIPTION_INDICE, subscriptionId, subscriptionData);
    return new ExecutionResult(code);
  }

  private SubscriptionData createSubscriptionData(SubscriptionPutRequest request) {
    SubscriptionData subscriptionData = new SubscriptionData();
    subscriptionData.setSubscriptionID(request.getSubscriptionId());
    subscriptionData.setData(request.getSubscriptionData());
    subscriptionData.setNoCond(request.getIndex().getNoCond());
    subscriptionData.setNfStatusNotificationUri(request.getIndex().getNfStatusNotificationUri());
    subscriptionData.setNfInstanceId(request.getIndex().getNfInstanceId());
    subscriptionData.setNfType(request.getIndex().getNfType());
    subscriptionData.setServiceName(request.getIndex().getServiceName());

    int id = 0;
    SubKeyStruct amfCond = request.getIndex().getAmfCond();
    if (amfCond != null) {

      KeyAggregation ka = new KeyAggregation();
      ka.setSubKey1(amfCond.getSubKey1());
      ka.setSubKey2(amfCond.getSubKey2());

      subscriptionData.addAmfCond(id, ka);
    }

    id = 0;
    for (SubKeyStruct ks : request.getIndex().getGuamiListList()) {

      KeyAggregation ka = new KeyAggregation();
      ka.setSubKey1(ks.getSubKey1());
      ka.setSubKey2(ks.getSubKey2());
      ka.setSubKey3(ks.getSubKey3());
      ka.setSubKey4(ks.getSubKey4());
      ka.setSubKey5(ks.getSubKey5());

      subscriptionData.addGuamiList(id, ka);
      id++;
    }

    id = 0;
    for (SubKeyStruct ks : request.getIndex().getSnssaiListList()) {

      KeyAggregation ka = new KeyAggregation();
      ka.setSubKey1(ks.getSubKey1());
      ka.setSubKey2(ks.getSubKey2());
      ka.setSubKey3(ks.getSubKey3());
      ka.setSubKey4(ks.getSubKey4());
      ka.setSubKey5(ks.getSubKey5());

      subscriptionData.addSnssaiList(id, ka);
      id++;
    }

    for (String nsi : request.getIndex().getNsiListList()) {
      subscriptionData.addNsiList(nsi);
    }

    id = 0;
    SubKeyStruct nfGroupCond = request.getIndex().getNfGroupCond();
    if (nfGroupCond != null) {

      KeyAggregation ka = new KeyAggregation();
      ka.setSubKey1(nfGroupCond.getSubKey1());
      ka.setSubKey2(nfGroupCond.getSubKey2());

      subscriptionData.addNfGroupCond(id, ka);
    }

    subscriptionData.setValidityTime(request.getIndex().getValidityTime());

    LOGGER.debug("Subscription Data : {} ", subscriptionData.toString());

    return subscriptionData;
  }

}
