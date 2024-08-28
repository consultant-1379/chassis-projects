package ericsson.core.nrf.dbproxy.helper.nrfprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfilePutRequestProto.*;
import org.apache.commons.lang3.RandomStringUtils;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import static org.junit.Assert.*;

public class NRFProfilePutHelperTest {
    NRFProfilePutHelper nrfProfilePutHelper;
    @Before
    public void setUp() throws Exception {
        nrfProfilePutHelper = NRFProfilePutHelper.getInstance();
    }

    @After
    public void tearDown() throws Exception {
    }

    @Test
    public void validateTest() {
        {
            NRFProfilePutRequest nrfProfilePutRequest = NRFProfilePutRequest.newBuilder().setNrfInstanceId("").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setNrfProfilePutRequest(nrfProfilePutRequest).build()).build()).build();
            int result = nrfProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_NRF_INSTANCE_ID);
        }
        {
            NRFProfilePutRequest nrfProfilePutRequest = NRFProfilePutRequest.newBuilder().setNrfInstanceId(RandomStringUtils.random(1025)).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setNrfProfilePutRequest(nrfProfilePutRequest).build()).build()).build();
            int result = nrfProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.NRF_INSTANCE_ID_LENGTH_EXCEED_MAX);
        }
        {
            NRFProfilePutRequest nrfProfilePutRequest = NRFProfilePutRequest.newBuilder().setNrfInstanceId("test").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setNrfProfilePutRequest(nrfProfilePutRequest).build()).build()).build();
            int result = nrfProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_RAW_NRF_PROFILE);
        }
        {
            NRFProfilePutRequest nrfProfilePutRequest = NRFProfilePutRequest.newBuilder().setNrfInstanceId("test").setRawNrfProfile(ByteString.copyFromUtf8("hahaha")).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setNrfProfilePutRequest(nrfProfilePutRequest).build()).build()).build();
            int result = nrfProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
    }

    @Test
    public void createResponseTest() {
        {
            NFMessage nfMessage = nrfProfilePutHelper.createResponse(Code.SUCCESS);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getPutResponse().getNrfProfilePutResponse().getCode(), Code.SUCCESS);
        }
        {
            NFMessage nfMessage = nrfProfilePutHelper.createResponse(Code.INTERNAL_ERROR);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getPutResponse().getNrfProfilePutResponse().getCode(), Code.INTERNAL_ERROR);
        }
        {
            NFMessage nfMessage = nrfProfilePutHelper.createResponse(Code.DATA_NOT_EXIST);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getPutResponse().getNrfProfilePutResponse().getCode(), Code.DATA_NOT_EXIST);
        }
    }
}
