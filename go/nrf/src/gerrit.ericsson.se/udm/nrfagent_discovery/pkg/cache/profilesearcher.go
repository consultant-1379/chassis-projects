package cache

import (
	"encoding/json"
	"regexp"
	"strconv"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/cache/provider"
	"github.com/buger/jsonparser"
	"github.com/deckarep/golang-set"
)

type profileSearcher struct {
	cache           *cache
	ids             mapset.Set
	searchParameter *SearchParameter
}

//Search is search function of ProfileSearcher
func (ps *profileSearcher) search() (mapset.Set, bool) {
	if ps.searchParameter.searchTargetNfInstanceID() {
		log.Infof("ProfileSearcher : cache search(exclude targetNfInstanceID) hit [%+v]", ps.ids)
		ps.searchByTargetNfInstanceID()
	}
	if ps.searchParameter.searchGpsi() {
		log.Infof("ProfileSearcher : cache search(exclude gpsi) hit [%+v]", ps.ids)
		ps.searchByGpsi()
	}
	if ps.searchParameter.searchExternalGroupIdentity() {
		log.Infof("ProfileSearcher : cache search(exclude externalGroupIdentity) hit [%+v]", ps.ids)
		ps.searchByExternalGroupIdentity()
	}
	if ps.searchParameter.searchDataSet() {
		log.Infof("ProfileSearcher : cache search(exclude dataSet) hit [%+v]", ps.ids)
		ps.searchByDataSet()
	}
	if ps.searchParameter.searchSupi() {
		log.Infof("ProfileSearcher : cache search(exclude supi) hit [%+v]", ps.ids)
		ps.searchBySupi()
	}
	if ps.searchParameter.searchSnssai() {
		log.Infof("ProfileSearcher : cache search(exclude snssai) hit [%+v]", ps.ids)
		ps.searchBySnssais()
	}
	if ps.searchParameter.searchAccessType() {
		log.Infof("ProfileSearcher : cache search(exclude accessType) hit [%+v]", ps.ids)
		ps.searchByAccessType()
	}
	if ps.searchParameter.searchChfSupportedPlmn() {
		log.Infof("ProfileSearcher : cache search(exclude chfSupportedPlmn) hit [%+v]", ps.ids)
		ps.searchByChfSupportedPlmn()
	}
	if ps.searchParameter.searchPreferredLocality() {
		log.Infof("ProfileSearcher : cache search(exclude preferredLocality) hit [%+v]", ps.ids)
		ps.searchByPreferredLocality()
	}
	if ps.searchParameter.searchDnaiList() {
		log.Infof("ProfileSearcher : cache search(exclude dnaiList) hit [%+v]", ps.ids)
		ps.searchByDnaiList()
	}

	if ps.ids == nil ||
		ps.ids.Cardinality() == 0 {
		return nil, false
	}
	log.Infof("search hit[%+v]", ps.ids)
	return ps.ids, true
}

func (ps *profileSearcher) searchByTargetNfInstanceID() {
	if ps.ids == nil ||
		ps.ids.Cardinality() == 0 {
		return
	}

	newIds := mapset.NewSet()
	if ps.ids.Contains(ps.searchParameter.targetNfInstanceID) {
		newIds.Add(ps.searchParameter.targetNfInstanceID)
	}
	ps.ids.Clear()
	ps.ids = newIds
}

func (ps *profileSearcher) searchByGpsi() {
	if ps.ids == nil ||
		ps.ids.Cardinality() == 0 {
		return
	}

	newIds := mapset.NewSet()
	it := ps.ids.Iterator()
	for elem := range it.C {
		id := elem.(string)
		if ps.profileGpsiSearch(id, ps.searchParameter.gpsi, ps.searchParameter.targetNfType) {
			newIds.Add(id)
		}
	}
	ps.ids.Clear()
	ps.ids = newIds
}

func (ps *profileSearcher) searchByExternalGroupIdentity() {
	if ps.ids == nil ||
		ps.ids.Cardinality() == 0 {
		return
	}

	newIds := mapset.NewSet()
	it := ps.ids.Iterator()
	for elem := range it.C {
		id := elem.(string)
		if ps.profileExternalGroupIdentitySearch(id,
			ps.searchParameter.externalGroupIdentity) {
			newIds.Add(id)
		}
	}
	ps.ids.Clear()
	ps.ids = newIds
}

