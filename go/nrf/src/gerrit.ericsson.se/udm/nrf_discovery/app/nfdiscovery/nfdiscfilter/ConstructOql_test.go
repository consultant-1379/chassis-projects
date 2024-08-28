package nfdiscfilter

import (
	"testing"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

func TestAndExpression(t *testing.T) {
	searchExpression := &SearchExpression{}
	searchExpression.where = "value.body.nfType"
	searchExpression.from = ""
	searchExpression.value = "AUSF"
	searchExpression.valueType = constvalue.ValueString
	searchExpression.operation = constvalue.EQ
	metaExpression := &MetaExpression{
		commonExpression: searchExpression,
		expressionType: constvalue.TypeSearchExpression,
	}

	searchExpression2 := &SearchExpression{
		where: "snssai.sd",
		from : "value.helper.sNssais snssai",
		valueType: constvalue.ValueString,
		value: "000000",
		operation: constvalue.EQ,
	}
	metaExpression2 := &MetaExpression{
		commonExpression: searchExpression2,
		expressionType: constvalue.TypeSearchExpression,
	}

	searchExpression3 := &SearchExpression{
		where: "snssai.sd",
		from : "value.helper.sNssais snssai",
		valueType: constvalue.ValueString,
		value: "RESERVED_EMPTY_SD",
		operation: constvalue.EQ,
	}
	metaExpression3 := &MetaExpression{
		commonExpression: searchExpression3,
		expressionType: constvalue.TypeSearchExpression,
	}

	var metaExpressionList []*MetaExpression
	metaExpressionList = append(metaExpressionList, metaExpression2, metaExpression3)

	orExpression := &ORExpression{
		orMetaExpression: metaExpressionList,
	}
	metaOrExpression := &MetaExpression{
		commonExpression: orExpression,
		expressionType: constvalue.TypeORExpression,
	}

	searchExpression4 := &SearchExpression{
		where: "snssai.sst",
		from : "value.helper.sNssais snssai",
		valueType: constvalue.ValueNum,
		value: 1,
		operation: constvalue.EQ,
	}
	metaExpression4 := &MetaExpression{
		commonExpression: searchExpression4,
		expressionType: constvalue.TypeSearchExpression,
	}

	var metaExpressionList2 []*MetaExpression
	metaExpressionList2 = append(metaExpressionList2, metaExpression4, metaOrExpression)

	andExpression := &AndExpression{
		andMetaExpression: metaExpressionList2,
	}
	metaAndExpression := &MetaExpression{
		commonExpression: andExpression,
		expressionType: constvalue.TypeAndExpression,
	}

	var metaExpressionList3 []*MetaExpression
	metaExpressionList3 = append(metaExpressionList3, metaAndExpression, metaExpression)

	andExpression2 := &AndExpression{
		andMetaExpression: metaExpressionList3,
	}
	var searchOQL string
	constructAndExpression(andExpression2, true, &searchOQL)
	if searchOQL != "SELECT DISTINCT value FROM (SELECT DISTINCT value FROM /ericsson-nrf-nfprofiles.entrySet, value.helper.sNssais snssai WHERE ((snssai.sst = 1 AND (snssai.sd = '000000' OR snssai.sd = 'RESERVED_EMPTY_SD')))) value WHERE (value.body.nfType = 'AUSF')" {
		t.Fatal("OQL should be matched, but fail")
	}
}

func TestBuildOQL(t *testing.T) {
	searchExpression := &SearchExpression{}
	searchExpression.where = "value.body.nfType"
	searchExpression.from = ""
	searchExpression.value = "AUSF"
	searchExpression.valueType = constvalue.ValueString
	searchExpression.operation = constvalue.EQ
	metaExpression := &MetaExpression{
		commonExpression: searchExpression,
		expressionType: constvalue.TypeSearchExpression,
	}


	var metaExpressionList2 []*MetaExpression
	metaExpressionList2 = append(metaExpressionList2, metaExpression)

	andExpression := &AndExpression{
		andMetaExpression: metaExpressionList2,
	}
	metaAndExpression := &MetaExpression{
		commonExpression: andExpression,
		expressionType: constvalue.TypeAndExpression,
	}

	var searchOQL string
	buildOql(metaAndExpression, &searchOQL)
	if searchOQL != "SELECT DISTINCT value FROM /ericsson-nrf-nfprofiles.entrySet WHERE (value.body.nfType = 'AUSF')" {
		t.Fatal("OQL should be matched, but not")
	}
}