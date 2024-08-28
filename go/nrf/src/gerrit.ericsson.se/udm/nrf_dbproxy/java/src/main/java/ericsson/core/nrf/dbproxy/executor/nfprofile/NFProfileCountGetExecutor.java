package ericsson.core.nrf.dbproxy.executor.nfprofile;

import java.util.List;
import java.util.ArrayList;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.common.ExecutionResult;
import ericsson.core.nrf.dbproxy.config.Attribute;
import ericsson.core.nrf.dbproxy.config.AttributeConfig;
import ericsson.core.nrf.dbproxy.executor.Executor;
import ericsson.core.nrf.dbproxy.helper.nfprofile.NFProfileCountGetHelper;
import ericsson.core.nrf.dbproxy.clientcache.ClientCacheService;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileCountGetRequestProto.NFProfileCountGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileFilterProto.NFProfileFilter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchExpressionProto.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchParameterProto.*;

public class NFProfileCountGetExecutor extends Executor
{
    private static final Logger logger = LogManager.getLogger(NFProfileCountGetExecutor.class);

    private static final String OQL_1 = ") value";
    private static final String OQL_2 = " AND ";
    private static final String OQL_3 = ".compareTo('";
    private static NFProfileCountGetExecutor instance = null;
    private String FROM_PREFIX;
    private String SELECT;
    private String SELECTCOUNT;
    private String WHERE;

    private NFProfileCountGetExecutor()
    {
        super(NFProfileCountGetHelper.getInstance());

	FROM_PREFIX = "/" + Code.NFPROFILE_INDICE + ".entrySet";
	SELECT = "SELECT DISTINCT value FROM ";
        SELECTCOUNT = "SELECT COUNT(*) FROM ";
	WHERE = " WHERE ";
	
    }

    public static synchronized NFProfileCountGetExecutor getInstance()
    {
        if(null == instance) {
        	instance = new NFProfileCountGetExecutor();
        }
        return instance;
    }

    protected ExecutionResult execute(NFMessage request)
    {
        NFProfileCountGetRequest get_request = request.getRequest().getGetRequest().getNfProfileCountGetRequest();
	if(get_request.hasFilter())
	{
            String query_string = getQueryString(get_request.getFilter());
            return ClientCacheService.getInstance().getByCountFilter(Code.NFPROFILE_INDICE, query_string);
	}
	else
	{
            return new ExecutionResult(Code.NFMESSAGE_PROTOCOL_ERROR);
	}
    }

    private String getQueryString(NFProfileFilter filter)
    {
	String[] expressions = {FROM_PREFIX, ""};
	if(filter.hasSearchExpression() && filter.getSearchExpression().hasAndExpression())
	{
	    if(filter.getSearchExpression().hasAndExpression()) {
			expressions = constructAndExpression(filter.getSearchExpression().getAndExpression(), true, false);
		} else if(filter.getSearchExpression().hasOrExpression()) {
			expressions = constructOrExpression(filter.getSearchExpression().getOrExpression());
		}
	}

	String query = buildQuery(expressions);

        query = addCustomInfo(filter, query);
		
	return CountQuery(query);
    }

    private String CountQuery(String query)
    {
        if(query.isEmpty())
            return SELECTCOUNT + FROM_PREFIX;
        else
            return SELECTCOUNT + "(" + query + OQL_1;
    }

    private String buildQuery(String[] expressions)
    {
	String query = "";
	if(expressions[1].isEmpty() == false)
	{
	    if(expressions[0].isEmpty()) {
			query = SELECT + FROM_PREFIX + WHERE + expressions[1];
		}
	    else if(expressions[0].indexOf("SELECT") == -1) {
			query = SELECT + FROM_PREFIX + ", " + expressions[0] + WHERE + expressions[1];
		}
	    else {
			query = SELECT + expressions[0] + WHERE + expressions[1];
		}
	}

	return query;
    }

    private String addCustomInfo(NFProfileFilter filter, String query)
    {
	boolean exist = false;
	StringBuilder where = new StringBuilder(WHERE);
	if(filter.hasExpiredTimeRange())
	{
	    long start = filter.getExpiredTimeRange().getStart();
	    long end = filter.getExpiredTimeRange().getEnd();
	    where.append("value.expiredTime >= " + Long.toString(start) + "L AND value.expiredTime <= " + Long.toString(end) + "L");
	    exist = true;
	}

	if(filter.hasLastUpdateTimeRange())
	{
	    long start = filter.getLastUpdateTimeRange().getStart();
	    long end = filter.getLastUpdateTimeRange().getEnd();
	    if(exist) where.append(OQL_2);
	    where.append("value.lastUpdateTime >= " + Long.toString(start) + "L AND value.lastUpdateTime <= " + Long.toString(end) + "L");
	    exist = true;
	}

	if(filter.getProvisioned() == Code.REGISTERED_ONLY || filter.getProvisioned() == Code.PROVISIONED_ONLY)
	{
	    if(exist) where.append(OQL_2);
	    where.append("value.provisioned = " + Integer.toString(filter.getProvisioned()));
	}

	if(exist == false) return query;

	if(query.isEmpty())
	    return SELECT + FROM_PREFIX + where.toString();
	else
	    return SELECT + "(" + query + OQL_1 + where.toString();
    }

