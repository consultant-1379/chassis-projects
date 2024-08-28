package ericsson.core.nrf.dbproxy.executor.common;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.clientcache.schema.GpsiprefixProfiles;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprefixprofile.GpsiprefixProfileProto.GpsiprefixProfile;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.HashMap;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.geode.cache.Region;

public class GpsiPrefixProfilesUtil {

  private static final Logger LOGGER = LogManager.getLogger(GpsiPrefixProfilesUtil.class);

  private GpsiPrefixProfilesUtil() {

  }

  public static int addGpsiprefixProfile(GpsiprefixProfile gpsiprefixProfile) {
    ExecutionResult getResult = ClientCacheService.getInstance()
        .getByID(Code.GPSIPREFIXPROFILE_INDICE, gpsiprefixProfile.getGpsiPrefix());
    if (getResult.getCode() == Code.DATA_NOT_EXIST) {
      GpsiprefixProfiles gpsiprefixProfiles = createGpsiprefixProfiles(
          gpsiprefixProfile.getGpsiPrefix(), gpsiprefixProfile.getValueInfo());
      int codePut = ClientCacheService.getInstance()
          .put(Code.GPSIPREFIXPROFILE_INDICE, gpsiprefixProfile.getGpsiPrefix(),
              gpsiprefixProfiles);
      if (codePut != Code.CREATED) {
        LOGGER.error("addGpsiprefixProfile: gpsiprefix {} created failure, code:{}",
            gpsiprefixProfile.getGpsiPrefix(), codePut);
        return codePut;
      }
    } else if (getResult.getCode() == Code.SUCCESS) {
      SearchResult searchResult = (SearchResult) getResult;
      List<GpsiprefixProfiles> gpsiProfilesList = getGpsiProfilesList(searchResult.getItems());
      GpsiprefixProfiles gpsiProfiles = gpsiProfilesList.get(0);
      gpsiProfiles.addValueInfo(gpsiprefixProfile.getValueInfo());
      LOGGER.debug("Begin add gpsiProfiles: {}", gpsiProfiles.toString());
      int codePut = ClientCacheService.getInstance()
          .put(Code.GPSIPREFIXPROFILE_INDICE, gpsiprefixProfile.getGpsiPrefix(), gpsiProfiles);
      if (codePut != Code.CREATED) {
        LOGGER.error("addGpsiprefixProfile: expend gpsiprefix {} failure, code:{}",
            gpsiprefixProfile.getGpsiPrefix(), codePut);
        return codePut;
      }

    } else {
      LOGGER.error("addGpsiprefixProfile: gpsi prefix get failure.");
      return getResult.getCode();
    }
    return Code.CREATED;
  }

  public static int delGpsiprefixProfile(GpsiprefixProfile gpsiprefixProfile) {
    ExecutionResult getResult = ClientCacheService.getInstance()
        .getByID(Code.GPSIPREFIXPROFILE_INDICE, gpsiprefixProfile.getGpsiPrefix());
    if (getResult.getCode() == Code.DATA_NOT_EXIST) {
      LOGGER.error("delGpsiprefixProfile gpsiprefix not exist: {} ",
          gpsiprefixProfile.getGpsiPrefix());
    } else if (getResult.getCode() == Code.SUCCESS) {
      SearchResult searchResult = (SearchResult) getResult;
      List<GpsiprefixProfiles> gpsiProfilesList = getGpsiProfilesList(searchResult.getItems());
      GpsiprefixProfiles gpsiProfiles = gpsiProfilesList.get(0);
      gpsiProfiles.rmValueInfo(gpsiprefixProfile.getValueInfo());
      LOGGER.debug("Begin change gpsiProfiles: {}", gpsiProfiles.toString());
      if (gpsiProfiles.getValueInfo().size() != 0) {
        int codePut = ClientCacheService.getInstance()
            .put(Code.GPSIPREFIXPROFILE_INDICE, gpsiprefixProfile.getGpsiPrefix(), gpsiProfiles);
        if (codePut != Code.CREATED) {
          LOGGER.error("delGpsiprefixProfile: reduce gpsiprefix {} failure, code:{}",
              gpsiprefixProfile.getGpsiPrefix(), codePut);
          return codePut;
        }
      } else {
        int codeDel = ClientCacheService.getInstance()
            .delete(Code.GPSIPREFIXPROFILE_INDICE, gpsiprefixProfile.getGpsiPrefix(), false);
        if (codeDel != Code.SUCCESS) {
          LOGGER.error("delGpsiprefixProfile: delete gpsiprefix {} failure, code:{}",
              gpsiprefixProfile.getGpsiPrefix(), codeDel);
          return codeDel;
        }
      }

    } else {
      LOGGER.error("delGpsiprefixProfile: gpsi prefix get failure.");
      return getResult.getCode();
    }
    return Code.SUCCESS;
  }

