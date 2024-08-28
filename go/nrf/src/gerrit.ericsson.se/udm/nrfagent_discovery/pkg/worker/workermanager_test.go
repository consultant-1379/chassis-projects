package worker

import (
	"fmt"
	"testing"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/cache"
)

/*
func TestPrepareDiscoveryAgent(t *testing.T) {
	t.Log("Execute case TestPrepareDiscoveryAgent")
}

func TestSyncCacheSubscribeInfo(t *testing.T) {
	t.Log("Execute case TestSyncCacheSubscribeInfo")
}

func TestPrepareNfRegister(t *testing.T) {
	t.Log("Execute case TestPrepareNfRegister")
}
*/
func TestHaveBacklog(t *testing.T) {
	nfType := "AUSF"

	rest := workerManager.haveBacklog(nfType)
	t.Logf("NfType:%s have inject task backlog result : %v", nfType, rest)
	if !rest {
		t.Fatalf("Expect nfType:%s have deploy task backlog, but not", nfType)
	}

	nfType = "UDM"
	rest = workerManager.haveBacklog(nfType)
	t.Logf("NfType:%s have inject task backlog result : %v", nfType, rest)
	if rest {
		t.Fatalf("Expect nfType:%s do not deploy task backlog, but not", nfType)
	}
}

/*
func TestWaitAgentReady(t *testing.T) {
	t.Log("Execute case TestWaitAgentReady")
}
*/
func TestDumpReady(t *testing.T) {
	nfType := "AUSF"
	key1 := fmt.Sprintf("%s-%s", "UDM", "udm-servicer-1")
	key2 := fmt.Sprintf("%s-%s", "UDM", "udm-servicer-2")
	key3 := fmt.Sprintf("%s", "UDM")

	defer workerManager.resetSubscribeBacklogTask("AUSF")
	defer workerManager.resetFetchProfileBacklogTask("AUSF")

	t.Logf("subscribe : %v", workerManager.subscribeContainer)
	t.Logf("fetchProfile : %v", workerManager.fetchProfileContainer)

	rest := workerManager.DumpReady(nfType)
	t.Logf("Dump ready result : %v", rest)
	if rest {
		t.Fatalf("Expect nfType:%s do not ready, but not", nfType)
	}

	workerManager.setSubscribeTaskStatus(nfType, key1, SuccessStatus)
	rest = workerManager.DumpReady(nfType)
	t.Logf("Dump ready result : %v", rest)
	if rest {
		t.Fatalf("Expect nfType:%s do not ready, but not", nfType)
	}

	workerManager.setSubscribeTaskStatus(nfType, key2, FailureStatus)
	rest = workerManager.DumpReady(nfType)
	t.Logf("Dump ready result : %v", rest)
	if rest {
		t.Fatalf("Expect nfType:%s do not ready, but not", nfType)
	}

	workerManager.setFetchProfileTaskStatus(nfType, key3, FailureStatus)
	rest = workerManager.DumpReady(nfType)
	t.Logf("Dump ready result : %v", rest)
	if !rest {
		t.Fatalf("Expect nfType:%s is ready, but not", nfType)
	}
}

/*
func TestLaunchLeftTask(t *testing.T) {
	t.Log("Execute case TestLaunchLeftTask")
}
*/
func TestInjectSuccessSubscribeTask(t *testing.T) {
	nfType := "AUSF"
	key := fmt.Sprintf("%s-%s", "UDM", "udm-servicer-1")

	defer workerManager.resetSubscribeBacklogTask("AUSF")
	defer workerManager.resetFetchProfileBacklogTask("AUSF")

	status := workerManager.fetchSubscribeTaskStatus(nfType, key)
	t.Logf("subscribe task:%s status : %v", key, status)
	if status != InitStatus {
		t.Fatalf("Expect subscribe task:%s is InitStatus, but not", key)
	}

	workerManager.InjectSuccessSubscribeTask(nfType, key)

	status = workerManager.fetchSubscribeTaskStatus(nfType, key)
	t.Logf("subscribe task:%s status : %v", key, status)
	if status != SuccessStatus {
		t.Fatalf("Expect subscribe task:%s is SuccessStatus, but not", key)
	}
}

