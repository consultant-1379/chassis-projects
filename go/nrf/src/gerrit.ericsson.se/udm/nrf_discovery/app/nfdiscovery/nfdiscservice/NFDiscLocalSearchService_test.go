package nfdiscservice

import (
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"testing"
)

func TestSupportParaForNRFProfile(t *testing.T) {
	local := &NFDiscLocalSearchService{}
	var DiscPara nfdiscrequest.DiscGetPara
	value := make(map[string][]string)
	DiscPara.InitMember(value)
	var snssaisList []string
	snssaisList = append(snssaisList, "{\"sst\":1,\"sd\":\"1\"}")
	snssaisList = append(snssaisList, "[{\"sst\":2,\"sd\":\"2\"}]")
	DiscPara.SetFlag(constvalue.SearchDataSnssais, true)
	DiscPara.SetValue(constvalue.SearchDataSnssais, snssaisList)
	if !local.supportParaForNRFProfile(DiscPara) {
		t.Fatal("para should support, but fail")
	}
	var dnnList []string
	dnnList = append(dnnList, "111111")
	dnnList = append(dnnList, "111111")
	DiscPara.SetFlag(constvalue.SearchDataDnn, true)
	DiscPara.SetValue(constvalue.SearchDataDnn, dnnList)
	if local.supportParaForNRFProfile(DiscPara) {
		t.Fatal("para should not support, but support")
	}
}
