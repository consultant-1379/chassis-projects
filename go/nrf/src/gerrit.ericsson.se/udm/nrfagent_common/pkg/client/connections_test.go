package client

import (
	//"fmt"
	"testing"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
)

var dataNrfService = []byte(`
    {
          "mode": "active-standby",
          "profile": [
            {
			 "id": "nrf-server-0",
			 "ipv4-address": ["127.0.0.1","10.0.0.1","10.0.0.2"],
              "service": [
                {
                  "id": 0,
                  "scheme": "http",
                  "version": [
                    {
                      "api-version-in-uri": "v1"
                    }
                  ],
                  "fqdn": "seliius03696.seli.gic.ericsson.se",
                  "ip-endpoint": [
                    {
                      "id": 0,
                      "ipv4-address": "127.0.0.1",
                      "ipv6-address": "",
                      "port": 3211,
                      "transport": "TCP"
                    }
                  ],
                  "name": "nnrf-disc"
                },
                {
                  "id": 1,
                  "scheme": "http",
                  "version": [
                    {
                      "api-version-in-uri": "v1"
                    }
                  ],
                  "fqdn": "",
                  "ip-endpoint": [
                    {
                      "id": 0,
                      "ipv4-address": "127.0.0.1",
                      "ipv6-address": "",
                      "port": 3212,
                      "transport": "TCP"
                    },
                    {
                      "id": 1,
                      "ipv4-address": "127.0.0.1",
                      "ipv6-address": "",
                      "port": 3213,
                      "transport": "TCP"
                    }
                  ],  
                  "name": "nnrf-nfm"
                }
              ]
            }
          ]
    }
`)

var dataNrfServiceErr = []byte(`
    {
          "usage": "active_standby",
          "nrfServerProfileList": [
            {
              "nrfServerServiceEndPoints": [
			    {
                  "id": 0,
                  "scheme": "http",
                  "versions": [
                    {
                      "apiVersionInUri": "v1"
                    }
                  ],
                  "fqdn": "",
                  "ipEndPoints": [
                    {
                      "id": 0,
                      "ipv4Address": "",
                      "ipv6Address": "",
                      "port": 3211,
                      "transport": "TCP"
                    }
                  ],
                  "apiPrefix": "nnrf-disc",
                  "serviceName": "nnrf-disc"
                },
                {
                  "id": 1,
                  "scheme": "http",
                  "versions": [
                    {
                      "apiVersionInUri": "v1"
                    }
                  ],
                  "fqdn": "",
                  "ipEndPoints": [
                    {
                      "id": 0,
                      "ipv4Address": "127.0.0.1",
                      "ipv6Address": "",
                      "port": 3212,
                      "transport": "TCP"
                    },
                    {
                      "id": 1,
                      "ipv4Address": "127.0.0.1",
                      "ipv6Address": "",
                      "port": 3213,
                      "transport": "TCP"
                    }
                  ],
                  "apiPrefix": "nnrf-nfm",
                  "serviceName": "nnrf-nfm"
                }
              ]
            }
          ]
    }
`)

var (
	nrfDiscBaseURL1 = "http://127.0.0.1:3211/nnrf-disc/v1/"
	nrfDiscBaseURL2 = "http://seliius03696.seli.gic.ericsson.se/nnrf-disc/v1/"
	nrfMgmtBaseURL1 = "http://127.0.0.1:3212/nnrf-nfm/v1/"
	nrfMgmtBaseURL2 = "http://127.0.0.1:3213/nnrf-nfm/v1/"
)

var primaryMessage = []byte(`{"monitorRole":"primaryRole","mgmtUrlPrefix":"","discUrlPrefix":""}`)
var secondaryMessage = []byte(` {"monitorRole":"secondaryRole","mgmtUrlPrefix":"","discUrlPrefix":""}`)
var wrongMessage = []byte(`{"Role":"secondaryRole","mgmtUrlPrefix":"","discUrlPrefix":""}`)

