package ericsson.core.nrf.dbproxy.helper.nfprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.clientcache.schema.NRFProfile;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.FragmentResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileGetRequestProto.*;
import org.apache.commons.lang.RandomStringUtils;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

public class NFProfileGetHelperTest {
    NFProfileGetHelper nfProfileGetHelper;

    @Before
    public void setUp() throws Exception {
        nfProfileGetHelper = NFProfileGetHelper.getInstance();
        assertNotNull(nfProfileGetHelper);
    }

    @After
    public void tearDown() throws Exception {
    }

    @Test
    public void validateTest() {
        {
            NFProfileGetRequest nfProfileGetRequest = NFProfileGetRequest.newBuilder().setFragmentSessionId("").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileGetRequest(nfProfileGetRequest).build()).build()).build();
            int result = nfProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.NFMESSAGE_PROTOCOL_ERROR);
        }
        {
            NFProfileGetRequest nfProfileGetRequest = NFProfileGetRequest.newBuilder().setFragmentSessionId("12312").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileGetRequest(nfProfileGetRequest).build()).build()).build();
            int result = nfProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
        {
            NFProfileGetRequest nfProfileGetRequest = NFProfileGetRequest.newBuilder().build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileGetRequest(nfProfileGetRequest).build()).build()).build();
            int result = nfProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.NFMESSAGE_PROTOCOL_ERROR);
        }
        {
            NFProfileGetRequest nfProfileGetRequest = NFProfileGetRequest.newBuilder().setFragmentSessionId("123").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileGetRequest(nfProfileGetRequest).build()).build()).build();
            int result = nfProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
    }

    @Test
    public void createResponseTest() {
        {
            NFMessage nfMessage = nfProfileGetHelper.createResponse(Code.SUCCESS);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getNfProfileGetResponse().getCode(), Code.SUCCESS);
        }
        {
            NFMessage nfMessage = nfProfileGetHelper.createResponse(Code.INTERNAL_ERROR);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getNfProfileGetResponse().getCode(), Code.INTERNAL_ERROR);
        }
        {
            NFMessage nfMessage = nfProfileGetHelper.createResponse(Code.DATA_NOT_EXIST);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getNfProfileGetResponse().getCode(), Code.DATA_NOT_EXIST);
        }
    }
}
