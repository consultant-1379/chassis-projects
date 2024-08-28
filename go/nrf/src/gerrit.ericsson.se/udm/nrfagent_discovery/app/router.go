package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/cache"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
)

func CheckHealth(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

// CheckReadiness return StatusOK when server is ready for service
func CheckReadiness(rw http.ResponseWriter, req *http.Request) {
	if ServerStatus == consts.ServerIsRunning {
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func GetEnvs(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	if body, err := getHostAllEnvInfof(); err == nil {
		_, _ = rw.Write(body)
	}
}

func Flush(rw http.ResponseWriter, req *http.Request) {
	cache.Instance().FlushAll()
	rw.WriteHeader(http.StatusOK)
	_, _ = rw.Write([]byte("Clean cache success"))
}

func GetOpts(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.WriteHeader(http.StatusOK)
	buf, err := json.MarshalIndent(cm.Opts, "", "  ")
	if err != nil {
		rw.WriteHeader(500)
		_, _ = rw.Write([]byte(fmt.Sprintf("Marshal encode failed: %s \n", err.Error())))
		return
	}
	_, _ = rw.Write(buf)
}

//func GetConfs(rw http.ResponseWriter, req *http.Request) {
//	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
//	rw.WriteHeader(http.StatusOK)
//	conflist := cm.GetInstance().GetAllConfigFiles()
//	buf, err := json.MarshalIndent(conflist, "", "  ")
//	if err != nil {
//		rw.WriteHeader(500)
//		_, _ = rw.Write([]byte(fmt.Sprintf("Marshal encode failed: %v \n", err)))
//		return
//	}
//	_, _ = rw.Write(buf)
//}

func getHostAllEnvInfof() ([]byte, error) {
	environment := make(map[string]string)
	for _, item := range os.Environ() {
		splits := strings.Split(item, "=")
		key := splits[0]
		val := strings.Join(splits[1:], "=")
		environment[key] = val
	}

	return json.MarshalIndent(environment, "", "  ")
}

////////////////////below///////////////////////
