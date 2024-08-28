package cache

import (
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
)

type searchIndexMapper struct {
	ServiceName      string           `json:"service-names,omitempty"`
	TargetNfType     string           `json:"target-nf-type,omitempty"`
	TargetPlmnList   []structs.PlmnID `json:"target-plmn-list,omitempty"`
	Dnn              string           `json:"dnn,omitempty"`
	SmfServingArea   string           `json:"smf-serving-area,omitempty"`
	SupportedFeature string           `json:"supported-features,omitempty"`
	RoutingIndicator string           `json:"routing-indicator,omitempty"`
	NsiList          string           `json:"nsi-list,omitempty"`
	GroupIDList      string           `json:"group-id-list,omitempty"`
	IPDomain         string           `json:"ip-domain,omitempty"`
	DnaiList         string           `json:"dnai-list,omitempty"`
	UpfIwkEpsInd     string           `json:"upf-iwk-eps-ind,omitempty"`
}

//indexedServiceName index ServiceName
func (sim *searchIndexMapper) indexedServiceName() bool {
	return sim.ServiceName != ""
}

//indexedTargetNfType indexed TargetNfType
func (sim *searchIndexMapper) indexedTargetNfType() bool {
	return sim.TargetNfType != ""
}

//indexedTargetPlmnList indexed TargetPlmnList
func (sim *searchIndexMapper) indexedTargetPlmnList() bool {
	return len(sim.TargetPlmnList) != 0
}

//indexedDnn indexed Dnn
func (sim *searchIndexMapper) indexedDnn() bool {
	return sim.Dnn != ""
}

//indexedSmfServingArea indexed SmfServingArea
func (sim *searchIndexMapper) indexedSmfServingArea() bool {
	return sim.SmfServingArea != ""
}

//indexedSupportedFeatures indexed SupportedFeature
func (sim *searchIndexMapper) indexedSupportedFeatures() bool {
	return sim.SupportedFeature != ""
}

//indexedRoutingIndicator routingIndicator of sim
func (sim *searchIndexMapper) indexedRoutingIndicator() bool {
	return sim.RoutingIndicator != ""
}

//indexedNsiList nsiList of sim
func (sim *searchIndexMapper) indexedNsiList() bool {
	return sim.NsiList != ""
}

//indexedIPDomain ipDomain of sim
func (sim *searchIndexMapper) indexedIPDomain() bool {
	return sim.IPDomain != ""
}

//indexedDnaiList dnaiList of sim
func (sim *searchIndexMapper) indexedDnaiList() bool {
	return sim.DnaiList != ""
}

//indexedUpfIwkEpsInd upfIwkEpsInd of sim
func (sim *searchIndexMapper) indexedUpfIwkEpsInd() bool {
	return sim.UpfIwkEpsInd != ""
}
