package ericsson.core.nrf.dbproxy.common;

import ericsson.core.nrf.dbproxy.clientcache.schema.GpsiProfile;
import ericsson.core.nrf.dbproxy.clientcache.schema.GroupProfile;
import ericsson.core.nrf.dbproxy.clientcache.schema.NRFProfile;
import ericsson.core.nrf.dbproxy.config.GeodeConfig;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.QueryResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.TraceInfo;
import java.util.ArrayList;
import java.util.List;
import org.apache.geode.cache.query.SelectResults;
import org.apache.geode.pdx.JSONFormatter;
import org.apache.geode.pdx.PdxInstance;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class FragmentUtil {

  private static final Logger LOGGER = LogManager.getLogger(FragmentUtil.class);

  private FragmentUtil() {
    throw new IllegalStateException("Utility class");
  }

  private static int getMaxTransmitFragmentSize() {
    int maxTransmitSize = GeodeConfig.getMaxTransmitFragmentSize();
    if (maxTransmitSize < 1 * 500 * 1024 || maxTransmitSize > 3 * 1024 * 1024) {
      maxTransmitSize = 3 * 1024 * 1024;
    }
    return maxTransmitSize;
  }

  public static boolean isNeedFragment(String regionName, SelectResults<Object> searchResults) {
    int maxTransmitFragmentSize = getMaxTransmitFragmentSize();
    int totalTransmitSize = 0;
    switch (regionName) {
      case Code.NFPROFILE_INDICE:
        for (Object object : searchResults) {
          PdxInstance pdxInstance = (PdxInstance) object;
          int profileSize = JSONFormatter.toJSON(pdxInstance).getBytes().length;
          totalTransmitSize += profileSize;
        }
        break;
      case Code.NRFPROFILE_INDICE:
        for (Object object : searchResults) {
          NRFProfile nrfProfile = (NRFProfile) object;
          int profileSize = nrfProfile.getRaw_data().size();
          totalTransmitSize += profileSize;
        }
        break;
      case Code.GROUPPROFILE_INDICE:
        for (Object object : searchResults) {
          GroupProfile groupProfile = (GroupProfile) object;
          int profileSize = groupProfile.getData().size();
          totalTransmitSize += profileSize;
        }
        break;
      case Code.GPSIPROFILE_INDICE:
        for (Object object : searchResults) {
          GpsiProfile gpsiProfile = (GpsiProfile) object;
          int profileSize = gpsiProfile.getData().size();
          totalTransmitSize += profileSize;
        }
        break;
      default:
        break;
    }
    LOGGER.debug("total transmit size=" + totalTransmitSize + ",Max transmit message size="
        + maxTransmitFragmentSize);
    if (totalTransmitSize > maxTransmitFragmentSize) {
      return true;
    }
    return false;
  }

  public static int transmitNumPerTime(FragmentResult fragmentResult, String regionName) {
    int maxTransmitFragmentSize = getMaxTransmitFragmentSize();
    int transmitNum = 0;
    int transmitSize = 0;
    int transmittedNumber = fragmentResult.getTransmittedNumber();
    List<Object> restFragment = fragmentResult.getItems()
        .subList(transmittedNumber, fragmentResult.getItems().size());
    switch (regionName) {
      case Code.NFPROFILE_INDICE:
        for (Object object : restFragment) {
          PdxInstance pdxInstance = (PdxInstance) object;
          int profileSize = JSONFormatter.toJSON(pdxInstance).getBytes().length;
          transmitSize += profileSize;
          if (transmitSize > maxTransmitFragmentSize) {
            break;
          }
          transmitNum++;
        }
        break;
      case Code.NRFPROFILE_INDICE:
        for (Object object : restFragment) {
          NRFProfile nrfProfile = (NRFProfile) object;
          int profileSize = nrfProfile.getRaw_data().size();
          transmitSize += profileSize;
          if (transmitSize > maxTransmitFragmentSize) {
            break;
          }
          transmitNum++;
        }
        break;
      case Code.GROUPPROFILE_INDICE:
        for (Object object : restFragment) {
          GroupProfile groupProfile = (GroupProfile) object;
          int profileSize = groupProfile.getData().size();
          transmitSize += profileSize;
          if (transmitSize > maxTransmitFragmentSize) {
            break;
          }
          transmitNum++;
        }
        break;
      case Code.GPSIPROFILE_INDICE:
        for (Object object : restFragment) {
          GpsiProfile gpsiProfile = (GpsiProfile) object;
          int profileSize = gpsiProfile.getData().size();
          transmitSize += profileSize;
          if (transmitSize > maxTransmitFragmentSize) {
            break;
          }
          transmitNum++;
        }
        break;
      default:
        break;
    }
    LOGGER.debug("Transmit profile num =" + transmitNum);
    return transmitNum;
  }

  public static boolean isNeedFragment(List<String> response) {
    int maxTransmitFragmentSize = getMaxTransmitFragmentSize();
    int totalTransmitSize = 0;
    for (int i = 0; i < response.size(); i++) {
      if (totalTransmitSize > maxTransmitFragmentSize) {
        return true;
      }
      totalTransmitSize += response.get(i).length();
    }
    if (totalTransmitSize > maxTransmitFragmentSize) {
      return true;
    }
    return false;
  }

  public static List<QueryResponse> getFragmentResponse(int code, boolean traceEnabled,
      TraceInfo traceInfo, List<String> response) {
    int maxTransmitFragmentSize = getMaxTransmitFragmentSize();
    int transmitStart = 0;
    int transmitNum = 0;
    int transmitMsgSize = 0;
    List<QueryResponse> responseList = new ArrayList<>();
    List<String> fragmentList = new ArrayList<>();
    for (int i = 0; i < response.size(); i++) {
      if (transmitMsgSize >= maxTransmitFragmentSize) {
        fragmentList.addAll(response.subList(transmitStart, transmitStart + transmitNum - 1));
        QueryResponse queryResponse;
        if (traceEnabled) {
          queryResponse = QueryResponse.newBuilder().setCode(code).addAllValue(fragmentList)
              .setTraceInfo(traceInfo).build();
        } else {
          queryResponse = QueryResponse.newBuilder().setCode(code).addAllValue(fragmentList)
              .build();
        }
        responseList.add(queryResponse);
        transmitStart = i - 1;
        transmitNum = 1;
        transmitMsgSize = response.get(i - 1).length();
        fragmentList.clear();
      }
      transmitNum++;
      transmitMsgSize += response.get(i).length();
      if (i == response.size() - 1) {
        fragmentList.addAll(response.subList(transmitStart, response.size()));
        QueryResponse queryResponse;
        if (traceEnabled) {
          queryResponse = QueryResponse.newBuilder().setCode(code).addAllValue(fragmentList)
              .setTraceInfo(traceInfo).build();
        } else {
          queryResponse = QueryResponse.newBuilder().setCode(code).addAllValue(fragmentList)
              .build();
        }
        responseList.add(queryResponse);
      }
    }
    return responseList;
  }
}
