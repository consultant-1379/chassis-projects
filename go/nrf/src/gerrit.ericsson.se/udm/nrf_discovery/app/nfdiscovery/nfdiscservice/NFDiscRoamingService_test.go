package nfdiscservice

import (
	"testing"

	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
)

func TestGetAddrFromNrfAddressData(t *testing.T) {
	roaming := &NFDiscRoamingService{}
	data := []byte(`{
		"address":{
			"fqdn": "seliius03696.seli.gic.ericsson.se",
			"ipv4Address": "127.0.0.1",
			"ipv6Address": "192.168.0.1"
		}
	}`)
	flag, addr := roaming.getAddrFromNrfAddressData(data)
	if !flag || addr != "seliius03696.seli.gic.ericsson.se" {
		t.Fatal("func getAddrFromNrfAddressData() should return true and addr, but fail")
	}
	data2 := []byte(`{
		"address":{
			"ipv4Address": "127.0.0.1",
			"ipv6Address": "192.168.0.1"
		}
	}`)
	flag2, addr2 := roaming.getAddrFromNrfAddressData(data2)
	if !flag2 || addr2 != "127.0.0.1" {
		t.Fatal("func getAddrFromNrfAddressData() should return true and addr, but fail")
	}
	data3 := []byte(`{
		"address":{
			"ipv6Address": "192.168.0.1"
		}
	}`)
	flag3, addr3 := roaming.getAddrFromNrfAddressData(data3)
	if !flag3 || addr3 != "192.168.0.1" {
		t.Fatal("func getAddrFromNrfAddressData() should return true and addr, but fail")
	}
}

func TestGetRemotePLMN(t *testing.T) {
	cm.NfProfile.PlmnID = append(cm.NfProfile.PlmnID, cm.TPLMN{Mnc: "000", Mcc: "460"})
	cm.NfProfile.PlmnID = append(cm.NfProfile.PlmnID, cm.TPLMN{Mnc: "111", Mcc: "460"})

	var DiscPara nfdiscrequest.DiscGetPara
	value := make(map[string][]string)
	DiscPara.InitMember(value)

	var plmnList []string
	plmnList = append(plmnList, "{\"mcc\":\"460\", \"mnc\":\"00\"}")
	plmnList = append(plmnList, "{\"mcc\":\"460\", \"mnc\":\"11\"}")
	DiscPara.SetFlag(constvalue.SearchDataTargetPlmnList, true)
	DiscPara.SetValue(constvalue.SearchDataTargetPlmnList, plmnList)
	roamingService := &NFDiscRoamingService{}
	remotePlmnList := roamingService.getRemotePLMN(DiscPara)

	if remotePlmnList[0] != "460011" {
		t.Fatal("func getRemotePLMN should get remote plmn 460011, but not")
	}
}
