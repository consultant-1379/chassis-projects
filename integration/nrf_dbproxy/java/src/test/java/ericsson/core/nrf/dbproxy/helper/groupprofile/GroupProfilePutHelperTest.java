package ericsson.core.nrf.dbproxy.helper.groupprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfilePutRequestProto.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileIndexProto.*;
import org.apache.commons.lang3.RandomStringUtils;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

public class GroupProfilePutHelperTest {
    GroupProfilePutHelper groupProfilePutHelper;
    @Before
    public void setUp() throws Exception {
        groupProfilePutHelper = GroupProfilePutHelper.getInstance();
    }

    @After
    public void tearDown() throws Exception {
    }
    @Test
    public void validateTest() {
        {
            GroupProfilePutRequest groupProfilePutRequest = GroupProfilePutRequest.newBuilder().setGroupProfileId("").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setGroupProfilePutRequest(groupProfilePutRequest).build()).build()).build();
            int result = groupProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_GROUP_PROFILE_ID);
        }
        {
            GroupProfilePutRequest groupProfilePutRequest = GroupProfilePutRequest.newBuilder().setGroupProfileId(RandomStringUtils.random(1025)).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setGroupProfilePutRequest(groupProfilePutRequest).build()).build()).build();
            int result = groupProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.GROUP_PROFILE_ID_LENGTH_EXCEED_MAX);
        }
        {
            GroupProfilePutRequest groupProfilePutRequest = GroupProfilePutRequest.newBuilder().setGroupProfileId("test").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setGroupProfilePutRequest(groupProfilePutRequest).build()).build()).build();
            int result = groupProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_GROUP_PROFILE_DATA);
        }
        {
            GroupProfilePutRequest groupProfilePutRequest = GroupProfilePutRequest.newBuilder().setGroupProfileId("test").setGroupProfileData(ByteString.copyFromUtf8("hahaha")).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setGroupProfilePutRequest(groupProfilePutRequest).build()).build()).build();
            int result = groupProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.GROUP_PROFILE_INVALID_PROFILE_TYPE);
        }
        {
			GroupProfileIndex groupProfileIndex = GroupProfileIndex.newBuilder().setProfileType(Code.PROFILE_TYPE_GROUPID).build();
            GroupProfilePutRequest groupProfilePutRequest = GroupProfilePutRequest.newBuilder().setGroupProfileId("test").setIndex(groupProfileIndex).setGroupProfileData(ByteString.copyFromUtf8("hahaha")).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setGroupProfilePutRequest(groupProfilePutRequest).build()).build()).build();
            int result = groupProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
    }

    @Test
    public void createResponseTest() {
        {
            NFMessage nfMessage = groupProfilePutHelper.createResponse(Code.SUCCESS);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getPutResponse().getGroupProfilePutResponse().getCode(), Code.SUCCESS);
        }
        {
            NFMessage nfMessage = groupProfilePutHelper.createResponse(Code.INTERNAL_ERROR);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getPutResponse().getGroupProfilePutResponse().getCode(), Code.INTERNAL_ERROR);
        }
        {
            NFMessage nfMessage = groupProfilePutHelper.createResponse(Code.DATA_NOT_EXIST);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getPutResponse().getGroupProfilePutResponse().getCode(), Code.DATA_NOT_EXIST);
        }
    }
}
