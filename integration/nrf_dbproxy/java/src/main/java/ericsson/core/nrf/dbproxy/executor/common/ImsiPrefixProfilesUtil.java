package ericsson.core.nrf.dbproxy.executor.common;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.clientcache.schema.ImsiprefixProfiles;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.common.SearchResult;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.imsiprefixprofile.ImsiprefixProfileProto.ImsiprefixProfile;
import java.util.ArrayList;
import java.util.List;
import java.util.HashMap;
import java.util.Map;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.geode.cache.Region;

public class ImsiPrefixProfilesUtil {

  private static final Logger LOGGER = LogManager.getLogger(ImsiPrefixProfilesUtil.class);

  private ImsiPrefixProfilesUtil() {
  }

  public static int addImsiprefixProfile(ImsiprefixProfile imsiprefixProfile) {
    ExecutionResult getResult = ClientCacheService.getInstance()
        .getByID(Code.IMSIPREFIXPROFILE_INDICE, imsiprefixProfile.getImsiPrefix());
    if (getResult.getCode() == Code.DATA_NOT_EXIST) {
      ImsiprefixProfiles imsiprefixProfiles = createImsiprefixProfiles(
          imsiprefixProfile.getImsiPrefix(), imsiprefixProfile.getValueInfo());
      int codePut = ClientCacheService.getInstance()
          .put(Code.IMSIPREFIXPROFILE_INDICE, imsiprefixProfile.getImsiPrefix(),
              imsiprefixProfiles);
      if (codePut != Code.CREATED) {
        LOGGER.error("addImsiprefixProfile: imsiprefix {} created failure, code:{}",
            imsiprefixProfile.getImsiPrefix(), codePut);
        return codePut;
      }
    } else if (getResult.getCode() == Code.SUCCESS) {
      SearchResult searchResult = (SearchResult) getResult;
      List<ImsiprefixProfiles> imsiProfilesList = getImsiProfilesList(searchResult.getItems());
      ImsiprefixProfiles imsiProfiles = imsiProfilesList.get(0);
      imsiProfiles.addValueInfo(imsiprefixProfile.getValueInfo());
      LOGGER.debug("Begin add imsiProfiles: {}", imsiProfiles.toString());
      int codePut = ClientCacheService.getInstance()
          .put(Code.IMSIPREFIXPROFILE_INDICE, imsiprefixProfile.getImsiPrefix(), imsiProfiles);
      if (codePut != Code.CREATED) {
        LOGGER.error("addImsiprefixProfile: expend imsiprefix {} failure, code:{}",
            imsiprefixProfile.getImsiPrefix(), codePut);
        return codePut;
      }

    } else {
      LOGGER.error("addImsiprefixProfile: imsi prefix get failure.");
      return getResult.getCode();
    }
    return Code.CREATED;
  }

  public static int delImsiprefixProfile(ImsiprefixProfile imsiprefixProfile) {
    ExecutionResult getResult = ClientCacheService.getInstance()
        .getByID(Code.IMSIPREFIXPROFILE_INDICE, imsiprefixProfile.getImsiPrefix());
    if (getResult.getCode() == Code.DATA_NOT_EXIST) {
      LOGGER.error("delImsiprefixProfile imsiprefix not exist: {} ",
          imsiprefixProfile.getImsiPrefix());
    } else if (getResult.getCode() == Code.SUCCESS) {
      SearchResult searchResult = (SearchResult) getResult;
      List<ImsiprefixProfiles> imsiProfilesList = getImsiProfilesList(searchResult.getItems());
      ImsiprefixProfiles imsiProfiles = imsiProfilesList.get(0);
      imsiProfiles.rmValueInfo(imsiprefixProfile.getValueInfo());
      LOGGER.debug("Begin change imsiProfiles: {}", imsiProfiles.toString());
      if (imsiProfiles.getValueInfo().size() != 0) {
        int codePut = ClientCacheService.getInstance()
            .put(Code.IMSIPREFIXPROFILE_INDICE, imsiprefixProfile.getImsiPrefix(), imsiProfiles);
        if (codePut != Code.CREATED) {
          LOGGER.error("delImsiprefixProfile: reduce imsiprefix {} failure, code:{}",
              imsiprefixProfile.getImsiPrefix(), codePut);
          return codePut;
        }
      } else {
        int codeDel = ClientCacheService.getInstance()
            .delete(Code.IMSIPREFIXPROFILE_INDICE, imsiprefixProfile.getImsiPrefix(), false);
        if (codeDel != Code.SUCCESS) {
          LOGGER.error("delImsiprefixProfile: delete imsiprefix {} failure, code:{}",
              imsiprefixProfile.getImsiPrefix(), codeDel);
          return codeDel;
        }
      }

    } else {
      LOGGER.error("delImsiprefixProfile: imsi prefix get failure.");
      return getResult.getCode();
    }
    return Code.SUCCESS;
  }

