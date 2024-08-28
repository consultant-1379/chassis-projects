package ericsson.core.nrf.dbproxy.executor.nrfaddress;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressFilterProto.NRFAddressFilter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nrfaddress.NRFAddressGetRequestProto.NRFAddressGetRequest;
import ericsson.core.nrf.dbproxy.helper.nrfaddress.NRFAddressGetHelper;
import java.util.List;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class NRFAddressGetExecutor extends Executor {

  private static final Logger LOGGER = LogManager.getLogger(NRFAddressGetExecutor.class);

  private static NRFAddressGetExecutor instance;

  static {
    instance = null;
  }

  private NRFAddressGetExecutor() {
    super(NRFAddressGetHelper.getInstance());
  }

  public static synchronized NRFAddressGetExecutor getInstance() {
    if (null == instance) {
      instance = new NRFAddressGetExecutor();
    }
    return instance;
  }

  protected ExecutionResult execute(NFMessage request) {
    NRFAddressGetRequest getRequest = request.getRequest().getGetRequest()
        .getNrfAddressGetRequest();
    switch (getRequest.getDataCase()) {
      case NRF_ADDRESS_ID:
        return ClientCacheService.getInstance()
            .getByID(Code.NRFADDRESS_INDICE, getRequest.getNrfAddressId());
      case FILTER:
        String queryString = getQueryString(getRequest.getFilter());
        return ClientCacheService.getInstance().getByFilter(Code.NRFADDRESS_INDICE, queryString);
      default:
        return new ExecutionResult(Code.NFMESSAGE_PROTOCOL_ERROR);
    }
  }

  private String getQueryString(NRFAddressFilter filter) {

    String regionName = Code.NRFADDRESS_INDICE;
    StringBuilder sb = new StringBuilder("SELECT * FROM /" + regionName + " p WHERE ");

    String operation = "OR";
    if (filter.getAndOperation()) {
      operation = "AND";
    }

    List<String> key1List = filter.getIndex().getNrfAddressKey1List();
    String key2 = filter.getIndex().getNrfAddressKey2();
    String key3 = filter.getIndex().getNrfAddressKey3();
    String key4 = filter.getIndex().getNrfAddressKey4();
    String key5 = filter.getIndex().getNrfAddressKey5();

    boolean keyExist = false;
    if (!key1List.isEmpty()) {
      boolean needOR = false;
      sb.append("(");
      for (String key1 : key1List) {
        if (needOR) {
          sb.append(" OR ");
        }
        sb.append("p.key1 = '" + key1 + "'");
        needOR = true;
        keyExist = true;
      }
      sb.append(")");
    }

    if (!key2.isEmpty()) {
      if (keyExist) {
        sb.append(" " + operation + " ");
      }
      sb.append("p.key2 = '" + key2 + "'");
      keyExist = true;
    }

    if (!key3.isEmpty()) {
      if (keyExist) {
        sb.append(" " + operation + " ");
      }
      sb.append("p.key3 = '" + key3 + "'");
      keyExist = true;
    }

    if (!key4.isEmpty()) {
      if (keyExist) {
        sb.append(" " + operation + " ");
      }
      sb.append("p.key4 = '" + key4 + "'");
      keyExist = true;
    }

    if (!key5.isEmpty()) {
      if (keyExist) {
        sb.append(" " + operation + " ");
      }
      sb.append("p.key5 = '" + key5 + "'");
    }

    String queryString = sb.toString();
    LOGGER.debug("OQL = {}", queryString);

    return queryString;

  }

}
