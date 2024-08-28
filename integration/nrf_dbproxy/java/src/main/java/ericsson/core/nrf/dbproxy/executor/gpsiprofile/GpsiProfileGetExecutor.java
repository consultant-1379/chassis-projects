package ericsson.core.nrf.dbproxy.executor.gpsiprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileFilterProto.GpsiProfileFilter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.gpsiprofile.GpsiProfileGetRequestProto.GpsiProfileGetRequest;
import ericsson.core.nrf.dbproxy.helper.gpsiprofile.GpsiProfileGetHelper;
import java.util.List;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class GpsiProfileGetExecutor extends Executor {

  private static final Logger LOGGER = LogManager.getLogger(GpsiProfileGetExecutor.class);

  private static GpsiProfileGetExecutor instance;

  static {
    instance = null;
  }

  private GpsiProfileGetExecutor() {
    super(GpsiProfileGetHelper.getInstance());
  }

  public static synchronized GpsiProfileGetExecutor getInstance() {
    if (null == instance) {
      instance = new GpsiProfileGetExecutor();
    }
    return instance;
  }

  protected ExecutionResult execute(NFMessage request) {
    GpsiProfileGetRequest getRequest = request.getRequest().getGetRequest()
        .getGpsiProfileGetRequest();
    switch (getRequest.getDataCase()) {
      case GPSI_PROFILE_ID:
        return ClientCacheService.getInstance()
            .getByID(Code.GPSIPROFILE_INDICE, getRequest.getGpsiProfileId());
      case FILTER:
        String queryString = getQueryString(getRequest.getFilter());
        return ClientCacheService.getInstance().getByFilter(Code.GPSIPROFILE_INDICE, queryString);
      case FRAGMENT_SESSION_ID:
        return ClientCacheService.getInstance()
            .getByFragSessionId(Code.GPSIPROFILE_INDICE, getRequest.getFragmentSessionId());
      default:
        return new ExecutionResult(Code.NFMESSAGE_PROTOCOL_ERROR);
    }
  }

  private String getQueryString(GpsiProfileFilter filter) {
    String regionName = Code.GPSIPROFILE_INDICE;
    StringBuilder sb = new StringBuilder("SELECT * FROM /" + regionName + " p WHERE ");

    String operation = "OR";
    if (filter.getAndOperation()) {
      operation = "AND";
    }

    List<String> nfTypeList = filter.getIndex().getNfTypeList();
    boolean nfTypeExist = false;
    boolean groupIdExist = false;
    if (!nfTypeList.isEmpty()) {

      boolean needOR = false;
      sb.append("(");
      for (String nfType : nfTypeList) {
        if (needOR) {
          sb.append(" OR ");
        }
        sb.append("p.nf_type['" + nfType + "'] = 1");
        needOR = true;
        nfTypeExist = true;
      }
      sb.append(")");
    }

    List<String> groupIdList = filter.getIndex().getGroupIndexList();
    if (!groupIdList.isEmpty()) {
      boolean needOR = false;
      if (nfTypeExist) {
        sb.append(" " + operation + " ");
      }
      sb.append("(");
      for (String groupId : groupIdList) {

        if (needOR) {
          sb.append(" OR ");
        }
        sb.append("p.group_id['" + groupId + "'] = 1");
        needOR = true;
        groupIdExist = true;
      }
      sb.append(")");
    }

    //For profile_type judgement, the operation must be "AND"
    int profileType = filter.getIndex().getProfileType();
    if (nfTypeExist || groupIdExist) {
      sb.append(" AND ");
    }
    sb.append("p.profile_type = " + profileType);

    String queryString = sb.toString();
    LOGGER.debug("OQL = {}", queryString);

    return queryString;

  }

}
