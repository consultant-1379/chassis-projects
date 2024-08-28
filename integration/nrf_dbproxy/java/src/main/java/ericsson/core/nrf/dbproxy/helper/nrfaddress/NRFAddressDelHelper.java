package ericsson.core.nrf.dbproxy.helper.nrfaddress;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.DelResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressDelRequestProto.NRFAddressDelRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressDelResponseProto.NRFAddressDelResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NRFAddressDelHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(NRFAddressDelHelper.class);

  private static NRFAddressDelHelper instance;

  private NRFAddressDelHelper() {
  }

  public static synchronized NRFAddressDelHelper getInstance() {
    if (null == instance) {
      instance = new NRFAddressDelHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    NRFAddressDelRequest request = message.getRequest().getDelRequest().getNrfAddressDelRequest();
    String nrfAddressId = request.getNrfAddressId();
    if (nrfAddressId.isEmpty()) {
      LOGGER.error("nrf_address_id field is empty in NRFAddressDelRequest");
      return Code.EMPTY_NRF_ADDRESS_ID;
    }

    return Code.VALID;
  }

  public NFMessage createResponse(int code) {

    NRFAddressDelResponse nrfAddressDelResponse = NRFAddressDelResponse.newBuilder()
        .setCode(code).build();
    DelResponse delResponse = DelResponse.newBuilder()
        .setNrfAddressDelResponse(nrfAddressDelResponse).build();
    return createNFMessage(delResponse);
  }
}
