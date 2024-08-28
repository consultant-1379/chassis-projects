package ericsson.core.nrf.dbproxy.helper.gpsiprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.clientcache.schema.GpsiProfile;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.FragmentResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileFilterProto;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileGetRequestProto.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileIndexProto;
import org.apache.commons.lang.RandomStringUtils;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

public class GpsiProfileGetHelperTest {
    GpsiProfileGetHelper gpsiProfileGetHelper;
    @Before
    public void setUp() throws Exception {
        gpsiProfileGetHelper = GpsiProfileGetHelper.getInstance();
    }

    @After
    public void tearDown() throws Exception {
    }
    
    @Test
    public void validateTest() {
        {
            GpsiProfileGetRequest gpsiProfileGetRequest = GpsiProfileGetRequest.newBuilder().setGpsiProfileId("").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGpsiProfileGetRequest(gpsiProfileGetRequest).build()).build()).build();
            int result = gpsiProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_GPSI_PROFILE_ID);
        }
        {
            GpsiProfileGetRequest gpsiProfileGetRequest = GpsiProfileGetRequest.newBuilder().setGpsiProfileId(RandomStringUtils.random(1025)).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGpsiProfileGetRequest(gpsiProfileGetRequest).build()).build()).build();
            int result = gpsiProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.GPSI_PROFILE_ID_LENGTH_EXCEED_MAX);
        }
        {
            GpsiProfileGetRequest gpsiProfileGetRequest = GpsiProfileGetRequest.newBuilder().setGpsiProfileId("123456789").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGpsiProfileGetRequest(gpsiProfileGetRequest).build()).build()).build();
            int result = gpsiProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
        {
            GpsiProfileGetRequest gpsiProfileGetRequest = GpsiProfileGetRequest.newBuilder().setFilter(GpsiProfileFilterProto.GpsiProfileFilter.newBuilder().setIndex(GpsiProfileIndexProto.GpsiProfileIndex.newBuilder().build()).build()).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGpsiProfileGetRequest(gpsiProfileGetRequest).build()).build()).build();
            int result = gpsiProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_GPSI_PROFILE_FILTER);
        }
        {
            GpsiProfileGetRequest gpsiProfileGetRequest = GpsiProfileGetRequest.newBuilder().setFilter(GpsiProfileFilterProto.GpsiProfileFilter.newBuilder().setIndex(GpsiProfileIndexProto.GpsiProfileIndex.newBuilder().addNfType("AMF").build()).build()).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGpsiProfileGetRequest(gpsiProfileGetRequest).build()).build()).build();
            int result = gpsiProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
        {
            GpsiProfileGetRequest gpsiProfileGetRequest = GpsiProfileGetRequest.newBuilder().setFragmentSessionId("").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGpsiProfileGetRequest(gpsiProfileGetRequest).build()).build()).build();
            int result = gpsiProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_FRAGMENT_SESSION_ID);
        }
        {
            GpsiProfileGetRequest gpsiProfileGetRequest = GpsiProfileGetRequest.newBuilder().setFragmentSessionId("12312").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGpsiProfileGetRequest(gpsiProfileGetRequest).build()).build()).build();
            int result = gpsiProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
        {
            GpsiProfileGetRequest gpsiProfileGetRequest = GpsiProfileGetRequest.newBuilder().build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGpsiProfileGetRequest(gpsiProfileGetRequest).build()).build()).build();
            int result = gpsiProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.NFMESSAGE_PROTOCOL_ERROR);
        }
        {
            GpsiProfileGetRequest gpsiProfileGetRequest = GpsiProfileGetRequest.newBuilder().setGpsiProfileId("123").setFragmentSessionId("123").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGpsiProfileGetRequest(gpsiProfileGetRequest).build()).build()).build();
            int result = gpsiProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
    }
    @Test
    public void createResponseTest() {
        {
            NFMessage nfMessage = gpsiProfileGetHelper.createResponse(Code.SUCCESS);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getGpsiProfileGetResponse().getCode(), Code.SUCCESS);
        }
        {
            NFMessage nfMessage = gpsiProfileGetHelper.createResponse(Code.INTERNAL_ERROR);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getGpsiProfileGetResponse().getCode(), Code.INTERNAL_ERROR);
        }
        {
            NFMessage nfMessage = gpsiProfileGetHelper.createResponse(Code.DATA_NOT_EXIST);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getGpsiProfileGetResponse().getCode(), Code.DATA_NOT_EXIST);
        }
        {
            FragmentResult fragmentResult = new FragmentResult();
            for (int i = 0; i < 100; i++) {
                GpsiProfile gpsiProfile = new GpsiProfile();
                gpsiProfile.setData(ByteString.copyFromUtf8("test"));
                gpsiProfile.setGpsiProfileID(10000 + i + "");
                fragmentResult.add(gpsiProfile);
            }
            NFMessage nfMessage = gpsiProfileGetHelper.createResponse(fragmentResult);
//            assertEquals(nfMessage.getResponse().getGetResponse().getGpsiProfileGetResponse().getGpsiProfileInfoCount(), Code.FRAGMENT_BLOCK_GPSI_PROFILE);
            assertEquals(nfMessage.getResponse().getGetResponse().getGpsiProfileGetResponse().getCode(), Code.SUCCESS);
        }
        {
            SearchResult searchResult = new SearchResult(false);
            for (int i = 0; i < 20; i++) {
                GpsiProfile gpsiProfile = new GpsiProfile();
                gpsiProfile.setData(ByteString.copyFromUtf8("test"));
                gpsiProfile.setGpsiProfileID(10000 + i + "");
                searchResult.add(gpsiProfile);
            }
            NFMessage nfMessage = gpsiProfileGetHelper.createResponse(searchResult);
//            assertEquals(nfMessage.getResponse().getGetResponse().getGpsiProfileGetResponse().getGpsiProfileInfoCount(), 20);
            assertEquals(nfMessage.getResponse().getGetResponse().getGpsiProfileGetResponse().getCode(), Code.SUCCESS);
        }
    }
}
