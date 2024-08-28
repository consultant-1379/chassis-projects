package cache

import (
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/cache/provider"
	"github.com/deckarep/golang-set"
)

type indexSearcher struct {
	cache           *cache
	searchParameter *SearchParameter
	indexMapper     searchIndexMapper
}

type searchPoint struct {
	indexCategory string
	searchKey     string
}

//Search search by index in cache
func (is *indexSearcher) search() (mapset.Set, bool) {
	result := is.meetIndexSearch()
	if !result {
		log.Warnf("Does not meet index search condition, will forward to query to NRF")
		return nil, false
	}

	//Instance().ShowIndexContent(is.searchParameter.requesterNfType, is.searchParameter.targetNfType)
	ids := is.indexSearch()
	if ids == nil ||
		ids.Cardinality() == 0 {
		return nil, false
	}

	log.Infof("search hit[%+v]", ids)
	return ids, true
}

///////////////////////////////private function////////////////////////////////

func (is *indexSearcher) meetIndexSearch() bool {
	if is.searchParameter.searchTargetNfType() && !is.indexMapper.indexedTargetNfType() {
		return false
	}
	if is.searchParameter.searchTargetPlmnList() && !is.indexMapper.indexedTargetPlmnList() {
		return false
	}
	if is.searchParameter.SearchServiceName() && !is.indexMapper.indexedServiceName() {
		return false
	}
	if is.searchParameter.searchDnn() && !is.indexMapper.indexedDnn() {
		return false
	}
	if is.searchParameter.searchSmfServingArea() && !is.indexMapper.indexedSmfServingArea() {
		return false
	}
	if is.searchParameter.searchRoutingIndicator() && !is.indexMapper.indexedRoutingIndicator() {
		return false
	}
	if is.searchParameter.searchNsiList() && !is.indexMapper.indexedNsiList() {
		return false
	}
	if is.searchParameter.searchIPDomain() && !is.indexMapper.indexedIPDomain() {
		return false
	}
	/*
		if is.searchParameter.searchDnaiList() && !is.indexMapper.indexedDnaiList() {
			return false
		}
	*/
	if is.searchParameter.searchUpfIwkEpsInd() && !is.indexMapper.indexedUpfIwkEpsInd() {
		return false
	}
	//	if is.searchParameter.searchSnssai() && !is.indexMapper.SearchSnssai() {
	//		return false
	//	}
	//	if is.searchParameter.searchSupi() && !is.indexMapper.SearchSupi() {
	//		return false
	//	}

	return true
}

