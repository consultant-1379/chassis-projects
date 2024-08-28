package nfdiscfilter

import (
	"fmt"
	"github.com/buger/jsonparser"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
)

//NFPostFilter to do action after filter one nfprofile success
type NFPostFilter struct {
	nfProfile  string
	bodyCommon string
	nfInfo     string
	nfServices string
}

func (p * NFPostFilter) filter(nfProfile []byte, queryForm *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool {
	if len(filterInfo.originProfile.NfServices) > 0 {
		p.nfServices = filterInfo.originProfile.NfServices
	} else {
		p.nfServices = ""
	}

	if len(filterInfo.originProfile.NfInfo) > 0 {
		p.nfInfo = filterInfo.originProfile.NfInfo
	} else {
		p.nfInfo = ""
	}

	p.bodyCommon = string(nfProfile)

	p.nfProfile = string(nfProfile)
	p.replaceFQDNAttr(queryForm)

	if !p.eliminateFilelds(queryForm){
		return false
	}

	if !p.eliminatn2INterfaceAmfInfo(queryForm){
		return false
	}

	guamiMatched := p.isMatchedAMFGuamiList(queryForm)
	localityMatched := isMatchedLocality([]byte(p.nfProfile), queryForm.GetNRFDiscPreferredLocality())

	log.Debugf("eliminated nfProfile=%v", p.nfProfile)
	if len(p.nfServices) > 0 {
		log.Debugf("eliminated nfServices=%v", p.nfServices)
		profile, err1 := jsonparser.Set([]byte(p.nfProfile), []byte(p.nfServices), constvalue.NfServices)
		if err1 != nil {
			log.Warnf("set value fail for %s,error:%v", constvalue.NfServices, err1)
			return false
		}
		p.nfProfile = string(profile)
	}

	if len(p.nfInfo) > 0 {
		log.Debugf("eliminated nfInfo=%v", p.nfInfo)
		if "" == constvalue.TargetNFInfo[queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)] {
			log.Warnf("target-nf-type= %s cannot find right nfinfo", queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType))
			return false
		}
		profile, err2 := jsonparser.Set([]byte(p.nfProfile), []byte(p.nfInfo), constvalue.TargetNFInfo[queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType)])

		if err2 != nil {
			log.Warnf("set value fail for nfInfo,error:%v", err2)
			return false
		}
		p.nfProfile = string(profile)
	}
	if localityMatched && guamiMatched {
		//master
		if filterInfo.newProfiles == "" {
			filterInfo.newProfiles = p.nfProfile
		} else {
			filterInfo.newProfiles = filterInfo.newProfiles + "," + p.nfProfile
		}
		if queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType) != "NRF" {
			p.getEtagValue(filterInfo, true)
		}
	} else {
		//backup
		if filterInfo.backupNewProfiles == "" {
			filterInfo.backupNewProfiles = p.nfProfile
		} else {
			filterInfo.backupNewProfiles = filterInfo.backupNewProfiles + "," + p.nfProfile
		}
		if queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType) != "NRF" {
			p.getEtagValue(filterInfo, false)
		}
	}

	return true

}


func (p *NFPostFilter) getEtagValue(filterInfo *FilterInfo, master bool) {
	profileInstID := filterInfo.originProfile.Key
	md5Sum := nfdiscutil.GetNFProfileMD5Sum([]byte(filterInfo.originProfile.MD5Sum), []byte(filterInfo.originProfile.NfServices))
	if master {
		if "" != md5Sum && "" != profileInstID {
			filterInfo.nfProfilesMd5Sum[profileInstID] = md5Sum
			filterInfo.etagKeys = append(filterInfo.etagKeys, profileInstID)
		} else {
			filterInfo.etagExist = false
		}
	} else {
		if "" != md5Sum && "" != profileInstID {
			filterInfo.backupNFProfileMd5Sum[profileInstID] = md5Sum
			filterInfo.backuupEtagKeys = append(filterInfo.backuupEtagKeys, profileInstID)
		} else {
			filterInfo.etagExist = false
		}
	}
}

func (p * NFPostFilter) isMatchedAMFGuamiList(queryForm *nfdiscrequest.DiscGetPara)  bool {
	if "AMF" != queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType) {
		return true
	}

	plmn, amfid := queryForm.GetNRFDiscGuamiType()
	if plmn == "" || amfid == "" {
		return true
	}

	matched := false

	_, err := jsonparser.ArrayEach([]byte(p.bodyCommon), func(value []byte, dataType jsonparser.ValueType, offset int, err1 error){
		if matched {
			return
		}
		mcc := ""
		mnc := ""
		amfID := ""
		parserErr := jsonparser.ObjectEach(value, func(key1 []byte, value1 []byte, dataType jsonparser.ValueType, offset int) error {
			if string(key1) == constvalue.SubDataPlmnId {
				parserErr2 := jsonparser.ObjectEach(value1, func(key2 []byte, value2 []byte, dataType jsonparser.ValueType, offset int) error {
					if string(key2) == constvalue.Mcc {
						mcc = string(value2)
					}

					if string(key2) == constvalue.Mnc {
						mnc = string(value2)
					}
					return nil
				})
				if parserErr2 != nil {
					log.Debugf("parse mcc, mnc error, err=%v", parserErr2)
				}
			}

			if string(key1) == constvalue.SearchDataAmfID {
				amfID = string(value1)
			}

			return nil
		})
		if parserErr != nil {
			log.Debugf("jsonparser ObjectEach error, err=%v", parserErr)
		}

		if plmn == (mcc+mnc) && amfid == amfID {
			matched = true
		}
	}, constvalue.GuamiList)
	if err != nil {
		log.Debugf("parse amfInfo guamiList error, err=%v", err)
	}

	if matched {
		return true
	}

	return false

}

