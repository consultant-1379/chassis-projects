package cache

import (
	"testing"
	"time"

	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
)

var contChfReg = []byte(`{
    "validityPeriod": 86400,
    "nfInstances": [{
        "nfInstanceId": "chf-5g-01",
        "nfType": "CHF",
        "plmnList": [
		   {
            "mcc": "460",
            "mnc": "000"
           },
		   {
            "mcc": "560",
            "mnc": "001"
           }
		],
        "sNssais": [{
                "sst": 2,
                "sd": "2"
            },
            {
                "sst": 4,
                "sd": "4"
            }
        ],
        "nsiList": ["100","101","102"],
        "fqdn": "seliius03695.seli.gic.ericsson.se",
        "ipv4Addresses": ["172.16.208.1"],
        "ipv6Addresses": ["FF01::1101"],
        "ipv6Prefixes": ["FF01"],
        "capacity": 100,
        "load": 50,
        "locality": "Shanghai",
        "priority": 1,
        "chfInfo" : {
              "plmnRangeList": [
              {
                  "start": "46000",
                  "end": "46011"
              },
              {
                  "pattern": "^46[3-4]{1}[0-9]{2,3}$"
              },
              {
                  "start": "460111",
                  "end": "460222"
              }]
              },
        "nfServices": [{
            "serviceInstanceId": "nchf-auth-01",
            "serviceName": "nchf-auth-01",
            "version": [{
                "apiVersionInUri": "v1Url",
                "apiFullVersion": "v1"
            }],
            "schema": "https://",
            "fqdn": "seliius03690.seli.gic.ericsson.se",
            "ipEndPoints": [{
                "ipv4Address": "10.210.121.75",
                "ipv6Address": "FF01::1101",
                "ipv6Prefix": "FF01",
                "transport": "TCP",
                "port": 30080
            }],
            "apiPrefix": "nudm-uecm",
            "defaultNotificationSubscriptions": [{
                "notificationType": "N1_MESSAGES",
                "callbackUri": "https://seliius03695.seli.gic.ericsson.se",
                "n1MessageClass": "5GMM",
                "n2InformationClass": "SM"
            }],
            "capacity": 0,
            "load": 50,
            "priority": 0,
            "supportedFeatures": "A1"
        }]
    }]
}
`)

func TestFetchProfileIds(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestFetchProfileIds: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestFetchProfileIds: Cached fail")
		}
	}
	t.Errorf("TestFetchProfileIds: Cached fail")
	ids := cacheManager.mcache["AUSF"].fetchProfileIDs("UDM")
	if len(ids) != 1 || ids[0] != "udm-5g-01" {
		t.Errorf("TestFetchProfileIds: FetchProfileIds fail")
	}
	requestNfType := "AUSF"
	targetServiceNames := []string{"nudm-auth"}
	targetNf := structs.TargetNf{
		RequesterNfType:          "AUSF",
		TargetNfType:             "UDM",
		TargetServiceNames:       targetServiceNames,
		SubscriptionValidityTime: 0,
	}
	cacheManager.SetTargetNf(requestNfType, targetNf)
	dumpData := structs.CacheDumpData{
		RequestNfType: "AUSF",
	}
	cacheManager.Dump("AUSF", &dumpData)
	if len(dumpData.CacheInfos) == 0 {
		t.Errorf("TestFetchProfileIds: Dump fail")
	}
	cacheManager.Flush("AUSF")
}

func TestFetchRoamingProfileIDs(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmRegRoam)
	if nfinstanceByte == nil {
		t.Errorf("TestFetchProfileIds: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, true)
		ok := cacheManager.ProbeRoam("AUSF", "UDM", "udm-5g-02")
		if !ok {
			t.Errorf("TestFetchProfileIds: Cached fail")
		}
	}
	ids := cacheManager.roamingCache["AUSF"].fetchProfileIDs("UDM")
	if len(ids) != 1 || ids[0] != "udm-5g-02" {
		t.Errorf("TestFetchProfileIds: FetchProfileIds fail")
	}
	requestNfType := "AUSF"
	targetServiceNames := []string{"nudm-auth"}
	targetNf := structs.TargetNf{
		RequesterNfType:          "AUSF",
		TargetNfType:             "UDM",
		TargetServiceNames:       targetServiceNames,
		SubscriptionValidityTime: 0,
	}
	cacheManager.SetTargetNf(requestNfType, targetNf)
	dumpData := structs.CacheDumpData{
		RequestNfType: "AUSF",
	}
	cacheManager.Dump("AUSF", &dumpData)
	if len(dumpData.RoamingCacheInfos) == 0 {
		t.Errorf("TestFetchProfileIds: Dump fail")
	}
	cacheManager.FlushRoam("AUSF")
}

