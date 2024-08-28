package ericsson.core.nrf.dbproxy.executor.cachenfprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.functionservice.RemoteCacheMonitorThread;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.cachenfprofile.CacheNFProfilePutRequestProto.CacheNFProfilePutRequest;
import ericsson.core.nrf.dbproxy.helper.cachenfprofile.CacheNFProfilePutHelper;
import org.apache.geode.pdx.JSONFormatter;
import org.apache.geode.pdx.PdxInstance;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class CacheNFProfilePutExecutor extends Executor {

  private static final Logger LOGGER = LogManager.getLogger(CacheNFProfilePutExecutor.class);

  private static CacheNFProfilePutExecutor instance;

  static {
    instance = null;
  }

  private CacheNFProfilePutExecutor() {
    super(CacheNFProfilePutHelper.getInstance());
  }

  public static synchronized CacheNFProfilePutExecutor getInstance() {
    if (null == instance) {
      instance = new CacheNFProfilePutExecutor();
    }

    return instance;
  }

  protected ExecutionResult execute(NFMessage request) {
    CacheNFProfilePutRequest putRequest = request.getRequest().getPutRequest()
        .getCacheNfProfilePutRequest();
    String cacheNfInstanceId = putRequest.getCacheNfInstanceId();
    String cacheNfProfile = putRequest.getRawCacheNfProfile();
    int code = Code.CREATED;
    try {
      PdxInstance pdxCacheNfProfile = JSONFormatter.fromJSON(cacheNfProfile.toString());
      ClientCacheService.getInstance()
          .put(Code.CACHENFPROFILE_INDICE, cacheNfInstanceId, pdxCacheNfProfile);
      RemoteCacheMonitorThread.getInstance().incCacheOperationCount();
    } catch (Exception e) {
      LOGGER.error(e.toString());
      code = Code.CACHE_NF_PROFILE_FORMAT_ERROR;
    }
    return new ExecutionResult(code);
  }

}
