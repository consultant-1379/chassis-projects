package ericsson.core.nrf.dbproxy.executor.protocolerror;

import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.helper.protocolerror.ProtocolErrorHelper;

public class ProtocolErrorExecutor extends Executor
{

    private static ProtocolErrorExecutor instance = null;

    private ProtocolErrorExecutor()
    {
        super(ProtocolErrorHelper.getInstance());
    }

    public static synchronized ProtocolErrorExecutor getInstance()
    {
        if(null == instance) {
            instance = new ProtocolErrorExecutor();
        }

        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        return null;
    }
}
