package ericsson.core.nrf.dbproxy.helper.nfprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfilePutRequestProto.*;
import org.apache.commons.lang.RandomStringUtils;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

public class NFProfilePutHelperTest {
    NFProfilePutHelper nfProfilePutHelper;
    @Before
    public void setUp() throws Exception {
        nfProfilePutHelper = NFProfilePutHelper.getInstance();
    }

    @After
    public void tearDown() throws Exception {
    }

    @Test
    public void validateTest() {
        {
            NFProfilePutRequest nfProfilePutRequest = NFProfilePutRequest.newBuilder().setNfInstanceId("").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setNfProfilePutRequest(nfProfilePutRequest).build()).build()).build();
            int result = nfProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_NF_INSTANCE_ID);
        }
        {
            NFProfilePutRequest nfProfilePutRequest = NFProfilePutRequest.newBuilder().setNfInstanceId(RandomStringUtils.random(1025)).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setNfProfilePutRequest(nfProfilePutRequest).build()).build()).build();
            int result = nfProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.NF_INSTANCE_ID_LENGTH_EXCEED_MAX);
        }
        {
            NFProfilePutRequest nfProfilePutRequest = NFProfilePutRequest.newBuilder().setNfInstanceId("test").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setNfProfilePutRequest(nfProfilePutRequest).build()).build()).build();
            int result = nfProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_NF_PROFILE);
        }
        {
            NFProfilePutRequest nfProfilePutRequest = NFProfilePutRequest.newBuilder().setNfInstanceId("test").setNfProfile("hahaha").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setPutRequest(PutRequest.newBuilder().setNfProfilePutRequest(nfProfilePutRequest).build()).build()).build();
            int result = nfProfilePutHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
    }

    @Test
    public void createResponseTest() {
        {
            NFMessage nfMessage = nfProfilePutHelper.createResponse(Code.SUCCESS);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getPutResponse().getNfProfilePutResponse().getCode(), Code.SUCCESS);
        }
        {
            NFMessage nfMessage = nfProfilePutHelper.createResponse(Code.INTERNAL_ERROR);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getPutResponse().getNfProfilePutResponse().getCode(), Code.INTERNAL_ERROR);
        }
        {
            NFMessage nfMessage = nfProfilePutHelper.createResponse(Code.DATA_NOT_EXIST);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getPutResponse().getNfProfilePutResponse().getCode(), Code.DATA_NOT_EXIST);
        }
    }
}
