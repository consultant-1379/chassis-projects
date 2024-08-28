package configmap

import (
	"path/filepath"
	"os"
	"io/ioutil"
	"testing"
)

var dbInfoConfFile = "eric-nrf-dbinfo-conf.json"
func TestDBInfoConf(t *testing.T) {
	jsonContent := `{
		"locator-server-name": "eric-nrf-kvdb-ag-locator",
		"region-names": "ericsson-nrf-nrfaddresses",
		"locator-server-port": 10334
        }`

	var buffer = []byte(jsonContent)
	currDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	AbsFileName := currDir + "/" + dbInfoConfFile

	err := ioutil.WriteFile(AbsFileName, buffer, 0666)
	if err != nil {
		t.Fatal(err.Error())
	}

	dbInfoConf := &DBInfoConf{
		FileName:     "",
		ConfInstance: &SubDBInfoConf{},
	}

	dbInfoConf.SetFileName(AbsFileName)

	if AbsFileName != dbInfoConf.GetFileName() {
		t.Fatalf(`attributesConf.GetFileName() didn't return expected FileName`)
	}

	dbInfoConf.LoadConf()
	if DBLocatorServerName != "eric-nrf-kvdb-ag-locator" {
		t.Fatalf("DBLocatorServerName should have the correct value, but not ")
	}
	if DBNrfAddressRegionName != "ericsson-nrf-nrfaddresses" {
		t.Fatalf("DBNrfAddressRegionName should have the correct value, but not ")
	}
	if DBNfprofileRegionName != "ericsson-nrf-nfprofiles" {
		t.Fatalf("DBNfprofileRegionName should have the correct value, but not ")
	}
	if DBSubscriptionRegionName != "ericsson-nrf-subscriptions" {
		t.Fatalf("DBSubscriptionRegionName should have the correct value, but not ")
	}
	if DBGroupProfileRegionName != "ericsson-nrf-groupprofiles" {
		t.Fatalf("DBSubscriptionRegionName should have the correct value, but not ")
	}
	if DBImsiprefixProfileRegionName != "ericsson-nrf-imsiprefixprofiles" {
		t.Fatalf("DBImsiprefixProfileRegionName should have the correct value, but not ")
	}
	if DBNrfprofileRegionName != "ericsson-nrf-nrfprofiles" {
		t.Fatalf("DBNrfprofileRegionName should have the correct value, but not ")
	}
	if DBGpsiProfileRegionName != "ericsson-nrf-gpsiprofiles" {
		t.Fatalf("DBGpsiProfileRegionName should have the correct value, but not ")
	}
	if DBGpsiprefixProfileRegionName != "ericsson-nrf-gpsiprefixprofiles" {
		t.Fatalf("DBGpsiprefixProfileRegionName should have the correct value, but not ")
	}
	if DBCachenfprofileRegionName != "ericsson-nrf-cachenfprofiles" {
		t.Fatalf(" should have the correct value, but not ")
	}
	if DBLocatorServerPort != 10334 {
		t.Fatalf(" should have the correct value, but not ")
	}
}
