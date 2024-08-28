package configmap

import (
	"path/filepath"
	"os"
	"io/ioutil"
	"testing"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

var attributesConfFile = "attributes.json"
func TestAttributesConf(t *testing.T) {
	jsonContent := `{
		"common":[
			{
			"parameter": "target-nf-type",
			"path": "body.nfType",
			"from": "",
			"where": "value.body.nfType",
			"exist_check": false
			},
			{
			"parameter": "nfStatus",
			"path": "body.nfStatus",
			"from": "",
			"where": "value.body.nfStatus",
			"exist_check": false
			}
		],
		"udr":[
			{
			"parameter": "supi",
			"path": "body.udrInfo.supiRanges.start",
			"from": "value.helper.udrInfo.supiRanges supi",
			"where": "supi.start",
			"exist_check": false
			}
		],
		"udm":[
			{
			"parameter": "supi",
			"path": "body.udmInfo.supiRanges.start",
			"from": "value.helper.udmInfo.supiRanges supi",
			"where": "supi.start",
			"exist_check": false
			}
		],
		"ausf":[
			{
			"parameter": "supi",
			"path": "body.ausfInfo.supiRanges.start",
			"from": "value.helper.ausfInfo.supiRanges supi",
			"where": "supi.start",
			"exist_check": false
			}
		],
		"amf":[
			{
			"parameter": "tai/plmnid/mcc",
			"path": "body.amfInfo.taiList.plmnId.mcc",
			"from": "value.helper.amfInfo.taiList tai",
			"where": "tai.plmnId.mcc",
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
			}
		],
		"upf":[
			{
			"parameter": "snssais/sst",
			"path": "body.upfInfo.sNssaiUpfInfoList.sNssai.sst",
			"from": "value.helper.upfInfo.sNssaiUpfInfoList snssai_upf_info",
			"where": "snssai_upf_info.sNssai.sst",
			"exist_check": false
			}
		],
		"pcf":[
			{
			"parameter": "dnn",
			"path": "body.pcfInfo.dnnList",
			"from": "value.body.pcfInfo.dnnList dnn",
			"where": "dnn",
			"exist_check": false
			}
		],
		"bsf":[
			{
			"parameter": "dnn",
			"path": "body.bsfInfo.dnnList",
			"from": "value.body.bsfInfo.dnnList dnn",
			"where": "dnn",
			"exist_check": false
			}
		],
		"chf":[
			{
			"parameter": "supi",
			"path": "body.chfInfo.supiRangeList.start",
			"from": "value.helper.chfInfo.supiRangeList supi",
			"where": "supi.start",
			"exist_check": false
			}
		]
        }`

	var buffer = []byte(jsonContent)
	currDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	AbsFileName := currDir + "/" + attributesConfFile

	err := ioutil.WriteFile(AbsFileName, buffer, 0666)
	if err != nil {
		t.Fatal(err.Error())
	}

	attributesConf := &AttributesConf{
		FileName:     "",
		ConfInstance: &SubAttributesConf{},
	}

	attributesConf.SetFileName(AbsFileName)

	if AbsFileName != attributesConf.GetFileName() {
		t.Fatalf(`attributesConf.GetFileName() didn't return expected FileName`)
	}

	attributesConf.LoadConf()
	if AttributesMap[constvalue.Common][constvalue.SearchDataTargetNfType].Path != "body.nfType" {
		t.Fatalf("AttributesMap should have the correct value, but not ")
	}
	if AttributesMap[constvalue.Common][constvalue.NfStatus].Path != "body.nfStatus" {
		t.Fatalf("AttributesMap should have the correct value, but not ")
	}
	if AttributesMap[constvalue.NfTypeCHF][constvalue.SearchDataSupi].Path != "body.chfInfo.supiRangeList.start" {
		t.Fatalf("AttributesMap should have the correct value, but not ")
	}
}