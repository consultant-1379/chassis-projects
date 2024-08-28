package nfdiscfilter

import (
	"regexp"
	"strconv"
	
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/utils"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
	"github.com/buger/jsonparser"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap"
	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
)

func isMatchedLocality(nfprofile []byte, preferredLocality string) bool {
	if preferredLocality != "" {
		preferredLocalityInProfile, err := jsonparser.GetString(nfprofile, constvalue.Locality)
		if err == nil && preferredLocalityInProfile == preferredLocality {
			return true
		}
		return false
	}
	return true
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

	start, err := strconv.ParseInt(s, 10, 64)
	end, err1 := strconv.ParseInt(e, 10, 64)

	//re = regexp.MustCompile("[0-9]{5,15}")
	supiInt64, err2 := strconv.ParseInt(nfdiscutil.Compile[constvalue.SupiRanges].FindString(supi), 10, 64)
	if err != nil || err1 != nil || err2 != nil {
		log.Debugf("ParseInt error, err=%v, err1=%v, err2=%v", err, err1, err2)
	}

	if supiInt64 >= start && supiInt64 <= end {
		log.Debugf("The supi is %v, range  is %v-%v, and the matched result is true\n", supi, start, end)
		return true
	}

	log.Debugf("The supi is %v, range  is %v-%v, and the matched result is false\n", supi, start, end)

	return false
}

func isMatchedGroupID(queryForm *nfdiscrequest.DiscGetPara, groupID []string, item []byte) nfdiscutil.MatchResult {
	//targetNFInfo := map[string]string{
	//	"UDM":  "udmInfo",
	//	"AUSF": "ausfInfo",
	//	"PCF":  "pcfInfo",
	//	"UDR":  "udrInfo",
	//}
	//if "" == targetNFInfo[queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)] {
	//	return nfdiscutil.ResultError
	//}
	groupIDInProfile, err := jsonparser.GetString(item, constvalue.GroupID)
	if err == nil {
		if len(groupID) == 0 {
			return nfdiscutil.ResultFoundNotMatch
		}
		for _, v := range groupID {
			if groupIDInProfile == v {
				return nfdiscutil.ResultFoundMatch
			}
		}
	}
	if nil == err {
		return nfdiscutil.ResultFoundNotMatch
	}
	return nfdiscutil.ResultError

}

