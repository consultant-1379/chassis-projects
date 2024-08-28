package ericsson.core.nrf.dbproxy.grpc;

import ericsson.core.nrf.dbproxy.functionservice.RemoteCacheMonitorThread;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import java.util.ArrayList;
import io.grpc.stub.StreamObserver;
import java.util.Map.Entry;
import java.util.Set;
import javax.print.attribute.standard.PDLOverrideSupported;
import org.apache.geode.cache.query.internal.StructImpl;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.geode.cache.Region;
import org.apache.geode.pdx.JSONFormatter;
import org.apache.geode.pdx.PdxInstance;
import org.apache.geode.cache.query.SelectResults;

import ericsson.core.nrf.dbproxy.clientcache.schema.GpsiprefixProfiles;
import ericsson.core.nrf.dbproxy.clientcache.schema.ImsiprefixProfiles;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;
import ericsson.core.nrf.dbproxy.executor.ExecutorManager;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.*;


public class NFDataManagementServiceImpl extends NFDataManagementServiceGrpc.NFDataManagementServiceImplBase {

    private static final Logger logger = LogManager.getLogger(NFDataManagementServiceImpl.class);

    @Override
    public void execute(NFMessage request, StreamObserver<NFMessage> responseObserver) {
        NFMessage response;
        if (ClientCacheService.getInstance().isAvailable() == false) {
            logger.warn("ClientCacheService is NOT initialized, KVDB access is NOT available");
            response = request;
        } else {
            Executor executor = ExecutorManager.getInstance().getExecutor(request);
            response = executor.process(request);
        }

        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }

  @Override
    public void transferParameter(ParaRequest request, StreamObserver<ParaResponse> responseObserver){
			long arrivalTime = System.currentTimeMillis();
			int code = Code.CREATED;
			if (request.getParameterValue().equals("")){
				code = Code.BAD_REQUEST;
			} else if (request.getParameterName().equals("local-cache-capacity")) {
				String value = request.getParameterValue();
				RemoteCacheMonitorThread.getInstance().setCapacity(Integer.parseInt(value));
			}

			long departureTime = System.currentTimeMillis();
			TraceInfo traceInfo = TraceInfo.newBuilder().setArrivalTime(arrivalTime).setDepartureTime(departureTime).build();

			ParaResponse response;
      if(request.getTraceEnable())
        response = ParaResponse.newBuilder().setCode(code).setTraceInfo(traceInfo).build();
      else
        response = ParaResponse.newBuilder().setCode(code).build();

			responseObserver.onNext(response);
			responseObserver.onCompleted();
		}
  @Override
    public void insert(InsertRequest request, StreamObserver<InsertResponse> responseObserver)
    {
      int code;
      if (request.getItemCount() == 0) {
        code = Code.BAD_REQUEST;
      } else {
        String region_name = request.getRegionName();
        KVItem item = request.getItem(0);
        String key = item.getKey();
        String value = item.getValue();
        logger.debug("insert to region: " + region_name + " key: " + key + " value: " + value);
        code = ClientCacheService.getInstance().put(region_name, key, value);
      }

      InsertResponse response = InsertResponse.newBuilder().setCode(code).build();
      responseObserver.onNext(response);
      responseObserver.onCompleted();
    }
  @Override
    public void remove(RemoveRequest request, StreamObserver<RemoveResponse> responseObserver)
    {
			long arrivalTime = System.currentTimeMillis();
			int code = Code.SUCCESS;

			List<String> keys = request.getKeyList();
			if(keys.size() == 0)
			{
				code = Code.BAD_REQUEST;
			} else {
			  String region_name = request.getRegionName();
				if (region_name.equals("ericsson-nrf-cachenfprofiles")){
				  try {
            Set<String> cakeyS = ClientCacheService.getInstance().getRegion(region_name).keySetOnServer();
            Map<String, Object> caEntrys = ClientCacheService.getInstance().getRegion(region_name).getAll(cakeyS);
            for (Map.Entry<String, Object>entry : caEntrys.entrySet()){
                String key = entry.getKey();
                String from = (String) ((PdxInstance) entry.getValue()).getField("from");
                for (Iterator<String> it = keys.iterator(); it.hasNext();){
                  if (from.equals(it.next().toString())){
                    ClientCacheService.getInstance().delete(region_name, key);
                    break;
                  }
                }
            }
          } catch (Exception e){
				    logger.error("Remove CacheNFProfiles from : " + keys.toString() + "  fail.  " + e.toString());
				    code = Code.INTERNAL_ERROR;
          }
				} else {
          String key = request.getKey(0);
          logger.debug("Remove data from region: " + region_name + " key: " + key);
          code = ClientCacheService.getInstance().delete(region_name, key);
        }
			}

			long departureTime = System.currentTimeMillis();
			TraceInfo traceInfo = TraceInfo.newBuilder().setArrivalTime(arrivalTime).setDepartureTime(departureTime).build();

			RemoveResponse response;
			if(request.getTraceEnabled())
				response = RemoveResponse.newBuilder().setCode(code).setTraceInfo(traceInfo).build();
			else
				response = RemoveResponse.newBuilder().setCode(code).build();


			responseObserver.onNext(response);
			responseObserver.onCompleted();
    }

