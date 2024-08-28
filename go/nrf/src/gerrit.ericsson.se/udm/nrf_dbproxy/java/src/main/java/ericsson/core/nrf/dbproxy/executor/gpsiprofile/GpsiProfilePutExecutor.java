package ericsson.core.nrf.dbproxy.executor.gpsiprofile;

import java.util.List;

import ericsson.core.nrf.dbproxy.executor.Executor;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.helper.gpsiprofile.GpsiProfilePutHelper;
import ericsson.core.nrf.dbproxy.clientcache.schema.GpsiProfile;
import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfilePutRequestProto.GpsiProfilePutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileIndexProto.GpsiProfileIndex;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileProto.GpsiprefixProfile;
import ericsson.core.nrf.dbproxy.executor.common.GpsiPrefixProfilesUtil;


public class GpsiProfilePutExecutor extends Executor
{
    private static final Logger logger = LogManager.getLogger(GpsiProfilePutExecutor.class);

    private static GpsiProfilePutExecutor instance = null;

    private GpsiProfilePutExecutor()
    {
        super(GpsiProfilePutHelper.getInstance());
    }

    public static synchronized GpsiProfilePutExecutor getInstance()
    {
        if(null == instance) {
            instance = new GpsiProfilePutExecutor();
        }
        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        int code = Code.SUCCESS;
        try {
            GpsiProfilePutRequest put_request = request.getRequest().getPutRequest().getGpsiProfilePutRequest();
            String gpsi_profile_id = put_request.getGpsiProfileId();
            GpsiProfile gpsi_profile = createGpsiProfile(put_request);
            List<GpsiprefixProfile>  profilePutList = put_request.getGpsiPrefixPutList();
            List<GpsiprefixProfile>  profileDelList = put_request.getGpsiPrefixDeleteList();
            ClientCacheService.getInstance().getCacheTransactionManager().begin();
            code = ClientCacheService.getInstance().put(Code.GPSIPROFILE_INDICE, gpsi_profile_id, gpsi_profile);
            if (Code.CREATED == code) {
                for(GpsiprefixProfile gpsiprefixProfile : profileDelList) {
                    int retCode = GpsiPrefixProfilesUtil.DelGpsiprefixProfile(gpsiprefixProfile);
                    if (Code.SUCCESS != retCode) {
                        ClientCacheService.getInstance().getCacheTransactionManager().rollback();
                        return new ExecutionResult(retCode);
                    }
                }
                for(GpsiprefixProfile gpsiprefixProfile : profilePutList) {
                    int retCode = GpsiPrefixProfilesUtil.AddGpsiprefixProfile(gpsiprefixProfile);
                    if (Code.CREATED != retCode) {
                        ClientCacheService.getInstance().getCacheTransactionManager().rollback();
                        return new ExecutionResult(retCode);
                    }
                }
            } else {
                logger.error("Rollback for gpsiprofile add failure.");
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

    private GpsiProfile createGpsiProfile(GpsiProfilePutRequest request)
    {
        GpsiProfileIndex index = request.getIndex();

        GpsiProfile gpsi_profile = new GpsiProfile();

        gpsi_profile.setGpsiProfileID(request.getGpsiProfileId());
        gpsi_profile.setData(request.getGpsiProfileData());
        gpsi_profile.setProfileType(request.getIndex().getProfileType());
        gpsi_profile.setGpsiVersion(request.getGpsiVersion());
        List<String> nf_type_list = index.getNfTypeList();
        if(nf_type_list.isEmpty() == false) {
            for(String nf_type : nf_type_list) {
                gpsi_profile.addNFType(nf_type);
            }
        }
        List<String> group_id_list = index.getGroupIndexList();
        if(group_id_list.isEmpty() == false) {
            for(String group_id : group_id_list) {
                gpsi_profile.addGroupID(group_id);
            }
        }

        logger.debug("Gpsi Profile : {}", gpsi_profile.toString());

        return gpsi_profile;
    }
}


