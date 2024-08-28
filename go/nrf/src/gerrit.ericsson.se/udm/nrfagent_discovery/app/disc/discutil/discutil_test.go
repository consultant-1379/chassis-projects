package discutil

import (
	"testing"

	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/cache"
)

var (
	cacheManager *cache.CacheManager
)

func TestPreComplieRegexp(t *testing.T) {
	//PreComplieRegexp()
	t.Log("execute test case PreComplieRegexp")
}

func TestServiceNameScopeVerify(t *testing.T) {
	targetNf := structs.TargetNf{
		RequesterNfType:    "AUSF",
		TargetNfType:       "UDM",
		TargetServiceNames: []string{"nudm-auth-01"},
	}

	serviceNames := make([]string, 0)
	serviceNames = append(serviceNames, "test")
	err := ServiceNameScopeVerify(serviceNames, &targetNf)
	if err == nil {
		t.Errorf("TestServiceNameScopeVerify not configure check error")
	}
	serviceNames = make([]string, 0)
	serviceNames = append(serviceNames, "nudm-auth-01")
	err = ServiceNameScopeVerify(serviceNames, &targetNf)
	if err != nil {
		t.Errorf("TestServiceNameScopeVerify configure check error")
	}
}

func TestRejectVerify(t *testing.T) {
	cache.SetCacheConfig("../../../build/config/cache-index.json")
	cacheManager = cache.Instance()
	cacheManager.InitCache("AUSF", "UDM")

	plmn460 := structs.PlmnID{
		Mcc: "460",
		Mnc: "000",
	}
	plmn470 := structs.PlmnID{
		Mcc: "470",
		Mnc: "000",
	}
	plmn480 := structs.PlmnID{
		Mcc: "480",
		Mnc: "000",
	}
	requestPlmns := make([]string, 0)
	targetPlmns := make([]string, 0)
	targetPlmns = append(targetPlmns, "470:000")
	oriqueryForm := make(map[string][]string)
	oriqueryForm[consts.SearchDataTargetPlmnList] = targetPlmns
	oriqueryForm[consts.SearchDataRequesterPlmnList] = requestPlmns
	requestNfTypes := make([]string, 0)
	requestNfTypes = append(requestNfTypes, "AUSF")
	oriqueryForm[consts.SearchDataRequesterNfType] = requestNfTypes

	localPlmns := make([]structs.PlmnID, 0)
	localPlmns = append(localPlmns, plmn460)
	localPlmns = append(localPlmns, plmn470)
	cache.Instance().SetRequesterPlmns("AUSF", localPlmns)

	discPlmns := make([]structs.PlmnID, 0)
	discPlmns = append(discPlmns, plmn460)
	err := RejectVerify(oriqueryForm, discPlmns)
	if err != nil {
		t.Errorf("TestRejectVerify check pass RejectVerify failure")
	}
	discPlmns2 := make([]structs.PlmnID, 0)
	discPlmns2 = append(discPlmns2, plmn480)
	err = RejectVerify(oriqueryForm, discPlmns2)
	if err == nil {
		t.Errorf("TestRejectVerify check RejectVerify failure")
	}
}