func (is *indexSearcher) indexSearch() mapset.Set {
	//	profileIds := make([]mapset.Set, 0)
	//	searchPointList := searchCondition.SearchPointInjection(cm.indexMapper)
	//	for _, searchPoint := range searchPointList {
	//		ids := searchPoint.search()
	//		profileIds = append(profileIds, ids)
	//	}
	//	intersactProfileIds := cm.intersection(profileIds)
	//	return intersactProfileIds

	intersectSets := make([]mapset.Set, 0)
	//	unionSets := make(map[string][]mapset.Set)
	serviceNameListSets := make([]mapset.Set, 0)
	mccMncUnionSets := make([]mapset.Set, 0)
	nsiListSets := make([]mapset.Set, 0)
	groupIDListSets := make([]mapset.Set, 0)
	dnaiListSets := make([]mapset.Set, 0)

	serviceNameExist := false
	mccMncExist := false
	nsiListExist := false
	dnaiListExist := false
	groupIDListExist := false

	searchPointList := is.searchParameter.IndexSearchPointInjection(is.indexMapper)
	log.Infof("search cacheIndex point list : %+v", searchPointList)
	for i := range searchPointList {
		ids := searchPointList[i].search(is.cache)
		indexCategory := searchPointList[i].indexCategory
		if indexCategory == "serviceName" {
			serviceNameExist = true
			if ids == nil ||
				ids.Cardinality() == 0 {
				continue
			}
			serviceNameListSets = append(serviceNameListSets, ids)
		} else if indexCategory == "mcc:mnc" {
			mccMncExist = true
			if ids == nil ||
				ids.Cardinality() == 0 {
				continue
			}
			mccMncUnionSets = append(mccMncUnionSets, ids)
		} else if indexCategory == "nsiList" {
			nsiListExist = true
			if ids == nil ||
				ids.Cardinality() == 0 {
				continue
			}
			nsiListSets = append(nsiListSets, ids)
		} else if indexCategory == "groupId" {
			groupIDListExist = true
			if ids == nil ||
				ids.Cardinality() == 0 {
				continue
			}
			groupIDListSets = append(groupIDListSets, ids)
		} else if indexCategory == "dnaiList" {
			dnaiListExist = true
			if ids == nil ||
				ids.Cardinality() == 0 {
				continue
			}
			dnaiListSets = append(dnaiListSets, ids)
		} else {
			if ids == nil ||
				ids.Cardinality() == 0 {
				return nil
			}
			intersectSets = append(intersectSets, ids)
		}
	}

	if len(serviceNameListSets) > 0 {
		serviceNameSet := serviceNameListSets[0]
		for i := 1; i < len(serviceNameListSets); i++ {
			serviceNameSet = serviceNameSet.Union(serviceNameListSets[i])
		}
		intersectSets = append(intersectSets, serviceNameSet)
	} else if serviceNameExist {
		return nil
	}

	if len(mccMncUnionSets) > 0 {
		mccMncSet := mccMncUnionSets[0]
		for i := 1; i < len(mccMncUnionSets); i++ {
			mccMncSet = mccMncSet.Union(mccMncUnionSets[i])
		}
		intersectSets = append(intersectSets, mccMncSet)
	} else if mccMncExist {
		return nil
	}

	if len(nsiListSets) > 0 {
		nsiSet := nsiListSets[0]
		for i := 1; i < len(nsiListSets); i++ {
			nsiSet = nsiSet.Union(nsiListSets[i])
		}
		intersectSets = append(intersectSets, nsiSet)
	} else if nsiListExist {
		return nil
	}

	if len(groupIDListSets) > 0 {
		groupIDSet := groupIDListSets[0]
		for i := 1; i < len(groupIDListSets); i++ {
			groupIDSet = groupIDSet.Union(groupIDListSets[i])
		}
		intersectSets = append(intersectSets, groupIDSet)
	} else if groupIDListExist {
		return nil
	}

	if len(dnaiListSets) > 0 {
		dnaiSet := dnaiListSets[0]
		for i := 1; i < len(dnaiListSets); i++ {
			dnaiSet = dnaiSet.Union(dnaiListSets[i])
		}
		intersectSets = append(intersectSets, dnaiSet)
	} else if dnaiListExist {
		return nil
	}

	//	for indexCategory := range unionSets {
	//		unionSet := unionSets[indexCategory][0]
	//		for i := 1; i < len(unionSets[indexCategory]); i++ {
	//			if !unionSets[indexCategory][i].IsSubset(unionSet) {
	//				unionSet = unionSet.Union(unionSets[indexCategory][i])
	//			}
	//		}
	//		intersectSets = append(intersectSets, unionSet)
	//	}

	log.Infof("search cacheIndex intersectSets: %+v", intersectSets)
	return is.intersection(intersectSets)
}

func (is *indexSearcher) intersection(intersectSets []mapset.Set) mapset.Set {
	if len(intersectSets) == 0 {
		return nil
	}
	// tuning: cache miss
	for i := range intersectSets {
		if intersectSets[i] == nil ||
			intersectSets[i].Cardinality() == 0 {
			return nil
		}
	}

	rest := intersectSets[0]
	for i := 1; i < len(intersectSets); i++ {
		rest = rest.Intersect(intersectSets[i])
		if rest.Cardinality() == 0 {
			return nil
		}
	}

	//	if mccMncSet.Cardinality() != 0 {
	//		rest = rest.Intersect(mccMncSet)
	//	}

	return rest
	// var profileIds []string
	// rest.Each(func(item interface{}) bool {
	// 	switch t := item.(type) {
	// 	case string:
	// 		profileIds = append(profileIds, t)
	// 	default:
	// 		log.Errorf("unsupported interface type")
	// 	}
	// 	return false
	// })

	// return profileIds
}

func (sp *searchPoint) search(c *cache) mapset.Set {
	if _, ok := c.cacheIndex[sp.indexCategory]; !ok {
		log.Errorf("no index-category[%s] in cache", sp.indexCategory)
		return mapset.NewSet()
	}
	var retSet = mapset.NewSet()
	matchKeySet, ok1 := c.cacheIndex[sp.indexCategory][sp.searchKey]
	if sp.indexCategory == "groupId" {
		matchAllSet, ok2 := c.cacheIndex[sp.indexCategory][provider.MatchAllGroupID]
		if ok1 {
			retSet = matchKeySet
		}
		if ok2 {
			it := matchAllSet.Iterator()
			for elem := range it.C {
				id := elem.(string)
				retSet.Add(id)
			}
		}
		if !ok1 && !ok2 {
			log.Warnf("no profiles in cache for index-category[%s] and index-value[%s]",
				sp.indexCategory, sp.searchKey)
		}
		return retSet
	}
	if !ok1 {
		log.Warnf("no profiles in cache for index-category[%s] and index-value[%s]",
			sp.indexCategory, sp.searchKey)
		return mapset.NewSet()
	}
	return c.cacheIndex[sp.indexCategory][sp.searchKey]
}
