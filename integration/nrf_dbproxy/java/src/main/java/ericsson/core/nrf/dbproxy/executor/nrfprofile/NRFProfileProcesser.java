package ericsson.core.nrf.dbproxy.executor.nrfprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.KVItem;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchBody;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchItem;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PatchRequest;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import org.apache.geode.cache.CacheTransactionManager;
import org.apache.geode.cache.CommitConflictException;
import org.apache.geode.cache.Region;
import org.apache.geode.cache.query.SelectResults;
import org.apache.geode.pdx.JSONFormatter;
import org.apache.geode.pdx.PdxInstance;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.json.JSONArray;
import org.json.JSONObject;
import org.json.JSONTokener;


public class NRFProfileProcesser {

  private static final Logger LOGGER = LogManager.getLogger(NRFProfileProcesser.class);
  private static final List<String> PATCHNFTYPES = getPatchNfTypes(false);
  private static final List<String> ALLPATCHTYPES = getPatchNfTypes(true);
  private static final String GET_NFINFOID_SQL = "select value.body.nfInstanceId from /ericsson-nrf-regionnfinfo.entrySet where value.body.nrfInstanceId=";
  private static final String NRF_REGION_NAME = "ericsson-nrf-nrfprofiles";
  private static final String NFINFO_REGION_NAME = "ericsson-nrf-regionnfinfo";
  private static final String NRFINFO = "nrfInfo";

  private NRFProfileProcesser() {
    throw new IllegalStateException("Utility class");
  }

  public static int putNRFProfile(KVItem kvItems) {
    Map<Object, Object> nfInfoMap = new HashMap<>();
    Region nrfProfileRegion = ClientCacheService.getInstance().getRegion(NRF_REGION_NAME);
    Region nfInfoRegion = ClientCacheService.getInstance().getRegion(NFINFO_REGION_NAME);
    String nrfInstanceId = kvItems.getKey();
    String nrfProfile = kvItems.getValue();
    JSONObject nrfJsonObject = new JSONObject(nrfProfile);
    JSONArray nfInfoArray = nrfJsonObject.getJSONArray(NRFINFO);
    for (int i = 0; i < nfInfoArray.length(); i++) {
      JSONObject infoObject = nfInfoArray.getJSONObject(i);
      String nfInstanceId = infoObject.getJSONObject("body").getString("nfInstanceId");
      PdxInstance pdxInstance = JSONFormatter.fromJSON(infoObject.toString());
      nfInfoMap.put(nfInstanceId, pdxInstance);
    }
    nrfJsonObject.remove(NRFINFO);
    PdxInstance nrfProfileToDB = JSONFormatter.fromJSON(nrfJsonObject.toString());
    LOGGER.debug("nrfProfile stored in db = {}", JSONFormatter.toJSON(nrfProfileToDB));

    int code = Code.INTERNAL_ERROR;
    CacheTransactionManager txManager = ClientCacheService.getInstance()
        .getCacheTransactionManager();
    boolean retryTransaction = false;
    int retryTime = 0;
    do {
      try {
        retryTime++;
        if (retryTime > 3) {
          break;
        }
        txManager.begin();
        SelectResults<String> results = (SelectResults<String>) nfInfoRegion.getRegionService()
            .getQueryService().newQuery(GET_NFINFOID_SQL + "'" + nrfInstanceId + "'").execute();
        nfInfoRegion.removeAll(results);
        nrfProfileRegion.remove(nrfInstanceId);
        nrfProfileRegion.put(nrfInstanceId, nrfProfileToDB);
        nfInfoRegion.putAll(nfInfoMap);
        txManager.commit();
        retryTransaction = false;
        code = Code.CREATED;
      } catch (CommitConflictException conflictException) {
        LOGGER.error(conflictException.toString());
        retryTransaction = true;
      } catch (Exception e) {
        LOGGER.error(e.toString());
      } finally {
        if (txManager.exists()) {
          txManager.rollback();
        }
      }
    } while (retryTransaction);
    return code;
  }

