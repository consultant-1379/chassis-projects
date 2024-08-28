package ericsson.core.nrf.dbproxy.executor.nfprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfilePutRequestProto.NFProfilePutRequest;
import ericsson.core.nrf.dbproxy.helper.nfprofile.NFProfilePutHelper;
import org.apache.geode.cache.CacheTransactionManager;
import org.apache.geode.cache.CommitConflictException;
import org.apache.geode.pdx.JSONFormatter;
import org.apache.geode.pdx.PdxInstance;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NFProfilePutExecutor extends Executor {

  private static final Logger LOGGER = LogManager.getLogger(NFProfilePutExecutor.class);

  private static NFProfilePutExecutor instance;

  static {
    instance = null;
  }

  private NFProfilePutExecutor() {
    super(NFProfilePutHelper.getInstance());
  }

  public static synchronized NFProfilePutExecutor getInstance() {
    if (null == instance) {
      instance = new NFProfilePutExecutor();
    }
    return instance;
  }

  protected ExecutionResult execute(NFMessage request) {
    NFProfilePutRequest putRequest = request.getRequest().getPutRequest().getNfProfilePutRequest();
    String nfInstanceId = putRequest.getNfInstanceId();
    String nfProfile = putRequest.getNfProfile();
    String nfHelper = putRequest.getNfHelperInfo();
    int code;
    if (nfHelper != null && nfHelper.length() > 0) {
      code = Code.INTERNAL_ERROR;
      CacheTransactionManager txManager = ClientCacheService.getInstance()
          .getCacheTransactionManager();
      boolean retryTransaction = false;
      int retryTime = 0;
      PdxInstance pdxInstanceNfprofile = JSONFormatter.fromJSON(nfProfile);
      PdxInstance pdxInstanceNfhelper = JSONFormatter.fromJSON(nfHelper);
      do {
        try {
          retryTime++;
          if (retryTime > 3) {
            break;
          }
          txManager.begin();
          ClientCacheService.getInstance()
              .put(Code.NFPROFILE_INDICE, nfInstanceId, pdxInstanceNfprofile);
          ClientCacheService.getInstance()
              .put(Code.NFHELPER_INDICE, nfInstanceId, pdxInstanceNfhelper);
          txManager.commit();
          retryTransaction = false;
          code = Code.CREATED;
        } catch (CommitConflictException conflictException) {
          LOGGER.error(conflictException.toString());
          retryTransaction = true;
        } catch (Exception e) {
          LOGGER.error(e.toString());
        } finally {
          if (txManager.exists()) {
            txManager.rollback();
          }
        }
      } while (retryTransaction);

    } else {
      try {
        PdxInstance pdxInstance = JSONFormatter.fromJSON(nfProfile);
        code = ClientCacheService.getInstance()
            .put(Code.NFPROFILE_INDICE, nfInstanceId, pdxInstance);
      } catch (Exception e) {
        LOGGER.error(e.toString());
        code = Code.NF_PROFILE_FORMAT_ERROR;
      }
    }

    return new ExecutionResult(code);
  }
}
