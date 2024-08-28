package ericsson.core.nrf.dbproxy.helper.gpsiprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.clientcache.schema.GpsiProfile;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.FragmentResult;
import ericsson.core.nrf.dbproxy.common.FragmentSessionManagement;
import ericsson.core.nrf.dbproxy.common.FragmentUtil;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileGetRequestProto.GpsiProfileGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileGetRequestProto.GpsiProfileGetRequest.DataCase;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileGetResponseProto.GpsiProfileGetResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileGetResponseProto.GpsiProfileInfo;
import ericsson.core.nrf.dbproxy.helper.Helper;
import java.util.ArrayList;
import java.util.List;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class GpsiProfileGetHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(GpsiProfileGetHelper.class);

  private static GpsiProfileGetHelper instance;

  private GpsiProfileGetHelper() {
  }

  public static synchronized GpsiProfileGetHelper getInstance() {
    if (null == instance) {
      instance = new GpsiProfileGetHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    GpsiProfileGetRequest request = message.getRequest().getGetRequest().getGpsiProfileGetRequest();
    DataCase dataCase = request.getDataCase();
    if (dataCase == DataCase.GPSI_PROFILE_ID) {

      String gpsiProfileId = request.getGpsiProfileId();
      if (gpsiProfileId.isEmpty()) {
        LOGGER.error("Empty gpsiProfileId is set in GpsiProfileGetRequest");
        return Code.EMPTY_GPSI_PROFILE_ID;
      } else if (gpsiProfileId.length() > Code.KEY_MAX_LENGTH) {
        LOGGER.error("gpsiProfileId length {} is too large, max length is {}",
            gpsiProfileId.length(), Code.KEY_MAX_LENGTH);
        return Code.GPSI_PROFILE_ID_LENGTH_EXCEED_MAX;
      }
    } else if (dataCase == DataCase.FILTER) {

      if (request.getFilter().getIndex().getNfTypeList().isEmpty() &&
          request.getFilter().getIndex().getGroupIndexList().isEmpty() &&
          request.getFilter().getIndex().getProfileType() == Code.PROFILE_TYPE_EMPTY) {
        LOGGER.error("Empty GpsiProfileFilter is set in filter of GpsiProfileGetRequest");
        return Code.EMPTY_GPSI_PROFILE_FILTER;
      }
    } else if (dataCase == DataCase.FRAGMENT_SESSION_ID) {

      String fragmentSessionId = request.getFragmentSessionId();
      if (fragmentSessionId.isEmpty()) {
        LOGGER.error("Empty fragmentSessionId is set in NFProfileGetRequest");
        return Code.EMPTY_FRAGMENT_SESSION_ID;
      }
    } else {
      LOGGER.error("Empty GpsiProfileGetRequest is received");
      return Code.NFMESSAGE_PROTOCOL_ERROR;
    }

    return Code.VALID;
  }

  public NFMessage createResponse(int code) {

    GpsiProfileGetResponse gpsiProfileGetResponse = GpsiProfileGetResponse.newBuilder()
        .setCode(code).build();
    GetResponse getResponse = GetResponse.newBuilder()
        .setGpsiProfileGetResponse(gpsiProfileGetResponse).build();
    return createNFMessage(getResponse);
  }

  @Override
  public NFMessage createResponse(ExecutionResult executionResult) {
    if (executionResult.getCode() != Code.SUCCESS) {
      return createResponse(executionResult.getCode());
    } else {
      SearchResult searchResult = (SearchResult) executionResult;
      if (searchResult.isFragmented()) {
        FragmentResult fragmentResult = (FragmentResult) searchResult;
        if (fragmentResult.getFragmentSessionID().isEmpty()) {
          int firstTransmitNum = FragmentUtil
              .transmitNumPerTime(fragmentResult, Code.GPSIPROFILE_INDICE);
          if (FragmentSessionManagement.getInstance().put(fragmentResult, firstTransmitNum)) {
            FragmentResult item = new FragmentResult();
            item.addAll(fragmentResult.getItems().subList(0, firstTransmitNum));
            item.setFragmentSessionID(fragmentResult.getFragmentSessionID());
            item.setTotalNumber(fragmentResult.getTotalNumber());
            item.setTransmittedNumber(fragmentResult.getTransmittedNumber());
            return createResponse(item);
          } else {
            return createResponse(Code.INTERNAL_ERROR);
          }
        } else {
          List<GpsiProfileInfo> gpsiProfileInfo = getGpsiProfileInfo(fragmentResult.getItems());

          String fragmentSessionId = fragmentResult.getFragmentSessionID();
          int totalNumber = fragmentResult.getTotalNumber();
          int transmittedNumber = fragmentResult.getTransmittedNumber();
          FragmentInfo fragmentInfo = FragmentInfo.newBuilder()
              .setFragmentSessionId(fragmentSessionId).setTotalNumber(totalNumber)
              .setTransmittedNumber(transmittedNumber).build();
          GpsiProfileGetResponse gpsiProfileGetResponse = GpsiProfileGetResponse.newBuilder()
              .setCode(searchResult.getCode()).addAllGpsiProfileInfo(gpsiProfileInfo)
              .setFragmentInfo(fragmentInfo).build();
          GetResponse getResponse = GetResponse.newBuilder()
              .setGpsiProfileGetResponse(gpsiProfileGetResponse).build();
          return createNFMessage(getResponse);
        }
      } else {
        List<GpsiProfileInfo> gpsiProfileInfo = getGpsiProfileInfo(searchResult.getItems());
        GpsiProfileGetResponse gpsiProfileGetResponse = GpsiProfileGetResponse.newBuilder()
            .setCode(searchResult.getCode()).addAllGpsiProfileInfo(gpsiProfileInfo).build();
        GetResponse getResponse = GetResponse.newBuilder()
            .setGpsiProfileGetResponse(gpsiProfileGetResponse).build();
        return createNFMessage(getResponse);
      }
    }
  }

  private List<GpsiProfileInfo> getGpsiProfileInfo(List<Object> items) {

    List<GpsiProfileInfo> gpsiProfileInfoList = new ArrayList<>();
    for (Object obj : items) {
      GpsiProfile item = (GpsiProfile) obj;
      ByteString data = item.getData();
      String gpsiProfileId = item.getGpsiProfileID();
      Long gpsiVersion = item.getGpsiVersion();
      GpsiProfileInfo gpsiProfileInfo = GpsiProfileInfo.newBuilder()
          .setGpsiProfileId(gpsiProfileId).setGpsiVersion(gpsiVersion).setGpsiProfileData(data)
          .build();
      gpsiProfileInfoList.add(gpsiProfileInfo);
    }
    return gpsiProfileInfoList;
  }
}