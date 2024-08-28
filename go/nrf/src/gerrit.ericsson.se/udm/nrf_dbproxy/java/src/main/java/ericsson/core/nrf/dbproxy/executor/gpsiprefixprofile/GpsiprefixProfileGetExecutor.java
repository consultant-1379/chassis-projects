package ericsson.core.nrf.dbproxy.executor.gpsiprefixprofile;

import java.util.List;
import java.util.ArrayList;
import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileGetRequestProto.GpsiprefixProfileGetRequest;
import ericsson.core.nrf.dbproxy.helper.gpsiprefixprofile.GpsiprefixProfileGetHelper;

public class GpsiprefixProfileGetExecutor extends Executor
{
    public static final int GPSI_PREFIX_SEARCH_COUNTER = 14;
	public static final int GPSI_PREFIX_MIN_VALUE = 10;

    private static GpsiprefixProfileGetExecutor instance = null;

    private GpsiprefixProfileGetExecutor()
    {
        super(GpsiprefixProfileGetHelper.getInstance());
    }

    public static synchronized GpsiprefixProfileGetExecutor getInstance()
    {
        if(null == instance) {
            instance = new GpsiprefixProfileGetExecutor();
        }
        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        GpsiprefixProfileGetRequest get_request = request.getRequest().getGetRequest().getGpsiprefixProfileGetRequest();
        return ClientCacheService.getInstance().getAllByID(Code.GPSIPREFIXPROFILE_INDICE, getSearchGpsiprefixList(get_request.getSearchGpsi()));
    }
    private List<Long> getSearchGpsiprefixList(Long gpsi)
    {
        List<Long> gpsiprefixList=new ArrayList<>();
		gpsiprefixList.add(gpsi);
        for(int i = 0; i < GPSI_PREFIX_SEARCH_COUNTER-1; i++) {
            gpsi = gpsi/10;
            if (gpsi < GPSI_PREFIX_MIN_VALUE) {
                break;
            }
            gpsiprefixList.add(gpsi);
        }
		return gpsiprefixList;
    }
}