package ericsson.core.nrf.dbproxy.executor.gpsiprofile;

import java.util.List;
import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileDelRequestProto.GpsiProfileDelRequest;
import ericsson.core.nrf.dbproxy.helper.gpsiprofile.GpsiProfileDelHelper;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileProto.GpsiprefixProfile;
import ericsson.core.nrf.dbproxy.executor.common.GpsiPrefixProfilesUtil;

public class GpsiProfileDeleteExecutor extends Executor
{
    private static final Logger logger = LogManager.getLogger(GpsiProfileDeleteExecutor.class);

    private static GpsiProfileDeleteExecutor instance = null;

    private GpsiProfileDeleteExecutor()
    {
        super(GpsiProfileDelHelper.getInstance());
    }

    public static synchronized GpsiProfileDeleteExecutor getInstance()
    {
        if(null == instance) {
            instance = new GpsiProfileDeleteExecutor();
        }
        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        int code = Code.SUCCESS;
        try {
            GpsiProfileDelRequest del_request = request.getRequest().getDelRequest().getGpsiProfileDelRequest();
            String gpsi_profile_id = del_request.getGpsiProfileId();
            List<GpsiprefixProfile>  profileDelList = del_request.getGpsiPrefixDeleteList();

            ClientCacheService.getInstance().getCacheTransactionManager().begin();
            for(GpsiprefixProfile gpsiprefixProfile : profileDelList) {
                int retCode = GpsiPrefixProfilesUtil.DelGpsiprefixProfile(gpsiprefixProfile);
                if (Code.SUCCESS != retCode) {
                    ClientCacheService.getInstance().getCacheTransactionManager().rollback();
                    return new ExecutionResult(retCode);
                }
            }
            code = ClientCacheService.getInstance().delete(Code.GPSIPROFILE_INDICE, gpsi_profile_id,false);
            if (Code.SUCCESS != code) {
                ClientCacheService.getInstance().getCacheTransactionManager().rollback();
                return new ExecutionResult(code);
            }
            ClientCacheService.getInstance().getCacheTransactionManager().commit();
        } catch (Exception e) {
            logger.error(e.toString());
            code = Code.INTERNAL_ERROR;
        }
        return new ExecutionResult(code);
    }

}