func (ps *profileSearcher) searchByDataSet() {
	if ps.ids == nil ||
		ps.ids.Cardinality() == 0 {
		return
	}

	newIds := mapset.NewSet()
	it := ps.ids.Iterator()
	for elem := range it.C {
		id := elem.(string)
		if ps.profileDataSetSearch(id, ps.searchParameter.dataSet) {
			newIds.Add(id)
		}
	}
	ps.ids.Clear()
	ps.ids = newIds
}

func (ps *profileSearcher) searchBySupi() {
	if ps.ids == nil ||
		ps.ids.Cardinality() == 0 {
		return
	}

	newIds := mapset.NewSet()
	it := ps.ids.Iterator()
	for elem := range it.C {
		id := elem.(string)
		if ps.profileSupiSearch(id,
			ps.searchParameter.supi, ps.searchParameter.targetNfType) {
			newIds.Add(id)
		}
	}
	ps.ids.Clear()
	ps.ids = newIds
}

func (ps *profileSearcher) searchBySnssais() {
	if ps.ids == nil ||
		ps.ids.Cardinality() == 0 {
		return
	}

	newIds := mapset.NewSet()
	it := ps.ids.Iterator()
	for elem := range it.C {
		id := elem.(string)
		if ps.profileSnssaisSearch(id, ps.searchParameter.snssai) {
			newIds.Add(id)
		}
	}
	ps.ids.Clear()
	ps.ids = newIds
}

func (ps *profileSearcher) searchByAccessType() {
	if ps.ids == nil ||
		ps.ids.Cardinality() == 0 {
		return
	}

	newIds := mapset.NewSet()
	it := ps.ids.Iterator()
	for elem := range it.C {
		id := elem.(string)
		if ps.profileAccessTypeSearch(id, ps.searchParameter.accessType) {
			newIds.Add(id)
		}
	}
	ps.ids.Clear()
	ps.ids = newIds
}

func (ps *profileSearcher) searchByChfSupportedPlmn() {
	if ps.ids == nil ||
		ps.ids.Cardinality() == 0 {
		return
	}

	newIds := mapset.NewSet()
	it := ps.ids.Iterator()
	for elem := range it.C {
		id := elem.(string)
		if ps.profileChfSupportedPlmnSearch(id, ps.searchParameter.chfSupportedPlmn) {
			newIds.Add(id)
		}
	}
	ps.ids.Clear()
	ps.ids = newIds
}

func (ps *profileSearcher) searchByPreferredLocality() {
	if ps.ids == nil ||
		ps.ids.Cardinality() == 0 {
		return
	}

	newIds := mapset.NewSet()
	it := ps.ids.Iterator()
	for elem := range it.C {
		id := elem.(string)
		if ps.profilePreferredLocalitySearch(id, ps.searchParameter.preferredLocality) {
			newIds.Add(id)
		}
	}
	ps.ids.Clear()
	ps.ids = newIds
}

func (ps *profileSearcher) searchByDnaiList() {
	if ps.ids == nil ||
		ps.ids.Cardinality() == 0 {
		return
	}

	newIds := mapset.NewSet()
	it := ps.ids.Iterator()
	for elem := range it.C {
		id := elem.(string)
		if ps.profileDnaiListSearch(id, ps.searchParameter.dnn, ps.searchParameter.dnaiList, ps.searchParameter.targetNfType) {
			newIds.Add(id)
		}
	}
	ps.ids.Clear()
	ps.ids = newIds
}

///////////////////////////////////

