package configmap

import (
	"encoding/json"
	"io/ioutil"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

var (
	MultisiteEnabled       = false
	MultisiteHeartbeatTime = 3
	MultisiteMonitorTime   = 5
	MultisiteExpireTime    = 10
)

type MultisiteInfoConf struct {
	FileName     string
	ConfInstance *SubMultisiteInfoConf
}

type SubMultisiteInfoConf struct {
	Enabled       bool
	HeartbeatTime int
	MonitorTime   int
	ExpireTime    int
}

// SetFileName set FileName
func (conf *MultisiteInfoConf) SetFileName(fileName string) {
	conf.FileName = fileName
}

// GetFileName return FileName
func (conf *MultisiteInfoConf) GetFileName() string {
	return conf.FileName
}

// ParseConf set the value of global variable from struct
func (conf *MultisiteInfoConf) ParseConf() {
	MultisiteEnabled = conf.ConfInstance.Enabled
	MultisiteHeartbeatTime = conf.ConfInstance.HeartbeatTime
	MultisiteMonitorTime = conf.ConfInstance.MonitorTime
	MultisiteExpireTime = conf.ConfInstance.ExpireTime
}

// LoadConf set the value of global variable from configuration file
func (conf *MultisiteInfoConf) LoadConf() {
	log.Debugf("Start to load config file %s", conf.FileName)
	bytes, err := ioutil.ReadFile(conf.FileName)
	if err != nil {
		log.Warnf("Failed to read config file:%s. %s", conf.FileName, err.Error())
		return
	}

	err = json.Unmarshal(bytes, conf.ConfInstance)
	if err != nil {
		log.Warnf("Failed to unmarshal config file:%s. %s", conf.FileName, err.Error())
		return
	}

	conf.ParseConf()
	log.Debugf("Successfully load config file %s", conf.FileName)
}

// ReloadConf reset the value of variable when configuration file changed
func (conf *MultisiteInfoConf) ReloadConf() {
	log.Debugf("reload config file")
	conf.LoadConf()
}
