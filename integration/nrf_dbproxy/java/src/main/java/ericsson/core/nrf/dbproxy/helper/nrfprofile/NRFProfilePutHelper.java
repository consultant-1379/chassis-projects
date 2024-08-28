package ericsson.core.nrf.dbproxy.helper.nrfprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PutResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfilePutRequestProto.NRFProfilePutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfilePutResponseProto.NRFProfilePutResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NRFProfilePutHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(NRFProfilePutHelper.class);

  private static NRFProfilePutHelper instance;

  private NRFProfilePutHelper() {
  }

  public static synchronized NRFProfilePutHelper getInstance() {
    if (null == instance) {
      instance = new NRFProfilePutHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    NRFProfilePutRequest request = message.getRequest().getPutRequest().getNrfProfilePutRequest();
    String nrfInstanceId = request.getNrfInstanceId();
    if (nrfInstanceId.isEmpty()) {
      LOGGER.error("nrf_instance_id field is empty in NRFProfilePutRequest");
      return Code.EMPTY_NRF_INSTANCE_ID;
    }

    if (nrfInstanceId.length() > Code.KEY_MAX_LENGTH) {
      LOGGER.error("nrf_instance_id length {} is too large, max length is {}",
          nrfInstanceId.length(), Code.KEY_MAX_LENGTH);
      return Code.NRF_INSTANCE_ID_LENGTH_EXCEED_MAX;
    }

    ByteString rawNrfProfile = request.getRawNrfProfile();
    if (rawNrfProfile.isEmpty()) {
      LOGGER.error("rawNrfProfile field is empty in NRFProfilePutRequest");
      return Code.EMPTY_RAW_NRF_PROFILE;
    }

    return Code.VALID;
  }

  public NFMessage createResponse(int code) {

    NRFProfilePutResponse nrfProfilePutResponse = NRFProfilePutResponse.newBuilder()
        .setCode(code).build();
    PutResponse putResponse = PutResponse.newBuilder()
        .setNrfProfilePutResponse(nrfProfilePutResponse).build();
    return createNFMessage(putResponse);
  }
}
