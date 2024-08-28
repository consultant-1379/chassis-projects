package nfdiscfilter

import (
	"testing"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

func TestBuildStringSearchParameter(t *testing.T) {
	searchMapping := configmap.SearchMapping{
		Parameter: "target-nf-type",
		Path: "body.nfType",
		From: "",
		Where: "value.body.nfType",
		ExistCheck: false,
	}
	metaExpression := buildStringSearchParameter(searchMapping, "UDM", constvalue.EQ)
	if metaExpression.expressionType != constvalue.TypeSearchExpression {
		t.Fatal("expression type should be search, but not")
	}
	searchExpression := metaExpression.commonExpression.(*SearchExpression)
	if searchExpression.value.(string) != "UDM" {
		t.Fatal("search value should be matched, but not")
	}
	if searchExpression.valueType != constvalue.ValueString {
		t.Fatal("value type should be string, but not")
	}
	if searchExpression.operation != constvalue.EQ {
		t.Fatal("operation should be matched, but not")
	}
	if searchExpression.where != "value.body.nfType" {
		t.Fatal("where should be matched, but not")
	}
}
func TestBuildIntegerSearchParameter(t *testing.T) {
	searchMapping := configmap.SearchMapping{
		Parameter: "snssais/sst",
		Path: "body.sNssais.sst",
		From: "value.helper.sNssais snssai",
		Where: "snssai.sst",
		ExistCheck: false,
	}
	metaExpression := buildIntegerSearchParameter(searchMapping, 123, constvalue.EQ)
	if metaExpression.expressionType != constvalue.TypeSearchExpression {
		t.Fatal("expression type should be search, but not")
	}
	searchExpression := metaExpression.commonExpression.(*SearchExpression)
	if searchExpression.value.(uint64) != 123 {
		t.Fatal("search value should be matched, but not")
	}
	if searchExpression.valueType != constvalue.ValueNum {
		t.Fatal("value type should be string, but not")
	}
	if searchExpression.operation != constvalue.EQ {
		t.Fatal("operation should be matched, but not")
	}
	if searchExpression.where != "snssai.sst" {
		t.Fatal("where should be matched, but not")
	}
	if searchExpression.from != "value.helper.sNssais snssai" {
		t.Fatal("from should be matched, but not")
	}
}

func TestBuildAndExpression(t *testing.T) {
	searchMapping := configmap.SearchMapping{
		Parameter: "target-nf-type",
		Path: "body.nfType",
		From: "",
		Where: "value.body.nfType",
		ExistCheck: false,
	}
	metaExpression := buildStringSearchParameter(searchMapping, "UDM", constvalue.EQ)
	searchMapping2 := configmap.SearchMapping{
		Parameter: "snssais/sst",
		Path: "body.sNssais.sst",
		From: "value.helper.sNssais snssai",
		Where: "snssai.sst",
		ExistCheck: false,
	}
	metaExpression2 := buildIntegerSearchParameter(searchMapping2, 123, constvalue.EQ)

	var metaExpressionList []*MetaExpression
	metaExpressionList = append(metaExpressionList, metaExpression, metaExpression2)
	andMetaExpression := buildAndExpression(metaExpressionList)
	if andMetaExpression.expressionType != constvalue.TypeAndExpression {
		t.Fatal("expression type should be and, but not")
	}
	andExpression := andMetaExpression.commonExpression.(*AndExpression)
	metaExpressionList2 := andExpression.andMetaExpression
	for key,value := range metaExpressionList2 {
		expression := value.commonExpression.(*SearchExpression)
		if key == 0 {
			if expression.value.(string) != "UDM" {
				t.Fatal("expression 0 is not matched")
			}
		}
		if key == 1 {
			if expression.value.(uint64) != 123 {
				t.Fatal("expression 1 is not matched")
			}
		}
	}
}

func TestBuildORExpression(t *testing.T) {
	searchMapping := configmap.SearchMapping{
		Parameter: "target-nf-type",
		Path: "body.nfType",
		From: "",
		Where: "value.body.nfType",
		ExistCheck: false,
	}
	metaExpression := buildStringSearchParameter(searchMapping, "UDM", constvalue.EQ)
	searchMapping2 := configmap.SearchMapping{
		Parameter: "snssais/sst",
		Path: "body.sNssais.sst",
		From: "value.helper.sNssais snssai",
		Where: "snssai.sst",
		ExistCheck: false,
	}
	metaExpression2 := buildIntegerSearchParameter(searchMapping2, 123, constvalue.EQ)

	var metaExpressionList []*MetaExpression
	metaExpressionList = append(metaExpressionList, metaExpression, metaExpression2)
	orMetaExpression := buildORExpression(metaExpressionList)
	if orMetaExpression.expressionType != constvalue.TypeORExpression {
		t.Fatal("expression type should be or, but not")
	}
	orExpression := orMetaExpression.commonExpression.(*ORExpression)
	metaExpressionList2 := orExpression.orMetaExpression
	for key,value := range metaExpressionList2 {
		expression := value.commonExpression.(*SearchExpression)
		if key == 0 {
			if expression.value.(string) != "UDM" {
				t.Fatal("expression 0 is not matched")
			}
		}
		if key == 1 {
			if expression.value.(uint64) != 123 {
				t.Fatal("expression 1 is not matched")
			}
		}
	}
}