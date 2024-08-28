package ericsson.core.nrf.dbproxy.executor.subscription;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionGetRequestProto.SubscriptionGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionIndexProto.SubKeyStruct;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionIndexProto.SubscriptionGetIndex;
import ericsson.core.nrf.dbproxy.helper.subscription.SubscriptionGetHelper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class SubscriptionGetExecutor extends Executor {

  private static final Logger LOGGER = LogManager.getLogger(SubscriptionGetExecutor.class);

  private static final String OQL_1 = "SELECT DISTINCT value FROM /";
  private static SubscriptionGetExecutor instance;

  static {
    instance = null;
  }

  private SubscriptionGetExecutor() {
    super(SubscriptionGetHelper.getInstance());
  }

  public static synchronized SubscriptionGetExecutor getInstance() {
    if (null == instance) {
      instance = new SubscriptionGetExecutor();
    }
    return instance;
  }

  protected ExecutionResult execute(NFMessage request) {
    SubscriptionGetRequest getRequest = request.getRequest().getGetRequest()
        .getSubscriptionGetRequest();
    switch (getRequest.getDataCase()) {
      case SUBSCRIPTION_ID:
        return ClientCacheService.getInstance()
            .getByID(Code.SUBSCRIPTION_INDICE, getRequest.getSubscriptionId());
      case FILTER:
        String queryString = getQueryString(getRequest.getFilter().getIndex());
        ExecutionResult result = ClientCacheService.getInstance()
            .getByFilter(Code.SUBSCRIPTION_INDICE, queryString);
        if (queryString.contains("validityTime") && result.getCode() == Code.SUCCESS) {
          result.setCode(Code.SUBSCRIPTION_MONITOR_SUCCESS);
        }
        return result;
      default:
        return new ExecutionResult(Code.NFMESSAGE_PROTOCOL_ERROR);
    }
  }

  private String getQueryString(SubscriptionGetIndex index) {
    String queryString = "";
    String regionName = Code.SUBSCRIPTION_INDICE;
    String from = regionName + ".entrySet";
    String where = "";

    String noCond = index.getNoCond();
    if (!noCond.isEmpty()) {
      if (where.isEmpty()) {
        where = "value.noCond = '" + noCond + "'";
      } else {
        where = where + " OR value.noCond = '" + noCond + "'";
      }
    }

    String nfStatusNotificationUri = index.getNfStatusNotificationUri();
    if (!nfStatusNotificationUri.isEmpty()) {
      if (where.isEmpty()) {
        where = "value.nfStatusNotificationUri = '" + nfStatusNotificationUri + "'";
      } else {
        where = where + " OR value.nfStatusNotificationUri = '" + nfStatusNotificationUri + "'";
      }
    }

    String nfInstanceId = index.getNfInstanceId();
    if (!nfInstanceId.isEmpty()) {
      if (where.isEmpty()) {
        where = "value.nfInstanceId = '" + nfInstanceId + "'";
      } else {
        where = where + " OR value.nfInstanceId = '" + nfInstanceId + "'";
      }
    }

    String nfType = index.getNfType();
    if (!nfType.isEmpty()) {
      if (where.isEmpty()) {
        where = "value.nfType = '" + nfType + "'";
      } else {
        where = where + " OR value.nfType = '" + nfType + "'";
      }
    }

    StringBuilder sb = new StringBuilder("");
    int serviceNameCount = 0;
    for (String serviceName : index.getServiceNamesList()) {
      if (!serviceName.isEmpty()) {
        if (serviceNameCount == 0) {
          sb.append("value.serviceName = '" + serviceName + "'");
        } else {
          sb.append(" OR value.serviceName = '" + serviceName + "'");
        }
        serviceNameCount++;
      }
    }

    if (sb.length() > 0) {
      if (where.isEmpty()) {
        where = "(" + sb.toString() + ")";
      } else {
        where = where + " OR " + "(" + sb.toString() + ")";
      }
    }

    String amfCondAlias = "amf";
    sb = new StringBuilder("");
    for (SubKeyStruct ks : index.getAmfCondsList()) {
      StringBuilder subSb = new StringBuilder("");

      String subKey1 = ks.getSubKey1();
      boolean subKeyExist = false;
      if (!subKey1.isEmpty()) {
        subSb.append(amfCondAlias + ".sub_key1 = '" + subKey1 + "'");
        subKeyExist = true;
      }

      String subKey2 = ks.getSubKey2();
      if (!subKey2.isEmpty()) {
        if (subKeyExist) {
          subSb.append(" AND ");
        }
        subSb.append(amfCondAlias + ".sub_key2 = '" + subKey2 + "'");
      }

      if (subSb.length() == 0) {
        continue;
      }

      if (sb.length() == 0) {
        sb.append("(" + subSb.toString() + ")");
      } else {
        sb.append(" OR (" + subSb.toString() + ")");
      }
    }

    if (sb.length() > 0) {
      from = from + ", value.amfCond.values " + amfCondAlias;
      if (where.isEmpty()) {
        where = "(" + sb.toString() + ")";
      } else {
        where = where + " OR (" + sb.toString() + ")";
      }
    }

    String guamiListAlias = "guami";
    sb = new StringBuilder("");
    for (SubKeyStruct ks : index.getGuamiListList()) {
      StringBuilder subSb = new StringBuilder("");

      String subKey1 = ks.getSubKey1();
      boolean subKeyExist = false;
      if (!subKey1.isEmpty()) {
        subSb.append(guamiListAlias + ".sub_key1 = '" + subKey1 + "'");
        subKeyExist = true;
      }

      String subKey2 = ks.getSubKey2();
      if (!subKey2.isEmpty()) {
        if (subKeyExist) {
          subSb.append(" AND ");
        }
        subSb.append(guamiListAlias + ".sub_key2 = '" + subKey2 + "'");
      }

      if (subSb.length() == 0) {
        continue;
      }

      if (sb.length() == 0) {
        sb.append("(" + subSb.toString() + ")");
      } else {
        sb.append(" OR (" + subSb.toString() + ")");
      }
    }

    if (sb.length() > 0) {
      from = from + ", value.guamiList.values " + guamiListAlias;
      if (where.isEmpty()) {
        where = "(" + sb.toString() + ")";
      } else {
        where = where + " OR (" + sb.toString() + ")";
      }
    }

    String snssaiListAlias = "snssai";
    sb = new StringBuilder("");
    for (SubKeyStruct ks : index.getSnssaiListList()) {
      StringBuilder subSb = new StringBuilder("");

      String subKey1 = ks.getSubKey1();
      boolean subKeyExist = false;
      if (!subKey1.isEmpty()) {
        subSb.append(snssaiListAlias + ".sub_key1 = '" + subKey1 + "'");
        subKeyExist = true;
      }

      String subKey2 = ks.getSubKey2();
      if (!subKey2.isEmpty()) {
        if (subKeyExist) {
          subSb.append(" AND ");
        }
        subSb.append(snssaiListAlias + ".sub_key2 = '" + subKey2 + "'");
      }

      if (subSb.length() == 0) {
        continue;
      }

      if (sb.length() == 0) {
        sb.append("(" + subSb.toString() + ")");
      } else {
        sb.append(" OR (" + subSb.toString() + ")");
      }
    }

    if (sb.length() > 0) {
      String fromTmp = "value.snssaiList.values " + snssaiListAlias;
      String whereTmp = sb.toString();

      String nsiListAlias = "nsi";
      sb = new StringBuilder("");
      int nsiCount = 0;
      for (String nsi : index.getNsiListList()) {
        if (!nsi.isEmpty()) {
          if (nsiCount == 0) {
            sb.append(nsiListAlias + " = '" + nsi + "'");
          } else {
            sb.append(" OR " + nsiListAlias + " = '" + nsi + "'");
          }
          nsiCount++;
        }
      }

      if (sb.length() > 0) {
        fromTmp = fromTmp + ", value.nsiList " + nsiListAlias;
        whereTmp = "(" + whereTmp + ") AND (" + sb.toString() + ")";
      }

      from = from + ", " + fromTmp;
      where = where + " OR (" + whereTmp + ")";
    }

    String nfGroupCondAlias = "nfGroup";
    sb = new StringBuilder("");
    SubKeyStruct nfGroupCond = index.getNfGroupCond();
    if (nfGroupCond != null) {
      boolean subKeyExist = false;
      String subKey1 = nfGroupCond.getSubKey1();
      if (!subKey1.isEmpty()) {
        sb.append(nfGroupCondAlias + ".sub_key1 = '" + subKey1 + "'");
        subKeyExist = true;
      }

      String subKey2 = nfGroupCond.getSubKey2();
      if (!subKey2.isEmpty()) {
        if (subKeyExist) {
          sb.append(" AND ");
        }
        sb.append(nfGroupCondAlias + ".sub_key2 = '" + subKey2 + "'");
      }
    }

    if (sb.length() > 0) {
      from = from + ", value.nfGroupCond.values " + nfGroupCondAlias;
      if (where.isEmpty()) {
        where = "(" + sb.toString() + ")";
      } else {
        where = where + " OR (" + sb.toString() + ")";
      }
    }

    long startValidityTime = index.getStartValidityTime();
    long endValidityTime = index.getEndValidityTime();
    if (startValidityTime < endValidityTime) {
      if (where.isEmpty()) {
        where =
            "value.validityTime >= " + startValidityTime + "L" + " AND " + "value.validityTime < "
                + endValidityTime + "L";
      } else {
        where =
            "(" + where + ") AND (" + "value.validityTime >= " + startValidityTime + "L" + " AND "
                + "value.validityTime < " + endValidityTime + "L)";
      }
    }

    queryString = OQL_1 + from + " where " + where;

    LOGGER.debug("OQL = {}", queryString);

    return queryString;
  }
}
