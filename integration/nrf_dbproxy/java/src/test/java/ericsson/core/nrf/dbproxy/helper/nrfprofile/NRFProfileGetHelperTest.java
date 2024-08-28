package ericsson.core.nrf.dbproxy.helper.nrfprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.clientcache.schema.NRFProfile;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.FragmentResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileGetRequestProto.*;
import org.apache.commons.lang3.RandomStringUtils;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;


import static org.junit.Assert.*;

public class NRFProfileGetHelperTest {
    NRFProfileGetHelper nrfProfileGetHelper;

    @Before
    public void setUp() throws Exception {
        nrfProfileGetHelper = NRFProfileGetHelper.getInstance();
        assertNotNull(nrfProfileGetHelper);
    }

    @After
    public void tearDown() throws Exception {
    }

    @Test
    public void validateTest() {
        {
            NRFProfileGetRequest nrfProfileGetRequest = NRFProfileGetRequest.newBuilder().setNrfInstanceId("").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNrfProfileGetRequest(nrfProfileGetRequest).build()).build()).build();
            int result = nrfProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_NRF_INSTANCE_ID);
        }
        {
            NRFProfileGetRequest nrfProfileGetRequest = NRFProfileGetRequest.newBuilder().setNrfInstanceId(RandomStringUtils.random(1025)).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNrfProfileGetRequest(nrfProfileGetRequest).build()).build()).build();
            int result = nrfProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.NRF_INSTANCE_ID_LENGTH_EXCEED_MAX);
        }
        {
            NRFProfileGetRequest nrfProfileGetRequest = NRFProfileGetRequest.newBuilder().setNrfInstanceId("123456789").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNrfProfileGetRequest(nrfProfileGetRequest).build()).build()).build();
            int result = nrfProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
        {
            NRFProfileGetRequest nrfProfileGetRequest = NRFProfileGetRequest.newBuilder().setFragmentSessionId("").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNrfProfileGetRequest(nrfProfileGetRequest).build()).build()).build();
            int result = nrfProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_FRAGMENT_SESSION_ID);
        }
        {
            NRFProfileGetRequest nrfProfileGetRequest = NRFProfileGetRequest.newBuilder().setFragmentSessionId("12312").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNrfProfileGetRequest(nrfProfileGetRequest).build()).build()).build();
            int result = nrfProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
        {
            NRFProfileGetRequest nrfProfileGetRequest = NRFProfileGetRequest.newBuilder().build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNrfProfileGetRequest(nrfProfileGetRequest).build()).build()).build();
            int result = nrfProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.NFMESSAGE_PROTOCOL_ERROR);
        }
        {
            NRFProfileGetRequest nrfProfileGetRequest = NRFProfileGetRequest.newBuilder().setNrfInstanceId("123").setFragmentSessionId("123").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNrfProfileGetRequest(nrfProfileGetRequest).build()).build()).build();
            int result = nrfProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
    }

    @Test
    public void createResponseTest() {
        {
            NFMessage nfMessage = nrfProfileGetHelper.createResponse(Code.SUCCESS);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getNrfProfileGetResponse().getCode(), Code.SUCCESS);
        }
        {
            NFMessage nfMessage = nrfProfileGetHelper.createResponse(Code.INTERNAL_ERROR);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getNrfProfileGetResponse().getCode(), Code.INTERNAL_ERROR);
        }
        {
            NFMessage nfMessage = nrfProfileGetHelper.createResponse(Code.DATA_NOT_EXIST);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getNrfProfileGetResponse().getCode(), Code.DATA_NOT_EXIST);
        }
        {
            FragmentResult fragmentResult = new FragmentResult();
            for (int i = 0; i < 100; i++) {
                NRFProfile nrfProfile = new NRFProfile();
                nrfProfile.setRaw_data(ByteString.copyFromUtf8("test"));
                nrfProfile.setKey1(10000 + i);
                fragmentResult.add(nrfProfile);
            }
            NFMessage nfMessage = nrfProfileGetHelper.createResponse(fragmentResult);
            assertEquals(nfMessage.getResponse().getGetResponse().getNrfProfileGetResponse().getNrfProfileCount(), 100);
            assertEquals(nfMessage.getResponse().getGetResponse().getNrfProfileGetResponse().getCode(), Code.SUCCESS);
        }
        {
            SearchResult searchResult = new SearchResult(false);
            for (int i = 0; i < 20; i++) {
                NRFProfile nrfProfile = new NRFProfile();
                nrfProfile.setRaw_data(ByteString.copyFromUtf8("test"));
                nrfProfile.setKey1(10000 + i);
                searchResult.add(nrfProfile);
            }
            NFMessage nfMessage = nrfProfileGetHelper.createResponse(searchResult);
            assertEquals(nfMessage.getResponse().getGetResponse().getNrfProfileGetResponse().getNrfProfileCount(), 20);
            assertEquals(nfMessage.getResponse().getGetResponse().getNrfProfileGetResponse().getCode(), Code.SUCCESS);
        }
    }

}
