package disc

import (
	"fmt"
	//"io/ioutil"
	"encoding/json"
	"net/http"
	"strings"

	"gerrit.ericsson.se/udm/common/pkg/httpserver"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/problemdetails"
	"gerrit.ericsson.se/udm/common/pkg/utils"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/common"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/cache"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/worker"
	//"github.com/buger/jsonparser"
)

func cacheSetup() {
	httpserver.Routes = append(httpserver.Routes,
		httpserver.Route{
			Name:        "cacheFlush",
			Method:      "DELETE",
			Pattern:     "/nrf-discovery-agent/v1/memcache/{nfType}",
			HandlerFunc: cacheFlush,
		},
		httpserver.Route{
			Name:        "cacheFlushRoam",
			Method:      "DELETE",
			Pattern:     "/nrf-discovery-agent/v1/memcache/{nfType}-roam",
			HandlerFunc: cacheFlushRoam,
		},
		httpserver.Route{
			Name:        "cacheDump",
			Method:      "GET",
			Pattern:     "/nrf-discovery-agent/v1/memcache/{nfType}",
			HandlerFunc: cacheDump,
		},
		httpserver.Route{
			Name:        "cacheSync",
			Method:      "GET",
			Pattern:     "/nrf-discovery-agent/v1/synccache/{nfType}",
			HandlerFunc: cacheSync,
		},
		httpserver.Route{
			Name:        "cacheFlushAll",
			Method:      "DELETE",
			Pattern:     "/nrf-discovery-agent/v1/memcache",
			HandlerFunc: cacheFlushAll,
		},
		httpserver.Route{
			Name:        "cacheDumpAll",
			Method:      "GET",
			Pattern:     "/nrf-discovery-agent/v1/memcache",
			HandlerFunc: cacheDumpAll,
		},
		//		httpserver.Route{
		//			Name:        "cacheProvisioning",
		//			Method:      "PUT",
		//			Pattern:     "/nrf-discovery-agent/v1/memcache",
		//			HandlerFunc: cacheProvisioningHandler,
		//		},
		httpserver.Route{
			Name:        "discWorkMode",
			Method:      "GET",
			Pattern:     "/nrf-discovery-agent/v1/work-mode",
			HandlerFunc: keepCacheStatusGetHandler,
		},
	)
}

func cacheDump(rw http.ResponseWriter, req *http.Request) {
	var sequenceID string = ""
	if log.GetLevel() >= log.DebugLevel {
		sequenceID = utils.GetSequenceId()
	}
	logcontent := &log.LogStruct{SequenceId: sequenceID}

	log.Debug("Cache dump cache query comes")

	uri := req.URL.RequestURI()
	arrays := strings.Split(uri, "/")
	nfType := strings.ToUpper(arrays[len(arrays)-1])
	log.Infof("Dump cache nfType = %s", nfType)

	contents := make([]byte, 0)
	dumpData := structs.CacheDumpData{
		RequestNfType: nfType,
	}
	cache.Instance().Dump(nfType, &dumpData)
	length := len(dumpData.CacheInfos)
	if length == 0 {
		errorInfo := "Cache is empty"
		problemDetails := &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
		logcontent.ResponseDescription = errorInfo
		handleCacheOperationFailure(rw, req, logcontent, http.StatusNotFound, problemDetails.ToString())

		return
	}

	contents, err := json.Marshal(dumpData)
	if err != nil {
		log.Errorf("cacheDump: Marshal dumpData failure: %s", err.Error())
		return
	}

	handleCacheOperationSuccess(rw, req, logcontent, http.StatusOK, string(contents))
}

