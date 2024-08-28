package ericsson.core.nrf.dbproxy.helper.nfprofile;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.helper.Helper;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfilePutRequestProto.NFProfilePutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfilePutResponseProto.NFProfilePutResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;

public class NFProfilePutHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(NFProfilePutHelper.class);

    private static NFProfilePutHelper instance;

    private NFProfilePutHelper() { }

    public static synchronized NFProfilePutHelper getInstance()
    {
        if(null == instance) {
            instance = new NFProfilePutHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        NFProfilePutRequest request = message.getRequest().getPutRequest().getNfProfilePutRequest();
        String nf_instance_id = request.getNfInstanceId();
        if(nf_instance_id.isEmpty() == true) {
            logger.error("nf_instance_id field is empty in NFProfilePutRequest");
            return Code.EMPTY_NF_INSTANCE_ID;
        }

        if(nf_instance_id.length() > Code.KEY_MAX_LENGTH) {
            logger.error("nf_instance_id length {} is too large, max length is {}",
                         nf_instance_id.length(), Code.KEY_MAX_LENGTH);
            return Code.NF_INSTANCE_ID_LENGTH_EXCEED_MAX;
        }

        String nf_profile = request.getNfProfile();
        if(nf_profile.isEmpty() == true) {
            logger.error("nf_profile field is empty in NFProfilePutRequest");
            return Code.EMPTY_NF_PROFILE;
        }

        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {

        NFProfilePutResponse nf_profile_put_response = NFProfilePutResponse.newBuilder().setCode(code).build();
        PutResponse put_response = PutResponse.newBuilder().setNfProfilePutResponse(nf_profile_put_response).build();
        return createNFMessage(put_response);
    }
}
