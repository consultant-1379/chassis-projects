package ericsson.core.nrf.dbproxy.grpc;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.clientcache.schema.GpsiprefixProfiles;
import ericsson.core.nrf.dbproxy.clientcache.schema.ImsiprefixProfiles;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.FragmentUtil;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.executor.ExecutorManager;
import ericsson.core.nrf.dbproxy.executor.nrfprofile.NRFProfileProcesser;
import ericsson.core.nrf.dbproxy.functionservice.RemoteCacheMonitorThread;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertRequest;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.InsertResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.KVItem;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaRequest;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.ParaResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchRequest;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryRequest;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveRequest;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.RemoveResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.TraceInfo;
import io.grpc.stub.StreamObserver;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import java.util.Set;
import org.apache.geode.cache.Region;
import org.apache.geode.cache.query.SelectResults;
import org.apache.geode.cache.query.internal.StructImpl;
import org.apache.geode.pdx.JSONFormatter;
import org.apache.geode.pdx.PdxInstance;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;


public class NFDataManagementServiceImpl extends
    NFDataManagementServiceGrpc.NFDataManagementServiceImplBase {

  private static final Logger LOGGER = LogManager.getLogger(NFDataManagementServiceImpl.class);

  @Override
  public void execute(NFMessage request, StreamObserver<NFMessage> responseObserver) {
    NFMessage response;
    if (!ClientCacheService.getInstance().isAvailable()) {
      LOGGER.warn("ClientCacheService is NOT initialized, KVDB access is NOT available");
      response = request;
    } else {
      Executor executor = ExecutorManager.getInstance().getExecutor(request);
      response = executor.process(request);
    }

    responseObserver.onNext(response);
    responseObserver.onCompleted();
  }

  @Override
  public void transferParameter(ParaRequest request,
      StreamObserver<ParaResponse> responseObserver) {
    long arrivalTime = System.currentTimeMillis();
    int code = Code.CREATED;
    if (request.getParameterValue().equals("")) {
      code = Code.BAD_REQUEST;
    } else if (request.getParameterName().equals("local-cache-capacity")) {
      String value = request.getParameterValue();
      RemoteCacheMonitorThread.getInstance().setCapacity(Integer.parseInt(value));
    }

    long departureTime = System.currentTimeMillis();
    TraceInfo traceInfo = TraceInfo.newBuilder().setArrivalTime(arrivalTime)
        .setDepartureTime(departureTime).build();

    ParaResponse response;
    if (request.getTraceEnable()) {
      response = ParaResponse.newBuilder().setCode(code).setTraceInfo(traceInfo).build();
    } else {
      response = ParaResponse.newBuilder().setCode(code).build();
    }

    responseObserver.onNext(response);
    responseObserver.onCompleted();
  }

  @Override
  public void insert(InsertRequest request, StreamObserver<InsertResponse> responseObserver) {
    long arrivalTime = System.currentTimeMillis();
    int code;
    String regionName = request.getRegionName();
    if (request.getItemCount() == 0) {
      code = Code.BAD_REQUEST;
    } else if (Code.NFPROFILE_INDICE.equals(regionName)) {
      LOGGER.error("inster infterface don't support region ericsson-nrf-nfprofiles");
      code = Code.INTERNAL_ERROR;
    } else if (Code.NRFPROFILE_INDICE.equals(regionName)) {
      List<KVItem> kvItems = request.getItemList();
      //just support put one nrfprofile now
      code = NRFProfileProcesser.putNRFProfile(kvItems.get(0));
    } else {
      List<KVItem> kvItems = request.getItemList();
      if (kvItems.size() == 1) {
        code = ClientCacheService.getInstance()
            .put(regionName, kvItems.get(0).getKey(), kvItems.get(0).getValue());
      } else {
        Map<Object, Object> kvMap = new HashMap<>();
        for (KVItem kvItem : kvItems) {
          LOGGER.debug(
              "insert to region: " + regionName + " key: " + kvItem.getKey() + " value: " + kvItem
                  .getValue());
          kvMap.put(kvItem.getKey(), kvItem.getValue());
        }
        code = ClientCacheService.getInstance().putAll(regionName, kvMap);
      }
    }
    long departureTime = System.currentTimeMillis();
    TraceInfo traceInfo = TraceInfo.newBuilder().setArrivalTime(arrivalTime)
        .setDepartureTime(departureTime).build();

    InsertResponse.Builder builder = InsertResponse.newBuilder().setCode(code);
    if (request.getTraceEnabled()) {
      builder.setTraceInfo(traceInfo);
    }

    InsertResponse response = builder.build();
    responseObserver.onNext(response);
    responseObserver.onCompleted();
  }

  @Override
  public void remove(RemoveRequest request, StreamObserver<RemoveResponse> responseObserver) {
    long arrivalTime = System.currentTimeMillis();
    int code = Code.SUCCESS;

    List<String> keys = request.getKeyList();
    if (keys.size() == 0) {
      code = Code.BAD_REQUEST;
    } else {
      String regionName = request.getRegionName();
      if (regionName.equals(Code.CACHENFPROFILE_INDICE)) {
        try {
          Set<String> cakeyS = ClientCacheService.getInstance().getRegion(regionName)
              .keySetOnServer();
          Map<String, Object> caEntrys = ClientCacheService.getInstance().getRegion(regionName)
              .getAll(cakeyS);
          for (Map.Entry<String, Object> entry : caEntrys.entrySet()) {
            String key = entry.getKey();
            String from = (String) ((PdxInstance) entry.getValue()).getField("from");
            for (Iterator<String> it = keys.iterator(); it.hasNext();) {
              if (from.equals(it.next().toString())) {
                ClientCacheService.getInstance().delete(regionName, key);
                break;
              }
            }
          }
        } catch (Exception e) {
          LOGGER.error(
              "Remove CacheNFProfiles from : " + keys.toString() + "  fail.  " + e.toString());
          code = Code.INTERNAL_ERROR;
        }
      } else if (Code.NRFPROFILE_INDICE.equals(regionName)) {
        //just support delete one nrfprofile one time
        code = NRFProfileProcesser.deleteNRFProfile(keys.get(0));
      } else {
        if (keys.size() == 1) {
          code = ClientCacheService.getInstance().delete(regionName, keys.get(0));
        } else {
          code = ClientCacheService.getInstance().deleteAll(regionName, keys);
        }
      }
    }

    long departureTime = System.currentTimeMillis();
    TraceInfo traceInfo = TraceInfo.newBuilder().setArrivalTime(arrivalTime)
        .setDepartureTime(departureTime).build();

    RemoveResponse response;
    RemoveResponse.Builder builder;
    if (request.getTraceEnabled()) {
      builder = RemoveResponse.newBuilder().setCode(code).setTraceInfo(traceInfo);
    } else {
      builder = RemoveResponse.newBuilder().setCode(code);
    }

    response = builder.build();
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

    String regionName = request.getRegionName();
    List<String> keys = request.getQueryList();
    List<String> values = new ArrayList<String>();
    if (keys.size() == 0) {
      code = Code.BAD_REQUEST;
    } else {
      if (regionName.equals(Code.GPSIPREFIXPROFILE_INDICE)) {
        List<Long> lkeys = getSearchGpsiprefixList(Long.parseLong(keys.get(0)));
        try {
          Region region = ClientCacheService.getInstance().getRegion(regionName);
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
          LOGGER.error(e.toString());
          code = Code.INTERNAL_ERROR;
        }
      } else if (regionName.equals(Code.IMSIPREFIXPROFILE_INDICE)) {
        List<Long> lkeys = getSearchImsiprefixList(Long.parseLong(keys.get(0)));
        try {
          Region region = ClientCacheService.getInstance().getRegion(regionName);
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
          LOGGER.error(e.toString());
          code = Code.INTERNAL_ERROR;
        }
      } else if (Code.NFPROFILE_INDICE.equals(regionName) || Code.NRFPROFILE_INDICE
          .equals(regionName) || Code.REGIONNFINFO_INDICE.equals(regionName)) {
        try {
          Region region = ClientCacheService.getInstance().getRegion(regionName);
          Map<String, Object> kvs = region.getAll(keys);
          for (Map.Entry<String, Object> kv : kvs.entrySet()) {
            Object value = kv.getValue();
            if (value != null && value instanceof PdxInstance) {
              values.add(JSONFormatter.toJSON((PdxInstance) value));
            }
          }
        } catch (Exception e) {
          LOGGER.error(e.toString());
          code = Code.INTERNAL_ERROR;
        }
      } else {
        try {
          Region region = ClientCacheService.getInstance().getRegion(regionName);
          Map<String, Object> kvs = region.getAll(keys);
          for (Map.Entry<String, Object> kv : kvs.entrySet()) {
            Object value = kv.getValue();
            if (value != null) {
              values.add((String) value);
            }
          }
        } catch (Exception e) {
          LOGGER.error(e.toString());
          code = Code.INTERNAL_ERROR;
        }
      }
    }
    if (values.size() <= 0 && code == Code.SUCCESS) {
      code = Code.DATA_NOT_EXIST;
    }
    long departureTime = System.currentTimeMillis();
    TraceInfo traceInfo = TraceInfo.newBuilder().setArrivalTime(arrivalTime)
        .setDepartureTime(departureTime).build();

    QueryResponse response;
    if (code != Code.SUCCESS) {
      if (request.getTraceEnabled()) {
        response = QueryResponse.newBuilder().setCode(code).setTraceInfo(traceInfo).build();
      } else {
        response = QueryResponse.newBuilder().setCode(code).build();
      }
      responseObserver.onNext(response);
    } else {
      if (FragmentUtil.isNeedFragment(values)) {
        List<QueryResponse> responseList = FragmentUtil
            .getFragmentResponse(code, request.getTraceEnabled(), traceInfo, values);
        LOGGER.debug("message is too large, need fragment, fragment size=" + responseList.size());
        for (int i = 0; i < responseList.size(); i++) {
          responseObserver.onNext(responseList.get(i));
        }
      } else {
        if (request.getTraceEnabled()) {
          response = QueryResponse.newBuilder().setCode(code).addAllValue(values)
              .setTraceInfo(traceInfo).build();
        } else {
          response = QueryResponse.newBuilder().setCode(code).addAllValue(values).build();
        }
        responseObserver.onNext(response);
      }
    }
    responseObserver.onCompleted();
  }

  @Override
  public void queryByFilter(QueryRequest request, StreamObserver<QueryResponse> responseObserver) {
    long arrivalTime = System.currentTimeMillis();
    int code = Code.SUCCESS;

    String regionName = request.getRegionName();
    List<String> filters = request.getQueryList();
    List<String> values = new ArrayList<String>();
    if (filters.size() == 0) {
      code = Code.BAD_REQUEST;
    } else {
      if (regionName.equals(Code.NRFPROFILE_INDICE+Code.REGIONNFINFO_INDICE)) {
        code = getMatchAllNRFProfile(filters, values);
      } else {
        try {
          Region region = ClientCacheService.getInstance().getRegion(regionName);
          for (String query : filters) {
            LOGGER.debug("queryByFilter oql:" + query);

            SelectResults<Object> results = (SelectResults<Object>) region.getRegionService()
                .getQueryService().newQuery(query).execute();
            parseSelectResult(results, values);
          }
        } catch (Exception e) {
          LOGGER.error(e.toString());
          code = Code.INTERNAL_ERROR;
        }
      }

      if (regionName.equals(Code.NFHELPER_INDICE) && values.size() > 0) {
        LOGGER.debug(
            "It's a query nfhelper case, need to query from nfprofile region to get nfInstanceId/profileUpdateTime pair!");
        try {
          Region region = ClientCacheService.getInstance().getRegion("ericsson-nrf-nfprofiles");
          Map<String, Object> kvs = region.getAll(values);
          values.clear();
          for (Map.Entry<String, Object> kv : kvs.entrySet()) {
            Object nfprofilePdxFormat = kv.getValue();
            if (nfprofilePdxFormat == null || !(nfprofilePdxFormat instanceof PdxInstance)) {
              continue;
            }
            String nfInstanceId = kv.getKey();
            Long profileUpdateTime = (Long) ((PdxInstance) nfprofilePdxFormat)
                .getField("profileUpdateTime");
            StringBuilder sb = new StringBuilder("{\"nfInstanceId\":\"" + nfInstanceId + "\",");
            sb.append("\"profileUpdateTime\":\"" + Long.toString(profileUpdateTime) + "\"}");
            LOGGER.debug("The result is: " + sb.toString());
            values.add(sb.toString());
          }
        } catch (Exception e) {
          LOGGER.error(e.toString());
          code = Code.INTERNAL_ERROR;
        }
      }
    }

    long departureTime = System.currentTimeMillis();
    TraceInfo traceInfo = TraceInfo.newBuilder().setArrivalTime(arrivalTime)
        .setDepartureTime(departureTime).build();

    QueryResponse response;
    if (code != Code.SUCCESS) {
      if (request.getTraceEnabled()) {
        response = QueryResponse.newBuilder().setCode(code).setTraceInfo(traceInfo).build();
      } else {
        response = QueryResponse.newBuilder().setCode(code).build();
      }
      responseObserver.onNext(response);
    } else {
      if (FragmentUtil.isNeedFragment(values)) {
        List<QueryResponse> responseList = FragmentUtil
            .getFragmentResponse(code, request.getTraceEnabled(), traceInfo, values);
        LOGGER.debug("message is too large, need fragment, fragment size=" + responseList.size());
        for (int i = 0; i < responseList.size(); i++) {
          responseObserver.onNext(responseList.get(i));
        }
      } else {
        if (request.getTraceEnabled()) {
          response = QueryResponse.newBuilder().setCode(code).addAllValue(values)
              .setTraceInfo(traceInfo).build();
        } else {
          response = QueryResponse.newBuilder().setCode(code).addAllValue(values).build();
        }
        responseObserver.onNext(response);
      }
    }
    responseObserver.onCompleted();
  }

  private int getMatchAllNRFProfile(List<String> filters, List<String> values) {
    try {
      Set<String> nrfInsts = ClientCacheService.getInstance().getRegion(Code.NRFPROFILE_INDICE).keySetOnServer();
      SelectResults<Object> regionNFInfo = (SelectResults<Object>) ClientCacheService.getInstance()
          .getRegion(Code.REGIONNFINFO_INDICE).getRegionService().getQueryService()
          .newQuery(filters.get(0)).execute();
      LOGGER.debug("Query OQL: " + filters.get(0));
      List<String> keys = mergeSelectResult(nrfInsts, regionNFInfo);
      if (keys.size() > 0) {
        Map<String, Object> kvs = ClientCacheService.getInstance().getRegion(Code.NRFPROFILE_INDICE)
            .getAll(keys);
        for (Map.Entry<String, Object> kv : kvs.entrySet()) {
          Object value = kv.getValue();
          if (value != null && value instanceof PdxInstance) {
            values.add(JSONFormatter.toJSON((PdxInstance) value));
          }
        }
      }
    } catch (Exception e) {
      LOGGER.debug(e.toString());
      return Code.INTERNAL_ERROR;
    }

    return Code.SUCCESS;
  }

  private List<String> mergeSelectResult(Set<String> nrfInsts,
      SelectResults<Object> regionNFInfo) {
    LOGGER.debug("Query Resoult: " + nrfInsts.toString() + "  \n" + regionNFInfo.toString());
    Map<String, Boolean> map = new HashMap<String, Boolean>();
    for (String inst : nrfInsts) {
      map.put(inst, true);
    }

    for (Object obj : regionNFInfo) {
      map.put((String) obj, false);
    }
    List<String> keys = new ArrayList<String>();
    Iterator<String> it = map.keySet().iterator();
    while (it.hasNext()) {
      String key = it.next();
      Boolean value = map.get(key);
      if (value.booleanValue()) {
        keys.add(key);
      }
    }

    return keys;
  }

  private void parseSelectResult(SelectResults<Object> results, List<String> values) {
    for (Object obj : results) {
      if (obj instanceof String) {
        values.add((String) obj);
      } else if (obj instanceof PdxInstance) {
        values.add(JSONFormatter.toJSON((PdxInstance) obj));
      } else if (obj instanceof StructImpl) {
        StringBuilder sb = new StringBuilder("{");
        int length = ((StructImpl) obj).getFieldNames().length;
        String[] names = ((StructImpl) obj).getFieldNames();
        Object[] objs = ((StructImpl) obj).getPdxFieldValues();
        for (int i = 0; i < length; i++) {
          if (objs[i] instanceof PdxInstance) {
            if (i > 0) {
              sb.append(",\"" + names[i] + "\":" + JSONFormatter.toJSON((PdxInstance) objs[i]));
            } else {
              sb.append("\"" + names[i] + "\":" + JSONFormatter.toJSON((PdxInstance) objs[i]));
            }

          } else {
            if (i > 0) {
              sb.append(",\"" + names[i] + "\":\"" + objs[i].toString() + "\"");
            } else {
              sb.append("\"" + names[i] + "\":\"" + objs[i].toString() + "\"");
            }
          }
        }
        sb.append("}");
        values.add(sb.toString());
      } else if (obj instanceof Integer) {
        values.add(String.valueOf(obj));
      } else {
        LOGGER.warn("no data type is matched");
      }
    }
  }

  @Override
  public void patchNrfProfile(PatchRequest request,
      StreamObserver<PatchResponse> responseObserver) {
    int code;
    if (!NRFProfileProcesser.validatePatchRequest(request)) {
      code = Code.BAD_REQUEST;
    } else {
      code = NRFProfileProcesser.processPatchRequest(request);
    }
    LOGGER.debug("patch nrfprofile return code ={}", code);
    PatchResponse.Builder builder = PatchResponse.newBuilder();
    builder.setCode(code);
    PatchResponse patchResponse = builder.build();
    responseObserver.onNext(patchResponse);
    responseObserver.onCompleted();
  }
}