/*
func TestLaunchAllLeftTask(t *testing.T) {
	t.Log("Execute case TestLaunchAllLeftTask")
}

func TestStopSubscribeWorker(t *testing.T) {
	t.Log("Execute case TestStopSubscribeWorker")
}

func TestStopFetchProfileWorker(t *testing.T) {
	t.Log("Execute case TestStopFetchProfileWorker")
}

func TestStopDumpCacheWorker(t *testing.T) {
	t.Log("Execute case TestStopDumpCacheWorker")
}

func TestStopAllWorker(t *testing.T) {
	t.Log("Execute case TestStopAllWorker")
}
*/
func TestAgentReady(t *testing.T) {
	nfType := "AUSF"
	key1 := fmt.Sprintf("%s-%s", "UDM", "udm-servicer-1")
	key2 := fmt.Sprintf("%s-%s", "UDM", "udm-servicer-2")
	key3 := fmt.Sprintf("%s", "UDM")

	defer workerManager.resetSubscribeBacklogTask("AUSF")
	defer workerManager.resetFetchProfileBacklogTask("AUSF")

	t.Logf("subscribe : %v", workerManager.subscribeContainer)
	t.Logf("fetchProfile : %v", workerManager.fetchProfileContainer)

	rest := workerManager.agentReady()
	t.Logf("agent ready result : %v", rest)
	if rest {
		t.Fatal("Expect agentReady is false, but not")
	}

	workerManager.setSubscribeTaskStatus(nfType, key1, SuccessStatus)
	rest = workerManager.agentReady()
	t.Logf("agent ready result : %v", rest)
	if rest {
		t.Fatal("Expect agentReady is false, but not")
	}

	t.Logf("subscribe : %v", workerManager.subscribeContainer)
	t.Logf("fetchProfile : %v", workerManager.fetchProfileContainer)

	workerManager.setSubscribeTaskStatus(nfType, key2, SuccessStatus)
	rest = workerManager.agentReady()
	t.Logf("agent ready result : %v", rest)
	if rest {
		t.Fatal("Expect agentReady is false, but not")
	}

	t.Logf("subscribe : %v", workerManager.subscribeContainer)
	t.Logf("fetchProfile : %v", workerManager.fetchProfileContainer)

	workerManager.setFetchProfileTaskStatus(nfType, key3, SuccessStatus)
	rest = workerManager.agentReady()
	t.Logf("agent ready result : %v", rest)
	if !rest {
		t.Fatal("Expect agentReady is true, but not")
	}

	t.Logf("subscribe : %v", workerManager.subscribeContainer)
	t.Logf("fetchProfile : %v", workerManager.fetchProfileContainer)
}

/*
func TestLaunchSubscribeWorker(t *testing.T) {
	t.Log("Execute case TestLaunchSubscribeWorker")
}

func TestLaunchFetchProfileWorker(t *testing.T) {
	t.Log("Execute case TestLaunchFetchProfileWorker")
}

func TestLaunchDumpCacheWorker(t *testing.T) {
	t.Log("Execute case TestLaunchDumpCacheWorker")
}
*/
func TestInjectSubscribeBacklogTask(t *testing.T) {
	nfType := "AUSF"
	targetNfs, ok := cache.Instance().GetTargetNfs(nfType)
	if !ok {
		t.Fatalf("Get targetNf by nfType:%s failure", nfType)
	}

	defer workerManager.resetSubscribeBacklogTask("AUSF")
	defer workerManager.resetFetchProfileBacklogTask("AUSF")

	log.Infof("NfType:%s TargetNfs:%v", nfType, targetNfs)

	workerManager.injectSubscribeBacklogTask(nfType, targetNfs)

	key := fmt.Sprintf("%s-%s", "UDM", "udm-servicer-1")
	status := workerManager.fetchSubscribeTaskStatus(nfType, key)
	t.Logf("subscribe status : %v", status)
	if status != InitStatus {
		t.Fatalf("Expect subscribe task:%s stauts is InitStatus, but not", key)
	}
}

