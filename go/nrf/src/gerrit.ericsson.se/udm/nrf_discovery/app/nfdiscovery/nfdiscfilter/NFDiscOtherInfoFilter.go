package nfdiscfilter

import (
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/common/pkg/log"
)
//NFOtherInfoFilter to process other info nftype filter
type NFOtherInfoFilter struct {

}

func (a *NFOtherInfoFilter) filter(nfprofile []byte, queryFrom *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool {
	return true

}

func (a *NFOtherInfoFilter) filterByKVDB(queryForm *nfdiscrequest.DiscGetPara) []*MetaExpression {
	log.Debugf("Enter NFUPFInfoFilter filterByTypeInfo %v", queryForm.GetValue())
	var metaExpressionList []*MetaExpression
	return metaExpressionList

}
