package ericsson.core.nrf.dbproxy.helper.nrfprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.clientcache.schema.NRFProfile;
import ericsson.core.nrf.dbproxy.common.*;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileGetRequestProto.NRFProfileGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileGetRequestProto.NRFProfileGetRequest.DataCase;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileGetResponseProto;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileGetResponseProto.FragmentNRFProfileInfo;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileGetResponseProto.NRFProfileGetResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileGetResponseProto.NRFProfileInfo;
import ericsson.core.nrf.dbproxy.helper.Helper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.util.ArrayList;
import java.util.List;

public class NRFProfileGetHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(NRFProfileGetHelper.class);

    private static NRFProfileGetHelper instance;

    private NRFProfileGetHelper()
    {
    }

    public static synchronized NRFProfileGetHelper getInstance()
    {
        if (null == instance) {
            instance = new NRFProfileGetHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {

        NRFProfileGetRequest request = message.getRequest().getGetRequest().getNrfProfileGetRequest();
        DataCase data_case = request.getDataCase();
        if (data_case == DataCase.NRF_INSTANCE_ID) {
            String nrf_instance_id = request.getNrfInstanceId();
            if (nrf_instance_id.isEmpty() == true) {
                logger.error("Empty nrf_instance_id is set in NRFProfileGetRequest");
                return Code.EMPTY_NRF_INSTANCE_ID;
            } else if (nrf_instance_id.length() > Code.KEY_MAX_LENGTH) {
                logger.error("nrf_instance_id length {} is too large, max length is {}",
                             nrf_instance_id.length(), Code.KEY_MAX_LENGTH);
                return Code.NRF_INSTANCE_ID_LENGTH_EXCEED_MAX;
            }
        } else if (data_case == DataCase.FILTER) {
            logger.trace("NRFProfile filter");
        } else if (data_case == DataCase.FRAGMENT_SESSION_ID) {

            String fragment_session_id = request.getFragmentSessionId();
            if (fragment_session_id.isEmpty()) {
                logger.error("Empty fragment_session_id is set in NRFProfileGetRequest");
                return Code.EMPTY_FRAGMENT_SESSION_ID;
            }
        } else {
            logger.error("Empty NRFProfileGetRequest is received");
            return Code.NFMESSAGE_PROTOCOL_ERROR;
        }

        return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {

        NRFProfileGetResponse nrf_profile_get_response = NRFProfileGetResponse.newBuilder().setCode(code).build();
        GetResponse get_response = GetResponse.newBuilder().setNrfProfileGetResponse(nrf_profile_get_response).build();
        return createNFMessage(get_response);
    }

    public NFMessage createResponse(ExecutionResult execution_result)
    {
        if (execution_result.getCode() != Code.SUCCESS) {
            return createResponse(execution_result.getCode());
        } else {
            SearchResult search_result = (SearchResult) execution_result;
            if (search_result.isFragmented()) {
                FragmentResult fragment_result = (FragmentResult) search_result;
                if (fragment_result.getFragmentSessionID().isEmpty()) {
                    int firstTransmitNum = FragmentUtil.transmitNumPerTime(fragment_result, Code.NRFPROFILE_INDICE);
                    if (FragmentSessionManagement.getInstance().put(fragment_result, firstTransmitNum)) {
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
                    List<NRFProfileInfo> nrf_profiles = getNRFProfileInfo(fragment_result.getItems());

                    String fragment_session_id = fragment_result.getFragmentSessionID();
                    int total_number = fragment_result.getTotalNumber();
                    int transmitted_number = fragment_result.getTransmittedNumber();
                    FragmentNRFProfileInfo fragment_info = FragmentNRFProfileInfo.newBuilder().setFragmentSessionId(fragment_session_id).setTotalNumber(total_number).setTransmittedNumber(transmitted_number).build();

                    NRFProfileGetResponseProto.NRFProfileGetResponse nrf_profile_get_response = NRFProfileGetResponse.newBuilder().setCode(fragment_result.getCode()).addAllNrfProfile(nrf_profiles).setFragmentNrfprofileInfo(fragment_info).build();
                    GetResponse get_response = GetResponse.newBuilder().setNrfProfileGetResponse(nrf_profile_get_response).build();
                    return createNFMessage(get_response);
                }
            } else {
                int provFlag = 0;
                if (search_result.getItems().size() == 1) {
                    provFlag = (int)getProvFlag(search_result.getItems().get(0));
                }

                List<NRFProfileInfo> nrf_profiles = getNRFProfileInfo(search_result.getItems());
                NRFProfileGetResponse nrf_profile_get_response = NRFProfileGetResponse.newBuilder().setCode(search_result.getCode()).addAllNrfProfile(nrf_profiles).setProvFlag(provFlag).build();
                GetResponse get_response = GetResponse.newBuilder().setNrfProfileGetResponse(nrf_profile_get_response).build();
                return createNFMessage(get_response);
            }
        }

    }

    private long getProvFlag(Object obj)
    {
        NRFProfile nrf_profile = (NRFProfile)obj;
        long prov_flag = nrf_profile.getProvFlag();

        return prov_flag;
    }

    private List<NRFProfileInfo> getNRFProfileInfo(List<Object> items)
    {

        List<NRFProfileInfo> nrf_profiles = new ArrayList<>();
        for (Object obj : items) {
            NRFProfile nrf_profile = (NRFProfile) obj;
            ByteString data = nrf_profile.getRaw_data();
            long expired_time = nrf_profile.getExpireTime();

            NRFProfileInfo item = NRFProfileInfo.newBuilder().setRawNrfProfile(data).setExpiredTime(expired_time).build();
            nrf_profiles.add(item);
        }

        return nrf_profiles;
    }


}
