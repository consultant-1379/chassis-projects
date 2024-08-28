package worker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/utils"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/client"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
)

func TestLoopFetchTargetNfs(t *testing.T) {
	nfType := "AUSF"
	loopFetchTargetNfs(nfType)
	t.Logf("Loop FetchTargetNfs by nfType:%s success", nfType)
}

func TestSubscribe(t *testing.T) {
	defer workerManager.resetSubscribeBacklogTask("AUSF")
	defer workerManager.resetFetchProfileBacklogTask("AUSF")

	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()

	nfType := "AUSF"

	StubHTTPDoToNrf("POST", http.StatusInternalServerError)
	rest := subscribe(nfType)
	t.Logf("subscribe result : %v", rest)
	if rest {
		t.Fatal("Expect subscribe result is failure, but success")
	}

	StubHTTPDoToNrf("POST", http.StatusCreated)
	rest = subscribe(nfType)
	t.Logf("subscribe result : %v", rest)
	if !rest {
		t.Fatal("Expect subscribe result is success, but failure")
	}
}

func TestSubscribeExecutor(t *testing.T) {
	defer workerManager.resetSubscribeBacklogTask("AUSF")
	defer workerManager.resetFetchProfileBacklogTask("AUSF")

	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()

	nfType := "AUSF"

	oneSubsData := structs.OneSubscriptionData{
		RequesterNfType:   nfType,
		TargetNfType:      "UDM",
		TargetServiceName: "udm-servicer-1",
		NotifCondition:    nil,
	}

	StubHTTPDoToNrf("POST", http.StatusInternalServerError)
	subscriptionID, timeStamp, err := subscribeExecutor(&oneSubsData)
	t.Logf("subscriptionID:%s, timeStamp=%v, err:%v", subscriptionID, timeStamp, err)
	if subscriptionID != "" {
		t.Fatal("Expect subscriptionID is empty, but not")
	}

	StubHTTPDoToNrf("POST", http.StatusCreated)
	subscriptionID, timeStamp, err = subscribeExecutor(&oneSubsData)
	fmt.Printf("subscriptionID:%s, timeStamp=%v, err:%v\n", subscriptionID, timeStamp, err)
	if subscriptionID == "" {
		t.Fatal("Expect subscriptionID is not empty, but not")
	}
}

/*
func TestSubscribeRoamExecutor(t *testing.T) {
	defer workerManager.resetSubscribeBacklogTask("AUSF")
	defer workerManager.resetFetchProfileBacklogTask("AUSF")

	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()

	nfType := "AUSF"

	oneSubsData := structs.OneSubscriptionData{
		RequesterNfType: nfType,
		TargetNfType:    "UDM",
		NfInstanceID:    "udm-5g-01",
		NotifCondition:  nil,
	}

	plmnID := structs.PlmnID{
		Mcc: "450",
		Mnc: "000",
	}

	validityTime := "2019-04-02T17:11:28+08:00"

	StubHTTPDoToNrf("POST", http.StatusInternalServerError)
	subscriptionID, timeStamp, err := subscribeRoamExecutor(&oneSubsData, &plmnID, validityTime)
	t.Logf("subscriptionID:%s, timeStamp=%v, err:%v", subscriptionID, timeStamp, err)
	if subscriptionID != "" {
		t.Fatal("Expect subscriptionID is empty, but not")
	}

	StubHTTPDoToNrf("POST", http.StatusCreated)
	subscriptionID, timeStamp, err = subscribeRoamExecutor(&oneSubsData, &plmnID, validityTime)
	fmt.Printf("subscriptionID:%s, timeStamp=%v, err:%v\n", subscriptionID, timeStamp, err)
	if subscriptionID == "" {
		t.Fatal("Expect subscriptionID is not empty, but not")
	}
}
*/
func TestSubscribeResponseParser(t *testing.T) {
	defer workerManager.resetSubscribeBacklogTask("AUSF")
	defer workerManager.resetFetchProfileBacklogTask("AUSF")

	respMgmt := &httpclient.HttpRespData{}

	respMgmt.StatusCode = http.StatusCreated
	subID, _ := utils.GetUUIDString()
	respMgmt.Location = "http://127.0.0.1:3212/nnrf-nfm/v1/subscriptions/" + subID
	respMgmt.Body = subscrRsp

	t.Logf("subscribe response : %v", respMgmt)

	subscriptionID, timeStamp, err := subscribeResponseParser(respMgmt)
	t.Logf("subscriptionID:%s, timeStamp=%v, err:%v\n", subscriptionID, timeStamp, err)
	if subscriptionID == "" {
		t.Fatal("Expect subscriptionID is not empty, but not")
	}
}

func TestFetchProfile(t *testing.T) {
	defer workerManager.resetSubscribeBacklogTask("AUSF")
	defer workerManager.resetFetchProfileBacklogTask("AUSF")

	backupHTTPDoToNrfMgmt := client.HTTPDoToNrfMgmt
	backupHTTPDoToNrfDisc := client.HTTPDoToNrfDisc
	defer func() {
		client.HTTPDoToNrfMgmt = backupHTTPDoToNrfMgmt
		client.HTTPDoToNrfDisc = backupHTTPDoToNrfDisc
	}()

	nfType := "AUSF"

	StubHTTPDoToNrf("GET", http.StatusInternalServerError)
	rest := fetchProfile(nfType)
	t.Logf("fetchProfile result : %v", rest)
	if rest {
		t.Fatal("Expect fetchProfile failure when NRF response is 500, but success")
	}

	StubHTTPDoToNrf("GET", http.StatusNotFound)
	rest = fetchProfile(nfType)
	t.Logf("fetchProfile result : %v", rest)
	if !rest {
		t.Fatal("Expect fetchProfile success when NRF response is 404, but failure")
	}

	StubHTTPDoToNrf("GET", http.StatusOK)
	rest = fetchProfile(nfType)
	t.Logf("fetchProfile result : %v", rest)
	if !rest {
		t.Fatal("Expect fetchProfile success when NRF reponse is 200, but failure")
	}
}