func TestCachedWithTTL(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestCachedWithTTL: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.CachedWithTTL("AUSF", "UDM", instance, 3, false)
		//ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		ok := false
		if !ok {
			t.Errorf("TestCachedWithTTL: CachedWithTTL cache fail")
		}
	}
	cacheManager.Flush("AUSF")
}

func TestSearchNormal(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestSearchNormal: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
	}

	searchConditionsWrong := SearchParameter{}
	searchConditionsWrong.targetNfType = "UDM"
	searchConditionsWrong.serviceNames = []string{"ttnudm-auth-01"} //wrong service name ttudm-5g-01
	searchConditionsWrong.requesterNfType = "udm"
	content, ok := cacheManager.Search("AUSF", "UDM", &searchConditionsWrong, false)
	if ok || len(content) != 0 {
		t.Errorf("TestSearchNormal: Search wrong failed")
	}

	searchConditions := SearchParameter{}
	searchConditions.targetNfType = "UDM"
	searchConditions.serviceNames = []string{"nudm-auth-01"}
	searchConditions.supportedFeatures = "A1"
	content, ok = cacheManager.Search("AUSF", "UDM", &searchConditions, false)
	if !ok || len(content) == 0 {
		t.Errorf("TestSearchNormal: Search failed")
	}
	cacheManager.Flush("AUSF")
}

//SupportFeature should align with ServiceName
func TestSearchSupportFeature(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestSearchSupportFeature: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
	}

	searchConditionsWrong := SearchParameter{}
	searchConditionsWrong.targetNfType = "UDM"
	searchConditionsWrong.serviceNames = []string{"ttnudm-auth-01"} //wrong service name ttudm-5g-01
	searchConditionsWrong.requesterNfType = "udm"
	content, ok := cacheManager.Search("AUSF", "UDM", &searchConditionsWrong, false)
	if ok || len(content) != 0 {
		t.Errorf("TestSearchSupportFeature: Search wrong failed")
	}

	searchConditions := SearchParameter{}
	searchConditions.targetNfType = "UDM"
	searchConditions.serviceNames = []string{"nudm-auth-01"}
	searchConditions.supportedFeatures = "A200"
	//Search should be failure because service-name nudm-auth-01 support feature value is A1.
	content, ok = cacheManager.Search("AUSF", "UDM", &searchConditions, false)
	if ok || len(content) != 0 {
		t.Errorf("TestSearch: Search result should be failed")
	}
	cacheManager.Flush("AUSF")
}

func TestSearchRoutingIndicator(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestSearchRoutingIndicator: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
	}

	// search routing-indicator mismatch
	searchConditions := SearchParameter{}
	searchConditions.targetNfType = "UDM" //udm-5g-01
	searchConditions.serviceNames = []string{"nudm-auth-01"}
	searchConditions.routingIndicator = "123"
	if content, ok := cacheManager.Search("AUSF", "UDM", &searchConditions, false); ok || len(content) != 0 {
		t.Errorf("TestSearchRoutingIndicator: Search routing-indicator mismatch failure")
	}
	// search routing-indicator match
	searchConditions.routingIndicator = "1111"
	if content, ok := cacheManager.Search("AUSF", "UDM", &searchConditions, false); !ok || len(content) == 0 {
		t.Errorf("TestSearchRoutingIndicator: Search failed")
	}
	searchConditions.routingIndicator = "1234"
	if content, ok := cacheManager.Search("AUSF", "UDM", &searchConditions, false); !ok || len(content) == 0 {
		t.Errorf("TestSearchRoutingIndicator: Search failed")
	}
	// search routing-indicator match
	searchConditions.routingIndicator = "5678"
	if content, ok := cacheManager.Search("AUSF", "UDM", &searchConditions, false); !ok || len(content) == 0 {
		t.Errorf("TestSearchRoutingIndicator: Search failed")
	}
	cacheManager.Flush("AUSF")
}

