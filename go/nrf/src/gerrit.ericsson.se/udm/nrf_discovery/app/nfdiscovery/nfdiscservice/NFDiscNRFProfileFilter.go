package nfdiscservice

import (
	"com/dbproxy/nfmessage/nrfprofile"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/utils"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
	"github.com/buger/jsonparser"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func setAMFFilterForNRFProfile(queryForm nfdiscrequest.DiscGetPara, nrfProfileIndex *nrfprofile.NRFProfileIndex) {

	plmnid, amfid := queryForm.GetNRFDiscGuamiType()
	if len(plmnid) == 5 {
		plmnArray := []rune(plmnid)
		plmnid = string(plmnArray[0:3]) + "0" + string(plmnArray[3:])
	}
	if plmnid != "" && amfid != "" {
		var subKey nrfprofile.NRFKeyStruct
		subKey.SubKey1 = plmnid
		subKey.SubKey2 = amfid
		var subKeyList []*nrfprofile.NRFKeyStruct
		subKeyList = append(subKeyList, &subKey)
		nrfProfileIndex.AmfKey1 = subKeyList
	}

	/*plmnid, tac := queryForm.GetNRFDiscTaiType()
	if plmnid != "" && tac != "" {
		var subKey nrfprofile.NRFKeyStruct
		subKey.SubKey1 = plmnid
		subKey.SubKey2 = tac
		var subKeyList []*nrfprofile.NRFKeyStruct
		subKeyList = append(subKeyList, &subKey)
		nrfProfileIndex.AmfKey2 = subKeyList
	}*/

	amfSetID := queryForm.GetNRFDiscAMFSetID()
	if "" != amfSetID {
		var subKey nrfprofile.NRFKeyStruct
		subKey.SubKey1 = amfSetID
		var subKeyList []*nrfprofile.NRFKeyStruct
		subKeyList = append(subKeyList, &subKey)
		nrfProfileIndex.AmfKey4 = subKeyList
	}

	amfRegionID := queryForm.GetNRFDiscAMFRegionID()
	if "" != amfRegionID {
		var subKey nrfprofile.NRFKeyStruct
		subKey.SubKey1 = amfRegionID
		var subKeyList []*nrfprofile.NRFKeyStruct
		subKeyList = append(subKeyList, &subKey)
		nrfProfileIndex.AmfKey3 = subKeyList
	}

}

func setSMFFilterForNRFProfile(queryForm nfdiscrequest.DiscGetPara, nrfProfileIndex *nrfprofile.NRFProfileIndex) {
	pgw := queryForm.GetNRFDiscPGW()
	if "" != pgw {
		var subKey nrfprofile.NRFKeyStruct
		subKey.SubKey1 = pgw
		var subKeyList []*nrfprofile.NRFKeyStruct
		subKeyList = append(subKeyList, &subKey)
		nrfProfileIndex.SmfKey2 = subKeyList
	}

	dnn := queryForm.GetNRFDiscDnnValue()
	if "" != dnn {
		var subKey nrfprofile.NRFKeyStruct
		subKey.SubKey1 = dnn
		var subKeyList []*nrfprofile.NRFKeyStruct
		subKeyList = append(subKeyList, &subKey)
		nrfProfileIndex.SmfKey1 = subKeyList
	}

	/*plmnid, tac := queryForm.GetNRFDiscTaiType()
	if plmnid != "" && tac != "" {
		var subKey nrfprofile.NRFKeyStruct
		subKey.SubKey1 = plmnid
		subKey.SubKey2 = tac
		var subKeyList []*nrfprofile.NRFKeyStruct
		subKeyList = append(subKeyList, &subKey)
		nrfProfileIndex.SmfKey3 = subKeyList
	}*/
}

func setUDMFilterForNRFProfile(queryForm nfdiscrequest.DiscGetPara, nrfProfileIndex *nrfprofile.NRFProfileIndex) {
	routingIndicator := queryForm.GetNRFDiscRoutingIndicator()
	if "" != routingIndicator {
		var subKey nrfprofile.NRFKeyStruct
		subKey.SubKey1 = routingIndicator
		var subKeyList []*nrfprofile.NRFKeyStruct
		subKeyList = append(subKeyList, &subKey)
		nrfProfileIndex.UdmKey2 = subKeyList
	}
	/*
		supi := queryForm.GetNRFDiscSupiValue()
		if supi != "" {
			groupid := getGroupIDfromDB(queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), supi)
			if nil != queryForm.GetNRFDiscGroupIDList() {
				groupid = append(groupid, queryForm.GetNRFDiscGroupIDList()...)
			}
			if len(groupid) > 0 {
				var subKeyList []*nrfprofile.NRFKeyStruct
				for _, v := range groupid {
					var subKey nrfprofile.NRFKeyStruct
					subKey.SubKey1 = v
					subKeyList = append(subKeyList, &subKey)
				}

				nrfProfileIndex.UdmKey1 = subKeyList
			}
		}*/

	gpsi := queryForm.GetNRFDiscGspi()
	if gpsi != "" {
		groupid, _ := nfdiscutil.GetGpsiGroupIDfromDB(queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), gpsi)
		if len(groupid) > 0 {
			var subKeyList []*nrfprofile.NRFKeyStruct
			for _, v := range groupid {
				var subKey nrfprofile.NRFKeyStruct
				subKey.SubKey1 = v
				subKeyList = append(subKeyList, &subKey)
			}

			nrfProfileIndex.UdmKey1 = subKeyList
		}
	}
}

