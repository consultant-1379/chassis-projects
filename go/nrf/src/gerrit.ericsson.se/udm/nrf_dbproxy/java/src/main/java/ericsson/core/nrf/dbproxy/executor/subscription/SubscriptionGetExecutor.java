package ericsson.core.nrf.dbproxy.executor.subscription;

import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionIndexProto.SubscriptionGetIndex;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionGetRequestProto.SubscriptionGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.subscription.SubscriptionIndexProto.SubKeyStruct;
import ericsson.core.nrf.dbproxy.helper.subscription.SubscriptionGetHelper;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.util.List;

public class SubscriptionGetExecutor extends Executor
{
    private static final Logger logger = LogManager.getLogger(SubscriptionGetExecutor.class);

    private static final String OQL_1 = "SELECT DISTINCT value FROM /";
    private static SubscriptionGetExecutor instance = null;

    private SubscriptionGetExecutor()
    {
        super(SubscriptionGetHelper.getInstance());
    }

    public static synchronized SubscriptionGetExecutor getInstance()
    {
        if(null == instance) {
            instance = new SubscriptionGetExecutor();
        }
        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        SubscriptionGetRequest get_request = request.getRequest().getGetRequest().getSubscriptionGetRequest();
        switch(get_request.getDataCase()) {
        case SUBSCRIPTION_ID:
            return ClientCacheService.getInstance().getByID(Code.SUBSCRIPTION_INDICE, get_request.getSubscriptionId());
        case FILTER:
            String query_string = getQueryString(get_request.getFilter().getIndex());
            ExecutionResult result = ClientCacheService.getInstance().getByFilter(Code.SUBSCRIPTION_INDICE, query_string);
            if(query_string.contains("validityTime") && result.getCode() == Code.SUCCESS) {
                result.setCode(Code.SUBSCRIPTION_MONITOR_SUCCESS);
            }
            return result;
        default:
            return new ExecutionResult(Code.NFMESSAGE_PROTOCOL_ERROR);
        }
    }

