package nfdiscfilter

import (
	"github.com/buger/jsonparser"
	"testing"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

func TestEliminatServices(t *testing.T) {
	filter := &NFServiceFilter{}
	nfProfile := []byte(`[
		{
		"serviceName": "srv1",
		"nfServiceStatus":"REGISTERED"
		},
		{
		"serviceName": "srv2",
		"nfServiceStatus":"REGISTERED"
		},
		{
		"serviceName": "srv3",
		"nfServiceStatus":"REGISTERED"
		}
		]`)

	var DiscPara nfdiscrequest.DiscGetPara
	value := make(map[string][]string)
	DiscPara.InitMember(value)

	var serviceNameList []string
	serviceNameList = append(serviceNameList, "srv1")
	serviceNameList = append(serviceNameList, "srv3")
	serviceNameList = append(serviceNameList, "srv5")
	DiscPara.SetValue(constvalue.SearchDataServiceName, serviceNameList)
	DiscPara.SetFlag(constvalue.SearchDataServiceName, true)

	filter.setFilterResult(nfProfile, true)
	result := filter.eliminatServices(&DiscPara)
	if !result {
		t.Fatal("Should not return error, but did !")
	}

	var serviceNameArray2 []string

	_, err := jsonparser.ArrayEach(filter.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		serviceName, errTemp1 := jsonparser.GetString(value, constvalue.NFServiceName)
		if errTemp1 != nil {
			t.Fatal("Should not return error, but did !")
		}

		serviceNameArray2 = append(serviceNameArray2, serviceName)
	})

	if err != nil {
		t.Fatal("Should not return error, but did !")
	}

	if len(serviceNameArray2) != 2 {
		t.Fatal("There should be two services, but not !")
	}

	if serviceNameArray2[0] == serviceNameArray2[1] {
		t.Fatal("serviceName1 should not be the same as serviceName2")
	}

	for _, item := range serviceNameArray2 {
		if item != "srv1" && item != "srv3" {
			t.Fatal("serviceName should be srv1 or srv3, but not !")
		}

	}
}

