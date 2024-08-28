package ericsson.core.nrf.dbproxy.helper.nfprofile;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.FragmentResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetRequest;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileFilterProto;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileFilterProto.NFProfileFilter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileFilterProto.Range;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileGetRequestProto.NFProfileGetRequest;
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
            NFProfileGetRequest nfProfileGetRequest = NFProfileGetRequest.newBuilder().setTargetNfInstanceId("").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileGetRequest(nfProfileGetRequest).build()).build()).build();
            int result = nfProfileGetHelper.validate(nfMessage);
            assertEquals(Code.NFMESSAGE_PROTOCOL_ERROR, result);
        }
        {
            NFProfileGetRequest nfProfileGetRequest = NFProfileGetRequest.newBuilder().setTargetNfInstanceId("4947a69a-f61b-4bc1-b9da-47c9c5d14b64").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileGetRequest(nfProfileGetRequest).build()).build()).build();
            int result = nfProfileGetHelper.validate(nfMessage);
            assertEquals(Code.VALID, result);
        }
        {
            NFProfileFilter filter = NFProfileFilter.newBuilder().setExpiredTimeRange(Range.newBuilder().setStart(123).setEnd(1234).build()).build();
            NFProfileGetRequest nfProfileGetRequest = NFProfileGetRequest.newBuilder().setFilter(filter).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileGetRequest(nfProfileGetRequest).build()).build()).build();
            int result = nfProfileGetHelper.validate(nfMessage);
            assertEquals(Code.VALID, result);
        }
        {
            NFProfileFilter filter = NFProfileFilter.newBuilder().setLastUpdateTimeRange(Range.newBuilder().setStart(123).setEnd(1234).build()).build();
            NFProfileGetRequest nfProfileGetRequest = NFProfileGetRequest.newBuilder().setFilter(filter).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileGetRequest(nfProfileGetRequest).build()).build()).build();
            int result = nfProfileGetHelper.validate(nfMessage);
            assertEquals(Code.VALID, result);
        }
        {
            NFProfileFilter filter = NFProfileFilter.newBuilder().setLastUpdateTimeRange(Range.newBuilder().setStart(12322).setEnd(1234).build()).build();
            NFProfileGetRequest nfProfileGetRequest = NFProfileGetRequest.newBuilder().setFilter(filter).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileGetRequest(nfProfileGetRequest).build()).build()).build();
            int result = nfProfileGetHelper.validate(nfMessage);
            assertEquals(Code.INVALID_RANGE, result);
        }
        {
            NFProfileFilter filter = NFProfileFilter.newBuilder().setProvisioned(1).setProvVersion(NFProfileFilterProto.ProvVersion.newBuilder().setSupiVersion(123).build()).build();
            NFProfileGetRequest nfProfileGetRequest = NFProfileGetRequest.newBuilder().setFilter(filter).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileGetRequest(nfProfileGetRequest).build()).build()).build();
            int result = nfProfileGetHelper.validate(nfMessage);
            assertEquals(Code.VALID, result);
        }
        {
            NFProfileFilter filter = NFProfileFilter.newBuilder().setProvisioned(3).build();
            NFProfileGetRequest nfProfileGetRequest = NFProfileGetRequest.newBuilder().setFilter(filter).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileGetRequest(nfProfileGetRequest).build()).build()).build();
            int result = nfProfileGetHelper.validate(nfMessage);
            assertEquals(Code.INVALID_PROVISIONED, result);

        }
        {
            NFProfileGetRequest nfProfileGetRequest = NFProfileGetRequest.newBuilder().setFragmentSessionId("").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileGetRequest(nfProfileGetRequest).build()).build()).build();
            int result = nfProfileGetHelper.validate(nfMessage);
            assertEquals(Code.NFMESSAGE_PROTOCOL_ERROR, result);
        }
        {
            NFProfileGetRequest nfProfileGetRequest = NFProfileGetRequest.newBuilder().setFragmentSessionId("12312").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileGetRequest(nfProfileGetRequest).build()).build()).build();
            int result = nfProfileGetHelper.validate(nfMessage);
            assertEquals(Code.VALID, result);
        }
        {
            NFProfileGetRequest nfProfileGetRequest = NFProfileGetRequest.newBuilder().build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileGetRequest(nfProfileGetRequest).build()).build()).build();
            int result = nfProfileGetHelper.validate(nfMessage);
            assertEquals(Code.NFMESSAGE_PROTOCOL_ERROR, result);
        }
        {
            NFProfileGetRequest nfProfileGetRequest = NFProfileGetRequest.newBuilder().setFragmentSessionId("123").build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileGetRequest(nfProfileGetRequest).build()).build()).build();
            int result = nfProfileGetHelper.validate(nfMessage);
            assertEquals(Code.VALID, result);
        }
    }

    @Test
    public void createResponseTest() {
        {
            NFMessage nfMessage = nfProfileGetHelper.createResponse(Code.SUCCESS);
            assertNotNull(nfMessage);
            assertEquals(Code.SUCCESS, nfMessage.getResponse().getGetResponse().getNfProfileGetResponse().getCode());
        }
        {
            NFMessage nfMessage = nfProfileGetHelper.createResponse(Code.INTERNAL_ERROR);
            assertNotNull(nfMessage);
            assertEquals(Code.INTERNAL_ERROR, nfMessage.getResponse().getGetResponse().getNfProfileGetResponse().getCode());
        }
        {
            NFMessage nfMessage = nfProfileGetHelper.createResponse(Code.DATA_NOT_EXIST);
            assertNotNull(nfMessage);
            assertEquals(Code.DATA_NOT_EXIST, nfMessage.getResponse().getGetResponse().getNfProfileGetResponse().getCode());
        }
        {
            SearchResult searchResult = new SearchResult(false);
            searchResult.add("test");
            NFMessage nfMessage = nfProfileGetHelper.createResponse(searchResult);
            assertNotNull(nfMessage);
            assertEquals(Code.SUCCESS, nfMessage.getResponse().getGetResponse().getNfProfileGetResponse().getCode());
        }
        {
            FragmentResult fragmentResult = new FragmentResult();
            fragmentResult.setFragmentSessionID("21312");
            fragmentResult.setTotalNumber(100);
            fragmentResult.setTransmittedNumber(20);
            NFMessage nfMessage = nfProfileGetHelper.createResponse(fragmentResult);
            assertNotNull(nfMessage);
            assertEquals(Code.SUCCESS, nfMessage.getResponse().getGetResponse().getNfProfileGetResponse().getCode());
        }
    }
}
