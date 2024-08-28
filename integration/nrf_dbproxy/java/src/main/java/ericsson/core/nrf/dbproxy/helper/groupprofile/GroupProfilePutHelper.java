package ericsson.core.nrf.dbproxy.helper.groupprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PutResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfilePutRequestProto.GroupProfilePutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfilePutResponseProto.GroupProfilePutResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class GroupProfilePutHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(GroupProfilePutHelper.class);

  private static GroupProfilePutHelper instance;

  private GroupProfilePutHelper() {
  }

  public static synchronized GroupProfilePutHelper getInstance() {
    if (null == instance) {
      instance = new GroupProfilePutHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    GroupProfilePutRequest request = message.getRequest().getPutRequest()
        .getGroupProfilePutRequest();
    String groupProfileId = request.getGroupProfileId();
    if (groupProfileId.isEmpty()) {
      LOGGER.error("groupProfileId field is empty in GroupProfilePutRequest");
      return Code.EMPTY_GROUP_PROFILE_ID;
    }

    if (groupProfileId.length() > Code.KEY_MAX_LENGTH) {
      LOGGER.error("groupProfileId length {} is too large, max length is {}",
          groupProfileId.length(), Code.KEY_MAX_LENGTH);
      return Code.GROUP_PROFILE_ID_LENGTH_EXCEED_MAX;
    }

    ByteString groupProfileData = request.getGroupProfileData();
    if (groupProfileData.isEmpty()) {
      LOGGER.error("groupProfileData field is empty in GroupProfilePutRequest");
      return Code.EMPTY_GROUP_PROFILE_DATA;
    }

    int profileType = request.getIndex().getProfileType();
    if (profileType != Code.PROFILE_TYPE_GROUPID && profileType != Code.PROFILE_TYPE_INSTANCEID) {
      LOGGER.error("profileType is not PROFILE_TYPE_GROUPID or PROFILE_TYPE_INSTANCEID.");
      return Code.GROUP_PROFILE_INVALID_PROFILE_TYPE;
    }
    return Code.VALID;
  }

  public NFMessage createResponse(int code) {

    GroupProfilePutResponse groupProfilePutResponse = GroupProfilePutResponse.newBuilder()
        .setCode(code).build();
    PutResponse putResponse = PutResponse.newBuilder()
        .setGroupProfilePutResponse(groupProfilePutResponse).build();
    return createNFMessage(putResponse);
  }
}
