package ericsson.core.nrf.dbproxy.executor.groupprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.clientcache.schema.GroupProfile;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.executor.common.ImsiPrefixProfilesUtil;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileIndexProto.GroupProfileIndex;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfilePutRequestProto.GroupProfilePutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.imsiprefixprofile.ImsiprefixProfileProto.ImsiprefixProfile;
import ericsson.core.nrf.dbproxy.helper.groupprofile.GroupProfilePutHelper;
import java.util.List;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.geode.cache.CacheTransactionManager;


public class GroupProfilePutExecutor extends Executor {

  private static final Logger LOGGER = LogManager.getLogger(GroupProfilePutExecutor.class);

  private static GroupProfilePutExecutor instance;

  static {
    instance = null;
  }

  private GroupProfilePutExecutor() {
    super(GroupProfilePutHelper.getInstance());
  }

  public static synchronized GroupProfilePutExecutor getInstance() {
    if (null == instance) {
      instance = new GroupProfilePutExecutor();
    }
    return instance;
  }

  protected ExecutionResult execute(NFMessage request) {
    int code = Code.SUCCESS;
    
    GroupProfilePutRequest putRequest = request.getRequest().getPutRequest()
        .getGroupProfilePutRequest();
    String groupProfileId = putRequest.getGroupProfileId();
    GroupProfile groupProfile = createGroupProfile(putRequest);
    List<ImsiprefixProfile> profilePutList = putRequest.getImsiPrefixPutList();
    List<ImsiprefixProfile> profileDelList = putRequest.getImsiPrefixDeleteList();
    CacheTransactionManager txManager = ClientCacheService.getInstance().getCacheTransactionManager();

    boolean retryTransaction = true;
    int retryTime = 0;
    do {
      try {
        retryTime++;
        if (retryTime > 3) {
          break;
        }

        txManager.begin();
        code = ClientCacheService.getInstance()
            .put(Code.GROUPPROFILE_INDICE, groupProfileId, groupProfile);
        if (Code.CREATED == code) {
          int retCode = ImsiPrefixProfilesUtil.delImsiprefixProfiles(profileDelList);
          if (Code.SUCCESS != retCode) {
            code = retCode;
            continue;
          }
          retCode = ImsiPrefixProfilesUtil.addImsiprefixProfiles(profilePutList);
          if (Code.SUCCESS != retCode) {
            code = retCode;
            continue;
          }
        } else {
          continue;
        }
        txManager.commit();
        retryTransaction = false;  
      } catch (Exception e) {
        LOGGER.error(e.toString());
        code = Code.INTERNAL_ERROR;
      } finally {
        if (txManager.exists()) {
          txManager.rollback();
        }
      }
    } while (retryTransaction);
    
    return new ExecutionResult(code);
  }

  private GroupProfile createGroupProfile(GroupProfilePutRequest request) {
    GroupProfileIndex index = request.getIndex();

    GroupProfile groupProfile = new GroupProfile();

    groupProfile.setGroupProfileID(request.getGroupProfileId());
    groupProfile.setData(request.getGroupProfileData());
    groupProfile.setProfileType(request.getIndex().getProfileType());
    groupProfile.setSupiVersion(request.getSupiVersion());
    List<String> nfTypeList = index.getNfTypeList();
    if (!nfTypeList.isEmpty()) {
      for (String nfType : nfTypeList) {
        groupProfile.addNFType(nfType);
      }
    }
    List<String> groupIdList = index.getGroupIndexList();
    if (!groupIdList.isEmpty()) {
      for (String groupId : groupIdList) {
        groupProfile.addGroupID(groupId);
      }
    }

    LOGGER.debug("Group Profile : {}", groupProfile.toString());

    return groupProfile;
  }
}


