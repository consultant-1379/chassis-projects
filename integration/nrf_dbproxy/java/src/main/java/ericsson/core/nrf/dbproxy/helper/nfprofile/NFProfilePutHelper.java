package ericsson.core.nrf.dbproxy.helper.nfprofile;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PutResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfilePutRequestProto.NFProfilePutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfilePutResponseProto.NFProfilePutResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NFProfilePutHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(NFProfilePutHelper.class);

  private static NFProfilePutHelper instance;

  private NFProfilePutHelper() {
  }

  public static synchronized NFProfilePutHelper getInstance() {
    if (null == instance) {
      instance = new NFProfilePutHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    NFProfilePutRequest request = message.getRequest().getPutRequest().getNfProfilePutRequest();
    String nfInstanceId = request.getNfInstanceId();
    if (nfInstanceId.isEmpty()) {
      LOGGER.error("nf_instance_id field is empty in NFProfilePutRequest");
      return Code.EMPTY_NF_INSTANCE_ID;
    }

    if (nfInstanceId.length() > Code.KEY_MAX_LENGTH) {
      LOGGER.error("nf_instance_id length {} is too large, max length is {}",
          nfInstanceId.length(), Code.KEY_MAX_LENGTH);
      return Code.NF_INSTANCE_ID_LENGTH_EXCEED_MAX;
    }

    String nfProfile = request.getNfProfile();
    if (nfProfile.isEmpty()) {
      LOGGER.error("nfProfile field is empty in NFProfilePutRequest");
      return Code.EMPTY_NF_PROFILE;
    }

    return Code.VALID;
  }

  public NFMessage createResponse(int code) {

    NFProfilePutResponse nfProfilePutResponse = NFProfilePutResponse.newBuilder().setCode(code)
        .build();
    PutResponse putResponse = PutResponse.newBuilder()
        .setNfProfilePutResponse(nfProfilePutResponse).build();
    return createNFMessage(putResponse);
  }
}
