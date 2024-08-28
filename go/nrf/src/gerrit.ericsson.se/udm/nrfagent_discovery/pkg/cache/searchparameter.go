package cache

import (
	"fmt"

	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
)

type SearchParameter struct {
	serviceNames          []string
	targetNfType          string
	requesterNfType       string
	targetPlmnList        []structs.PlmnID
	requesterPlmnList     []structs.PlmnID
	snssai                []SNssai
	dnn                   string
	smfServingArea        string
	tai                   Tai
	ecgi                  Ecgi
	ncgi                  Ncgi
	supi                  string
	supportedFeatures     string
	routingIndicator      string
	nsiList               []string
	gpsi                  string
	externalGroupIdentity string
	dataSet               string
	targetNfInstanceID    string
	groupIDList           []string
	ipDomain              string
	dnaiList              []string
	upfIwkEpsInd          string
	chfSupportedPlmn      structs.PlmnID
	preferredLocality     string
	accessType            string
}

//SetServiceNames
func (sp *SearchParameter) SetServiceNames(serviceNames []string) {
	sp.serviceNames = serviceNames
}

//SetGroupIDList set groupIDList to searchParameter
func (sp *SearchParameter) SetGroupIDList(groupIDList []string) {
	sp.groupIDList = groupIDList
}

//SetIPDomain set ipDomain to searchParameter
func (sp *SearchParameter) SetIPDomain(ipDomain string) {
	sp.ipDomain = ipDomain
}

//SetDnaiList set dnaiList to searchParameter
func (sp *SearchParameter) SetDnaiList(dnaiList []string) {
	sp.dnaiList = dnaiList
}

//SetUpfIwkEpsInd set upfIwkEpsInd to searchParameter
func (sp *SearchParameter) SetUpfIwkEpsInd(upfIwkEpsInd string) {
	sp.upfIwkEpsInd = upfIwkEpsInd
}

//SetTargetNfType
func (sp *SearchParameter) SetTargetNfType(targetNfType string) {
	sp.targetNfType = targetNfType
}

//SetRequesterNfType
func (sp *SearchParameter) SetRequesterNfType(requesterNfType string) {
	sp.requesterNfType = requesterNfType
}

//SetTargetPlmn
func (sp *SearchParameter) SetTargetPlmnList(targetPlmnList []structs.PlmnID) {
	sp.targetPlmnList = targetPlmnList
}

//SetRequesterPlmn
func (sp *SearchParameter) SetRequesterPlmnList(requesterPlmnList []structs.PlmnID) {
	sp.requesterPlmnList = requesterPlmnList
}

//SetSnssai
func (sp *SearchParameter) SetSnssai(snssai []SNssai) {
	sp.snssai = snssai
}

//SetDnn
func (sp *SearchParameter) SetDnn(dnn string) {
	sp.dnn = dnn
}

//SetSmfServingArea
func (sp *SearchParameter) SetSmfServingArea(smfServingArea string) {
	sp.smfServingArea = smfServingArea
}

//SetTai
func (sp *SearchParameter) SetTai(tai Tai) {
	sp.tai = tai
}

//SetEcgi
func (sp *SearchParameter) SetEcgi(ecgi Ecgi) {
	sp.ecgi = ecgi
}

//SetNcgi
func (sp *SearchParameter) SetNcgi(ncgi Ncgi) {
	sp.ncgi = ncgi
}

//SetSupi
func (sp *SearchParameter) SetSupi(supi string) {
	sp.supi = supi
}

//SetSupportedFeatures is ...
func (sp *SearchParameter) SetSupportedFeatures(supportedFeatures string) {
	sp.supportedFeatures = supportedFeatures
}

//SetRoutingIndicator
func (sp *SearchParameter) SetRoutingIndicator(routingIndicator string) {
	sp.routingIndicator = routingIndicator
}

//SetNsiList
func (sp *SearchParameter) SetNsiList(nsiList []string) {
	sp.nsiList = nsiList
}

//SetGpsi
func (sp *SearchParameter) SetGpsi(gpsi string) {
	sp.gpsi = gpsi
}

//SetExternalGroupIdentity
func (sp *SearchParameter) SetExternalGroupIdentity(externalGroupIdentity string) {
	sp.externalGroupIdentity = externalGroupIdentity
}

//SetDataSet
func (sp *SearchParameter) SetDataSet(dataSet string) {
	sp.dataSet = dataSet
}

//SetTargetNfInstanceID
func (sp *SearchParameter) SetTargetNfInstanceID(targetNfInstanceID string) {
	sp.targetNfInstanceID = targetNfInstanceID
}