func TestSearchGroupIDList(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestSearchRoutingIndicator: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
	}

	// search routing-indicator mismatch
	searchConditions := SearchParameter{}
	searchConditions.targetNfType = "UDM" //udm-5g-01
	searchConditions.serviceNames = []string{"nausf-auth-01"}
	searchConditions.groupIDList = []string{"udmxxx"}
	content, ok := cacheManager.Search("AUSF", "UDM", &searchConditions, false)
	if ok || len(content) != 0 {
		t.Errorf("TestSearchGroupIDList: Search group-id-list mismatch failure")
	}

	searchConditions2 := SearchParameter{}
	searchConditions2.targetNfType = "UDM" //udm-5g-01
	searchConditions2.serviceNames = []string{"nudm-auth-01"}
	searchConditions2.groupIDList = []string{"udmtest"}
	content2, ok2 := cacheManager.Search("AUSF", "UDM", &searchConditions2, false)
	if !ok2 || len(content2) == 0 {
		t.Errorf("TestSearchRoutingIndicator: Search failed")
	}
	cacheManager.Flush("AUSF")
}

func TestAssembleResponseContents(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestAssembleResponseContents: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
	}

	ids := cacheManager.mcache["AUSF"].fetchProfileIDs("UDM")

	searchConditionsWrong := SearchParameter{}
	searchConditionsWrong.targetNfType = "udm"                      //udm-5g-01
	searchConditionsWrong.serviceNames = []string{"ttnudm-auth-01"} //wrong service name ttudm-5g-01
	searchConditionsWrong.requesterNfType = "udm"
	content := cacheManager.assembleResponseContents("AUSF", "UDM", ids, &searchConditionsWrong, false, false)
	if len(content) != 0 {
		t.Errorf("TestAssembleResponseContents: assembleResponseContents wrong failed")
	}

	searchConditions := SearchParameter{}
	searchConditions.targetNfType = "udm" //udm-5g-01
	searchConditions.serviceNames = []string{"nudm-auth-01"}
	searchConditions.requesterNfType = "udm"
	content = cacheManager.assembleResponseContents("AUSF", "UDM", ids, &searchConditions, false, false)
	if len(content) == 0 {
		t.Errorf("TestAssembleResponseContents: assembleResponseContents failed")
	}
	cacheManager.Flush("AUSF")
}

func TestSearchChfSupportedPlmn(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contChfReg)
	if nfinstanceByte == nil {
		t.Errorf("TestSearchRoutingIndicator: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
	}

	// search routing-indicator mismatch
	searchConditions := SearchParameter{}
	searchConditions.targetNfType = "CHF"
	searchConditions.serviceNames = []string{"nchf-auth-01"}
	searchConditions.chfSupportedPlmn = structs.PlmnID{
		Mcc: "460",
		Mnc: "10",
	}
	content, ok := cacheManager.Search("AUSF", "UDM", &searchConditions, false)
	if !ok || len(content) == 0 {
		t.Errorf("TestSearchGroupIDList: Search group-id-list match failure")
	}

	searchConditions.chfSupportedPlmn = structs.PlmnID{
		Mcc: "463",
		Mnc: "999",
	}
	content, ok = cacheManager.Search("AUSF", "UDM", &searchConditions, false)

	if !ok || len(content) == 0 {
		t.Errorf("TestSearchGroupIDList: Search group-id-list mismatch failure")
	}

	searchConditions.chfSupportedPlmn = structs.PlmnID{
		Mcc: "460",
		Mnc: "331",
	}
	content, ok = cacheManager.Search("AUSF", "UDM", &searchConditions, false)
	if ok || len(content) != 0 {
		t.Errorf("TestSearchGroupIDList: Search group-id-list mismatch failure")
	}

	searchConditions.chfSupportedPlmn = structs.PlmnID{
		Mcc: "465",
		Mnc: "99",
	}
	content, ok = cacheManager.Search("AUSF", "UDM", &searchConditions, false)
	if ok || len(content) != 0 {
		t.Errorf("TestSearchGroupIDList: Search group-id-list mismatch failure")
	}
	cacheManager.Flush("AUSF")
}

