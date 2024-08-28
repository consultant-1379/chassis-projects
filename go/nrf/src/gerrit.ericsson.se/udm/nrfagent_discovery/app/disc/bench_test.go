package disc

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	//"gerrit.ericsson.se/udm/common/pkg/log"
	//"gerrit.ericsson.se/udm/nrfagent_common/pkg/common"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
)

func StubSetupBenchmark() {
	cacheManager.SetTargetNf("UDM", structs.TargetNf{
		RequesterNfType:    "UDM",
		TargetNfType:       "UDR",
		TargetServiceNames: []string{"nudr-dr"},
	})
	cacheManager.SetTargetNf("UDM", structs.TargetNf{
		RequesterNfType:    "UDM",
		TargetNfType:       "AMF",
		TargetServiceNames: []string{"namf-evts", "namf-mt", "namf-loc"},
	})
	cacheManager.SetTargetNf("AUSF", structs.TargetNf{
		RequesterNfType:    "AUSF",
		TargetNfType:       "UDM",
		TargetServiceNames: []string{"nudm-auth-01", "nudm-uecm"},
	})

	cacheManager.SetRequesterFqdn("UDM", "seliius03696.seli.gic.ericsson.se")
	cacheManager.SetRequesterFqdn("AUSF", "seliius03696.seli.gic.ericsson.se")
}

func TestHandleNfDiscovery(t *testing.T) {
	StubSetupBenchmark()
	StubHTTPDoToNrf("GET", http.StatusOK)

	resp := httptest.NewRecorder()

	req := &http.Request{
		Method: "GET",
		Form:   nil,
		URL: &url.URL{
			RawQuery: "service-names=nudm-uecm&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&target-nf-type=UDM&requester-nf-type=AUSF&supi=imsi-600<imsi>9999&snssais=%7B%22sst%22%3A+0,%22sd%22%3A+%22000000%22%7D&snssais=%7B%22sst%22%3A+1,%22sd%22%3A+%22111111%22%7D",
		},
	}
	t.Logf("%v", req)
	//nfDiscoveryRequestHandler(nil, req)
	nfDiscoveryRequestHandler(resp, req)
}

//func benchHandleNfDiscoveryRequest(rawQuery string, b *testing.B) {
//	StubSetupBenchmark()
//	StubHTTPDoToNrf("GET", http.StatusOK)

//	resp := httptest.NewRecorder()

//	req := &http.Request{
//		Method: "GET",
//		Form:   nil,
//		URL: &url.URL{
//			RawQuery: rawQuery,
//		},
//	}
//	b.Logf("%v", req)
//	for n := 0; n < b.N; n++ {
//		b.Logf("%d", n)
//		//nfDiscoveryRequestHandler(nil, req)
//		nfDiscoveryRequestHandler(resp, req)
//	}
//}

//func BenchmarkNfDiscBase(b *testing.B) {
//	benchHandleNfDiscoveryRequest("requester-nf-type=UDM&target-nf-type=AMF", b)
//}
//func BenchmarkNfDiscWithServiceNames(b *testing.B) {
//	benchHandleNfDiscoveryRequest("requester-nf-type=UDM&target-nf-type=UDR&service-names=nudr-dr", b)
//}
//func BenchmarkNfDiscWithMultiServiceNames(b *testing.B) {
//	benchHandleNfDiscoveryRequest("requester-nf-type=UDM&target-nf-type=UDR&service-names=namf-evts&service-names=namf-mt&service-names=namf-loc", b)
//}
//func BenchmarkNfDiscNST(b *testing.B) {
//	benchHandleNfDiscoveryRequest("service-names=nudm-auth-01&requester-nf-instance-fqdn=seliius03696.seli.gic.ericsson.se&target-nf-type=UDM&requester-nf-type=AUSF&supi=imsi-600<imsi>9999&snssais=%7B%22sst%22%3A+0,%22sd%22%3A+%22000000%22%7D&snssais=%7B%22sst%22%3A+1,%22sd%22%3A+%22111111%22%7D", b)
//}

//func benchHandleDiscoveryRequest(targetNf *structs.TargetNf, nfInstanceID string, b *testing.B) {
//	// stub comm env and targetNfProfile
//	StubSetupBenchmark()
//	// set logging level to fatal
//	// log.SetLevel(log.DebugLevel)
//	log.SetLevel(log.FatalLevel)
//	// stub requester NF FQDN
//	cacheManager.SetReqFqdn(targetNf.RequesterNfType, "seliius03696.seli.gic.ericsson.se")
//	// stub msgbus
//	common.SetDiscMsgbus(nil)
//	// stub http2Nrf
//	StubHTTPDoToNrf("GET", http.StatusOK)
//	// stub cache
//	cacheManager.Flush(targetNf.RequesterNfType)

//	for n := 0; n < b.N; n++ {
//		handleDiscoveryRequest(targetNf, nfInstanceID)
//	}
//}

//func BenchmarkNrfDiscBase(b *testing.B) {
//	targetNf := &structs.TargetNf{
//		RequesterNfType:    "NAUSF",
//		TargetNfType:       "UDM",
//		TargetServiceNames: []string{"nudm-auth-01"},
//	}
//	benchHandleDiscoveryRequest(targetNf, "", b)
//}

//func BenchmarkNrfDiscWithCache(b *testing.B) {
//	targetNf := &structs.TargetNf{
//		RequesterNfType:    "AUSF",
//		TargetNfType:       "UDM",
//		TargetServiceNames: []string{"nudm-auth-01"},
//	}
//	benchHandleDiscoveryRequest(targetNf, "", b)
//}