var sendMmVerify monitorMessage

func sendMonitorMessageStub() {
	sendMonitorMessage = func(msg *monitorMessage) error {
		sendMmVerify.Role = msg.Role
		sendMmVerify.MgmtPrefix = msg.MgmtPrefix
		sendMmVerify.DiscPrefix = msg.DiscPrefix
		sendMmVerify.ForceReq = msg.ForceReq
		return nil
	}
}

func resetSendMmVerify() {
	sendMmVerify.Role = ""
	sendMmVerify.MgmtPrefix = ""
	sendMmVerify.DiscPrefix = ""
	sendMmVerify.ForceReq = false
}

func TestMonitorMessageHandler(t *testing.T) {
	backupSendMonitorMessage := sendMonitorMessage
	defer func() {
		sendMonitorMessage = backupSendMonitorMessage
	}()
	sendMonitorMessageStub()

	backupHb2NrfServer := hb2NrfServer
	defer func() {
		hb2NrfServer = backupHb2NrfServer
	}()
	hb2NrfServerStub()

	t.Run("TestPrimaryHandler", func(t *testing.T) {
		setNrfServerPrefix("", "")
		mm := &monitorMessage{
			Role: SeconaryMonitor,
		}
		primaryMonitorMessageHandler(mm)

		setNrfServerPrefix(nrfMgmtBaseURL2, nrfDiscBaseURL2)
		primaryMonitorMessageHandler(mm)
	})

	t.Run("TestSeconaryHandler", func(t *testing.T) {
		var mp, dp string
		setNrfServerPrefix("", "")
		mm := &monitorMessage{
			Role: PrimaryMonitor,
		}
		secondaryMonitorMessageHandler(mm)
		mp, dp = getNrfServerPrefix()
		nrfConnStatus := GetNRFConnStatus()
		if mp != "" ||
			dp != "" ||
			nrfConnStatus != NRFConnUnknown {
			t.Errorf("TestSeconaryHandler failed")
		}
		mm.MgmtPrefix = nrfMgmtBaseURL1
		mm.DiscPrefix = nrfDiscBaseURL1
		secondaryMonitorMessageHandler(mm)
		mp, dp = getNrfServerPrefix()
		nrfConnStatus = GetNRFConnStatus()
		if mp == "" ||
			dp == "" ||
			nrfConnStatus != NRFConnNormal {
			t.Errorf("TestSeconaryHandler failed")
		}

		mm.MgmtPrefix = ""
		setNrfServerPrefix("", "")
		secondaryMonitorMessageHandler(mm)
		mp, dp = getNrfServerPrefix()
		nrfConnStatus = GetNRFConnStatus()
		if mp != "" ||
			dp != "" ||
			nrfConnStatus != NRFConnLost {
			t.Errorf("TestSeconaryHandler failed")
		}
	})

	t.Run("TestMessageHandler", func(t *testing.T) {
		monitorRole = ""
		monitorMessageHandler(primaryMessage)

		monitorRole = ""
		monitorMessageHandler(wrongMessage)

		monitorRole = PrimaryMonitor
		monitorMessageHandler(secondaryMessage)

		monitorRole = SeconaryMonitor
		monitorMessageHandler(primaryMessage)
	})
}