func TestCalculateTtl(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestCalculateTtl: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.CachedWithTTL("AUSF", "UDM", instance, 60, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestCalculateTtl: CachedWithTTL cache fail")
		}
	}
	time.Sleep(1 * time.Second)
	ids := cacheManager.mcache["AUSF"].fetchProfileIDs("UDM")
	leftTime := cacheManager.calculateTTL("AUSF", "UDM", ids, false)
	if leftTime != 60-1 {
		t.Errorf("TestCalculateTtl: CalculateTtl left time is not correct.")
	}
	cacheManager.Flush("AUSF")
}

func TestDeCachedByNfType(t *testing.T) {
	requestNfType := "AUSF"
	targetServiceNames := []string{"nudm-auth"}
	targetNf := structs.TargetNf{
		RequesterNfType:          "AUSF",
		TargetNfType:             "UDM",
		TargetServiceNames:       targetServiceNames,
		SubscriptionValidityTime: 0,
	}
	cacheManager.SetTargetNf(requestNfType, targetNf)
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestDeCachedByNfType: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestDeCachedByNfType: Cached fail")
		}
	}
	cacheManager.DeCachedByNfType("UDM")
	cacheManager.DeCachedByNfType("AUSF")
	ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
	if ok {
		t.Errorf("TestDeCachedByNfType: TestDeCachedByNfType fail")
	}
	cacheManager.DeCachedByNfType("AUSF")
	cacheManager.Flush("AUSF")
}

func TestDump(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestDump: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestDump: Cached fail")
		}
	}

	dumpData1 := structs.CacheDumpData{
		RequestNfType: "AUSF",
	}
	cacheManager.Dump("AUSF", &dumpData1)

	if len(dumpData1.CacheInfos) != 1 {
		t.Fatal("TestDump: Dump AUSF check fail")
	}

	dumpData2 := structs.CacheDumpData{
		RequestNfType: "UDM",
	}
	cacheManager.Dump("UDM", &dumpData2)
	if len(dumpData2.CacheInfos) != 0 {
		t.Fatal("TestDump: Dump UDM check fail")
	}

	dumpDataList1 := make([]structs.CacheDumpData, 0)
	cacheManager.DumpAll(&dumpDataList1)
	if len(dumpDataList1) != 1 {
		t.Fatal("TestDump: DumpAll succ check fail")
	}

	requestNfType := "AUSF"
	targetServiceNames := []string{"nudm-auth"}
	targetNf := structs.TargetNf{
		RequesterNfType:          "AUSF",
		TargetNfType:             "UDM",
		TargetServiceNames:       targetServiceNames,
		SubscriptionValidityTime: 0,
	}
	cacheManager.SetTargetNf(requestNfType, targetNf)
	cacheManager.DeCachedByNfType("AUSF")

	dumpDataList2 := make([]structs.CacheDumpData, 0)
	cacheManager.DumpAll(&dumpDataList2)
	if len(dumpDataList2[0].CacheInfos[0].NfProfiles) != 0 {
		t.Errorf("TestDump: DumpAll not found check fail")
	}

	cacheManager.Flush("AUSF")
}

func TestSync(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestSync: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.Cached("AUSF", "UDM", instance, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Fatal("TestSync: Cached fail")
		}
	}

	syncData1 := structs.CacheSyncData{
		RequestNfType: "AUSF",
	}
	cacheManager.Sync("AUSF", &syncData1)
	if len(syncData1.CacheInfos) != 1 {
		t.Fatal("TestSync: sync ausf check fail")
	}

	syncData2 := structs.CacheSyncData{
		RequestNfType: "UDM",
	}
	cacheManager.Sync("UDM", &syncData2)
	if len(syncData2.CacheInfos) != 0 {
		t.Fatal("TestDump: sync UDM check fail")
	}

	requestNfType := "AUSF"
	targetServiceNames := []string{"nudm-auth"}
	targetNf := structs.TargetNf{
		RequesterNfType:          "AUSF",
		TargetNfType:             "UDM",
		TargetServiceNames:       targetServiceNames,
		SubscriptionValidityTime: 0,
	}
	cacheManager.SetTargetNf(requestNfType, targetNf)
	cacheManager.DeCachedByNfType("AUSF")
	syncData3 := structs.CacheSyncData{
		RequestNfType: "AUSF",
	}
	cacheManager.Sync("AUSF", &syncData3)
	if len(syncData3.CacheInfos[0].NfProfiles) != 0 {
		t.Fatal("TestSync: sync ausf check fail")
	}

	cacheManager.Flush("AUSF")
}

