package nfdiscfilter

import (
	"fmt"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"strings"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
)

const (
	//FROMPREFIX is source of region name
	FROMPREFIX = "/ericsson-nrf-nfprofiles.entrySet";
	//SELECT is select sentence
	SELECT = "SELECT DISTINCT value FROM ";
	//HELPERSELECT is to select helper from
	HELPERSELECT = "SELECT DISTINCT value.helper FROM ";
	//WHERE is where sentence
	WHERE = " WHERE ";
)

//buildOql is to get OQL string from metaExpression
func buildOql(metaExpression *MetaExpression, searchOql *string) {
	if metaExpression.expressionType == constvalue.TypeAndExpression {
		constructAndExpression(metaExpression.commonExpression.(*AndExpression), true, searchOql)
	} else {
		log.Error("metaExpression is not and expression")
	}
}

//buildInstIDOql is to get nfInstanceId OQL string from metaExpression
func buildInstIDOql(metaExpression *MetaExpression, searchOql *string) {
	if metaExpression.expressionType == constvalue.TypeAndExpression {
		constructAndExpression(metaExpression.commonExpression.(*AndExpression), true, searchOql)
		if strings.Index(*searchOql, SELECT) != -1 {
			*searchOql = strings.Replace(*searchOql, SELECT, "SELECT DISTINCT value.nfInstanceId,value.profileUpdateTime FROM ", 1)
		} else {
			*searchOql = strings.Replace(*searchOql, HELPERSELECT, "SELECT DISTINCT value.helper.nfInstanceId,value.helper.profileUpdateTime FROM ", 1)
		}
	} else {
		log.Error("metaExpression is not and expression")
	}
}

//constructOql is to generate each layer oql
func constructOql(fromExpressionList []string, whereExpression string, searchOQL *string) string {
	innerFlag := false
	if *searchOQL == "" {
		innerFlag = true
		if internalconf.DiscCacheEnable {
			*searchOQL = HELPERSELECT + FROMPREFIX
		} else {
			*searchOQL = SELECT + FROMPREFIX
		}
	} else {
		*searchOQL = SELECT + "(" + *searchOQL + ") value"
	}
	for _, value := range (fromExpressionList) {
		if !innerFlag && internalconf.DiscCacheEnable {
			*searchOQL += ", " + strings.Replace(value, "value.helper", "value", -1)
		} else {
			*searchOQL += ", " + value
		}
	}
	if !innerFlag && internalconf.DiscCacheEnable {
		*searchOQL += WHERE + strings.Replace(whereExpression, "value.helper", "value", -1)
	} else {
		*searchOQL += WHERE + whereExpression
	}
	return *searchOQL
}

//constructAndExpression is to generate and expression oql
func constructAndExpression(andExpressions *AndExpression, isFirstLayer bool, searchOQL *string) ([]string, string) {
	var fromExpressionList []string
	var whereExpressionList []string
	for _, value := range (andExpressions.andMetaExpression) {
		var froms []string
		var wheres string
		if value.expressionType == constvalue.TypeSearchExpression {
			search := value.commonExpression.(*SearchExpression)
			froms, wheres = constructSearchExpression(search)
		} else if value.expressionType == constvalue.TypeAndExpression {
			and := value.commonExpression.(*AndExpression)
			froms, wheres = constructAndExpression(and, false, searchOQL)
		} else if value.expressionType == constvalue.TypeORExpression {
			or := value.commonExpression.(*ORExpression)
			froms, wheres = constructOrExpression(or, searchOQL)
		} else {
			log.Error("error type")
		}
		for _, fromvalue := range froms {
			fromArray := strings.Split(fromvalue, ",")
			for _, value := range fromArray {
				exist := false
				for _, value2 := range fromExpressionList {
					if value == value2 {
						exist = true
						break
					}
				}
				if !exist {
					fromExpressionList = append(fromExpressionList, value)
				}
			}
		}

		whereExpressionList = append(whereExpressionList, wheres)
		if isFirstLayer {
			whereExpression := "("
			for index, value := range whereExpressionList {
				if index == len(whereExpressionList) - 1 {
					whereExpression += value
				} else {
					whereExpression += value + " AND "
				}

			}
			whereExpression += ")"
			*searchOQL = constructOql(fromExpressionList, whereExpression, searchOQL)
			fromExpressionList = []string{}
			whereExpressionList = []string{}
		}
	}
	whereExpression := "("
	for index, value := range whereExpressionList {
		if index == len(whereExpressionList) - 1 {
			whereExpression += value
		} else {
			whereExpression += value + " AND "
		}

	}
	whereExpression += ")"

	return fromExpressionList, whereExpression
}

