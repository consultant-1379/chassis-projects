package cm

import (
	//"bytes"
	"encoding/json"
	//"io/ioutil"
	"os"
	"path/filepath"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/utils"
)

// FsNotifyOpt for fsnotify option
type FsNotifyOpt struct {
	fileDir  string
	fileName string
}

// Configuration for internal configuration
type Configuration struct {
	EnableIpv6             bool
	EnableIpv6ToIpv4       bool
	DefaultValidityPeriod  int
	DefaultRequesterNfList []string
	LoadThreshold          int
	RegWaitTime            int
	HTTP2Conns             int
	HTTP2Timeout           int
	KeepCacheRetryCount    int
}

var configFile string
var instance *Configuration
var fsWatcher *utils.FsNotify

// GetInstance for Configuration instance
func GetInstance() *Configuration {
	if instance == nil {
		instance = &Configuration{
			EnableIpv6:             false,
			EnableIpv6ToIpv4:       false,
			DefaultValidityPeriod:  86400,
			DefaultRequesterNfList: []string{"AUSF", "UDM", "UDR", "NEF", "PCF", "AMF", "SMF", "SMSF", "NSSF", "UPF", "LMF", "GMLC", "5G_EIR", "SEPP", "N3IWF", "AF", "UDSF", "BSF", "CHF"},
			LoadThreshold:          10,
			RegWaitTime:            3,
			HTTP2Conns:             2,
			HTTP2Timeout:           60,
			KeepCacheRetryCount:    10,
		}
	}
	return instance
}

// GetFileName for fsnotify callback
func (f *FsNotifyOpt) GetFileName() string {
	return f.fileDir
}

// Handler for fsnotify callback
func (f *FsNotifyOpt) Handler(name, op string) {
	// For configmap, config file is a link to link.
	if op == "REMOVE" {
		log.Info("Reload configuration!")
		reLoad()
	}
}

func reLoad() {
	instance.loadFromFile()
}

func (c *Configuration) loadFromFile() {
	file, err := os.Open(configFile)

	if err == nil {
		decoder := json.NewDecoder(file)
		err = decoder.Decode(c)
		//instance.EnableIpv6 = c.EnableIpv6
		if err != nil {
			log.Errorf("Load configuration failed ! (%v)", err)
		}
	} else {
		log.Errorf("Open config file failed! (%v)", err)
	}
	log.Infof("%s,%+v", configFile, *c)
}

// InitCM for CM initialize
func InitCM() {
	fsWatcher = utils.InitFsWatcher()
	configFile = Opts.configFile
	GetInstance().loadFromFile()

	dir, err := filepath.Abs(filepath.Dir(configFile))
	if err != nil {
		log.Error(err)
	}
	file, _ := filepath.Abs(configFile)

	fs := &FsNotifyOpt{dir, file}
	if err := fsWatcher.AddFileToFsWatcher(fs); err != nil {
		log.Error("Can not add directory " + dir + " to fsnotify, err " + err.Error())
	}
}

// IsEnableIpv6 for open ipv6 swith
func IsEnableIpv6() bool {
	return GetInstance().EnableIpv6
}

// IsEnableConvertIpv6ToIpv4 for open Convert Ipv4 To Ipv6 swith
func IsEnableConvertIpv6ToIpv4() bool {
	return GetInstance().EnableIpv6ToIpv4
}

// GetDefaultRequesterNfList ..
func GetDefaultRequesterNfList() []string {
	return GetInstance().DefaultRequesterNfList
}

// GetDefaultValidityPeriod ..
func GetDefaultValidityPeriod() int {
	return GetInstance().DefaultValidityPeriod
}

//GetLoadThreshold for load upload patch threshold
func GetLoadThreshold() int {
	return GetInstance().LoadThreshold
}

//GetRegWaitTime for get register waiting time of all services up under same nfType
func GetRegWaitTime() int {
	return GetInstance().RegWaitTime
}

//GetHTTP2Conns ..
func GetHTTP2Conns() int {
	conns := GetInstance().HTTP2Conns
	if conns < 1 {
		conns = 1
	}
	return conns
}

//GetHTTP2Timeout ..
func GetHTTP2Timeout() int {
	to := GetInstance().HTTP2Timeout
	if to < 30 {
		to = 30
	}
	return to
}

//GetKeepCacheRetryCount ..
func GetKeepCacheRetryCount() int {
	return GetInstance().KeepCacheRetryCount
}
