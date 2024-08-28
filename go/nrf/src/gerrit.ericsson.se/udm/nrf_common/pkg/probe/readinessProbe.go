package probe

import (
	"com/dbproxy/nfmessage/nfprofile"
	"net/http"

	"os"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/fm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
	"com/dbproxy"
)

var (
	shutDownFlag = false
	isNeedWarmUp = true
	warmUpNum = 1000
	tryTime      = 1500
	//isConfigReady = false
	// CmConfigFlag is to notify  managment server that configuration in CM is ready for managment service
	//CmConfigFlag = make(chan bool)
)

func SetShutDownFlag(flag bool) {
	shutDownFlag = flag
}

// ReadinessProbe_Handler only for nrf_disc & nrf_prov
func ReadinessProbe_Handler(rw http.ResponseWriter, req *http.Request) {
	log.Debugf("readinessProbe comes")

	ready := isDBReady()
	workMode := os.Getenv("WORK_MODE")
	if workMode == constvalue.APP_WORKMODE_NRF_DISC {
		if isNeedWarmUp && ready {
			doWarmUp()
			isNeedWarmUp = false
		}
		handleDatabaseConnectionAlarm(ready)
	}

	if shutDownFlag || !ready {
		isNeedWarmUp = true
		rw.WriteHeader(http.StatusInternalServerError)
		log.Warningf("readinessProbe result %d", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	log.Debugf("readinessProbe result %d", http.StatusOK)
}

func isDBReady() bool {
	workMode := os.Getenv("WORK_MODE")
	if workMode == constvalue.APP_WORKMODE_NRF_DISC {
		return isNfProfileGetSuccess()
	}
	return isNfProfileGetSuccess() && isNfProfilePutSuccess() && isNfProfileDeleteSuccess()
}

func handleDatabaseConnectionAlarm(ready bool) {
        if ready {
                fm.ClearNRFDatabaseConnectionFailureAlarm()
        } else {
                fm.RaiseNRFDatabaseConnectionFailureAlarm()
        }
}

func isNfProfileGetSuccess() bool {
	dbReady := false
	nfInstanceID := os.Getenv("HOSTNAME")
	nfProfileKey := &nfprofile.NFProfileGetRequest_TargetNfInstanceId{
		TargetNfInstanceId: nfInstanceID,
	}

	getReq := &nfprofile.NFProfileGetRequest{
		Data: nfProfileKey,
	}

	getRes, err := dbmgmt.GetNFProfile(getReq)
	if err == nil && (getRes.GetCode() == dbmgmt.DbDataNotExist || getRes.GetCode() == dbmgmt.DbGetSuccess) {
		dbReady = true
	} else {
		log.Warning("isNfProfileGetSuccess failed")
	}

	return dbReady
}

func isNfProfilePutSuccess() bool {
	dbReady := false
	nfInstanceID := os.Getenv("HOSTNAME")
	body := `{
		"expiredTime": 12345678,
		"lastUpdateTime": 12345678,
		"provisioned": 0,
		"md5sum": {
			"nfProfile": "aaaaa",
			"nausf-auth-01": "bbbbb"
		},
		"body": {
  			"capacity": 100,
  			"fqdn": "seliius03696.seli.gic.ericsson.se",
  			"nfInstanceId": "` + nfInstanceID + `",
  			"nfServices": [
    			{
      				"allowedNfTypes": [
        				"AUSF",
        				"NRF",
        				"AMF"
      				],
      				"allowedPlmns": [
        			{
          				"mcc": "460",
          				"mnc": "000"
        			}
      				],
      				"capacity": 100,
      				"fqdn": "seliius03696.seli.gic.ericsson.se",
      				"ipEndPoints": [
        			{
          				"ipv4Address": "172.16.208.1",
          				"port": 30088
        			}
      				],
      				"schema": "http://",
      				"serviceInstanceId": "nausf-auth-01",
      				"serviceName": "nausf-auth",
      				"version": [
        			{
          				"apiFullVersion": "1.R15.1.1 ",
          				"apiVersionInUri": "v1",
          				"expiry": "2020-07-06T02:54:32Z"
        			}
      				]
    			}
  			],
  			"nfStatus": "SUSPENDED",
  			"nfType": "readinessProbeTest",
  			"plmn": {
    				"mcc": "460",
    				"mnc": "00"
  			},
  			"sNssais": [
    			{
      				"sd": "1",
      				"sst": 1
    			},
    			{
      				"sd": "0",
      				"sst": 0
    			}
  			],
  			"udrInfo": {
    				"supiRanges": [
      				{
        				"end": "20000",
        				"start": "10000"
      				}
    				]
  			}
		}
	}`

	putReq := &nfprofile.NFProfilePutRequest{NfInstanceId: nfInstanceID, NfProfile: body}
	putresp, err := dbmgmt.PutNFProfile(putReq)

	if err == nil && (putresp.GetCode() == dbmgmt.DbPutSuccess) {
		dbReady = true
	} else {
		log.Warning("isNfProfilePutSuccess failed")
	}
	return dbReady
}

func isNfProfileDeleteSuccess() bool {
	dbReady := false
	nfInstanceID := os.Getenv("HOSTNAME")
	deleteReq := &nfprofile.NFProfileDelRequest{NfInstanceId: nfInstanceID}
	deleteResp, err := dbmgmt.DeleteNFProfile(deleteReq)

	if err == nil && (deleteResp.GetCode() == dbmgmt.DbDataNotExist || deleteResp.GetCode() == dbmgmt.DbDeleteSuccess) {
		dbReady = true
	} else {
		log.Warning("isNfProfileDeleteSuccess failed")
	}
	return dbReady
}

//doWarmUp is used to warm up jvm to decrease the latency when pod first start
func doWarmUp() {
	log.Warn("start Warm up")
	var i int
	for i = 0; i < warmUpNum; i++ {
		isNfProfileQuerySuccess()
	}
	log.Warn("warm up success!")
}

//isNfProfileQuerySuccess is to monitor nfprofile query process
func isNfProfileQuerySuccess() bool {
	dbReady := false

	getReq := &dbproxy.QueryRequest{
		RegionName: "ericsson-nrf-nfprofiles",
		Query: []string{"SELECT DISTINCT value.nfInstanceId,value.profileUpdateTime FROM (SELECT DISTINCT value.helper FROM /ericsson-nrf-nfprofiles.entrySet, value.helper.smfInfo.taiList tai, value.helper.smfInfo.taiRangeList tai_range,  tai_range.tacRangeList tac_range WHERE (((((tai.plmnId.mcc = '310' AND tai.plmnId.mnc = '010') AND tai.tac = 'abc101') OR ((tai_range.plmnId.mcc = '310' AND tai_range.plmnId.mnc = '010') AND (((tac_range.start <= 'abc101' AND tac_range.start.length() <= 6) AND (tac_range.end >= 'abc101' AND tac_range.end.length() >= 6)) OR 'abc101'.matches(tac_range.pattern.toString()) = true))) OR (tai.tac = 'RESERVED_EMPTY_TAC' AND tac_range.pattern = 'RESERVED_EMPTY_TAC_RANGE_PATTERN')))) value WHERE (value.nfType = 'SMF')"},
	}

	queryRes, err := dbmgmt.QueryWithFilter(getReq)
	if err == nil && (queryRes.GetCode() == dbmgmt.DbDataNotExist || queryRes.GetCode() == dbmgmt.DbGetSuccess) {
		dbReady = true
	} else {
		log.Warning("isNfProfileQuerySuccess failed")
	}

	return dbReady
}