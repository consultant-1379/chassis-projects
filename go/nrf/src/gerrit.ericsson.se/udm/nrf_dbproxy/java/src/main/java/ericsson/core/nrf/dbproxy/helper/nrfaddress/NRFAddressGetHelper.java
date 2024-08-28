package ericsson.core.nrf.dbproxy.helper.nrfaddress;

import java.util.List;
import java.util.ArrayList;

import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import com.google.protobuf.ByteString;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.clientcache.schema.NRFAddress;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressGetRequestProto.NRFAddressGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressGetResponseProto.NRFAddressGetResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressGetRequestProto.NRFAddressGetRequest.DataCase;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;

public class NRFAddressGetHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(NRFAddressGetHelper.class);

    private static NRFAddressGetHelper instance;

    private NRFAddressGetHelper() { }

    public static synchronized NRFAddressGetHelper getInstance()
    {
        if(null == instance) {
            instance = new NRFAddressGetHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        NRFAddressGetRequest request = message.getRequest().getGetRequest().getNrfAddressGetRequest();
        DataCase data_case = request.getDataCase();
        if(data_case == DataCase.NRF_ADDRESS_ID) {

            String nrf_address_id = request.getNrfAddressId();
            if(nrf_address_id.isEmpty() == true) {
                logger.error("Empty nrf_address_id is set in NRFAddressGetRequest");
                return Code.EMPTY_NRF_ADDRESS_ID;
            } else if(nrf_address_id.length() > Code.KEY_MAX_LENGTH) {
                logger.error("nrf_address_id length {} is too large, max length is {}",
                             nrf_address_id.length(), Code.KEY_MAX_LENGTH);
                return Code.NRF_ADDRESS_ID_LENGTH_EXCEED_MAX;
            }
        } else if(data_case == DataCase.FILTER) {

            if(request.getFilter().getIndex().getNrfAddressKey1List().isEmpty() &&
               request.getFilter().getIndex().getNrfAddressKey2().isEmpty() &&
               request.getFilter().getIndex().getNrfAddressKey3().isEmpty() &&
               request.getFilter().getIndex().getNrfAddressKey4().isEmpty() &&
               request.getFilter().getIndex().getNrfAddressKey5().isEmpty()) {
                logger.error("Empty NRFAddressFilter is set in filter of NRFAddressGetRequest");
                return Code.EMPTY_NRF_ADDRESS_FILTER;
            }
        } else {
            logger.error("Empty NRFAddressGetRequest is received");
            return Code.NFMESSAGE_PROTOCOL_ERROR;
        }

        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {

        NRFAddressGetResponse nrf_address_get_response = NRFAddressGetResponse.newBuilder().setCode(code).build();
        GetResponse get_response = GetResponse.newBuilder().setNrfAddressGetResponse(nrf_address_get_response).build();
        return createNFMessage(get_response);
    }

    public NFMessage createResponse(ExecutionResult execution_result)
    {

        if(execution_result.getCode() != Code.SUCCESS) {
            return createResponse(execution_result.getCode());
        } else {
            SearchResult search_result = (SearchResult)execution_result;

            List<ByteString> nrf_addresses = new ArrayList<>();
            for(Object obj : search_result.getItems()) {
                NRFAddress item = (NRFAddress)obj;
                nrf_addresses.add(item.getData());
            }
            NRFAddressGetResponse nrf_address_get_response = NRFAddressGetResponse.newBuilder().setCode(execution_result.getCode()).addAllNrfAddressData(nrf_addresses).build();
            GetResponse get_response = GetResponse.newBuilder().setNrfAddressGetResponse(nrf_address_get_response).build();
            return createNFMessage(get_response);
        }

    }
}