//SetChfSupportedPlmn is for set chfSupportedPlmn value
func (sp *SearchParameter) SetChfSupportedPlmn(chfSupportedPlmn structs.PlmnID) {
	sp.chfSupportedPlmn = chfSupportedPlmn
}

//SetPreferredLocality is for set preferredLocality
func (sp *SearchParameter) SetPreferredLocality(preferredLocality string) {
	sp.preferredLocality = preferredLocality
}

//SetAccessType is for set accessType
func (sp *SearchParameter) SetAccessType(accessType string) {
	sp.accessType = accessType
}

//SearchServiceName search by serviceNames
func (sp *SearchParameter) SearchServiceName() bool {
	return len(sp.serviceNames) != 0
}

//searchGroupIDList search by groupIDList
func (sp *SearchParameter) searchGroupIDList() bool {
	return len(sp.groupIDList) != 0
}

//searchTargetNfType search by targetNfType
func (sp *SearchParameter) searchTargetNfType() bool {
	return sp.targetNfType != ""
}

//searchRequesterNfType search by requesterNfType
func (sp *SearchParameter) searchRequesterNfType() bool {
	return sp.requesterNfType != ""
}

//searchSnssai search by snssai
func (sp *SearchParameter) searchSnssai() bool {
	return len(sp.snssai) != 0
}

//searchPreferredLocality search by preferredLocality
func (sp *SearchParameter) searchPreferredLocality() bool {
	return sp.preferredLocality != ""
}

//searchAccessType search by accessType
func (sp *SearchParameter) searchAccessType() bool {
	return sp.accessType != ""
}

//searchChfSupportedPlmn search by chfSupportedPlmn
func (sp *SearchParameter) searchChfSupportedPlmn() bool {
	prober := structs.PlmnID{}
	return sp.chfSupportedPlmn != prober
}

//searchTargetPlmnList search by targetPlmn
func (sp *SearchParameter) searchTargetPlmnList() bool {
	return len(sp.targetPlmnList) != 0
}

//searchRequesterPlmnList search by requesterplmn
func (sp *SearchParameter) searchRequesterPlmnList() bool {
	return len(sp.requesterPlmnList) != 0
}

//searchDnn search by dnn
func (sp *SearchParameter) searchDnn() bool {
	return sp.dnn != ""
}

func (sp *SearchParameter) searchSmfServingArea() bool {
	return sp.smfServingArea != ""
}

//SearchSupportedFeatures search by supportedFeatures
func (sp *SearchParameter) SearchSupportedFeatures() bool {
	return sp.supportedFeatures != ""
}

//searchRoutingIndicator search by routingIndicator
func (sp *SearchParameter) searchRoutingIndicator() bool {
	return sp.routingIndicator != ""
}

//searchNsiList search by nsiList
func (sp *SearchParameter) searchNsiList() bool {
	return len(sp.nsiList) != 0
}

//searchSupi search by supi
func (sp *SearchParameter) searchSupi() bool {
	return sp.supi != ""
}

//searchIPDomain search by ipDomain
func (sp *SearchParameter) searchIPDomain() bool {
	return sp.ipDomain != ""
}

//searchDnaiList search by dnaiList
func (sp *SearchParameter) searchDnaiList() bool {
	return len(sp.dnaiList) != 0
}

//searchUpfIwkEpsInd search by upfIwkEpsInd
func (sp *SearchParameter) searchUpfIwkEpsInd() bool {
	return sp.upfIwkEpsInd != ""
}

//searchTargetNfInstanceID search by targetNfInstanceID
func (sp *SearchParameter) searchTargetNfInstanceID() bool {
	return sp.targetNfInstanceID != ""
}

//searchGpsi search by gpsi
func (sp *SearchParameter) searchGpsi() bool {
	return sp.gpsi != ""
}

//searchExternalGroupIdentity search by externalGroupIdentity
func (sp *SearchParameter) searchExternalGroupIdentity() bool {
	return sp.externalGroupIdentity != ""
}

//searchDataSet search by dataSet
func (sp *SearchParameter) searchDataSet() bool {
	return sp.dataSet != ""
}

//ProfileSearchNecessary profile search is necessary
func (sp *SearchParameter) ProfileSearchNecessary() bool {
	result := sp.searchTargetNfInstanceID() ||
		sp.searchGpsi() ||
		sp.searchExternalGroupIdentity() ||
		sp.searchDataSet() ||
		sp.searchSupi() ||
		sp.searchSnssai() ||
		sp.searchPreferredLocality() ||
		sp.searchAccessType() ||
		sp.searchChfSupportedPlmn() ||
		sp.searchDnaiList()

	return result
}

