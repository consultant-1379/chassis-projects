package configmap

import (
	"fmt"
	"os"

	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/utils"
)

type ConfigMap interface {
	SetFileName(fileName string)
	GetFileName() string
	LoadConf()
	ReloadConf()
	ParseConf()
}

// InternalMgmtConfInst points to the internal config content for NRF Management
var InternalMgmtConfInst *internalconf.InternalMgmtConf

// InternalDiscConfInst points to the internal config content for NRF Discovery
var InternalDiscConfInst *internalconf.InternalDiscConf

// InternalProvConfInst points to the internal config content for NRF Provision
var InternalProvConfInst *internalconf.InternalProvConf

// InternalNotifyConfInst points to the internal config content for NRF Notification
var InternalNotifyConfInst *internalconf.InternalNotifyConf

var MultisiteInfoConfInst *SubMultisiteInfoConf

// AttributesConfInst is used for discovery search Attributes
var AttributesConfInst *SubAttributesConf

var DBInfoConfInst *SubDBInfoConf

var ConfigMapMap map[string]ConfigMap

func InitConfigMap() error {

	workMode := os.Getenv("WORK_MODE")

	ConfigMapMap = make(map[string]ConfigMap)

	confDir := os.Getenv("INTERNAL_CONF_DIR")

	if !utils.FileExist(confDir) {
		return fmt.Errorf("directory %s doesn't exist", confDir)
	}

	configFile := os.Getenv("INTERNAL_CONF_FILE")

	AbsFileName := confDir + "/" + configFile

	internalConfInst := &internalconf.InternalNrfConf{
		FileName: AbsFileName,
	}

	if workMode == constvalue.APP_WORKMODE_NRF_MGMT {
		InternalMgmtConfInst = &internalconf.InternalMgmtConf{}
		internalConfInst.ConfInstance = InternalMgmtConfInst

	} else if workMode == constvalue.APP_WORKMODE_NRF_DISC {
		InternalDiscConfInst = &internalconf.InternalDiscConf{}
		internalConfInst.ConfInstance = InternalDiscConfInst

	} else if workMode == constvalue.APP_WORKMODE_NRF_PROV {
		InternalProvConfInst = &internalconf.InternalProvConf{}
		internalConfInst.ConfInstance = InternalProvConfInst

	} else if workMode == constvalue.AppWorkmodeNrfNotif {
		InternalNotifyConfInst = &internalconf.InternalNotifyConf{}
		internalConfInst.ConfInstance = InternalNotifyConfInst
	}

	ConfigMapMap[AbsFileName] = internalConfInst

	confDir = os.Getenv("DB_MULTISITE_INFO_CONF_DIR")
	if !utils.FileExist(confDir) {
		return fmt.Errorf("directory %s doesn't exist", confDir)
	}

	configFile = os.Getenv("DB_MULTISITE_INFO_CONF_FILE")

	AbsFileName = confDir + "/" + configFile

	MultisiteInfoConfInst = &SubMultisiteInfoConf{}

	ConfigMapMap[AbsFileName] = &MultisiteInfoConf{
		FileName:     AbsFileName,
		ConfInstance: MultisiteInfoConfInst,
	}

	configFile = os.Getenv("DB_INFO_CONF")
	if !utils.FileExist(configFile) {
		return fmt.Errorf("directory %s doesn't exist", configFile)
	}

	AbsFileName = configFile

	DBInfoConfInst = &SubDBInfoConf{}

	ConfigMapMap[AbsFileName] = &DBInfoConf{
		FileName:     AbsFileName,
		ConfInstance: DBInfoConfInst,
	}

	if workMode == constvalue.APP_WORKMODE_NRF_DISC {
		configFile = os.Getenv("DBPROXY_ATTRIBUTES_CONF")
		if !utils.FileExist(configFile) {
			return fmt.Errorf("directory %s doesn't exist", configFile)
		}

		AbsFileName = configFile

		AttributesConfInst = &SubAttributesConf{}

		ConfigMapMap[AbsFileName] = &AttributesConf{
			FileName:     AbsFileName,
			ConfInstance: AttributesConfInst,
		}
	}

	for _, v := range ConfigMapMap {
		v.LoadConf()
	}

	return nil
}
