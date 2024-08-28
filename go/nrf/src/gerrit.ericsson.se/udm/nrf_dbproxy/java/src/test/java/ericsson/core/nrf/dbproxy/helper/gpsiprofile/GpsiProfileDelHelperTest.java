package ericsson.core.nrf.dbproxy.helper.gpsiprofile;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileDelRequestProto.*;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import static org.junit.Assert.*;

public class GpsiProfileDelHelperTest {
    GpsiProfileDelHelper gpsiProfileDelHelper;
    @Before
    public void setUp() throws Exception {
        gpsiProfileDelHelper = GpsiProfileDelHelper.getInstance();
    }

    @After
    public void tearDown() throws Exception {
    }

    @Test
    public void validateTest() {
        {
            GpsiProfileDelRequest gpsiProfileDelRequest = GpsiProfileDelRequest.newBuilder().setGpsiProfileId("").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setDelRequest(DelRequest.newBuilder().setGpsiProfileDelRequest(gpsiProfileDelRequest).build()).build()).build();
            int result = gpsiProfileDelHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_GPSI_PROFILE_ID);
        }
        {
            GpsiProfileDelRequest gpsiProfileDelRequest = GpsiProfileDelRequest.newBuilder().setGpsiProfileId("123456789").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setDelRequest(DelRequest.newBuilder().setGpsiProfileDelRequest(gpsiProfileDelRequest).build()).build()).build();
            int result = gpsiProfileDelHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
    }

    @Test
    public void createResponseTest() {
        {
            NFMessage nfMessage = gpsiProfileDelHelper.createResponse(Code.SUCCESS);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getDelResponse().getGpsiProfileDelResponse().getCode(), Code.SUCCESS);
        }
        {
            NFMessage nfMessage = gpsiProfileDelHelper.createResponse(Code.INTERNAL_ERROR);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getDelResponse().getGpsiProfileDelResponse().getCode(), Code.INTERNAL_ERROR);
        }
        {
            NFMessage nfMessage = gpsiProfileDelHelper.createResponse(Code.DATA_NOT_EXIST);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getDelResponse().getGpsiProfileDelResponse().getCode(), Code.DATA_NOT_EXIST);
        }
    }

}
