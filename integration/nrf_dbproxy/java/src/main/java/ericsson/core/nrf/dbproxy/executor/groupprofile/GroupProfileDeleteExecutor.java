package ericsson.core.nrf.dbproxy.executor.groupprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.executor.common.ImsiPrefixProfilesUtil;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileDelRequestProto.GroupProfileDelRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.imsiprefixprofile.ImsiprefixProfileProto.ImsiprefixProfile;
import ericsson.core.nrf.dbproxy.helper.groupprofile.GroupProfileDelHelper;
import java.util.List;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class GroupProfileDeleteExecutor extends Executor {

  private static final Logger LOGGER = LogManager.getLogger(GroupProfileDeleteExecutor.class);

  private static GroupProfileDeleteExecutor instance;

  static {
    instance = null;
  }

  private GroupProfileDeleteExecutor() {
    super(GroupProfileDelHelper.getInstance());
  }

  public static synchronized GroupProfileDeleteExecutor getInstance() {
    if (null == instance) {
      instance = new GroupProfileDeleteExecutor();
    }
    return instance;
  }

  protected ExecutionResult execute(NFMessage request) {
    int code = Code.SUCCESS;
    try {
      GroupProfileDelRequest delRequest = request.getRequest().getDelRequest()
          .getGroupProfileDelRequest();
      String groupProfileId = delRequest.getGroupProfileId();
      List<ImsiprefixProfile> profileDelList = delRequest.getImsiPrefixDeleteList();
      ClientCacheService.getInstance().getCacheTransactionManager().begin();
      for (ImsiprefixProfile imsiprefixProfile : profileDelList) {
        int retCode = ImsiPrefixProfilesUtil.delImsiprefixProfile(imsiprefixProfile);
        if (Code.SUCCESS != retCode) {
          ClientCacheService.getInstance().getCacheTransactionManager().rollback();
          return new ExecutionResult(retCode);
        }
      }
      code = ClientCacheService.getInstance()
          .delete(Code.GROUPPROFILE_INDICE, groupProfileId, false);
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