func (p * NFPostFilter) replaceFQDNAttr(queryForm *nfdiscrequest.DiscGetPara){
	newprofile := p.bodyCommon
	if !nfdiscutil.IsPlmnMatchHomePlmn(queryForm.GetNRFDiscPlmnListValue(constvalue.SearchDataRequesterPlmnList)) {
		interPlmnFqdn, err := jsonparser.GetString([]byte(p.bodyCommon), constvalue.InterPlmnFqdn)
		if err == nil {
			_, err = jsonparser.GetString([]byte(p.bodyCommon), constvalue.Fqdn)
			if err == nil {
				tempProfile, err := jsonparser.Set([]byte(p.bodyCommon), []byte("\""+interPlmnFqdn+"\""), constvalue.Fqdn)
				newprofile = string(tempProfile)
				if err != nil {
					newprofile = p.bodyCommon
				}
			}
		}

		newservices := make([][]byte, 1)
		_, parseErr := jsonparser.ArrayEach([]byte(p.nfServices), func(value1 []byte, dataType jsonparser.ValueType, offset int, err1 error) {
			interPlmnFqdn, err2 := jsonparser.GetString(value1, constvalue.NFServiceInterPlmnFqdn)
			if err2 == nil {
				_, err2 := jsonparser.GetString(value1, constvalue.NFServiceFqdn)
				if err2 == nil {
					var service []byte
					service, err3 := jsonparser.Set(value1, []byte("\""+interPlmnFqdn+"\""), constvalue.NFServiceFqdn)
					if err3 != nil {
						newservices = append(newservices, value1)
					} else {
						newservices = append(newservices, service)
					}
				} else {
					newservices = append(newservices, value1)
				}
			} else {
				newservices = append(newservices, value1)
			}

		})
		if parseErr != nil {
			log.Debugf("parse nfServices error, err=%v", parseErr)
		}
		services := ""
		for _, value := range newservices {
			if services == "" {
				services = string(value)
			} else {
				services = services + "," + string(value)
			}
		}
		services = "[" + services + "]"

		p.nfServices = services
		p.nfProfile = newprofile
		return
	}
	return
}

func (p* NFPostFilter) eliminateFilelds(queryForm *nfdiscrequest.DiscGetPara) bool {
	log.Debugf("NFPostFilter nfprofile: %s", p.bodyCommon)
	newNfProfile := jsonparser.Delete(
		jsonparser.Delete(
			jsonparser.Delete(
				jsonparser.Delete(
					jsonparser.Delete(
						jsonparser.Delete(
							[]byte(p.bodyCommon), constvalue.InterPlmnFqdn),
						constvalue.HeartBeatTimer),
					constvalue.AllowedNfDomains),
				constvalue.AllowedNFTypes),
			constvalue.AllowedPlmns),
		constvalue.AllowedNssais)

	serviceLength := 0
	_, err := jsonparser.ArrayEach([]byte(p.nfServices), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		serviceLength++
	})
	if err != nil {
		log.Debugf("parsing arry fail for %s,error:%v", constvalue.NfServices, err)
		p.nfProfile = string(newNfProfile)
		return true
	}

	newNfServices := []byte(p.nfServices)
	for i := 0; i < serviceLength; i++ {
		index := fmt.Sprintf("[%d]", i)
		newNfServices = jsonparser.Delete(
			jsonparser.Delete(
				jsonparser.Delete(
					jsonparser.Delete(
						jsonparser.Delete(
							newNfServices, index, constvalue.NFServiceInterPlmnFqdn),
						 index, constvalue.NFServiceAllowedPlmns),
					 index, constvalue.NFServiceAllowedNFTypes),
				 index, constvalue.NFServiceAllowedDomains),
			 index, constvalue.NFServiceAllowedNssais)
	}
	p.nfServices = string(newNfServices)
	p.nfProfile = string(newNfProfile)
	return true

}

func (p *NFPostFilter) eliminatn2INterfaceAmfInfo(queryForm *nfdiscrequest.DiscGetPara) bool {
	if "AMF" == queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType) {
		newNFInfo := jsonparser.Delete([]byte(p.nfInfo), constvalue.N2InterfaceAmfInfo)
		p.nfInfo = string(newNFInfo)
	}
	return true
}
