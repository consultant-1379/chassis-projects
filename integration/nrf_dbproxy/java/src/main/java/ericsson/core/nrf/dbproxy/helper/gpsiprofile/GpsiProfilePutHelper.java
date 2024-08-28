package ericsson.core.nrf.dbproxy.helper.gpsiprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PutResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfilePutRequestProto.GpsiProfilePutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfilePutResponseProto.GpsiProfilePutResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class GpsiProfilePutHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(GpsiProfilePutHelper.class);

  private static GpsiProfilePutHelper instance;

  private GpsiProfilePutHelper() {
  }

  public static synchronized GpsiProfilePutHelper getInstance() {
    if (null == instance) {
      instance = new GpsiProfilePutHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    GpsiProfilePutRequest request = message.getRequest().getPutRequest().getGpsiProfilePutRequest();
    String gpsiProfileId = request.getGpsiProfileId();
    if (gpsiProfileId.isEmpty()) {
      LOGGER.error("gpsi_profile_id field is empty in GpsiProfilePutRequest");
      return Code.EMPTY_GPSI_PROFILE_ID;
    }

    if (gpsiProfileId.length() > Code.KEY_MAX_LENGTH) {
      LOGGER.error("gpsi_profile_id length {} is too large, max length is {}",
          gpsiProfileId.length(), Code.KEY_MAX_LENGTH);
      return Code.GPSI_PROFILE_ID_LENGTH_EXCEED_MAX;
    }

    ByteString gpsiProfileData = request.getGpsiProfileData();
    if (gpsiProfileData.isEmpty()) {
      LOGGER.error("gpsiProfileData field is empty in GpsiProfilePutRequest");
      return Code.EMPTY_GPSI_PROFILE_DATA;
    }

    int profileType = request.getIndex().getProfileType();
    if (profileType != Code.PROFILE_TYPE_GROUPID && profileType != Code.PROFILE_TYPE_INSTANCEID) {
      LOGGER.error("profileType is not PROFILE_TYPE_GROUPID or PROFILE_TYPE_INSTANCEID.");
      return Code.GPSI_PROFILE_INVALID_PROFILE_TYPE;
    }

    return Code.VALID;
  }

  public NFMessage createResponse(int code) {

    GpsiProfilePutResponse gpsiProfilePutResponse = GpsiProfilePutResponse.newBuilder()
        .setCode(code).build();
    PutResponse putResponse = PutResponse.newBuilder()
        .setGpsiProfilePutResponse(gpsiProfilePutResponse).build();
    return createNFMessage(putResponse);
  }
}