func cacheSync(rw http.ResponseWriter, req *http.Request) {
	var sequenceID string = ""
	if log.GetLevel() >= log.DebugLevel {
		sequenceID = utils.GetSequenceId()
	}
	logcontent := &log.LogStruct{SequenceId: sequenceID}

	log.Debug("Cache sync cache query comes")

	uri := req.URL.RequestURI()
	arrays := strings.Split(uri, "/")
	nfType := strings.ToUpper(arrays[len(arrays)-1])
	log.Infof("Sync cache nfType = %s", nfType)

	ok := worker.Instance().DumpReady(nfType)
	if !ok {
		errorInfo := fmt.Sprintf("Cache not ready for nfType:%s", nfType)
		problemDetails := &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
		logcontent.ResponseDescription = errorInfo
		handleCacheOperationFailure(rw, req, logcontent, http.StatusForbidden, problemDetails.ToString())

		return
	}

	contents := make([]byte, 0)
	syncData := structs.CacheSyncData{
		RequestNfType: nfType,
	}
	cache.Instance().Sync(nfType, &syncData)
	length := len(syncData.CacheInfos)
	if length == 0 {
		errorInfo := "Cache is empty"
		problemDetails := &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
		logcontent.ResponseDescription = errorInfo
		handleCacheOperationFailure(rw, req, logcontent, http.StatusNotFound, problemDetails.ToString())

		return
	}

	contents, err := json.Marshal(syncData)
	if err != nil {
		log.Errorf("cacheSync: Marshal syncData failure: %s", err.Error())
		return
	}

	handleCacheOperationSuccess(rw, req, logcontent, http.StatusOK, string(contents))
}

func cacheDumpAll(rw http.ResponseWriter, req *http.Request) {
	var sequenceID string = ""
	if log.GetLevel() >= log.DebugLevel {
		sequenceID = utils.GetSequenceId()
	}
	logcontent := &log.LogStruct{SequenceId: sequenceID}

	log.Debug("Cache dump cache all query comes")
	dumpDatas := make([]structs.CacheDumpData, 0)
	cache.Instance().DumpAll(&dumpDatas)
	length := len(dumpDatas)
	if length == 0 {
		errorInfo := "Cache is empty"
		problemDetails := &problemdetails.ProblemDetails{
			Title: errorInfo,
		}
		logcontent.ResponseDescription = errorInfo
		handleCacheOperationFailure(rw, req, logcontent, http.StatusNotFound, problemDetails.ToString())

		return
	}

	contents, err := json.Marshal(dumpDatas)
	if err != nil {
		log.Errorf("cacheDump: Marshal dumpData failure: %s", err.Error())
		return
	}

	handleCacheOperationSuccess(rw, req, logcontent, http.StatusOK, string(contents))
}

func cacheFlush(rw http.ResponseWriter, req *http.Request) {
	var sequenceID string
	if log.GetLevel() >= log.DebugLevel {
		sequenceID = utils.GetSequenceId()
	}
	logcontent := &log.LogStruct{SequenceId: sequenceID}

	log.Debug("Cache flush query comes")

	uri := req.URL.RequestURI()
	arrays := strings.Split(uri, "/")
	nfType := strings.ToUpper(arrays[len(arrays)-1])

	cache.Instance().Flush(nfType)
	log.Infof("Flush cache nfType = %s", nfType)

	respData := fmt.Sprintf("Clean %s cache success", nfType)
	handleCacheOperationSuccess(rw, req, logcontent, http.StatusOK, respData)
}

func cacheFlushRoam(rw http.ResponseWriter, req *http.Request) {
	var sequenceID string
	if log.GetLevel() >= log.DebugLevel {
		sequenceID = utils.GetSequenceId()
	}
	logcontent := &log.LogStruct{SequenceId: sequenceID}

	log.Debug("Cache flushRoam query comes")

	uri := req.URL.RequestURI()
	arrays := strings.Split(uri, "/")
	nfType := common.GetReqNfTypeForRoam(strings.ToUpper(arrays[len(arrays)-1]))

	cache.Instance().FlushRoam(nfType)
	log.Infof("FlushRoam cache nfType = %s", nfType)

	respData := fmt.Sprintf("Clean %s Roamcache success", nfType)
	handleCacheOperationSuccess(rw, req, logcontent, http.StatusOK, respData)
}

func cacheFlushAll(rw http.ResponseWriter, req *http.Request) {
	var sequenceID string
	if log.GetLevel() >= log.DebugLevel {
		sequenceID = utils.GetSequenceId()
	}
	logcontent := &log.LogStruct{SequenceId: sequenceID}

	log.Debug("Cache flush all query comes")

	cache.Instance().FlushAll()
	log.Info("Flush cache all")

	respData := "Clean cache all success"
	handleCacheOperationSuccess(rw, req, logcontent, http.StatusOK, respData)
}

