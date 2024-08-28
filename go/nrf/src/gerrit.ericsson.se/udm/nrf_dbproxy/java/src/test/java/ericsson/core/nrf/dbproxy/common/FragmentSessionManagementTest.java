package ericsson.core.nrf.dbproxy.common;

import ericsson.core.nrf.dbproxy.clientcache.schema.GroupProfile;
import org.junit.After;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;

public class FragmentSessionManagementTest {
    FragmentSessionManagement fragmentSessionManagement;
    @Before
    public void setUp() throws Exception {
        fragmentSessionManagement = FragmentSessionManagement.getInstance();
    }

    @After
    public void tearDown() throws Exception {
    }

    @Test
    public void put(){

	{
        	Assert.assertEquals(false, fragmentSessionManagement.put(null, 0));
	}

	{
        	FragmentResult fragmentResult = new FragmentResult();
        	Assert.assertEquals(true, fragmentSessionManagement.put(fragmentResult, 0));
	}

	{
        	FragmentResult fragmentResult = new FragmentResult();
		for(int i = 0; i < 100; i++)
		{
			fragmentResult.add(new Object());
		}
        	Assert.assertEquals(true, fragmentSessionManagement.put(fragmentResult, 100));
	}

	{
        	FragmentResult fragmentResult = new FragmentResult();
		for(int i = 0; i < 100; i++)
		{
			fragmentResult.add(new Object());
		}
        	Assert.assertEquals(false, fragmentSessionManagement.put(fragmentResult, 101));
	}
    }

    @Test
    public void get(){

 	FragmentResult fragmentResult = new FragmentResult();
	for(int i = 0; i < 101; i++)
	{
		fragmentResult.add(new GroupProfile());
	}
        Assert.assertEquals(true, fragmentSessionManagement.put(fragmentResult, 50));
	FragmentResult result = (FragmentResult)fragmentSessionManagement.get(Code.GROUPPROFILE_INDICE, fragmentResult.getFragmentSessionID());
	Assert.assertEquals(result.getFragmentSessionID(), fragmentResult.getFragmentSessionID());
	Assert.assertEquals(result.getTotalNumber(), fragmentResult.getTotalNumber());
	Assert.assertEquals(101, result.getTotalNumber());
	Assert.assertEquals(result.getTransmittedNumber(), fragmentResult.getTransmittedNumber());
	Assert.assertEquals(result.getTransmittedNumber(), fragmentResult.getTotalNumber());
    }

    @Test
    public void validate() {

    }
}