  public static int addImsiprefixProfiles(List<ImsiprefixProfile> profilePutList) {
    Map<Object, Object> addMap = new HashMap<>();
    Region imsiRegion = ClientCacheService.getInstance().getRegion(Code.IMSIPREFIXPROFILE_INDICE);

    for (ImsiprefixProfile imsiprefixProfile : profilePutList) {
      ExecutionResult getResult = ClientCacheService.getInstance()
          .getByID(Code.IMSIPREFIXPROFILE_INDICE, imsiprefixProfile.getImsiPrefix());
      if (getResult.getCode() == Code.DATA_NOT_EXIST) {
        ImsiprefixProfiles imsiprefixProfiles = createImsiprefixProfiles(
            imsiprefixProfile.getImsiPrefix(), imsiprefixProfile.getValueInfo());
        addMap.put(imsiprefixProfile.getImsiPrefix(), imsiprefixProfiles);
      } else if (getResult.getCode() == Code.SUCCESS) {
        SearchResult searchResult = (SearchResult) getResult;
        List<ImsiprefixProfiles> imsiProfilesList = getImsiProfilesList(searchResult.getItems());
        ImsiprefixProfiles imsiprefixProfiles = imsiProfilesList.get(0);
        imsiprefixProfiles.addValueInfo(imsiprefixProfile.getValueInfo());
        addMap.put(imsiprefixProfile.getImsiPrefix(), imsiprefixProfiles);
      } else {
        LOGGER.error("addImsiprefixProfiles: imsi prefix get failure.");
        return getResult.getCode();
      }
    }

    imsiRegion.putAll(addMap);
    return Code.SUCCESS;
  }

  public static int delImsiprefixProfiles(List<ImsiprefixProfile> profileDelList) {
    List<Object> delList = new ArrayList<>();
    Map<Object, Object> updateMap = new HashMap<>();
    Region imsiRegion = ClientCacheService.getInstance().getRegion(Code.IMSIPREFIXPROFILE_INDICE);

    for (ImsiprefixProfile imsiprefixProfile : profileDelList) {
      ExecutionResult getResult = ClientCacheService.getInstance()
          .getByID(Code.IMSIPREFIXPROFILE_INDICE, imsiprefixProfile.getImsiPrefix());
      if (getResult.getCode() == Code.DATA_NOT_EXIST) {
        LOGGER.error("delImsiprefixProfiles imsiprefix not exist: {} ", imsiprefixProfile.getImsiPrefix());
        continue;
      } else if (getResult.getCode() == Code.SUCCESS) {
        SearchResult searchResult = (SearchResult) getResult;
        List<ImsiprefixProfiles> imsiProfilesList = getImsiProfilesList(searchResult.getItems());
        ImsiprefixProfiles imsiprefixProfiles = imsiProfilesList.get(0);
        imsiprefixProfiles.rmValueInfo(imsiprefixProfile.getValueInfo());
        if (imsiprefixProfiles.getValueInfo().size() != 0) {
          updateMap.put(imsiprefixProfile.getImsiPrefix(), imsiprefixProfiles);
        } else {
          delList.add(imsiprefixProfile.getImsiPrefix());
        }
      } else {
        LOGGER.error("delImsiprefixProfiles: imsi prefix get failure.");
        return getResult.getCode();
      }
    }

    imsiRegion.removeAll(delList);
    imsiRegion.putAll(updateMap);
    return Code.SUCCESS;
  }

  private static ImsiprefixProfiles createImsiprefixProfiles(Long imsiPrefix, String valueInfo) {
    ImsiprefixProfiles imsiprefixProfiles = new ImsiprefixProfiles();

    imsiprefixProfiles.setImsiprefix(imsiPrefix);
    imsiprefixProfiles.addValueInfo(valueInfo);

    LOGGER.debug("Imsiprefix Profile : {}", imsiprefixProfiles.toString());

    return imsiprefixProfiles;
  }

  private static List<ImsiprefixProfiles> getImsiProfilesList(List<Object> items) {
    List<ImsiprefixProfiles> imsiprofilesList = new ArrayList<>();
    for (Object obj : items) {
      ImsiprefixProfiles imsiProfiles = (ImsiprefixProfiles) obj;
      imsiprofilesList.add(imsiProfiles);
      LOGGER.debug("getValueInfoList: value_info {}", imsiProfiles.toString());
    }
    return imsiprofilesList;
  }
}