func TestCacheProfileDiffHandler(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestCacheProfileDiffHandler: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.CachedWithTTL("AUSF", "UDM", instance, 3, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestCacheProfileDiffHandler: CachedWithTTL cache fail")
		}
	}

	ret := cacheManager.cacheProfileDiffHandler("AUSF", "UDM", nil)
	if ret {
		t.Errorf("TestCacheProfileDiffHandler: newCache nil check failure")
	}

	newCache := cacheManager.buildCache("AUSF", "UDM", nil, homeCache)
	ret = cacheManager.cacheProfileDiffHandler("AUSF", "UDM", newCache)
	if !ret {
		t.Errorf("TestCacheProfileDiffHandler: newCache send dereg nfprofile failure")
	}
	cacheManager.Flush("AUSF")
}

func TestEnterNormalWorkModeMaster(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestEnterNormalWorkMode: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.CachedWithTTL("AUSF", "UDM", instance, 3, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestEnterNormalWorkMode: CachedWithTTL cache fail")
		}
	}
	requestNfType := "AUSF"
	targetServiceNames := []string{"nudm-auth"}
	targetNf := structs.TargetNf{
		RequesterNfType:          "AUSF",
		TargetNfType:             "UDM",
		TargetServiceNames:       targetServiceNames,
		SubscriptionValidityTime: 0,
	}
	cacheManager.SetTargetNf(requestNfType, targetNf)
	cacheManager.EnterKeepCacheWorkMode()
	cacheManager.EnterNormalWorkMode()
	ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
	if ok {
		t.Errorf("TestEnterNormalWorkMode: master cache not cleanup when enter nroaml work mode")
	}
	cacheManager.Flush("AUSF")
}

func TestEnterNormalWorkModeSlave(t *testing.T) {
	activeLeaderMock(false)
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestEnterNormalWorkMode: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.CachedWithTTL("AUSF", "UDM", instance, 3, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestEnterNormalWorkMode: CachedWithTTL cache fail")
		}
	}
	requestNfType := "AUSF"
	targetServiceNames := []string{"nudm-auth"}
	targetNf := structs.TargetNf{
		RequesterNfType:          "AUSF",
		TargetNfType:             "UDM",
		TargetServiceNames:       targetServiceNames,
		SubscriptionValidityTime: 0,
	}
	cacheManager.SetTargetNf(requestNfType, targetNf)
	cacheManager.EnterKeepCacheWorkMode()
	cacheManager.EnterNormalWorkMode()
	ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
	if !ok {
		t.Errorf("TestEnterNormalWorkMode: slave cache is cleanup when enter nroaml work mode")
	}
	activeLeaderMock(true)
	cacheManager.Flush("AUSF")
}

func TestSuperviseSubscription(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestCachedWithTTL: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.CachedWithTTL("AUSF", "UDM", instance, 3, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestCachedWithTTL: CachedWithTTL cache fail")
		}
	}

	requesterNfType := "AUSF"
	targetNfType := "UDM"
	subscriptionID := "20bd0bb9-edc1-4c74-8ec5-74e4fed79ac8"
	timepoint := time.Now()

	cacheManager.SuperviseSubscription(requesterNfType, targetNfType, subscriptionID, timepoint)
}

func TestSuperviseRoamingSubscription(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmRegRoam)
	if nfinstanceByte == nil {
		t.Errorf("TestCachedWithTTL: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.CachedWithTTL("AUSF", "UDM", instance, 3, true)
		ok := cacheManager.ProbeRoam("AUSF", "UDM", "udm-5g-02")
		if !ok {
			t.Errorf("TestCachedWithTTL: CachedWithTTL cache fail")
		}
	}

	requesterNfType := "AUSF"
	targetNfType := "UDM"
	subscriptionID := "20bd0bb9-edc1-4c74-8ec5-74e4fed79ac9"
	timepoint := time.Now()

	cacheManager.SuperviseRoamingSubscription(requesterNfType, targetNfType, subscriptionID, timepoint)
}

