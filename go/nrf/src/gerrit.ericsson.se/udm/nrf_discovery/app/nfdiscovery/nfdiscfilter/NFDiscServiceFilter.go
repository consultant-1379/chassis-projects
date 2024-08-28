package nfdiscfilter

import (
	"github.com/buger/jsonparser"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"encoding/hex"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
)
//NFServiceFilter to process nfservice filter in nfprofile
type NFServiceFilter struct {
	nfServices []byte
}

func (s *NFServiceFilter)filter(nfServices []byte, queryForm *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool {

	if !s.setFilterResult(nfServices, true) {
		return false
	}

	if len(nfServices) == 0 {
		log.Debugf("nfServices not exist, isSearchNfService = %v", isSearchNfService(queryForm))
		if isSearchNfService(queryForm) {
			return s.setFilterResult([]byte(""), false)
		}
		return s.setFilterResult(s.nfServices, true)
	}

	if !s.isMatchedSupportFetature(queryForm) {
		return false
	}

	if !s.eliminatServices(queryForm) {
		return false
	}

	if !s.filterNfServicesByRequesterNfType(queryForm, filterInfo) {
		return false
	}

	if !s.filterNfServicesByRequesterPLMN(queryForm, filterInfo) {
		return false
	}

	if !s.filterNfServicesByRequesterNfFQDN(queryForm, filterInfo) {
		return false
	}
	log.Debugf("Service filter succss nfprofile: %s", string(s.nfServices))
	filterInfo.originProfile.NfServices = string(s.nfServices)

	return true

}

func (s *NFServiceFilter)setFilterResult(nfprofile []byte, ret bool) bool {
	s.nfServices = nfprofile
	return ret
}

func (s *NFServiceFilter)supportedFeatureMask(supportFeatureList []byte, featureInProfileList []byte) bool {
	i := len(supportFeatureList)
	j := len(featureInProfileList)

	for i >= 1 && j >= 1 {

		if (supportFeatureList[i - 1] & featureInProfileList[j - 1]) != 0 {
			return true
		}
		i = i - 1
		j = j - 1
	}

	return false
}

func (s *NFServiceFilter)stringToByteArray(supportFeature string) (list []byte, ret bool) {
	if len(supportFeature) % 2 != 0 {
		supportFeature = "0" + supportFeature
	}
	supportFeatureList, err := hex.DecodeString(supportFeature)
	if err != nil {
		return nil, false
	}

	return supportFeatureList, true
}

func (s *NFServiceFilter) isMatchedSupportFetature(queryForm *nfdiscrequest.DiscGetPara) bool {
	supportFeature := queryForm.GetNRFDiscSupportedFeatures()
	if supportFeature == "" {
		return s.setFilterResult(s.nfServices, true)
	}
	log.Debugf("Search nfProfile with supported-features: %s", supportFeature)
	newNfServices := ""
	supportFeatureList, ret := s.stringToByteArray(supportFeature)
	if ret == false {
		return s.setFilterResult([]byte(""), false)
	}
	_, err := jsonparser.ArrayEach(s.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		featureInProfile, err := jsonparser.GetString(value, "supportedFeatures")
		if err == nil {
			featureInProfileList, ret := s.stringToByteArray(featureInProfile)
			if ret == false {
				return
			}
			if s.supportedFeatureMask(supportFeatureList, featureInProfileList) {
				if newNfServices == "" {
					newNfServices = string(value[:])
				} else {
					newNfServices = newNfServices + "," + string(value[:])
				}
			}
		}
	})

	if err != nil || newNfServices == "" {
		return s.setFilterResult([]byte(""), false)
	}
	newNfServices = "[" + newNfServices + "]"

	return s.setFilterResult([]byte(newNfServices), true)

}

func (s *NFServiceFilter)eliminatServices(queryForm *nfdiscrequest.DiscGetPara) bool {
	serviceNameArray := queryForm.GetNRFDiscServiceName()
	log.Debugf("Enter eliminatServices")
	newNfServices := ""
	_, err := jsonparser.ArrayEach(s.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		matched := true
		//eliminat by service name
		if nil != serviceNameArray {
			name, err1 := jsonparser.GetString(value, constvalue.NFServiceName)
			if err1 != nil {
				log.Debugf("parse nfServiceName error, err=%v", err1)
			}
			for _, v := range serviceNameArray {
				if name == v {
					matched = true
					break
				}
				matched = false
			}
		}

		//eliminat by service Status
		if matched {
			status, err1 := jsonparser.GetString(value, constvalue.NFServiceStatus)
			if err1 != nil {
				log.Debugf("parse nfServiceName error, err=%v", err1)
			}
			if status != constvalue.NFServiceStatusRegistered {
				matched = false
			}
		}

		if matched {
			if newNfServices == "" {
				newNfServices = string(value[:])
			} else {
				newNfServices = newNfServices + "," + string(value[:])
			}
		}
	})
	if err != nil {
		log.Warnf("parsing arry fail for %s,error:%v", constvalue.NfServices, err)
		return s.setFilterResult([]byte(""), false)
	}

	if newNfServices == "" {
		log.Debugf("No services matched")
		return s.setFilterResult([]byte(""), false)
	}

	newNfServices = "[" + newNfServices + "]"

	return s.setFilterResult([]byte(newNfServices), true)
}

