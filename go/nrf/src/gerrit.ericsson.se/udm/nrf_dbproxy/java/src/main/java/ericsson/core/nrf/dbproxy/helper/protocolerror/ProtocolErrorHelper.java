package ericsson.core.nrf.dbproxy.helper.protocolerror;

import ericsson.core.nrf.dbproxy.helper.Helper;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;

public class ProtocolErrorHelper extends Helper
{

    private static ProtocolErrorHelper instance;

    private ProtocolErrorHelper() { }

    public static synchronized ProtocolErrorHelper getInstance()
    {
        if(null == instance) {
            instance = new ProtocolErrorHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {
        return Code.NFMESSAGE_PROTOCOL_ERROR;
    }

    public NFMessage createResponse(int code)
    {

        NFProtocolError protocol_error = NFProtocolError.newBuilder().setCode(code).build();
        return NFMessage.newBuilder().setProtocolError(protocol_error).build();
    }
}
