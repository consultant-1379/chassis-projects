package ericsson.core.nrf.dbproxy.executor.nrfaddress;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressDelRequestProto.NRFAddressDelRequest;
import ericsson.core.nrf.dbproxy.helper.nrfaddress.NRFAddressDelHelper;

public class NRFAddressDeleteExecutor extends Executor
{
    private static NRFAddressDeleteExecutor instance = null;

    private NRFAddressDeleteExecutor()
    {
        super(NRFAddressDelHelper.getInstance());
    }

    public static synchronized NRFAddressDeleteExecutor getInstance()
    {
        if(null == instance) {
            instance = new NRFAddressDeleteExecutor();
        }
        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        NRFAddressDelRequest del_request = request.getRequest().getDelRequest().getNrfAddressDelRequest();
        String nrf_address_id = del_request.getNrfAddressId();
        int code = ClientCacheService.getInstance().delete(Code.NRFADDRESS_INDICE, nrf_address_id);
        return new ExecutionResult(code);
    }

}
