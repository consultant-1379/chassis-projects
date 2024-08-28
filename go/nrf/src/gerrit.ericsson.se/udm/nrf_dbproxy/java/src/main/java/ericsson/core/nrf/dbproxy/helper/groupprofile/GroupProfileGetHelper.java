package ericsson.core.nrf.dbproxy.helper.groupprofile;

import java.util.List;
import java.util.ArrayList;

import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import com.google.protobuf.ByteString;

import ericsson.core.nrf.dbproxy.common.*;
import ericsson.core.nrf.dbproxy.clientcache.schema.GroupProfile;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileGetRequestProto.GroupProfileGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileGetResponseProto.GroupProfileGetResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileGetResponseProto.GroupProfileInfo;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileGetRequestProto.GroupProfileGetRequest.DataCase;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;

public class GroupProfileGetHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(GroupProfileGetHelper.class);

    private static GroupProfileGetHelper instance;

    private GroupProfileGetHelper() { }

    public static synchronized GroupProfileGetHelper getInstance()
    {
        if(null == instance) {
            instance = new GroupProfileGetHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        GroupProfileGetRequest request = message.getRequest().getGetRequest().getGroupProfileGetRequest();
        DataCase data_case = request.getDataCase();
        if(data_case == DataCase.GROUP_PROFILE_ID) {

            String group_profile_id = request.getGroupProfileId();
            if(group_profile_id.isEmpty() == true) {
                logger.error("Empty group_profile_id is set in GroupProfileGetRequest");
                return Code.EMPTY_GROUP_PROFILE_ID;
            } else if(group_profile_id.length() > Code.KEY_MAX_LENGTH) {
                logger.error("group_profile_id length {} is too large, max length is {}",
                             group_profile_id.length(), Code.KEY_MAX_LENGTH);
                return Code.GROUP_PROFILE_ID_LENGTH_EXCEED_MAX;
            }
        } else if(data_case == DataCase.FILTER) {

            if(request.getFilter().getIndex().getNfTypeList().isEmpty() &&
			    request.getFilter().getIndex().getGroupIndexList().isEmpty() && 
			    request.getFilter().getIndex().getProfileType() == Code.PROFILE_TYPE_EMPTY) {
                logger.error("Empty GroupProfileFilter is set in filter of GroupProfileGetRequest");
                return Code.EMPTY_GROUP_PROFILE_FILTER;
            }
        } else if(data_case == DataCase.FRAGMENT_SESSION_ID) {

            String fragment_session_id = request.getFragmentSessionId();
            if(fragment_session_id.isEmpty()) {
                logger.error("Empty fragment_session_id is set in NFProfileGetRequest");
                return Code.EMPTY_FRAGMENT_SESSION_ID;
            }
        } else {
            logger.error("Empty GroupProfileGetRequest is received");
            return Code.NFMESSAGE_PROTOCOL_ERROR;
        }

        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {

        GroupProfileGetResponse group_profile_get_response = GroupProfileGetResponse.newBuilder().setCode(code).build();
        GetResponse get_response = GetResponse.newBuilder().setGroupProfileGetResponse(group_profile_get_response).build();
        return createNFMessage(get_response);
    }

    public NFMessage createResponse(ExecutionResult execution_result)
    {
        if(execution_result.getCode() != Code.SUCCESS) {
            return createResponse(execution_result.getCode());
        } else {
            SearchResult search_result = (SearchResult)execution_result;
            if(search_result.isFragmented()) {
                FragmentResult fragment_result = (FragmentResult)search_result;
                if(fragment_result.getFragmentSessionID().isEmpty()) {
                    int firstTransmitNum = FragmentUtil.transmitNumPerTime(fragment_result, Code.GROUPPROFILE_INDICE);
                    if(FragmentSessionManagement.getInstance().put(fragment_result, firstTransmitNum)) {
                        FragmentResult item = new FragmentResult();
                        item.addAll(fragment_result.getItems().subList(0, firstTransmitNum));
                        item.setFragmentSessionID(fragment_result.getFragmentSessionID());
                        item.setTotalNumber(fragment_result.getTotalNumber());
                        item.setTransmittedNumber(fragment_result.getTransmittedNumber());
                        return createResponse(item);
                    } else {
                        return createResponse(Code.INTERNAL_ERROR);
                    }
                } else {
                    List<GroupProfileInfo> group_profile_info = getGroupProfileInfo(fragment_result.getItems());

                    String fragment_session_id = fragment_result.getFragmentSessionID();
                    int total_number = fragment_result.getTotalNumber();
                    int transmitted_number = fragment_result.getTransmittedNumber();
                    FragmentInfo fragment_info = FragmentInfo.newBuilder().setFragmentSessionId(fragment_session_id).setTotalNumber(total_number).setTransmittedNumber(transmitted_number).build();
                    GroupProfileGetResponse group_profile_get_response = GroupProfileGetResponse.newBuilder().setCode(search_result.getCode()).addAllGroupProfileInfo(group_profile_info).setFragmentInfo(fragment_info).build();
                    GetResponse get_response = GetResponse.newBuilder().setGroupProfileGetResponse(group_profile_get_response).build();
                    return createNFMessage(get_response);
                }
            } else {
                List<GroupProfileInfo> group_profile_info = getGroupProfileInfo(search_result.getItems());
                GroupProfileGetResponse group_profile_get_response = GroupProfileGetResponse.newBuilder().setCode(search_result.getCode()).addAllGroupProfileInfo(group_profile_info).build();
                GetResponse get_response = GetResponse.newBuilder().setGroupProfileGetResponse(group_profile_get_response).build();
                return createNFMessage(get_response);
            }
        }
    }

    private List<GroupProfileInfo> getGroupProfileInfo(List<Object> items)
    {

        List<GroupProfileInfo> group_profile_info_list = new ArrayList<>();
        for(Object obj : items) {
            GroupProfile item = (GroupProfile)obj;
            ByteString data = item.getData();
            String group_profile_id = item.getGroupProfileID();
            Long supi_version = item.getSupiVersion();
            GroupProfileInfo group_profile_info = GroupProfileInfo.newBuilder().setGroupProfileId(group_profile_id).setSupiVersion(supi_version).setGroupProfileData(data).build();
            group_profile_info_list.add(group_profile_info);
        }
        return group_profile_info_list;
    }
}