  public static int addGpsiprefixProfiles(List<GpsiprefixProfile> profilePutList) {
    Map<Object, Object> addMap = new HashMap<>();
    Region gpsiRegion = ClientCacheService.getInstance().getRegion(Code.GPSIPREFIXPROFILE_INDICE);

    for (GpsiprefixProfile gpsiprefixProfile : profilePutList) {
      ExecutionResult getResult = ClientCacheService.getInstance()
          .getByID(Code.GPSIPREFIXPROFILE_INDICE, gpsiprefixProfile.getGpsiPrefix());
      if (getResult.getCode() == Code.DATA_NOT_EXIST) {
        GpsiprefixProfiles gpsiprefixProfiles = createGpsiprefixProfiles(
            gpsiprefixProfile.getGpsiPrefix(), gpsiprefixProfile.getValueInfo());
        addMap.put(gpsiprefixProfile.getGpsiPrefix(), gpsiprefixProfiles);
      } else if (getResult.getCode() == Code.SUCCESS) {
        SearchResult searchResult = (SearchResult) getResult;
        List<GpsiprefixProfiles> gpsiProfilesList = getGpsiProfilesList(searchResult.getItems());
        GpsiprefixProfiles gpsiprefixProfiles = gpsiProfilesList.get(0);
        gpsiprefixProfiles.addValueInfo(gpsiprefixProfile.getValueInfo());
        addMap.put(gpsiprefixProfile.getGpsiPrefix(), gpsiprefixProfiles);
      } else {
        LOGGER.error("addGpsiprefixProfiles: gpsi prefix get failure.");
        return getResult.getCode();
      }
    }

    gpsiRegion.putAll(addMap);
    return Code.SUCCESS;
  }

  public static int delGpsiprefixProfiles(List<GpsiprefixProfile> profileDelList) {
    List<Object> delList = new ArrayList<>();
    Map<Object, Object> updateMap = new HashMap<>();
    Region gpsiRegion = ClientCacheService.getInstance().getRegion(Code.GPSIPREFIXPROFILE_INDICE);

    for (GpsiprefixProfile gpsiprefixProfile : profileDelList) {
      ExecutionResult getResult = ClientCacheService.getInstance()
          .getByID(Code.GPSIPREFIXPROFILE_INDICE, gpsiprefixProfile.getGpsiPrefix());
      if (getResult.getCode() == Code.DATA_NOT_EXIST) {
        LOGGER.error("delGpsiprefixProfiles gpsiprefix not exist: {} ", gpsiprefixProfile.getGpsiPrefix());
        continue;
      } else if (getResult.getCode() == Code.SUCCESS) {
        SearchResult searchResult = (SearchResult) getResult;
        List<GpsiprefixProfiles> gpsiProfilesList = getGpsiProfilesList(searchResult.getItems());
        GpsiprefixProfiles gpsiprefixProfiles = gpsiProfilesList.get(0);
        gpsiprefixProfiles.rmValueInfo(gpsiprefixProfile.getValueInfo());
        if (gpsiprefixProfiles.getValueInfo().size() != 0) {
          updateMap.put(gpsiprefixProfile.getGpsiPrefix(), gpsiprefixProfiles);
        } else {
          delList.add(gpsiprefixProfile.getGpsiPrefix());
        }
      } else {
        LOGGER.error("delGpsiprefixProfile: gpsi prefix get failure.");
        return getResult.getCode();
      }
    }

    gpsiRegion.removeAll(delList);
    gpsiRegion.putAll(updateMap);
    return Code.SUCCESS;
  }

  private static GpsiprefixProfiles createGpsiprefixProfiles(Long gpsiPrefix, String valueInfo) {
    GpsiprefixProfiles gpsiprefixProfiles = new GpsiprefixProfiles();

    gpsiprefixProfiles.setGpsiprefix(gpsiPrefix);
    gpsiprefixProfiles.addValueInfo(valueInfo);

    LOGGER.debug("Gpsiprefix Profile : {}", gpsiprefixProfiles.toString());

    return gpsiprefixProfiles;
  }

  private static List<GpsiprefixProfiles> getGpsiProfilesList(List<Object> items) {
    List<GpsiprefixProfiles> gpsiprofilesList = new ArrayList<>();
    for (Object obj : items) {
      GpsiprefixProfiles gpsiProfiles = (GpsiprefixProfiles) obj;
      gpsiprofilesList.add(gpsiProfiles);
      LOGGER.debug("getValueInfoList: value_info {}", gpsiProfiles.toString());
    }
    return gpsiprofilesList;
  }
}
