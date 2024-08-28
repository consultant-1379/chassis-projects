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

public class NRFProfilePutHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(NRFProfilePutHelper.class);

    private static NRFProfilePutHelper instance;

    private NRFProfilePutHelper()
    {
    }

    public static synchronized NRFProfilePutHelper getInstance()
    {
        if (null == instance) {
            instance = new NRFProfilePutHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        NRFProfilePutRequest request = message.getRequest().getPutRequest().getNrfProfilePutRequest();
        String nrf_instance_id = request.getNrfInstanceId();
        if (nrf_instance_id.isEmpty() == true) {
            logger.error("nrf_instance_id field is empty in NRFProfilePutRequest");
            return Code.EMPTY_NRF_INSTANCE_ID;
        }

        if (nrf_instance_id.length() > Code.KEY_MAX_LENGTH) {
            logger.error("nrf_instance_id length {} is too large, max length is {}",
                         nrf_instance_id.length(), Code.KEY_MAX_LENGTH);
            return Code.NRF_INSTANCE_ID_LENGTH_EXCEED_MAX;
        }

        ByteString raw_nrf_profile = request.getRawNrfProfile();
        if (raw_nrf_profile.isEmpty() == true) {
            logger.error("raw_nrf_profile field is empty in NRFProfilePutRequest");
            return Code.EMPTY_RAW_NRF_PROFILE;
        }

        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {

        NRFProfilePutResponse nrf_profile_put_response = NRFProfilePutResponse.newBuilder().setCode(code).build();
        PutResponse put_response = PutResponse.newBuilder().setNrfProfilePutResponse(nrf_profile_put_response).build();
        return createNFMessage(put_response);
    }
}