func (ps *profileSearcher) profileGpsiSearch(id string, number string, nfType string) bool {
	nfInfo := ps.cache.fetchProfileByID(id)
	if nfInfo == nil {
		return false
	}

	result := false
	var groupID string
	var gpsiRangeList = make([]identity, 0)
	if nfType == "UDR" {
		var udrInfo udrInfo
		err := json.Unmarshal(nfInfo, &udrInfo)
		if err != nil {
			log.Errorf("Unmarshal udrInfo gpsi failure, please check the [%s] profile, Error: %s", id, err.Error())
			return false
		}
		if udrInfo.UdrInfo == nil {
			log.Warnf("No udrInfo in [%s] profile.", id)
			return false
		}
		groupID = udrInfo.UdrInfo.GroupID
		gpsiRangeList = udrInfo.UdrInfo.GpsiRanges
	} else if nfType == "UDM" {
		var udmInfo udmInfo
		err := json.Unmarshal(nfInfo, &udmInfo)
		if err != nil {
			log.Errorf("Unmarshal udmInfo gpsi failure, please check the [%s] profile, Error: %s", id, err.Error())
			return false
		}
		if udmInfo.UdmInfo == nil {
			log.Warnf("No udmInfo in [%s] profile.", id)
			return false
		}
		groupID = udmInfo.UdmInfo.GroupID
		gpsiRangeList = udmInfo.UdmInfo.GpsiRanges
	} else {
		log.Warnf("profileGpsiSearch: gpsiRange is not support for %s,  please check the [%s] profile.", nfType, id)
		return false
	}

	if len(gpsiRangeList) == 0 && len(groupID) == 0 {
		log.Debugf("profileGpsiSearch: gpsiRange and groupID not exist, [%s] profile support all gpsi.", id)
		return true
	} else {
		for _, gpsi := range gpsiRangeList {
			if !gpsi.check() {
				log.Infof("[%s] profile one gpsi config is invalid, will skip to check", id)
				continue
			}
			if gpsi.cover(number) {
				result = true
				log.Infof("[%s] profile gpsi cover number[%s]", id, number)
				break
			}
		}
	}
	return result
}

func (ps *profileSearcher) profileExternalGroupIdentitySearch(id string, identity string) bool {
	nfInfo := ps.cache.fetchProfileByID(id)
	if nfInfo == nil {
		return false
	}
	ret := false
	info := []string{"udrInfo", "udmInfo"}

	for _, i := range info {
		item, _, _, err := jsonparser.Get(nfInfo, i)
		log.Debugf("item:%s", string(item))
		if err == nil {
			_, err1 := jsonparser.ArrayEach(item, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				if ret {
					return
				}
				pattern, err := jsonparser.GetString(value, "pattern")
				if err == nil {
					matched, _ := regexp.MatchString(pattern, identity)
					log.Debugf("The externalGroupID: %s, pattern : %s, matched result: %v", identity, pattern, matched)
					if matched {
						ret = true
						return
					}
				}

			}, "externalGroupIdentifiersRanges")

			if err1 != nil {
				ret = false
				log.Debugf("Parsering array fail for externalGroupIdentifiersRanges, error: %v", err)
			}
			if ret {
				return ret
			}
		}
	}

	return ret
}

func (ps *profileSearcher) profileDataSetSearch(id string, dataSet string) bool {
	matched := false
	nfInfo := ps.cache.fetchProfileByID(id)
	if nfInfo == nil {
		return false
	}

	log.Debugf("nfInfo: %s", string(nfInfo))
	_, err := jsonparser.ArrayEach(nfInfo, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if matched {
			return
		}
		log.Debugf("dataset value: %s", string(value[:]))
		dataSetInProfile := string(value[:])
		if dataSet == dataSetInProfile {
			matched = true
			return
		}
	}, "udrInfo", "supportedDataSets")

	if err != nil {
		matched = false
	}

	return matched
}