func setAUSFFilterForNRFProfile(queryForm nfdiscrequest.DiscGetPara, nrfProfileIndex *nrfprofile.NRFProfileIndex) {
	routingIndicator := queryForm.GetNRFDiscRoutingIndicator()
	if "" != routingIndicator {
		var subKey nrfprofile.NRFKeyStruct
		subKey.SubKey1 = routingIndicator
		var subKeyList []*nrfprofile.NRFKeyStruct
		subKeyList = append(subKeyList, &subKey)
		nrfProfileIndex.AusfKey2 = subKeyList
	}
	/*
		supi := queryForm.GetNRFDiscSupiValue()
		if supi != "" {
			groupid := getGroupIDfromDB(queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), supi)
			if nil != queryForm.GetNRFDiscGroupIDList() {
				groupid = append(groupid, queryForm.GetNRFDiscGroupIDList()...)
			}
			if len(groupid) > 0 {
				var subKeyList []*nrfprofile.NRFKeyStruct
				for _, v := range groupid {
					var subKey nrfprofile.NRFKeyStruct
					subKey.SubKey1 = v
					subKeyList = append(subKeyList, &subKey)
				}

				nrfProfileIndex.AusfKey1 = subKeyList
			}
		}*/
}

func setPCFFilterForNRFProfile(queryForm nfdiscrequest.DiscGetPara, nrfProfileIndex *nrfprofile.NRFProfileIndex) {
	dnn := queryForm.GetNRFDiscDnnValue()
	if "" != dnn {
		var subKey nrfprofile.NRFKeyStruct
		subKey.SubKey1 = dnn
		var subKeyList []*nrfprofile.NRFKeyStruct
		subKeyList = append(subKeyList, &subKey)
		nrfProfileIndex.PcfKey1 = subKeyList
	}
	/*
		supi := queryForm.GetNRFDiscSupiValue()
		if supi != "" {
			groupid := getGroupIDfromDB(queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), supi)
			if nil != queryForm.GetNRFDiscGroupIDList() {
				groupid = append(groupid, queryForm.GetNRFDiscGroupIDList()...)
			}
			if len(groupid) > 0 {
				var subKeyList []*nrfprofile.NRFKeyStruct
				for _, v := range groupid {
					var subKey nrfprofile.NRFKeyStruct
					subKey.SubKey1 = v
					subKeyList = append(subKeyList, &subKey)
				}

				nrfProfileIndex.PcfKey2 = subKeyList
			}
		}*/
}

