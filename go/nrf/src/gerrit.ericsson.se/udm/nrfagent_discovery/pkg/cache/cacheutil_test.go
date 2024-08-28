package cache

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"gerrit.ericsson.se/udm/nrfagent_common/pkg/client"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/election"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/k8sapiclient"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
)

func stubNotifIPEndPoint() bool {
	ipEndPoint := structs.StatusNotifIPEndPoint{
		Ipv4Address: "192.168.110.112",
		Port:        12345,
	}

	ipEndPointData, err := json.Marshal(ipEndPoint)
	if err != nil {
		return false
	}

	return structs.UpdateStatusNotifIPEndPoint(ipEndPointData)
}

func StubGetLeader(ID string) {
	election.GetLeader = func() string {
		return ID
	}
}

func TestProlongSubscriptionFromNrf(t *testing.T) {
	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()

	StubHTTPDoToNrf("PATCHSubscr", http.StatusOK)

	requestNfType := "AUSF"
	targetServiceNames := []string{"nudm-auth"}
	targetNf := structs.TargetNf{
		RequesterNfType:          "AUSF",
		TargetNfType:             "UDM",
		TargetServiceNames:       targetServiceNames,
		SubscriptionValidityTime: 0,
	}
	cacheManager.SetTargetNf(requestNfType, targetNf)

	subscriptionID := "subscriptionTest"

	oneSubsData := &structs.OneSubscriptionData{
		RequesterNfType:   "AUSF",
		TargetNfType:      "PCF",
		TargetServiceName: "nudm-auth",
	}

	time := prolongSubscriptionFromNrf(oneSubsData, subscriptionID)
	if time != nil {
		t.Errorf("TestProlongSubscriptionFromNrf: check return time failure for wring TargetNfType.")
	}

	oneSubsData = &structs.OneSubscriptionData{
		RequesterNfType:   "AUSF",
		TargetNfType:      "UDM",
		TargetServiceName: "nudm-auth",
	}
	StubHTTPDoToNrf("PATCHSubscr", http.StatusInternalServerError)
	time = prolongSubscriptionFromNrf(oneSubsData, subscriptionID)
	if time != nil {
		t.Errorf("TestProlongSubscriptionFromNrf: check return time failure for wring TargetNfType.")
	}

	StubHTTPDoToNrf("PATCHSubscr", http.StatusNoContent)
	time = prolongSubscriptionFromNrf(oneSubsData, subscriptionID)
	if time == nil {
		t.Errorf("TestProlongSubscriptionFromNrf: check return time failure for StatusNoContent")
	}
	StubHTTPDoToNrf("PATCHSubscr", http.StatusOK)
	time = prolongSubscriptionFromNrf(oneSubsData, subscriptionID)
	if time == nil || !strings.HasPrefix(time.String(), "2019-04-02 17:11:23") {
		t.Errorf("TestProlongSubscriptionFromNrf: check return time failure. time(%v)", time)
	}
}

func TestProberNfProfile(t *testing.T) {
	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()
	StubHTTPDoToNrf("GET", http.StatusNotModified)

	requestNfType := "AUSF"

	validityPeriod, ret := proberNfProfile("AUSF", "UDM", "InstanceIDTest")
	if validityPeriod != 0 || ret {
		t.Errorf("TestProberNfProfile: GetTargetNf failure check return value.")
	}

	targetServiceNames := []string{"nudm-auth"}
	targetNf := structs.TargetNf{
		RequesterNfType:          "AUSF",
		TargetNfType:             "UDM",
		TargetServiceNames:       targetServiceNames,
		SubscriptionValidityTime: 0,
	}
	cacheManager.SetTargetNf(requestNfType, targetNf)

	validityPeriod, ret = proberNfProfile("AUSF", "UDM", "InstanceIDTest")

	if validityPeriod != 0 || ret {
		t.Errorf("TestProberNfProfile: GetRequesterFqdn failure check return value.")
	}
	ausfFqdn := "seliius03696.seli.gic.ericsson.se"
	cacheManager.SetRequesterFqdn("AUSF", ausfFqdn)

	validityPeriod, ret = proberNfProfile("AUSF", "UDM", "InstanceIDTest")
	if validityPeriod != 86400 || !ret {
		t.Errorf("TestProberNfProfile: check retCode StatusNotModified failure.")
	}

	StubHTTPDoToNrf("GET", http.StatusOK)
	validityPeriod, ret = proberNfProfile("AUSF", "UDM", "InstanceIDTest")
	if validityPeriod != 43200 || !ret {
		t.Errorf("TestProberNfProfile: check validityPeriod ret value failure.")
	}
}

