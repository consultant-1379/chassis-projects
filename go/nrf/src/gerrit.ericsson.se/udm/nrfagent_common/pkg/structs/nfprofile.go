/*
package structs version:29510-f00
*/
package structs

import (
	"time"
)

//SNSSAI data structure
type SNSSAI struct {
	SST int32  `json:"sst"`
	SD  string `json:"sd,omitempty"`
}

//SupiRange
type SupiRange struct {
	Start   string `json:"start,omitempty"`
	End     string `json:"end,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}

//IdentityRange definition
type IdentityRange struct {
	Start   string `json:"start,omitempty"`
	End     string `json:"end,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}

//PlmnRange definition
type PlmnRange struct {
	Start   string `json:"start,omitempty"`
	End     string `json:"end,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}

//UDRINFO definition
type UDRINFO struct {
	GroupID            string          `json:"groupId,omitempty"`
	SupiRanges         []SupiRange     `json:"supiRanges,omitempty"`
	GpsiRanges         []IdentityRange `json:"gpsiRanges,omitempty"`
	ExtGroupIdfsRanges []IdentityRange `json:"externalGroupIdentifiersRanges,omitempty"`
	SupportedDataSets  []string        `json:"supportedDataSets,omitempty"`
}

//UDMINFO definition
type UDMINFO struct {
	GroupID            string          `json:"groupId,omitempty"`
	SupiRanges         []SupiRange     `json:"supiRanges,omitempty"`
	GpsiRanges         []IdentityRange `json:"gpsiRanges,omitempty"`
	ExtGroupIdfsRanges []IdentityRange `json:"externalGroupIdentifiersRanges,omitempty"`
	RoutingIndicators  []string        `json:"routingIndicators,omitempty"`
}

//UdmInfoSum definition
type UdmInfoSum struct {
	GroupIDList          []string        `json:"groupIdList,omitempty"`
	SupiRanges           []SupiRange     `json:"supiRanges,omitempty"`
	GpsiRanges           []IdentityRange `json:"gpsiRanges,omitempty"`
	ExtGroupIdfsRanges   []IdentityRange `json:"externalGroupIdentityfiersRanges,omitempty"`
	RoutingIndicatorList []string        `json:"routingIndicatorList,omitempty"`
}

//AUSFINFO definition
type AUSFINFO struct {
	GroupID           string      `json:"groupId,omitempty"`
	SupiRanges        []SupiRange `json:"supiRanges,omitempty"`
	RoutingIndicators []string    `json:"routingIndicators,omitempty"`
}

//AusfInfoSum definition
type AusfInfoSum struct {
	GroupIDList          string      `json:"groupIdList,omitempty"`
	SupiRanges           []SupiRange `json:"supiRanges,omitempty"`
	RoutingIndicatorList []string    `json:"routingIndicatorList,omitempty"`
}

//IPEndPoint definition
type IPEndPoint struct {
	Ipv4Address string `json:"ipv4Address,omitempty"`
	Ipv6Address string `json:"ipv6Address,omitempty"`
	Transport   string `json:"transport,omitempty"`
	Port        int32  `json:"port,omitempty"`
}

/*
//SubscriptionData defined in RFC29510
type SubscriptionData struct {
	NfStatusNotificationURI string   `json:"nfStatusNotificationUri"`
	SubscriptionID          string   `json:"subscriptionId,omitempty"`
	ValidityTime            string   `json:"validityTime,omitempty"`
	ReqNotifEvents          []string `json:"reqNotifEvents,omitempty"`
	PlmnID                  *PlmnID  `json:"plmnId,omitempty"`
	NfInstanceID            string   `json:"nfInstanceId,omitempty"`
	NfType                  string   `json:"nfType,omitempty"`
	ServiceName             string   `json:"serviceName,omitempty"`
	AmfSetID                string   `json:"amfSetId,omitempty"`
	AmfRegionID             string   `json:"amfRegionId,omitempty"`
	GuamiList               []Guami  `json:"guamiList,omitempty"`
}
*/

//SubscriptionData defined in RFC29510
type SubscriptionData struct {
	NfStatusNotificationURI string          `json:"nfStatusNotificationUri"`
	SubscrCond              interface{}     `json:"subscrCond,omitempty"`
	SubscriptionID          string          `json:"subscriptionId,omitempty"`
	ValidityTime            string          `json:"validityTime,omitempty"`
	ReqNotifEvents          []string        `json:"reqNotifEvents,omitempty"`
	PlmnID                  *PlmnID         `json:"plmnId,omitempty"`
	NotifCondition          *NotifCondition `json:"notifCondition,omitempty"`
	ReqNfType               string          `json:"reqNfType,omitempty"`
	ReqNfFqdn               string          `json:"reqNfFqdn,omitempty"`
}

//NotifCondition notifCondition subscription info
type NotifCondition struct {
	MonitoredAttributes   []string `json:"monitoredAttributes,omitempty"`
	UnmonitoredAttributes []string `json:"unmonitoredAttributes,omitempty"`
}

//NfInstanceIDCond inInstanceId subscription info
type NfInstanceIDCond struct {
	NfInstanceID string `json:"nfInstanceId"`
}

//NfTypeCond nfType subscription info
type NfTypeCond struct {
	NfType string `json:"nfType"`
}

//ServiceNameCond serviceName subscription info
type ServiceNameCond struct {
	ServiceName string `json:"serviceName"`
}

//AmfCond amfCond subscription info
type AmfCond struct {
	AmfSetID    string `json:"amfSetId,omitempty"`
	AmfRegionID string `json:"amfRegionId,omitempty"`
}

//GuamiListCond guamiList subscription info
type GuamiListCond struct {
	GuamiList []Guami `json:"guamiList"`
}

//NetworkSliceCond networkSlice subscription info
type NetworkSliceCond struct {
	SnssaiList []SNSSAI `json:"snssaiList"`
	NsiList    []string `json:"nsiList,omitempty"`
}

//NfGroupCond nfGroup subscription info
type NfGroupCond struct {
	NfType    string `json:"nfType"`
	NfGroupID string `json:"nfGroupId"`
}

//DefaultNtfSubscrp definition
type DefaultNtfSubscrp struct {
	NotificationType string `json:"notificationType"`
	CallbackURI      string `json:"callbackUri"`
	N1MsgClass       string `json:"n1MessageClass,omitempty"`
	N2InfoClass      string `json:"n2InformationClass,omitempty"`
}

//PlmnID
type PlmnID struct {
	Mcc string `json:"mcc"`
	Mnc string `json:"mnc"`
}

//Guami definition
type Guami struct {
	PLMNID *PlmnID `json:"plmnId"`
	AMFID  string  `json:"amfId"`
}

//Tai definition
type Tai struct {
	PLMNID *PlmnID `json:"plmnId"`
	Tac    string  `json:"tac"` //      pattern: '(^[A-Fa-f0-9]{4}$)|(^[A-Fa-f0-9]{6}$)'
}

//TaiRange definition
type TaiRange struct {
	PLMNID       *PlmnID    `json:"plmnId"`
	TacRangeList []TacRange `json:"tacRangeList"` //      pattern: '(^[A-Fa-f0-9]{4}$)|(^[A-Fa-f0-9]{6}$)'
}

//TacRange definition
type TacRange struct {
	Start   string `json:"start,omitempty"`
	End     string `json:"end,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}

//NFServiceVersion definition
type NFServiceVersion struct {
	APIVersionInUrI string `json:"apiVersionInUri"`
	APIFullVersion  string `json:"apiFullVersion"`
	Expiry          string `json:"expiry,omitempty"`
}

//ChfServiceInfoChfServiceInfo  definition
type ChfServiceInfoChfServiceInfo struct {
	PrimaryChfServiceInstance   []string `json:"primaryChfServiceInstance,omitempty"`
	SecondaryChfServiceInstance []string `json:"secondaryChfServiceInstance,omitempty"`
}

//NFService definition
//TODO: nfServiceStatus is mandatary in 15.1.0, please remove the property omitempty in the later version
type NFService struct {
	SrvID              string                        `json:"serviceInstanceId"`
	SrvName            string                        `json:"serviceName"`
	Versions           []NFServiceVersion            `json:"versions"`
	Scheme             string                        `json:"scheme"`
	NfSrvStatus        string                        `json:"nfServiceStatus"`
	FQDN               string                        `json:"fqdn,omitempty"`
	InterPlmnFqdn      string                        `json:"interPlmnFqdn,omitempty"`
	IPEndPoints        []IPEndPoint                  `json:"ipEndPoints,omitempty"`
	APIPrefix          string                        `json:"apiPrefix,omitempty"`
	DefaultNtfSubscrps []DefaultNtfSubscrp           `json:"defaultNotificationSubscriptions,omitempty"`
	AllowedPLMNs       []PlmnID                      `json:"allowedPlmns,omitempty"`
	AllowedNFTypes     []string                      `json:"allowedNfTypes,omitempty"`
	AllowedDomains     []string                      `json:"allowedNfDomains,omitempty"`
	AllowedNSSAIs      []SNSSAI                      `json:"allowedNssais,omitempty"`
	Priority           *int32                        `json:"priority,omitempty"`
	Capacity           *int32                        `json:"capacity,omitempty"`
	Load               *int32                        `json:"load,omitempty"`
	RecoveryTime       string                        `json:"recoveryTime,omitempty"`
	ChfServiceInfo     *ChfServiceInfoChfServiceInfo `json:"chfServiceInfo,omitempty"`
	SupportedFeatures  string                        `json:"supportedFeatures,omitempty"`
}

//AMFInfo definition
type AMFInfo struct {
	AMFSetID           string              `json:"amfSetId"`
	AMFRegionID        string              `json:"amfRegionId"`
	GuamiList          []Guami             `json:"guamiList"`
	TaiList            []Tai               `json:"taiList,omitempty"`
	TaiRangeList       []TaiRange          `json:"taiRangeList,omitempty"`
	BcpInfoAmfFailure  []Guami             `json:"backupInfoAmfFailure,omitempty"`
	BcpInfoAmfRemoval  []Guami             `json:"backupInfoAmfRemoval,omitempty"`
	N2InterfaceAmfInfo *N2InterfaceAmfInfo `json:"n2InterfaceAmfInfo,omitempty"`
}

//AmfInfoSum ..
type AmfInfoSum struct {
	AmfSetIDList    []string   `json:"amfSetIdList"`
	AmfRegionIDList []string   `json:"amfRegionIdList"`
	GuamiList       []Guami    `json:"guamiList"`
	TaiList         []Tai      `json:"taiList,omitempty"`
	TaiRangeList    []TacRange `json:"taiRangeList,omitempty"`
}

//DnnSmfInfoItem  definition
type DnnSmfInfoItem struct {
	Dnn string `json:"dnn"`
}

//SnssaiSmfInfoItem definition
type SnssaiSmfInfoItem struct {
	SNssai         *SNSSAI          `json:"sNssai"`
	DnnSmfInfoList []DnnSmfInfoItem `json:"dnnSmfInfoList"`
}

//SMFInfo ..
type SMFInfo struct {
	SNssaiSmfInfoList []SnssaiSmfInfoItem `json:"sNssaiSmfInfoList"`
	TaiList           []Tai               `json:"taiList,omitempty"`
	TaiRangeList      []TaiRange          `json:"taiRangeList,omitempty"`
	PgwFqdn           string              `json:"pgwFqdn,omitempty"`
	AccessType        []string            `json:"accessType,omitempty"`
}

//SmfInfoSum ..
type SmfInfoSum struct {
	DnnList      []string   `json:"dnnList"`
	TaiList      []Tai      `json:"taiList,omitempty"`
	TaiRangeList []TaiRange `json:"taiRangeList,omitempty"`
	PgwFqdnList  []string   `json:"pgwFqdnList,omitempty"`
}

//UPFInfo ..
type UPFInfo struct {
	SNssaiUpfInfoList    []SnssaiUpfInfoItem    `json:"sNssaiUpfInfoList"`
	SmfServingArea       []string               `json:"smfServingArea,omitempty"`
	InterfaceUpfInfoList []InterfaceUpfInfoItem `json:"interfaceUpfInfoList,omitempty"`
	IwkEpsInd            *bool                  `json:"iwkEpsInd,omitempty"`
}

//PCFInfo definition
type PCFInfo struct {
	Dnnlist     []string    `json:"dnnList,omitempty"`
	SupiRanges  []SupiRange `json:"supiRanges,omitempty"`
	RxDiamHost  string      `json:"rxDiamHost,omitempty"`
	RxDiamRealm string      `json:"rxDiamRealm,omitempty"`
}

//PcfInfoSum definition
type PcfInfoSum struct {
	Dnnlist    []string    `json:"dnnList"`
	SupiRanges []SupiRange `json:"supiRanges,omitempty"`
}

//BSFInfo definition
type BSFInfo struct {
	Dnnlist           []string           `json:"dnnList,omitempty"`
	IPDomainList      []string           `json:"ipDomainList,omitempty"`
	IPv4AddressRanges []Ipv4AddressRange `json:"ipv4AddressRanges,omitempty"`
	IPv6PrefixRanges  []Ipv6PrefixRange  `json:"ipv6PrefixRanges,omitempty"`
}

//CHFInfo definition
type CHFInfo struct {
	SupiRangeList []SupiRange     `json:"supiRangeList,omitempty"`
	GpsiRangeList []IdentityRange `json:"gpsiRangeList,omitempty"`
	PlmnRangeList []PlmnRange     `json:"plmnRangeList,omitempty"`
}

//NRFInfo definition 29510f20
//type NRFInfo struct {
//	ServedUdrInfo  *UDRINFO  `json:"servedUdrInfo,omitempty"`
//	ServedUdmInfo  *UDMINFO  `json:"servedUdmInfo,omitempty"`
//	ServedAusfInfo *AUSFINFO `json:"servedAusfInfo,omitempty"`
//	ServedAmfInfo  *AMFInfo  `json:"servedAmfInfo,omitempty"`
//	ServedSmffInfo *SMFInfo  `json:"servedSmfInfo,omitempty"`
//	ServedUpfInfo  *UPFInfo  `json:"servedUpfInfo,omitempty"`
//	ServedPcfInfo  *PCFInfo  `json:"servedPcfInfo,omitempty"`
//	ServedBsfInfo  *BSFInfo  `json:"servedBsfInfo,omitempty"`
//	ServedChfInfo  *CHFInfo  `json:"servedChfInfo,omitempty"`
//}

//NRFInfo definition
type NRFInfo struct {
	AmfInfoSum  *AmfInfoSum  `json:"amfInfoSum,omitempty"`
	SmfInfoSum  *SmfInfoSum  `json:"smfInfoSum,omitempty"`
	UdmInfoSum  *UdmInfoSum  `json:"udmInfoSum,omitempty"`
	AusfInfoSum *AusfInfoSum `json:"ausfInfoSum,omitempty"`
	PcfInfoSum  *PcfInfoSum  `json:"pcfInfoSum,omitempty"`
}

//CustomInfo definition
type CustomInfo struct {
}

//Ipv4AddressRange definition
type Ipv4AddressRange struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

//Ipv6PrefixRange definition
type Ipv6PrefixRange struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

//InterfaceUpfInfoItem definition
type InterfaceUpfInfoItem struct {
	InterfaceType         string   `json:"interfaceType"`
	IPv4EndpointAddresses []string `json:"ipv4EndpointAddresses,omitempty"`
	IPv6EndpointAddresses []string `json:"ipv6EndpointAddresses,omitempty"`
	EndpointFqdn          string   `json:"endpointFqdn,omitempty"`
	NetworkInstance       string   `json:"networkInstance,omitempty"`
}

//N2InterfaceAmfInfo definition
type N2InterfaceAmfInfo struct {
	IPv4EndpointAddress []string `json:"ipv4EndpointAddress,omitempty"`
	IPv6EndpointAddress []string `json:"ipv6EndpointAddress,omitempty"`
	AmfName             string   `json:"amfName,omitempty"`
}

//SnssaiUpfInfoItem definition
type SnssaiUpfInfoItem struct {
	SNssai         *SNSSAI          `json:"sNssai"`
	DNNUpfInfoList []DnnUpfInfoItem `json:"dnnUpfInfoList"`
}

//DnnUpfInfoItem definition
type DnnUpfInfoItem struct {
	DNN      string   `json:"dnn"`
	DnaiList []string `json:"dnaiList,omitempty"`
}

//NfProfile definition
type NfProfile struct {
	NfInstanceID         string      `json:"nfInstanceId"`
	NfType               string      `json:"nfType"`
	NfStatus             string      `json:"nfStatus"`
	HeartBeatTimer       *int32      `json:"heartBeatTimer,omitempty"`
	PlmnList             []PlmnID    `json:"plmnList,omitempty"`
	SNssais              []SNSSAI    `json:"sNssais,omitempty"`
	NsiList              []string    `json:"nsiList,omitempty"`
	FQDN                 string      `json:"fqdn,omitempty"`
	InterPlmnFqdn        string      `json:"interPlmnFqdn,omitempty"`
	Ipv4Addresses        []string    `json:"ipv4Addresses,omitempty"`
	Ipv6Addresses        []string    `json:"ipv6Addresses,omitempty"`
	AllowedPLMNs         []PlmnID    `json:"allowedPlmns,omitempty"`
	AllowedNFTypes       []string    `json:"allowedNfTypes,omitempty"`
	AllowedDomains       []string    `json:"allowedNfDomains,omitempty"`
	AllowedNSSAIs        []SNSSAI    `json:"allowedNssais,omitempty"`
	Priority             *int32      `json:"priority,omitempty"`
	Capacity             *int32      `json:"capacity,omitempty"`
	Load                 *int32      `json:"load,omitempty"`
	Locality             string      `json:"locality,omitempty"`
	UdrInfo              *UDRINFO    `json:"udrInfo,omitempty"`
	UdmInfo              *UDMINFO    `json:"udmInfo,omitempty"`
	AusfInfo             *AUSFINFO   `json:"ausfInfo,omitempty"`
	AmfInfo              *AMFInfo    `json:"amfInfo,omitempty"`
	SmfInfo              *SMFInfo    `json:"smfInfo,omitempty"`
	UpfInfo              *UPFInfo    `json:"upfInfo,omitempty"`
	PcfInfo              *PCFInfo    `json:"pcfInfo,omitempty"`
	BsfInfo              *BSFInfo    `json:"bsfInfo,omitempty"`
	ChfInfo              *CHFInfo    `json:"chfInfo,omitempty"`
	NrfInfo              *NRFInfo    `json:"nrfInfo,omitempty"`
	CustomInfo           interface{} `json:"customInfo,omitempty"`
	RecoveryTime         string      `json:"recoveryTime,omitempty"`
	NfServicePersistence *bool       `json:"nfServicePersistence,omitempty"`
	NfSrvList            []NFService `json:"nfServices,omitempty"`
}

//NotificationData definition
type NotificationData struct {
	NfEvent        string               `json:"event"`
	NfInstanceURI  string               `json:"nfInstanceUri,omitempty"`
	NfProfile      *NfProfile           `json:"nfProfile,omitempty"`
	ProfileChanges []NfProfilePatchData `json:"profileChanges,omitempty"`
}

//NfProfilePatchData is a patch item to notify nfProfile changes
type NfProfilePatchData struct {
	Op        string      `json:"op"`
	Path      string      `json:"path"`
	From      string      `json:"from,omitempty"`
	OrigValue interface{} `json:"origValue,omitempty"`
	NewValue  interface{} `json:"newValue,omitempty"`
}

//NfProfilePatchApplyData is a patch item to apply patch on nfProfile
type NfProfilePatchApplyData struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	From  string      `json:"from,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

//NotificationMsg is for messageBus notify-agent notification message
type NotificationMsg struct {
	NfEvent         string        `json:"event"`
	NfType          string        `json:"nfType,omitempty"`
	NfInstanceID    string        `json:"nfInstanceId"`
	AgentProducerID string        `json:"agentProducerId"`
	MessageBody     *SearchResult `json:"messageBody,omitempty"`
}

//NtfDiscInnerMsg is for messageBus notify disc inner message
type NtfDiscInnerMsg struct {
	NfEvent         string               `json:"event"`
	NfType          string               `json:"nfType,omitempty"`
	ReqNfType       string               `json:"reqNfType,omitempty"`
	NfInstanceID    string               `json:"nfInstanceId"`
	AgentProducerID string               `json:"agentProducerId"`
	MessageBody     []NfProfilePatchData `json:"messageBody,omitempty"`
	NfProfile       *NfProfile           `json:"nfProfile,omitempty"`
}

//CacheSyncData is for master-slave sync cache
type CacheSyncData struct {
	RequestNfType     string          `json:"requestNfType"`
	CacheInfos        []CacheSyncInfo `json:"cacheInfos"`
	RoamingCacheInfos []CacheSyncInfo `json:"roamingCacheInfos"`
}

//CacheSyncInfo definition
type CacheSyncInfo struct {
	TargetNfType      string             `json:"targetNfType"`
	NfProfiles        [][]byte           `json:"nfProfiles"`
	TtlInfos          []TtlInfo          `json:"ttlInfos"`
	SubscriptionInfos []SubscriptionInfo `json:"subscriptionInfos"`
	EtagInfos         []EtagInfo         `json:"etagInfos,omitempty"`
}

//CacheDumpData is for dump cache for oam
type CacheDumpData struct {
	RequestNfType     string          `json:"requestNfType"`
	CacheInfos        []CacheDumpInfo `json:"cacheInfos"`
	RoamingCacheInfos []CacheDumpInfo `json:"roamingCacheInfos"`
}

//CacheDumpInfo definition
type CacheDumpInfo struct {
	TargetNfType      string             `json:"targetNfType"`
	NfProfiles        []string           `json:"nfProfiles"`
	TtlInfos          []TtlInfo          `json:"ttlInfos"`
	SubscriptionInfos []SubscriptionInfo `json:"subscriptionInfos"`
	EtagInfos         []EtagInfo         `json:"etagInfos,omitempty"`
}

//TtlInfo definition
type TtlInfo struct {
	NfInstanceID string    `json:"nfInstanceId"`
	ValidityTime time.Time `json:"validityTime"`
}

//EtagInfo definition
type EtagInfo struct {
	NfInstanceID string `json:"nfInstanceId"`
	FingerPrint  string `json:"fingerPrint"`
}