func plmnDiscGRPCGetRequestFilter(queryForm nfdiscrequest.DiscGetPara, nrfProfileGetRequest *nrfprofile.NRFProfileGetRequest) {
	var nrfProfileIndex nrfprofile.NRFProfileIndex
	switch queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType) {
	case "AMF":
		setAMFFilterForNRFProfile(queryForm, &nrfProfileIndex)
	case "SMF":
		setSMFFilterForNRFProfile(queryForm, &nrfProfileIndex)
	case "UDM":
		setUDMFilterForNRFProfile(queryForm, &nrfProfileIndex)
	case "AUSF":
		setAUSFFilterForNRFProfile(queryForm, &nrfProfileIndex)
	case "PCF":
		setPCFFilterForNRFProfile(queryForm, &nrfProfileIndex)
	default:

	}

	nrfProfileIndex.Key1 = uint64(time.Now().Unix()) * 1000
	nrfProfileIndex.Key2 = math.MaxInt64

	nrfProfileFilter := &nrfprofile.NRFProfileFilter{
		AndOperation: true,
		Index:        &nrfProfileIndex,
	}

	nrfProfileFilterData := &nrfprofile.NRFProfileGetRequest_Filter{
		Filter: nrfProfileFilter,
	}
	nrfProfileGetRequest.Data = nrfProfileFilterData

}
func matchByPatternForSupi(supi string, value []byte) bool {
	pattern, err := jsonparser.GetString(value, "pattern")
	if err == nil {
		matched, err2 := regexp.MatchString(pattern, supi)
		if err2 != nil {
			log.Debugf("supi regex match error, err=%v", err2)
		}
		log.Debugf("The supi is %s, pattern regex is %s, and the matched result is %v\n", supi, pattern, matched)
		if matched {
			return true
		}
	}
	return false
}

func matchbyStartEndForSupi(supi string, value []byte) bool {
	s, sOk := jsonparser.GetString(value, "start")
	e, endOk := jsonparser.GetString(value, "end")
	if sOk != nil || endOk != nil {
		return false
	}
	//supi format imsi-[0-9]{5,15}|nai-.+|suci-[0-9]{5-15}
	if !utils.IsDigit(s) || !utils.IsDigit(e) {
		log.Errorf("The start or end of supi range is not digit")
		return false
	}
	//re := regexp.MustCompile("imsi-[0-9]{5,15}|suci-[0-9]{5,15}")
	//matched := re.MatchString(supi)
	matched := nfdiscutil.Compile[constvalue.SupiFormat].MatchString(supi)
	if !matched {
		log.Debugf("supi %s is not imsi or suci format", supi)
		return false
	}

	start, err1 := strconv.ParseInt(s, 10, 64)
	end, err2 := strconv.ParseInt(e, 10, 64)

	//re = regexp.MustCompile("[0-9]{5,15}")
	supiInt64, err3 := strconv.ParseInt(nfdiscutil.Compile[constvalue.SupiRanges].FindString(supi), 10, 64)
	if err1 != nil || err2 != nil || err3 != nil {
		log.Debugf("supi parse int error, err1=%v,err2=%v,err3=%v", err1, err2, err3)
	}

	if supiInt64 >= start && supiInt64 <= end {
		log.Debugf("The supi is %v, range  is %v-%v, and the matched result is true\n", supi, start, end)
		return true
	}
	log.Debugf("The supi is %v, range  is %v-%v, and the matched result is false\n", supi, start, end)

	return false
}

