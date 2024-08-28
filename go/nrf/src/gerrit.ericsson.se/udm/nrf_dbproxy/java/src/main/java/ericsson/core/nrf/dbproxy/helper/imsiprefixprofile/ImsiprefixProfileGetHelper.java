package ericsson.core.nrf.dbproxy.helper.imsiprefixprofile;

import java.util.List;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;
import java.util.Iterator;

import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import com.google.protobuf.ByteString;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.clientcache.schema.ImsiprefixProfiles;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.imsiprefixprofile.ImsiprefixProfileGetRequestProto.ImsiprefixProfileGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.imsiprefixprofile.ImsiprefixProfileGetResponseProto.ImsiprefixProfileGetResponse;

import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;

public class ImsiprefixProfileGetHelper extends Helper
{

    public static final long SEARCH_IMSI_MAX_VALUE = 999999999999999L;

    private static final Logger logger = LogManager.getLogger(ImsiprefixProfileGetHelper.class);

    private static ImsiprefixProfileGetHelper instance;

    private ImsiprefixProfileGetHelper() { }

    public static synchronized ImsiprefixProfileGetHelper getInstance()
    {
        if(null == instance) {
            instance = new ImsiprefixProfileGetHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        ImsiprefixProfileGetRequest request = message.getRequest().getGetRequest().getImsiprefixProfileGetRequest();

        Long search_imsi = request.getSearchImsi();
        if(search_imsi == 0L) {
            logger.error("value 0 search_imsi is set in ImsiperfixProfileGetRequest");
            return Code.EMPTY_SEARCH_IMSI;
        }else if(search_imsi > SEARCH_IMSI_MAX_VALUE) {
            logger.error("SEARCH_IMSI_MAX_VALUE search_imsi is set in ImsiperfixProfileGetRequest");
            return Code.SEARCH_IMSI_LENGTH_EXCEED_MAX;
        }
        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {
        ImsiprefixProfileGetResponse nrf_address_get_response = ImsiprefixProfileGetResponse.newBuilder().setCode(code).build();
        GetResponse get_response = GetResponse.newBuilder().setImsiprefixProfileGetResponse(nrf_address_get_response).build();
        return createNFMessage(get_response);
    }

    public NFMessage createResponse(ExecutionResult execution_result)
    {
        if(execution_result.getCode() != Code.SUCCESS) {
            return createResponse(execution_result.getCode());
        } else {
            List<String> value_info_list = new ArrayList<>();
            SearchResult search_result = (SearchResult)execution_result;
            for(Object obj : search_result.getItems()) {
                Map<Long,ImsiprefixProfiles>  search_result_map = (Map<Long,ImsiprefixProfiles>)obj;
                for (Map.Entry<Long,ImsiprefixProfiles> entry : search_result_map.entrySet()) {
                    Long key = entry.getKey();
                    if (search_result_map.get(key) != null) {
                        ImsiprefixProfiles imsiprefixProfiles = search_result_map.get(key);
                        Iterator iter = imsiprefixProfiles.getValueInfo().keySet().iterator();
                        while (iter.hasNext()) {
                            String value = (String)(iter.next());
                            value_info_list.add(value);
                        }
                    }
                }
            }
            if (value_info_list.isEmpty()) {
                return createResponse(Code.DATA_NOT_EXIST);
            }			
            ImsiprefixProfileGetResponse imsiprefix_profile_get_response = ImsiprefixProfileGetResponse.newBuilder().setCode(execution_result.getCode()).addAllValueInfo(value_info_list).build();
            GetResponse get_response = GetResponse.newBuilder().setImsiprefixProfileGetResponse(imsiprefix_profile_get_response).build();
            return createNFMessage(get_response);
        }

    }
}
