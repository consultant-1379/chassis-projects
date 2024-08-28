package ericsson.core.nrf.dbproxy.executor.nrfaddress;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.clientcache.schema.NRFAddress;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressIndexProto.NRFAddressIndex;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressPutRequestProto.NRFAddressPutRequest;
import ericsson.core.nrf.dbproxy.helper.nrfaddress.NRFAddressPutHelper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NRFAddressPutExecutor extends Executor {

  private static final Logger LOGGER = LogManager.getLogger(NRFAddressPutExecutor.class);

  private static NRFAddressPutExecutor instance;

  static {
    instance = null;
  }

  private NRFAddressPutExecutor() {
    super(NRFAddressPutHelper.getInstance());
  }

  public static synchronized NRFAddressPutExecutor getInstance() {
    if (null == instance) {
      instance = new NRFAddressPutExecutor();
    }

    return instance;
  }

  protected ExecutionResult execute(NFMessage request) {
    NRFAddressPutRequest putRequest = request.getRequest().getPutRequest()
        .getNrfAddressPutRequest();
    String nrfAddressId = putRequest.getNrfAddressId();
    NRFAddress nrfAddress = createNRFAddress(putRequest);
    int code = ClientCacheService.getInstance()
        .put(Code.NRFADDRESS_INDICE, nrfAddressId, nrfAddress);
    return new ExecutionResult(code);
  }

  private NRFAddress createNRFAddress(NRFAddressPutRequest request) {
    NRFAddressIndex index = request.getIndex();

    NRFAddress nrfAddress = new NRFAddress();

    nrfAddress.setNRFAddressID(request.getNrfAddressId());
    nrfAddress.setData(request.getNrfAddressData());
    nrfAddress.setKey1(index.getNrfAddressKey1());
    nrfAddress.setKey2(index.getNrfAddressKey2());
    nrfAddress.setKey3(index.getNrfAddressKey3());
    nrfAddress.setKey4(index.getNrfAddressKey4());
    nrfAddress.setKey5(index.getNrfAddressKey5());

    LOGGER.debug("NRF Address : {}", nrfAddress.toString());

    return nrfAddress;
  }

}
