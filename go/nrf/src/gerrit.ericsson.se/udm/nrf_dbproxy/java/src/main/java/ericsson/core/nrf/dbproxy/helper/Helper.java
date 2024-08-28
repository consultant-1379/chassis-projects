package ericsson.core.nrf.dbproxy.helper;

import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;

public abstract class Helper
{
    public Helper() {}

    public abstract int validate(NFMessage message);
    public abstract NFMessage createResponse(int code);
    public NFMessage createResponse(ExecutionResult execution_result)
    {
        return createResponse(execution_result.getCode());
    }

    protected NFMessage createNFMessage(PutResponse put_response)
    {
        NFResponse response = NFResponse.newBuilder().setPutResponse(put_response).build();
        return NFMessage.newBuilder().setResponse(response).build();
    }

    protected NFMessage createNFMessage(GetResponse get_response)
    {
        NFResponse response = NFResponse.newBuilder().setGetResponse(get_response).build();
        return NFMessage.newBuilder().setResponse(response).build();
    }

    protected NFMessage createNFMessage(DelResponse del_response)
    {
        NFResponse response = NFResponse.newBuilder().setDelResponse(del_response).build();
        return NFMessage.newBuilder().setResponse(response).build();
    }


}
