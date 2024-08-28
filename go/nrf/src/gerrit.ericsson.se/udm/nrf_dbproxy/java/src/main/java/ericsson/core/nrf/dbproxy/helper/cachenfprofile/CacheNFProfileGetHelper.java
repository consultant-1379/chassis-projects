package ericsson.core.nrf.dbproxy.helper.cachenfprofile;

import com.google.protobuf.ByteString;
import com.google.protobuf.util.JsonFormat;
import ericsson.core.nrf.dbproxy.clientcache.schema.CacheNFProfile;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.cachenfprofile.CacheNFProfileGetRequestProto.CacheNFProfileGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.cachenfprofile.CacheNFProfileGetResponseProto.CacheNFProfileGetResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.geode.pdx.JSONFormatter;
import org.apache.geode.pdx.PdxInstance;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class CacheNFProfileGetHelper extends Helper {

    private static final Logger logger = LogManager.getLogger(CacheNFProfileGetHelper.class);

    private static CacheNFProfileGetHelper instance;

    private CacheNFProfileGetHelper() {
    }

    public static synchronized CacheNFProfileGetHelper getInstance() {
        if (null == instance) {
            instance = new CacheNFProfileGetHelper();
        }
        return instance;
    }

    public int validate(NFMessage message) {
        CacheNFProfileGetRequest request = message.getRequest().getGetRequest().getCacheNfProfileGetRequest();
        String cache_nf_instance_id = request.getCacheNfInstanceId();
        if (cache_nf_instance_id.isEmpty()) {
            logger.error("Empty cache_nf_instance_id is set in CacheNFProfileGetRequest");
            return Code.EMPTY_CACHE_NF_INSTANCE_ID;
        } else if (cache_nf_instance_id.length() > Code.KEY_MAX_LENGTH) {
            logger.error("cache_nf_instance_id length {} is too large, max length is {}",
                    cache_nf_instance_id.length(), Code.KEY_MAX_LENGTH);
            return Code.CACHE_NF_INSTANCE_ID_LENGTH_EXCEED_MAX;
        }
        return Code.VALID;
    }

    public NFMessage createResponse(int code) {
        CacheNFProfileGetResponse cache_nf_profile_get_response = CacheNFProfileGetResponse.newBuilder().setCode(code).build();
        GetResponse get_response = GetResponse.newBuilder().setCacheNfProfileGetResponse(cache_nf_profile_get_response).build();
        return createNFMessage(get_response);
    }

    @Override
    public NFMessage createResponse(ExecutionResult execution_result) {
        if (execution_result.getCode() != Code.SUCCESS) {
            return createResponse(execution_result.getCode());
        } else {
            SearchResult search_result = (SearchResult)execution_result;
            if (search_result.getItems().size() > 0 ) {
                Object obj = search_result.getItems().get(0);
                String cache_nf_profile = JSONFormatter.toJSON((PdxInstance)obj);
                CacheNFProfileGetResponse cache_nf_profile_get_response = CacheNFProfileGetResponse.newBuilder().setCode(search_result.getCode()).setCacheNfProfile(cache_nf_profile).build();
                GetResponse get_response = GetResponse.newBuilder().setCacheNfProfileGetResponse(cache_nf_profile_get_response).build();
                return createNFMessage(get_response);
            } else {
                return createResponse(Code.DATA_NOT_EXIST);
            }

        }

    }


}
