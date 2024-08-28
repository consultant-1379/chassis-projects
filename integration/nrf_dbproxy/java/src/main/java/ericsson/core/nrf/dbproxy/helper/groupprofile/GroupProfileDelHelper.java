package ericsson.core.nrf.dbproxy.helper.groupprofile;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.DelResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileDelRequestProto.GroupProfileDelRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileDelResponseProto.GroupProfileDelResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class GroupProfileDelHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(GroupProfileDelHelper.class);

  private static GroupProfileDelHelper instance;

  private GroupProfileDelHelper() {
  }

  public static synchronized GroupProfileDelHelper getInstance() {
    if (null == instance) {
      instance = new GroupProfileDelHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    GroupProfileDelRequest request = message.getRequest().getDelRequest()
        .getGroupProfileDelRequest();
    String groupProfileId = request.getGroupProfileId();
    if (groupProfileId.isEmpty()) {
      LOGGER.error("group_profile_id field is empty in GroupProfileDelRequest");
      return Code.EMPTY_GROUP_PROFILE_ID;
    }

    return Code.VALID;
  }

  public NFMessage createResponse(int code) {

    GroupProfileDelResponse groupProfileDelResponse = GroupProfileDelResponse.newBuilder()
        .setCode(code).build();
    DelResponse delResponse = DelResponse.newBuilder()
        .setGroupProfileDelResponse(groupProfileDelResponse).build();
    return createNFMessage(delResponse);
  }
}
