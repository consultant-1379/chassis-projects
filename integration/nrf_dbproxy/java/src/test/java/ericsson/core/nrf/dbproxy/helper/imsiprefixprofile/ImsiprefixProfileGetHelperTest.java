package ericsson.core.nrf.dbproxy.helper.imsiprefixprofile;

import com.google.protobuf.ByteString;
import ericsson.core.nrf.dbproxy.clientcache.schema.ImsiprefixProfiles;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.FragmentResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.imsiprefixprofile.ImsiprefixProfileGetRequestProto.*;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import java.util.Map;
import java.util.HashMap;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

public class ImsiprefixProfileGetHelperTest {
    ImsiprefixProfileGetHelper imsiprefixProfileGetHelper;
    @Before
    public void setUp() throws Exception {
        imsiprefixProfileGetHelper = ImsiprefixProfileGetHelper.getInstance();
    }

    @After
    public void tearDown() throws Exception {
    }
    
    @Test
    public void validateTest() {
        {
            ImsiprefixProfileGetRequest imsiprefixProfileGetRequest = ImsiprefixProfileGetRequest.newBuilder().setSearchImsi(0).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setImsiprefixProfileGetRequest(imsiprefixProfileGetRequest).build()).build()).build();
            int result = imsiprefixProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.EMPTY_SEARCH_IMSI);
        }
        {
            ImsiprefixProfileGetRequest imsiprefixProfileGetRequest = ImsiprefixProfileGetRequest.newBuilder().setSearchImsi(1000000000000000L).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setImsiprefixProfileGetRequest(imsiprefixProfileGetRequest).build()).build()).build();
            int result = imsiprefixProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.SEARCH_IMSI_LENGTH_EXCEED_MAX);
        }
        {
            ImsiprefixProfileGetRequest imsiprefixProfileGetRequest = ImsiprefixProfileGetRequest.newBuilder().setSearchImsi(123456789).build();
            NFMessage nfMessage = NFMessage.newBuilder().setRequest(NFRequest.newBuilder().setGetRequest(GetRequest.newBuilder().setImsiprefixProfileGetRequest(imsiprefixProfileGetRequest).build()).build()).build();
            int result = imsiprefixProfileGetHelper.validate(nfMessage);
            assertEquals(result, Code.VALID);
        }
    }
    @Test
    public void createResponseTest() {
        {
            NFMessage nfMessage = imsiprefixProfileGetHelper.createResponse(Code.SUCCESS);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getImsiprefixProfileGetResponse().getCode(), Code.SUCCESS);
        }
        {
            NFMessage nfMessage = imsiprefixProfileGetHelper.createResponse(Code.INTERNAL_ERROR);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getImsiprefixProfileGetResponse().getCode(), Code.INTERNAL_ERROR);
        }
        {
            NFMessage nfMessage = imsiprefixProfileGetHelper.createResponse(Code.DATA_NOT_EXIST);
            assertNotNull(nfMessage);
            assertEquals(nfMessage.getResponse().getGetResponse().getImsiprefixProfileGetResponse().getCode(), Code.DATA_NOT_EXIST);
        }
        {
            SearchResult searchResult = new SearchResult(false);
            for (long i = 0L; i < 20L; i++) {
				Map<Long,ImsiprefixProfiles>  search_result_map = new HashMap<Long,ImsiprefixProfiles>();
                ImsiprefixProfiles imsiprefixProfile = new ImsiprefixProfiles();
                imsiprefixProfile.setImsiprefix(10000L + i);
                imsiprefixProfile.addValueInfo("15_gid_gid01_UDM");
                search_result_map.put(10000L+i,imsiprefixProfile);
                searchResult.add(search_result_map);
            }
            NFMessage nfMessage = imsiprefixProfileGetHelper.createResponse(searchResult);
            assertEquals(nfMessage.getResponse().getGetResponse().getImsiprefixProfileGetResponse().getValueInfoCount(), 20);
            assertEquals(nfMessage.getResponse().getGetResponse().getImsiprefixProfileGetResponse().getCode(), Code.SUCCESS);
        }
    }
}
