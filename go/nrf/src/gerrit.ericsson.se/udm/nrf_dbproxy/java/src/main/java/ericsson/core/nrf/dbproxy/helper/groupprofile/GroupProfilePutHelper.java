package ericsson.core.nrf.dbproxy.helper.groupprofile;

import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import com.google.protobuf.ByteString;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfilePutRequestProto.GroupProfilePutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfilePutResponseProto.GroupProfilePutResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;

public class GroupProfilePutHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(GroupProfilePutHelper.class);

    private static GroupProfilePutHelper instance;

    private GroupProfilePutHelper() { }

    public static synchronized GroupProfilePutHelper getInstance()
    {
        if(null == instance) {
            instance = new GroupProfilePutHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        GroupProfilePutRequest request = message.getRequest().getPutRequest().getGroupProfilePutRequest();
        String group_profile_id = request.getGroupProfileId();
        if(group_profile_id.isEmpty() == true) {
            logger.error("group_profile_id field is empty in GroupProfilePutRequest");
            return Code.EMPTY_GROUP_PROFILE_ID;
        }

        if(group_profile_id.length() > Code.KEY_MAX_LENGTH) {
            logger.error("group_profile_id length {} is too large, max length is {}",
                         group_profile_id.length(), Code.KEY_MAX_LENGTH);
            return Code.GROUP_PROFILE_ID_LENGTH_EXCEED_MAX;
        }

        ByteString group_profile_data = request.getGroupProfileData();
        if(group_profile_data.isEmpty() == true) {
            logger.error("group_profile_data field is empty in GroupProfilePutRequest");
            return Code.EMPTY_GROUP_PROFILE_DATA;
        }
		
		int profile_type = request.getIndex().getProfileType();
		if (profile_type != Code.PROFILE_TYPE_GROUPID && profile_type != Code.PROFILE_TYPE_INSTANCEID) {
            logger.error("profile_type is not PROFILE_TYPE_GROUPID or PROFILE_TYPE_INSTANCEID.");
            return Code.GROUP_PROFILE_INVALID_PROFILE_TYPE;
		}
        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {

        GroupProfilePutResponse group_profile_put_response = GroupProfilePutResponse.newBuilder().setCode(code).build();
        PutResponse put_response = PutResponse.newBuilder().setGroupProfilePutResponse(group_profile_put_response).build();
        return createNFMessage(put_response);
    }
}