//constructOrExpression is to generate or expression oql
func constructOrExpression(orExpressions *ORExpression, searchOQL *string) ([]string, string) {
	var fromExpressionList []string
	var whereExpressionList []string
	for _, value := range (orExpressions.orMetaExpression) {
		var froms []string
		var wheres string
		if value.expressionType == constvalue.TypeSearchExpression {
			search := value.commonExpression.(*SearchExpression)
			froms, wheres = constructSearchExpression(search)
		} else if value.expressionType == constvalue.TypeAndExpression {
			and := value.commonExpression.(*AndExpression)
			froms, wheres = constructAndExpression(and, false, searchOQL)
		} else if value.expressionType == constvalue.TypeORExpression {
			or := value.commonExpression.(*ORExpression)
			froms, wheres = constructOrExpression(or, searchOQL)
		} else {
			log.Error("error type")
		}
		for _, fromvalue := range froms {
			fromArray := strings.Split(fromvalue, ",")
			for _, value := range fromArray {
				exist := false
				for _, value2 := range fromExpressionList {
					if value == value2 {
						exist = true
						break
					}
				}
				if !exist {
					fromExpressionList = append(fromExpressionList, value)
				}
			}
		}
		whereExpressionList = append(whereExpressionList, wheres)

	}
	whereExpression := "("
	for index, value := range whereExpressionList {
		if index == len(whereExpressionList) - 1 {
			whereExpression += value
		} else {
			whereExpression += value + " OR "
		}

	}
	whereExpression += ")"
	return fromExpressionList, whereExpression
}

//constructSearchExpression is to generate oql match expression
func constructSearchExpression(searchExpression *SearchExpression) ([]string, string) {
	var whereExpression string
	if searchExpression.valueType == constvalue.ValueNum {
		var operation string
		op := searchExpression.operation;
		switch(op) {
		case constvalue.LT:
			operation = " < "
		case constvalue.LE:
			operation = " <= "
		case constvalue.EQ:
			operation = " = "
		case constvalue.GE:
			operation = " >= "
		case constvalue.GT:
			operation = " > "
		default:
			log.Errorf("%s%v%s%v\n", "Invalid operation = ", operation, ", ignore this attribute ", searchExpression.value);
		}
		whereExpression = fmt.Sprintf("%s%s%v", searchExpression.where, operation, searchExpression.value)
	} else if searchExpression.valueType == constvalue.ValueString {
		var operation string
		op := searchExpression.operation;
		switch(op) {
		case constvalue.LT:
			operation = " < "
			whereExpression = fmt.Sprintf("%v%v%v%v%v", searchExpression.where, operation, "'", searchExpression.value, "'")
		case constvalue.LE:
			operation = " <= "
			whereExpression = fmt.Sprintf("%v%v%v%v%v", searchExpression.where, operation, "'", searchExpression.value, "'")
		case constvalue.EQ:
			operation = " = "
			whereExpression = fmt.Sprintf("%v%v%v%v%v", searchExpression.where, operation, "'", searchExpression.value, "'")
		case constvalue.GE:
			operation = " >= "
			whereExpression = fmt.Sprintf("%v%v%v%v%v", searchExpression.where, operation, "'", searchExpression.value, "'")
		case constvalue.GT:
			operation = " > "
			whereExpression = fmt.Sprintf("%v%v%v%v%v", searchExpression.where, operation, "'", searchExpression.value, "'")
		case constvalue.REGEX:
			whereExpression = fmt.Sprintf("%v%v%v%v%v", "'", searchExpression.value, "'.matches(", searchExpression.where, ".toString()) = true")
		default:
			log.Errorf("%v%v%v%v\n", "Invalid operation = ", operation, ", ignore this attribute ", searchExpression.value);
		}

	}
	if searchExpression.from != "" {
		return []string{searchExpression.from}, whereExpression
	}
	return []string{}, whereExpression

}