func isMatchedSupiForNRFProfile(queryForm *nfdiscrequest.DiscGetPara, nrfInfo []byte) bool {
	supi := queryForm.GetNRFDiscSupiValue()
	targetNFInfo := map[string]string{
		"UDM":  "udmInfoSum",
		"AUSF": "ausfInfoSum",
		"PCF":  "pcfInfoSum",
	}

	if "" == targetNFInfo[queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)] {
		return false
	}

	ret := false
	num := 0
	_, err := jsonparser.ArrayEach(nrfInfo, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		num = num + 1
		if ret == true {
			return
		}
		_, err1 := jsonparser.GetString(value, "pattern")
		_, err2 := jsonparser.GetString(value, "start")
		_, err3 := jsonparser.GetString(value, "end")
		if err1 != nil && err2 != nil && err3 != nil {
			log.Debugf("supiranges not have start & end & pattern, match each supi")
			ret = true
			return
		}
		//match with pattern of supirange
		if matchByPatternForSupi(supi, value) {
			ret = true
			return
		}

		//match with start and end of supirange
		if matchbyStartEndForSupi(supi, value) {
			ret = true
			return
		}
	}, constvalue.NrfInfo, targetNFInfo[queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)], constvalue.SupiRanges)
	//No supirange in NFProfile, the NFProfile is matched all
	if err == nil && ret == true {
		return true
	}

	if err != nil {
		_, err1 := jsonparser.GetString(nrfInfo, constvalue.NrfInfo, targetNFInfo[queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)], constvalue.GroupIDList)
		if err1 != nil {
			log.Debugf("supiRanges&groupId both not exist, match all supi")
			return true
		}
		log.Debugf("supiRanges not exist, but groupid exist, not match all supi")
		return false
	}

	if num == 0 && err == nil {
		log.Debugf("supiRanges is [], this will match all supi")
		return true
	}
	return ret
}

func isMatchedGpsiForNRFPRofile(queryForm *nfdiscrequest.DiscGetPara, nrfInfo []byte) bool {
	ret := false
	gpsi := queryForm.GetNRFDiscGspi()
	targetNFInfo := map[string]string{
		"UDM": "udmInfoSum",
	}

	if "" == targetNFInfo[queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)] {
		return false
	}
	num := 0
	_, err := jsonparser.ArrayEach(nrfInfo, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		num = num + 1
		if ret {
			return
		}
		pattern, err1 := jsonparser.GetString(value, "pattern")
		s, err2 := jsonparser.GetString(value, "start")
		e, err3 := jsonparser.GetString(value, "end")
		if err1 != nil && err2 != nil && err3 != nil {
			log.Debugf("gpsiranges not have start & end & pattern, match each supi")
			ret = true
			return
		}

		if err1 == nil {
			matched, err := regexp.MatchString(pattern, gpsi)
			if err != nil {
				log.Debugf("gpsi regex match error, err=%v", err)
			}
			log.Debugf("The gpsi: %s, pattern : %s, matched result: %v", gpsi, pattern, matched)
			if matched {
				ret = true
				return
			}
		}

		if err2 != nil || err3 != nil {
			return
		}

		if !utils.IsDigit(s) || !utils.IsDigit(e) {
			log.Errorf("The start or end of gpsiranges is not digit")
		}

		start, err1 := strconv.ParseInt(s, 10, 64)
		end, err2 := strconv.ParseInt(e, 10, 64)

		//re := regexp.MustCompile("[0-9]{5,15}")
		gpsiInt64, err3 := strconv.ParseInt(nfdiscutil.Compile[constvalue.GpsiRanges].FindString(gpsi), 10, 64)
		if err1 != nil || err2 != nil || err3 != nil {
			log.Debugf("gpsi parse int error, err1=%v,err2=%v,err3=%v", err1, err2, err3)
		}
		if gpsiInt64 >= start && gpsiInt64 <= end {
			ret = true
			log.Debugf("the gpsi is %v, range is %v-%v matched result is %v", gpsi, start, end, ret)
			return
		}

		log.Debugf("the gpsi is %v, range is %v-%v matched result is %v", gpsi, start, end, ret)

	}, constvalue.NrfInfo, targetNFInfo[queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)], constvalue.GpsiRanges)

	if err != nil || ret == true {
		return true
	}

	if num == 0 && err == nil {
		log.Debugf("gpsiRanges is [], this will match all gpsi")
		return true
	}
	return ret
}