func TestHandleDiscoveryResponse(t *testing.T) {
	defer workerManager.resetSubscribeBacklogTask("AUSF")
	defer workerManager.resetFetchProfileBacklogTask("AUSF")

	nfType := "AUSF"
	targetNfType := "UDM"

	respDisc := &httpclient.HttpRespData{}

	respDisc.StatusCode = http.StatusNotFound
	respDisc.Body = []byte(`Can not find the nfInstance info`)

	err := handleDiscoveryResponse(nfType, targetNfType, respDisc)
	t.Logf("handleDiscoveryResponse result:%s", err.Error())
	if err == nil {
		t.Fatal("Expect parse discovery response failure, but success")
	}

	respDisc.StatusCode = http.StatusOK
	respDisc.Body = searchResultUDM

	err = handleDiscoveryResponse(nfType, targetNfType, respDisc)
	t.Logf("handleDiscoveryResponse result:%v", err)
	if err != nil {
		t.Fatal("Expect parse discovery response success, but failure")
	}
}

func TestDumpCacheFromMaster(t *testing.T) {
	defer workerManager.resetSubscribeBacklogTask("AUSF")
	defer workerManager.resetFetchProfileBacklogTask("AUSF")

	nfType := "AUSF"

	rest := dumpCacheFromMaster(nfType)
	t.Logf("dump cache rest : %v", rest)
	if rest {
		t.Fatal("Expect dumpCacheFromMaster failure, but not")
	}
}

func TestHandleDumpResponse(t *testing.T) {
	defer workerManager.resetSubscribeBacklogTask("AUSF")
	defer workerManager.resetFetchProfileBacklogTask("AUSF")

	now := time.Now()
	validityTime := now.Add(time.Duration(3600) * time.Second)

	cacheInfos := make([]structs.CacheSyncInfo, 0)

	nfProfiles := make([][]byte, 0)
	nfProfiles = append(nfProfiles, udmProfile)

	ttlInfos := make([]structs.TtlInfo, 0)
	ttlInfo := structs.TtlInfo{
		NfInstanceID: "5g-udm-01",
		ValidityTime: validityTime,
	}
	ttlInfos = append(ttlInfos, ttlInfo)

	subscriptionInfos := make([]structs.SubscriptionInfo, 0)
	subscriptionInfo := structs.SubscriptionInfo{
		RequesterNfType:   "AUSF",
		TargetNfType:      "UDM",
		TargetServiceName: "udm-servicer-1",
		NotifCondition:    nil,
		SubscriptionID:    "subscriptions/subscription-000000000001",
		ValidityTime:      validityTime,
	}
	subscriptionInfos = append(subscriptionInfos, subscriptionInfo)

	etagInfos := make([]structs.EtagInfo, 0)
	etagInfo := structs.EtagInfo{
		NfInstanceID: "5g-udm-01",
		FingerPrint:  "askfdja-akfja-aflakjd",
	}
	etagInfos = append(etagInfos, etagInfo)

	cacheSyncInfo := structs.CacheSyncInfo{
		TargetNfType:      "UDM",
		NfProfiles:        nfProfiles,
		TtlInfos:          ttlInfos,
		SubscriptionInfos: subscriptionInfos,
		EtagInfos:         etagInfos,
	}

	cacheInfos = append(cacheInfos, cacheSyncInfo)

	cacheSync := structs.CacheSyncData{
		RequestNfType: "AUSF",
		CacheInfos:    cacheInfos,
	}

	cacheSyncData, err := json.Marshal(cacheSync)
	if err != nil {
		fmt.Printf("Marsh fail, err:%s\n", err.Error())
	}

	nfType := "AUSF"

	respDump := &httpclient.HttpRespData{}
	respDump.StatusCode = http.StatusOK
	respDump.Body = cacheSyncData

	err = handleDumpResponse(nfType, respDump)
	if err != nil {
		t.Fatalf("Expect handler dump response success, but fail, err:%s", err.Error())
	}
}

func TestGetTtlTimeStamp(t *testing.T) {
	nfInstanceID1 := "udm-5g-01"
	nfInstanceID2 := "udm-5g-02"

	timeStamp := time.Now()

	ttlInfos := make([]structs.TtlInfo, 0)
	ttlInfo1 := structs.TtlInfo{
		NfInstanceID: nfInstanceID1,
		ValidityTime: timeStamp,
	}
	ttlInfo2 := structs.TtlInfo{
		NfInstanceID: nfInstanceID2,
		ValidityTime: timeStamp,
	}
	ttlInfos = append(ttlInfos, ttlInfo1)
	ttlInfos = append(ttlInfos, ttlInfo2)

	ttlTime1 := getTtlTimeStamp(ttlInfos, nfInstanceID1)
	if ttlTime1 != timeStamp {
		t.Fatal("Expect get ttl timestamp success, but failure")
	}

	ttlTime2 := getTtlTimeStamp(ttlInfos, nfInstanceID2)
	if ttlTime2 != timeStamp {
		t.Fatal("Expect get ttl timestamp success, but failure")
	}

	probeTime := time.Time{}
	nfInstanceID3 := "no-exist"
	ttlTime3 := getTtlTimeStamp(ttlInfos, nfInstanceID3)
	if ttlTime3 != probeTime {
		t.Fatal("Expect get ttl timestamp failure, but not")
	}
}
