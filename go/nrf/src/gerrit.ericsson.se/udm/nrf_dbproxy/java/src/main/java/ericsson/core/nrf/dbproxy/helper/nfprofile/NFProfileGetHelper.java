package ericsson.core.nrf.dbproxy.helper.nfprofile;

import java.util.List;
import java.util.ArrayList;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.apache.geode.pdx.JSONFormatter;
import org.apache.geode.pdx.PdxInstance;

import ericsson.core.nrf.dbproxy.common.*;
import ericsson.core.nrf.dbproxy.config.AttributeConfig;
import ericsson.core.nrf.dbproxy.helper.Helper;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileGetRequestProto.NFProfileGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileGetResponseProto.NFProfileGetResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.FragmentInfoProto.FragmentInfo;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileFilterProto.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchExpressionProto.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchParameterProto.*;

public class NFProfileGetHelper extends Helper
{

    private static final Logger logger = LogManager.getLogger(NFProfileGetHelper.class);

    private static NFProfileGetHelper instance;

    private NFProfileGetHelper() { }

    public static synchronized NFProfileGetHelper getInstance()
    {
        if(null == instance) {
        	instance = new NFProfileGetHelper();
        }
        return instance;
    }

    public int validate(NFMessage message)
    {
        NFProfileGetRequest request = message.getRequest().getGetRequest().getNfProfileGetRequest();
	if(request.getTargetNfInstanceId().isEmpty() == false)
	{
	    return validateTargetNfInstanceId(request.getTargetNfInstanceId());
	}	
	else if(request.hasFilter())
	{
	    return validateFilter(request.getFilter());
	}
        else if(request.getFragmentSessionId().isEmpty() == false)
	{
            return Code.VALID;
        } 
	else
	{
            logger.error("Empty NFProfileGetRequest is received");
            return Code.NFMESSAGE_PROTOCOL_ERROR;
        }
    }

    private int validateTargetNfInstanceId(String target_nf_instance_id)
    {
	int code = Code.VALID;
        if(target_nf_instance_id.length() > Code.KEY_MAX_LENGTH) code = Code.NF_INSTANCE_ID_LENGTH_EXCEED_MAX;
	return code;
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
	
	if(filter.hasProvVersion())
	{
	    empty_filter = false;
	    int code = validateProvVersion(filter.getProvVersion());
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
	
    private int validateProvVersion(ProvVersion provVersion)
    {
	int code = Code.VALID;
	if(provVersion.getSupiVersion() < 0 || provVersion.getGpsiVersion() < 0) code = Code.INVALID_PROV_VERSION;
	return code;
    }
	
    private int validateSearchExpression(SearchExpression search_expression)
    {
	if(search_expression.hasAndExpression()) return validateAndExpression(search_expression.getAndExpression());

	if(search_expression.hasOrExpression()) return validateORExpression(search_expression.getOrExpression());

	return Code.EMPTY_SEARCH_EXPRESSION;
    }

	private int validateAndExpression(AndExpression and_expression) {
		if (and_expression.getMetaExpressionCount() == 0) return Code.EMPTY_AND_EXPRESSION;

		int code = Code.VALID;
		for (MetaExpression meta_expression : and_expression.getMetaExpressionList()) {
			if (meta_expression.hasSearchParameter()) {
				code = validateSearchParameter(meta_expression.getSearchParameter());
			}
			else if (meta_expression.hasAndExpression()) {
				code = validateAndExpression(meta_expression.getAndExpression());
			}
			else if (meta_expression.hasOrExpression()) {
				code = validateORExpression(meta_expression.getOrExpression());
			}
			else {
				code = Code.EMPTY_META_EXPRESSION;
			}

			if (code != Code.VALID) break;
		}

		return code;
	}

     private int validateORExpression(ORExpression or_expression)
    {
	if(or_expression.getMetaExpressionCount() == 0) return Code.EMPTY_OR_EXPRESSION;

	int code = Code.VALID;
	for(MetaExpression meta_expression : or_expression.getMetaExpressionList())
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

	    if(code != Code.VALID) break;
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

        NFProfileGetResponse nf_profile_get_response = NFProfileGetResponse.newBuilder().setCode(code).build();
        GetResponse get_response = GetResponse.newBuilder().setNfProfileGetResponse(nf_profile_get_response).build();
        return createNFMessage(get_response);
    }

    @Override
    public NFMessage createResponse(ExecutionResult execution_result) {
        if (execution_result.getCode() != Code.SUCCESS) {
            return createResponse(execution_result.getCode());
        } else {
            SearchResult search_result = (SearchResult) execution_result;
            if (search_result.isFragmented()) {
                FragmentResult fragment_result = (FragmentResult) search_result;
                if (fragment_result.getFragmentSessionID().isEmpty()) {
                    int firstTransmitNum = FragmentUtil.transmitNumPerTime(fragment_result, Code.NFPROFILE_INDICE);
                    if (FragmentSessionManagement.getInstance().put(fragment_result, firstTransmitNum)) {
                        FragmentResult item = new FragmentResult();
                        item.addAll(fragment_result.getItems().subList(0, firstTransmitNum));
                        item.setFragmentSessionID(fragment_result.getFragmentSessionID());
                        item.setTotalNumber(fragment_result.getTotalNumber());
                        item.setTransmittedNumber(fragment_result.getTransmittedNumber());
                        return createResponse(item);
                    } else {
                        return createResponse(Code.INTERNAL_ERROR);
                    }
                } else {
                    String fragment_session_id = fragment_result.getFragmentSessionID();
                    int total_number = fragment_result.getTotalNumber();
                    int transmitted_number = fragment_result.getTransmittedNumber();
                    FragmentInfo fragment_info = FragmentInfo.newBuilder().setFragmentSessionId(fragment_session_id).setTotalNumber(total_number).setTransmittedNumber(transmitted_number).build();

                    List<String> nf_profile_list = getNFProfile(fragment_result);

                    NFProfileGetResponse nf_profile_get_response = NFProfileGetResponse.newBuilder().setCode(fragment_result.getCode()).addAllNfProfile(nf_profile_list).setFragmentInfo(fragment_info).build();
                    GetResponse get_response = GetResponse.newBuilder().setNfProfileGetResponse(nf_profile_get_response).build();
                    return createNFMessage(get_response);
                }
            } else {
                List<String> nf_profile_list = getNFProfile(search_result);
                NFProfileGetResponse nf_profile_get_response = NFProfileGetResponse.newBuilder().setCode(search_result.getCode()).addAllNfProfile(nf_profile_list).build();
                GetResponse get_response = GetResponse.newBuilder().setNfProfileGetResponse(nf_profile_get_response).build();
                return createNFMessage(get_response);
            }
        }
    }

    private List<String> getNFProfile(SearchResult search_result)
    {
	List<String> nf_profile_list = new ArrayList<>();
	for(Object obj : search_result.getItems())
	{
	    try
	    {
		nf_profile_list.add(JSONFormatter.toJSON((PdxInstance)obj));
	    }
	    catch(Exception e)
	    {
		logger.error("Fail to format to JSON, " + e.toString());
	    }
	}
	return nf_profile_list;
    }
}
