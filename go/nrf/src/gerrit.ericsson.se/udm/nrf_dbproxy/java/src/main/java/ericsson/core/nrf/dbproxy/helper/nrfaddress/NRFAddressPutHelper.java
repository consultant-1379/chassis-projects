package ericsson.core.nrf.dbproxy.helper.nrfaddress;

import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import com.google.protobuf.ByteString;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressPutRequestProto.NRFAddressPutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressPutResponseProto.NRFAddressPutResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;

public class NRFAddressPutHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(NRFAddressPutHelper.class);

    private static NRFAddressPutHelper instance;

    private NRFAddressPutHelper() { }

    public static synchronized NRFAddressPutHelper getInstance()
    {
        if(null == instance) {
            instance = new NRFAddressPutHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        NRFAddressPutRequest request = message.getRequest().getPutRequest().getNrfAddressPutRequest();
        String nrf_address_id = request.getNrfAddressId();
        if(nrf_address_id.isEmpty() == true) {
            logger.error("nrf_address_id field is empty in NRFAddressPutRequest");
            return Code.EMPTY_NRF_ADDRESS_ID;
        }

        if(nrf_address_id.length() > Code.KEY_MAX_LENGTH) {
            logger.error("nrf_address_id length {} is too large, max length is {}",
                         nrf_address_id.length(), Code.KEY_MAX_LENGTH);
            return Code.NRF_ADDRESS_ID_LENGTH_EXCEED_MAX;
        }

        ByteString nrf_address_data = request.getNrfAddressData();
        if(nrf_address_data.isEmpty() == true) {
            logger.error("nrf_address_data field is empty in NRFAddressPutRequest");
            return Code.EMPTY_NRF_ADDRESS_DATA;
        }

        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {

        NRFAddressPutResponse nrf_address_put_response = NRFAddressPutResponse.newBuilder().setCode(code).build();
        PutResponse put_response = PutResponse.newBuilder().setNrfAddressPutResponse(nrf_address_put_response).build();
        return createNFMessage(put_response);
    }
}
