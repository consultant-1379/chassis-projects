package ericsson.core.nrf.dbproxy.executor.nrfprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileDelRequestProto.NRFProfileDelRequest;
import ericsson.core.nrf.dbproxy.helper.nrfprofile.NRFProfileDelHelper;

public class NRFProfileDeleteExecutor extends Executor {

  private static NRFProfileDeleteExecutor instance;

  static {
    instance = null;
  }

  private NRFProfileDeleteExecutor() {
    super(NRFProfileDelHelper.getInstance());
  }

  public static synchronized NRFProfileDeleteExecutor getInstance() {
    if (null == instance) {
      instance = new NRFProfileDeleteExecutor();
    }

    return instance;
  }

  protected ExecutionResult execute(NFMessage request) {
    NRFProfileDelRequest delRequest = request.getRequest().getDelRequest()
        .getNrfProfileDelRequest();
    String nrfInstanceId = delRequest.getNrfInstanceId();
    int code = ClientCacheService.getInstance().delete(Code.NRFPROFILE_INDICE, nrfInstanceId);
    return new ExecutionResult(code);
  }
}
