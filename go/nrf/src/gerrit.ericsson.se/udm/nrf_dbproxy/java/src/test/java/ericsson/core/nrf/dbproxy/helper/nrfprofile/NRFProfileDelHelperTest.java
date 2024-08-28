package ericsson.core.nrf.dbproxy.helper.nrfprofile;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileDelRequestProto.*;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import static org.junit.Assert.*;

public class NRFProfileDelHelperTest {
    NRFProfileDelHelper nrfProfileDelHelper;
    @Before
    public void setUp() throws Exception {
        nrfProfileDelHelper = NRFProfileDelHelper.getInstance();
        assertNotNull(nrfProfileDelHelper);
    }

    @After
    public void tearDown() throws Exception {
    }

    @Test
    public void validateTest() {
        {
            NRFProfileDelRequest nrfProfileDelRequest = NRFProfileDelRequest.newBuilder().setNrfInstanceId("").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setDelRequest(DelRequest.newBuilder().setNrfProfileDelRequest(nrfProfileDelRequest).build()).build()).build();
            int result = nrfProfileDelHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_NRF_INSTANCE_ID);
        }
        {
            NRFProfileDelRequest nrfProfileDelRequest = NRFProfileDelRequest.newBuilder().setNrfInstanceId("123456789").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setDelRequest(DelRequest.newBuilder().setNrfProfileDelRequest(nrfProfileDelRequest).build()).build()).build();
            int result = nrfProfileDelHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
    }

    @Test
    public void createResponseTest() {
        {
            NFMessage nfMessage = nrfProfileDelHelper.createResponse(Code.SUCCESS);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getDelResponse().getNrfProfileDelResponse().getCode(), Code.SUCCESS);
        }
        {
            NFMessage nfMessage = nrfProfileDelHelper.createResponse(Code.INTERNAL_ERROR);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getDelResponse().getNrfProfileDelResponse().getCode(), Code.INTERNAL_ERROR);
        }
        {
            NFMessage nfMessage = nrfProfileDelHelper.createResponse(Code.DATA_NOT_EXIST);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getDelResponse().getNrfProfileDelResponse().getCode(), Code.DATA_NOT_EXIST);
        }
    }
}
