package ericsson.core.nrf.dbproxy.helper;

import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.DelResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PutResponse;

public abstract class Helper {

  public Helper() {
  }

  public abstract int validate(NFMessage message);

  public abstract NFMessage createResponse(int code);

  public NFMessage createResponse(ExecutionResult executionResult) {
    return createResponse(executionResult.getCode());
  }

  protected NFMessage createNFMessage(PutResponse putResponse) {
    NFResponse response = NFResponse.newBuilder().setPutResponse(putResponse).build();
    return NFMessage.newBuilder().setResponse(response).build();
  }

  protected NFMessage createNFMessage(GetResponse getResponse) {
    NFResponse response = NFResponse.newBuilder().setGetResponse(getResponse).build();
    return NFMessage.newBuilder().setResponse(response).build();
  }

  protected NFMessage createNFMessage(DelResponse delResponse) {
    NFResponse response = NFResponse.newBuilder().setDelResponse(delResponse).build();
    return NFMessage.newBuilder().setResponse(response).build();
  }


}
