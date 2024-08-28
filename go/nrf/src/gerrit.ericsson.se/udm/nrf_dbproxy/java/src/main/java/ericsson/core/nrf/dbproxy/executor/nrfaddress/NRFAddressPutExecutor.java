package ericsson.core.nrf.dbproxy.executor.nrfaddress;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.clientcache.schema.NRFAddress;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressIndexProto.NRFAddressIndex;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressPutRequestProto.NRFAddressPutRequest;
import ericsson.core.nrf.dbproxy.helper.nrfaddress.NRFAddressPutHelper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NRFAddressPutExecutor extends Executor
{
    private static final Logger logger = LogManager.getLogger(NRFAddressPutExecutor.class);

    private static NRFAddressPutExecutor instance = null;

    private NRFAddressPutExecutor()
    {
        super(NRFAddressPutHelper.getInstance());
    }

    public static synchronized NRFAddressPutExecutor getInstance()
    {
        if(null == instance) {
            instance = new NRFAddressPutExecutor();
        }

        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        NRFAddressPutRequest put_request = request.getRequest().getPutRequest().getNrfAddressPutRequest();
        String nrf_address_id = put_request.getNrfAddressId();
        NRFAddress nrf_address = createNRFAddress(put_request);
        int code = ClientCacheService.getInstance().put(Code.NRFADDRESS_INDICE, nrf_address_id, nrf_address);
        return new ExecutionResult(code);
    }

    private NRFAddress createNRFAddress(NRFAddressPutRequest request)
    {
        NRFAddressIndex index = request.getIndex();

        NRFAddress nrf_address = new NRFAddress();

        nrf_address.setNRFAddressID(request.getNrfAddressId());
        nrf_address.setData(request.getNrfAddressData());
        nrf_address.setKey1(index.getNrfAddressKey1());
        nrf_address.setKey2(index.getNrfAddressKey2());
        nrf_address.setKey3(index.getNrfAddressKey3());
        nrf_address.setKey4(index.getNrfAddressKey4());
        nrf_address.setKey5(index.getNrfAddressKey5());

        logger.debug("NRF Address : {}", nrf_address.toString());

        return nrf_address;
    }

}