  public static int deleteNRFProfile(String nrfInstanceId) {
    Region nrfProfileRegion = ClientCacheService.getInstance().getRegion(NRF_REGION_NAME);
    Region nfInfoRegion = ClientCacheService.getInstance().getRegion(NFINFO_REGION_NAME);
    int code = Code.INTERNAL_ERROR;
    CacheTransactionManager txManager = ClientCacheService.getInstance()
        .getCacheTransactionManager();
    boolean retryTransaction = false;
    int retryTime = 0;
    do {
      try {
        retryTime++;
        if (retryTime > 3) {
          break;
        }
        txManager.begin();
        SelectResults<String> results = (SelectResults<String>) nfInfoRegion.getRegionService()
            .getQueryService().newQuery(GET_NFINFOID_SQL + "'" + nrfInstanceId + "'").execute();
        nfInfoRegion.removeAll(results);
        nrfProfileRegion.remove(nrfInstanceId);
        txManager.commit();
        retryTransaction = false;
        code = Code.SUCCESS;
      } catch (CommitConflictException conflictException) {
        LOGGER.error(conflictException.toString());
        retryTransaction = true;
      } catch (Exception e) {
        LOGGER.error(e.toString());
      } finally {
        if (txManager.exists()) {
          txManager.rollback();
        }
      }
    } while (retryTransaction);
    return code;
  }

  public static boolean validatePatchRequest(PatchRequest patchRequest) {
    LOGGER.debug("patchRequest={}", patchRequestToString(patchRequest));
    if (patchRequest.getAllPatchCount() == 0) {
      return false;
    }
    for (PatchBody patchBody : patchRequest.getAllPatchList()) {
      String patchType = patchBody.getPatchType();
      int patchOperation = patchBody.getOperation();
      if (!ALLPATCHTYPES.contains(patchType)) {
        LOGGER.error("unknow patch nftype=" + patchType);
        return false;
      }
      if (PATCHNFTYPES.contains(patchType) && patchBody.getPatchNrfInstId().isEmpty()) {
        LOGGER.error("nrfIntsanceId should not be null when patch type=" + patchType);
        return false;
      }
      if (patchOperation != Code.PATCH_ADD && patchOperation != Code.PATCH_REPLACE
          && patchOperation != Code.PATCH_REMOVE) {
        LOGGER.error("unknow patch operation=" + patchOperation);
        return false;
      }
      List<PatchItem> patchItems = patchBody.getPatchItemList();
      for (PatchItem patchItem : patchItems) {
        if (patchItem.getPatchKey().isEmpty()) {
          LOGGER.error("empty patch key");
          return false;
        }
      }
    }
    return true;
  }

  public static String patchRequestToString(PatchRequest patchRequest) {
    StringBuilder sb = new StringBuilder();
    sb.append("[");
    int index = 0;
    for (PatchBody patchBody : patchRequest.getAllPatchList()) {
      if (index > 0) {
        sb.append(",{");
      }
      sb.append("{");
      sb.append("patchType=" + patchBody.getPatchType());
      sb.append(",patchNrfInstanceId=" + patchBody.getPatchNrfInstId());
      sb.append(",patchOperation=" + patchBody.getOperation());
      sb.append(",patchItem=[");
      List<PatchItem> patchItems = patchBody.getPatchItemList();
      for (PatchItem patchItem : patchItems) {
        sb.append("{patchKey=" + patchItem.getPatchKey());
        sb.append(",patchValue=" + patchItem.getPatchValue() + "}");
      }
      sb.append("]");
      sb.append("}");
      sb.append("}");
      index++;
    }
    sb.append("]");
    return sb.toString();
  }

  public static int processPatchRequest(PatchRequest patchRequest) {
    int code = Code.INTERNAL_ERROR;
    CacheTransactionManager txManager = ClientCacheService.getInstance()
        .getCacheTransactionManager();
    boolean retryTransaction = false;
    int retryTime = 0;
    do {
      try {
        retryTime++;
        if (retryTime > 3) {
          LOGGER.debug("patch transaction commit 3 times, stop trying");
          break;
        }
        txManager.begin();
        for (PatchBody patchBody : patchRequest.getAllPatchList()) {
          processPatchBody(patchBody);
        }
        txManager.commit();
        code = Code.CREATED;
        LOGGER.debug("patch transaction commit successful, stop trying, retryTime={}", retryTime);
        retryTransaction = false;
      } catch (CommitConflictException ex) {
        retryTransaction = true;
        LOGGER.error("commit patch request fail, err=" + ex.toString());
      } catch (Exception ex) {
        LOGGER.error(ex.toString());
      } finally {
        if (txManager.exists()) {
          LOGGER.error("transaction still exist, will rollback");
          txManager.rollback();
        }
      }
    } while (retryTransaction);
    return code;
  }

