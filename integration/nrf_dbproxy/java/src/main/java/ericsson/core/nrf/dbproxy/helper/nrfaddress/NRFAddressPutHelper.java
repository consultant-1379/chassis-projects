package ericsson.core.nrf.dbproxy.helper.nrfaddress;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PutResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressPutRequestProto.NRFAddressPutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressPutResponseProto.NRFAddressPutResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NRFAddressPutHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(NRFAddressPutHelper.class);

  private static NRFAddressPutHelper instance;

  private NRFAddressPutHelper() {
  }

  public static synchronized NRFAddressPutHelper getInstance() {
    if (null == instance) {
      instance = new NRFAddressPutHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    NRFAddressPutRequest request = message.getRequest().getPutRequest().getNrfAddressPutRequest();
    String nrfAddressId = request.getNrfAddressId();
    if (nrfAddressId.isEmpty()) {
      LOGGER.error("nrfAddressId field is empty in NRFAddressPutRequest");
      return Code.EMPTY_NRF_ADDRESS_ID;
    }

    if (nrfAddressId.length() > Code.KEY_MAX_LENGTH) {
      LOGGER.error("nrfAddressId length {} is too large, max length is {}",
          nrfAddressId.length(), Code.KEY_MAX_LENGTH);
      return Code.NRF_ADDRESS_ID_LENGTH_EXCEED_MAX;
    }

    ByteString nrfAddressData = request.getNrfAddressData();
    if (nrfAddressData.isEmpty()) {
      LOGGER.error("nrfAddressData field is empty in NRFAddressPutRequest");
      return Code.EMPTY_NRF_ADDRESS_DATA;
    }

    return Code.VALID;
  }

  public NFMessage createResponse(int code) {

    NRFAddressPutResponse nrfAddressPutResponse = NRFAddressPutResponse.newBuilder()
        .setCode(code).build();
    PutResponse putResponse = PutResponse.newBuilder()
        .setNrfAddressPutResponse(nrfAddressPutResponse).build();
    return createNFMessage(putResponse);
  }
}
