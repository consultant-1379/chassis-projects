package ericsson.core.nrf.dbproxy.executor.imsiprefixprofile;

import java.util.List;
import java.util.ArrayList;
import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.imsiprefixprofile.ImsiprefixProfileGetRequestProto.ImsiprefixProfileGetRequest;
import ericsson.core.nrf.dbproxy.helper.imsiprefixprofile.ImsiprefixProfileGetHelper;

public class ImsiprefixProfileGetExecutor extends Executor
{
    public static final int IMSI_PREFIX_SEARCH_COUNTER = 11;
    public static final int IMSI_PREFIX_MIN_VALUE = 10000;


    private static ImsiprefixProfileGetExecutor instance = null;

    private ImsiprefixProfileGetExecutor()
    {
        super(ImsiprefixProfileGetHelper.getInstance());
    }

    public static synchronized ImsiprefixProfileGetExecutor getInstance()
    {
        if(null == instance) {
            instance = new ImsiprefixProfileGetExecutor();
        }
        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        ImsiprefixProfileGetRequest get_request = request.getRequest().getGetRequest().getImsiprefixProfileGetRequest();
        return ClientCacheService.getInstance().getAllByID(Code.IMSIPREFIXPROFILE_INDICE, getSearchImsiprefixList(get_request.getSearchImsi()));
    }
    private List<Long> getSearchImsiprefixList(Long imsi)
    {
        List<Long> imsiprefixList=new ArrayList<>();
        imsiprefixList.add(imsi);
        for(int i = 0; i < IMSI_PREFIX_SEARCH_COUNTER-1; i++) {
            imsi = imsi/10;
            if (imsi < IMSI_PREFIX_MIN_VALUE) {
                break;
            }
            imsiprefixList.add(imsi);
        }
        return imsiprefixList;
    }
}