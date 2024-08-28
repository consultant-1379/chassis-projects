package ericsson.core.nrf.dbproxy.helper.groupprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.clientcache.schema.GroupProfile;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.FragmentResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileFilterProto;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileGetRequestProto.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileIndexProto;
import org.apache.commons.lang3.RandomStringUtils;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

public class GroupProfileGetHelperTest {
    GroupProfileGetHelper groupProfileGetHelper;
    @Before
    public void setUp() throws Exception {
        groupProfileGetHelper = GroupProfileGetHelper.getInstance();
    }

    @After
    public void tearDown() throws Exception {
    }
    
    @Test
    public void validateTest() {
        {
            GroupProfileGetRequest groupProfileGetRequest = GroupProfileGetRequest.newBuilder().setGroupProfileId("").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGroupProfileGetRequest(groupProfileGetRequest).build()).build()).build();
            int result = groupProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_GROUP_PROFILE_ID);
        }
        {
            GroupProfileGetRequest groupProfileGetRequest = GroupProfileGetRequest.newBuilder().setGroupProfileId(RandomStringUtils.random(1025)).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGroupProfileGetRequest(groupProfileGetRequest).build()).build()).build();
            int result = groupProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.GROUP_PROFILE_ID_LENGTH_EXCEED_MAX);
        }
        {
            GroupProfileGetRequest groupProfileGetRequest = GroupProfileGetRequest.newBuilder().setGroupProfileId("123456789").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGroupProfileGetRequest(groupProfileGetRequest).build()).build()).build();
            int result = groupProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
        {
            GroupProfileGetRequest groupProfileGetRequest = GroupProfileGetRequest.newBuilder().setFilter(GroupProfileFilterProto.GroupProfileFilter.newBuilder().setIndex(GroupProfileIndexProto.GroupProfileIndex.newBuilder().build()).build()).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGroupProfileGetRequest(groupProfileGetRequest).build()).build()).build();
            int result = groupProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_GROUP_PROFILE_FILTER);
        }
        {
            GroupProfileGetRequest groupProfileGetRequest = GroupProfileGetRequest.newBuilder().setFilter(GroupProfileFilterProto.GroupProfileFilter.newBuilder().setIndex(GroupProfileIndexProto.GroupProfileIndex.newBuilder().addNfType("AMF").build()).build()).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGroupProfileGetRequest(groupProfileGetRequest).build()).build()).build();
            int result = groupProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
        {
            GroupProfileGetRequest groupProfileGetRequest = GroupProfileGetRequest.newBuilder().setFragmentSessionId("").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGroupProfileGetRequest(groupProfileGetRequest).build()).build()).build();
            int result = groupProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_FRAGMENT_SESSION_ID);
        }
        {
            GroupProfileGetRequest groupProfileGetRequest = GroupProfileGetRequest.newBuilder().setFragmentSessionId("12312").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGroupProfileGetRequest(groupProfileGetRequest).build()).build()).build();
            int result = groupProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
        {
            GroupProfileGetRequest groupProfileGetRequest = GroupProfileGetRequest.newBuilder().build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGroupProfileGetRequest(groupProfileGetRequest).build()).build()).build();
            int result = groupProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.NFMESSAGE_PROTOCOL_ERROR);
        }
        {
            GroupProfileGetRequest groupProfileGetRequest = GroupProfileGetRequest.newBuilder().setGroupProfileId("123").setFragmentSessionId("123").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGroupProfileGetRequest(groupProfileGetRequest).build()).build()).build();
            int result = groupProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
    }
    @Test
    public void createResponseTest() {
        {
            NFMessage nfMessage = groupProfileGetHelper.createResponse(Code.SUCCESS);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getGroupProfileGetResponse().getCode(), Code.SUCCESS);
        }
        {
            NFMessage nfMessage = groupProfileGetHelper.createResponse(Code.INTERNAL_ERROR);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getGroupProfileGetResponse().getCode(), Code.INTERNAL_ERROR);
        }
        {
            NFMessage nfMessage = groupProfileGetHelper.createResponse(Code.DATA_NOT_EXIST);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getGroupProfileGetResponse().getCode(), Code.DATA_NOT_EXIST);
        }
        {
            FragmentResult fragmentResult = new FragmentResult();
            for (int i = 0; i < 100; i++) {
                GroupProfile groupProfile = new GroupProfile();
                groupProfile.setData(ByteString.copyFromUtf8("test"));
                groupProfile.setGroupProfileID(10000 + i + "");
                fragmentResult.add(groupProfile);
            }
            NFMessage nfMessage = groupProfileGetHelper.createResponse(fragmentResult);
//            assertEquals(nfMessage.getResponse().getGetResponse().getGroupProfileGetResponse().getGroupProfileInfoCount(), Code.FRAGMENT_BLOCK_GROUP_PROFILE);
            assertEquals(nfMessage.getResponse().getGetResponse().getGroupProfileGetResponse().getCode(), Code.SUCCESS);
        }
        {
            SearchResult searchResult = new SearchResult(false);
            for (int i = 0; i < 20; i++) {
                GroupProfile groupProfile = new GroupProfile();
                groupProfile.setData(ByteString.copyFromUtf8("test"));
                groupProfile.setGroupProfileID(10000 + i + "");
                searchResult.add(groupProfile);
            }
            NFMessage nfMessage = groupProfileGetHelper.createResponse(searchResult);
//            assertEquals(nfMessage.getResponse().getGetResponse().getGroupProfileGetResponse().getGroupProfileInfoCount(), 20);
            assertEquals(nfMessage.getResponse().getGetResponse().getGroupProfileGetResponse().getCode(), Code.SUCCESS);
        }
    }
}