func TestInjectFetchProfileBacklogTask(t *testing.T) {
	nfType := "AUSF"
	targetNfs, ok := cache.Instance().GetTargetNfs(nfType)
	if !ok {
		log.Warnf("Get targetNfProfiles fail, don't deploy %s targetNf configmap", nfType)
		return
	}

	defer workerManager.resetSubscribeBacklogTask("AUSF")
	defer workerManager.resetFetchProfileBacklogTask("AUSF")

	log.Infof("NfType:%s TargetNfs:%v", nfType, targetNfs)

	workerManager.injectFetchProfileBacklogTask(nfType, targetNfs)
	key := fmt.Sprintf("%s", "UDM")
	status := workerManager.fetchFetchProfileTaskStatus(nfType, key)
	t.Logf("fetchProfile status : %v", status)
	if status != InitStatus {
		t.Fatalf("Expect fetchProfile task:%s stauts is InitStatus, but not", key)
	}
}

func TestResetSubscribeBacklogTask(t *testing.T) {
	nfType := "AUSF"
	key := fmt.Sprintf("%s-%s", "UDM", "udm-servicer-1")
	status := SuccessStatus
	workerManager.setSubscribeTaskStatus(nfType, key, status)

	defer workerManager.resetSubscribeBacklogTask("AUSF")
	defer workerManager.resetFetchProfileBacklogTask("AUSF")

	checkStatus := workerManager.checkSubscribeTaskSuccess(nfType, key)
	t.Logf("subscribe status : %v", checkStatus)
	if !checkStatus {
		t.Fatalf("Expect subscribe task:%s success, but not", key)
	}

	rest := workerManager.resetSubscribeBacklogTask(nfType)
	if !rest {
		t.Fatalf("Expect resetSubscribeBacklogTask:%s success, but not", key)
	}

	checkStatus = workerManager.checkSubscribeTaskSuccess(nfType, key)
	t.Logf("subscribe status : %v", checkStatus)
	if checkStatus {
		t.Fatalf("After reset subscribe backlog stask, Expect subscribe task:%s success failure, but not", key)
	}
}

func TestResetFetchProfileBacklogTask(t *testing.T) {
	nfType := "AUSF"
	key := fmt.Sprintf("%s", "UDM")
	status := SuccessStatus
	workerManager.setFetchProfileTaskStatus(nfType, key, status)

	defer workerManager.resetSubscribeBacklogTask("AUSF")
	defer workerManager.resetFetchProfileBacklogTask("AUSF")

	checkStatus := workerManager.checkFetchProfileTaskSuccess(nfType, key)
	t.Logf("fetchProfile status : %v", checkStatus)
	if !checkStatus {
		t.Fatalf("Expect fetchProfile task:%s success, but not", key)
	}

	rest := workerManager.resetFetchProfileBacklogTask(nfType)
	if !rest {
		t.Fatalf("Expect resetFetchProfileBacklogTask:%s success, but not", key)
	}

	checkStatus = workerManager.checkFetchProfileTaskSuccess(nfType, key)
	t.Logf("fetchProfile status : %v", checkStatus)
	if checkStatus {
		t.Fatalf("After reset fetchProfile backlog stask, Expect fetchProfile task:%s success failure, but not", key)
	}
}

func TestSetSubscribeTaskStatus(t *testing.T) {
	nfType := "AUSF"
	key := fmt.Sprintf("%s-%s", "UDM", "udm-servicer-1")
	status := SuccessStatus
	workerManager.setSubscribeTaskStatus(nfType, key, status)

	defer workerManager.resetSubscribeBacklogTask("AUSF")
	defer workerManager.resetFetchProfileBacklogTask("AUSF")

	checkStatus := workerManager.checkSubscribeTaskSuccess(nfType, key)
	t.Logf("subscribe status : %v", status)
	if !checkStatus {
		t.Fatalf("Expect subscribe task:%s success, but not", key)
	}
}

func TestSetFetchProfileTaskStatus(t *testing.T) {
	nfType := "AUSF"
	key := fmt.Sprintf("%s", "UDM")
	status := SuccessStatus
	workerManager.setFetchProfileTaskStatus(nfType, key, status)

	defer workerManager.resetSubscribeBacklogTask("AUSF")
	defer workerManager.resetFetchProfileBacklogTask("AUSF")

	checkStatus := workerManager.checkFetchProfileTaskSuccess(nfType, key)
	t.Logf("fetchProfile status : %v", status)
	if !checkStatus {
		t.Fatalf("Expect fetchProfile task:%s success, but not", key)
	}
}

func TestCheckSubscribeAllTaskSuccess(t *testing.T) {
	nfType := "AUSF"
	status := workerManager.checkSubscribeAllTaskSuccess(nfType)
	t.Logf("subscribe all task success status : %v", status)
	if status {
		t.Fatal("Expect subscribe all task success is not success, but not")
	}
}

