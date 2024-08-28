package ericsson.core.nrf.dbproxy.executor.nfprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileDelRequestProto.NFProfileDelRequest;
import ericsson.core.nrf.dbproxy.helper.nfprofile.NFProfileDelHelper;

public class NFProfileDeleteExecutor extends Executor
{
    private static NFProfileDeleteExecutor instance = null;

    private NFProfileDeleteExecutor()
    {
        super(NFProfileDelHelper.getInstance());
    }

    public static synchronized NFProfileDeleteExecutor getInstance()
    {
        if(null == instance) {
            instance = new NFProfileDeleteExecutor();
        }
        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        NFProfileDelRequest del_request = request.getRequest().getDelRequest().getNfProfileDelRequest();
        String nf_instance_id = del_request.getNfInstanceId();
        int code = ClientCacheService.getInstance().delete(Code.NFPROFILE_INDICE, nf_instance_id);
        return new ExecutionResult(code);
    }
}
