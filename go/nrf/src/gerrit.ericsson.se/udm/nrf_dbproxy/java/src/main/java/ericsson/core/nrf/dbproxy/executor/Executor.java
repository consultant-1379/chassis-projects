package ericsson.core.nrf.dbproxy.executor;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.helper.Helper;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;

public abstract class Executor
{
    private static final Logger logger = LogManager.getLogger(Executor.class);

    protected Helper helper;

    public Executor(Helper helper)
    {
        this.helper = helper;
    }

    protected  abstract ExecutionResult execute(NFMessage request);

    public NFMessage process(NFMessage request)
    {

        logger.trace("\n" + request.toString());

        int code = helper.validate(request);
        if(code != Code.VALID)
	{
            NFMessage response = helper.createResponse(code);
            logger.trace("\n" + response.toString());
	    return response;
	}

        ExecutionResult execution_result = execute(request);
        NFMessage response = helper.createResponse(execution_result);

        logger.trace("Result: " + execution_result.toString());
        logger.trace("\n" + response.toString());

        return response;
    }
}
