/*package ericsson.core.nrf.dbproxy.dbproxy_client;

import java.util.List;
import java.util.ArrayList;
import ericsson.core.nrf.dbproxy.common.Code;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFMessage;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.NFRequest;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.PutRequest;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.GetRequest;
import ericsson.core.nrf.dbproxy.grpc.NFDataManagementServiceOuterClass.DelRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfilePutRequestProto.NFProfilePutRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfilePutResponseProto.NFProfilePutResponse;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileGetRequestProto.NFProfileGetRequest;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileFilterProto.NFProfileFilter;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.nfprofile.NFProfileFilterProto.Range;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchExpressionProto.*;
import ericsson.core.nrf.dbproxy.grpc.nfmessage.common.SearchParameterProto.*;

public class RequestBuilder
{
    public static NFMessage buildNFProfilePutRequest()
    {
	String[] json = JSONUtil.readFile("/pdu/code/udm1_registration.json");
	String nf_instance_id = json[0];
	String nf_profile = json[1];
	NFProfilePutRequest nf_profile_put_request = NFProfilePutRequest.newBuilder().setNfInstanceId(nf_instance_id).setNfProfile(nf_profile).build();
	PutRequest put_request = PutRequest.newBuilder().setNfProfilePutRequest(nf_profile_put_request).build();
	return buildPutRequest(put_request);
    }

    public static NFMessage buildPutRequest(PutRequest put_request)
    {
	NFRequest nf_request = NFRequest.newBuilder().setPutRequest(put_request).build();
	NFMessage request = NFMessage.newBuilder().setRequest(nf_request).build();
	return request;
    }

    public static MetaExpression buildStringSearchParameter(String name, String value, int operation)
    {
	SearchAttribute search_attribute = SearchAttribute.newBuilder().setName(name).setOperation(operation).build();
	StringValue str = StringValue.newBuilder().setValue(value).build();
	SearchValue search_value = SearchValue.newBuilder().setStr(str).build();

	SearchParameter search_parameter = SearchParameter.newBuilder().setAttribute(search_attribute).setValue(search_value).build();

	MetaExpression meta_expression = MetaExpression.newBuilder().setSearchParameter(search_parameter).build();
	return meta_expression;
    }

    public static MetaExpression buildIntegerSearchParameter(String name, int value, int operation)
    {
	SearchAttribute search_attribute = SearchAttribute.newBuilder().setName(name).setOperation(operation).build();
	IntegerValue num = IntegerValue.newBuilder().setValue(value).build();
	SearchValue search_value = SearchValue.newBuilder().setNum(num).build();

	SearchParameter search_parameter = SearchParameter.newBuilder().setAttribute(search_attribute).setValue(search_value).build();

	MetaExpression meta_expression = MetaExpression.newBuilder().setSearchParameter(search_parameter).build();
	return meta_expression;
    }

    public static MetaExpression buildSnssai(String path)
    {
	SearchAttribute sst_attribute = SearchAttribute.newBuilder().setName(path+".sst").setOperation(Code.OPERATOR_EQ).build();
	SearchAttribute sd_attribute = SearchAttribute.newBuilder().setName(path+".sd").setOperation(Code.OPERATOR_EQ).build();

	List<MetaExpression> meta_expression_list =  new ArrayList<>();
	{
	    IntegerValue sst_num = IntegerValue.newBuilder().setValue(10).build();
	    SearchValue sst_value = SearchValue.newBuilder().setNum(sst_num).build();
	    SearchParameter sst_parameter = SearchParameter.newBuilder().setAttribute(sst_attribute).setValue(sst_value).build();
	    MetaExpression sst_expression = MetaExpression.newBuilder().setSearchParameter(sst_parameter).build();

	    StringValue sd_str = StringValue.newBuilder().setValue("aa").build();
	    SearchValue sd_value = SearchValue.newBuilder().setStr(sd_str).build();
	    SearchParameter sd_parameter = SearchParameter.newBuilder().setAttribute(sd_attribute).setValue(sd_value).build();
	    MetaExpression sd_expression = MetaExpression.newBuilder().setSearchParameter(sd_parameter).build();

	    AndExpression and_expression = AndExpression.newBuilder().addMetaExpression(sst_expression).addMetaExpression(sd_expression).build();
	    MetaExpression meta_expression = MetaExpression.newBuilder().setAndExpression(and_expression).build();
	    meta_expression_list.add(meta_expression);
        }

	{
	    IntegerValue sst_num = IntegerValue.newBuilder().setValue(20).build();
	    SearchValue sst_value = SearchValue.newBuilder().setNum(sst_num).build();
	    SearchParameter sst_parameter = SearchParameter.newBuilder().setAttribute(sst_attribute).setValue(sst_value).build();
	    MetaExpression sst_expression = MetaExpression.newBuilder().setSearchParameter(sst_parameter).build();

	    StringValue sd_str = StringValue.newBuilder().setValue("bb").build();
	    SearchValue sd_value = SearchValue.newBuilder().setStr(sd_str).build();
	    SearchParameter sd_parameter = SearchParameter.newBuilder().setAttribute(sd_attribute).setValue(sd_value).build();
	    MetaExpression sd_expression = MetaExpression.newBuilder().setSearchParameter(sd_parameter).build();

	    AndExpression and_expression = AndExpression.newBuilder().addMetaExpression(sst_expression).addMetaExpression(sd_expression).build();
	    MetaExpression meta_expression = MetaExpression.newBuilder().setAndExpression(and_expression).build();
	    meta_expression_list.add(meta_expression);
	}

	{
	    IntegerValue sst_num = IntegerValue.newBuilder().setValue(30).build();
	    SearchValue sst_value = SearchValue.newBuilder().setNum(sst_num).build();
	    SearchParameter sst_parameter = SearchParameter.newBuilder().setAttribute(sst_attribute).setValue(sst_value).build();
	    MetaExpression sst_expression = MetaExpression.newBuilder().setSearchParameter(sst_parameter).build();

	    meta_expression_list.add(sst_expression);

	}

	ORExpression or_expression = ORExpression.newBuilder().addAllMetaExpression(meta_expression_list).build();
	return MetaExpression.newBuilder().setOrExpression(or_expression).build();
    }

    public static MetaExpression buildGpsi(String path)
    {
	String gpsi = "gpsi-10000";
	MetaExpression range_expression = buildRangeExpression(path, gpsi);
	MetaExpression pattern_expression = buildStringSearchParameter(path + ".pattern", gpsi, Code.OPERATOR_REGEX);

	ORExpression or_expression = ORExpression.newBuilder().addMetaExpression(range_expression).addMetaExpression(pattern_expression).build();

	return MetaExpression.newBuilder().setOrExpression(or_expression).build();
    }

    public static MetaExpression buildRangeExpression(String path, String value)
    {
	MetaExpression start_meta_expression = buildStringSearchParameter(path + ".start", value, Code.OPERATOR_LE);
	MetaExpression end_meta_expression = buildStringSearchParameter(path + ".end", value, Code.OPERATOR_GE);

	AndExpression and_expression = AndExpression.newBuilder().addMetaExpression(start_meta_expression).addMetaExpression(end_meta_expression).build();

	return MetaExpression.newBuilder().setAndExpression(and_expression).build();
    }

    public static MetaExpression buildTargetNFType()
    {
        return buildStringSearchParameter("body.nfType", "UDM", Code.OPERATOR_EQ);
    }

    public static MetaExpression buildNFStatus()
    {
	return buildStringSearchParameter("body.nfStatus", "REGISTERED", Code.OPERATOR_EQ);
    }

    public static MetaExpression buildNSIList()
    {
	List<MetaExpression> meta_expression_list = new ArrayList<>();
	for(int i = 1; i < 4; i++)
	{
	    meta_expression_list.add(buildStringSearchParameter("body.nsiList", "nsi-" + Integer.toString(i), Code.OPERATOR_EQ));
	}

	ORExpression or_expression = ORExpression.newBuilder().addAllMetaExpression(meta_expression_list).build();

	return MetaExpression.newBuilder().setOrExpression(or_expression).build();
    }

    public static MetaExpression buildRequesterNFType()
    {
	return buildStringSearchParameter("body.nfServices.allowedNfTypes", "AMF", Code.OPERATOR_EQ);
    }

    public static MetaExpression buildServiceNames()
    {
	List<MetaExpression> meta_expression_list = new ArrayList<>();
	for(int i = 1; i < 4; i++)
	{
	    meta_expression_list.add(buildStringSearchParameter("body.nfServices.serviceName", "service-name-" + Integer.toString(i), Code.OPERATOR_EQ));
	}

	ORExpression or_expression = ORExpression.newBuilder().addAllMetaExpression(meta_expression_list).build();

	return MetaExpression.newBuilder().setOrExpression(or_expression).build();
    }

    public static MetaExpression buildRequesterNFInstanceFQDN()
    {
	return buildStringSearchParameter("body.nfServices.allowedNfDomains", "nrf.rocket@ericsson.com", Code.OPERATOR_REGEX);
    }

    public static MetaExpression buildTargetPLMN()
    {
	MetaExpression mcc_meta_expression = buildStringSearchParameter("body.plmn.mcc", "460", Code.OPERATOR_EQ);
	MetaExpression mnc_meta_expression = buildStringSearchParameter("body.plmn.mnc", "010", Code.OPERATOR_EQ);

	AndExpression and_expression = AndExpression.newBuilder().addMetaExpression(mcc_meta_expression).addMetaExpression(mnc_meta_expression).build();

	return MetaExpression.newBuilder().setAndExpression(and_expression).build();
    }

    public static MetaExpression buildRequesterPLMN()
    {
	MetaExpression mcc_meta_expression = buildStringSearchParameter("body.nfServices.allowedPlmns.mcc", "460", Code.OPERATOR_EQ);
	MetaExpression mnc_meta_expression = buildStringSearchParameter("body.nfServices.allowedPlmns.mnc", "010", Code.OPERATOR_EQ);

	AndExpression and_expression = AndExpression.newBuilder().addMetaExpression(mcc_meta_expression).addMetaExpression(mnc_meta_expression).build();

	return MetaExpression.newBuilder().setAndExpression(and_expression).build();
    }

    public static MetaExpression buildTai()
    {
	MetaExpression mcc_meta_expression = buildStringSearchParameter("body.amfInfo.taiList.plmnId.mcc", "460", Code.OPERATOR_EQ);
	MetaExpression mnc_meta_expression = buildStringSearchParameter("body.amfInfo.taiList.plmnId.mnc", "010", Code.OPERATOR_EQ);
	MetaExpression tac_expression = buildStringSearchParameter("body.amfInfo.taiList.tac", "tac-001", Code.OPERATOR_EQ);

	AndExpression and_expression = AndExpression.newBuilder().addMetaExpression(mcc_meta_expression).addMetaExpression(mnc_meta_expression).addMetaExpression(tac_expression).build();

	return MetaExpression.newBuilder().setAndExpression(and_expression).build();

    }

    public static List<MetaExpression> buildSearchParameters()
    {
	List<MetaExpression> meta_expression_list = new ArrayList<>();

	meta_expression_list.add(buildTargetNFType());
	meta_expression_list.add(buildTai());

	meta_expression_list.add(buildRequesterNFType());
	meta_expression_list.add(buildServiceNames());
	meta_expression_list.add(buildNFStatus());
	meta_expression_list.add(buildRequesterNFInstanceFQDN());
	meta_expression_list.add(buildTargetPLMN());
	meta_expression_list.add(buildRequesterPLMN());

	meta_expression_list.add(buildNSIList());
	meta_expression_list.add(buildSnssai("body.sNssais"));
	meta_expression_list.add(buildSnssai("body.upfInfo.sNssaiUpfInfoList.sNssai"));
	meta_expression_list.add(buildGpsi("body.udmInfo.gpsiRanges"));



	return meta_expression_list;
    }

    public static NFMessage buildNFProfileGetRequest()
    {
	List<MetaExpression> meta_expression_list = buildSearchParameters();
	AndExpression and_expression = AndExpression.newBuilder().addAllMetaExpression(meta_expression_list).build();
	SearchExpression search_expression = SearchExpression.newBuilder().setAndExpression(and_expression).build();

	Range expired_time_range = Range.newBuilder().setStart(100).setEnd(2000).build();
	Range last_update_time_range = Range.newBuilder().setStart(50).setEnd(5000).build();
	NFProfileFilter filter = NFProfileFilter.newBuilder().setSearchExpression(search_expression).setExpiredTimeRange(expired_time_range).setLastUpdateTimeRange(last_update_time_range).build();
	NFProfileGetRequest nf_profile_get_Request = NFProfileGetRequest.newBuilder().setFilter(filter).build();
	GetRequest get_request = GetRequest.newBuilder().setNfProfileGetRequest(nf_profile_get_Request).build();
	return buildGetRequest(get_request);
    }

    public static NFMessage buildGetRequest(GetRequest get_request)
    {
	NFRequest nf_request = NFRequest.newBuilder().setGetRequest(get_request).build();
	NFMessage request = NFMessage.newBuilder().setRequest(nf_request).build();
	return request;
    }
}
*/
