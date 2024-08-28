package ericsson.core.nrf.dbproxy.helper.groupprofile;

import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileDelRequestProto.GroupProfileDelRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileDelResponseProto.GroupProfileDelResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;

public class GroupProfileDelHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(GroupProfileDelHelper.class);

    private static GroupProfileDelHelper instance;

    private GroupProfileDelHelper() { }

    public static synchronized GroupProfileDelHelper getInstance()
    {
        if(null == instance) {
            instance = new GroupProfileDelHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        GroupProfileDelRequest request = message.getRequest().getDelRequest().getGroupProfileDelRequest();
        String group_profile_id = request.getGroupProfileId();
        if(group_profile_id.isEmpty() == true) {
            logger.error("group_profile_id field is empty in GroupProfileDelRequest");
            return Code.EMPTY_GROUP_PROFILE_ID;
        }

        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {

        GroupProfileDelResponse group_profile_del_response = GroupProfileDelResponse.newBuilder().setCode(code).build();
        DelResponse del_response = DelResponse.newBuilder().setGroupProfileDelResponse(group_profile_del_response).build();
        return createNFMessage(del_response);
    }
}
