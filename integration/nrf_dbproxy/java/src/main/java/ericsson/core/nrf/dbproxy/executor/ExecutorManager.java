package ericsson.core.nrf.dbproxy.executor;

import ericsson.core.nrf.dbproxy.executor.cachenfprofile.CacheNFProfileGetExecutor;
import ericsson.core.nrf.dbproxy.executor.cachenfprofile.CacheNFProfilePutExecutor;
import ericsson.core.nrf.dbproxy.executor.gpsiprefixprofile.GpsiprefixProfileGetExecutor;
import ericsson.core.nrf.dbproxy.executor.gpsiprofile.GpsiProfileDeleteExecutor;
import ericsson.core.nrf.dbproxy.executor.gpsiprofile.GpsiProfileGetExecutor;
import ericsson.core.nrf.dbproxy.executor.gpsiprofile.GpsiProfilePutExecutor;
import ericsson.core.nrf.dbproxy.executor.groupprofile.GroupProfileDeleteExecutor;
import ericsson.core.nrf.dbproxy.executor.groupprofile.GroupProfileGetExecutor;
import ericsson.core.nrf.dbproxy.executor.groupprofile.GroupProfilePutExecutor;
import ericsson.core.nrf.dbproxy.executor.imsiprefixprofile.ImsiprefixProfileGetExecutor;
import ericsson.core.nrf.dbproxy.executor.nfprofile.NFProfileCountGetExecutor;
import ericsson.core.nrf.dbproxy.executor.nfprofile.NFProfileDeleteExecutor;
import ericsson.core.nrf.dbproxy.executor.nfprofile.NFProfileGetExecutor;
import ericsson.core.nrf.dbproxy.executor.nfprofile.NFProfilePutExecutor;
import ericsson.core.nrf.dbproxy.executor.nrfaddress.NRFAddressDeleteExecutor;
import ericsson.core.nrf.dbproxy.executor.nrfaddress.NRFAddressGetExecutor;
import ericsson.core.nrf.dbproxy.executor.nrfaddress.NRFAddressPutExecutor;
import ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfileDeleteExecutor;
import ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfileGetExecutor;
import ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfilePutExecutor;
import ericsson.core.nrf.dbproxy.executor.protocolerror.ProtocolErrorExecutor;
import ericsson.core.nrf.dbproxy.executor.subscription.SubscriptionDeleteExecutor;
import ericsson.core.nrf.dbproxy.executor.subscription.SubscriptionGetExecutor;
import ericsson.core.nrf.dbproxy.executor.subscription.SubscriptionPutExecutor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.DelRequest;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetRequest;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFRequest;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PutRequest;

public class ExecutorManager {

  private static ExecutorManager instance;

  static {
    instance = null;
  }

  private ExecutorManager() {
  }

  public static synchronized ExecutorManager getInstance() {
    if (null == instance) {
      instance = new ExecutorManager();
    }
    return instance;
  }

  public Executor getExecutor(NFMessage request) {
    if (!request.hasRequest()) {
      return ProtocolErrorExecutor.getInstance();
    }

    NFRequest nfRequest = request.getRequest();

    if (nfRequest.hasPutRequest()) {

      PutRequest putRequest = nfRequest.getPutRequest();

      if (putRequest.hasNfProfilePutRequest()) {
        return NFProfilePutExecutor.getInstance();
      } else if (putRequest.hasSubscriptionPutRequest()) {
        return SubscriptionPutExecutor.getInstance();
      } else if (putRequest.hasNrfAddressPutRequest()) {
        return NRFAddressPutExecutor.getInstance();
      } else if (putRequest.hasGroupProfilePutRequest()) {
        return GroupProfilePutExecutor.getInstance();
      } else if (putRequest.hasNrfProfilePutRequest()) {
        return NRFProfilePutExecutor.getInstance();
      } else if (putRequest.hasGpsiProfilePutRequest()) {
        return GpsiProfilePutExecutor.getInstance();
      } else if (putRequest.hasCacheNfProfilePutRequest()) {
        return CacheNFProfilePutExecutor.getInstance();
      } else {
        return ProtocolErrorExecutor.getInstance();
      }
    } else if (nfRequest.hasGetRequest()) {

      GetRequest getRequest = nfRequest.getGetRequest();

      if (getRequest.hasNfProfileGetRequest()) {
        return NFProfileGetExecutor.getInstance();
      } else if (getRequest.hasSubscriptionGetRequest()) {
        return SubscriptionGetExecutor.getInstance();
      } else if (getRequest.hasNrfAddressGetRequest()) {
        return NRFAddressGetExecutor.getInstance();
      } else if (getRequest.hasGroupProfileGetRequest()) {
        return GroupProfileGetExecutor.getInstance();
      } else if (getRequest.hasNrfProfileGetRequest()) {
        return NRFProfileGetExecutor.getInstance();
      } else if (getRequest.hasImsiprefixProfileGetRequest()) {
        return ImsiprefixProfileGetExecutor.getInstance();
      } else if (getRequest.hasGpsiProfileGetRequest()) {
        return GpsiProfileGetExecutor.getInstance();
      } else if (getRequest.hasGpsiprefixProfileGetRequest()) {
        return GpsiprefixProfileGetExecutor.getInstance();
      } else if (getRequest.hasNfProfileCountGetRequest()) {
        return NFProfileCountGetExecutor.getInstance();
      } else if (getRequest.hasCacheNfProfileGetRequest()) {
        return CacheNFProfileGetExecutor.getInstance();
      } else {
        return ProtocolErrorExecutor.getInstance();
      }
    } else if (nfRequest.hasDelRequest()) {

      DelRequest delRequest = nfRequest.getDelRequest();

      if (delRequest.hasNfProfileDelRequest()) {
        return NFProfileDeleteExecutor.getInstance();
      } else if (delRequest.hasSubscriptionDelRequest()) {
        return SubscriptionDeleteExecutor.getInstance();
      } else if (delRequest.hasNrfAddressDelRequest()) {
        return NRFAddressDeleteExecutor.getInstance();
      } else if (delRequest.hasGroupProfileDelRequest()) {
        return GroupProfileDeleteExecutor.getInstance();
      } else if (delRequest.hasNrfProfileDelRequest()) {
        return NRFProfileDeleteExecutor.getInstance();
      } else if (delRequest.hasGpsiProfileDelRequest()) {
        return GpsiProfileDeleteExecutor.getInstance();
      } else {
        return ProtocolErrorExecutor.getInstance();
      }
    } else {
      return ProtocolErrorExecutor.getInstance();
    }
  }
}
