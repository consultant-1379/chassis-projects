package ericsson.core.nrf.dbproxy.helper.imsiprefixprofile;

import ericsson.core.nrf.dbproxy.clientcache.schema.ImsiprefixProfiles;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.imsiprefixprofile.ImsiprefixProfileGetRequestProto.ImsiprefixProfileGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.imsiprefixprofile.ImsiprefixProfileGetResponseProto.ImsiprefixProfileGetResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class ImsiprefixProfileGetHelper extends Helper {

  public static final long SEARCH_IMSI_MAX_VALUE = 999999999999999L;

  private static final Logger LOGGER = LogManager.getLogger(ImsiprefixProfileGetHelper.class);

  private static ImsiprefixProfileGetHelper instance;

  private ImsiprefixProfileGetHelper() {
  }

  public static synchronized ImsiprefixProfileGetHelper getInstance() {
    if (null == instance) {
      instance = new ImsiprefixProfileGetHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    ImsiprefixProfileGetRequest request = message.getRequest().getGetRequest()
        .getImsiprefixProfileGetRequest();

    Long searchImsi = request.getSearchImsi();
    if (searchImsi == 0L) {
      LOGGER.error("value 0 searchImsi is set in ImsiperfixProfileGetRequest");
      return Code.EMPTY_SEARCH_IMSI;
    } else if (searchImsi > SEARCH_IMSI_MAX_VALUE) {
      LOGGER.error("SEARCH_IMSI_MAX_VALUE searchImsi is set in ImsiperfixProfileGetRequest");
      return Code.SEARCH_IMSI_LENGTH_EXCEED_MAX;
    }
    return Code.VALID;
  }

  public NFMessage createResponse(int code) {
    ImsiprefixProfileGetResponse nrfAddressGetResponse = ImsiprefixProfileGetResponse
        .newBuilder().setCode(code).build();
    GetResponse getResponse = GetResponse.newBuilder()
        .setImsiprefixProfileGetResponse(nrfAddressGetResponse).build();
    return createNFMessage(getResponse);
  }

  @Override
  public NFMessage createResponse(ExecutionResult executionResult) {
    if (executionResult.getCode() != Code.SUCCESS) {
      return createResponse(executionResult.getCode());
    } else {
      List<String> valueInfoList = new ArrayList<>();
      SearchResult searchResult = (SearchResult) executionResult;
      for (Object obj : searchResult.getItems()) {
        Map<Long, ImsiprefixProfiles> searchResultMap = (Map<Long, ImsiprefixProfiles>) obj;
        for (Map.Entry<Long, ImsiprefixProfiles> entry : searchResultMap.entrySet()) {
          Long key = entry.getKey();
          if (searchResultMap.get(key) != null) {
            ImsiprefixProfiles imsiprefixProfiles = searchResultMap.get(key);
            Iterator iter = imsiprefixProfiles.getValueInfo().keySet().iterator();
            while (iter.hasNext()) {
              String value = (String) (iter.next());
              valueInfoList.add(value);
            }
          }
        }
      }
      if (valueInfoList.isEmpty()) {
        return createResponse(Code.DATA_NOT_EXIST);
      }
      ImsiprefixProfileGetResponse imsiprefixProfileGetResponse = ImsiprefixProfileGetResponse
          .newBuilder().setCode(executionResult.getCode()).addAllValueInfo(valueInfoList)
          .build();
      GetResponse getResponse = GetResponse.newBuilder()
          .setImsiprefixProfileGetResponse(imsiprefixProfileGetResponse).build();
      return createNFMessage(getResponse);
    }

  }
}
