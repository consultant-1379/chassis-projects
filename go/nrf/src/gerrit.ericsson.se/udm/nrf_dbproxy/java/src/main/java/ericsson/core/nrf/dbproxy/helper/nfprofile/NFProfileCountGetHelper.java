package ericsson.core.nrf.dbproxy.helper.nfprofile;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import ericsson.core.nrf.dbproxy.common.*;
import ericsson.core.nrf.dbproxy.config.AttributeConfig;
import ericsson.core.nrf.dbproxy.helper.Helper;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileCountGetRequestProto.NFProfileCountGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileCountGetResponseProto.NFProfileCountGetResponse;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileFilterProto.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchExpressionProto.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchParameterProto.*;

public class NFProfileCountGetHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(NFProfileCountGetHelper.class);

    private static NFProfileCountGetHelper instance;

    private NFProfileCountGetHelper() { }

    public static synchronized NFProfileCountGetHelper getInstance()
    {
        if(null == instance) {
        	instance = new NFProfileCountGetHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {
        NFProfileCountGetRequest request = message.getRequest().getGetRequest().getNfProfileCountGetRequest();
	if(request.hasFilter())
	{
	    return validateFilter(request.getFilter());
	} 
	else
	{
            logger.error("Empty NFProfileCountGetRequest is received");
            return Code.NFMESSAGE_PROTOCOL_ERROR;
        }
    }

    private int validateFilter(NFProfileFilter filter)
    {
	boolean empty_filter = true;

	if(filter.hasExpiredTimeRange()) 
	{
	    empty_filter = false;
	    int code = validateRange(filter.getExpiredTimeRange());
	    if(code != Code.VALID) return code;
	}

	if(filter.hasLastUpdateTimeRange())
	{
	    empty_filter = false;
	    int code = validateRange(filter.getLastUpdateTimeRange());
	    if(code != Code.VALID) return code;
	}

	{
	    int code = validateProvisioned(filter.getProvisioned());
	    if(code != Code.VALID) return code;
	}

	if(filter.hasSearchExpression())
	{
	    empty_filter = false;
	    int code = validateSearchExpression(filter.getSearchExpression());
	    if(code != Code.VALID) return code;
	}
	
	return empty_filter ? Code.EMPTY_NFPROFILE_FILTER : Code.VALID;
    }

    private int validateRange(Range range)
    {
	int code = Code.VALID;
	if(range.getStart() > range.getEnd()) code = Code.INVALID_RANGE;
	return code;
    }

    private int validateProvisioned(int provisioned)
    {
	int code = Code.VALID;
	if(provisioned < Code.REGISTERED_PROVISIONED || provisioned > Code.PROVISIONED_ONLY) code = Code.INVALID_PROVISIONED;
	return code;
    }

    private int validateSearchExpression(SearchExpression search_expression)
    {
	if(search_expression.hasAndExpression()) return validateAndExpression(search_expression.getAndExpression());

	if(search_expression.hasOrExpression()) return validateORExpression(search_expression.getOrExpression());

	return Code.EMPTY_SEARCH_EXPRESSION;
    }

    private int validateAndExpression(AndExpression and_expression)
    {
	if(and_expression.getMetaExpressionCount() == 0) return Code.EMPTY_AND_EXPRESSION;

	int code = Code.VALID;
	for(MetaExpression meta_expression : and_expression.getMetaExpressionList())
	{
	    if(meta_expression.hasSearchParameter()) {
			code = validateSearchParameter(meta_expression.getSearchParameter());
		}
	    else if(meta_expression.hasAndExpression()) {
			code = validateAndExpression(meta_expression.getAndExpression());
		}
	    else if(meta_expression.hasOrExpression()) {
			code = validateORExpression(meta_expression.getOrExpression());
		}
	    else {
			code = Code.EMPTY_META_EXPRESSION;
		}

	    if(code != Code.VALID) {break;}
	}
	
	return code;
    }

	private int validateORExpression(ORExpression or_expression) {
		if (or_expression.getMetaExpressionCount() == 0) return Code.EMPTY_OR_EXPRESSION;

		int code = Code.VALID;
		for (MetaExpression meta_expression : or_expression.getMetaExpressionList()) {
			if (meta_expression.hasSearchParameter()) {
				code = validateSearchParameter(meta_expression.getSearchParameter());
			} else if (meta_expression.hasAndExpression()) {
				code = validateAndExpression(meta_expression.getAndExpression());
			} else if (meta_expression.hasOrExpression()) {
				code = validateORExpression(meta_expression.getOrExpression());
			} else {
				code = Code.EMPTY_META_EXPRESSION;
			}

			if (code != Code.VALID) {
				break;
			}
		}

		return code;
	}

    private int validateSearchParameter(SearchParameter search_parameter)
    {
	if(search_parameter.hasAttribute() == false) return Code.SEARCH_ATTRIBUTE_MISSED;

	SearchAttribute attribute = search_parameter.getAttribute();
	if(attribute.getName().isEmpty()) return Code.EMPTY_ATTRIBUTE_NAME;
	if(attribute.getOperation() < Code.OPERATOR_LT || attribute.getOperation() > Code.OPERATOR_REGEX) return Code.INVALID_ATTRIBUTE_OPERATOR;
	if(AttributeConfig.getInstance().get(attribute.getName()) == null) return Code.ATTRIBUTE_NOT_KNOWN;

	if(search_parameter.hasValue() == false) return Code.SEARCH_VALUE_MISSED;
	SearchValue value = search_parameter.getValue();
	if(value.hasNum() == false && value.hasStr() == false) return Code.EMPTY_SEARCH_VALUE;

	return Code.VALID;
    }

    public NFMessage createResponse(int code)
    {

        NFProfileCountGetResponse nf_profile_count_get_response = NFProfileCountGetResponse.newBuilder().setCode(code).build();
        GetResponse get_response = GetResponse.newBuilder().setNfProfileCountGetResponse(nf_profile_count_get_response).build();
        return createNFMessage(get_response);
    }

    public NFMessage createResponse(ExecutionResult execution_result)
    {
        if(execution_result.getCode() != Code.SUCCESS) 
	{
            return createResponse(execution_result.getCode());
        } 
	 
	
        SearchResult search_result = (SearchResult)execution_result;
	    
	int count = getNFProfileCount(search_result);
        if (count == -1)
        {
            return createResponse(Code.INTERNAL_ERROR);
        } 

        NFProfileCountGetResponse nf_profile_count_get_response = NFProfileCountGetResponse.newBuilder().setCode(search_result.getCode()).setCount(count).build();
        GetResponse get_response = GetResponse.newBuilder().setNfProfileCountGetResponse(nf_profile_count_get_response).build();
        return createNFMessage(get_response);        
    }

    private int getNFProfileCount(SearchResult search_result)
    {
        int count = -1;
        
        if (search_result.getItems().size() < 1)
        {
            return count;
        }

	try
        {
            count = Integer.parseInt(String.valueOf(search_result.getItems().get(0)));
        }
        catch(Exception e)
        {
            logger.error("Fail to format to int, " + e.toString());
        }
     
        return count;
    }        
}
