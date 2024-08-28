package ericsson.core.nrf.dbproxy.executor;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public abstract class Executor {

  private static final Logger LOGGER = LogManager.getLogger(Executor.class);

  protected Helper helper;

  public Executor(Helper helper) {
    this.helper = helper;
  }

  protected abstract ExecutionResult execute(NFMessage request);

  public NFMessage process(NFMessage request) {

    LOGGER.debug("\n" + request.toString());

    int code = helper.validate(request);
    if (code != Code.VALID) {
      NFMessage response = helper.createResponse(code);
      LOGGER.debug("\n" + response.toString());
      return response;
    }

    ExecutionResult executionResult = execute(request);
    NFMessage response = helper.createResponse(executionResult);

    LOGGER.debug("Result: " + executionResult.toString());
    LOGGER.debug("\n" + response.toString());

    return response;
  }
}
