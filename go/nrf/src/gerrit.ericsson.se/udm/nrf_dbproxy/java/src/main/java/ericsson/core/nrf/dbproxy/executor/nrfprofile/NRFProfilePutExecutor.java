package ericsson.core.nrf.dbproxy.executor.nrfprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.clientcache.schema.KeyAggregation;
import ericsson.core.nrf.dbproxy.clientcache.schema.NRFProfile;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileIndexProto.NRFKeyStruct;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfilePutRequestProto.NRFProfilePutRequest;
import ericsson.core.nrf.dbproxy.helper.nrfprofile.NRFProfilePutHelper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.util.HashMap;
import java.util.List;

public class NRFProfilePutExecutor extends Executor
{
    private static final Logger logger = LogManager.getLogger(NRFProfilePutExecutor.class);

    private static NRFProfilePutExecutor instance = null;

    private NRFProfilePutExecutor()
    {
        super(NRFProfilePutHelper.getInstance());
    }

    public static synchronized NRFProfilePutExecutor getInstance()
    {
        if (null == instance) {
            instance = new NRFProfilePutExecutor();
        }

        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        NRFProfilePutRequest put_request = request.getRequest().getPutRequest().getNrfProfilePutRequest();
        String nrf_instance_id = put_request.getNrfInstanceId();
        NRFProfile nrf_profile = createNRFProfile(put_request);
        int code = ClientCacheService.getInstance().put(Code.NRFPROFILE_INDICE, nrf_instance_id, nrf_profile);
        return new ExecutionResult(code);
    }

    private NRFProfile createNRFProfile(NRFProfilePutRequest request)
    {
        NRFProfile nrf_profile = new NRFProfile();
        nrf_profile.setRaw_data(request.getRawNrfProfile());


        nrf_profile.setAMFKey1(processNRFKeyStructList(request.getIndex().getAmfKey1List()));
        nrf_profile.setAMFKey2(processNRFKeyStructList(request.getIndex().getAmfKey2List()));
        nrf_profile.setAMFKey3(processNRFKeyStructList(request.getIndex().getAmfKey3List()));
        nrf_profile.setAMFKey4(processNRFKeyStructList(request.getIndex().getAmfKey4List()));

        nrf_profile.setSMFKey1(processNRFKeyStructList(request.getIndex().getSmfKey1List()));
        nrf_profile.setSMFKey2(processNRFKeyStructList(request.getIndex().getSmfKey2List()));
        nrf_profile.setSMFKey3(processNRFKeyStructList(request.getIndex().getSmfKey3List()));

        nrf_profile.setUDMKey1(processNRFKeyStructList(request.getIndex().getUdmKey1List()));
        nrf_profile.setUDMKey2(processNRFKeyStructList(request.getIndex().getUdmKey2List()));

        nrf_profile.setAUSFKey1(processNRFKeyStructList(request.getIndex().getAusfKey1List()));
        nrf_profile.setAUSFKey2(processNRFKeyStructList(request.getIndex().getAusfKey2List()));

        nrf_profile.setPCFKey1(processNRFKeyStructList(request.getIndex().getPcfKey1List()));
        nrf_profile.setPCFKey2(processNRFKeyStructList(request.getIndex().getPcfKey2List()));


        nrf_profile.setKey1(request.getIndex().getKey1());
        nrf_profile.setKey3(request.getIndex().getKey3());

        logger.trace("NRF Profile : {}", nrf_profile.toString());

        return nrf_profile;
    }

    private HashMap processNRFKeyStructList(List<NRFKeyStruct> nrfKeyStructs)
    {
        HashMap result = new HashMap();
        int id = 0;
        for (NRFKeyStruct nks : nrfKeyStructs) {

            KeyAggregation ka = new KeyAggregation();
            ka.setSubKey1(nks.getSubKey1());
            ka.setSubKey2(nks.getSubKey2());
            ka.setSubKey3(nks.getSubKey3());
            ka.setSubKey4(nks.getSubKey4());
            ka.setSubKey5(nks.getSubKey5());

            result.put(id, ka);
            id++;
        }
        return result;
    }
}
