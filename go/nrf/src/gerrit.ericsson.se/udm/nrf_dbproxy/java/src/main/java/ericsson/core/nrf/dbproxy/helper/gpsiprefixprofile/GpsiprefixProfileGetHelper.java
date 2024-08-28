package ericsson.core.nrf.dbproxy.helper.gpsiprefixprofile;

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
import ericsson.core.nrf.dbproxy.clientcache.schema.GpsiprefixProfiles;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetResponseProto.GpsiprefixProfileGetResponse;

import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;

public class GpsiprefixProfileGetHelper extends Helper
{
    public static final long SEARCH_GPSI_MAX_VALUE = 999999999999999L;

    private static final Logger logger = LogManager.getLogger(GpsiprefixProfileGetHelper.class);

    private static GpsiprefixProfileGetHelper instance;

    private GpsiprefixProfileGetHelper() { }

    public static synchronized GpsiprefixProfileGetHelper getInstance()
    {
        if(null == instance) {
            instance = new GpsiprefixProfileGetHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        GpsiprefixProfileGetRequest request = message.getRequest().getGetRequest().getGpsiprefixProfileGetRequest();
       
        Long search_gpsi = request.getSearchGpsi();
        if(search_gpsi == 0L) {
            logger.error("value 0 search_gpsi is set in GpsiperfixProfileGetRequest");
            return Code.EMPTY_SEARCH_GPSI;
        } else if(search_gpsi > SEARCH_GPSI_MAX_VALUE) {
            logger.error("SEARCH_GPSI_MAX_VALUE search_gpsi is set in GpsiperfixProfileGetRequest");
            return Code.SEARCH_GPSI_LENGTH_EXCEED_MAX;
        }
        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {
        GpsiprefixProfileGetResponse nrf_address_get_response = GpsiprefixProfileGetResponse.newBuilder().setCode(code).build();
        GetResponse get_response = GetResponse.newBuilder().setGpsiprefixProfileGetResponse(nrf_address_get_response).build();
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
                Map<Long,GpsiprefixProfiles>  search_result_map = (Map<Long,GpsiprefixProfiles>)obj;
                for (Map.Entry<Long,GpsiprefixProfiles> entry : search_result_map.entrySet()){
					Long key = entry.getKey();
					if (search_result_map.get(key) != null) {
						GpsiprefixProfiles gpsiprefixProfiles = search_result_map.get(key);						
	                     Iterator iter = gpsiprefixProfiles.getValueInfo().keySet().iterator();
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
            GpsiprefixProfileGetResponse gpsiprefix_profile_get_response = GpsiprefixProfileGetResponse.newBuilder().setCode(execution_result.getCode()).addAllValueInfo(value_info_list).build();
            GetResponse get_response = GetResponse.newBuilder().setGpsiprefixProfileGetResponse(gpsiprefix_profile_get_response).build();
            return createNFMessage(get_response);
        }

    }
}
