package ericsson.core.nrf.dbproxy.helper.nrfaddress;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.clientcache.schema.NRFAddress;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressGetRequestProto.NRFAddressGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressGetRequestProto.NRFAddressGetRequest.DataCase;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressGetResponseProto.NRFAddressGetResponse;
import ericsson.core.nrf.dbproxy.helper.Helper;
import java.util.ArrayList;
import java.util.List;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NRFAddressGetHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(NRFAddressGetHelper.class);

  private static NRFAddressGetHelper instance;

  private NRFAddressGetHelper() {
  }

  public static synchronized NRFAddressGetHelper getInstance() {
    if (null == instance) {
      instance = new NRFAddressGetHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    NRFAddressGetRequest request = message.getRequest().getGetRequest().getNrfAddressGetRequest();
    DataCase dataCase = request.getDataCase();
    if (dataCase == DataCase.NRF_ADDRESS_ID) {

      String nrfAddressId = request.getNrfAddressId();
      if (nrfAddressId.isEmpty()) {
        LOGGER.error("Empty nrfAddressId is set in NRFAddressGetRequest");
        return Code.EMPTY_NRF_ADDRESS_ID;
      } else if (nrfAddressId.length() > Code.KEY_MAX_LENGTH) {
        LOGGER.error("nrfAddressId length {} is too large, max length is {}",
            nrfAddressId.length(), Code.KEY_MAX_LENGTH);
        return Code.NRF_ADDRESS_ID_LENGTH_EXCEED_MAX;
      }
    } else if (dataCase == DataCase.FILTER) {

      if (request.getFilter().getIndex().getNrfAddressKey1List().isEmpty() &&
          request.getFilter().getIndex().getNrfAddressKey2().isEmpty() &&
          request.getFilter().getIndex().getNrfAddressKey3().isEmpty() &&
          request.getFilter().getIndex().getNrfAddressKey4().isEmpty() &&
          request.getFilter().getIndex().getNrfAddressKey5().isEmpty()) {
        LOGGER.error("Empty NRFAddressFilter is set in filter of NRFAddressGetRequest");
        return Code.EMPTY_NRF_ADDRESS_FILTER;
      }
    } else {
      LOGGER.error("Empty NRFAddressGetRequest is received");
      return Code.NFMESSAGE_PROTOCOL_ERROR;
    }

    return Code.VALID;
  }

  public NFMessage createResponse(int code) {

    NRFAddressGetResponse nrfAddressGetResponse = NRFAddressGetResponse.newBuilder()
        .setCode(code).build();
    GetResponse getResponse = GetResponse.newBuilder()
        .setNrfAddressGetResponse(nrfAddressGetResponse).build();
    return createNFMessage(getResponse);
  }

  @Override
  public NFMessage createResponse(ExecutionResult executionResult) {

    if (executionResult.getCode() != Code.SUCCESS) {
      return createResponse(executionResult.getCode());
    } else {
      SearchResult searchResult = (SearchResult) executionResult;

      List<ByteString> nrfAddresses = new ArrayList<>();
      for (Object obj : searchResult.getItems()) {
        NRFAddress item = (NRFAddress) obj;
        nrfAddresses.add(item.getData());
      }
      NRFAddressGetResponse nrfAddressGetResponse = NRFAddressGetResponse.newBuilder()
          .setCode(executionResult.getCode()).addAllNrfAddressData(nrfAddresses).build();
      GetResponse getResponse = GetResponse.newBuilder()
          .setNrfAddressGetResponse(nrfAddressGetResponse).build();
      return createNFMessage(getResponse);
    }

  }
}
