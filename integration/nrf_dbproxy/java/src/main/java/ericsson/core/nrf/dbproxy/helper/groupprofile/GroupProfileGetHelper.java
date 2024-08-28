package ericsson.core.nrf.dbproxy.helper.groupprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.clientcache.schema.GroupProfile;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.FragmentResult;
import ericsson.core.nrf.dbproxy.common.FragmentSessionManagement;
import ericsson.core.nrf.dbproxy.common.FragmentUtil;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileGetRequestProto.GroupProfileGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileGetRequestProto.GroupProfileGetRequest.DataCase;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileGetResponseProto.GroupProfileGetResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileGetResponseProto.GroupProfileInfo;
import ericsson.core.nrf.dbproxy.helper.Helper;
import java.util.ArrayList;
import java.util.List;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class GroupProfileGetHelper extends Helper {

  private static final Logger LOGGER = LogManager.getLogger(GroupProfileGetHelper.class);

  private static GroupProfileGetHelper instance;

  private GroupProfileGetHelper() {
  }

  public static synchronized GroupProfileGetHelper getInstance() {
    if (null == instance) {
      instance = new GroupProfileGetHelper();
    }
    return instance;
  }

  public int validate(NFMessage message) {

    GroupProfileGetRequest request = message.getRequest().getGetRequest()
        .getGroupProfileGetRequest();
    DataCase dataCase = request.getDataCase();
    if (dataCase == DataCase.GROUP_PROFILE_ID) {

      String groupProfileId = request.getGroupProfileId();
      if (groupProfileId.isEmpty()) {
        LOGGER.error("Empty group_profile_id is set in GroupProfileGetRequest");
        return Code.EMPTY_GROUP_PROFILE_ID;
      } else if (groupProfileId.length() > Code.KEY_MAX_LENGTH) {
        LOGGER.error("group_profile_id length {} is too large, max length is {}",
            groupProfileId.length(), Code.KEY_MAX_LENGTH);
        return Code.GROUP_PROFILE_ID_LENGTH_EXCEED_MAX;
      }
    } else if (dataCase == DataCase.FILTER) {

      if (request.getFilter().getIndex().getNfTypeList().isEmpty() &&
          request.getFilter().getIndex().getGroupIndexList().isEmpty() &&
          request.getFilter().getIndex().getProfileType() == Code.PROFILE_TYPE_EMPTY) {
        LOGGER.error("Empty GroupProfileFilter is set in filter of GroupProfileGetRequest");
        return Code.EMPTY_GROUP_PROFILE_FILTER;
      }
    } else if (dataCase == DataCase.FRAGMENT_SESSION_ID) {

      String fragmentSessionId = request.getFragmentSessionId();
      if (fragmentSessionId.isEmpty()) {
        LOGGER.error("Empty fragmentSessionId is set in NFProfileGetRequest");
        return Code.EMPTY_FRAGMENT_SESSION_ID;
      }
    } else {
      LOGGER.error("Empty GroupProfileGetRequest is received");
      return Code.NFMESSAGE_PROTOCOL_ERROR;
    }

    return Code.VALID;
  }

  public NFMessage createResponse(int code) {

    GroupProfileGetResponse groupProfileGetResponse = GroupProfileGetResponse.newBuilder()
        .setCode(code).build();
    GetResponse getResponse = GetResponse.newBuilder()
        .setGroupProfileGetResponse(groupProfileGetResponse).build();
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
              .transmitNumPerTime(fragmentResult, Code.GROUPPROFILE_INDICE);
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
          List<GroupProfileInfo> groupProfileInfo = getGroupProfileInfo(
              fragmentResult.getItems());

          String fragmentSessionId = fragmentResult.getFragmentSessionID();
          int totalNumber = fragmentResult.getTotalNumber();
          int transmittedNumber = fragmentResult.getTransmittedNumber();
          FragmentInfo fragmentInfo = FragmentInfo.newBuilder()
              .setFragmentSessionId(fragmentSessionId).setTotalNumber(totalNumber)
              .setTransmittedNumber(transmittedNumber).build();
          GroupProfileGetResponse groupProfileGetResponse = GroupProfileGetResponse.newBuilder()
              .setCode(searchResult.getCode()).addAllGroupProfileInfo(groupProfileInfo)
              .setFragmentInfo(fragmentInfo).build();
          GetResponse getResponse = GetResponse.newBuilder()
              .setGroupProfileGetResponse(groupProfileGetResponse).build();
          return createNFMessage(getResponse);
        }
      } else {
        List<GroupProfileInfo> groupProfileInfo = getGroupProfileInfo(searchResult.getItems());
        GroupProfileGetResponse groupProfileGetResponse = GroupProfileGetResponse.newBuilder()
            .setCode(searchResult.getCode()).addAllGroupProfileInfo(groupProfileInfo).build();
        GetResponse getResponse = GetResponse.newBuilder()
            .setGroupProfileGetResponse(groupProfileGetResponse).build();
        return createNFMessage(getResponse);
      }
    }
  }

  private List<GroupProfileInfo> getGroupProfileInfo(List<Object> items) {

    List<GroupProfileInfo> groupProfileInfoList = new ArrayList<>();
    for (Object obj : items) {
      GroupProfile item = (GroupProfile) obj;
      ByteString data = item.getData();
      String groupProfileId = item.getGroupProfileID();
      Long supiVersion = item.getSupiVersion();
      GroupProfileInfo groupProfileInfo = GroupProfileInfo.newBuilder()
          .setGroupProfileId(groupProfileId).setSupiVersion(supiVersion)
          .setGroupProfileData(data).build();
      groupProfileInfoList.add(groupProfileInfo);
    }
    return groupProfileInfoList;
  }
}