func isMatchedExternalGroupIDForNRFProfile(queryForm *nfdiscrequest.DiscGetPara, nrfInfo []byte) bool {
	ret := false
	targetNFInfo := map[string]string{
		"UDM": "udmInfoSum",
	}

	if "" == targetNFInfo[queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)] {
		return false
	}
	num := 0
	externalGroupID := queryForm.GetNRFDiscExterGroupID()
	_, err1 := jsonparser.ArrayEach(nrfInfo, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		num = num + 1
		if ret {
			return
		}
		pattern, err1 := jsonparser.GetString(value, "pattern")
		_, err2 := jsonparser.GetString(value, "start")
		_, err3 := jsonparser.GetString(value, "end")
		if err1 != nil && err2 != nil && err3 != nil {
			log.Debugf("externalidentityranges not have start & end & pattern, match each supi")
			ret = true
			return
		}

		if err1 == nil {
			matched, err := regexp.MatchString(pattern, externalGroupID)
			if err != nil {
				log.Errorf("externalGroupID regex match error, err=%v", err)
			}
			log.Debugf("The externalGroupID: %s, pattern : %s, matched result: %v", externalGroupID, pattern, matched)
			if matched {
				ret = true
				return
			}
		}
		//TODO when externalGroupIdentityfiersRanges change to be externalGroupIdentifiersRanges in nrfprofile, below need modify
	}, constvalue.NrfInfo, targetNFInfo[queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)], "externalGroupIdentityfiersRanges") //constvalue.ExternalGroupIdentityfiersRanges)

	if err1 != nil || ret == true {
		return true
	}
	if num == 0 && err1 == nil {
		log.Debugf("externalGroupIdentifiersRanges is [], this will match all externalGroupIdentity")
		return true
	}
	return ret
}

func getRegionNRFAddrFromProfile(nrfInfo []byte) map[string]string {
	nrfaddr := make(map[string]string)
	_, err := jsonparser.ArrayEach(nrfInfo, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		serviceName, err := jsonparser.GetString(value, constvalue.NFServiceName)
		if serviceName != constvalue.NNRFDISC || err != nil {
			return
		}
		schema, err := jsonparser.GetString(value, constvalue.NFServiceScheme)
		if err != nil {
			return
		}
		if !strings.Contains(schema, "://") {
			schema += "://"
		}

		fqdn, err := jsonparser.GetString(value, constvalue.NFServiceFqdn)
		if err != nil {
			fqdn = ""
			log.Warnf("Not include fqdn in NRFProfile.nfService")
		}

		priority, err := jsonparser.GetInt(value, constvalue.NFServicePriority)
		if err != nil {
			priority = math.MaxInt64
		}
		apiVersionInURIList := make(map[string]string)
		_, err1 := jsonparser.ArrayEach(value, func(value1 []byte, dataType jsonparser.ValueType, offset int, err error) {
			apiVersionInURI, err := jsonparser.GetString(value1, constvalue.ApiVersionInUri)
			if err != nil {
				return
			}
			apiVersionInURIList[apiVersionInURI] = ""
		}, constvalue.NFServiceVersions)
		if err1 != nil {
			return
		}

		_, err = jsonparser.ArrayEach(value, func(value1 []byte, dataType jsonparser.ValueType, offset int, err error) {
			port, err := jsonparser.GetInt(value1, constvalue.IPEndPointPort)
			if err != nil {
				return
			}

			v4addr, err1 := jsonparser.GetString(value1, constvalue.IPEndPointIpv4Address)
			if err1 == nil {
				for version := range apiVersionInURIList {
					nrfaddr[schema+v4addr+":"+strconv.Itoa(int(port))+"/"+constvalue.DiscoveryServiceName+"/"+version] = strconv.FormatInt(priority, 10) + ",1"
				}
			}

			v6addr, err2 := jsonparser.GetString(value1, constvalue.IPEndPointIpv6Address)
			if err2 == nil {
				for version := range apiVersionInURIList {
					nrfaddr[schema+"["+v6addr+"]"+":"+strconv.Itoa(int(port))+"/"+constvalue.DiscoveryServiceName+"/"+version] = strconv.FormatInt(priority, 10) + ",2"
				}
			}
                        if fqdn != "" {
				for version := range apiVersionInURIList {
					nrfaddr[schema + fqdn + ":" + strconv.Itoa(int(port)) + "/" + constvalue.DiscoveryServiceName + "/" + version] = strconv.FormatInt(priority, 10) + ",3"
				}
			}

		}, constvalue.NFServiceIPEndPoints)
		if err != nil {
			log.Errorf("jsonparser parser nrfprofile ipEndPoints fail")
			if fqdn != "" {
				if strings.Contains(schema, "https") {
					for version := range apiVersionInURIList {
						nrfaddr[schema + fqdn + ":443" + "/" + constvalue.DiscoveryServiceName + "/" + version] = strconv.FormatInt(priority, 10) + ",4"
					}
				} else {
					for version := range apiVersionInURIList {
						nrfaddr[schema + fqdn + ":80" + "/" + constvalue.DiscoveryServiceName + "/" + version] = strconv.FormatInt(priority, 10) + ",4"
					}
				}
			}

		}
	}, constvalue.NfServices)

	if err != nil {
		log.Errorf("PlmnNRF get RegionNRF address fail from nrfprofile")
	}

	return nrfaddr
}

