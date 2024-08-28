package structs

import (
	"fmt"
	"strings"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

//TargetNf get from TargetNfProfile
type TargetNf struct {
	RequesterNfType          string
	TargetNfType             string
	TargetServiceNames       []string
	NotifCondition           *NotifCondition
	SubscriptionValidityTime int
}

func (tn TargetNf) Info() string {
	serviceNames := strings.Join(tn.TargetServiceNames, ",")
	infoStr := fmt.Sprintf("requesterNftype[%s], targetNfType[%s], targetServiceNames[%s], subscriptionValidityTime[%d]", tn.RequesterNfType, tn.TargetNfType, serviceNames, tn.SubscriptionValidityTime)

	return infoStr
}

//OneSubscriptionData one subscribe data
type OneSubscriptionData struct {
	RequesterNfType   string
	TargetNfType      string
	TargetServiceName string
	NfInstanceID      string
	NotifCondition    *NotifCondition
}

//PatchItem patch item
type PatchItem struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	From  string      `json:"from,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

//Show is for show the content of TargetNf
func (tg *TargetNf) Show() {
	log.Infof("RequesterNfType[%s], TargetNfType[%s], support-services:%v\n", tg.RequesterNfType, tg.TargetNfType, tg.TargetServiceNames)
}

//NfTypeOfTarget ..
type NfTypeOfTarget struct {
	NfTypeOfTarget string `json:"nfTypeOfTarget"`
}

//SearchResult ..
type SearchResult struct {
	ValidityPeriod int32                   `json:"validityPeriod"`
	NfInstances    []SearchResultNFProfile `json:"nfInstances,omitempty"`
}

//SearchResultNFProfile ..
type SearchResultNFProfile struct {
	NfInstanceID         string                  `json:"nfInstanceId"`
	NfType               string                  `json:"nfType"`
	NfStatus             string                  `json:"nfStatus"`
	PLMN                 []PlmnID                `json:"plmnList,omitempty"`
	SNSSAI               []SNSSAI                `json:"sNssais,omitempty"`
	NsiList              []string                `json:"nsiList,omitempty"`
	FQDN                 string                  `json:"fqdn,omitempty"`
	Ipv4Addresses        []string                `json:"ipv4Addresses,omitempty"`
	Ipv6Addresses        []string                `json:"ipv6Addresses,omitempty"`
	Capacity             *int32                  `json:"capacity,omitempty"`
	Load                 *int32                  `json:"load,omitempty"`
	Locality             string                  `json:"locality,omitempty"`
	Priority             *int32                  `json:"priority,omitempty"`
	UdrInfo              *UDRINFO                `json:"udrInfo,omitempty"`
	UdmInfo              *UDMINFO                `json:"udmInfo,omitempty"`
	AusfInfo             *AUSFINFO               `json:"ausfInfo,omitempty"`
	AmfInfo              *AMFInfo                `json:"amfInfo,omitempty"`
	SmfInfo              *SMFInfo                `json:"smfInfo,omitempty"`
	UpfInfo              *UPFInfo                `json:"upfInfo,omitempty"`
	PcfInfo              *PCFInfo                `json:"pcfInfo,omitempty"`
	BsfInfo              *BSFInfo                `json:"bsfInfo,omitempty"`
	ChfInfo              *CHFInfo                `json:"chfInfo,omitempty"`
	CustomInfo           interface{}             `json:"customInfo,omitempty"`
	RecoveryTime         string                  `json:"recoveryTime,omitempty"`
	NfServicePersistence *bool                   `json:"nfServicePersistence,omitempty"`
	NfSrvList            []SearchResultNFService `json:"nfServices,omitempty"`
}

//SearchResultNFService ..
type SearchResultNFService struct {
	SrvID              string                        `json:"serviceInstanceId"`
	SrvName            string                        `json:"serviceName"`
	Versions           []NFServiceVersion            `json:"versions"`
	Scheme             string                        `json:"scheme"`
	NfSrvStatus        string                        `json:"nfServiceStatus"`
	FQDN               string                        `json:"fqdn,omitempty"`
	IPEndPoints        []IPEndPoint                  `json:"ipEndPoints,omitempty"`
	APIPrefix          string                        `json:"apiPrefix,omitempty"`
	DefaultNtfSubscrps []DefaultNtfSubscrp           `json:"defaultNotificationSubscriptions,omitempty"`
	Capacity           *int32                        `json:"capacity,omitempty"`
	Load               *int32                        `json:"load,omitempty"`
	Priority           *int32                        `json:"priority,omitempty"`
	RecoveryTime       string                        `json:"recoveryTime,omitempty"`
	ChfServiceInfo     *ChfServiceInfoChfServiceInfo `json:"chfServiceInfo,omitempty"`
	SupportedFeatures  string                        `json:"supportedFeatures,omitempty"`
}

//RegDiscInnerMsg for message between register-agent and discovery-agent
type RegDiscInnerMsg struct {
	EventType    string   `json:"eventType"`
	NfInstanceID string   `json:"nfInstanceId"`
	NfType       string   `json:"nfType,omitempty"`
	FQDN         string   `json:"fqdn,omitempty"`
	Plmns        []PlmnID `json:"plmns,omitempty"`
}

//DiscDiscInnerMsg for message between master-slave discovery-agent
type DiscDiscInnerMsg struct {
	EventType       string           `json:"eventType"`
	AgentProducerID string           `json:"agentProducerId"`
	SubscrInfo      SubscriptionInfo `json:"subscriptionInfo"`
}

type NfInfoForRegDisc struct {
	NfInstanceID    string   `json:"nfInstanceId"`
	RequesterNfFqdn string   `json:"requesterNfFqdn"`
	RequesterNfType string   `json:"requesterNfType"`
	RequesterPlmns  []PlmnID `json:"requesterPlmns,omitempty"`
}
