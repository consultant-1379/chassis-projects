package ericsson.core.nrf.dbproxy.executor.cachenfprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.clientcache.schema.CacheNFProfile;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.functionservice.RemoteCacheMonitorThread;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.cachenfprofile.CacheNFProfilePutRequestProto.CacheNFProfilePutRequest;
import ericsson.core.nrf.dbproxy.helper.cachenfprofile.CacheNFProfilePutHelper;
import org.apache.geode.pdx.JSONFormatter;
import org.apache.geode.pdx.PdxInstance;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class CacheNFProfilePutExecutor extends Executor
{
    private static final Logger logger = LogManager.getLogger(CacheNFProfilePutExecutor.class);

    private static CacheNFProfilePutExecutor instance = null;

    private CacheNFProfilePutExecutor()
    {
        super(CacheNFProfilePutHelper.getInstance());
    }

    public static synchronized CacheNFProfilePutExecutor getInstance()
    {
        if (null == instance) {
            instance = new CacheNFProfilePutExecutor();
        }

        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        CacheNFProfilePutRequest put_request = request.getRequest().getPutRequest().getCacheNfProfilePutRequest();
        String cache_nf_instance_id = put_request.getCacheNfInstanceId();
        String cache_nf_profile = put_request.getRawCacheNfProfile();
        int code = Code.CREATED;
        try {
            PdxInstance pdx_cache_nf_profile = JSONFormatter.fromJSON(cache_nf_profile.toString());
            ClientCacheService.getInstance().put(Code.CACHENFPROFILE_INDICE, cache_nf_instance_id, pdx_cache_nf_profile);
            RemoteCacheMonitorThread.getInstance().incCacheOperationCount();
        } catch (Exception e) {
            logger.error(e.toString());
            code = Code.CACHE_NF_PROFILE_FORMAT_ERROR;
        }
        return new ExecutionResult(code);
    }

}