func TestDelSubscriptionMonitor(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestCachedWithTTL: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.CachedWithTTL("AUSF", "UDM", instance, 3, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestCachedWithTTL: CachedWithTTL cache fail")
		}
	}

	requesterNfType := "AUSF"
	targetNfType := "UDM"
	subscriptionID := "20bd0bb9-edc1-4c74-8ec5-74e4fed79ac8"

	cacheManager.DelSubscriptionMonitor(requesterNfType, targetNfType, subscriptionID)
}

func TestDelRoamingSubscriptionMonitor(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmRegRoam)
	if nfinstanceByte == nil {
		t.Errorf("TestCachedWithTTL: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.CachedWithTTL("AUSF", "UDM", instance, 3, true)
		ok := cacheManager.ProbeRoam("AUSF", "UDM", "udm-5g-02")
		if !ok {
			t.Errorf("TestCachedWithTTL: CachedWithTTL cache fail")
		}
	}

	requesterNfType := "AUSF"
	targetNfType := "UDM"
	subscriptionID := "20bd0bb9-edc1-4c74-8ec5-74e4fed79ac9"

	cacheManager.DelRoamingSubscriptionMonitor(requesterNfType, targetNfType, subscriptionID)
}

func TestProbeAllCache(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestCachedWithTTL: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.CachedWithTTL("AUSF", "UDM", instance, 3, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestCachedWithTTL: CachedWithTTL cache fail")
		}
	}

	requesterNfType := "AUSF"
	targetNfType := "UDM"
	nfInstanceID := "udm-5g-01"

	exist, isRoam := cacheManager.ProbeAllCache(requesterNfType, targetNfType, nfInstanceID)
	if !exist || isRoam {
		t.Fatal("Expect probe nfInstanceID in normal cache")
	}

	nfInstanceID = "udm-5g-03"
	exist, isRoam = cacheManager.ProbeAllCache(requesterNfType, targetNfType, nfInstanceID)
	if exist || isRoam {
		t.Fatal("Expect probe no nfInstanceID in normal cache and roaming cache")
	}

	nfinstanceRoamByte, _, _ := SpliteSeachResult(contUdmRegRoam)
	if nfinstanceRoamByte == nil {
		t.Errorf("TestCachedWithTTL: SpliteSeachResult fail")
	}
	for _, instanceRoam := range nfinstanceRoamByte {
		cacheManager.CachedWithTTL("AUSF", "UDM", instanceRoam, 3, true)
		ok := cacheManager.ProbeRoam("AUSF", "UDM", "udm-5g-02")
		if !ok {
			t.Errorf("TestCachedWithTTL: CachedWithTTL cache fail")
		}
	}

	nfInstanceID = "udm-5g-02"
	exist, isRoam = cacheManager.ProbeAllCache(requesterNfType, targetNfType, nfInstanceID)
	if !exist || !isRoam {
		t.Fatal("Expect probe nfInstanceID in roaming cache")
	}
}

func TestGetProfileByID(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmReg)
	if nfinstanceByte == nil {
		t.Errorf("TestCachedWithTTL: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.CachedWithTTL("AUSF", "UDM", instance, 3, false)
		ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
		if !ok {
			t.Errorf("TestCachedWithTTL: CachedWithTTL cache fail")
		}
	}

	requesterNfType := "AUSF"
	targetNfType := "UDM"
	nfInstanceID := "udm-5g-01"
	isRoam := false

	profile := cacheManager.GetProfileByID(requesterNfType, targetNfType, nfInstanceID, isRoam)

	if profile == nil {
		t.Fatal("Expect normal profile, but not")
	}

	nfinstanceRoamByte, _, _ := SpliteSeachResult(contUdmRegRoam)
	if nfinstanceRoamByte == nil {
		t.Errorf("TestCachedWithTTL: SpliteSeachResult fail")
	}
	for _, instanceRoam := range nfinstanceRoamByte {
		cacheManager.CachedWithTTL("AUSF", "UDM", instanceRoam, 3, true)
		ok := cacheManager.ProbeRoam("AUSF", "UDM", "udm-5g-02")
		if !ok {
			t.Errorf("TestCachedWithTTL: CachedWithTTL cache fail")
		}
	}

	nfInstanceID = "udm-5g-02"
	isRoam = true

	profile = cacheManager.GetProfileByID(requesterNfType, targetNfType, nfInstanceID, isRoam)

	if profile == nil {
		t.Fatal("Expect roaming profile, but not")
	}
}

