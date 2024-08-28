package nfdiscfilter

import (
	"testing"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"net/url"
)
func TestIsMatchedUpfIwkEpsInd(t *testing.T) {
	filter := &NFUPFInfoFilter{}
	nfProfile := []byte(`{
	     "sNssaiUpfInfoList": [
		 {
		     "sNssai": {
			 "sst": 3,
			 "sd": "sd3"
		     },
		     "dnnUpfInfoList": [
			 {
			     "dnn": "dnn4"
			 },
			 {
			     "dnn": "dnn5"
			 }
		     ]
		 }
	     ],
	     "iwkEpsInd": true
	 }`)

	ok := filter.isMatchedUpfIwkEpsInd(false, nfProfile)
	if ok {
		t.Fatal("upf-iwk-eps-ind should not be matched, but matched !")
	}
	ok2 := filter.isMatchedUpfIwkEpsInd(true, nfProfile)
	if !ok2 {
		t.Fatal("upf-iwk-eps-ind should be matched, but not !")
	}

	nfProfile2 := []byte(`{
	     "sNssaiUpfInfoList": [
		 {
		     "sNssai": {
			 "sst": 3,
			 "sd": "sd3"
		     },
		     "dnnUpfInfoList": [
			 {
			     "dnn": "dnn4"
			 },
			 {
			     "dnn": "dnn5"
			 }
		     ]
		 }
	     ]
	 }`)

	ok3 := filter.isMatchedUpfIwkEpsInd(false, nfProfile2)
	if !ok3 {
		t.Fatal("upf-iwk-eps-ind should not be matched, but matched !")
	}
}

func TestUPFFilterByKVDB(t *testing.T) {
	InitAttributes(t)
	nfdiscutil.PreComplieRegexp()
	var queryform nfdiscrequest.DiscGetPara
	oriqueryForm, err := url.ParseQuery(`smf-serving-area=123&requester-nf-type=AUSF&target-nf-type=UDM`)
	if err != nil {
		t.Fatal("url parse error")
	}
	queryform.InitMember(oriqueryForm)
	queryform.ValidateNRFDiscovery()
	filter := &NFUPFInfoFilter{}
	metaExpression := filter.filterByKVDB(&queryform)
	andExpression := buildAndExpression(metaExpression)
	var result string
	andExpression.metaExpressionToString(&result)
	if result != "AND{{where=smf_serving_area,value=123,operation=0}}" {
		t.Fatal("upf filter by by kvdb fail.")
	}
}