func isMatchedTaiList(nfProfile []byte, plmnID string, tac string, nftype string) bool {
	matched := false
	targetNFInfo := map[string]string{
		"AMF": "amfInfoSum",
		"SMF": "smfInfoSum",
	}

	if "" == targetNFInfo[nftype] {
		return false
	}
	_, err := jsonparser.ArrayEach(nfProfile, func(value []byte, dataType jsonparser.ValueType, offset int, err1 error) {
		if matched {
			return
		}
		mccInProfile := ""
		mncInProfile := ""
		tacInProfile := ""
		err2 := jsonparser.ObjectEach(value, func(key1 []byte, value1 []byte, dataType jsonparser.ValueType, offset int) error {
			if string(key1) == constvalue.SearchDataPlmnID {
				err3 := jsonparser.ObjectEach(value1, func(key2 []byte, value2 []byte, dataType jsonparser.ValueType, offset int) error {
					if string(key2) == constvalue.Mcc {
						mccInProfile = string(value2)
					}

					if string(key2) == constvalue.Mnc {
						mncInProfile = string(value2)
					}
					return nil
				})

				if err3 != nil {
					return err3
				}
			}

			if string(key1) == constvalue.SearchDataTac {
				tacInProfile = string(value1)
			}

			return nil
		})

		if err2 == nil && plmnID == (mccInProfile+mncInProfile) && tac == tacInProfile {
			matched = true
		}

	}, constvalue.NrfInfo, targetNFInfo[nftype], constvalue.TaiList)

	if err != nil {
		matched = false
	}
	if matched {
		return matched
	}
	return matched
}

func isMatchedTacRangeList(tacRange []byte, tac string) bool {
	ret := false
	_, err := jsonparser.ArrayEach(tacRange, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if ret {
			return
		}
		pattern, err := jsonparser.GetString(value, "pattern")
		if err == nil {
			matched, regexErr := regexp.MatchString(pattern, tac)
			if regexErr != nil {
				log.Errorf("tac regex match error,err=%v", regexErr)
			}
			if matched {
				ret = true
			}
			return
		}

		start, err := jsonparser.GetString(value, "start")
		if err != nil {
			ret = false
			return
		}

		end, err := jsonparser.GetString(value, "end")
		if err != nil {
			ret = false
			return
		}

		if len(tac) != len(start) || len(start) != len(end) {
			ret = false
			return
		}

		if strings.Compare(tac, start) >= 0 && strings.Compare(end, tac) >= 0 {
			ret = true
			return
		}
	})
	if err != nil {
		return false
	}
	return ret
}

