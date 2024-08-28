package ericsson.core.nrf.dbproxy.helper.nrfprofile;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.DelResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileDelRequestProto.NRFProfileDelRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileDelResponseProto.NRFProfileDelResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NRFProfileDelHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(NRFProfileDelHelper.class);

    private static NRFProfileDelHelper instance;

    private NRFProfileDelHelper()
    {
    }

    public static synchronized NRFProfileDelHelper getInstance()
    {
        if (null == instance) {
            instance = new NRFProfileDelHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        NRFProfileDelRequest request = message.getRequest().getDelRequest().getNrfProfileDelRequest();
        String nrf_instance_id = request.getNrfInstanceId();
        if (nrf_instance_id.isEmpty() == true) {
            logger.error("nrf_instance_id field is empty in NFProfileDelRequest");
            return Code.EMPTY_NRF_INSTANCE_ID;
        }

        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {
        NRFProfileDelResponse nrf_profile_del_response = NRFProfileDelResponse.newBuilder().setCode(code).build();
        DelResponse del_response = DelResponse.newBuilder().setNrfProfileDelResponse(nrf_profile_del_response).build();
        return createNFMessage(del_response);
    }
}