func isMatchedSupi(queryForm *nfdiscrequest.DiscGetPara, nfInfo []byte) bool {
	supi := queryForm.GetNRFDiscSupiValue()
	//targetNFInfo := map[string]string{
	//	"UDM":  "udmInfo",
	//	"AUSF": "ausfInfo",
	//	"PCF":  "pcfInfo",
	//	"UDR":  "udrInfo",
	//	"CHF":  "chfInfo",
	//}
	//
	//if "" == targetNFInfo[queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)] {
	//	return false
	//}
	supiRanges := constvalue.SupiRanges
	if "CHF" == queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType) {
		supiRanges = "supiRangeList"
	}
	ret := false
	num := 0
	_, err := jsonparser.ArrayEach(nfInfo, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
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
	}, supiRanges)
	//No supirange&groupid in NFProfile, the NFProfile is matched all
	if err == nil && ret == true {
		return true
	}
	if err != nil {
		_, err1 := jsonparser.GetString(nfInfo, constvalue.GroupID)
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

func isMatchedGpsi(queryForm *nfdiscrequest.DiscGetPara, nfInfo []byte) bool {
	ret := false
	gpsi := queryForm.GetNRFDiscGspi()
	//targetNFInfo := map[string]string{
	//	"UDM": "udmInfo",
	//	"UDR": "udrInfo",
	//	"CHF": "chfInfo",
	//}
	//
	//if "" == targetNFInfo[queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)] {
	//	return false
	//}

	gpsiRanges := constvalue.GpsiRanges
	if "CHF" == queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType) {
		gpsiRanges = "gpsiRangeList"
	}
	num := 0
	_, err := jsonparser.ArrayEach(nfInfo, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
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

		start, parseErr := strconv.ParseInt(s, 10, 64)
		end, parseErr1 := strconv.ParseInt(e, 10, 64)

		//re := regexp.MustCompile("[0-9]{5,15}")
		gpsiInt64, parseErr2 := strconv.ParseInt(nfdiscutil.Compile[constvalue.GpsiRanges].FindString(gpsi), 10, 64)
		if parseErr != nil || parseErr1 != nil || parseErr2 != nil {
			log.Debugf("parsrInt error, err=%v, err1=%v, err2=%v", parseErr, parseErr1, parseErr2)
		}

		if gpsiInt64 >= start && gpsiInt64 <= end {
			ret = true
			log.Debugf("the gpsi is %v, range is %v-%v matched result is %v", gpsi, start, end, ret)
			return
		}

		log.Debugf("the gpsi is %v, range is %v-%v matched result is %v", gpsi, start, end, ret)

	}, gpsiRanges)

	if err != nil || ret == true {
		return true
	}

	if num == 0 && err == nil {
		log.Debugf("gpsiRanges is [], this will match all gpsi")
		return true
	}
	return ret
}
func isMatchedExternalGroupID(queryForm *nfdiscrequest.DiscGetPara, nfInfo []byte) bool {
	ret := false
	//targetNFInfo := map[string]string{
	//	"UDM": "udmInfo",
	//	"UDR": "udrInfo",
	//}
	//
	//if "" == targetNFInfo[queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)] {
	//	return false
	//}
	num := 0
	externalGroupID := queryForm.GetNRFDiscExterGroupID()
	_, err1 := jsonparser.ArrayEach(nfInfo, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
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
				log.Debugf("externalGroupID regex match error, err=%v", err)
			}
			log.Debugf("The externalGroupID: %s, pattern : %s, matched result: %v", externalGroupID, pattern, matched)
			if matched {
				ret = true
				return
			}
		}

	}, constvalue.ExternalGroupIdentityfiersRanges)

	if err1 != nil || ret == true {
		return true
	}
	if num == 0 && err1 == nil {
		log.Debugf("externalGroupIdentifiersRanges is [], this will match all externalGroupIdentity")
		return true
	}
	return ret
}

func getParamSearchPath(nfType string, parameter string) configmap.SearchMapping {
	searchMapping, ok := configmap.AttributesMap[nfType][parameter]
	if !ok {
		log.Errorf("Fail to get attribute path with nf-type=%s, parameter=%s", nfType, parameter)
	}
	return searchMapping
}

//isSnssaisParaOnly is used to decide whether parameter only has snssai
func isSnssaisParaOnly(queryForm *nfdiscrequest.DiscGetPara) bool {
	if queryForm.GetExistFlag(constvalue.SearchDataSnssais) && !queryForm.GetExistFlag(constvalue.SearchDataDnn) && !queryForm.GetExistFlag(constvalue.SearchDataDnaiList) {
		return true
	}
	return false
}

//isSearchNfService is used to decide whether search nfService
func isSearchNfService(queryForm *nfdiscrequest.DiscGetPara) bool {
	if queryForm.GetExistFlag(constvalue.SearchDataServiceName) || queryForm.GetExistFlag(constvalue.SearchDataSupportedFeatures) {
		return true;
	}
	return false;
}



func isAllowedRequesterPlmn(serviceOrProfile []byte, requesterPlmnList []string, searchPath string, profileCommon []byte) bool {
	_, _, _, err  := jsonparser.Get(serviceOrProfile, searchPath)
	if err == nil {
		if nfdiscutil.IsAllowedPLMN(serviceOrProfile, requesterPlmnList, searchPath) {
			log.Debugf("reuqester-plmn-list is allowed to access nfprofile in field allowedPlmns")
			return true
		}
	} else {
		return true
	}

	_, _, _, err = jsonparser.Get(profileCommon, constvalue.PlmnList)
	if err == nil {
		if nfdiscutil.IsAllowedPLMN(profileCommon, requesterPlmnList, constvalue.PlmnList) {
			log.Debugf("reuqester-plmn-list is allowed to access nfprofile in field plmnList")
			return true
		}
	} else {
		for _, plmn := range cm.NfProfile.PlmnID {
			for _, requesterPlmn := range requesterPlmnList {
				if requesterPlmn == plmn.Mcc + plmn.Mnc {
					log.Debugf("reuqester-plmn-list is allowed to access nfprofile in cm plmnList")
					return true
				}
			}
		}
	}

	return false
}