func isMatchedTaiRangeList(nfProfile []byte, plmnID string, tac string, nftype string) bool {
	matched := false
	targetNFInfo := map[string]string{
		"AMF": "amfInfoSum",
		"SMF": "smfInfoSum",
	}

	if "" == targetNFInfo[nftype] {
		return false
	}
	_, err := jsonparser.ArrayEach(nfProfile, func(value []byte, dataType jsonparser.ValueType, offset int, err1 error) {
		if matched {
			return
		}
		mccInProfile := ""
		mncInProfile := ""
		var tacRangeInProfile []byte
		err2 := jsonparser.ObjectEach(value, func(key1 []byte, value1 []byte, dataType jsonparser.ValueType, offset int) error {
			if string(key1) == constvalue.SearchDataPlmnID {
				err3 := jsonparser.ObjectEach(value1, func(key2 []byte, value2 []byte, dataType jsonparser.ValueType, offset int) error {
					if string(key2) == constvalue.Mcc {
						mccInProfile = string(value2)
					}

					if string(key2) == constvalue.Mnc {
						mncInProfile = string(value2)
					}
					return nil
				})

				if err3 != nil {
					return err3
				}
			}

			if string(key1) == constvalue.TacRangeList {
				tacRangeInProfile = value1
			}

			return nil
		})

		if err2 == nil && plmnID == (mccInProfile+mncInProfile) && isMatchedTacRangeList(tacRangeInProfile, tac) {
			matched = true
		}

	}, constvalue.NrfInfo, targetNFInfo[nftype], constvalue.TaiRangeList)

	if err != nil {
		matched = false
	}
	if matched {
		return matched
	}
	return matched
}

func isTaiMatchedException(item []byte, nftype string) bool {
	targetNFInfo := map[string]string{
		"AMF": "amfInfoSum",
		"SMF": "smfInfoSum",
	}
	if "" == targetNFInfo[nftype] {
		return false
	}

	_, datatype1, _, err := jsonparser.Get(item, constvalue.NrfInfo, targetNFInfo[nftype], constvalue.TaiList)
	_, datatype2, _, err1 := jsonparser.Get(item, constvalue.NrfInfo, targetNFInfo[nftype], constvalue.TaiRangeList)
	if datatype1 == jsonparser.NotExist && datatype2 == jsonparser.NotExist && err == nil && err1 == nil {
		return true
	}

	return false
}

func isMatchedGroupIDForNRFProfile(queryForm *nfdiscrequest.DiscGetPara, groupID []string, item []byte) nfdiscutil.MatchResult {
	targetNRFInfo := map[string]string{
		"UDM":  "udmInfoSum",
		"AUSF": "ausfInfoSum",
		"PCF":  "pcfInfoSum",
	}
	if "" == targetNRFInfo[queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)] {
		return nfdiscutil.ResultError
	}
	matched := false
	_, err := jsonparser.ArrayEach(item, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if matched {
			return
		}
		groupIDInProfile := string(value[:])
		for _, v := range groupID {
			if groupIDInProfile == v {
				matched = true
				return
			}
		}
	}, constvalue.NrfInfo, targetNRFInfo[queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)], constvalue.GroupIDList)

	if err != nil {
		return nfdiscutil.ResultFoundNotMatch
	}
	if matched {
		return nfdiscutil.ResultFoundMatch
	}
	return nfdiscutil.ResultFoundNotMatch
}

