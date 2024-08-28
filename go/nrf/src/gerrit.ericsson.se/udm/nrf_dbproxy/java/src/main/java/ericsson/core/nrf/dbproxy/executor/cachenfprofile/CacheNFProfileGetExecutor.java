package ericsson.core.nrf.dbproxy.executor.cachenfprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.cachenfprofile.CacheNFProfileGetRequestProto.CacheNFProfileGetRequest;
import ericsson.core.nrf.dbproxy.helper.cachenfprofile.CacheNFProfileGetHelper;


public class CacheNFProfileGetExecutor extends Executor
{

    private static CacheNFProfileGetExecutor instance = null;

    private CacheNFProfileGetExecutor()
    {
        super(CacheNFProfileGetHelper.getInstance());
    }

    public static synchronized CacheNFProfileGetExecutor getInstance()
    {
        if (null == instance) {
            instance = new CacheNFProfileGetExecutor();
        }
        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        CacheNFProfileGetRequest get_request = request.getRequest().getGetRequest().getCacheNfProfileGetRequest();
        return ClientCacheService.getInstance().getByID(Code.CACHENFPROFILE_INDICE, get_request.getCacheNfInstanceId());
    }
}