    private String getQueryString(SubscriptionGetIndex index)
    {
        String query_string = "";
        String region_name = Code.SUBSCRIPTION_INDICE;
        String FROM = region_name + ".entrySet";
        String WHERE = "";
		
        String noCond = index.getNoCond();
        if (!noCond.isEmpty()) {
            if (WHERE.isEmpty()) {
                WHERE = "value.noCond = '" + noCond + "'";
            } else {
                WHERE = WHERE + " OR value.noCond = '" + noCond + "'";
            }
        }
		
        String nfStatusNotificationUri = index.getNfStatusNotificationUri();
        if (!nfStatusNotificationUri.isEmpty()) {
            if (WHERE.isEmpty()) {
                WHERE = "value.nfStatusNotificationUri = '" + nfStatusNotificationUri + "'";
            } else {
                WHERE = WHERE + " OR value.nfStatusNotificationUri = '" + nfStatusNotificationUri + "'";
            }
        }
		
        String nfInstanceId = index.getNfInstanceId();
        if (!nfInstanceId.isEmpty()) {
            if (WHERE.isEmpty()) {
                WHERE = "value.nfInstanceId = '" + nfInstanceId + "'";
            } else {
                WHERE = WHERE + " OR value.nfInstanceId = '" + nfInstanceId + "'";
            }
        }
		
        String nfType = index.getNfType();
        if (!nfType.isEmpty()) {
            if (WHERE.isEmpty()) {
                WHERE = "value.nfType = '" + nfType + "'";
            } else {
                WHERE = WHERE + " OR value.nfType = '" + nfType + "'";
            }
        }     
		
        StringBuilder sb = new StringBuilder("");
        int serviceNameCount = 0;
        for(String serviceName : index.getServiceNamesList()) {
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
            if (WHERE.isEmpty()) {
                WHERE = "(" + sb.toString() + ")";
            } else {
                WHERE = WHERE + " OR " + "(" + sb.toString() + ")";
            }
        }
		
        String amfCondAlias = "amf";
        sb = new StringBuilder("");
        for(SubKeyStruct ks : index.getAmfCondsList()) {
            StringBuilder sub_sb = new StringBuilder("");

            String sub_key1 = ks.getSubKey1();
            boolean sub_key_exist = false;
            if(!sub_key1.isEmpty()) {
                sub_sb.append(amfCondAlias + ".sub_key1 = '" + sub_key1 + "'");
                sub_key_exist = true;
            }

            String sub_key2 = ks.getSubKey2();
            if(!sub_key2.isEmpty()) {
                if (sub_key_exist) {
                    sub_sb.append(" AND ");
                }
                sub_sb.append(amfCondAlias + ".sub_key2 = '" + sub_key2 + "'");
            }

            if(sub_sb.length() == 0) continue;

            if(sb.length() == 0)
                sb.append("(" + sub_sb.toString() + ")");
            else
                sb.append(" OR (" + sub_sb.toString() + ")");
        }

        if(sb.length() > 0) {
            FROM = FROM + ", value.amfCond.values " + amfCondAlias;
            if (WHERE.isEmpty()) {
                WHERE = "(" + sb.toString() + ")";
            } else {
                WHERE = WHERE + " OR ("  + sb.toString() + ")";
            }
        }
		
        String guamiListAlias = "guami";
        sb = new StringBuilder("");
        for(SubKeyStruct ks : index.getGuamiListList()) {
            StringBuilder sub_sb = new StringBuilder("");

            String sub_key1 = ks.getSubKey1();
            boolean sub_key_exist = false;
            if(!sub_key1.isEmpty()) {
                sub_sb.append(guamiListAlias + ".sub_key1 = '" + sub_key1 + "'");
                sub_key_exist = true;
            }

            String sub_key2 = ks.getSubKey2();
            if(!sub_key2.isEmpty()) {
                if (sub_key_exist) {
                    sub_sb.append(" AND ");
                } 
                sub_sb.append(guamiListAlias + ".sub_key2 = '" + sub_key2 + "'");
            }

            if(sub_sb.length() == 0) continue;

            if(sb.length() == 0)
                sb.append("(" + sub_sb.toString() + ")");
            else
                sb.append(" OR (" + sub_sb.toString() + ")");
        }

        if(sb.length() > 0) {
            FROM = FROM + ", value.guamiList.values " + guamiListAlias;
            if (WHERE.isEmpty()) {
                WHERE = "(" + sb.toString() + ")";
            } else {
                WHERE = WHERE + " OR (" + sb.toString() + ")";
            }
        }
		
        String snssaiListAlias = "snssai";
        sb = new StringBuilder("");
        for(SubKeyStruct ks : index.getSnssaiListList()) {
            StringBuilder sub_sb = new StringBuilder("");

            String sub_key1 = ks.getSubKey1();
            boolean sub_key_exist = false;
            if(!sub_key1.isEmpty()) {
                sub_sb.append(snssaiListAlias + ".sub_key1 = '" + sub_key1 + "'");
                sub_key_exist = true;
            }

            String sub_key2 = ks.getSubKey2();
            if(!sub_key2.isEmpty()) {
                if (sub_key_exist) {
                    sub_sb.append(" AND ");
                }
                sub_sb.append(snssaiListAlias + ".sub_key2 = '" + sub_key2 + "'");
            }

            if(sub_sb.length() == 0) continue;

            if(sb.length() == 0)
                sb.append("(" + sub_sb.toString() + ")");
            else
                sb.append(" OR (" + sub_sb.toString() + ")");
        }

        if(sb.length() > 0) {
            String FROMTMP = "value.snssaiList.values " + snssaiListAlias;
            String WHERETMP = sb.toString();

            String nsiListAlias = "nsi";
            sb = new StringBuilder("");
            int nsiCount = 0;
            for(String nsi : index.getNsiListList()) {
                if(!nsi.isEmpty()) {
                    if(nsiCount == 0) {
                        sb.append(nsiListAlias + " = '" + nsi + "'");
                    } else {
                        sb.append(" OR " + nsiListAlias + " = '" + nsi + "'");
                    }
                    nsiCount++;
                }
            }

            if (sb.length() > 0) {
                FROMTMP = FROMTMP + ", value.nsiList " + nsiListAlias;
                WHERETMP = "(" + WHERETMP + ") AND (" + sb.toString() + ")"; 
            }

            FROM = FROM + ", " + FROMTMP;
            WHERE = WHERE + " OR (" + WHERETMP + ")";
        }

        String nfGroupCondAlias = "nfGroup";
        sb = new StringBuilder("");
        SubKeyStruct nfGroupCond = index.getNfGroupCond();
        if (nfGroupCond != null) {
            boolean sub_key_exist = false;
            String sub_key1 = nfGroupCond.getSubKey1();
            if (!sub_key1.isEmpty()) {
                sb.append(nfGroupCondAlias + ".sub_key1 = '" + sub_key1 + "'");
                sub_key_exist = true;
            }

            String sub_key2 = nfGroupCond.getSubKey2();
            if (!sub_key2.isEmpty()) {
                if (sub_key_exist) {
                    sb.append(" AND ");
                }
                sb.append(nfGroupCondAlias + ".sub_key2 = '" + sub_key2 + "'");
            }
        }

        if (sb.length() > 0) {
            FROM = FROM + ", value.nfGroupCond.values " + nfGroupCondAlias;
            if (WHERE.isEmpty()) {
                WHERE = "(" + sb.toString() + ")";
            } else {
                WHERE = WHERE + " OR (" + sb.toString() + ")";
            }
        }
		
        long startValidityTime = index.getStartValidityTime();
        long endValidityTime = index.getEndValidityTime();
        if(startValidityTime < endValidityTime) {
            if (WHERE.isEmpty()) {
                WHERE = "value.validityTime >= " + startValidityTime + "L" + " AND " + "value.validityTime < " + endValidityTime + "L";
            } else {
                WHERE = "(" + WHERE + ") AND (" + "value.validityTime >= " + startValidityTime + "L" + " AND " + "value.validityTime < " + endValidityTime + "L)";
            }		
        }

        query_string = OQL_1 + FROM + " WHERE " + WHERE;

        logger.trace("OQL = {}", query_string);

        return query_string;
    }
}
