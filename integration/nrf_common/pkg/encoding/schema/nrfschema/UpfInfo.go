package nrfschema

import (
	"fmt"
	"strings"

	"encoding/json"

	"com/dbproxy/nfmessage/subscription"

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
	} else {
		smfServingArea = fmt.Sprintf(`"smfServingArea":["%s"]`, constvalue.EmptySmfServingArea)
	}
	return fmt.Sprintf(`"upfInfo":{%s, %s}`, sNssaiUpfInfoList, smfServingArea)
}

//GenerateNfGroupCond generate NfGroupCond for subscription
func (u *TUpfInfo) GenerateNfGroupCond() *subscription.SubKeyStruct {
	return nil
}

// IsEqual is to check if NFInfo is equal
func (u *TUpfInfo) IsEqual(c TNfInfo) bool {

	b := c.(*TUpfInfo)

	if (u.IwkEpsInd != nil && b.IwkEpsInd == nil) || (u.IwkEpsInd == nil && b.IwkEpsInd != nil) {
		return false
	}

	if u.IwkEpsInd != nil && b.IwkEpsInd != nil {
		if *(u.IwkEpsInd) != *(b.IwkEpsInd) {
			return false
		}
	}

	if len(u.SmfServingArea) != len(b.SmfServingArea) {
		return false
	}

	if len(u.InterfaceUpfInfoList) != len(b.InterfaceUpfInfoList) {
		return false
	}

	if len(u.SNssaiUpfInfoList) != len(b.SNssaiUpfInfoList) {
		return false
	}

	for k, item := range u.SmfServingArea {
		if item != b.SmfServingArea[k] {
			return false
		}
	}

	for k, item := range u.InterfaceUpfInfoList {
		bb := b.InterfaceUpfInfoList[k]
		if item.EndpointFqdn != bb.EndpointFqdn || item.InterfaceType != bb.InterfaceType || item.NetworkInstance != bb.NetworkInstance {
			return false
		}

		if len(item.Ipv4EndpointAddresses) != len(bb.Ipv4EndpointAddresses) {
			return false
		}

		if len(item.Ipv6EndpointAddresses) != len(bb.Ipv6EndpointAddresses) {
			return false
		}

		for j, ipv4 := range item.Ipv4EndpointAddresses {
			if ipv4 != bb.Ipv4EndpointAddresses[j] {
				return false
			}
		}

		for j, ipv6 := range item.Ipv6EndpointAddresses {
			if ipv6 != bb.Ipv6EndpointAddresses[j] {
				return false
			}
		}
	}

	for k, item := range u.SNssaiUpfInfoList {
		bb := b.SNssaiUpfInfoList[k]
		if item.SNssai.Sst != bb.SNssai.Sst || item.SNssai.Sd != bb.SNssai.Sd {
			return false
		}

		if len(item.DnnUpfInfoList) != len(bb.DnnUpfInfoList) {
			return false
		}

		for j, inItem := range item.DnnUpfInfoList {
			bbb := bb.DnnUpfInfoList[j]
			if inItem.Dnn != bbb.Dnn {
				return false
			}

			if len(inItem.DnaiList) != len(bbb.DnaiList) {
				return false
			}
			for i, dnai := range inItem.DnaiList {
				if dnai != bbb.DnaiList[i] {
					return false
				}
			}
		}
	}

	return true
}