func (ps *profileSearcher) profileSupiSearch(id string, number string, nfType string) bool {
	nfInfo := ps.cache.fetchProfileByID(id)
	if nfInfo == nil {
		return false
	}

	result := false

	var groupID string
	var supiRangeList = make([]identity, 0)
	if nfType == "UDR" {
		var udrInfo udrInfo
		err := json.Unmarshal(nfInfo, &udrInfo)
		if err != nil {
			log.Warnf("Unmarshal udrInfo supi failure, please check the [%s] profile, Error: %s", id, err.Error())
			return false
		}
		if udrInfo.UdrInfo == nil {
			log.Warnf("No udrInfo in [%s] profile.", id)
			return false
		}
		groupID = udrInfo.UdrInfo.GroupID
		supiRangeList = udrInfo.UdrInfo.SupiRanges
	} else if nfType == "UDM" {
		var udmInfo udmInfo
		err := json.Unmarshal(nfInfo, &udmInfo)
		if err != nil {
			log.Warnf("Unmarshal udmInfo supi failure, please check the [%s] profile, Error: %s", id, err.Error())
			return false
		}
		if udmInfo.UdmInfo == nil {
			log.Warnf("No udmInfo in [%s] profile.", id)
			return false
		}
		groupID = udmInfo.UdmInfo.GroupID
		supiRangeList = udmInfo.UdmInfo.SupiRanges
	} else if nfType == "AUSF" {
		var ausfInfo ausfInfo
		err := json.Unmarshal(nfInfo, &ausfInfo)
		if err != nil {
			log.Errorf("Unmarshal ausfInfo failure, please check the [%s] profile, Error: %s", id, err.Error())
			return false
		}
		if ausfInfo.AusfInfo == nil {
			log.Warnf("No ausfInfo in [%s] profile.", id)
			return false
		}
		groupID = ausfInfo.AusfInfo.GroupID
		supiRangeList = ausfInfo.AusfInfo.SupiRanges
	} else if nfType == "PCF" {
		var pcfInfo pcfInfo
		err := json.Unmarshal(nfInfo, &pcfInfo)
		if err != nil {
			log.Errorf("Unmarshal pcfInfo failure, please check the [%s] profile, Error: %s", id, err.Error())
			return false
		}
		if pcfInfo.PcfInfo == nil {
			log.Warnf("No pcfInfo in [%s] profile.", id)
			return false
		}
		supiRangeList = pcfInfo.PcfInfo.SupiRanges
	} else {
		log.Warnf("profileSupiSearch: supiRange is not support for %s,  please check the [%s] profile.", nfType, id)
		return false
	}

	if len(supiRangeList) == 0 && len(groupID) == 0 {
		log.Debugf("profileSupiSearch: supiRange and groupID not exist, [%s] profile match all supi.", id)
		return true
	} else {
		for _, supi := range supiRangeList {
			if !supi.check() {
				log.Infof("[%s] profile one supi config is invalid, will skip to check", id)
				continue
			}
			if supi.cover(number) {
				result = true
				log.Infof("[%s] profile supi cover number[%s]", id, number)
				break
			}
		}
	}

	return result
}

func (ps *profileSearcher) profileSnssaisSearch(id string, snssai []SNssai) bool {
	nfInfo := ps.cache.fetchProfileByID(id)
	if nfInfo == nil {
		return false
	}

	var nfProfile structs.SearchResultNFProfile

	err := json.Unmarshal(nfInfo, &nfProfile)
	if err != nil {
		log.Infof("fetchContentByServName Unmarshal content fail, Error: %s", err.Error())
		return false
	}

	sNssaiArray := ps.searchParameter.snssai
	if len(nfProfile.SNSSAI) == 0 {
		return false
	}

	result := false
	for _, nfProfileSnssai := range nfProfile.SNSSAI {
		for _, sNssai := range sNssaiArray {
			if nfProfileSnssai.SD == "" && sNssai.Sd == "" {
				if nfProfileSnssai.SST == sNssai.Sst {
					result = true
					break
				}
			}
			if nfProfileSnssai.SD != "" && sNssai.Sd != "" {
				if nfProfileSnssai.SST == sNssai.Sst && nfProfileSnssai.SD == sNssai.Sd {
					result = true
					break
				}
			}
		}
		if result == true {
			break
		}
	}

	log.Infof("fetchContentBySnssais result %+v", result)
	return result
}

func (ps *profileSearcher) profileAccessTypeSearch(id string, accessType string) bool {
	nfInfo := ps.cache.fetchProfileByID(id)
	if nfInfo == nil {
		return false
	}

	matched := false
	_, err := jsonparser.ArrayEach(nfInfo, func(value []byte, dataType jsonparser.ValueType, offset int, err1 error) {
		accessTypeInProfile := string(value[:])
		if accessType == accessTypeInProfile {
			matched = true
			return
		}
	}, provider.SmfInfo, provider.AccessType)
	if err != nil {
		matched = false
	}

	return matched
}