func (sp *SearchParameter) IndexSearchPointInjection(indexMapper searchIndexMapper) []searchPoint {
	searchPointList := make([]searchPoint, 0)

	if sp.SearchServiceName() {
		for _, serviceName := range sp.serviceNames {
			point := searchPoint{
				indexCategory: indexMapper.ServiceName,
				searchKey:     serviceName,
			}
			searchPointList = append(searchPointList, point)
		}
	}

	if sp.searchGroupIDList() {
		for _, groupID := range sp.groupIDList {
			point := searchPoint{
				indexCategory: indexMapper.GroupIDList,
				searchKey:     groupID,
			}
			searchPointList = append(searchPointList, point)
		}
	}

	if sp.searchTargetNfType() {
		point := searchPoint{
			indexCategory: indexMapper.TargetNfType,
			searchKey:     sp.targetNfType,
		}
		searchPointList = append(searchPointList, point)
	}

	if sp.searchTargetPlmnList() {
		for _, plmnIDItem := range sp.targetPlmnList {
			point := searchPoint{
				indexCategory: fmt.Sprintf("%s:%s", indexMapper.TargetPlmnList[0].Mcc, indexMapper.TargetPlmnList[0].Mnc),
				searchKey:     fmt.Sprintf("%s:%s", plmnIDItem.Mcc, plmnIDItem.Mnc),
			}
			searchPointList = append(searchPointList, point)
		}
	}
	//	if sp.searchTargetPlmn() {
	//		point := SearchPoint{
	//			indexCategory: fmt.Sprintf("%s:%s", indexMapper.TargetPlmnList[0].Mcc, indexMapper.TargetPlmnList[0].Mnc),
	//			searchKey:     fmt.Sprintf("%s:%s", sp.targetPlmn[0].Mcc, sp.targetPlmn[0].Mnc),
	//		}
	//		searchPointList = append(searchPointList, point)
	//	}

	if sp.searchDnn() {
		point := searchPoint{
			indexCategory: indexMapper.Dnn,
			searchKey:     sp.dnn,
		}
		searchPointList = append(searchPointList, point)
	}

	if sp.searchSmfServingArea() {
		point := searchPoint{
			indexCategory: indexMapper.SmfServingArea,
			searchKey:     sp.smfServingArea,
		}
		searchPointList = append(searchPointList, point)
	}

	if sp.searchRoutingIndicator() {
		point := searchPoint{
			indexCategory: indexMapper.RoutingIndicator,
			searchKey:     sp.routingIndicator,
		}
		searchPointList = append(searchPointList, point)
	}

	if sp.searchNsiList() {
		for _, nsi := range sp.nsiList {
			point := searchPoint{
				indexCategory: indexMapper.NsiList,
				searchKey:     nsi,
			}
			searchPointList = append(searchPointList, point)
		}
	}

	if sp.searchIPDomain() {
		point := searchPoint{
			indexCategory: indexMapper.IPDomain,
			searchKey:     sp.ipDomain,
		}
		searchPointList = append(searchPointList, point)
	}

	/*
		if sp.searchDnaiList() {
			for _, dnai := range sp.dnaiList {
				point := searchPoint{
					indexCategory: indexMapper.DnaiList,
					searchKey:     dnai,
				}
				searchPointList = append(searchPointList, point)
			}
		}
	*/

	if sp.searchUpfIwkEpsInd() {
		point := searchPoint{
			indexCategory: indexMapper.UpfIwkEpsInd,
			searchKey:     sp.upfIwkEpsInd,
		}
		searchPointList = append(searchPointList, point)
	}

	/*
		if sp.searchSnssai() {
			point := searchPoint{
				indexCategory: fmt.Sprintf("%s:%s", indexMapper.Snssai.Sst, indexMapper.Snssai.Sd),
				searchKey:     fmt.Sprintf("%s:%s", sim.Snssai.Sst, sim.Snssai.Sd),
			}
			searchPointList = append(searchPointList, point)
		}
	*/

	/*
		if sp.searchNfInstanceID() {
			point := searchPoint{
				indexCategory: indexMapper.NfInstanceID,
				searchKey:     sim.NfInstanceID,
			}
			searchPointList = append(searchPointList, point)
		}
	*/

	/*
		if sp.searchSupi() {
			point := searchPoint{
				indexCategory: indexMapper.Supi,
				searchKey:     sim.Supi,
			}
			searchPointList = append(searchPointList, point)
		}
	*/

	return searchPointList
}
