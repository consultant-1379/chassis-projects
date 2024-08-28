package ericsson.core.nrf.dbproxy.helper.gpsiprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfilePutRequestProto.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileIndexProto.*;
import org.apache.commons.lang.RandomStringUtils;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

public class GpsiProfilePutHelperTest {
    GpsiProfilePutHelper gpsiProfilePutHelper;
    @Before
    public void setUp() throws Exception {
        gpsiProfilePutHelper = GpsiProfilePutHelper.getInstance();
    }

    @After
    public void tearDown() throws Exception {
    }
    @Test
    public void validateTest() {
        {
            GpsiProfilePutRequest gpsiProfilePutRequest = GpsiProfilePutRequest.newBuilder().setGpsiProfileId("").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setGpsiProfilePutRequest(gpsiProfilePutRequest).build()).build()).build();
            int result = gpsiProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_GPSI_PROFILE_ID);
        }
        {
            GpsiProfilePutRequest gpsiProfilePutRequest = GpsiProfilePutRequest.newBuilder().setGpsiProfileId(RandomStringUtils.random(1025)).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setGpsiProfilePutRequest(gpsiProfilePutRequest).build()).build()).build();
            int result = gpsiProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.GPSI_PROFILE_ID_LENGTH_EXCEED_MAX);
        }
        {
            GpsiProfilePutRequest gpsiProfilePutRequest = GpsiProfilePutRequest.newBuilder().setGpsiProfileId("test").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setGpsiProfilePutRequest(gpsiProfilePutRequest).build()).build()).build();
            int result = gpsiProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_GPSI_PROFILE_DATA);
        }
		{
            GpsiProfilePutRequest gpsiProfilePutRequest = GpsiProfilePutRequest.newBuilder().setGpsiProfileId("test").setGpsiProfileData(ByteString.copyFromUtf8("hahaha")).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setGpsiProfilePutRequest(gpsiProfilePutRequest).build()).build()).build();
            int result = gpsiProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.GPSI_PROFILE_INVALID_PROFILE_TYPE);
        }
        {
            GpsiProfileIndex gpsiProfileIndex = GpsiProfileIndex.newBuilder().setProfileType(Code.PROFILE_TYPE_GROUPID).build();
            GpsiProfilePutRequest gpsiProfilePutRequest = GpsiProfilePutRequest.newBuilder().setGpsiProfileId("test").setIndex(gpsiProfileIndex).setGpsiProfileData(ByteString.copyFromUtf8("hahaha")).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setGpsiProfilePutRequest(gpsiProfilePutRequest).build()).build()).build();
            int result = gpsiProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
    }

    @Test
    public void createResponseTest() {
        {
            NFMessage nfMessage = gpsiProfilePutHelper.createResponse(Code.SUCCESS);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getPutResponse().getGpsiProfilePutResponse().getCode(), Code.SUCCESS);
        }
        {
            NFMessage nfMessage = gpsiProfilePutHelper.createResponse(Code.INTERNAL_ERROR);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getPutResponse().getGpsiProfilePutResponse().getCode(), Code.INTERNAL_ERROR);
        }
        {
            NFMessage nfMessage = gpsiProfilePutHelper.createResponse(Code.DATA_NOT_EXIST);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getPutResponse().getGpsiProfilePutResponse().getCode(), Code.DATA_NOT_EXIST);
        }
    }
}
