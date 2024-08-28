package ericsson.core.nrf.dbproxy.executor.nrfprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileDelRequestProto.NRFProfileDelRequest;
import ericsson.core.nrf.dbproxy.helper.nrfprofile.NRFProfileDelHelper;

public class NRFProfileDeleteExecutor extends Executor
{
    private static NRFProfileDeleteExecutor instance = null;

    private NRFProfileDeleteExecutor()
    {
        super(NRFProfileDelHelper.getInstance());
    }

    public static synchronized NRFProfileDeleteExecutor getInstance()
    {
        if (null == instance) {
            instance = new NRFProfileDeleteExecutor();
        }

        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        NRFProfileDelRequest del_request = request.getRequest().getDelRequest().getNrfProfileDelRequest();
        String nrf_instance_id = del_request.getNrfInstanceId();
        int code = ClientCacheService.getInstance().delete(Code.NRFPROFILE_INDICE, nrf_instance_id);
        return new ExecutionResult(code);
    }
}
