package ericsson.core.nrf.dbproxy.helper.nfprofile;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetRequest;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileFilterProto.NFProfileFilter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileFilterProto.Range;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileCountGetRequestProto.NFProfileCountGetRequest;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

public class NFProfileCountGetHelperTest {
    NFProfileCountGetHelper nfProfileCountGetHelper;

    @Before
    public void setUp() throws Exception {
        nfProfileCountGetHelper = NFProfileCountGetHelper.getInstance();
        assertNotNull(nfProfileCountGetHelper);
    }

    @After
    public void tearDown() throws Exception {
    }

    @Test
    public void validateTest() {
        {
            NFProfileFilter filter = NFProfileFilter.newBuilder().setExpiredTimeRange(Range.newBuilder().setStart(123).setEnd(1234).build()).build();
            NFProfileCountGetRequest nfProfileCountGetRequest = NFProfileCountGetRequest.newBuilder().setFilter(filter).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileCountGetRequest(nfProfileCountGetRequest).build()).build()).build();
            int result = nfProfileCountGetHelper.validate(nfMessage);
            assertEquals(Code.VALID, result);
        }
        {
            NFProfileFilter filter = NFProfileFilter.newBuilder().setLastUpdateTimeRange(Range.newBuilder().setStart(123).setEnd(1234).build()).build();
            NFProfileCountGetRequest nfProfileCountGetRequest = NFProfileCountGetRequest.newBuilder().setFilter(filter).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileCountGetRequest(nfProfileCountGetRequest).build()).build()).build();
            int result = nfProfileCountGetHelper.validate(nfMessage);
            assertEquals(Code.VALID, result);
        }
        {
            NFProfileFilter filter = NFProfileFilter.newBuilder().setLastUpdateTimeRange(Range.newBuilder().setStart(12322).setEnd(1234).build()).build();
            NFProfileCountGetRequest nfProfileCountGetRequest = NFProfileCountGetRequest.newBuilder().setFilter(filter).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileCountGetRequest(nfProfileCountGetRequest).build()).build()).build();
            int result = nfProfileCountGetHelper.validate(nfMessage);
            assertEquals(Code.INVALID_RANGE, result);
        }
        {
            NFProfileFilter filter = NFProfileFilter.newBuilder().setProvisioned(3).build();
            NFProfileCountGetRequest nfProfileCountGetRequest = NFProfileCountGetRequest.newBuilder().setFilter(filter).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileCountGetRequest(nfProfileCountGetRequest).build()).build()).build();
            int result = nfProfileCountGetHelper.validate(nfMessage);
            assertEquals(Code.INVALID_PROVISIONED, result);

        }
        {
            NFProfileCountGetRequest nfProfileCountGetRequest = NFProfileCountGetRequest.newBuilder().build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setNfProfileCountGetRequest(nfProfileCountGetRequest).build()).build()).build();
            int result = nfProfileCountGetHelper.validate(nfMessage);
            assertEquals(Code.NFMESSAGE_PROTOCOL_ERROR, result);
        }
    }

    @Test
    public void createResponseTest() {
        {
            NFMessage nfMessage = nfProfileCountGetHelper.createResponse(Code.SUCCESS);
            assertNotNull(nfMessage);
            assertEquals(Code.SUCCESS, nfMessage.getResponse().getGetResponse().getNfProfileCountGetResponse().getCode());
        }
        {
            NFMessage nfMessage = nfProfileCountGetHelper.createResponse(Code.INTERNAL_ERROR);
            assertNotNull(nfMessage);
            assertEquals(Code.INTERNAL_ERROR, nfMessage.getResponse().getGetResponse().getNfProfileCountGetResponse().getCode());
        }
        {
            NFMessage nfMessage = nfProfileCountGetHelper.createResponse(Code.DATA_NOT_EXIST);
            assertNotNull(nfMessage);
            assertEquals(Code.DATA_NOT_EXIST, nfMessage.getResponse().getGetResponse().getNfProfileCountGetResponse().getCode());
        }
        {
            SearchResult searchResult = new SearchResult(false);
            searchResult.add(100);
            NFMessage nfMessage = nfProfileCountGetHelper.createResponse(searchResult);
            assertNotNull(nfMessage);
            assertEquals(Code.SUCCESS, nfMessage.getResponse().getGetResponse().getNfProfileCountGetResponse().getCode());
        }
    }
}