func (ps *profileSearcher) profileChfSupportedPlmnSearch(id string, supportedPlmn structs.PlmnID) bool {
	nfInfo := ps.cache.fetchProfileByID(id)
	if nfInfo == nil {
		return false
	}

	ret := false
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
			log.Debugf("plmnRangeList not have start & end & pattern, match each plmn")
			ret = true
			return
		}
		if err1 == nil {
			matched, _ := regexp.MatchString(pattern, supportedPlmn.Mcc+supportedPlmn.Mnc)
			log.Debugf("The chf-supported-plmn: %v, pattern : %s, matched result: %v", supportedPlmn, pattern, matched)
			if matched {
				ret = true
				return
			}
		}

		if err2 != nil || err3 != nil {
			return
		}
		if ps.isMccMatched(s, e, supportedPlmn.Mcc) && ps.isMncMatched(s, e, supportedPlmn.Mnc) {
			ret = true
			log.Debugf("the chf-supported-plmn is %v, range is %s-%s matched result is true", supportedPlmn, s, e)
			return
		}
		log.Debugf("the chf-supported-plmn is %v, range is %s-%s matched result is %t", supportedPlmn, s, e, ret)

	}, provider.ChfInfo, provider.PlmnRangeList)

	if err != nil || ret == true {
		return true
	}

	if num == 0 && err == nil {
		log.Debugf("plmnRangeList is [], this will match all chf-supported-plmn")
		return true
	}
	return ret
}

//mcc match
func (ps *profileSearcher) isMccMatched(start, end, mccStr string) bool {
	s, _ := strconv.ParseInt(start[0:3], 10, 64)
	e, _ := strconv.ParseInt(end[0:3], 10, 64)
	mcc, _ := strconv.ParseInt(mccStr, 10, 64)

	if mcc >= s && mcc <= e {
		return true
	}

	return false
}

//mnc match
func (ps *profileSearcher) isMncMatched(start, end, mncStr string) bool {
	mncStart := start[3:]
	mncEnd := end[3:]
	if len(mncStart) > len(mncStr) || len(mncEnd) < len(mncStr) {
		return false
	}
	s, _ := strconv.ParseInt(mncStart, 10, 64)
	e, _ := strconv.ParseInt(mncEnd, 10, 64)
	mnc, _ := strconv.ParseInt(mncStr, 10, 64)
	if mnc >= s && mnc <= e {
		return true
	}

	return false
}

func (ps *profileSearcher) profilePreferredLocalitySearch(id string, preferredLocality string) bool {
	nfInfo := ps.cache.fetchProfileByID(id)
	if nfInfo == nil {
		return false
	}

	if preferredLocality != "" {
		preferredLocalityInProfile, err := jsonparser.GetString(nfInfo, provider.Locality)
		if err == nil && preferredLocalityInProfile == preferredLocality {
			return true
		}
		return false
	}
	return true
}

func (ps *profileSearcher) profileDnaiListSearch(id string, dnn string, dnaiList []string, nfType string) bool {
	nfInfo := ps.cache.fetchProfileByID(id)
	if nfInfo == nil {
		return false
	}

	if nfType == "UPF" {
		upfInfoData, _, _, err := jsonparser.Get(nfInfo, "upfInfo")
		if err != nil {
			log.Warnf("upfInfo not include UPF instance %s", id)
			return false
		}

		var upfInfo structs.UPFInfo
		err = json.Unmarshal(upfInfoData, &upfInfo)
		if err != nil {
			log.Warnf("Unmarshal upfInfo failure, please check the [%s] profile, Error: %s", id, err.Error())
			return false
		}
		for _, sNssaiUpfInfo := range upfInfo.SNssaiUpfInfoList {
			for _, dnnUpfInfo := range sNssaiUpfInfo.DNNUpfInfoList {
				if dnn != "" {
					if dnnUpfInfo.DNN == dnn {
						if len(dnnUpfInfo.DnaiList) == 0 {
							log.Debugf("match all dnai when dnai not exist.dnn(%s), id(%s)", dnn, id)
							return true
						}
						for _, dnaiItem := range dnnUpfInfo.DnaiList {
							for _, dnai := range dnaiList {
								if dnaiItem == dnai {
									return true
								}
							}
						}
					}
				} else {
					if dnnUpfInfo.DNN != "" && len(dnnUpfInfo.DnaiList) == 0 {
						log.Debugf("match all dnai when dnai not exist, id(%s)", id)
						return true
					}
					if dnnUpfInfo.DnaiList != nil {
						for _, dnaiItem := range dnnUpfInfo.DnaiList {
							for _, dnai := range dnaiList {
								if dnaiItem == dnai {
									return true
								}
							}
						}
					}
				}
			}
		}
	}

	return false
}