func plmnDiscNRFProfileFilter(nrfresp *nrfprofile.NRFProfileGetResponse, queryForm nfdiscrequest.DiscGetPara) ([]string, map[string]string){
	log.Debugf("Filter nrfprofile")
	nrfProfilesInfo := nrfresp.GetNrfProfile()
	var nrfProfiles [][]byte
	for _, item := range nrfProfilesInfo {
		nrfProfiles = append(nrfProfiles, item.RawNrfProfile)
	}

	var nrfAddrList []string
	nrfAddrPriorityMap := make(map[string]string)
	nrfAddrInstanceIDMap := make(map[string]string)
	var groupIDList []string
	if queryForm.GetNRFDiscSupiValue() != "" {
		groupIDList, _ = nfdiscutil.GetGroupIDfromDB(queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType), queryForm.GetNRFDiscSupiValue())
	}

	if queryForm.GetNRFDiscGroupIDList() != nil {
		groupIDList = append(groupIDList, queryForm.GetNRFDiscGroupIDList()...)
	}

	for _, item := range nrfProfiles {
		log.Debugf("nrfprofile: %s", string(item))
		nfStatus, err := jsonparser.GetString(item, constvalue.NfStatus)
		if err != nil || nfStatus != constvalue.NFStatusRegistered {
			log.Debugf("Get nfStatus from nrfprofile fail or nfStatus not REGISTERED")
			continue
		}
		if queryForm.GetNRFDiscSupiValue() != "" {
			log.Debugf("Search nrfProfile with supi %s", queryForm.GetNRFDiscSupiValue())
			matchResult := isMatchedGroupIDForNRFProfile(&queryForm, groupIDList, item)
			if (matchResult == nfdiscutil.ResultFoundMatch) || isMatchedSupiForNRFProfile(&queryForm, item) {
				log.Debugf("Matched nrfProfile is Found with supi %s", queryForm.GetNRFDiscSupiValue())
			} else {
				log.Debugf("No Matched nrfProfile with supi %s", queryForm.GetNRFDiscSupiValue())
				continue
			}
		}

		if queryForm.GetNRFDiscGroupIDList() != nil {
			log.Debugf("Search nrfProfile with groupidlist")
			matchedResult := isMatchedGroupIDForNRFProfile(&queryForm, groupIDList, item)
			if matchedResult == nfdiscutil.ResultFoundMatch {
				log.Debugf("Matched nrfProfile is Found with groupid list")
			} else {
				log.Debugf("No Matched nrfProfile with groupid list")
				continue
			}
		}

		if queryForm.GetNRFDiscGspi() != "" {
			if !isMatchedGpsiForNRFPRofile(&queryForm, item) {
				log.Debugf("No matched nrfprofile with gpsi: %s", queryForm.GetNRFDiscGspi())
				continue
			} else {
				log.Debugf("Matched nrfProfile with supi : %s", queryForm.GetNRFDiscGspi())
			}
		}

		if queryForm.GetNRFDiscExterGroupID() != "" {
			if !isMatchedExternalGroupIDForNRFProfile(&queryForm, item) {
				log.Debugf("No matched nrfprofile with externalGroupID :%s", queryForm.GetNRFDiscExterGroupID())
				continue
			} else {
				log.Debugf("Matched nrfprofile with externalGroupID : %s", queryForm.GetNRFDiscExterGroupID())
			}
		}

		plmnid, tac := queryForm.GetNRFDiscTaiType()
		if plmnid != "" && tac != "" {
			if isTaiMatchedException(item, queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)) || isMatchedTaiList(item, plmnid, tac, queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)) || isMatchedTaiRangeList(item, plmnid, tac, queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)) {
				log.Debugf("Matched nrfprofile with tac plmnid: %s, tac: %s", plmnid, tac)
			} else {
				log.Debugf("No matched nrfprofile with tai plmnid: %s, tac: %s", plmnid, tac)
				continue
			}
		}

		instanceid, err := jsonparser.GetString(item, constvalue.NfInstanceId);
                if err != nil {
			log.Warnf("Get NRFProfile instanceID fail")
		}
		nrfAddrPriority := getRegionNRFAddrFromProfile(item)
		for k, v := range nrfAddrPriority {
			nrfAddrPriorityMap[k] = v
			nrfAddrInstanceIDMap[k] = instanceid
		}
	}
	sortAddr := sortMapByValue(nrfAddrPriorityMap)
	for _, v := range sortAddr {
		nrfAddrList = append(nrfAddrList, v.Key)
	}
	return nrfAddrList, nrfAddrInstanceIDMap
}
