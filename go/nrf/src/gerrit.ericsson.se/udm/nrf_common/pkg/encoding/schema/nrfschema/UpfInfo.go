package nrfschema

import (
	"fmt"
	"strings"

	"encoding/json"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
)

// GetInvalidInterfaceUpfInfoIndexs return invalid interfaceUpfInfoList index
func (u *TUpfInfo) GetInvalidInterfaceUpfInfoIndexs() []string {
	var invalidInterfaceUpfInfoIndexs []string
	if u.InterfaceUpfInfoList != nil {
		index := 0
		for _, item := range u.InterfaceUpfInfoList {
			if !item.IsValid() {
				invalidInterfaceUpfInfoIndex := fmt.Sprintf("%s[%d]", constvalue.InterfaceUpfInfoList, index)
				invalidInterfaceUpfInfoIndexs = append(invalidInterfaceUpfInfoIndexs, invalidInterfaceUpfInfoIndex)
			}
			index++
		}
	}

	return invalidInterfaceUpfInfoIndexs
}

func (u *TUpfInfo) createNfInfo() string {
	var sNssaiUpfInfoList string
	if u.SNssaiUpfInfoList != nil && len(u.SNssaiUpfInfoList) > 0 {
		sNssaiDnnList := ""
		for _, item := range u.SNssaiUpfInfoList {
			sst := item.SNssai.Sst
			sd := strings.ToLower(item.SNssai.Sd)
			if sd == "" {
				sd = constvalue.EmptySd
			}
			sNssaisList := fmt.Sprintf(`"sNssai":{"sst":%d,"sd":"%s"}`, sst, sd)
			dnnUpfInfoList := ""
			//dnnUpfInfoList, _ := json.Marshal(item.DnnUpfInfoList)
			for _, dnnUpfinfoItem := range item.DnnUpfInfoList {
				dnn := dnnUpfinfoItem.Dnn
				dnaiList := ""
				if dnnUpfinfoItem.DnaiList == nil {
					dnaiList = fmt.Sprintf(`["%s"]`, constvalue.EmptyDnai)
				} else {
					dnaiListByte, _ := json.Marshal(dnnUpfinfoItem.DnaiList)
					dnaiList = string(dnaiListByte)
				}
				if dnnUpfInfoList != "" {
					dnnUpfInfoList += fmt.Sprintf(`,{"dnn":"%s","dnaiList":%s}`, dnn, dnaiList)
				} else {
					dnnUpfInfoList += fmt.Sprintf(`{"dnn":"%s","dnaiList":%s}`, dnn, dnaiList)
				}

			}
			if sNssaiDnnList != "" {
				sNssaiDnnList += fmt.Sprintf(`,{%s,"dnnUpfInfoList":[%s]}`, sNssaisList, dnnUpfInfoList)
			} else {
				sNssaiDnnList += fmt.Sprintf(`{%s,"dnnUpfInfoList":[%s]}`, sNssaisList, dnnUpfInfoList)
			}
		}
		sNssaiUpfInfoList = fmt.Sprintf(`"sNssaiUpfInfoList":[%s]`, sNssaiDnnList)
	}

	var smfServingArea string
	if u.SmfServingArea != nil && len(u.SmfServingArea) > 0 {
		for _, v := range u.SmfServingArea {
			if smfServingArea == "" {
				smfServingArea = fmt.Sprintf(`"%v"`, v)
			} else {
				smfServingArea = fmt.Sprintf(`%s,"%v"`, smfServingArea, v)
			}
		}
		smfServingArea = fmt.Sprintf(`"smfServingArea":[%s]`, smfServingArea)
	}
	if smfServingArea != "" {
		return fmt.Sprintf(`"upfInfo":{%s, %s}`, sNssaiUpfInfoList, smfServingArea)
	}
	return fmt.Sprintf(`"upfInfo":{%s}`, sNssaiUpfInfoList)

}
