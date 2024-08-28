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
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.util.List;

public class NRFProfileGetExecutor extends Executor
{
    private static final Logger logger = LogManager.getLogger(NRFProfileGetExecutor.class);

    private static final String OQL_1 = " AND ";
    private static NRFProfileGetExecutor instance = null;

    private NRFProfileGetExecutor()
    {
        super(NRFProfileGetHelper.getInstance());
    }

    public static synchronized NRFProfileGetExecutor getInstance()
    {
        if (null == instance) {
            instance = new NRFProfileGetExecutor();
        }

        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        NRFProfileGetRequest get_request = request.getRequest().getGetRequest().getNrfProfileGetRequest();
        switch (get_request.getDataCase()) {
        case NRF_INSTANCE_ID:
            return ClientCacheService.getInstance().getByID(Code.NRFPROFILE_INDICE, get_request.getNrfInstanceId());
        case FILTER:
            String query_string = getQueryString(get_request.getFilter());
            return ClientCacheService.getInstance().getByFilter(Code.NRFPROFILE_INDICE, query_string);
        case FRAGMENT_SESSION_ID:
            return ClientCacheService.getInstance().getByFragSessionId(Code.NRFPROFILE_INDICE, get_request.getFragmentSessionId());
        default:
            return new ExecutionResult(Code.NFMESSAGE_PROTOCOL_ERROR);
        }
    }

    private String getQueryString(NRFProfileFilter filter)
    {
        String region_name = Code.NRFPROFILE_INDICE;
        StringBuilder sb_previous = new StringBuilder("SELECT DISTINCT p FROM /" + region_name + " p");
        StringBuilder sb = new StringBuilder("");

        String operation = "OR";
        if (filter.getAndOperation() == true) operation = "AND";

        boolean key_exist = false;

        key_exist = AddKeyStructQueryString(filter.getIndex().getAmfKey1List(), "amf_key1", "AMFK1", sb, sb_previous, operation, key_exist);
        key_exist = AddKeyStructQueryString(filter.getIndex().getAmfKey2List(), "amf_key2", "AMFK2", sb, sb_previous, operation, key_exist);
        key_exist = AddKeyStructQueryString(filter.getIndex().getAmfKey3List(), "amf_key3", "AMFK3", sb, sb_previous, operation, key_exist);
        key_exist = AddKeyStructQueryString(filter.getIndex().getAmfKey4List(), "amf_key4", "AMFK4", sb, sb_previous, operation, key_exist);

        key_exist = AddKeyStructQueryString(filter.getIndex().getSmfKey1List(), "smf_key1", "SMFK1", sb, sb_previous, operation, key_exist);
        key_exist = AddKeyStructQueryString(filter.getIndex().getSmfKey2List(), "smf_key2", "SMFK2", sb, sb_previous, operation, key_exist);
        key_exist = AddKeyStructQueryString(filter.getIndex().getSmfKey3List(), "smf_key3", "SMFK3", sb, sb_previous, operation, key_exist);

        key_exist = AddKeyStructQueryString(filter.getIndex().getUdmKey1List(), "udm_key1", "UDMK1", sb, sb_previous, operation, key_exist);
        key_exist = AddKeyStructQueryString(filter.getIndex().getUdmKey2List(), "udm_key2", "UDMK2", sb, sb_previous, operation, key_exist);

        key_exist = AddKeyStructQueryString(filter.getIndex().getAusfKey1List(), "ausf_key1", "AUSFK1", sb, sb_previous, operation, key_exist);
        key_exist = AddKeyStructQueryString(filter.getIndex().getAusfKey2List(), "ausf_key2", "AUSFK2", sb, sb_previous, operation, key_exist);

        key_exist = AddKeyStructQueryString(filter.getIndex().getPcfKey1List(), "pcf_key1", "PCFK1", sb, sb_previous, operation, key_exist);
        key_exist = AddKeyStructQueryString(filter.getIndex().getPcfKey2List(), "pcf_key2", "PCFK2", sb, sb_previous, operation, key_exist);

        long key1 = filter.getIndex().getKey1();
        long key2 = filter.getIndex().getKey2();
        if (key1 < key2) {
            if (key_exist)
                sb.append(" " + operation + " ");
            sb.append("(p.key1 >= " + key1 + "L" + OQL_1 + "p.key1 < " + key2 + "L)");
        }

        sb_previous.append(" WHERE ");
        String query_string = sb_previous.toString() + sb.toString();
		query_string = getCommonQueryString(filter, query_string, key_exist);

        logger.trace("OQL = {}", query_string);

        return query_string;
    }

    private boolean AddKeyStructQueryString(List<NRFKeyStruct> nks_list, String key_name, String index_name, StringBuilder sb, StringBuilder sb_previous, String operation, boolean key_exist)
    {

        StringBuilder ksb = new StringBuilder("");
        for (NRFKeyStruct nks : nks_list) {

            String sub_key1 = nks.getSubKey1();
            String sub_key2 = nks.getSubKey2();
            String sub_key3 = nks.getSubKey3();
            String sub_key4 = nks.getSubKey4();
            String sub_key5 = nks.getSubKey5();

            StringBuilder sub_sb = new StringBuilder("");
            boolean sub_key_exist = false;

            if (!sub_key1.isEmpty()) {
                sub_sb.append(index_name + ".sub_key1 = '" + sub_key1 + "'");
                sub_key_exist = true;
            }

            if (!sub_key2.isEmpty()) {
                if (sub_key_exist)
                    sub_sb.append(OQL_1);
                sub_sb.append(index_name + ".sub_key2 = '" + sub_key2 + "'");
                sub_key_exist = true;
            }

            if (!sub_key3.isEmpty()) {
                if (sub_key_exist)
                    sub_sb.append(OQL_1);
                sub_sb.append(index_name + ".sub_key3 = '" + sub_key3 + "'");
                sub_key_exist = true;
            }

            if (!sub_key4.isEmpty()) {
                if (sub_key_exist)
                    sub_sb.append(OQL_1);
                sub_sb.append(index_name + ".sub_key4 = '" + sub_key4 + "'");
                sub_key_exist = true;
            }

            if (!sub_key5.isEmpty()) {
                if (sub_key_exist)
                    sub_sb.append(OQL_1);
                sub_sb.append(index_name + ".sub_key5 = '" + sub_key5 + "'");
            }

            if (sub_sb.length() != 0) {
                if (ksb.length() == 0)
                    ksb.append("(" + sub_sb.toString() + ")");
                else
                    ksb.append(" OR (" + sub_sb.toString() + ")");
            }
        }

        if (ksb.length() != 0) {

            sb_previous.append(", p." + key_name + ".values " + index_name);

            if (key_exist)
                sb.append(" " + operation + " ");
            sb.append("(" + ksb.toString() + ")");
            key_exist = true;
        }
        return key_exist;
    }

	private String getCommonQueryString(NRFProfileFilter filter, String inner_query, boolean key_exist)
    {
        StringBuilder sb = new StringBuilder("");

        long key3 = filter.getIndex().getKey3();
        if(key3==1 || key3==2) {
            if(key_exist) sb.append(OQL_1);
            sb.append("value.key3 = " + key3);
        }

        if(sb.length() == 0) return inner_query;
        
        if(inner_query.isEmpty()) {
            String region_name = Code.NFPROFILE_INDICE;
            return "SELECT DISTINCT value FROM /" + region_name + ".entrySet WHERE " + sb.toString();
        } else {
            return "SELECT DISTINCT value FROM (" + inner_query + ") value WHERE " + sb.toString();
        }
    }
}