func TestCheckFetchProfileAllTaskSuccess(t *testing.T) {
	nfType := "AUSF"
	status := workerManager.checkFetchProfileAllTaskSuccess(nfType)
	t.Logf("fetchProfile all task success status : %v", status)
	if status {
		t.Fatal("Expect fetchProfile all task success is not success, but not")
	}
}

func TestCheckSubscribeAllTaskDone(t *testing.T) {
	nfType := "AUSF"
	status := workerManager.checkSubscribeAllTaskDone(nfType)
	t.Logf("subscribe all task done status : %v", status)
	if status {
		t.Fatal("Expect subscribe all task done is not success, but not")
	}
}

func TestCheckFetchProfileAllTaskDone(t *testing.T) {
	nfType := "AUSF"
	status := workerManager.checkFetchProfileAllTaskDone(nfType)
	t.Logf("fetchProfile all task done status : %v", status)
	if status {
		t.Fatal("Expect fetchProfile all task done is not success, but not")
	}
}

func TestCheckSubscribeTaskSuccess(t *testing.T) {
	nfType := "AUSF"
	key := fmt.Sprintf("%s-%s", "UDM", "udm-servicer-1")
	status := workerManager.checkSubscribeTaskSuccess(nfType, key)
	t.Logf("subscribe status : %v", status)
	if status {
		t.Fatalf("Expect subscribe task:%s is not success, but not", key)
	}
}

func TestCheckFetchProfileTaskSuccess(t *testing.T) {
	nfType := "AUSF"
	key := fmt.Sprintf("%s", "UDM")
	status := workerManager.checkFetchProfileTaskSuccess(nfType, key)
	t.Logf("fetchProfile status : %v", status)
	if status {
		t.Fatalf("Expect fetchProfile task:%s is not success, but not", key)
	}
}

func TestFetchSubscribeTaskStatus(t *testing.T) {
	nfType := "AUSF"
	key := fmt.Sprintf("%s-%s", "UDM", "udm-servicer-1")
	status := workerManager.fetchSubscribeTaskStatus(nfType, key)
	t.Logf("subscribe status : %v", status)
	if status != InitStatus {
		t.Fatalf("Expect subscribe task:%s status is InitStatus, but not", key)
	}

	nfType = "No-such"
	status = workerManager.fetchSubscribeTaskStatus(nfType, key)
	t.Logf("subscribe status : %v", status)
	if status != UnknowStatus {
		t.Fatalf("Expect subscribe task:%s status is UnknowStatus, but not", key)
	}

	key = fmt.Sprintf("%s-%s", "UDM", "udm-servicer-no-such")
	status = workerManager.fetchSubscribeTaskStatus(nfType, key)
	t.Logf("subscribe status : %v", status)
	if status != UnknowStatus {
		t.Fatalf("Expect subscribe task:%s status is UnknowStatus, but not", key)
	}
}

func TestFetchFetchProfileTaskStatus(t *testing.T) {
	nfType := "AUSF"
	key := fmt.Sprintf("%s", "UDM")
	status := workerManager.fetchFetchProfileTaskStatus(nfType, key)
	t.Logf("fetchProfile status : %v", status)
	if status != InitStatus {
		t.Fatalf("Expect fetchProfile task:%s status is InitStatus, but not", key)
	}

	nfType = "No-such"
	status = workerManager.fetchFetchProfileTaskStatus(nfType, key)
	t.Logf("fetchProfile status : %v", status)
	if status != UnknowStatus {
		t.Fatalf("Expect fetchProfile task:%s status is UnknowStatus, but not", key)
	}

	nfType = "AUSF"
	key = fmt.Sprintf("%s", "no-exist")
	status = workerManager.fetchFetchProfileTaskStatus(nfType, key)
	t.Logf("fetchProfile status : %v", status)
	if status != UnknowStatus {
		t.Fatalf("Expect fetchProfile task:%s status is UnknowStatus, but not", key)
	}
}

func TestStartWorkModeMonitor(t *testing.T) {
	StartWorkModeMonitor(2)
	time.Sleep(2 * time.Second)
	if IsKeepCacheMode() {
		t.Fatalf("TestStartWorkModeMonitor keep cache status should be false.")
	}
}