func TestFilterNfServicesByRequesterNfType(t *testing.T) {
	filter := &NFServiceFilter{}
	filterInfo := &FilterInfo{nfProfilesMd5Sum:make(map[string]string, 1)}
	nfProfile := []byte(`[
						    {
							    "serviceName": "srv1",
								"serviceName": "srv1",
                                 "version": [],
                                 "schema": "schema1",
						    },
							{
						        "serviceInstanceId": "srv1",
                                 "serviceName": "srv1",
                                 "version": [],
                                 "schema": "schema1",
							    "allowedPlmns": [
                                     {
                                         "mcc": "000",
                                         "mnc": "00"
                                     },
								    {
                                         "mcc": "001",
                                         "mnc": "01"
                                     }
                                 ],
                                 "allowedNfTypes": [
                                     "NRF",
                                     "UDM"
                                 ],
					        },
							{
						        "serviceInstanceId": "srv1",
                                 "serviceName": "srv1",
                                 "version": [],
                                 "schema": "schema1",
							    "allowedPlmns": [
                                     {
                                         "mcc": "000",
                                         "mnc": "00"
                                     },
								    {
                                         "mcc": "001",
                                         "mnc": "01"
                                     }
                                 ],
                                 "allowedNfTypes": [
                                     "NRF",
                                     "AUSF"
                                 ],
					        }
						]`)
	var DiscPara nfdiscrequest.DiscGetPara
	value := make(map[string][]string)
	DiscPara.InitMember(value)

	var requesterNfTypeArray []string
	requesterNfTypeArray = append(requesterNfTypeArray, "NRF")
	DiscPara.SetFlag(constvalue.SearchDataRequesterNfType, true)
	DiscPara.SetValue(constvalue.SearchDataRequesterNfType, requesterNfTypeArray)

	filter.setFilterResult(nfProfile, true)
	retCode := filter.filterNfServicesByRequesterNfType(&DiscPara, filterInfo)
	if !retCode {
		t.Fatalf("should return 0, but not !")
	}

	serviceCount := 0

	_, _ = jsonparser.ArrayEach(filter.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		serviceCount++
	})

	if serviceCount != 3 {
		t.Fatalf("There should be 3 allowed service, but not !")
	}

	var requesterNfTypeArray2 []string
	requesterNfTypeArray2 = append(requesterNfTypeArray2, "UDM")
	DiscPara.SetFlag(constvalue.SearchDataRequesterNfType, true)
	DiscPara.SetValue(constvalue.SearchDataRequesterNfType, requesterNfTypeArray2)
	filter.setFilterResult(nfProfile, true)
	retCode = filter.filterNfServicesByRequesterNfType(&DiscPara, filterInfo)
	if !retCode {
		t.Fatalf("should return 0, but not !")
	}

	serviceCount = 0

	_, _ = jsonparser.ArrayEach(filter.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		serviceCount++
	})

	if serviceCount != 2 {
		t.Fatalf("There should be 2 allowed service, but not !")
	}

	var requesterNfTypeArray3 []string
	requesterNfTypeArray3 = append(requesterNfTypeArray3, "AUSF")
	DiscPara.SetFlag(constvalue.SearchDataRequesterNfType, true)
	DiscPara.SetValue(constvalue.SearchDataRequesterNfType, requesterNfTypeArray3)
	filter.setFilterResult(nfProfile, true)
	retCode = filter.filterNfServicesByRequesterNfType(&DiscPara, filterInfo)
	if !retCode {
		t.Fatalf("should return 0, but not !")
	}

	serviceCount = 0

	_, _ = jsonparser.ArrayEach(filter.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		serviceCount++
	})

	if serviceCount != 2 {
		t.Fatalf("There should be 2 allowed service, but not !")
	}

	var requesterNfTypeArray4 []string
	requesterNfTypeArray4 = append(requesterNfTypeArray4, "UDR")
	DiscPara.SetFlag(constvalue.SearchDataRequesterNfType, true)
	DiscPara.SetValue(constvalue.SearchDataRequesterNfType, requesterNfTypeArray4)
	filter.setFilterResult(nfProfile, true)
	retCode = filter.filterNfServicesByRequesterNfType(&DiscPara, filterInfo)
	if !retCode {
		t.Fatalf("should return 0, but not !")
	}

	serviceCount = 0

	_, _ = jsonparser.ArrayEach(filter.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		serviceCount++
	})

	if serviceCount != 1 {
		t.Fatalf("There should be 1 allowed service, but not !")
	}

	nfProfile = []byte(`[
							{
						        "serviceInstanceId": "srv1",
                                 "serviceName": "srv1",
                                 "version": [],
                                 "schema": "schema1",
							    "allowedPlmns": [
                                     {
                                         "mcc": "000",
                                         "mnc": "00"
                                     },
								    {
                                         "mcc": "001",
                                         "mnc": "01"
                                     }
                                 ],
                                 "allowedNfTypes": [
                                     "NRF",
                                     "UDM"
                                 ],
					        },
							{
						        "serviceInstanceId": "srv1",
                                 "serviceName": "srv1",
                                 "version": [],
                                 "schema": "schema1",
							    "allowedPlmns": [
                                     {
                                         "mcc": "000",
                                         "mnc": "00"
                                     },
								    {
                                         "mcc": "001",
                                         "mnc": "01"
                                     }
                                 ],
                                 "allowedNfTypes": [
                                     "NRF",
                                     "AUSF"
                                 ],
					        }
						]`)

	var requesterNfTypeArray5 []string
	requesterNfTypeArray5 = append(requesterNfTypeArray5, "UDR")
	DiscPara.SetFlag(constvalue.SearchDataRequesterNfType, true)
	DiscPara.SetValue(constvalue.SearchDataRequesterNfType, requesterNfTypeArray5)
	filter.setFilterResult(nfProfile, true)
	retCode = filter.filterNfServicesByRequesterNfType(&DiscPara, filterInfo)
	if !retCode {
		t.Fatalf("should return 1, but not !")
	}
	nfServiceCount := 0
	_, err := jsonparser.ArrayEach(filter.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		nfServiceCount++
	})
	if err != nil {
		t.Fatalf("nfService field should have, but not")
	}
	if nfServiceCount > 0 {
		t.Fatalf("nfService should be [], but have value")
	}
}