func TestGetBaseNrfURLs(t *testing.T) {
	//test format value
	structs.UpdateNrfServerList(dataNrfService)

	var nrfServers structs.NrfServerList
	if !structs.GetNrfServerList(&nrfServers) {
		t.Fatalf("failed to get NrfServerList")
	}
	if len(nrfServers.NrfServerProfileList) == 0 {
		t.Fatalf("failed to get NrfServerProfileList")
	}
	if len(nrfServers.NrfServerProfileList[0].NrfServiceEndPoints) == 0 {
		t.Fatalf("failed to get NrfServiceEndPoints")

	}

	t.Run("TestGetBaseNrfURLs_FetchServiceLevelAddr", func(t *testing.T) {
		//fmt.Println("--->nrfServers:", nrfServers.NrfServerProfileList[0].NrfServiceEndPoints[0])
		baseURL := getBaseNrfURLs(&nrfServers.NrfServerProfileList[0].NrfServiceEndPoints[0], &nrfServers.NrfServerProfileList[0], false)
		//fmt.Println("--->baseURL:", baseURL)
		if len(baseURL) == 0 {
			t.Errorf("TestGetBaseNrfURLs: getBaseNrfURLs format check failure.")
		}

		if baseURL[0] != nrfDiscBaseURL1 {
			t.Errorf("TestGetBaseNrfURLs: getBaseNrfURLs check failure.")
		}

		baseURL = getBaseNrfURLs(&nrfServers.NrfServerProfileList[0].NrfServiceEndPoints[1], &nrfServers.NrfServerProfileList[0], false)
		//fmt.Println("--->baseURL:", baseURL)
		if len(baseURL) == 0 {
			t.Errorf("TestGetBaseNrfURLs: getBaseNrfURLs format check failure.")
		}

		if len(baseURL) != 2 {
			t.Errorf("TestGetBaseNrfURLs: getBaseNrfURLs check failure.")
		}

	})

	t.Run("TestGetBaseNrfURLs_FetchNFLevelAddr", func(t *testing.T) {
		//fmt.Println("--->nrfServers nnrf-disc:", nrfServers.NrfServerProfileList[0].NrfServiceEndPoints[0])
		baseURL := getBaseNrfURLs(&nrfServers.NrfServerProfileList[0].NrfServiceEndPoints[0], &nrfServers.NrfServerProfileList[0], true)
		//fmt.Println("--->baseURL:", baseURL)
		//		baseURL1 := getBaseNrfURLs(&nrfServers.NrfServerProfileList[0].NrfServiceEndPoints[1], &nrfServers.NrfServerProfileList[0], true)
		//		fmt.Println("--->baseURL1:", baseURL1)
		if len(baseURL) == 0 {
			t.Errorf("TestGetBaseNrfURLs: getBaseNrfURLs format check failure.")
		}

		if baseURL[0] != nrfDiscBaseURL1 {
			t.Errorf("TestGetBaseNrfURLs: getBaseNrfURLs check failure.")
		}
	})

}

func TestConnectNrfServer(t *testing.T) {
	backupHb2NrfServer := hb2NrfServer
	defer func() {
		hb2NrfServer = backupHb2NrfServer
	}()

	hb2NrfServer = func(url string, serviceName string) (*httpclient.HttpRespData, bool) {
		if url == nrfDiscBaseURL2 ||
			url == nrfMgmtBaseURL2 ||
			url == nrfDiscBaseURL1 ||
			url == nrfMgmtBaseURL1 {
			return nil, true
		}
		return nil, false
	}

	structs.UpdateNrfServerList(dataNrfService)

	var nrfServers structs.NrfServerList
	if !structs.GetNrfServerList(&nrfServers) {
		t.Fatalf("failed to get NrfServerList")
	}
	if len(nrfServers.NrfServerProfileList) == 0 {
		t.Fatalf("failed to get NrfServerProfileList")
	}
	if len(nrfServers.NrfServerProfileList[0].NrfServiceEndPoints) == 0 {
		t.Fatalf("failed to get NrfServiceEndPoints")

	}

	SyncNrfStatus := connectNrfServer(&nrfServers.NrfServerProfileList[0], false)
	if SyncNrfStatusNOK == SyncNrfStatus {
		t.Fatalf("Service Level connectNrfServer failed")
	}

	SyncNrfStatus = connectNrfServer(&nrfServers.NrfServerProfileList[0], true)
	if SyncNrfStatusNOK == SyncNrfStatus {
		t.Fatalf("NF Level connectNrfServer failed")
	}

}

func hb2NrfServerStub() {
	hb2NrfServer = func(url string, serviceName string) (*httpclient.HttpRespData, bool) {
		if url == nrfDiscBaseURL2 ||
			url == nrfMgmtBaseURL2 {
			return nil, true
		}
		return nil, false
	}
}