    public String[] constructAndExpression(AndExpression and_expression, boolean in_and_expression, boolean inner_expression_exist)
    {
	String[] expressions = {"", ""};
	List<String> from_expression_list = new ArrayList<String>();
	for(MetaExpression meta_expression : and_expression.getMetaExpressionList())
        {
	    String[] sub_expressions = {"", ""};
	    if(meta_expression.hasSearchParameter()) {
			sub_expressions = constructSearchParameterExpression(meta_expression.getSearchParameter());
		}
	    else if(meta_expression.hasAndExpression())
	    {
		if(expressions[1].isEmpty())
		    sub_expressions = constructAndExpression(meta_expression.getAndExpression(), in_and_expression, false);
		else
		    sub_expressions = constructAndExpression(meta_expression.getAndExpression(), in_and_expression, true);
	    }
	    else if(meta_expression.hasOrExpression()) {
			sub_expressions = constructOrExpression(meta_expression.getOrExpression());
		}
	    else {
			logger.error("Empty MetaExpression in the AndExpression = " + and_expression.toString());
		}

	    if(sub_expressions[1].isEmpty()) continue;

	    if(in_and_expression && !inner_expression_exist)
	    {
		String inner_query = buildQuery(expressions);
		if(inner_query.isEmpty())
		{
		    expressions = sub_expressions;
		}
		else 
		{
		    if(sub_expressions[0].isEmpty())
		        expressions[0] = "(" + inner_query + OQL_1;
		    else
		        expressions[0] = "(" + inner_query + ") value, " + sub_expressions[0];
		    
		    expressions[1] = sub_expressions[1];
		}
  	    }
	    else
	    {
		if(sub_expressions[0].isEmpty() == false)
		{
		    String[] from_expressions = sub_expressions[0].split(",");
		    for(String from : from_expressions)
		    {
			if(from.isEmpty()) continue;
			if(from_expression_list.contains(from)) continue;
			from_expression_list.add(from);
			if(expressions[0].isEmpty()) expressions[0] = from;
			else expressions[0] += ", " + from;
		    }
		}	

		if(expressions[1].isEmpty() == false) 
		    expressions[1] += OQL_2;
		expressions[1] += sub_expressions[1];
	    }
        }

	if(in_and_expression == false && expressions[1].isEmpty() == false)
	    expressions[1] = "(" + expressions[1] + ")";

	return expressions;
    }

    public String[] constructOrExpression(ORExpression or_expression)
    {
	String[] expressions = {"", ""};
	List<String> from_expression_list = new ArrayList<String>();
	for(MetaExpression meta_expression : or_expression.getMetaExpressionList())
        {
	    String[] sub_expressions = {"", ""};
	    if(meta_expression.hasSearchParameter()) {
			sub_expressions = constructSearchParameterExpression(meta_expression.getSearchParameter());
		}
	    else if(meta_expression.hasAndExpression()) {
			sub_expressions = constructAndExpression(meta_expression.getAndExpression(), false, true);
		}
	    else if(meta_expression.hasOrExpression()) {
			sub_expressions = constructOrExpression(meta_expression.getOrExpression());
		}
	    else {
			logger.error("Empty MetaExpression in the ORExpression = " + or_expression.toString());
		}

	    if(sub_expressions[1].isEmpty()) continue;

	    if(sub_expressions[0].isEmpty() == false)
	    {
		String[] from_expressions = sub_expressions[0].split(",");
		for(String from : from_expressions)
		{
		    if(from.isEmpty()) continue;
		    if(from_expression_list.contains(from)) continue;
		    from_expression_list.add(from);
		    if(expressions[0].isEmpty()) expressions[0] = from;
		    else expressions[0] += ", " + from;
		}
	    }	

	    if(expressions[1].isEmpty() == false) {
			expressions[1] += " OR ";
		}
	    expressions[1] += sub_expressions[1];
        }

	if(expressions[1].isEmpty() == false)
	    expressions[1] = "(" + expressions[1] + ")";

	return expressions;
    }


    private String[] constructSearchParameterExpression(SearchParameter search_parameter)
    {
	SearchAttribute search_attribute = search_parameter.getAttribute();
	Attribute attribute = AttributeConfig.getInstance().get(search_attribute.getName());
	String name = attribute.getWhere();
	int operation = search_attribute.getOperation();

	String[] expressions = {"", ""};
	SearchValue search_value = search_parameter.getValue();
	if(search_value.hasNum())
	{
		String op = "";
		switch(operation)
		{
		    case Code.OPERATOR_LT: op = " < "; break;
		    case Code.OPERATOR_LE: op = " <= "; break;
		    case Code.OPERATOR_EQ: op = " = "; break;
		    case Code.OPERATOR_GE: op = " >= "; break;
		    case Code.OPERATOR_GT: op = " > "; break;
		    default:
			logger.warn("Invalid operation = " + Long.toString(operation) + ", ignore this attribute " + search_attribute.getName());
			return expressions;
		}
		expressions[0] = attribute.getFrom();
		expressions[1] = name + op + Long.toString(search_value.getNum().getValue()) + "L";
	}
	else if(search_value.hasStr())
	{
		String str = search_value.getStr().getValue();
		switch(operation)
		{
		    case Code.OPERATOR_LT: expressions[1] = name + OQL_3 + str + "') < 0"; break;
		    case Code.OPERATOR_LE: expressions[1] = name + OQL_3 + str + "') <= 0"; break;
		    case Code.OPERATOR_EQ: expressions[1] = name + " = '" + str + "'"; break;
		    case Code.OPERATOR_GE: expressions[1] = name + OQL_3 + str + "') >= 0"; break;
		    case Code.OPERATOR_GT: expressions[1] = name + OQL_3 + str + "') > 0"; break;
		    case Code.OPERATOR_REGEX: expressions[1] = "'" + str + "'.matches(" + name + ".toString()) = true"; break;
		    default:
			logger.warn("Invalid operation = " + Long.toString(operation) + ", ignore this attribute");
			return expressions;
		}
		expressions[0] = attribute.getFrom();
	}
	else
	{
	    logger.debug("Empty search value for the attribute " + search_attribute.getName());
	}

	return expressions;
    }
}
