package ericsson.core.nrf.dbproxy.helper.nrfaddress;

import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressDelRequestProto.NRFAddressDelRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressDelResponseProto.NRFAddressDelResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;

public class NRFAddressDelHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(NRFAddressDelHelper.class);

    private static NRFAddressDelHelper instance;

    private NRFAddressDelHelper() { }

    public static synchronized NRFAddressDelHelper getInstance()
    {
        if(null == instance) {
            instance = new NRFAddressDelHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        NRFAddressDelRequest request = message.getRequest().getDelRequest().getNrfAddressDelRequest();
        String nrf_address_id = request.getNrfAddressId();
        if(nrf_address_id.isEmpty() == true) {
            logger.error("nrf_address_id field is empty in NRFAddressDelRequest");
            return Code.EMPTY_NRF_ADDRESS_ID;
        }

        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {

        NRFAddressDelResponse nrf_address_del_response = NRFAddressDelResponse.newBuilder().setCode(code).build();
        DelResponse del_response = DelResponse.newBuilder().setNrfAddressDelResponse(nrf_address_del_response).build();
        return createNFMessage(del_response);
    }
}
