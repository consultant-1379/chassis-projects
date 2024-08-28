package common

import (
	"testing"
	"time"

	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
)

////GetSelfUUID test case
//func TestGetSelfUUID(t *testing.T) {
//	fmt.Println("Test GetSelfUUID")
//}

////SetSelfUUID test case
//func TestSetSelfUUID(t *testing.T) {
//	fmt.Println("Test SetSelfUUID")
//}

////GetDiscMsgbus test case
//func TestGetDiscMsgbus(t *testing.T) {
//	fmt.Println("Test GetDiscMsgbus")
//}

////SetDiscMsgbus test case
//func TestSetDiscMsgbus(t *testing.T) {
//	fmt.Println("Test SetDiscMsgbus")
//}

////CmNrfAgentLogHandler test case
//func TestCmNrfAgentLogHandler(t *testing.T) {
//	fmt.Println("Test CmNrfAgentLogHandler")
//}

//TestConvertIpv6ToIpv4
func TestConvertIpv6ToIpv4(t *testing.T) {
	v6Addr1 := "fe80:1234:0000:0000:0000:0000:0000:0001"
	v6Addr2 := "fe80:1234::0001"
	v6Addr3 := "fe80:1234::192.168.1.1"
	v6Addr4 := "fe80:1234::192.168.1"
	v6Addr5 := "fe80:1234:0000:0000:xxoo:0000:0000:0001"

	v4Addr := ""
	t.Run("case1", func(t *testing.T) {
		v4Addr = ConvertIpv6ToIpv4(v6Addr1)
		if v4Addr != "0.0.0.1" {
			t.Fatalf("failed to convert Ipv6Address %s", v6Addr1)
		}
	})
	t.Run("case2", func(t *testing.T) {
		v4Addr = ConvertIpv6ToIpv4(v6Addr2)
		if v4Addr != "0.0.0.1" {
			t.Fatalf("failed to convert Ipv6Address %s", v6Addr2)
		}
	})
	t.Run("case3", func(t *testing.T) {
		v4Addr = ConvertIpv6ToIpv4(v6Addr3)
		if v4Addr != "192.168.1.1" {
			t.Fatalf("failed to convert Ipv6Address %s", v6Addr3)
		}
	})
	t.Run("case4", func(t *testing.T) {
		v4Addr = ConvertIpv6ToIpv4(v6Addr4)
		if v4Addr != "" {
			t.Fatalf("failed to convert Ipv6Address %s", v6Addr4)
		}
	})
	t.Run("case5", func(t *testing.T) {
		v4Addr = ConvertIpv6ToIpv4(v6Addr5)
		if v4Addr != "" {
			t.Fatalf("failed to convert Ipv6Address %s", v6Addr5)
		}
	})
}

func TestConvertIpv6ToIpv4InSearchResult(t *testing.T) {
	content0 := []byte(`{"validityPeriod":86400,"nfInstances":{"nfInstanceId":"","nfType":"","nfStatus":"REGISERED","ipv4Addresses":[]}}`)
	content1 := []byte(`{"validityPeriod":86400,"nfInstances":[{"nfInstanceId":"","nfType":"","nfStatus":"REGISERED","ipv4Addresses":["127.0.0.1","127.0.0.2","127.0.0.3"]}]}`)
	content2 := []byte(`{"validityPeriod":86400,"nfInstances":[{"nfInstanceId":"","nfType":"","nfStatus":"REGISERED","ipv4Addresses":["127.0.0.1","127.0.0.2","127.0.0.3"],"ipv6Addresses":["fe80:1234:0000:0000:xxoo:0000:0000:0001","fe80:1234:0000:0000:0000:0000:0000:0002","fe80:1234:0000:0000:0000:0000:0000:0003"]}]}`)
	cooked2 := []byte(`{"validityPeriod":86400,"nfInstances":[{"nfInstanceId":"","nfType":"","nfStatus":"REGISERED","ipv4Addresses":["127.0.0.1","127.0.0.2","127.0.0.3","0.0.0.2","0.0.0.3"]}]}`)
	content3 := []byte(`{"validityPeriod":86400,"nfInstances":[{"ipv4Addresses":["127.0.0.1","127.0.0.2","127.0.0.3"],"ipv6Addresses":["fe80:1234:0000:0000:xxoo:0000:0000:0001","fe80:1234:0000:0000:0000:0000:0000:0002","fe80:1234:0000:0000:0000:0000:0000:0003"],"nfServices":[{"ipEndPoints":[{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0007","transport":"TCP","port":30088},{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0008","transport":"TCP","port":30089},{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0009","transport":"TCP","port":30090},{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0010","transport":"TCP","port":30089},{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0020","transport":"TCP","port":30089},{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0030","transport":"TCP","port":30089}]}]}]}`)

	var cooked []byte
	var err error
	t.Run("case1", func(t *testing.T) {
		cooked, err = ConvertIpv6ToIpv4InSearchResult(content0, false)
		s1 := string(cooked)
		s2 := string(content0)
		if err != nil || s1 != s2 {
			t.Fatalf("failed to disable convert")
		}
	})
	t.Run("case1", func(t *testing.T) {
		cooked, err = ConvertIpv6ToIpv4InSearchResult(content0, true)
		if err == nil {
			t.Fatalf("failed to convert content0")
		}
	})
	t.Run("case2", func(t *testing.T) {
		cooked, err = ConvertIpv6ToIpv4InSearchResult(content1, true)
		s1 := string(cooked)
		s2 := string(content1)
		if err != nil || s1 != s2 {
			t.Fatalf("failed to convert %s to %s", s1, s2)
		}
	})
	t.Run("case3", func(t *testing.T) {
		cooked, err = ConvertIpv6ToIpv4InSearchResult(content2, true)
		s1 := string(cooked)
		s2 := string(cooked2)
		if err != nil || s1 != s2 {
			t.Fatalf("failed to convert %s to %s", s1, s2)
		}
	})
	t.Run("case4", func(t *testing.T) {
		cooked, err = ConvertIpv6ToIpv4InSearchResult(content3, true)
		if err != nil {
			t.Fatalf("failed to convert content3")
		}
	})
}

