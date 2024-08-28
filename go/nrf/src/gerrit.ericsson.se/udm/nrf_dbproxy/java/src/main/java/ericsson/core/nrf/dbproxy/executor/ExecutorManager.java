package ericsson.core.nrf.dbproxy.executor;

import ericsson.core.nrf.dbproxy.executor.cachenfprofile.CacheNFProfileGetExecutor;
import ericsson.core.nrf.dbproxy.executor.cachenfprofile.CacheNFProfilePutExecutor;
import ericsson.core.nrf.dbproxy.executor.groupprofile.*;
import ericsson.core.nrf.dbproxy.executor.nfprofile.*;
import ericsson.core.nrf.dbproxy.executor.nrfaddress.*;
import ericsson.core.nrf.dbproxy.executor.gpsiprofile.*;
import ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfileDeleteExecutor;
import ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfileGetExecutor;
import ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfilePutExecutor;
import ericsson.core.nrf.dbproxy.executor.imsiprefixprofile.ImsiprefixProfileGetExecutor;
import ericsson.core.nrf.dbproxy.executor.gpsiprefixprofile.GpsiprefixProfileGetExecutor;
import ericsson.core.nrf.dbproxy.executor.protocolerror.ProtocolErrorExecutor;
import ericsson.core.nrf.dbproxy.executor.subscription.*;

import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFRequest;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PutRequest;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetRequest;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.DelRequest;

public class ExecutorManager
{

    private static ExecutorManager instance = null;

    private ExecutorManager() { }

    public static synchronized ExecutorManager getInstance()
    {
        if(null == instance) {
            instance = new ExecutorManager();
        }
        return instance;
    }

    public Executor getExecutor(NFMessage request)
    {
        if(!request.hasRequest()) {
            return ProtocolErrorExecutor.getInstance();
        }

        NFRequest nf_request = request.getRequest();

        if(nf_request.hasPutRequest()) {

            PutRequest put_request = nf_request.getPutRequest();

            if(put_request.hasNfProfilePutRequest()) {
                return NFProfilePutExecutor.getInstance();
            } else if(put_request.hasSubscriptionPutRequest()) {
                return SubscriptionPutExecutor.getInstance();
            } else if(put_request.hasNrfAddressPutRequest()) {
                return NRFAddressPutExecutor.getInstance();
            } else if(put_request.hasGroupProfilePutRequest()) {
                return GroupProfilePutExecutor.getInstance();
            } else if(put_request.hasNrfProfilePutRequest()) {
                return NRFProfilePutExecutor.getInstance();
            } else if(put_request.hasGpsiProfilePutRequest()) {
                return GpsiProfilePutExecutor.getInstance();
            } else if(put_request.hasCacheNfProfilePutRequest()) {
                return CacheNFProfilePutExecutor.getInstance();
            } else {
                return ProtocolErrorExecutor.getInstance();
            }
        } else if(nf_request.hasGetRequest()) {

            GetRequest get_request = nf_request.getGetRequest();

            if(get_request.hasNfProfileGetRequest()) {
                return NFProfileGetExecutor.getInstance();
            } else if(get_request.hasSubscriptionGetRequest()) {
                return SubscriptionGetExecutor.getInstance();
            } else if(get_request.hasNrfAddressGetRequest()) {
                return NRFAddressGetExecutor.getInstance();
            } else if(get_request.hasGroupProfileGetRequest()) {
                return GroupProfileGetExecutor.getInstance();
            } else if(get_request.hasNrfProfileGetRequest()) {
                return NRFProfileGetExecutor.getInstance();
            } else if(get_request.hasImsiprefixProfileGetRequest()) {
                return ImsiprefixProfileGetExecutor.getInstance();
            } else if(get_request.hasGpsiProfileGetRequest()) {
                return GpsiProfileGetExecutor.getInstance();
            } else if(get_request.hasGpsiprefixProfileGetRequest()) {
                return GpsiprefixProfileGetExecutor.getInstance();
            } else if(get_request.hasNfProfileCountGetRequest()) {
                return NFProfileCountGetExecutor.getInstance();
            } else if(get_request.hasCacheNfProfileGetRequest()) {
                return CacheNFProfileGetExecutor.getInstance();
            } else {
                return ProtocolErrorExecutor.getInstance();
            }
        } else if(nf_request.hasDelRequest()) {

            DelRequest del_request = nf_request.getDelRequest();

            if(del_request.hasNfProfileDelRequest()) {
                return NFProfileDeleteExecutor.getInstance();
            } else if(del_request.hasSubscriptionDelRequest()) {
                return SubscriptionDeleteExecutor.getInstance();
            } else if(del_request.hasNrfAddressDelRequest()) {
                return NRFAddressDeleteExecutor.getInstance();
            } else if(del_request.hasGroupProfileDelRequest()) {
                return GroupProfileDeleteExecutor.getInstance();
            } else if(del_request.hasNrfProfileDelRequest()) {
                return NRFProfileDeleteExecutor.getInstance();
            } else if(del_request.hasGpsiProfileDelRequest()) {
                return GpsiProfileDeleteExecutor.getInstance();
            } else {
                return ProtocolErrorExecutor.getInstance();
            }
        } else {
            return ProtocolErrorExecutor.getInstance();
        }
    }
}