func TestSelectNrfServer(t *testing.T) {
	backupHb2NrfServer := hb2NrfServer
	defer func() {
		hb2NrfServer = backupHb2NrfServer
	}()
	hb2NrfServerStub()

	backupSendMonitorMessage := sendMonitorMessage
	defer func() {
		sendMonitorMessage = backupSendMonitorMessage
	}()
	sendMonitorMessageStub()

	setNrfServerPrefix("", "")
	monitorRole = PrimaryMonitor
	resetSendMmVerify()
	structs.UpdateNrfServerList(dataNrfService)
	selectNrfServer()

	var mp, dp string
	mp, dp = getNrfServerPrefix()
	//fmt.Println("--->mp:", mp, "\tdp:", dp)
	if dp != nrfDiscBaseURL2 ||
		mp != nrfMgmtBaseURL2 ||
		sendMmVerify.Role != PrimaryMonitor ||
		sendMmVerify.DiscPrefix != nrfDiscBaseURL2 ||
		sendMmVerify.MgmtPrefix != nrfMgmtBaseURL2 {
		t.Errorf("TestSelectNrfServer: selectNrfServer failure.")
	}
	resetSendMmVerify()

}

func TestSendResponseForForceReq(t *testing.T) {
	backupSendMonitorMessage := sendMonitorMessage
	defer func() {
		sendMonitorMessage = backupSendMonitorMessage
	}()
	sendMonitorMessageStub()

	monitorRole = PrimaryMonitor
	resetSendMmVerify()
	mm := &monitorMessage{
		Role:       monitorRole,
		MgmtPrefix: "",
		DiscPrefix: "",
		ForceReq:   true,
	}
	ret := sendResponseForForceReq(mm)
	if !ret || sendMmVerify.Role != PrimaryMonitor || sendMmVerify.ForceReq {
		t.Errorf("TestSendResponseForForceReq: ForceReq true monitorMessage check failure.")
	}
	resetSendMmVerify()

	mm = &monitorMessage{
		Role:       monitorRole,
		MgmtPrefix: "test.com",
		DiscPrefix: "test.com",
		ForceReq:   false,
	}
	ret = sendResponseForForceReq(mm)
	if ret {
		t.Errorf("TestSendResponseForForceReq: ForceReq false monitorMessage check failure.")
	}
}

func TestMonitorHandler(t *testing.T) {
	backupSendMonitorMessage := sendMonitorMessage
	defer func() {
		sendMonitorMessage = backupSendMonitorMessage
	}()
	sendMonitorMessageStub()

	backupHb2NrfServer := hb2NrfServer
	defer func() {
		hb2NrfServer = backupHb2NrfServer
	}()
	hb2NrfServerStub()

	setNrfServerPrefix("", "")
	monitorRole = SeconaryMonitor
	resetSendMmVerify()
	monitorHandler(true)
	if sendMmVerify.Role != SeconaryMonitor || !sendMmVerify.ForceReq {
		t.Errorf("TestSendResponseForForceReq: ForceReq true monitorMessage check failure.")
	}
	resetSendMmVerify()
	monitorHandler(false)
	if sendMmVerify.Role != SeconaryMonitor || sendMmVerify.ForceReq {
		t.Errorf("TestSendResponseForForceReq: ForceReq true monitorMessage check failure.")
	}
	monitorRole = PrimaryMonitor
	resetSendMmVerify()
	structs.UpdateNrfServerList(dataNrfService)
	monitorHandler(false)
	if sendMmVerify.Role != PrimaryMonitor ||
		sendMmVerify.DiscPrefix != nrfDiscBaseURL2 ||
		sendMmVerify.MgmtPrefix != nrfMgmtBaseURL2 {
		t.Errorf("TestSendResponseForForceReq: ForceReq true monitorMessage check failure.")
	}
	resetSendMmVerify()
}
