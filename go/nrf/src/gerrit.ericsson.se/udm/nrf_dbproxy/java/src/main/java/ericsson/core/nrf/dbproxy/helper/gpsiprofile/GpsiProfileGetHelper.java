package ericsson.core.nrf.dbproxy.helper.gpsiprofile;

import java.util.List;
import java.util.ArrayList;

import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import com.google.protobuf.ByteString;

import ericsson.core.nrf.dbproxy.common.*;
import ericsson.core.nrf.dbproxy.clientcache.schema.GpsiProfile;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileGetRequestProto.GpsiProfileGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileGetResponseProto.GpsiProfileGetResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileGetResponseProto.GpsiProfileInfo;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileGetRequestProto.GpsiProfileGetRequest.DataCase;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;

public class GpsiProfileGetHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(GpsiProfileGetHelper.class);

    private static GpsiProfileGetHelper instance;

    private GpsiProfileGetHelper() { }

    public static synchronized GpsiProfileGetHelper getInstance()
    {
        if(null == instance) {
            instance = new GpsiProfileGetHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        GpsiProfileGetRequest request = message.getRequest().getGetRequest().getGpsiProfileGetRequest();
        DataCase data_case = request.getDataCase();
        if(data_case == DataCase.GPSI_PROFILE_ID) {

            String gpsi_profile_id = request.getGpsiProfileId();
            if(gpsi_profile_id.isEmpty() == true) {
                logger.error("Empty gpsi_profile_id is set in GpsiProfileGetRequest");
                return Code.EMPTY_GPSI_PROFILE_ID;
            } else if(gpsi_profile_id.length() > Code.KEY_MAX_LENGTH) {
                logger.error("gpsi_profile_id length {} is too large, max length is {}",
                             gpsi_profile_id.length(), Code.KEY_MAX_LENGTH);
                return Code.GPSI_PROFILE_ID_LENGTH_EXCEED_MAX;
            }
        } else if(data_case == DataCase.FILTER) {

            if(request.getFilter().getIndex().getNfTypeList().isEmpty() && 
			   request.getFilter().getIndex().getGroupIndexList().isEmpty() && 
			   request.getFilter().getIndex().getProfileType() == Code.PROFILE_TYPE_EMPTY) {
                logger.error("Empty GpsiProfileFilter is set in filter of GpsiProfileGetRequest");
                return Code.EMPTY_GPSI_PROFILE_FILTER;
            }
        } else if(data_case == DataCase.FRAGMENT_SESSION_ID) {

            String fragment_session_id = request.getFragmentSessionId();
            if(fragment_session_id.isEmpty()) {
                logger.error("Empty fragment_session_id is set in NFProfileGetRequest");
                return Code.EMPTY_FRAGMENT_SESSION_ID;
            }
        } else {
            logger.error("Empty GpsiProfileGetRequest is received");
            return Code.NFMESSAGE_PROTOCOL_ERROR;
        }

        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {

        GpsiProfileGetResponse gpsi_profile_get_response = GpsiProfileGetResponse.newBuilder().setCode(code).build();
        GetResponse get_response = GetResponse.newBuilder().setGpsiProfileGetResponse(gpsi_profile_get_response).build();
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
                    int firstTransmitNum = FragmentUtil.transmitNumPerTime(fragment_result, Code.GPSIPROFILE_INDICE);
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
                    List<GpsiProfileInfo> gpsi_profile_info = getGpsiProfileInfo(fragment_result.getItems());

                    String fragment_session_id = fragment_result.getFragmentSessionID();
                    int total_number = fragment_result.getTotalNumber();
                    int transmitted_number = fragment_result.getTransmittedNumber();
                    FragmentInfo fragment_info = FragmentInfo.newBuilder().setFragmentSessionId(fragment_session_id).setTotalNumber(total_number).setTransmittedNumber(transmitted_number).build();
                    GpsiProfileGetResponse gpsi_profile_get_response = GpsiProfileGetResponse.newBuilder().setCode(search_result.getCode()).addAllGpsiProfileInfo(gpsi_profile_info).setFragmentInfo(fragment_info).build();
                    GetResponse get_response = GetResponse.newBuilder().setGpsiProfileGetResponse(gpsi_profile_get_response).build();
                    return createNFMessage(get_response);
                }
            } else {
                List<GpsiProfileInfo> gpsi_profile_info = getGpsiProfileInfo(search_result.getItems());
                GpsiProfileGetResponse gpsi_profile_get_response = GpsiProfileGetResponse.newBuilder().setCode(search_result.getCode()).addAllGpsiProfileInfo(gpsi_profile_info).build();
                GetResponse get_response = GetResponse.newBuilder().setGpsiProfileGetResponse(gpsi_profile_get_response).build();
                return createNFMessage(get_response);
            }
        }
    }

    private List<GpsiProfileInfo> getGpsiProfileInfo(List<Object> items)
    {

        List<GpsiProfileInfo> gpsi_profile_info_list = new ArrayList<>();
        for(Object obj : items) {
            GpsiProfile item = (GpsiProfile)obj;
            ByteString data = item.getData();
            String gpsi_profile_id = item.getGpsiProfileID();
            Long gpsi_version = item.getGpsiVersion();
            GpsiProfileInfo gpsi_profile_info = GpsiProfileInfo.newBuilder().setGpsiProfileId(gpsi_profile_id).setGpsiVersion(gpsi_version).setGpsiProfileData(data).build();
            gpsi_profile_info_list.add(gpsi_profile_info);
        }
        return gpsi_profile_info_list;
    }
}