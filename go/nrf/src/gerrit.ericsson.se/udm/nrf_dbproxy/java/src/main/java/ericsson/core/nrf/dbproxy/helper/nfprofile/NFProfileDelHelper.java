package ericsson.core.nrf.dbproxy.helper.nfprofile;

import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileDelRequestProto.NFProfileDelRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileDelResponseProto.NFProfileDelResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;

public class NFProfileDelHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(NFProfileDelHelper.class);

    private static NFProfileDelHelper instance;

    private NFProfileDelHelper() { }

    public static synchronized NFProfileDelHelper getInstance()
    {
        if(null == instance) {
            instance = new NFProfileDelHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        NFProfileDelRequest request = message.getRequest().getDelRequest().getNfProfileDelRequest();
        String nf_instance_id = request.getNfInstanceId();
        if(nf_instance_id.isEmpty() == true) {
            logger.error("nf_instance_id field is empty in NFProfileDelRequest");
            return Code.EMPTY_NF_INSTANCE_ID;
        }

        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {

        NFProfileDelResponse nf_profile_del_response = NFProfileDelResponse.newBuilder().setCode(code).build();
        DelResponse del_response = DelResponse.newBuilder().setNfProfileDelResponse(nf_profile_del_response).build();
        return createNFMessage(del_response);
    }
}
