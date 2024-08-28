package nfdiscfilter


import (
	"testing"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"net/url"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap"
	"path/filepath"
	"os"
	"io/ioutil"
)

var attributes = []byte(`{
		"common":[
			{
			    "parameter": "snssais/sst",
			    "path": "body.sNssais.sst",
			    "from": "value.helper.sNssais snssai",
			    "where": "snssai.sst",
			    "exist_check": false
			},
			{
			    "parameter": "snssais/sd",
			    "path": "body.sNssais.sd",
			    "from": "value.helper.sNssais snssai",
			    "where": "snssai.sd",
			    "exist_check": false
			}
		],
		"smf":[
			{
			    "parameter": "snssais/sst",
			    "path": "body.smfInfo.sNssaiSmffInfoList.sNssai.sst",
			    "from": "value.helper.smfInfo.sNssaiSmfInfoList snssai_smf_info",
			    "where": "snssai_smf_info.sNssai.sst",
			    "exist_check": false
			},
			{
			    "parameter": "snssais/sd",
			    "path": "body.smfInfo.sNssaiSmfInfoList.sNssai.sd",
			    "from": "value.helper.smfInfo.sNssaiSmfInfoList snssai_smf_info",
			    "where": "snssai_smf_info.sNssai.sd",
			    "exist_check": false
			},
			{
			    "parameter": "dnn",
			    "path": "body.smfInfo.sNssaiSmfInfoList.dnnSmfInfoList.dnn",
			    "from": "value.helper.smfInfo.sNssaiSmfInfoList snssai_smf_info, snssai_smf_info.dnnSmfInfoList dnn_smf_info",
			    "where": "dnn_smf_info.dnn",
			    "exist_check": false
			},
			{
			    "parameter": "pgw",
			    "path": "body.smfInfo.pgwFqdn",
			    "from": "",
			    "where": "value.body.smfInfo.pgwFqdn",
			    "exist_check": false
			}
		],
		"upf": [
			{
			    "parameter": "smf-serving-area",
			    "path": "body.upfInfo.smfServingArea",
			    "from": "value.body.upfInfo.smfServingArea smf_serving_area",
			    "where": "smf_serving_area",
			    "exist_check": false
			}
		],
		"udr": [
			{
			    "parameter": "data-set",
			    "path": "body.udrInfo.supportedDataSets",
			    "from": "value.body.udrInfo.supportedDataSets data_set",
			    "where": "data_set",
			    "exist_check": false
			}
		]
        }`)
var attributesConfFile = "attributes.json"
func InitAttributes(t *testing.T)  {
	var buffer = []byte(attributes)
	currDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	AbsFileName := currDir + "/" + attributesConfFile

	err := ioutil.WriteFile(AbsFileName, buffer, 0666)
	if err != nil {
		t.Fatal(err.Error())
	}

	attributesConf := &configmap.AttributesConf{
		FileName:     "",
		ConfInstance: &configmap.SubAttributesConf{},
	}

	attributesConf.SetFileName(AbsFileName)

	if AbsFileName != attributesConf.GetFileName() {
		t.Fatalf(`attributesConf.GetFileName() didn't return expected FileName`)
	}

	attributesConf.LoadConf()
}
func TestIsMatchedByPGWInd(t *testing.T){
	filter := &NFSMFInfoFilter{}
        nfprofile := []byte(`{
            "pgwFqdn":"seliius19972.ericsson.se"
        }`)

	nfprofile2 := []byte(`{
            "pgwFqdn":" "
        }`)

	req := &nfdiscrequest.DiscGetPara{}
	var pgwind1 []string
	pgwind1 = append(pgwind1, "true")
	value := make(map[string][]string)
	value[constvalue.SearchDataPGWInd] = pgwind1
	req.InitMember(value)
	req.SetFlag(constvalue.SearchDataPGWInd ,true)
	if !filter.isMatchedByPGWInd(nfprofile, req){
		t.Fatal("pgwind is true, pgwFqdn has value, should match, but not!")
	}

	if filter.isMatchedByPGWInd(nfprofile2, req) {
		t.Fatal("pgwind is true, pgwFqdn no value, should not match, but match")
	}

	req2 := &nfdiscrequest.DiscGetPara{}
	var pgwind2 []string
	pgwind2 = append(pgwind2, "false")
	value2 := make(map[string][]string)
	value2[constvalue.SearchDataPGWInd] = pgwind2
	req2.InitMember(value2)
	req2.SetFlag(constvalue.SearchDataPGWInd ,true)
	if filter.isMatchedByPGWInd(nfprofile, req2){
		t.Fatal("pgwind is false, pgwFqdn has value, should not match, but match!")
	}

	if !filter.isMatchedByPGWInd(nfprofile2, req2) {
		t.Fatal("pgwind is false, pgwfqdn no value, should match, but not")
	}
}

func TestIsMatchedAccessType(t *testing.T) {
	filter := &NFSMFInfoFilter{}
	nfprofile := []byte(`{
            "accessType": ["3GPP_ACCESS"]
        }`)

	nfprofile2 := []byte(`{

        }`)

	nfprofile3 := []byte(`{
            "accessType": ["3GPP_ACCESS", "NON_3GPP_ACCESS"]
        }`)


	if !filter.isMatchedAccessType(nfprofile, constvalue.Access3GPP){
		t.Fatal("access-type should match, but not!")
	}

	if !filter.isMatchedAccessType(nfprofile2, constvalue.Access3GPP) {
		t.Fatal("access-type should match, but not")
	}

	if filter.isMatchedAccessType(nfprofile3, "TEST"){
		t.Fatal("access-type should not match, but match!")
	}

}

func TestSMFFilterByKVDB(t *testing.T) {
	InitAttributes(t)
	nfdiscutil.PreComplieRegexp()
	var queryform nfdiscrequest.DiscGetPara
	oriqueryForm, err := url.ParseQuery(`pgw=123&requester-nf-type=AUSF&target-nf-type=UDM`)
	if err != nil {
		t.Fatal("url parse error")
	}
	queryform.InitMember(oriqueryForm)
	queryform.ValidateNRFDiscovery()
	filter := &NFSMFInfoFilter{}
	metaExpression := filter.filterByKVDB(&queryform)
	andExpression := buildAndExpression(metaExpression)
	var result string
	andExpression.metaExpressionToString(&result)
	if result != "AND{{where=value.body.smfInfo.pgwFqdn,value=123,operation=0}}" {
		t.Fatal("smf filter by by kvdb fail.")
	}
}