func TestFilterNfServicesByRequesterPLMN(t *testing.T) {
	filter := &NFServiceFilter{}
	filterInfo := &FilterInfo{nfProfilesMd5Sum:make(map[string]string, 1)}
	nfProfile := []byte(`[
							    {
								    "serviceName": "srv1",
									"serviceName": "srv1",
	                                 "version": [],
	                                 "schema": "schema1",
							    },
								{
							        "serviceInstanceId": "srv1",
	                                 "serviceName": "srv1",
	                                 "version": [],
	                                 "schema": "schema1",
								    "allowedPlmns": [
	                                     {
	                                         "mcc": "000",
	                                         "mnc": "00"
	                                     },
									    {
	                                         "mcc": "001",
	                                         "mnc": "01"
	                                     }
	                                 ],
	                                 "allowedNfTypes": [
	                                     "NRF",
	                                     "UDM"
	                                 ],
						        },
								{
							        "serviceInstanceId": "srv1",
	                                 "serviceName": "srv1",
	                                 "version": [],
	                                 "schema": "schema1",
								    "allowedPlmns": [
	                                     {
	                                         "mcc": "000",
	                                         "mnc": "00"
	                                     },
									    {
	                                         "mcc": "002",
	                                         "mnc": "02"
	                                     }
	                                 ],
	                                 "allowedNfTypes": [
	                                     "NRF",
	                                     "AUSF"
	                                 ],
						        }
							]`)
	var DiscPara nfdiscrequest.DiscGetPara
	value := make(map[string][]string)
	DiscPara.InitMember(value)

	var requesterPlmnArray []string
	requesterPlmnArray = append(requesterPlmnArray, "{\"mcc\": \"000\", \"mnc\": \"00\"}")
	DiscPara.SetFlag(constvalue.SearchDataRequesterPlmnList, true)
	DiscPara.SetValue(constvalue.SearchDataRequesterPlmnList, requesterPlmnArray)

	filter.setFilterResult(nfProfile, true)
	retCode := filter.filterNfServicesByRequesterPLMN(&DiscPara, filterInfo)
	if !retCode {
		t.Fatalf("should return 0, but not !")
	}

	serviceCount := 0

	_, _ = jsonparser.ArrayEach(filter.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		serviceCount++
	})

	if serviceCount != 3 {
		t.Fatalf("There should be 3 allowed service, but not !")
	}

	var requesterPlmnArray2 []string
	requesterPlmnArray2 = append(requesterPlmnArray2, "{\"mcc\": \"001\", \"mnc\": \"01\"}")
	DiscPara.SetFlag(constvalue.SearchDataRequesterPlmnList, true)
	DiscPara.SetValue(constvalue.SearchDataRequesterPlmnList, requesterPlmnArray2)

	filter.setFilterResult(nfProfile, true)
	retCode = filter.filterNfServicesByRequesterPLMN(&DiscPara, filterInfo)
	if !retCode {
		t.Fatalf("should return 0, but not !")
	}

	serviceCount = 0

	_, _ = jsonparser.ArrayEach(filter.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		serviceCount++
	})

	if serviceCount != 2 {
		t.Fatalf("There should be 2 allowed service, but not !")
	}

	var requesterPlmnArray3 []string
	requesterPlmnArray3 = append(requesterPlmnArray3, "{\"mcc\": \"002\", \"mnc\": \"02\"}")
	DiscPara.SetFlag(constvalue.SearchDataRequesterPlmnList, true)
	DiscPara.SetValue(constvalue.SearchDataRequesterPlmnList, requesterPlmnArray3)

	filter.setFilterResult(nfProfile, true)
	retCode = filter.filterNfServicesByRequesterPLMN(&DiscPara, filterInfo)
	if !retCode {
		t.Fatalf("should return 0, but not !")
	}

	serviceCount = 0

	_, _ = jsonparser.ArrayEach(filter.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		serviceCount++
	},)

	if serviceCount != 2 {
		t.Fatalf("There should be 2 allowed service, but not !")
	}

	var requesterPlmnArray4 []string
	requesterPlmnArray4 = append(requesterPlmnArray4, "{\"mcc\": \"003\", \"mnc\": \"03\"}")
	DiscPara.SetFlag(constvalue.SearchDataRequesterPlmnList, true)
	DiscPara.SetValue(constvalue.SearchDataRequesterPlmnList, requesterPlmnArray4)

	filter.setFilterResult(nfProfile, true)
	retCode = filter.filterNfServicesByRequesterPLMN(&DiscPara, filterInfo)
	if !retCode {
		t.Fatalf("should return 0, but not !")
	}

	serviceCount = 0

	_, _ = jsonparser.ArrayEach(filter.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		serviceCount++
	})

	if serviceCount != 1 {
		t.Fatalf("There should be 1 allowed service, but not !")
	}

	nfProfile = []byte(`[
								{
							        "serviceInstanceId": "srv1",
	                                 "serviceName": "srv1",
	                                 "version": [],
	                                 "schema": "schema1",
								    "allowedPlmns": [
	                                     {
	                                         "mcc": "000",
	                                         "mnc": "00"
	                                     },
									    {
	                                         "mcc": "001",
	                                         "mnc": "01"
	                                     }
	                                 ],
	                                 "allowedNfTypes": [
	                                     "NRF",
	                                     "UDM"
	                                 ],
						        },
								{
							        "serviceInstanceId": "srv1",
	                                 "serviceName": "srv1",
	                                 "version": [],
	                                 "schema": "schema1",
								    "allowedPlmns": [
	                                     {
	                                         "mcc": "000",
	                                         "mnc": "00"
	                                     },
									    {
	                                         "mcc": "001",
	                                         "mnc": "01"
	                                     }
	                                 ],
	                                 "allowedNfTypes": [
	                                     "NRF",
	                                     "AUSF"
	                                 ],
						        }
							]`)

	var requesterPlmnArray5 []string
	requesterPlmnArray5 = append(requesterPlmnArray5, "{\"mcc\": \"003\", \"mnc\": \"03\"}")
	DiscPara.SetFlag(constvalue.SearchDataRequesterPlmnList, true)
	DiscPara.SetValue(constvalue.SearchDataRequesterPlmnList, requesterPlmnArray3)

	filter.setFilterResult(nfProfile, true)
	retCode = filter.filterNfServicesByRequesterPLMN(&DiscPara, filterInfo)
	if !retCode {
		t.Fatalf("should return 1, but not !")
	}
	nfServiceCount := 0
	_, err := jsonparser.ArrayEach(filter.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		nfServiceCount++
	})
	if err != nil {
		t.Fatalf("nfService field should have, but not")
	}
	if nfServiceCount > 0 {
		t.Fatalf("nfService should be [], but have value")
	}
}

