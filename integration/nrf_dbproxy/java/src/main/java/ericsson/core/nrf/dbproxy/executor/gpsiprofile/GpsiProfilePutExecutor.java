package ericsson.core.nrf.dbproxy.executor.gpsiprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.clientcache.schema.GpsiProfile;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.executor.common.GpsiPrefixProfilesUtil;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileProto.GpsiprefixProfile;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileIndexProto.GpsiProfileIndex;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfilePutRequestProto.GpsiProfilePutRequest;
import ericsson.core.nrf.dbproxy.helper.gpsiprofile.GpsiProfilePutHelper;
import java.util.List;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.geode.cache.CacheTransactionManager;


public class GpsiProfilePutExecutor extends Executor {

  private static final Logger LOGGER = LogManager.getLogger(GpsiProfilePutExecutor.class);

  private static GpsiProfilePutExecutor instance;

  static {
    instance = null;
  }

  private GpsiProfilePutExecutor() {
    super(GpsiProfilePutHelper.getInstance());
  }

  public static synchronized GpsiProfilePutExecutor getInstance() {
    if (null == instance) {
      instance = new GpsiProfilePutExecutor();
    }
    return instance;
  }

  protected ExecutionResult execute(NFMessage request) {
    int code = Code.SUCCESS;

    GpsiProfilePutRequest putRequest = request.getRequest().getPutRequest()
        .getGpsiProfilePutRequest();
    String gpsiProfileId = putRequest.getGpsiProfileId();
    GpsiProfile gpsiProfile = createGpsiProfile(putRequest);
    List<GpsiprefixProfile> profilePutList = putRequest.getGpsiPrefixPutList();
    List<GpsiprefixProfile> profileDelList = putRequest.getGpsiPrefixDeleteList();
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
            .put(Code.GPSIPROFILE_INDICE, gpsiProfileId, gpsiProfile);
        if (Code.CREATED == code) {
          int retCode = GpsiPrefixProfilesUtil.delGpsiprefixProfiles(profileDelList);
          if (Code.SUCCESS != retCode) {
            code = retCode;
            continue;
          }
          retCode = GpsiPrefixProfilesUtil.addGpsiprefixProfiles(profilePutList);
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

  private GpsiProfile createGpsiProfile(GpsiProfilePutRequest request) {
    GpsiProfileIndex index = request.getIndex();

    GpsiProfile gpsiProfile = new GpsiProfile();

    gpsiProfile.setGpsiProfileID(request.getGpsiProfileId());
    gpsiProfile.setData(request.getGpsiProfileData());
    gpsiProfile.setProfileType(request.getIndex().getProfileType());
    gpsiProfile.setGpsiVersion(request.getGpsiVersion());
    List<String> nfTypeList = index.getNfTypeList();
    if (!nfTypeList.isEmpty()) {
      for (String nfType : nfTypeList) {
        gpsiProfile.addNFType(nfType);
      }
    }
    List<String> groupIdList = index.getGroupIndexList();
    if (!groupIdList.isEmpty()) {
      for (String groupId : groupIdList) {
        gpsiProfile.addGroupID(groupId);
      }
    }

    LOGGER.debug("Gpsi Profile : {}", gpsiProfile.toString());

    return gpsiProfile;
  }
}


