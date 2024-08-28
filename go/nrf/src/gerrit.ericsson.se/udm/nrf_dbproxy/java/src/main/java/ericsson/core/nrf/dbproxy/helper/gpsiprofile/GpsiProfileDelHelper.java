package ericsson.core.nrf.dbproxy.helper.gpsiprofile;

import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileDelRequestProto.GpsiProfileDelRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileDelResponseProto.GpsiProfileDelResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;

public class GpsiProfileDelHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(GpsiProfileDelHelper.class);

    private static GpsiProfileDelHelper instance;

    private GpsiProfileDelHelper() { }

    public static synchronized GpsiProfileDelHelper getInstance()
    {
        if(null == instance) {
            instance = new GpsiProfileDelHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        GpsiProfileDelRequest request = message.getRequest().getDelRequest().getGpsiProfileDelRequest();
        String gpsi_profile_id = request.getGpsiProfileId();
        if(gpsi_profile_id.isEmpty() == true) {
            logger.error("gpsi_profile_id field is empty in GpsiProfileDelRequest");
            return Code.EMPTY_GPSI_PROFILE_ID;
        }

        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {

        GpsiProfileDelResponse gpsi_profile_del_response = GpsiProfileDelResponse.newBuilder().setCode(code).build();
        DelResponse del_response = DelResponse.newBuilder().setGpsiProfileDelResponse(gpsi_profile_del_response).build();
        return createNFMessage(del_response);
    }
}
