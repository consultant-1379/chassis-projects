package nfdiscfilter

import (
	"testing"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"net/url"
)

func TestIsMatchedDataSet(t *testing.T){
	filter := &NFUDRInfoFilter{}
	dataset1 := "SUBSCRIPTION"
	dataset2 := "APPLICATION"
	nfprofile := []byte(`{
	      "supportedDataSets":["SUBSCRIPTION","POLICY","EXPOSURE"]
	}`)
	if !filter.isMatchedDataSet(dataset1, nfprofile){
		t.Fatal("dataset should be matched, but not !")
	}

	if filter.isMatchedDataSet(dataset2, nfprofile) {
		t.Fatal("dataset should not matche, but matched !")
	}
}


func TestUDRFilterByKVDB(t *testing.T) {
	InitAttributes(t)
	nfdiscutil.PreComplieRegexp()
	var queryform nfdiscrequest.DiscGetPara
	oriqueryForm, err := url.ParseQuery(`data-set=APPLICATION&requester-nf-type=AUSF&target-nf-type=UDM`)
	if err != nil {
		t.Fatal("url parse error")
	}
	queryform.InitMember(oriqueryForm)
	problem := queryform.ValidateNRFDiscovery()
	if problem != nil {
		t.Fatal("validate request fail")
	}
	filter := &NFUDRInfoFilter{}
	metaExpression := filter.filterByKVDB(&queryform)
	andExpression := buildAndExpression(metaExpression)
	var result string
	andExpression.metaExpressionToString(&result)
	if result != "AND{{where=data_set,value=APPLICATION,operation=0}}" {
		t.Fatal("udr filter by by kvdb fail.")
	}
}