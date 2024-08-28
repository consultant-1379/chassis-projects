package ericsson.core.nrf.dbproxy.helper.cachenfprofile;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.cachenfprofile.CacheNFProfileGetRequestProto.CacheNFProfileGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.cachenfprofile.CacheNFProfileGetResponseProto.CacheNFProfileGetResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.geode.pdx.JSONFormatter;
import org.apache.geode.pdx.PdxInstance;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class CacheNFProfileGetHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(CacheNFProfileGetHelper.class);

  private static CacheNFProfileGetHelper instance;

  private CacheNFProfileGetHelper() {
  }

  public static synchronized CacheNFProfileGetHelper getInstance() {
    if (null == instance) {
      instance = new CacheNFProfileGetHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {
    CacheNFProfileGetRequest request = message.getRequest().getGetRequest()
        .getCacheNfProfileGetRequest();
    String cacheNfInstanceId = request.getCacheNfInstanceId();
    if (cacheNfInstanceId.isEmpty()) {
      LOGGER.error("Empty cacheNfInstanceId is set in CacheNFProfileGetRequest");
      return Code.EMPTY_CACHE_NF_INSTANCE_ID;
    } else if (cacheNfInstanceId.length() > Code.KEY_MAX_LENGTH) {
      LOGGER.error("cacheNfInstanceId length {} is too large, max length is {}",
          cacheNfInstanceId.length(), Code.KEY_MAX_LENGTH);
      return Code.CACHE_NF_INSTANCE_ID_LENGTH_EXCEED_MAX;
    }
    return Code.VALID;
  }

  public NFMessage createResponse(int code) {
    CacheNFProfileGetResponse cacheNfProfileGetResponse = CacheNFProfileGetResponse.newBuilder()
        .setCode(code).build();
    GetResponse getResponse = GetResponse.newBuilder()
        .setCacheNfProfileGetResponse(cacheNfProfileGetResponse).build();
    return createNFMessage(getResponse);
  }

  @Override
  public NFMessage createResponse(ExecutionResult executionResult) {
    if (executionResult.getCode() != Code.SUCCESS) {
      return createResponse(executionResult.getCode());
    } else {
      SearchResult searchResult = (SearchResult) executionResult;
      if (searchResult.getItems().size() > 0) {
        Object obj = searchResult.getItems().get(0);
        String cacheNfProfile = JSONFormatter.toJSON((PdxInstance) obj);
        CacheNFProfileGetResponse cacheNfProfileGetResponse = CacheNFProfileGetResponse
            .newBuilder().setCode(searchResult.getCode()).setCacheNfProfile(cacheNfProfile)
            .build();
        GetResponse getResponse = GetResponse.newBuilder()
            .setCacheNfProfileGetResponse(cacheNfProfileGetResponse).build();
        return createNFMessage(getResponse);
      } else {
        return createResponse(Code.DATA_NOT_EXIST);
      }

    }

  }


}
