package ericsson.core.nrf.dbproxy.executor.groupprofile;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileFilterProto.GroupProfileFilter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.groupprofile.GroupProfileGetRequestProto.GroupProfileGetRequest;
import ericsson.core.nrf.dbproxy.helper.groupprofile.GroupProfileGetHelper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.util.List;

public class GroupProfileGetExecutor extends Executor
{
    private static final Logger logger = LogManager.getLogger(GroupProfileGetExecutor.class);

    private static GroupProfileGetExecutor instance = null;

    private GroupProfileGetExecutor()
    {
        super(GroupProfileGetHelper.getInstance());
    }

    public static synchronized GroupProfileGetExecutor getInstance()
    {
        if(null == instance) {
            instance = new GroupProfileGetExecutor();
        }
        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        GroupProfileGetRequest get_request = request.getRequest().getGetRequest().getGroupProfileGetRequest();
        switch(get_request.getDataCase()) {
        case GROUP_PROFILE_ID:
            return ClientCacheService.getInstance().getByID(Code.GROUPPROFILE_INDICE, get_request.getGroupProfileId());
        case FILTER:
            String query_string = getQueryString(get_request.getFilter());
            return ClientCacheService.getInstance().getByFilter(Code.GROUPPROFILE_INDICE, query_string);
        case FRAGMENT_SESSION_ID:
            return ClientCacheService.getInstance().getByFragSessionId(Code.GROUPPROFILE_INDICE, get_request.getFragmentSessionId());
        default:
            return new ExecutionResult(Code.NFMESSAGE_PROTOCOL_ERROR);
        }
    }

    private String getQueryString(GroupProfileFilter filter)
    {
        String region_name = Code.GROUPPROFILE_INDICE;
        StringBuilder sb = new StringBuilder("SELECT * FROM /" + region_name + " p WHERE ");

        String operation = "OR";
        if(filter.getAndOperation() == true) operation = "AND";

        List<String> nf_type_list = filter.getIndex().getNfTypeList();
		boolean nf_type_exist = false;
		boolean group_id_exist = false;
        if (nf_type_list.isEmpty() == false) {

            boolean needOR = false;
            sb.append("(");
            for(String nf_type: nf_type_list) {
                if(needOR == true)
                    sb.append(" OR ");
                sb.append("p.nf_type['" + nf_type + "'] = 1");
                needOR = true;
                nf_type_exist = true;
            }
            sb.append(")");
        }
		
         List<String> group_id_list = filter.getIndex().getGroupIndexList();
        if (group_id_list.isEmpty() == false) {
            boolean needOR = false;
            if (nf_type_exist) {
                sb.append(" " + operation + " ");
            }
			sb.append("(");	
            for(String group_id: group_id_list) {

                if(needOR == true)
                    sb.append(" OR ");
                sb.append("p.group_id['" + group_id + "'] = 1");
                needOR = true;
                group_id_exist = true;
            }
            sb.append(")");
        }
		
		//For profile_type judgement, the operation must be "AND"
		int profile_type = filter.getIndex().getProfileType();
		if (nf_type_exist || group_id_exist) {
			sb.append(" AND ");
		}
		sb.append("p.profile_type = " + profile_type);

        String query_string = sb.toString();
        logger.debug("OQL = {}", query_string);

        return query_string;
    }
}