func TestConvertIpv6ToIpv4InSearchResult01(t *testing.T) {
	//	content1 := []byte(`{"validityPeriod":86400,"nfInstances":[{"nfInstanceId":"","nfType":"","nfStatus":"REGISERED","ipv4Address":["127.0.0.1","127.0.0.2","127.0.0.3"]}]}`)
	//	content2 := []byte(`{"validityPeriod":86400,"nfInstances":[{"nfInstanceId":"","nfType":"","nfStatus":"REGISERED","ipv4Address":["127.0.0.1","127.0.0.2","127.0.0.3"],"ipv6Address":["fe80:1234:0000:0000:xxoo:0000:0000:0001","fe80:1234:0000:0000:0000:0000:0000:0002","fe80:1234:0000:0000:0000:0000:0000:0003"]}]}`)
	//	cooked2 := []byte(`{"validityPeriod":86400,"nfInstances":[{"nfInstanceId":"","nfType":"","nfStatus":"REGISERED","ipv4Address":["127.0.0.1","127.0.0.2","127.0.0.3","0.0.0.2","0.0.0.3"]}]}`)
	//	content3 := []byte(`{"validityPeriod":86400,"nfInstances":[{"ipv4Address":["127.0.0.1","127.0.0.2","127.0.0.3"],"ipv6Address":["fe80:1234:0000:0000:xxoo:0000:0000:0001","fe80:1234:0000:0000:0000:0000:0000:0002","fe80:1234:0000:0000:0000:0000:0000:0003"],"nfServices":[{"ipEndPoints":[{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0007","transport":"TCP","port":30088},{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0008","transport":"TCP","port":30089},{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0009","transport":"TCP","port":30090},{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0010","transport":"TCP","port":30089},{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0020","transport":"TCP","port":30089},{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0030","transport":"TCP","port":30089}]}]}]}`)

	content1 := []byte(`{"validityPeriod":86400,"nfInstances":[{"nfInstanceId":"","nfType":"","nfStatus":"REGISERED","ipv4Addresses":["127.0.0.1","127.0.0.2","127.0.0.3"]}]}`)
	content2 := []byte(`{"validityPeriod":86400,"nfInstances":[{"nfInstanceId":"","nfType":"","nfStatus":"REGISERED","ipv4Addresses":["127.0.0.1","127.0.0.2","127.0.0.3"],"ipv6Addresses":["fe80:1234:0000:0000:xxoo:0000:0000:0001","fe80:1234:0000:0000:0000:0000:0000:0002","fe80:1234:0000:0000:0000:0000:0000:0003"]}]}`)
	cooked2 := []byte(`{"validityPeriod":86400,"nfInstances":[{"nfInstanceId":"","nfType":"","nfStatus":"REGISERED","ipv4Addresses":["127.0.0.1","127.0.0.2","127.0.0.3","0.0.0.2","0.0.0.3"]}]}`)
	content3 := []byte(`{"validityPeriod":86400,"nfInstances":[{"ipv4Address":["127.0.0.1","127.0.0.2","127.0.0.3"],"ipv6Addresses":["fe80:1234:0000:0000:xxoo:0000:0000:0001","fe80:1234:0000:0000:0000:0000:0000:0002","fe80:1234:0000:0000:0000:0000:0000:0003"],"nfServices":[{"ipEndPoints":[{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0007","transport":"TCP","port":30088},{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0008","transport":"TCP","port":30089},{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0009","transport":"TCP","port":30090},{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0010","transport":"TCP","port":30089},{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0020","transport":"TCP","port":30089},{"ipv4Address":"","ipv6Address":"fe80:1234:0000:0000:0000:0000:0000:0030","transport":"TCP","port":30089}]}]}]}`)

	var cooked []byte
	var err error
	t.Run("case2", func(t *testing.T) {
		cooked, err = ConvertIpv6ToIpv4InSearchResult(content1, true)
		s1 := string(cooked)
		s2 := string(content1)
		if err != nil || s1 != s2 {
			t.Fatalf("failed to convert %s to %s", s1, s2)
		}
	})
	t.Run("case3", func(t *testing.T) {
		cooked, err = ConvertIpv6ToIpv4InSearchResult(content2, true)
		s1 := string(cooked)
		s2 := string(cooked2)
		if err != nil || s1 != s2 {
			t.Fatalf("failed to convert %s to %s", s1, s2)
		}
	})
	t.Run("case4", func(t *testing.T) {
		cooked, err = ConvertIpv6ToIpv4InSearchResult(content3, true)
		if err != nil {
			t.Fatalf("failed to convert content3")
		}
	})
}

func TestDispatchSubscrInfoToMessageBus(t *testing.T) {
	subInfo := structs.SubscriptionInfo{
		RequesterNfType:   "AUSF",
		TargetNfType:      "UDM",
		TargetServiceName: "udm-01",
		SubscriptionID:    "123-456-789",
		ValidityTime:      time.Time{},
	}

	ok := DispatchSubscrInfoToMessageBus(subInfo)
	if ok {
		t.Errorf("TestDispatchSubscrInfoToMessageBus: check failure")
	}
}
