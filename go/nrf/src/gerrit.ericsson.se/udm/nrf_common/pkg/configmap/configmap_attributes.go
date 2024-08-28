package configmap

import (
	"encoding/json"
	"io/ioutil"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

var (
	//AttributesMap is to store Attributes.json into map
	AttributesMap map[string]map[string]SearchMapping
)

//AttributesConf is struct of Attribute
type AttributesConf struct {
	FileName     string
	ConfInstance *SubAttributesConf
}

//SubAttributesConf is format of Attributes.json
type SubAttributesConf struct {
	Common []SearchMapping `json:"common"`
	Udr    []SearchMapping `json:"udr"`
	Udm    []SearchMapping `json:"udm"`
	Ausf   []SearchMapping `json:"ausf"`
	Amf    []SearchMapping `json:"amf"`
	Smf    []SearchMapping `json:"smf"`
	Upf    []SearchMapping `json:"upf"`
	Pcf    []SearchMapping `json:"pcf"`
	Bsf    []SearchMapping `json:"bsf"`
	Chf    []SearchMapping `json:"chf"`
}

//SearchMapping is Attribute type
type SearchMapping struct {
	Parameter  string `json:"parameter"`
	Path       string `json:"path"`
	From       string `json:"from"`
	Where      string `json:"where"`
	ExistCheck bool   `json:"exist_check"`
}

// SetFileName set FileName
func (conf *AttributesConf) SetFileName(fileName string) {
	conf.FileName = fileName
}

// GetFileName return FileName
func (conf *AttributesConf) GetFileName() string {
	return conf.FileName
}

// ParseConf set the value of global variable from struct
func (conf *AttributesConf) ParseConf() {
	AttributesMap = make(map[string]map[string]SearchMapping)
	commonMap := make(map[string]SearchMapping)
	for _, value := range (conf.ConfInstance.Common) {
		commonMap[value.Parameter] = value
	}
	AttributesMap[constvalue.Common] = commonMap

	udrMap := make(map[string]SearchMapping)
	for _, value := range (conf.ConfInstance.Udr) {
		udrMap[value.Parameter] = value
	}
	AttributesMap[constvalue.NfTypeUDR] = udrMap

	udmMap := make(map[string]SearchMapping)
	for _, value := range (conf.ConfInstance.Udm) {
		udmMap[value.Parameter] = value
	}
	AttributesMap[constvalue.NfTypeUDM] = udmMap

	ausfMap := make(map[string]SearchMapping)
	for _, value := range (conf.ConfInstance.Ausf) {
		ausfMap[value.Parameter] = value
	}
	AttributesMap[constvalue.NfTypeAUSF] = ausfMap

	amfMap := make(map[string]SearchMapping)
	for _, value := range (conf.ConfInstance.Amf) {
		amfMap[value.Parameter] = value
	}
	AttributesMap[constvalue.NfTypeAMF] = amfMap

	smfMap := make(map[string]SearchMapping)
	for _, value := range (conf.ConfInstance.Smf) {
		smfMap[value.Parameter] = value
	}
	AttributesMap[constvalue.NfTypeSMF] = smfMap

	upfMap := make(map[string]SearchMapping)
	for _, value := range (conf.ConfInstance.Upf) {
		upfMap[value.Parameter] = value
	}
	AttributesMap[constvalue.NfTypeUPF] = upfMap

	pcfMap := make(map[string]SearchMapping)
	for _, value := range (conf.ConfInstance.Pcf) {
		pcfMap[value.Parameter] = value
	}
	AttributesMap[constvalue.NfTypePCF] = pcfMap

	bsfMap := make(map[string]SearchMapping)
	for _, value := range (conf.ConfInstance.Bsf) {
		bsfMap[value.Parameter] = value
	}
	AttributesMap[constvalue.NfTypeBSF] = bsfMap

	chfMap := make(map[string]SearchMapping)
	for _, value := range (conf.ConfInstance.Chf) {
		chfMap[value.Parameter] = value
	}
	AttributesMap[constvalue.NfTypeCHF] = chfMap
}

// LoadConf set the value of global variable from configuration file
func (conf *AttributesConf) LoadConf() {
	log.Debugf("Start to load attributes config file %s", conf.FileName)
	bytes, err := ioutil.ReadFile(conf.FileName)
	if err != nil {
		log.Warnf("Failed to read attributes config file:%s. %s", conf.FileName, err.Error())
		return
	}

	err = json.Unmarshal(bytes, conf.ConfInstance)
	if err != nil {
		log.Warnf("Failed to unmarshal attributes config file:%s. %s", conf.FileName, err.Error())
		return
	}

	conf.ParseConf()
	log.Debugf("Successfully load attributes config file %s", conf.FileName)
}

// ReloadConf reset the value of variable when configuration file changed
func (conf *AttributesConf) ReloadConf() {
	log.Debugf("reload attributes config file")
	conf.LoadConf()
}
