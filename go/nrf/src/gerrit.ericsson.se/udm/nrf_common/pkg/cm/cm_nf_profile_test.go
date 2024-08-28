package cm

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTNfProfileParseConf(t *testing.T) {
	data := []byte(`{"instance-id": "0c765084-9cc5-49c6-9876-ae2f5fa2a63f",
        "type": "nrf",
        "status": "registered",
        "requested-heartbeat-timer": 120,
        "plmn-id": [],
        "snssai": [],
        "nsi": [],
        "fqdn": "",
        "inter-plmn-fqdn": "",
        "ipv4-address": [],
        "ipv6-address": [],
        "allowed-nf-domain": [],
        "allowed-plmn": [],
        "allowed-nf-type": [],
        "allowed-nssai": [],
        "priority": 10,
        "capacity": 100,
        "locality": "",
        "service-persistence": false,
        "service": [
          {
            "instance-id": "nnrf-nfm-01",
            "name": "nnrf-nfm",
            "version": [
              {
                "api-full-version": "1.R15.1.1",
                "api-version-in-uri": "v1",
                "expiry": "2020-07-06T02:54:32Z"
              }
            ],
            "scheme": "",
            "status": "registered",
            "fqdn": "",
            "inter-plmn-fqdn": "",
            "ip-endpoint": [
              {
                "id": 1,
                "transport": "tcp",
                "ipv4-address": "192.168.1.5",
                "port": 0
              }
            ],
            "api-prefix": "",
            "allowed-plmn": [],
            "allowed-nf-type": ["amf", "ausf"],
            "allowed-nf-domain": [],
            "allowed-nssai": [],
            "priority": 5,
            "capacity": 100,
            "supported-features": ""
          },
          {
            "instance-id": "nnrf-disc-01",
            "name": "nnrf-disc",
            "version": [
              {
                "api-full-version": "1.R15.1.1",
                "api-version-in-uri": "v1",
                "expiry": "2020-07-06T02:54:32Z"
              }
            ],
            "scheme": "",
            "status": "registered",
            "fqdn": "",
            "inter-plmn-fqdn": "",
            "ip-endpoint": [
              {
                "id": 1,
                "transport": "tcp",
                "ipv4-address": "192.168.1.6",
                "port": 0
              }
            ],
            "api-prefix": "",
            "allowed-plmn": [],
            "allowed-nf-type": ["pcf", "udm"],
            "allowed-nf-domain": [],
            "allowed-nssai": [],
            "priority": 5,
            "capacity": 100,
            "supported-features": ""
          }
        ]
	}`)

	nfProfile := &TNfProfile{}
	err := json.Unmarshal(data, nfProfile)
	if err != nil {
		t.Fatalf("Unmarshal error, %v", err)
	}

	nfProfile.ParseConf()

	if len(ManagementNfServices) != 1 || len(DiscoveryNfServices) != 1 {
		t.Fatalf("TNfProfile.ParseConf didn't work correctly !")
	}
}

