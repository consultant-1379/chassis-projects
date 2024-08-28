package ericsson.core.nrf.dbproxy.helper.gpsiprefixprofile;

import ericsson.core.nrf.dbproxy.clientcache.schema.GpsiprefixProfiles;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetResponseProto.GpsiprefixProfileGetResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class GpsiprefixProfileGetHelper extends Helper {

  public static final long SEARCH_GPSI_MAX_VALUE = 999999999999999L;

  private static final Logger LOGGER = LogManager.getLogger(GpsiprefixProfileGetHelper.class);

  private static GpsiprefixProfileGetHelper instance;

  private GpsiprefixProfileGetHelper() {
  }

  public static synchronized GpsiprefixProfileGetHelper getInstance() {
    if (null == instance) {
      instance = new GpsiprefixProfileGetHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    GpsiprefixProfileGetRequest request = message.getRequest().getGetRequest()
        .getGpsiprefixProfileGetRequest();

    Long searchGpsi = request.getSearchGpsi();
    if (searchGpsi == 0L) {
      LOGGER.error("value 0 searchGpsi is set in GpsiperfixProfileGetRequest");
      return Code.EMPTY_SEARCH_GPSI;
    } else if (searchGpsi > SEARCH_GPSI_MAX_VALUE) {
      LOGGER.error("SEARCH_GPSI_MAX_VALUE searchGpsi is set in GpsiperfixProfileGetRequest");
      return Code.SEARCH_GPSI_LENGTH_EXCEED_MAX;
    }
    return Code.VALID;
  }

  public NFMessage createResponse(int code) {
    GpsiprefixProfileGetResponse nrfAddressGetResponse = GpsiprefixProfileGetResponse
        .newBuilder().setCode(code).build();
    GetResponse getResponse = GetResponse.newBuilder()
        .setGpsiprefixProfileGetResponse(nrfAddressGetResponse).build();
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
        Map<Long, GpsiprefixProfiles> searchResultMap = (Map<Long, GpsiprefixProfiles>) obj;
        for (Map.Entry<Long, GpsiprefixProfiles> entry : searchResultMap.entrySet()) {
          Long key = entry.getKey();
          if (searchResultMap.get(key) != null) {
            GpsiprefixProfiles gpsiprefixProfiles = searchResultMap.get(key);
            Iterator iter = gpsiprefixProfiles.getValueInfo().keySet().iterator();
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
      GpsiprefixProfileGetResponse gpsiprefixProfileGetResponse = GpsiprefixProfileGetResponse
          .newBuilder().setCode(executionResult.getCode()).addAllValueInfo(valueInfoList)
          .build();
      GetResponse getResponse = GetResponse.newBuilder()
          .setGpsiprefixProfileGetResponse(gpsiprefixProfileGetResponse).build();
      return createNFMessage(getResponse);
    }

  }
}