//func cacheProvisioningHandler(rw http.ResponseWriter, req *http.Request) {
//	var sequenceID string
//	if log.GetLevel() >= log.DebugLevel {
//		sequenceID = utils.GetSequenceId()
//	}
//	logcontent := &log.LogStruct{SequenceId: sequenceID}

//	log.Debug("Cache provisioning query comes")

//	body, err := ioutil.ReadAll(req.Body)
//	if err != nil {
//		errorInfo := "Provision body is empty"
//		problemDetails := &problemdetails.ProblemDetails{
//			Title: errorInfo,
//		}
//		logcontent.ResponseDescription = errorInfo
//		handleCacheOperationFailure(rw, req, logcontent, http.StatusBadRequest, problemDetails.ToString())

//		return
//	}
//	defer close(req)

//	validityPeriod, err := jsonparser.GetInt(body, "validityPeriod")
//	if err != nil {
//		errorInfo := "Provision data miss validityPeriod"
//		problemDetails := &problemdetails.ProblemDetails{
//			Title: errorInfo,
//		}
//		logcontent.ResponseDescription = errorInfo
//		handleCacheOperationFailure(rw, req, logcontent, http.StatusBadRequest, problemDetails.ToString())

//		return
//	}

//	requestDescription := fmt.Sprintf(`{"validityPeriod":"%d"}`, validityPeriod)
//	log.Debugf(requestLogFormat, sequenceID, req.URL, req.Method, requestDescription)

//	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
//		if nfType, e := jsonparser.GetString(value, "nfType"); e == nil {
//			cache.Instance().Cached(nfType, nfType, value)
//		} else {
//			log.Errorf("cacheProvisioningHandler: NF profile in search result does not contain nfType")
//		}
//	}, "nfInstances")
//	if err != nil {
//		errorInfo := err.Error()
//		problemDetails := &problemdetails.ProblemDetails{
//			Title: errorInfo,
//		}
//		logcontent.ResponseDescription = errorInfo
//		handleCacheOperationFailure(rw, req, logcontent, http.StatusBadRequest, problemDetails.ToString())

//		return
//	}

//	respInfo := "Provision data to cache success"
//	handleCacheOperationSuccess(rw, req, logcontent, http.StatusCreated, respInfo)
//}

func handleCacheOperationFailure(rw http.ResponseWriter, req *http.Request, logcontent *log.LogStruct, statusCode int, body string) {
	log.Debugf(consts.REQUEST_LOG_FORMAT, logcontent.SequenceId, req.URL, req.Method, logcontent.RequestDescription)
	cacheOperationResponseHander(rw, req, statusCode, body)
	log.Errorf(consts.RESPONSE_LOG_FORMAT, logcontent.SequenceId, statusCode, logcontent.ResponseDescription)
}

func handleCacheOperationSuccess(rw http.ResponseWriter, req *http.Request, logcontent *log.LogStruct, statusCode int, body string) {
	log.Debugf(consts.REQUEST_LOG_FORMAT, logcontent.SequenceId, req.URL, req.Method, logcontent.RequestDescription)
	cacheOperationResponseHander(rw, req, statusCode, body)
	log.Debugf(consts.RESPONSE_LOG_FORMAT, logcontent.SequenceId, statusCode, logcontent.ResponseDescription)
}

func cacheOperationResponseHander(rw http.ResponseWriter, req *http.Request, statusCode int, body string) {
	if statusCode == http.StatusOK {
		rw.Header().Set("Content-Type", httpContentTypeJSON)
	} else {
		rw.Header().Set("Content-Type", httpContentTypeProblemJSON)
	}

	rw.WriteHeader(statusCode)

	if body != "" {
		_, err := rw.Write([]byte(body))
		if err != nil {
			log.Warnf("%v", err)
		}
	}
}

func keepCacheStatusGetHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", httpContentTypeJSON)
	rw.WriteHeader(http.StatusOK)
	retContent := make([]byte, 0)
	if worker.IsKeepCacheMode() {
		retContent = []byte(`Discovery agent work mode is keep cache mode.`)
	} else {
		retContent = []byte(`Discovery agent work mode is normal mode.`)
	}

	_, _ = rw.Write(retContent)
}
