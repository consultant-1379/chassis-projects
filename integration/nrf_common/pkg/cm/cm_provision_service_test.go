package cm

import (
	"encoding/json"
	"testing"

	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

func TestParseConf(t *testing.T) {

}

func TestConstructIngressAddress(t *testing.T) {
	// case 1: work in IPv4 stack, only IPv6 address configured, ingress shall be empty
	IPStackMode = constvalue.IPStackv4
	provSvcPro := []byte(`{
        "prov-address": [
          {
            "fqdn": "",
            "id": 1,
            "ipv6-address": "fd09::1320:a7a3:6f88:98de",
            "ipv4-address": "",
            "scheme": "https",
            "port": 3000
          }
        ]
      }
	`)
	var provService TProvisionService
	err := json.Unmarshal(provSvcPro, &provService)
	if err != nil {
		t.Errorf("TestConstructIngressAddress: json unmarshal error %s", err)
	}
	provService.ParseConf()

	if GetProvIngressAddress() != "" {
		t.Errorf("TestConstructIngressAddress: locationPrefix is error")
	}

	// case 2: work in IPv4 stack, only fqdn address configured, ingress shall be fqnd
	IPStackMode = constvalue.IPStackv4
	provSvcPro = []byte(`{
        "prov-address": [
          {
            "fqdn": "www.example.com",
            "id": 1,
            "ipv6-address": "",
            "ipv4-address": "",
            "scheme": "https",
            "port": 3000
          }
        ]
      }
	`)
	var provService1 TProvisionService
	err = json.Unmarshal(provSvcPro, &provService1)
	if err != nil {
		t.Errorf("TestConstructIngressAddress: json unmarshal error %s", err)
	}
	provService1.ParseConf()

	if GetProvIngressAddress() != "https://www.example.com:3000" {
		t.Errorf("TestConstructIngressAddress: locationPrefix is error")
	}

	// case 3: work in IPv4 stack, only IPv4 address configured, ingress shall be IPv4
	IPStackMode = constvalue.IPStackv4
	provSvcPro = []byte(`{
        "prov-address": [
          {
            "fqdn": "",
            "id": 1,
            "ipv6-address": "",
            "ipv4-address": "10.10.10.10",
            "scheme": "https",
            "port": 3000
          }
        ]
      }
	`)
	var provService3 TProvisionService
	err = json.Unmarshal(provSvcPro, &provService3)
	if err != nil {
		t.Errorf("TestConstructIngressAddress: json unmarshal error %s", err)
	}

	provService3.ParseConf()

	if GetProvIngressAddress() != "https://10.10.10.10:3000" {
		t.Errorf("TestConstructIngressAddress: locationPrefix is error")
	}

	// case 4: work in IPv4 stack, both IPv4 and IPv6 and fqdn address configured, ingress shall be IPv4
	IPStackMode = constvalue.IPStackv4
	provSvcPro = []byte(`{
        "prov-address": [
          {
            "fqdn": "www.example.com",
            "id": 1,
            "ipv6-address": "fd09::1320:a7a3:6f88:98de",
            "ipv4-address": "10.10.10.10",
            "scheme": "https",
            "port": 3000
          }
        ]
      }
	`)
	var provService4 TProvisionService
	err = json.Unmarshal(provSvcPro, &provService4)
	if err != nil {
		t.Errorf("TestConstructIngressAddress: json unmarshal error %s", err)
	}
	provService4.ParseConf()

	if GetProvIngressAddress() != "https://10.10.10.10:3000" {
		t.Errorf("TestConstructIngressAddress: locationPrefix is error")
	}

	// case 5: work in IPv4 stack, two items are configured,
	// one is IPv4, the other is IPv6, ingress shall be IPv4
	IPStackMode = constvalue.IPStackv4
	provSvcPro = []byte(`{
        "prov-address": [
          {
            "fqdn": "",
            "id": 1,
            "ipv6-address": "fd09::1320:a7a3:6f88:98de",
            "ipv4-address": "",
            "scheme": "https",
            "port": 3000
          },
          {
            "fqdn": "",
            "id": 1,
            "ipv6-address": "",
            "ipv4-address": "10.10.10.10",
            "scheme": "http",
            "port": 3000
          }
        ]
      }
	`)
	var provService5 TProvisionService
	err = json.Unmarshal(provSvcPro, &provService5)
	if err != nil {
		t.Errorf("TestConstructIngressAddress: json unmarshal error %s", err)
	}

	provService5.ParseConf()

	if GetProvIngressAddress() != "http://10.10.10.10:3000" {
		t.Errorf("TestConstructIngressAddress: locationPrefix is error")
	}

	// case 6: work in IPv6 stack, only IPv4 address configured, ingress shall be empty
	IPStackMode = constvalue.IPStackv6
	provSvcPro = []byte(`{
        "prov-address": [
          {
            "fqdn": "",
            "id": 1,
            "ipv6-address": "",
            "ipv4-address": "10.10.10.10",
            "scheme": "https",
            "port": 3000
          }
        ]
      }
	`)
	var provService6 TProvisionService
	err = json.Unmarshal(provSvcPro, &provService6)
	if err != nil {
		t.Errorf("TestConstructIngressAddress: json unmarshal error %s", err)
	}

	provService6.ParseConf()

	if GetProvIngressAddress() != "" {
		t.Errorf("TestConstructIngressAddress: locationPrefix is error")
	}

	// case 7: work in IPv6 stack, only fqdn address configured, ingress shall be fqnd
	IPStackMode = constvalue.IPStackv6
	provSvcPro = []byte(`{
        "prov-address": [
          {
            "fqdn": "www.example.com",
            "id": 1,
            "ipv6-address": "",
            "ipv4-address": "",
            "scheme": "https",
            "port": 3000
          }
        ]
      }
	`)
	var provService7 TProvisionService
	err = json.Unmarshal(provSvcPro, &provService7)
	if err != nil {
		t.Errorf("TestConstructIngressAddress: json unmarshal error %s", err)
	}

	provService7.ParseConf()

	if GetProvIngressAddress() != "https://www.example.com:3000" {
		t.Errorf("TestConstructIngressAddress: locationPrefix is error")
	}

	// case 8: work in IPv6 stack, only IPv6 address configured, ingress shall be IPv6
	IPStackMode = constvalue.IPStackv6
	provSvcPro = []byte(`{
        "prov-address": [
          {
            "fqdn": "",
            "id": 1,
            "ipv6-address": "fd09::1320:a7a3:6f88:98de",
            "ipv4-address": "",
            "scheme": "https",
            "port": 3000
          }
        ]
      }
	`)
	var provService8 TProvisionService
	err = json.Unmarshal(provSvcPro, &provService8)
	if err != nil {
		t.Errorf("TestConstructIngressAddress: json unmarshal error %s", err)
	}

	provService8.ParseConf()

	if GetProvIngressAddress() != "https://[fd09::1320:a7a3:6f88:98de]:3000" {
		t.Errorf("TestConstructIngressAddress: locationPrefix is error")
	}

	// case 9: work in IPv6 stack, both IPv4 and IPv6 and fqdn address configured, ingress shall be IPv6
	IPStackMode = constvalue.IPStackv6
	provSvcPro = []byte(`{
        "prov-address": [
          {
            "fqdn": "www.example.com",
            "id": 1,
            "ipv6-address": "fd09::1320:a7a3:6f88:98de",
            "ipv4-address": "10.10.10.10",
            "scheme": "https",
            "port": 3000
          }
        ]
      }
	`)
	var provService9 TProvisionService
	err = json.Unmarshal(provSvcPro, &provService9)
	if err != nil {
		t.Errorf("TestConstructIngressAddress: json unmarshal error %s", err)
	}

	provService9.ParseConf()

	if GetProvIngressAddress() != "https://[fd09::1320:a7a3:6f88:98de]:3000" {
		t.Errorf("TestConstructIngressAddress: locationPrefix is error")
	}

	// case 10: work in IPv6 stack, two items are configured,
	// one is IPv4, the other is IPv6, ingress shall be IPv6
	IPStackMode = constvalue.IPStackv6
	provSvcPro = []byte(`{
        "prov-address": [
          {
            "fqdn": "",
            "id": 1,
            "ipv6-address": "",
            "ipv4-address": "10.10.10.10",
            "scheme": "https",
            "port": 3000
          },
          {
            "fqdn": "",
            "id": 1,
            "ipv6-address": "fd09::1320:a7a3:6f88:98de",
            "ipv4-address": "",
            "scheme": "http",
            "port": 3000
          }
        ]
      }
	`)
	var provService10 TProvisionService
	err = json.Unmarshal(provSvcPro, &provService10)
	if err != nil {
		t.Errorf("TestConstructIngressAddress: json unmarshal error %s", err)
	}

	provService10.ParseConf()

	if GetProvIngressAddress() != "http://[fd09::1320:a7a3:6f88:98de]:3000" {
		t.Errorf("TestConstructIngressAddress: locationPrefix is error")
	}
}
