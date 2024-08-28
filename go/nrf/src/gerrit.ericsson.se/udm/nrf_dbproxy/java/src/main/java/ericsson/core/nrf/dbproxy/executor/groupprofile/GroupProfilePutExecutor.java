package ericsson.core.nrf.dbproxy.executor.groupprofile;

import java.util.List;

import ericsson.core.nrf.dbproxy.executor.Executor;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.helper.groupprofile.GroupProfilePutHelper;
import ericsson.core.nrf.dbproxy.clientcache.schema.GroupProfile;
import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfilePutRequestProto.GroupProfilePutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileIndexProto.GroupProfileIndex;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.imsiprefixprofile.ImsiprefixProfileProto.ImsiprefixProfile;
import ericsson.core.nrf.dbproxy.executor.common.ImsiPrefixProfilesUtil;


public class GroupProfilePutExecutor extends Executor
{
    private static final Logger logger = LogManager.getLogger(GroupProfilePutExecutor.class);

    private static GroupProfilePutExecutor instance = null;

    private GroupProfilePutExecutor()
    {
        super(GroupProfilePutHelper.getInstance());
    }

    public static synchronized GroupProfilePutExecutor getInstance()
    {
        if(null == instance) {
            instance = new GroupProfilePutExecutor();
        }
        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        int code = Code.SUCCESS;
        try {
            GroupProfilePutRequest put_request = request.getRequest().getPutRequest().getGroupProfilePutRequest();
            String group_profile_id = put_request.getGroupProfileId();
            GroupProfile group_profile = createGroupProfile(put_request);
            List<ImsiprefixProfile>  profilePutList = put_request.getImsiPrefixPutList();
            List<ImsiprefixProfile>  profileDelList = put_request.getImsiPrefixDeleteList();
            ClientCacheService.getInstance().getCacheTransactionManager().begin();
            code = ClientCacheService.getInstance().put(Code.GROUPPROFILE_INDICE, group_profile_id, group_profile);
            if (Code.CREATED == code) {
                for(ImsiprefixProfile imsiprefixProfile : profileDelList) {
                    int retCode = ImsiPrefixProfilesUtil.DelImsiprefixProfile(imsiprefixProfile);
                    if (Code.SUCCESS != retCode) {
                        ClientCacheService.getInstance().getCacheTransactionManager().rollback();
                        return new ExecutionResult(retCode);
                    }
                }
                for(ImsiprefixProfile imsiprefixProfile : profilePutList) {
                    int retCode = ImsiPrefixProfilesUtil.AddImsiprefixProfile(imsiprefixProfile);
                    if (Code.CREATED != retCode) {
                        ClientCacheService.getInstance().getCacheTransactionManager().rollback();
                        return new ExecutionResult(retCode);
                    }
                }
            } else {
                logger.error("Rollback for groupprofile add failure.");
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

    private GroupProfile createGroupProfile(GroupProfilePutRequest request)
    {
        GroupProfileIndex index = request.getIndex();

        GroupProfile group_profile = new GroupProfile();

        group_profile.setGroupProfileID(request.getGroupProfileId());
        group_profile.setData(request.getGroupProfileData());
        group_profile.setProfileType(request.getIndex().getProfileType());
        group_profile.setSupiVersion(request.getSupiVersion());
        List<String> nf_type_list = index.getNfTypeList();
        if(nf_type_list.isEmpty() == false) {
            for(String nf_type : nf_type_list) {
                group_profile.addNFType(nf_type);
            }
        }
        List<String> group_id_list = index.getGroupIndexList();
        if(group_id_list.isEmpty() == false) {
            for(String group_id : group_id_list) {
                group_profile.addGroupID(group_id);
            }
        }

        logger.debug("Group Profile : {}", group_profile.toString());

        return group_profile;
    }
}


