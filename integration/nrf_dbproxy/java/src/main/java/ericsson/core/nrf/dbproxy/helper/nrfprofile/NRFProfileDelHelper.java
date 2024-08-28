package ericsson.core.nrf.dbproxy.helper.nrfprofile;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.DelResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileDelRequestProto.NRFProfileDelRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileDelResponseProto.NRFProfileDelResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NRFProfileDelHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(NRFProfileDelHelper.class);

  private static NRFProfileDelHelper instance;

  private NRFProfileDelHelper() {
  }

  public static synchronized NRFProfileDelHelper getInstance() {
    if (null == instance) {
      instance = new NRFProfileDelHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    NRFProfileDelRequest request = message.getRequest().getDelRequest().getNrfProfileDelRequest();
    String nrfInstanceId = request.getNrfInstanceId();
    if (nrfInstanceId.isEmpty()) {
      LOGGER.error("nrfInstanceId field is empty in NFProfileDelRequest");
      return Code.EMPTY_NRF_INSTANCE_ID;
    }

    return Code.VALID;
  }

  public NFMessage createResponse(int code) {
    NRFProfileDelResponse nrfProfileDelResponse = NRFProfileDelResponse.newBuilder()
        .setCode(code).build();
    DelResponse delResponse = DelResponse.newBuilder()
        .setNrfProfileDelResponse(nrfProfileDelResponse).build();
    return createNFMessage(delResponse);
  }
}
