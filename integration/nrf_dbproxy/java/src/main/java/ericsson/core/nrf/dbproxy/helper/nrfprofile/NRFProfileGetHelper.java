package ericsson.core.nrf.dbproxy.helper.nrfprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.clientcache.schema.NRFProfile;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.FragmentResult;
import ericsson.core.nrf.dbproxy.common.FragmentSessionManagement;
import ericsson.core.nrf.dbproxy.common.FragmentUtil;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileGetRequestProto.NRFProfileGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileGetRequestProto.NRFProfileGetRequest.DataCase;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileGetResponseProto;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileGetResponseProto.FragmentNRFProfileInfo;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileGetResponseProto.NRFProfileGetResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileGetResponseProto.NRFProfileInfo;
import ericsson.core.nrf.dbproxy.helper.Helper;
import java.util.ArrayList;
import java.util.List;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NRFProfileGetHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(NRFProfileGetHelper.class);

  private static NRFProfileGetHelper instance;

  private NRFProfileGetHelper() {
  }

  public static synchronized NRFProfileGetHelper getInstance() {
    if (null == instance) {
      instance = new NRFProfileGetHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    NRFProfileGetRequest request = message.getRequest().getGetRequest().getNrfProfileGetRequest();
    DataCase dataCase = request.getDataCase();
    if (dataCase == DataCase.NRF_INSTANCE_ID) {
      String nrfInstanceId = request.getNrfInstanceId();
      if (nrfInstanceId.isEmpty()) {
        LOGGER.error("Empty nrfInstanceId is set in NRFProfileGetRequest");
        return Code.EMPTY_NRF_INSTANCE_ID;
      } else if (nrfInstanceId.length() > Code.KEY_MAX_LENGTH) {
        LOGGER.error("nrfInstanceId length {} is too large, max length is {}",
            nrfInstanceId.length(), Code.KEY_MAX_LENGTH);
        return Code.NRF_INSTANCE_ID_LENGTH_EXCEED_MAX;
      }
    } else if (dataCase == DataCase.FILTER) {
      LOGGER.debug("NRFProfile filter");
    } else if (dataCase == DataCase.FRAGMENT_SESSION_ID) {

      String fragmentSessionId = request.getFragmentSessionId();
      if (fragmentSessionId.isEmpty()) {
        LOGGER.error("Empty fragmentSessionId is set in NRFProfileGetRequest");
        return Code.EMPTY_FRAGMENT_SESSION_ID;
      }
    } else {
      LOGGER.error("Empty NRFProfileGetRequest is received");
      return Code.NFMESSAGE_PROTOCOL_ERROR;
    }

    return Code.VALID;
  }

  public NFMessage createResponse(int code) {

    NRFProfileGetResponse nrfProfileGetResponse = NRFProfileGetResponse.newBuilder()
        .setCode(code).build();
    GetResponse getResponse = GetResponse.newBuilder()
        .setNrfProfileGetResponse(nrfProfileGetResponse).build();
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
              .transmitNumPerTime(fragmentResult, Code.NRFPROFILE_INDICE);
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
          List<NRFProfileInfo> nrfProfiles = getNRFProfileInfo(fragmentResult.getItems());

          String fragmentSessionId = fragmentResult.getFragmentSessionID();
          int totalNumber = fragmentResult.getTotalNumber();
          int transmittedNumber = fragmentResult.getTransmittedNumber();
          FragmentNRFProfileInfo fragmentInfo = FragmentNRFProfileInfo.newBuilder()
              .setFragmentSessionId(fragmentSessionId).setTotalNumber(totalNumber)
              .setTransmittedNumber(transmittedNumber).build();

          NRFProfileGetResponseProto.NRFProfileGetResponse nrfProfileGetResponse = NRFProfileGetResponse
              .newBuilder().setCode(fragmentResult.getCode()).addAllNrfProfile(nrfProfiles)
              .setFragmentNrfprofileInfo(fragmentInfo).build();
          GetResponse getResponse = GetResponse.newBuilder()
              .setNrfProfileGetResponse(nrfProfileGetResponse).build();
          return createNFMessage(getResponse);
        }
      } else {
        int provFlag = 0;
        if (searchResult.getItems().size() == 1) {
          provFlag = (int) getProvFlag(searchResult.getItems().get(0));
        }

        List<NRFProfileInfo> nrfProfiles = getNRFProfileInfo(searchResult.getItems());
        NRFProfileGetResponse nrfProfileGetResponse = NRFProfileGetResponse.newBuilder()
            .setCode(searchResult.getCode()).addAllNrfProfile(nrfProfiles).setProvFlag(provFlag)
            .build();
        GetResponse getResponse = GetResponse.newBuilder()
            .setNrfProfileGetResponse(nrfProfileGetResponse).build();
        return createNFMessage(getResponse);
      }
    }

  }

  private long getProvFlag(Object obj) {
    NRFProfile nrfProfile = (NRFProfile) obj;
    long provFlag = nrfProfile.getProvFlag();

    return provFlag;
  }

  private List<NRFProfileInfo> getNRFProfileInfo(List<Object> items) {

    List<NRFProfileInfo> nrfProfiles = new ArrayList<>();
    for (Object obj : items) {
      NRFProfile nrfProfile = (NRFProfile) obj;
      ByteString data = nrfProfile.getRaw_data();
      long expiredTime = nrfProfile.getExpireTime();

      NRFProfileInfo item = NRFProfileInfo.newBuilder().setRawNrfProfile(data)
          .setExpiredTime(expiredTime).build();
      nrfProfiles.add(item);
    }

    return nrfProfiles;
  }


}
