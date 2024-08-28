package ericsson.core.nrf.dbproxy.helper.groupprofile;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileDelRequestProto.*;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import static org.junit.Assert.*;

public class GroupProfileDelHelperTest {
    GroupProfileDelHelper groupProfileDelHelper;
    @Before
    public void setUp() throws Exception {
        groupProfileDelHelper = GroupProfileDelHelper.getInstance();
    }

    @After
    public void tearDown() throws Exception {
    }

    @Test
    public void validateTest() {
        {
            GroupProfileDelRequest groupProfileDelRequest = GroupProfileDelRequest.newBuilder().setGroupProfileId("").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setDelRequest(DelRequest.newBuilder().setGroupProfileDelRequest(groupProfileDelRequest).build()).build()).build();
            int result = groupProfileDelHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_GROUP_PROFILE_ID);
        }
        {
            GroupProfileDelRequest groupProfileDelRequest = GroupProfileDelRequest.newBuilder().setGroupProfileId("123456789").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setDelRequest(DelRequest.newBuilder().setGroupProfileDelRequest(groupProfileDelRequest).build()).build()).build();
            int result = groupProfileDelHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
    }

    @Test
    public void createResponseTest() {
        {
            NFMessage nfMessage = groupProfileDelHelper.createResponse(Code.SUCCESS);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getDelResponse().getGroupProfileDelResponse().getCode(), Code.SUCCESS);
        }
        {
            NFMessage nfMessage = groupProfileDelHelper.createResponse(Code.INTERNAL_ERROR);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getDelResponse().getGroupProfileDelResponse().getCode(), Code.INTERNAL_ERROR);
        }
        {
            NFMessage nfMessage = groupProfileDelHelper.createResponse(Code.DATA_NOT_EXIST);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getDelResponse().getGroupProfileDelResponse().getCode(), Code.DATA_NOT_EXIST);
        }
    }

}