    private List<Long> getSearchImsiprefixList(Long imsi) {
        List<Long> imsiprefixList = new ArrayList<>();
        imsiprefixList.add(imsi);
        for (int i = 0; i < 10; i++) {
            imsi = imsi / 10;
            if (imsi < 10000) {
                break;
            }
            imsiprefixList.add(imsi);
        }
        return imsiprefixList;
    }

    private List<Long> getSearchGpsiprefixList(Long gpsi) {
        List<Long> gpsiprefixList = new ArrayList<>();
        gpsiprefixList.add(gpsi);
        for (int i = 0; i < 13; i++) {
            gpsi = gpsi / 10;
            if (gpsi < 10) {
                break;
            }
            gpsiprefixList.add(gpsi);
        }
        return gpsiprefixList;
    }

    @Override
    public void queryByKey(QueryRequest request, StreamObserver<QueryResponse> responseObserver) {
        long arrivalTime = System.currentTimeMillis();
        int code = Code.SUCCESS;

        String region_name = request.getRegionName();
        List<String> keys = request.getQueryList();
        List<String> values = new ArrayList<String>();
        if (keys.size() == 0) {
            code = Code.BAD_REQUEST;
        } else {
            if (region_name.equals("ericsson-nrf-gpsiprefixprofiles")) {
                List<Long> lkeys = getSearchGpsiprefixList(Long.parseLong(keys.get(0)));
                try {
                    Region region = ClientCacheService.getInstance().getRegion(region_name);
                    Map<String, Object> kvs = region.getAll(lkeys);
                    for (Map.Entry<String, Object> kv : kvs.entrySet()) {
                        GpsiprefixProfiles value = (GpsiprefixProfiles) kv.getValue();
                        if (value != null) {
                            Iterator iter = value.getValueInfo().keySet().iterator();
                            while (iter.hasNext()) {
                                String v = (String) (iter.next());
                                values.add(v);
                            }
                        }
                    }
                } catch (Exception e) {
                    logger.error(e.toString());
                    code = Code.INTERNAL_ERROR;
                }
            }else if (region_name.equals("ericsson-nrf-imsiprefixprofiles")) {
                List<Long> lkeys = getSearchImsiprefixList(Long.parseLong(keys.get(0)));
                try {
                    Region region = ClientCacheService.getInstance().getRegion(region_name);
                    Map<String, Object> kvs = region.getAll(lkeys);
                    for (Map.Entry<String, Object> kv : kvs.entrySet()) {
                        ImsiprefixProfiles value = (ImsiprefixProfiles) kv.getValue();
                        if (value != null) {
                            Iterator iter = value.getValueInfo().keySet().iterator();
                            while (iter.hasNext()) {
                                String v = (String) (iter.next());
                                values.add(v);
                            }
                        }
                    }
                } catch (Exception e) {
                    logger.error(e.toString());
                    code = Code.INTERNAL_ERROR;
                }
            } else if (region_name.equals("ericsson-nrf-nfprofiles")) {
                try {
                    Region region = ClientCacheService.getInstance().getRegion(region_name);
                    Map<String, Object> kvs = region.getAll(keys);
                    for (Map.Entry<String, Object> kv : kvs.entrySet()) {
                        Object value = kv.getValue();
                        if (value != null && value instanceof PdxInstance) {
                            values.add(JSONFormatter.toJSON((PdxInstance) value));
                        }
                    }
                } catch (Exception e) {
                    logger.error(e.toString());
                    code = Code.INTERNAL_ERROR;
                }
            } else {
                try {
                    Region region = ClientCacheService.getInstance().getRegion(region_name);
                    Map<String, Object> kvs = region.getAll(keys);
                    for (Map.Entry<String, Object> kv : kvs.entrySet()) {
                        Object value = kv.getValue();
                        if (value != null) {
                            values.add((String)(value));
                        }
                    }
                } catch (Exception e) {
                    logger.error(e.toString());
                    code = Code.INTERNAL_ERROR;
                }
            }
        }
        long departureTime = System.currentTimeMillis();
        TraceInfo traceInfo = TraceInfo.newBuilder().setArrivalTime(arrivalTime).setDepartureTime(departureTime).build();

        logger.debug("queryByID cost time=" + (departureTime - arrivalTime) + "ms");
        QueryResponse response;
        if (code != Code.SUCCESS) {
            if (request.getTraceEnabled())
                response = QueryResponse.newBuilder().setCode(code).setTraceInfo(traceInfo).build();
            else
                response = QueryResponse.newBuilder().setCode(code).build();
            responseObserver.onNext(response);
        } else {
            if (FragmentUtil.isNeedFragment(values)) {
                List<QueryResponse> responseList = FragmentUtil.getFragmentResponse(code, request.getTraceEnabled(), traceInfo, values);
                logger.debug("message is too large, need fragment, fragment size=" + responseList.size());
                for (int i = 0; i < responseList.size(); i++) {
                    responseObserver.onNext(responseList.get(i));
                }
            } else {
                if (request.getTraceEnabled())
                    response = QueryResponse.newBuilder().setCode(code).addAllValue(values).setTraceInfo(traceInfo).build();
                else
                    response = QueryResponse.newBuilder().setCode(code).addAllValue(values).build();
                responseObserver.onNext(response);
            }
        }
        responseObserver.onCompleted();
    }

