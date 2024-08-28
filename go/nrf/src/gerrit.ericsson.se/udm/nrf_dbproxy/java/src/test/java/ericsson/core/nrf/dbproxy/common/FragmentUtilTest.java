package ericsson.core.nrf.dbproxy.common;

import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.TraceInfo;
import org.apache.commons.lang.RandomStringUtils;
import org.junit.Assert;
import org.junit.Test;

import java.util.ArrayList;
import java.util.List;

public class FragmentUtilTest {
    @Test
    public void isNeedFragmentTest() {
        List<String> values = new ArrayList<>();
        for (int i=0;i<1600;i++) {
            values.add(RandomStringUtils.random(102400));
        }
        Assert.assertTrue(FragmentUtil.isNeedFragment(values));
        List<String> values2 = new ArrayList<>();
        for (int i=0;i<500;i++) {
            values2.add(RandomStringUtils.random(1024));
        }
        Assert.assertFalse(FragmentUtil.isNeedFragment(values2));
    }

    @Test
    public void getFragmentResponseTest() {
        List<String> values = new ArrayList<>();
        for (int i=0;i<1600;i++) {
            values.add(RandomStringUtils.random(102400));
        }
        Assert.assertTrue(FragmentUtil.isNeedFragment(values));
        TraceInfo traceInfo = TraceInfo.newBuilder().build();
        List<QueryResponse> result = FragmentUtil.getFragmentResponse(2000, false, traceInfo, values);
        Assert.assertEquals(54, result.size());
        int count=0;
        for (int i=0;i<result.size();i++){
            count+=result.get(i).getValueList().size();
        }
        Assert.assertEquals(1600, count);
    }
}
