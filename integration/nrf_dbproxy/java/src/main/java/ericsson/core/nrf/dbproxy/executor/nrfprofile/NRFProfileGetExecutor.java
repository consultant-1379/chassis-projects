package ericsson.core.nrf.dbproxy.executor.nrfprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileFilterProto.NRFProfileFilter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileGetRequestProto.NRFProfileGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfprofile.NRFProfileIndexProto.NRFKeyStruct;
import ericsson.core.nrf.dbproxy.helper.nrfprofile.NRFProfileGetHelper;
import java.util.List;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NRFProfileGetExecutor extends Executor {

  private static final Logger LOGGER = LogManager.getLogger(NRFProfileGetExecutor.class);

  private static final String OQL_1 = " AND ";
  private static NRFProfileGetExecutor instance;

  static {
    instance = null;
  }

  private NRFProfileGetExecutor() {
    super(NRFProfileGetHelper.getInstance());
  }

  public static synchronized NRFProfileGetExecutor getInstance() {
    if (null == instance) {
      instance = new NRFProfileGetExecutor();
    }

    return instance;
  }

  protected ExecutionResult execute(NFMessage request) {
    NRFProfileGetRequest getRequest = request.getRequest().getGetRequest()
        .getNrfProfileGetRequest();
    switch (getRequest.getDataCase()) {
      case NRF_INSTANCE_ID:
        return ClientCacheService.getInstance()
            .getByID(Code.NRFPROFILE_INDICE, getRequest.getNrfInstanceId());
      case FILTER:
        String queryString = getQueryString(getRequest.getFilter());
        return ClientCacheService.getInstance().getByFilter(Code.NRFPROFILE_INDICE, queryString);
      case FRAGMENT_SESSION_ID:
        return ClientCacheService.getInstance()
            .getByFragSessionId(Code.NRFPROFILE_INDICE, getRequest.getFragmentSessionId());
      default:
        return new ExecutionResult(Code.NFMESSAGE_PROTOCOL_ERROR);
    }
  }

  private String getQueryString(NRFProfileFilter filter) {
    String regionName = Code.NRFPROFILE_INDICE;
    StringBuilder sbPrevious = new StringBuilder("SELECT DISTINCT p FROM /" + regionName + " p");
    StringBuilder sb = new StringBuilder("");

    String operation = "OR";
    if (filter.getAndOperation()) {
      operation = "AND";
    }

    boolean keyExist = false;

    keyExist = addKeyStructQueryString(filter.getIndex().getAmfKey1List(), "amf_key1", "AMFK1", sb,
        sbPrevious, operation, keyExist);
    keyExist = addKeyStructQueryString(filter.getIndex().getAmfKey2List(), "amf_key2", "AMFK2", sb,
        sbPrevious, operation, keyExist);
    keyExist = addKeyStructQueryString(filter.getIndex().getAmfKey3List(), "amf_key3", "AMFK3", sb,
        sbPrevious, operation, keyExist);
    keyExist = addKeyStructQueryString(filter.getIndex().getAmfKey4List(), "amf_key4", "AMFK4", sb,
        sbPrevious, operation, keyExist);

    keyExist = addKeyStructQueryString(filter.getIndex().getSmfKey1List(), "smf_key1", "SMFK1", sb,
        sbPrevious, operation, keyExist);
    keyExist = addKeyStructQueryString(filter.getIndex().getSmfKey2List(), "smf_key2", "SMFK2", sb,
        sbPrevious, operation, keyExist);
    keyExist = addKeyStructQueryString(filter.getIndex().getSmfKey3List(), "smf_key3", "SMFK3", sb,
        sbPrevious, operation, keyExist);

    keyExist = addKeyStructQueryString(filter.getIndex().getUdmKey1List(), "udm_key1", "UDMK1", sb,
        sbPrevious, operation, keyExist);
    keyExist = addKeyStructQueryString(filter.getIndex().getUdmKey2List(), "udm_key2", "UDMK2", sb,
        sbPrevious, operation, keyExist);

    keyExist = addKeyStructQueryString(filter.getIndex().getAusfKey1List(), "ausf_key1", "AUSFK1",
        sb, sbPrevious, operation, keyExist);
    keyExist = addKeyStructQueryString(filter.getIndex().getAusfKey2List(), "ausf_key2", "AUSFK2",
        sb, sbPrevious, operation, keyExist);

    keyExist = addKeyStructQueryString(filter.getIndex().getPcfKey1List(), "pcf_key1", "PCFK1", sb,
        sbPrevious, operation, keyExist);
    keyExist = addKeyStructQueryString(filter.getIndex().getPcfKey2List(), "pcf_key2", "PCFK2", sb,
        sbPrevious, operation, keyExist);

    long key1 = filter.getIndex().getKey1();
    long key2 = filter.getIndex().getKey2();
    if (key1 < key2) {
      if (keyExist) {
        sb.append(" " + operation + " ");
      }
      sb.append("(p.key1 >= " + key1 + "L" + OQL_1 + "p.key1 < " + key2 + "L)");
    }

    sbPrevious.append(" WHERE ");
    String queryString = sbPrevious.toString() + sb.toString();
    queryString = getCommonQueryString(filter, queryString, keyExist);

    LOGGER.debug("OQL = {}", queryString);

    return queryString;
  }

  private boolean addKeyStructQueryString(List<NRFKeyStruct> nksList, String keyName,
      String indexName, StringBuilder sb, StringBuilder sbPrevious, String operation,
      boolean keyExist) {

    StringBuilder ksb = new StringBuilder("");
    for (NRFKeyStruct nks : nksList) {

      String subKey1 = nks.getSubKey1();
      String subKey2 = nks.getSubKey2();
      String subKey3 = nks.getSubKey3();
      String subKey4 = nks.getSubKey4();
      String subKey5 = nks.getSubKey5();

      StringBuilder subSb = new StringBuilder("");
      boolean subKeyExist = false;

      if (!subKey1.isEmpty()) {
        subSb.append(indexName + ".sub_key1 = '" + subKey1 + "'");
        subKeyExist = true;
      }

      if (!subKey2.isEmpty()) {
        if (subKeyExist) {
          subSb.append(OQL_1);
        }
        subSb.append(indexName + ".sub_key2 = '" + subKey2 + "'");
        subKeyExist = true;
      }

      if (!subKey3.isEmpty()) {
        if (subKeyExist) {
          subSb.append(OQL_1);
        }
        subSb.append(indexName + ".sub_key3 = '" + subKey3 + "'");
        subKeyExist = true;
      }

      if (!subKey4.isEmpty()) {
        if (subKeyExist) {
          subSb.append(OQL_1);
        }
        subSb.append(indexName + ".sub_key4 = '" + subKey4 + "'");
        subKeyExist = true;
      }

      if (!subKey5.isEmpty()) {
        if (subKeyExist) {
          subSb.append(OQL_1);
        }
        subSb.append(indexName + ".sub_key5 = '" + subKey5 + "'");
      }

      if (subSb.length() != 0) {
        if (ksb.length() == 0) {
          ksb.append("(" + subSb.toString() + ")");
        } else {
          ksb.append(" OR (" + subSb.toString() + ")");
        }
      }
    }

    if (ksb.length() != 0) {

      sbPrevious.append(", p." + keyName + ".values " + indexName);

      if (keyExist) {
        sb.append(" " + operation + " ");
      }
      sb.append("(" + ksb.toString() + ")");
      keyExist = true;
    }
    return keyExist;
  }

  private String getCommonQueryString(NRFProfileFilter filter, String innerQuery,
      boolean keyExist) {
    StringBuilder sb = new StringBuilder("");

    long key3 = filter.getIndex().getKey3();
    if (key3 == 1 || key3 == 2) {
      if (keyExist) {
        sb.append(OQL_1);
      }
      sb.append("value.key3 = " + key3);
    }

    if (sb.length() == 0) {
      return innerQuery;
    }

    if (innerQuery.isEmpty()) {
      String regionName = Code.NFPROFILE_INDICE;
      return "SELECT DISTINCT value FROM /" + regionName + ".entrySet WHERE " + sb.toString();
    } else {
      return "SELECT DISTINCT value FROM (" + innerQuery + ") value WHERE " + sb.toString();
    }
  }
}
