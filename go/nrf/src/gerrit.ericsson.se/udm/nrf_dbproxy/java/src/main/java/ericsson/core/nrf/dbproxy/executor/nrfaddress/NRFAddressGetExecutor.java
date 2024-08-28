package ericsson.core.nrf.dbproxy.executor.nrfaddress;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressFilterProto.NRFAddressFilter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressGetRequestProto.NRFAddressGetRequest;
import ericsson.core.nrf.dbproxy.helper.nrfaddress.NRFAddressGetHelper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import java.util.List;

public class NRFAddressGetExecutor extends Executor
{
    private static final Logger logger = LogManager.getLogger(NRFAddressGetExecutor.class);

    private static NRFAddressGetExecutor instance = null;

    private NRFAddressGetExecutor()
    {
        super(NRFAddressGetHelper.getInstance());
    }

    public static synchronized NRFAddressGetExecutor getInstance()
    {
        if(null == instance) {
            instance = new NRFAddressGetExecutor();
        }
        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        NRFAddressGetRequest get_request = request.getRequest().getGetRequest().getNrfAddressGetRequest();
        switch(get_request.getDataCase()) {
        case NRF_ADDRESS_ID:
            return ClientCacheService.getInstance().getByID(Code.NRFADDRESS_INDICE, get_request.getNrfAddressId());
        case FILTER:
            String query_string = getQueryString(get_request.getFilter());
            return ClientCacheService.getInstance().getByFilter(Code.NRFADDRESS_INDICE, query_string);
        default:
            return new ExecutionResult(Code.NFMESSAGE_PROTOCOL_ERROR);
        }
    }

    private String getQueryString(NRFAddressFilter filter)
    {

        String region_name = Code.NRFADDRESS_INDICE;
        StringBuilder sb = new StringBuilder("SELECT * FROM /" + region_name + " p WHERE ");

        String operation = "OR";
        if(filter.getAndOperation() == true) operation = "AND";

        List<String> key1_list = filter.getIndex().getNrfAddressKey1List();
        String key2 = filter.getIndex().getNrfAddressKey2();
        String key3 = filter.getIndex().getNrfAddressKey3();
        String key4 = filter.getIndex().getNrfAddressKey4();
        String key5 = filter.getIndex().getNrfAddressKey5();

		boolean key_exist = false;
        if (key1_list.isEmpty() == false) {
            boolean needOR = false;
            sb.append("(");
            for(String key1: key1_list) {
                if(needOR == true)
                    sb.append(" OR ");
                sb.append("p.key1 = '" + key1 + "'");
                needOR = true;
                key_exist = true;
            }
            sb.append(")");
        }
		
        if(!key2.isEmpty()) {
            if(key_exist)
                sb.append(" " + operation + " ");
            sb.append("p.key2 = '" + key2 + "'");
            key_exist = true;
        }

        if(!key3.isEmpty()) {
            if(key_exist)
                sb.append(" " + operation + " ");
            sb.append("p.key3 = '" + key3 + "'");
            key_exist = true;
        }

        if(!key4.isEmpty()) {
            if(key_exist)
                sb.append(" " + operation + " ");
            sb.append("p.key4 = '" + key4 + "'");
            key_exist = true;
        }

        if(!key5.isEmpty()) {
            if(key_exist)
                sb.append(" " + operation + " ");
            sb.append("p.key5 = '" + key5 + "'");
        }

        String query_string = sb.toString();
        logger.debug("OQL = {}", query_string);

        return query_string;

    }

}
