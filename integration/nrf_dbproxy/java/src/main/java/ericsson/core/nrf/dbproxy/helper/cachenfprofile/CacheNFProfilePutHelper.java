package ericsson.core.nrf.dbproxy.helper.cachenfprofile;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PutResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.cachenfprofile.CacheNFProfilePutRequestProto.CacheNFProfilePutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.cachenfprofile.CacheNFProfilePutResponseProto.CacheNFProfilePutResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class CacheNFProfilePutHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(CacheNFProfilePutHelper.class);

  private static CacheNFProfilePutHelper instance;

  private CacheNFProfilePutHelper() {
  }

  public static synchronized CacheNFProfilePutHelper getInstance() {
    if (null == instance) {
      instance = new CacheNFProfilePutHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    CacheNFProfilePutRequest request = message.getRequest().getPutRequest()
        .getCacheNfProfilePutRequest();
    String cacheNfInstanceId = request.getCacheNfInstanceId();
    if (cacheNfInstanceId.isEmpty()) {
      LOGGER.error("cache_f_instance_id field is empty in NRFProfilePutRequest");
      return Code.EMPTY_CACHE_NF_INSTANCE_ID;
    }

    if (cacheNfInstanceId.length() > Code.KEY_MAX_LENGTH) {
      LOGGER.error("cacheNfInstanceId length {} is too large, max length is {}",
          cacheNfInstanceId.length(), Code.KEY_MAX_LENGTH);
      return Code.CACHE_NF_INSTANCE_ID_LENGTH_EXCEED_MAX;
    }

    String rawCacheNfProfile = request.getRawCacheNfProfile();
    if (rawCacheNfProfile.isEmpty()) {
      LOGGER.error("rawCacheNfProfile field is empty in CacheNFProfilePutRequest");
      return Code.EMPTY_RAW_CACHE_NF_PROFILE;
    }

    return Code.VALID;
  }

  public NFMessage createResponse(int code) {

    CacheNFProfilePutResponse cacheNfProfilePutResponse = CacheNFProfilePutResponse.newBuilder()
        .setCode(code).build();
    PutResponse putResponse = PutResponse.newBuilder()
        .setCacheNfProfilePutResponse(cacheNfProfilePutResponse).build();
    return createNFMessage(putResponse);
  }
}
