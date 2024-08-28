package nfdiscfilter

import (
	"fmt"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

//MetaExpression is common struct of and,or,search expression
type MetaExpression struct {
	commonExpression interface{}
	expressionType int32
}

//AndExpression is and relationship
type AndExpression struct {
	andMetaExpression []*MetaExpression
}

//ORExpression is or relationship
type ORExpression struct {
	orMetaExpression []*MetaExpression
}

//SearchExpression is search expression
type SearchExpression struct {
	from string
	where string
	value interface{}
	operation int32
	valueType interface{}
}

//metaExpressionToString is to print metaExpression
func (m *MetaExpression) metaExpressionToString(result *string) {
	if m.expressionType == constvalue.TypeSearchExpression {
		m.commonExpression.(*SearchExpression).searchExpressionToString(result)
	} else if m.expressionType == constvalue.TypeAndExpression {
		m.commonExpression.(*AndExpression).andExpressionToString(result)
	} else if m.expressionType == constvalue.TypeORExpression {
		m.commonExpression.(*ORExpression).orExpressionToString(result)
	}
}

//andExpressionToString is to print andExpression
func (a *AndExpression)andExpressionToString(result *string) {
	*result = fmt.Sprintf("%sAND{", *result)
	for _, value := range (a.andMetaExpression) {
		if value.expressionType == constvalue.TypeSearchExpression {
			search := value.commonExpression.(*SearchExpression)
			search.searchExpressionToString(result)
		} else if value.expressionType == constvalue.TypeAndExpression {
			and := value.commonExpression.(*AndExpression)
			and.andExpressionToString(result)
		} else if value.expressionType == constvalue.TypeORExpression {
			or := value.commonExpression.(*ORExpression)
			or.orExpressionToString(result)
		}
	}
	*result = fmt.Sprintf("%s}", *result)
}

//orExpressionToString is to print orExpression
func (o *ORExpression)orExpressionToString(result *string) {
	*result = fmt.Sprintf("%sOR{", *result)
	for _, value := range (o.orMetaExpression) {
		if value.expressionType == constvalue.TypeSearchExpression {
			search := value.commonExpression.(*SearchExpression)
			search.searchExpressionToString(result)
		} else if value.expressionType == constvalue.TypeAndExpression {
			and := value.commonExpression.(*AndExpression)
			and.andExpressionToString(result)
		} else if value.expressionType == constvalue.TypeORExpression {
			or := value.commonExpression.(*ORExpression)
			or.orExpressionToString(result)
		}
	}
	*result = fmt.Sprintf("%s}", *result)
}

//searchExpressionToString is to print searchExpression
func (s *SearchExpression)searchExpressionToString(result *string) {
	*result = fmt.Sprintf("%s{where=%s,value=%v,operation=%v}", *result, s.where, s.value, s.operation)

}