package disc

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/cache"
)

func TestCacheSetup(t *testing.T) {
	cacheSetup()
}

func TestCacheDump(t *testing.T) {
	cacheManager.Flush("AUSF")
	log.SetLevel(log.LevelUint("DEBUG"))

	t.Run("TestCacheDump_AUSF_Empty", func(t *testing.T) {
		cacheManager.Flush("AUSF")
		var nobody = []byte("")
		resp := httptest.NewRecorder()
		resp.Code = http.StatusOK

		req := httptest.NewRequest("GET", "/nrf-discovery-agent/v1/memcache/AUSF", bytes.NewBuffer(nobody))
		req.Header.Set("Content-Type", "application/json")

		cacheDump(resp, req)
		if resp.Code != http.StatusOK {
			t.Errorf("TestCacheDump: cache dump failed, response code is %d", resp.Code)
		}
	})

	t.Run("TestCacheDump_AUSF_NonEmpty", func(t *testing.T) {
		nfinstanceByte, _, _ := cache.SpliteSeachResult(searchResultUDM)
		if nfinstanceByte == nil {
			t.Errorf("TestCacheDump: SpliteSeachResult fail")
		}
		for _, instance := range nfinstanceByte {
			cacheManager.Cached("AUSF", "UDM", instance, false)
			ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
			if !ok {
				t.Errorf("TestCacheDump: Cached fail")
			}
		}

		var nobody = []byte("")
		resp := httptest.NewRecorder()
		resp.Code = http.StatusOK

		req := httptest.NewRequest("GET", "/nrf-discovery-agent/v1/memcache/AUSF", bytes.NewBuffer(nobody))
		req.Header.Set("Content-Type", "application/json")

		cacheDump(resp, req)
		if resp.Code != http.StatusOK {
			t.Errorf("TestCacheDump: cache dump failed, response code is %d", resp.Code)
		}

		cacheManager.Flush("AUSF")
	})

}

func TestCacheDumpAll(t *testing.T) {
	log.SetLevel(log.LevelUint("ERROR"))

	t.Run("TestCacheDumpAll_AUSF_Empty", func(t *testing.T) {
		cacheManager.Flush("AUSF")
		var nobody = []byte("")
		resp := httptest.NewRecorder()
		resp.Code = http.StatusOK

		req := httptest.NewRequest("GET", "/nrf-discovery-agent/v1/memcache", bytes.NewBuffer(nobody))
		req.Header.Set("Content-Type", "application/json")

		cacheDumpAll(resp, req)
		if resp.Code != http.StatusOK {
			t.Errorf("TestCacheDumpAll_AUSF_Empty: cache dump all failed, response code is %d", resp.Code)
		}
	})

	t.Run("TestCacheDumpAll_AUSF_NonEmpty", func(t *testing.T) {
		cacheManager.Flush("AUSF")
		nfinstanceByte, _, _ := cache.SpliteSeachResult(searchResultUDM)
		if nfinstanceByte == nil {
			t.Errorf("TestCacheDumpAll_AUSF_NonEmpty: SpliteSeachResult fail")
		}
		for _, instance := range nfinstanceByte {
			cacheManager.Cached("AUSF", "UDM", instance, false)
			ok := cacheManager.Probe("AUSF", "UDM", "udm-5g-01")
			if !ok {
				t.Errorf("TestCacheDumpAll_AUSF_NonEmpty: Cached fail")
			}
		}

		var nobody = []byte("")
		resp := httptest.NewRecorder()
		resp.Code = http.StatusOK

		req := httptest.NewRequest("GET", "/nrf-discovery-agent/v1/memcache/AUSF", bytes.NewBuffer(nobody))
		req.Header.Set("Content-Type", "application/json")

		cacheDumpAll(resp, req)
		if resp.Code != http.StatusOK {
			t.Errorf("TestCacheDumpAll_AUSF_NonEmpty: cache dump failed, response code is %d", resp.Code)
		}

		cacheManager.Flush("AUSF")
	})
}

func TestCacheSync(t *testing.T) {
	cacheManager.Flush("AUSF")
	log.SetLevel(log.LevelUint("DEBUG"))

	t.Run("TestCacheSync_Dump_Not_Ready", func(t *testing.T) {
		cacheManager.Flush("AUSF")
		var nobody = []byte("")
		resp := httptest.NewRecorder()
		resp.Code = http.StatusOK

		req := httptest.NewRequest("GET", "/nrf-discovery-agent/v1/synccache/AUSF", bytes.NewBuffer(nobody))
		req.Header.Set("Content-Type", "application/json")

		cacheSync(resp, req)
		if resp.Code != http.StatusOK {
			t.Errorf("TestCacheSync_Dump_Not_Ready: cache sync failed, response code is %d", resp.Code)
		}
	})
}

func TestCacheFlush(t *testing.T) {
	var nobody = []byte("")
	resp := httptest.NewRecorder()
	resp.Code = http.StatusOK

	req := httptest.NewRequest("DELETE", "/nrf-discovery-agent/v1/memcache/AUSF", bytes.NewBuffer(nobody))
	req.Header.Set("Content-Type", "application/json")

	cacheFlush(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("TestCacheFlush: cache flush failed, response code is %d", resp.Code)
	}
}

func TestCacheFlushRoam(t *testing.T) {
	var nobody = []byte("")
	resp := httptest.NewRecorder()
	resp.Code = http.StatusOK

	req := httptest.NewRequest("DELETE", "/nrf-discovery-agent/v1/memcache/AUSF-roam", bytes.NewBuffer(nobody))
	req.Header.Set("Content-Type", "application/json")

	cacheFlushRoam(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("TestCacheFlushRoam: cache flush failed, response code is %d", resp.Code)
	}
}

func TestCacheFlushAll(t *testing.T) {
	var nobody = []byte("")
	resp := httptest.NewRecorder()
	resp.Code = http.StatusOK

	req := httptest.NewRequest("DELETE", "/nrf-discovery-agent/v1/memcache/AUSF-roam", bytes.NewBuffer(nobody))
	req.Header.Set("Content-Type", "application/json")

	cacheFlushAll(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("TestCacheFlushAll: cache flush failed, response code is %d", resp.Code)
	}
}

func TestHandleCacheOperation(t *testing.T) {
	var nobody = []byte("")
	resp := httptest.NewRecorder()
	resp.Code = http.StatusOK

	req := httptest.NewRequest("DELETE", "/nrf-discovery-agent/v1/memcache/AUSF-roam", bytes.NewBuffer(nobody))
	req.Header.Set("Content-Type", "application/json")
	logcontent := &log.LogStruct{SequenceId: "pod-1"}
	logcontent.ResponseDescription = "Handle cache"

	t.Run("TestHandleCacheOperationFailure", func(t *testing.T) {
		handleCacheOperationFailure(resp, req, logcontent, 404, "")
	})

	t.Run("TestCacheOperationResponseHander", func(t *testing.T) {
		cacheOperationResponseHander(resp, req, 200, "cache handler")
	})

	t.Run("TestKeepCacheStatusGetHandler", func(t *testing.T) {
		keepCacheStatusGetHandler(resp, req)
	})
}