  private static void processPatchBody(PatchBody patchBody) {
    Region nrfProfileRegion = ClientCacheService.getInstance().getRegion(NRF_REGION_NAME);
    Region nfInfoRegion = ClientCacheService.getInstance().getRegion(NFINFO_REGION_NAME);
    int patchOperation = patchBody.getOperation();
    String patchType = patchBody.getPatchType();
    List<PatchItem> patchItems = patchBody.getPatchItemList();
    String nrfInstanceId = patchBody.getPatchNrfInstId();

    Map<Object, Object> patchMap = new HashMap<>();
    List<Object> patchKeys = new ArrayList<>();
    try {
      for (PatchItem patchItem : patchItems) {
        if (!patchItem.getPatchValue().isEmpty()) {
          Object typeObject = new JSONTokener(patchItem.getPatchValue()).nextValue();
          if (typeObject instanceof JSONObject) {
            PdxInstance pdxInstance = JSONFormatter.fromJSON(patchItem.getPatchValue());
            patchMap.put(nrfInstanceId, pdxInstance);
          } else if (typeObject instanceof JSONArray) {
            JSONArray nfInfoArray = (JSONArray) typeObject;
            for (int i = 0; i < nfInfoArray.length(); i++) {
              JSONObject infoObject = nfInfoArray.getJSONObject(i);
              String nfInstanceId = infoObject.getJSONObject("body").getString("nfInstanceId");
              PdxInstance pdxInstance = JSONFormatter.fromJSON(infoObject.toString());
              patchMap.put(nfInstanceId, pdxInstance);
            }
          } else {
            LOGGER.error("not a valid patch value={}", patchItem.getPatchValue());
          }
        }
        patchKeys.add(patchItem.getPatchKey());
      }
      switch (patchType) {
        case Code.NFTYPE_NRF:
          switch (patchOperation) {
            case Code.PATCH_ADD:
            case Code.PATCH_REPLACE:
              if (patchMap.size() > 0) {
                nrfProfileRegion.putAll(patchMap);
              }
              break;
            case Code.PATCH_REMOVE:
              nrfProfileRegion.removeAll(patchKeys);
              break;
            default:
              LOGGER.error("unknown operation code");
          }
          break;
        case Code.PATCHINSTID:
          switch (patchOperation) {
            case Code.PATCH_ADD:
            case Code.PATCH_REPLACE:
              if (patchMap.size() > 0) {
                nfInfoRegion.putAll(patchMap);
              }
              break;
            case Code.PATCH_REMOVE:
              nfInfoRegion.removeAll(patchKeys);
              break;
            default:
              LOGGER.error("unknown operation code");
          }
          break;
        case Code.NFTYPE_NRFINFO:
          SelectResults<String> allInfos = (SelectResults<String>) nfInfoRegion.getRegionService()
              .getQueryService().newQuery(GET_NFINFOID_SQL + "'" + nrfInstanceId + "'").execute();
          nfInfoRegion.removeAll(allInfos);
          if (patchMap.size() > 0) {
            nfInfoRegion.putAll(patchMap);
          }
          break;
        case Code.NFTYPE_AMF:
        case Code.NFTYPE_AUSF:
        case Code.NFTYPE_BSF:
        case Code.NFTYPE_CHF:
        case Code.NFTYPE_PCF:
        case Code.NFTYPE_SMF:
        case Code.NFTYPE_UDM:
        case Code.NFTYPE_UDR:
        case Code.NFTYPE_UPF:
          SelectResults<String> nfInfos = (SelectResults<String>) nfInfoRegion.getRegionService()
              .getQueryService().newQuery(
                  GET_NFINFOID_SQL + "'" + nrfInstanceId + "' AND value.body.nfType='" + patchType
                      + "'").execute();
          nfInfoRegion.removeAll(nfInfos);
          if (patchMap.size() > 0) {
            nfInfoRegion.putAll(patchMap);
          }
          break;
        default:
          LOGGER.error("unknown patch type");
      }
    } catch (Exception ex) {
      LOGGER.error(ex.toString());
    }
  }

  private static List<String> getPatchNfTypes(boolean isAll) {
    List<String> patchNfTypes = new ArrayList<>();
    patchNfTypes.add(Code.NFTYPE_NRFINFO);
    patchNfTypes.add(Code.NFTYPE_AMF);
    patchNfTypes.add(Code.NFTYPE_AUSF);
    patchNfTypes.add(Code.NFTYPE_BSF);
    patchNfTypes.add(Code.NFTYPE_CHF);
    patchNfTypes.add(Code.NFTYPE_PCF);
    patchNfTypes.add(Code.NFTYPE_SMF);
    patchNfTypes.add(Code.NFTYPE_UDM);
    patchNfTypes.add(Code.NFTYPE_UDR);
    patchNfTypes.add(Code.NFTYPE_UPF);
    if (isAll) {
      patchNfTypes.add(Code.NFTYPE_NRF);
      patchNfTypes.add(Code.PATCHINSTID);
    }
    return patchNfTypes;
  }
}
