package ericsson.core.nrf.dbproxy.executor.groupprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileFilterProto.GroupProfileFilter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileGetRequestProto.GroupProfileGetRequest;
import ericsson.core.nrf.dbproxy.helper.groupprofile.GroupProfileGetHelper;
import java.util.List;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

public class GroupProfileGetExecutor extends Executor {

  private static final Logger LOGGER = LogManager.getLogger(GroupProfileGetExecutor.class);

  private static GroupProfileGetExecutor instance;

  static {
    instance = null;
  }

  private GroupProfileGetExecutor() {
    super(GroupProfileGetHelper.getInstance());
  }

  public static synchronized GroupProfileGetExecutor getInstance() {
    if (null == instance) {
      instance = new GroupProfileGetExecutor();
    }
    return instance;
  }

  protected ExecutionResult execute(NFMessage request) {
    GroupProfileGetRequest getRequest = request.getRequest().getGetRequest()
        .getGroupProfileGetRequest();
    switch (getRequest.getDataCase()) {
      case GROUP_PROFILE_ID:
        return ClientCacheService.getInstance()
            .getByID(Code.GROUPPROFILE_INDICE, getRequest.getGroupProfileId());
      case FILTER:
        String queryString = getQueryString(getRequest.getFilter());
        return ClientCacheService.getInstance().getByFilter(Code.GROUPPROFILE_INDICE, queryString);
      case FRAGMENT_SESSION_ID:
        return ClientCacheService.getInstance()
            .getByFragSessionId(Code.GROUPPROFILE_INDICE, getRequest.getFragmentSessionId());
      default:
        return new ExecutionResult(Code.NFMESSAGE_PROTOCOL_ERROR);
    }
  }

  private String getQueryString(GroupProfileFilter filter) {
    String regionName = Code.GROUPPROFILE_INDICE;
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

    //For profileType judgement, the operation must be "AND"
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
