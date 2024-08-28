package ericsson.core.nrf.dbproxy.helper.cachenfprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PutResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.cachenfprofile.CacheNFProfilePutRequestProto.CacheNFProfilePutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.cachenfprofile.CacheNFProfilePutResponseProto.CacheNFProfilePutResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class CacheNFProfilePutHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(CacheNFProfilePutHelper.class);

    private static CacheNFProfilePutHelper instance;

    private CacheNFProfilePutHelper()
    {
    }

    public static synchronized CacheNFProfilePutHelper getInstance()
    {
        if (null == instance) {
            instance = new CacheNFProfilePutHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        CacheNFProfilePutRequest request = message.getRequest().getPutRequest().getCacheNfProfilePutRequest();
        String cache_nf_instance_id = request.getCacheNfInstanceId();
        if (cache_nf_instance_id.isEmpty() == true) {
            logger.error("cache_f_instance_id field is empty in NRFProfilePutRequest");
            return Code.EMPTY_CACHE_NF_INSTANCE_ID;
        }

        if (cache_nf_instance_id.length() > Code.KEY_MAX_LENGTH) {
            logger.error("cache_nf_instance_id length {} is too large, max length is {}",
                         cache_nf_instance_id.length(), Code.KEY_MAX_LENGTH);
            return Code.CACHE_NF_INSTANCE_ID_LENGTH_EXCEED_MAX;
        }

        String raw_cache_nf_profile = request.getRawCacheNfProfile();
        if (raw_cache_nf_profile.isEmpty() == true) {
            logger.error("raw_cache_nf_profile field is empty in CacheNFProfilePutRequest");
            return Code.EMPTY_RAW_CACHE_NF_PROFILE;
        }

        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {

        CacheNFProfilePutResponse cache_nf_profile_put_response = CacheNFProfilePutResponse.newBuilder().setCode(code).build();
        PutResponse put_response = PutResponse.newBuilder().setCacheNfProfilePutResponse(cache_nf_profile_put_response).build();
        return createNFMessage(put_response);
    }
}