func TestFilterNfServicesByRequesterNfFQDN(t *testing.T) {
	filter := &NFServiceFilter{}
	filterInfo := &FilterInfo{nfProfilesMd5Sum:make(map[string]string, 1)}
	nfProfile := []byte(`[
						    {
						      "serviceInstanceId": "nupf-test-01",
						      "serviceName": "nupf-test02",
						      "version": [],
						      "schema": "https",
						      "fqdn": "seliius03695.seli.gic.ericsson.se",
						      "allowedNfTypes": [
							"UDM"
						      ],
						      "allowedNfDomains":[".+ericsson.se"]
						    },
						     {
						      "serviceInstanceId": "nupf-test-02",
						      "serviceName": "nupf-test02",
						      "version": [],
						      "schema": "https",
						      "fqdn": "seliius03696.seli.gic.ericsson.se",
						      "allowedNfTypes": [
							"UDM"
						      ],
						      "allowedNfDomains":["^seliius\\d{5}.seli.gic.ericsson.se$", "^seliius\\d{4}.seli.gic.ericsson.se$"]
						    },
						     {
						      "serviceInstanceId": "nupf-test-03",
						      "serviceName": "nupf-test02",
						      "version": [],
						      "schema": "https",
						      "fqdn": "seliius03697.seli.gic.ericsson.se",
						      "allowedNfTypes": [
							"UDM"
						      ],
						      "allowedNfDomains":["seliius03669\\.seli.gic.ericsson.se", "seliius03696.seli.gic.ericsson.se"]
						    }
						  ]`)
	var DiscPara nfdiscrequest.DiscGetPara
	value := make(map[string][]string)
	DiscPara.InitMember(value)

	var requesterNfFQDNArray []string
	requesterNfFQDNArray = append(requesterNfFQDNArray, "seliius03696.seli.gic.ericsson.se")
	DiscPara.SetFlag(constvalue.SearchDataRequesterNFInstFQDN, true)
	DiscPara.SetValue(constvalue.SearchDataRequesterNFInstFQDN, requesterNfFQDNArray)

	filter.setFilterResult(nfProfile, true)
	retCode := filter.filterNfServicesByRequesterNfFQDN(&DiscPara, filterInfo)
	if !retCode {
		t.Fatalf("should return 0, but not !")
	}

	serviceCount := 0

	_, _ = jsonparser.ArrayEach(filter.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		serviceCount++
	})
	if serviceCount != 3 {
		t.Fatalf("There should be 3 allowed service, but not !")
	}

	var requesterNfFQDNArray2 []string
	requesterNfFQDNArray2 = append(requesterNfFQDNArray2, "seliius03697.seli.gic.ericsson.se")
	DiscPara.SetFlag(constvalue.SearchDataRequesterNFInstFQDN, true)
	DiscPara.SetValue(constvalue.SearchDataRequesterNFInstFQDN, requesterNfFQDNArray2)

	filter.setFilterResult(nfProfile, true)
	retCode = filter.filterNfServicesByRequesterNfFQDN(&DiscPara, filterInfo)
	if !retCode {
		t.Fatalf("should return 0, but not !")
	}

	serviceCount = 0

	_, _ = jsonparser.ArrayEach(filter.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		serviceCount++
	})

	if serviceCount != 2 {
		t.Fatalf("There should be 2 allowed service, but not !")
	}

	var requesterNfFQDNArray3 []string
	requesterNfFQDNArray3 = append(requesterNfFQDNArray3, "swqeqw.seli.gic.ericsson.se")
	DiscPara.SetFlag(constvalue.SearchDataRequesterNFInstFQDN, true)
	DiscPara.SetValue(constvalue.SearchDataRequesterNFInstFQDN, requesterNfFQDNArray3)

	filter.setFilterResult(nfProfile, true)
	retCode = filter.filterNfServicesByRequesterNfFQDN(&DiscPara, filterInfo)
	if !retCode {
		t.Fatalf("should return 0, but not !")
	}

	serviceCount = 0

	_, _ = jsonparser.ArrayEach(filter.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		serviceCount++
	})

	if serviceCount != 1 {
		t.Fatalf("There should be 1 allowed service, but not !")
	}

}

