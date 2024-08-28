package ericsson.core.nrf.dbproxy.helper.gpsiprofile;

import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import com.google.protobuf.ByteString;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfilePutRequestProto.GpsiProfilePutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfilePutResponseProto.GpsiProfilePutResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;

public class GpsiProfilePutHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(GpsiProfilePutHelper.class);

    private static GpsiProfilePutHelper instance;

    private GpsiProfilePutHelper() { }

    public static synchronized GpsiProfilePutHelper getInstance()
    {
        if(null == instance) {
            instance = new GpsiProfilePutHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        GpsiProfilePutRequest request = message.getRequest().getPutRequest().getGpsiProfilePutRequest();
        String gpsi_profile_id = request.getGpsiProfileId();
        if(gpsi_profile_id.isEmpty() == true) {
            logger.error("gpsi_profile_id field is empty in GpsiProfilePutRequest");
            return Code.EMPTY_GPSI_PROFILE_ID;
        }

        if(gpsi_profile_id.length() > Code.KEY_MAX_LENGTH) {
            logger.error("gpsi_profile_id length {} is too large, max length is {}",
                         gpsi_profile_id.length(), Code.KEY_MAX_LENGTH);
            return Code.GPSI_PROFILE_ID_LENGTH_EXCEED_MAX;
        }

        ByteString gpsi_profile_data = request.getGpsiProfileData();
        if(gpsi_profile_data.isEmpty() == true) {
            logger.error("gpsi_profile_data field is empty in GpsiProfilePutRequest");
            return Code.EMPTY_GPSI_PROFILE_DATA;
        }
		
		int profile_type = request.getIndex().getProfileType();
		if (profile_type != Code.PROFILE_TYPE_GROUPID && profile_type != Code.PROFILE_TYPE_INSTANCEID) {
            logger.error("profile_type is not PROFILE_TYPE_GROUPID or PROFILE_TYPE_INSTANCEID.");
            return Code.GPSI_PROFILE_INVALID_PROFILE_TYPE;
		}
		
        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {

        GpsiProfilePutResponse gpsi_profile_put_response = GpsiProfilePutResponse.newBuilder().setCode(code).build();
        PutResponse put_response = PutResponse.newBuilder().setGpsiProfilePutResponse(gpsi_profile_put_response).build();
        return createNFMessage(put_response);
    }
}
