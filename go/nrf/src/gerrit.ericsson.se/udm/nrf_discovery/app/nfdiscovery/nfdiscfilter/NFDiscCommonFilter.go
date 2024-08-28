package nfdiscfilter

import (
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscrequest"
	"github.com/buger/jsonparser"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
)

//NFCommonFilter to process common  nfprofile filter
type NFCommonFilter struct {
}

//filter is to match nftype adn target-nf-fqdn in nfprofile
func (c *NFCommonFilter) filter(nfprofile []byte, queryForm *nfdiscrequest.DiscGetPara, filterInfo *FilterInfo) bool {
	log.Debugf("Enter Common Filter")
	if !filterInfo.KVDBSearch {
		nfType, err := jsonparser.GetString(nfprofile, constvalue.NfType)
		if err != nil || nfType != queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataTargetNfType) {
			return false
		}
	}

	if !filterInfo.KVDBSearch && queryForm.GetNRFDiscTargetNFFQDN() != "" {
		if !c.isMatchedNFProfileFQDN(queryForm.GetNRFDiscTargetNFFQDN(), nfprofile) {
			log.Debugf("No Matched nfProfile with Target-nf-fqdn: %s", queryForm.GetNRFDiscTargetNFFQDN())
			return false
		}
	}

	if !c.isMatchedNFProfilesStatus(nfprofile){
		return false
	}

	/*targetplmn, _ := queryForm.GetNRFDiscPlmnValue(constvalue.SearchDataTargetPlmn)
	if targetplmn != "" {
		if !c.isMatchedTargetPlmn(targetplmn, nfprofile){
			return false
		}
	}*/

	nsiArray := queryForm.GetNRFDiscNsiList()
	if nsiArray != nil {
		if !c.isMatchedNsiList(nsiArray, nfprofile) {
			return false
		}
	}
	requesterNfType := queryForm.GetNRFDiscNFTypeValue(constvalue.SearchDataRequesterNfType)
	if requesterNfType != "" {
		_, _, _, err := jsonparser.Get(nfprofile, constvalue.AllowedNFTypes)
		if err == nil {
			if !nfdiscutil.IsAllowedNfType(nfprofile, requesterNfType, constvalue.AllowedNFTypes) {
				log.Debugf("reuqester-nf-type is forbidden to access nfprofile in field allowedNfTypes")
				filterInfo.nfTypeForbiddenInProfile = true
				return false
			}
		}
	}
	requesterPlmnList := queryForm.GetNRFDiscPlmnListValue(constvalue.SearchDataRequesterPlmnList)
	if len(requesterPlmnList) > 0 {

		if !isAllowedRequesterPlmn(nfprofile, requesterPlmnList, constvalue.AllowedPlmns, []byte(filterInfo.originProfile.BodyCommon)){
			filterInfo.plmnForbiddenInProfile = true
			return false
		}

	}
	targetPlmnList := queryForm.GetNRFDiscPlmnListValue(constvalue.SearchDataTargetPlmnList)
	if len(targetPlmnList) > 0 {
		_, _, _, err := jsonparser.Get(nfprofile, constvalue.PlmnList)
		if err == nil {
			if !nfdiscutil.IsAllowedPLMN(nfprofile, targetPlmnList, constvalue.PlmnList) {
				log.Debugf("target-plmn-list is not matched with nfprofile in field plmnList")
				return false
			}
		}
	}
	requesterNfInstanceFqdn := queryForm.GetNRFDiscRequesterNFInstFQDN()
	if requesterNfInstanceFqdn != "" {
		_, _, _, err := jsonparser.Get(nfprofile, constvalue.AllowedNfDomains)
		if err == nil {
			if !nfdiscutil.IsAllowedNfFQDN(nfprofile, requesterNfInstanceFqdn, constvalue.AllowedNfDomains) {
				log.Debugf("reuqester-nf-instance-fqdn is forbidden to access nfprofile in field allowedNfDomains")
				filterInfo.domainForbiddenInProfile = true
				return false
			}
		}
	}
	return true
}

//isMatchedNFProfileFQDN is to match fqdn in nfprofile
func (c *NFCommonFilter)isMatchedNFProfileFQDN(fqdn string, item []byte) bool {
	fqdnInProfile, err := jsonparser.GetString(item, constvalue.Fqdn)
	if err == nil {
		if fqdn == fqdnInProfile {
			return true
		}
	}

	return false
}

func (c *NFCommonFilter)isMatchedNFProfilesStatus(item []byte) bool{
	status, err := jsonparser.GetString(item, constvalue.NfStatus)
	if err == nil {
		if status == constvalue.NFStatusRegistered{
			return true
		}
	}

	return false
}

func (c *NFCommonFilter) isMatchedNsiList(nsiList []string, item []byte) bool {
	var matched bool

	_, err := jsonparser.ArrayEach(item, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if matched {
			return
		}
		for _, v := range nsiList{
			if v == string(value){
				matched = true
				return
			}
		}

	}, constvalue.NsiList)
        if matched && err == nil {
		return true
	}
	return false
}
