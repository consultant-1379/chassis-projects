package nfdiscfilter

import (
	"testing"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

func TestMetaExpressionToString(t *testing.T) {
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
	var result string
	andExpression2.andExpressionToString(&result)
	if result != "AND{AND{{where=snssai.sst,value=1,operation=0}OR{{where=snssai.sd,value=000000,operation=0}{where=snssai.sd,value=RESERVED_EMPTY_SD,operation=0}}}{where=value.body.nfType,value=AUSF,operation=0}}" {
		t.Fatal("AndExpression toString should be matched, but fail")
	}
}