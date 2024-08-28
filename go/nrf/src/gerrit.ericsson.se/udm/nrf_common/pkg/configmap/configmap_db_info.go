package configmap

import (
	"encoding/json"
	"fmt"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"io/ioutil"
)

var (
	DBLocatorServerName = "eric-nrf-kvdb-ag-locator"
	DBLocatorServerPort = 10334

	DBNrfAddressRegionName        = "ericsson-nrf-nrfaddresses"
	DBNfprofileRegionName         = "ericsson-nrf-nfprofiles"
	DBSubscriptionRegionName      = "ericsson-nrf-subscriptions"
	DBGroupProfileRegionName      = "ericsson-nrf-groupprofiles"
	DBImsiprefixProfileRegionName = "ericsson-nrf-imsiprefixprofiles"
	DBNrfprofileRegionName        = "ericsson-nrf-nrfprofiles"
	DBGpsiProfileRegionName       = "ericsson-nrf-gpsiprofiles"
	DBGpsiprefixProfileRegionName = "ericsson-nrf-gpsiprefixprofiles"
	DBCachenfprofileRegionName    = "ericsson-nrf-cachenfprofiles"
)

type DBInfoConf struct {
	FileName     string
	ConfInstance *SubDBInfoConf
}

type SubDBInfoConf struct {
	LocatorServerName string `json:"locator-server-name"`
	LocatorServerPort int    `json:"locator-server-port"`
	RegionNames       string `json:"region-names"`
}

// SetFileName set FileName
func (conf *DBInfoConf) SetFileName(fileName string) {
	conf.FileName = fileName
}

// GetFileName return FileName
func (conf *DBInfoConf) GetFileName() string {
	return conf.FileName
}

// ParseConf set the value of global variable from struct
func (conf *DBInfoConf) ParseConf() {
	DBLocatorServerName = conf.ConfInstance.LocatorServerName
	DBLocatorServerPort = conf.ConfInstance.LocatorServerPort
}

// LoadConf set the value of global variable from configuration file
func (conf *DBInfoConf) LoadConf() {
	log.Debugf("Start to load db info config file %s", conf.FileName)
	bytes, err := ioutil.ReadFile(conf.FileName)
	if err != nil {
		log.Warnf("Failed to read db info config file:%s. %s", conf.FileName, err.Error())
		return
	}

	err = json.Unmarshal(bytes, conf.ConfInstance)
	if err != nil {
		log.Warnf("Failed to unmarshal db info config file:%s. %s", conf.FileName, err.Error())
		return
	}

	conf.ParseConf()
	log.Debugf("Successfully load db info config file %s", conf.FileName)
	conf.ShowConf()
}

// ReloadConf reset the value of variable when configuration file changed
func (conf *DBInfoConf) ReloadConf() {
	log.Debugf("reload db info config file")
	conf.LoadConf()
}

// ShowConf shows the value of global variable
func (conf *DBInfoConf) ShowConf() {
	fmt.Printf("locator server name: %s\n", DBLocatorServerName)
	fmt.Printf("locator server port: %d\n", DBLocatorServerPort)

	fmt.Printf("nrf address region name: %s\n", DBNrfAddressRegionName)
	fmt.Printf("nf profile region name: %s\n", DBNfprofileRegionName)
	fmt.Printf("subscript region name: %s\n", DBSubscriptionRegionName)
	fmt.Printf("group profile region name: %s\n", DBGroupProfileRegionName)
	fmt.Printf("imsiprefix profile region name: %s\n", DBImsiprefixProfileRegionName)
	fmt.Printf("nrf profile region name: %s\n", DBNrfprofileRegionName)
	fmt.Printf("gpsi profile region name: %s\n", DBGpsiProfileRegionName)
	fmt.Printf("gpsi prefix profile region name: %s\n", DBGpsiprefixProfileRegionName)
	fmt.Printf("cache nf profile region name: %s\n", DBCachenfprofileRegionName)
}