func (s *NFServiceFilter)filterNfServicesByRequesterNfType(queryForm *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool {
	requesterNfType := queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataRequesterNfType)

	log.Debugf("Enter filterNfServicesByRequesterNfType")
	newNfServices := ""
	_, err := jsonparser.ArrayEach(s.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		ok := false
		_, valueType, _, err1 := jsonparser.Get(value, constvalue.NFServiceAllowedNFTypes)
		if valueType == jsonparser.NotExist || err1 != nil {
			ok = true
		} else if valueType == jsonparser.Array {
			if requesterNfType != "" {
				ok = nfdiscutil.IsAllowedNfType(value, requesterNfType, constvalue.NFServiceAllowedNFTypes)
			} else {
				ok = false
			}
		}
		if ok {
			if newNfServices == "" {
				newNfServices = string(value[:])
			} else {
				newNfServices = newNfServices + "," + string(value[:])
			}
		}

	})
	if err != nil {
		log.Warnf("parsing arry fail for %s,error:%v", constvalue.NfServices, err)
		return s.setFilterResult([]byte(""), false)
	}

	if newNfServices == "" && isSearchNfService(queryForm) {
		log.Debugf("consumer want to search nfservice, but reuqester-nf-type is forbidden to access nfservice in nfprofile")
		filterInfo.nfTypeForbiddenInService = true
		return s.setFilterResult([]byte(""), false)
	}
	newNfServices = "[" + newNfServices + "]"

	return s.setFilterResult([]byte(newNfServices), true)
}

func (s *NFServiceFilter)filterNfServicesByRequesterPLMN(queryForm *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool {
	requesterPlmnList := queryForm.GetNRFDiscPlmnListValue(constvalue.SearchDataRequesterPlmnList)

	log.Debugf("Enter filterNfServicesByRequesterPLMN with requester-plmn-list = %v", requesterPlmnList)
	newNfServices := ""
	_, err := jsonparser.ArrayEach(s.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var ok bool
		_, _, _, err1 := jsonparser.Get(value, constvalue.NFServiceAllowedPlmns)
		if  err1 != nil {
			ok = true
		} else {
			if len(requesterPlmnList) > 0 {
				ok = isAllowedRequesterPlmn(value, requesterPlmnList, constvalue.NFServiceAllowedPlmns, []byte(filterInfo.originProfile.BodyCommon))
			} else {
				ok = false
			}
		}


		if ok {
			if newNfServices == "" {
				newNfServices = string(value[:])
			} else {
				newNfServices = newNfServices + "," + string(value[:])
			}
		}

	})
	if err != nil {
		log.Warnf("parsing arry fail for %s,error:%v", constvalue.NfServices, err)
		return s.setFilterResult([]byte(""), false)
	}

	if newNfServices == "" && isSearchNfService(queryForm) {
		log.Debugf("consumer want to search nfservice, but requester-plmn-list is forbidden to access nfservice in nfprofile")
		filterInfo.plmnForbiddenInService = true
		return s.setFilterResult([]byte(""), false)
	}
	newNfServices = "[" + newNfServices + "]"

	return s.setFilterResult([]byte(newNfServices), true)
}

func (s *NFServiceFilter)filterNfServicesByRequesterNfFQDN(queryForm *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool {
	fqdn := queryForm.GetNRFDiscRequesterNFInstFQDN()

	log.Debugf("Enter filterNfServicesByRequesterNfFQDN Parmeter Fqdn: %s", fqdn)
	newNfServices := ""
	_, err := jsonparser.ArrayEach(s.nfServices, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		ok := false
		_, valueType, _, err1 := jsonparser.Get(value, constvalue.NFServiceAllowedDomains)
		if valueType == jsonparser.NotExist || err1 != nil {
			log.Debugf("allowedDomains not exist")
			ok = true
		} else if valueType == jsonparser.Array {
			if fqdn != "" {
				ok = nfdiscutil.IsAllowedNfFQDN(value, fqdn, constvalue.NFServiceAllowedDomains)
			} else {
				ok = false
			}
		}
		if ok {
			if newNfServices == "" {
				newNfServices = string(value[:])
			} else {
				newNfServices = newNfServices + "," + string(value[:])
			}
		}

	})
	if err != nil {
		log.Warnf("parsing arry fail for %s,error:%v", constvalue.NfServices, err)
		return s.setFilterResult([]byte(""), false)
	}

	if newNfServices == "" && isSearchNfService(queryForm) {
		log.Debugf("consumer want to search nfservice, but requester-nf-instance-fqdn is forbidden to access nfservice in nfprofile")
		filterInfo.domainForbiddenInService = true
		return s.setFilterResult([]byte(""), false)
	}
	newNfServices = "[" + newNfServices + "]"

	return s.setFilterResult([]byte(newNfServices), true)
}

