package ericsson.core.nrf.dbproxy.helper.nfprofile;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileDelRequestProto.*;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;


public class NFProfileDelHelperTest {
    NFProfileDelHelper nfProfileDelHelper;
    @Before
    public void setUp() throws Exception {
        nfProfileDelHelper = NFProfileDelHelper.getInstance();
    }

    @After
    public void tearDown() throws Exception {
    }

    @Test
    public void validateTest() {
        {
            NFProfileDelRequest nfProfileDelRequest = NFProfileDelRequest.newBuilder().setNfInstanceId("").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setDelRequest(DelRequest.newBuilder().setNfProfileDelRequest(nfProfileDelRequest).build()).build()).build();
            int result = nfProfileDelHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_NF_INSTANCE_ID);
        }
        {
            NFProfileDelRequest nfProfileDelRequest = NFProfileDelRequest.newBuilder().setNfInstanceId("123456789").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setDelRequest(DelRequest.newBuilder().setNfProfileDelRequest(nfProfileDelRequest).build()).build()).build();
            int result = nfProfileDelHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
    }


    @Test
    public void createResponseTest() {
        {
            NFMessage nfMessage = nfProfileDelHelper.createResponse(Code.SUCCESS);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getDelResponse().getNfProfileDelResponse().getCode(), Code.SUCCESS);
        }
        {
            NFMessage nfMessage = nfProfileDelHelper.createResponse(Code.INTERNAL_ERROR);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getDelResponse().getNfProfileDelResponse().getCode(), Code.INTERNAL_ERROR);
        }
        {
            NFMessage nfMessage = nfProfileDelHelper.createResponse(Code.DATA_NOT_EXIST);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getDelResponse().getNfProfileDelResponse().getCode(), Code.DATA_NOT_EXIST);
        }
    }
}
