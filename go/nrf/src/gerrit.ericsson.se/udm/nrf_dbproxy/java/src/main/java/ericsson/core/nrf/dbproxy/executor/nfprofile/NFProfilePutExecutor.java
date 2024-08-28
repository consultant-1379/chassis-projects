package ericsson.core.nrf.dbproxy.executor.nfprofile;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.helper.nfprofile.NFProfilePutHelper;
import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfilePutRequestProto.NFProfilePutRequest;

import org.apache.geode.pdx.JSONFormatter;
import org.apache.geode.pdx.PdxInstance;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NFProfilePutExecutor extends Executor
{
    private static final Logger logger = LogManager.getLogger(NFProfilePutExecutor.class);

    private static NFProfilePutExecutor instance = null;

    private NFProfilePutExecutor()
    {
        super(NFProfilePutHelper.getInstance());
    }

    public static synchronized NFProfilePutExecutor getInstance()
    {
        if(null == instance) {
            instance = new NFProfilePutExecutor();
        }
        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        NFProfilePutRequest put_request = request.getRequest().getPutRequest().getNfProfilePutRequest();
        String nf_instance_id = put_request.getNfInstanceId();
        String nf_profile = put_request.getNfProfile();
        int code = Code.CREATED;
        try {
            PdxInstance pdx_instance = JSONFormatter.fromJSON(nf_profile);
            code = ClientCacheService.getInstance().put(Code.NFPROFILE_INDICE, nf_instance_id, pdx_instance);
        } catch(Exception e) {
            logger.error(e.toString());
            code = Code.NF_PROFILE_FORMAT_ERROR;
        }
        return new ExecutionResult(code);
    }
}