    @Override
    public void queryByFilter(QueryRequest request, StreamObserver<QueryResponse> responseObserver) {
        long arrivalTime = System.currentTimeMillis();
        int code = Code.SUCCESS;

        String region_name = request.getRegionName();
        List<String> filters = request.getQueryList();
        List<String> values = new ArrayList<String>();
        if (filters.size() == 0) {
            code = Code.BAD_REQUEST;
        } else {
            try {
            Region region = ClientCacheService.getInstance().getRegion(region_name);
                for (String query : filters) {
                    SelectResults<Object> results = (SelectResults<Object>) region.getRegionService().getQueryService().newQuery(query).execute();
                    for (Object obj : results) {
                        if (obj instanceof String)
                            values.add((String) obj);
                        else if (obj instanceof PdxInstance)
                            values.add(JSONFormatter.toJSON((PdxInstance) obj));
                        else if (obj instanceof StructImpl) {
                            values.add(obj.toString());
                        }
                    }
                }
            } catch (Exception e) {
                logger.error(e.toString());
                code = Code.INTERNAL_ERROR;
            }
        }
        long departureTime = System.currentTimeMillis();
        TraceInfo traceInfo = TraceInfo.newBuilder().setArrivalTime(arrivalTime).setDepartureTime(departureTime).build();

        logger.debug("queryByFilter cost time=" + (departureTime - arrivalTime) + "ms");
        QueryResponse response;
        if (code != Code.SUCCESS) {
            if (request.getTraceEnabled())
                response = QueryResponse.newBuilder().setCode(code).setTraceInfo(traceInfo).build();
            else
                response = QueryResponse.newBuilder().setCode(code).build();
            responseObserver.onNext(response);
        } else {
            if (FragmentUtil.isNeedFragment(values)) {
                List<QueryResponse> responseList = FragmentUtil.getFragmentResponse(code, request.getTraceEnabled(), traceInfo, values);
                logger.debug("message is too large, need fragment, fragment size=" + responseList.size());
                for (int i = 0; i < responseList.size(); i++) {
                    responseObserver.onNext(responseList.get(i));
                }
            } else {
                if (request.getTraceEnabled())
                    response = QueryResponse.newBuilder().setCode(code).addAllValue(values).setTraceInfo(traceInfo).build();
                else
                    response = QueryResponse.newBuilder().setCode(code).addAllValue(values).build();
                responseObserver.onNext(response);
            }
        }
        responseObserver.onCompleted();
    }
}
