package ericsson.core.nrf.dbproxy.executor.gpsiprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.executor.common.GpsiPrefixProfilesUtil;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileProto.GpsiprefixProfile;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileDelRequestProto.GpsiProfileDelRequest;
import ericsson.core.nrf.dbproxy.helper.gpsiprofile.GpsiProfileDelHelper;
import java.util.List;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class GpsiProfileDeleteExecutor extends Executor {

  private static final Logger LOGGER = LogManager.getLogger(GpsiProfileDeleteExecutor.class);

  private static GpsiProfileDeleteExecutor instance;

  static {
    instance = null;
  }

  private GpsiProfileDeleteExecutor() {
    super(GpsiProfileDelHelper.getInstance());
  }

  public static synchronized GpsiProfileDeleteExecutor getInstance() {
    if (null == instance) {
      instance = new GpsiProfileDeleteExecutor();
    }
    return instance;
  }

  protected ExecutionResult execute(NFMessage request) {
    int code = Code.SUCCESS;
    try {
      GpsiProfileDelRequest delRequest = request.getRequest().getDelRequest()
          .getGpsiProfileDelRequest();
      String gpsiProfileId = delRequest.getGpsiProfileId();
      List<GpsiprefixProfile> profileDelList = delRequest.getGpsiPrefixDeleteList();

      ClientCacheService.getInstance().getCacheTransactionManager().begin();
      for (GpsiprefixProfile gpsiprefixProfile : profileDelList) {
        int retCode = GpsiPrefixProfilesUtil.delGpsiprefixProfile(gpsiprefixProfile);
        if (Code.SUCCESS != retCode) {
          ClientCacheService.getInstance().getCacheTransactionManager().rollback();
          return new ExecutionResult(retCode);
        }
      }
      code = ClientCacheService.getInstance()
          .delete(Code.GPSIPROFILE_INDICE, gpsiProfileId, false);
      if (Code.SUCCESS != code) {
        ClientCacheService.getInstance().getCacheTransactionManager().rollback();
        return new ExecutionResult(code);
      }
      ClientCacheService.getInstance().getCacheTransactionManager().commit();
    } catch (Exception e) {
      LOGGER.error(e.toString());
      code = Code.INTERNAL_ERROR;
    }
    return new ExecutionResult(code);
  }

}