func TestCachedWithTtlTimestamp(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmRegRoam)
	if nfinstanceByte == nil {
		t.Errorf("TestCachedWithTTL: SpliteSeachResult fail")
	}

	timepoint := time.Now()
	for _, instance := range nfinstanceByte {
		cacheManager.CachedWithTtlTimestamp("AUSF", "UDM", instance, timepoint, true)
		ok := cacheManager.ProbeRoam("AUSF", "UDM", "udm-5g-02")
		if !ok {
			t.Errorf("TestCachedWithTTL: CachedWithTTL cache fail")
		}
	}
}

func TestDeCachedCacheManager(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmRegRoam)
	if nfinstanceByte == nil {
		t.Errorf("TestCachedWithTTL: SpliteSeachResult fail")
	}

	timepoint := time.Now()
	for _, instance := range nfinstanceByte {
		cacheManager.CachedWithTtlTimestamp("AUSF", "UDM", instance, timepoint, true)
		ok := cacheManager.ProbeRoam("AUSF", "UDM", "udm-5g-02")
		if !ok {
			t.Errorf("TestCachedWithTTL: CachedWithTTL cache fail")
		}
	}

	requesterNfType := "AUSF"
	targetNfType := "UDM"
	nfInstanceID := "udm-5g-02"
	isRoaming := true
	cacheManager.DeCached(requesterNfType, targetNfType, nfInstanceID, isRoaming)
}

func TestReCached(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmRegRoam)
	if nfinstanceByte == nil {
		t.Errorf("TestCachedWithTTL: SpliteSeachResult fail")
	}

	timepoint := time.Now()
	for _, instance := range nfinstanceByte {
		cacheManager.CachedWithTtlTimestamp("AUSF", "UDM", instance, timepoint, true)
		ok := cacheManager.ProbeRoam("AUSF", "UDM", "udm-5g-02")
		if !ok {
			t.Errorf("TestCachedWithTTL: CachedWithTTL cache fail")
		}
	}

	nfInstanceID := "udm-5g-02"
	nfinstanceByte, _, _ = SpliteSeachResult(contUdmRegRoam)
	if nfinstanceByte == nil {
		t.Errorf("TestCachedWithTTL: SpliteSeachResult fail")
	}
	for _, instance := range nfinstanceByte {
		cacheManager.ReCached("AUSF", "UDM", nfInstanceID, instance, true)
		ok := cacheManager.ProbeRoam("AUSF", "UDM", "udm-5g-02")
		if !ok {
			t.Errorf("TestCachedWithTTL: CachedWithTTL cache fail")
		}
	}
}

func TestSearchRoamingCache(t *testing.T) {
	nfinstanceByte, _, _ := SpliteSeachResult(contUdmRegRoam)
	if nfinstanceByte == nil {
		t.Errorf("TestCachedWithTTL: SpliteSeachResult fail")
	}

	timepoint := time.Now()
	for _, instance := range nfinstanceByte {
		cacheManager.CachedWithTtlTimestamp("AUSF", "UDM", instance, timepoint, true)
	}

	searchConditionsWrong := SearchParameter{}
	searchConditionsWrong.targetNfType = "UDM"
	searchConditionsWrong.serviceNames = []string{"ttnudm-auth-02"} //wrong service name ttudm-5g-02
	searchConditionsWrong.requesterNfType = "udm"
	content, ok := cacheManager.SearchRoamingCache("AUSF", "UDM", &searchConditionsWrong, false)
	if ok || len(content) != 0 {
		t.Errorf("TestSearchRoaming: Search wrong failed")
	}

	searchConditions := SearchParameter{}
	searchConditions.targetNfType = "UDM"
	searchConditions.serviceNames = []string{"nudm-auth-02"}
	searchConditions.supportedFeatures = "A1"
	content, ok = cacheManager.SearchRoamingCache("AUSF", "UDM", &searchConditions, false)
	if !ok || len(content) == 0 {
		t.Errorf("TestSearchRoaming: Search failed")
	}
	cacheManager.FlushRoam("AUSF")
}