func TestDoSubscriptionToNrf(t *testing.T) {
	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()

	StubHTTPDoToNrf("POST", http.StatusInternalServerError)
	stubNotifIPEndPoint()

	subscrIDPrefix := "subscriptions"

	oneSubsData := &structs.OneSubscriptionData{
		RequesterNfType:   "AUSF",
		TargetNfType:      "PCF",
		TargetServiceName: "nudm-auth",
	}
	ausfFqdn := "seliius03696.seli.gic.ericsson.se"
	cacheManager.SetRequesterFqdn("AUSF", ausfFqdn)

	subscrID, time := doSubscriptionToNrf(oneSubsData)
	if subscrID != "" || time != nil {
		t.Errorf("TestDoSubscriptionToNrf: check return value failure.")
	}

	StubHTTPDoToNrf("POST", http.StatusCreated)
	subscrID, time = doSubscriptionToNrf(oneSubsData)
	if !strings.HasPrefix(subscrID, subscrIDPrefix) || time == nil || !strings.HasPrefix(time.String(), "2019-04-02 17:11:23") {
		t.Errorf("TestDoSubscriptionToNrf: check return value failure.")
	}
}

func TestFetchSubscriptionInfoFromMaster(t *testing.T) {
	fHTTPDo := client.HTTPDo
	defer func() {
		client.HTTPDo = fHTTPDo
	}()
	StubGetLeader("127.0.0.1")

	subscrInfo := &structs.SubscriptionInfo{
		RequesterNfType:   "AUSF",
		TargetNfType:      "PCF",
		TargetServiceName: "nudm-auth",
		SubscriptionID:    "SubscriptionTest",
	}

	StubHTTPDoToMaster("GETSubscr", http.StatusInternalServerError)

	subscrID, time := fetchSubscriptionInfoFromMaster(subscrInfo, false)
	if subscrID != "" || time != nil {
		t.Errorf("TestFetchSubscriptionInfoFromMaster: check return value failure.")
	}

	StubHTTPDoToMaster("GETSubscr", http.StatusOK)

	subscrID, time = fetchSubscriptionInfoFromMaster(subscrInfo, false)
	if subscrID != "subscriptionTest" || time == nil || !strings.HasPrefix(time.String(), "2019-04-02 17:11:26") {
		t.Errorf("TestFetchSubscriptionInfoFromMaster: check return value failure.")
	}
}

func TestSyncNrfData(t *testing.T) {
	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()

	StubHTTPDoToNrf("GET", http.StatusInternalServerError)

	targetServiceNames := []string{"nudm-auth"}
	targetNf := &structs.TargetNf{
		RequesterNfType:          "AUSF",
		TargetNfType:             "UDM",
		TargetServiceNames:       targetServiceNames,
		SubscriptionValidityTime: 0,
	}
	ret := SyncNrfData(targetNf, false, nil)
	if ret {
		t.Errorf("TestSyncNrfData: check RequesterFqdn result code failure.")
	}

	ausfFqdn := "seliius03696.seli.gic.ericsson.se"
	cacheManager.SetRequesterFqdn("AUSF", ausfFqdn)
	ret = SyncNrfData(targetNf, false, nil)
	if ret {
		t.Errorf("TestSyncNrfData: result code StatusInternalServerError check failure.")
	}

	StubHTTPDoToNrf("GET", http.StatusOK)
	ret = SyncNrfData(targetNf, false, nil)
	if !ret {
		t.Errorf("TestSyncNrfData: check result code failure.")
	}
}

func TestUpdateConfigmapStorage(t *testing.T) {
	ret := updateConfigmapStorage(nil)
	if ret {
		t.Logf("TestUpdateConfigmapStorage: input nil ptr check failure.")
	}
	subscrInfo := structs.SubscriptionInfo{
		RequesterNfType:   "AUSF",
		TargetNfType:      "PCF",
		TargetServiceName: "nudm-auth",
		SubscriptionID:    "SubscriptionTest",
	}
	subscrInfoMap := map[string]structs.SubscriptionInfo{}
	subscrInfoMap["AUSF"] = subscrInfo
	ret = updateConfigmapStorage(subscrInfoMap)
	if ret {
		t.Errorf("TestUpdateConfigmapStorage: input PatchConfigMap failure check.")
	}
	k8sapiclient.PatchConfigMapStub(nil)
	ret = updateConfigmapStorage(subscrInfoMap)
	if !ret {
		t.Errorf("TestUpdateConfigmapStorage: return value check failure.")
	}
	cacheManager.Flush("AUSF")
}

func TestGetNfType(t *testing.T) {
	nfType, err := getNfType("AUSF", "udm-5g-01")
	if nfType != "" || err == nil {
		t.Errorf("TestGetNfType: no cache check failure.")
	}
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestGetNfType: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestGetNfType: Cached fail")
		}
	}
	nfType, err = getNfType("AUSF", "udm-5g-01")
	t.Logf("TestGetNfType: %s,err(%v)", nfType, err)
	if nfType != "UDM" || err != nil {
		t.Errorf("TestGetNfType: check return value failure.")
	}
	cacheManager.Flush("AUSF")
}

/*
func TestCacheClean(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestCacheClean: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestCacheClean: Cached fail")
		}
	}
	cacheClean("AUSF", "UDM", "udm-5g-01")
	ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
	if ok {
		t.Errorf("TestCacheClean: check cache clean fail")
	}
	cacheManager.Flush("AUSF")
}
*/
