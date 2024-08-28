package ericsson.core.nrf.dbproxy.helper.gpsiprefixprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.clientcache.schema.GpsiprefixProfiles;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.FragmentResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.*;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import java.util.Map;
import java.util.HashMap;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

public class GpsiprefixProfileGetHelperTest {
    GpsiprefixProfileGetHelper gpsiprefixProfileGetHelper;
    @Before
    public void setUp() throws Exception {
        gpsiprefixProfileGetHelper = GpsiprefixProfileGetHelper.getInstance();
    }

    @After
    public void tearDown() throws Exception {
    }
    
    @Test
    public void validateTest() {
        {
            GpsiprefixProfileGetRequest gpsiprefixProfileGetRequest = GpsiprefixProfileGetRequest.newBuilder().setSearchGpsi(0).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGpsiprefixProfileGetRequest(gpsiprefixProfileGetRequest).build()).build()).build();
            int result = gpsiprefixProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_SEARCH_GPSI);
        }
        {
            GpsiprefixProfileGetRequest gpsiprefixProfileGetRequest = GpsiprefixProfileGetRequest.newBuilder().setSearchGpsi(1000000000000000L).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGpsiprefixProfileGetRequest(gpsiprefixProfileGetRequest).build()).build()).build();
            int result = gpsiprefixProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.SEARCH_GPSI_LENGTH_EXCEED_MAX);
        }
        {
            GpsiprefixProfileGetRequest gpsiprefixProfileGetRequest = GpsiprefixProfileGetRequest.newBuilder().setSearchGpsi(123456789).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setGpsiprefixProfileGetRequest(gpsiprefixProfileGetRequest).build()).build()).build();
            int result = gpsiprefixProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
    }
    @Test
    public void createResponseTest() {
        {
            NFMessage nfMessage = gpsiprefixProfileGetHelper.createResponse(Code.SUCCESS);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getGpsiprefixProfileGetResponse().getCode(), Code.SUCCESS);
        }
        {
            NFMessage nfMessage = gpsiprefixProfileGetHelper.createResponse(Code.INTERNAL_ERROR);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getGpsiprefixProfileGetResponse().getCode(), Code.INTERNAL_ERROR);
        }
        {
            NFMessage nfMessage = gpsiprefixProfileGetHelper.createResponse(Code.DATA_NOT_EXIST);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getGpsiprefixProfileGetResponse().getCode(), Code.DATA_NOT_EXIST);
        }
        {
            SearchResult searchResult = new SearchResult(false);
            for (long i = 0L; i < 20L; i++) {
				Map<Long,GpsiprefixProfiles>  search_result_map = new HashMap<Long,GpsiprefixProfiles>();
                GpsiprefixProfiles gpsiprefixProfile = new GpsiprefixProfiles();
                gpsiprefixProfile.setGpsiprefix(10000L + i);
                gpsiprefixProfile.addValueInfo("15_gid_gid01_UDM");
                search_result_map.put(10000L+i,gpsiprefixProfile);
                searchResult.add(search_result_map);
            }
            NFMessage nfMessage = gpsiprefixProfileGetHelper.createResponse(searchResult);
            assertEquals(nfMessage.getResponse().getGetResponse().getGpsiprefixProfileGetResponse().getValueInfoCount(), 20);
            assertEquals(nfMessage.getResponse().getGetResponse().getGpsiprefixProfileGetResponse().getCode(), Code.SUCCESS);
        }
    }
}