func TestTNfProfileToUpper(t *testing.T) {
	nfProfile := TNfProfile{
		InstanceID:    "nrf01",
		Type:          "nrf",
		Status:        "registered",
		AllowedNfType: []string{"ausf", "amf"},
		Service: []TNfService{
			TNfService{
				InstanceID:    "nnrf-nfm-01",
				Name:          "nnrf-nfm",
				Status:        "registered",
				AllowedNfType: []string{"pcf", "udm"},
			},
			TNfService{
				InstanceID:                      "nnrf-nfm-02",
				Name:                            "nnrf-nfm",
				Status:                          "registered",
				AllowedNfType:                   []string{"pcf", "udm"},
				IPEndpoint:                      []TIPEndpoint{},
				DefaultNotificationSubscription: []TDefaultNotificationSubscription{},
			},
			TNfService{
				InstanceID:    "nnrf-disc-01",
				Name:          "nnrf-disc",
				Status:        "registered",
				AllowedNfType: []string{"pcf", "udm"},
				IPEndpoint: []TIPEndpoint{
					TIPEndpoint{
						ID:          1,
						Transport:   "tcp",
						Ipv4Address: "10.10.10.10",
						Port:        3000,
					},
				},
				DefaultNotificationSubscription: []TDefaultNotificationSubscription{
					TDefaultNotificationSubscription{
						NotificationType:   "n1-messages",
						CallbackURI:        "test",
						N1MessageClass:     "5gmm",
						N2InformationClass: "pws-bcal",
					},
				},
			},
			TNfService{
				InstanceID:    "nnrf-disc-02",
				Name:          "nnrf-disc",
				Status:        "registered",
				AllowedNfType: []string{"pcf", "udm"},
				IPEndpoint: []TIPEndpoint{
					TIPEndpoint{
						ID:          1,
						Transport:   "tcp",
						Ipv4Address: "10.10.10.10",
						Port:        3000,
					},
					TIPEndpoint{
						ID:          2,
						Transport:   "tcp",
						Ipv4Address: "10.10.10.20",
						Port:        3000,
					},
				},
				DefaultNotificationSubscription: []TDefaultNotificationSubscription{
					TDefaultNotificationSubscription{
						NotificationType:   "n1-messages",
						CallbackURI:        "test",
						N1MessageClass:     "5gmm",
						N2InformationClass: "pws-bcal",
					},
					TDefaultNotificationSubscription{
						NotificationType:   "n2-information",
						CallbackURI:        "test",
						N1MessageClass:     "updp",
						N2InformationClass: "nrppa",
					},
				},
			},
		},
	}

	originalNfProfile := nfProfile

	nfProfile.toUpper()

	if nfProfile.Type != "NRF" || nfProfile.Status != "REGISTERED" {
		t.Fatalf("TNfProfile.toUpper didn't return right value !")
	}

	for index := range nfProfile.AllowedNfType {
		if nfProfile.AllowedNfType[index] != strings.ToUpper(originalNfProfile.AllowedNfType[index]) {
			t.Fatalf("TNfProfile.toUpper didn't return right value !")
		}
	}

	for index := range nfProfile.Service {
		if nfProfile.Service[index].Status != strings.ToUpper(originalNfProfile.Service[index].Status) {
			t.Fatalf("TNfProfile.toUpper didn't return right value !")
		}

		for subIndex := range nfProfile.Service[index].IPEndpoint {
			if nfProfile.Service[index].IPEndpoint[subIndex].Transport != strings.ToUpper(originalNfProfile.Service[index].IPEndpoint[subIndex].Transport) {
				t.Fatalf("TNfProfile.toUpper didn't return right value !")
			}
		}

		for subIndex := range nfProfile.Service[index].AllowedNfType {
			if nfProfile.Service[index].AllowedNfType[subIndex] != strings.ToUpper(originalNfProfile.Service[index].AllowedNfType[subIndex]) {
				t.Fatalf("TNfProfile.toUpper didn't return right value !")
			}
		}

		for subIndex := range nfProfile.Service[index].DefaultNotificationSubscription {
			if nfProfile.Service[index].DefaultNotificationSubscription[subIndex].NotificationType != strings.ToUpper(originalNfProfile.Service[index].DefaultNotificationSubscription[subIndex].NotificationType) {
				t.Fatalf("TNfProfile.toUpper didn't return right value !")
			}
			if nfProfile.Service[index].DefaultNotificationSubscription[subIndex].N1MessageClass != strings.ToUpper(originalNfProfile.Service[index].DefaultNotificationSubscription[subIndex].N1MessageClass) {
				t.Fatalf("TNfProfile.toUpper didn't return right value !")
			}
			if nfProfile.Service[index].DefaultNotificationSubscription[subIndex].N2InformationClass != strings.ToUpper(originalNfProfile.Service[index].DefaultNotificationSubscription[subIndex].N2InformationClass) {
				t.Fatalf("TNfProfile.toUpper didn't return right value !")
			}
		}
	}

}
