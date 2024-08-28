package nrfschema

type TUpfInfo struct {
	InterfaceUpfInfoList []TInterfaceUpfInfoItem `json:"interfaceUpfInfoList,omitempty"`
	SNssaiUpfInfoList    []TSnssaiUpfInfoItem    `json:"sNssaiUpfInfoList"`
	SmfServingArea       []string                `json:"smfServingArea,omitempty"`
	IwkEpsInd            *bool                   `json:"iwkEpsInd,omitempty"`
}
type TProblemDetails struct {
	Status        int             `json:"status,omitempty"`
	Title         string          `json:"title,omitempty"`
	Type          string          `json:"type,omitempty"`
	Cause         string          `json:"cause,omitempty"`
	Instance      string          `json:"instance,omitempty"`
	InvalidParams []TInvalidParam `json:"invalidParams,omitempty"`
}
type TInvalidParam struct {
	Param  string `json:"param"`
	Reason string `json:"reason,omitempty"`
}
type TDnnUpfInfoItem struct {
	Dnn      string   `json:"dnn"`
	DnaiList []string `json:"dnaiList,omitempty"`
}
type TTai struct {
	PlmnId TPlmnId `json:"plmnId"`
	Tac    string  `json:"tac"`
}
type TGuami struct {
	PlmnId TPlmnId `json:"plmnId"`
	AmfId  string  `json:"amfId"`
}
type TAusfInfo struct {
	GroupId           string       `json:"groupId,omitempty"`
	RoutingIndicators []string     `json:"routingIndicators,omitempty"`
	SupiRanges        []TSupiRange `json:"supiRanges,omitempty"`
}
type TSnssaiUpfInfoItem struct {
	DnnUpfInfoList []TDnnUpfInfoItem `json:"dnnUpfInfoList"`
	SNssai         TSnssai           `json:"sNssai"`
}
type TNFRegistrationData struct {
	NfProfile      TNFProfile `json:"nfProfile"`
	HeartBeatTimer int        `json:"heartBeatTimer"`
}
type TIdentityRange struct {
	Pattern string `json:"pattern,omitempty"`
	Start   string `json:"start,omitempty"`
	End     string `json:"end,omitempty"`
}
type TBsfInfo struct {
	DnnList           []string            `json:"dnnList,omitempty"`
	IpDomainList      []string            `json:"ipDomainList,omitempty"`
	Ipv6PrefixRanges  []TIpv6PrefixRange  `json:"ipv6PrefixRanges,omitempty"`
	Ipv4AddressRanges []TIpv4AddressRange `json:"ipv4AddressRanges,omitempty"`
}
type TChfInfo struct {
	SupiRangeList []TSupiRange     `json:"supiRangeList"`
	GpsiRangeList []TIdentityRange `json:"gpsiRangeList"`
	PlmnRangeList []TPlmnRange     `json:"plmnRangeList"`
}
type TPlmnRange struct {
	Start   string `json:"start,omitempty"`
	End     string `json:"end,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}
type TSupiRange struct {
	Start   string `json:"start,omitempty"`
	End     string `json:"end,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}
type TSmfInfo struct {
	SNssaiSmfInfoList []TSnssaiSmfInfoItem `json:"sNssaiSmfInfoList"`
	AccessType        []string             `json:"accessType,omitempty"`
	PgwFqdn           string               `json:"pgwFqdn,omitempty"`
	TaiList           []TTai               `json:"taiList,omitempty"`
	TaiRangeList      []TTaiRange          `json:"taiRangeList,omitempty"`
}
type TSnssaiSmfInfoItem struct {
	SNssai         TSnssai           `json:"sNssai"`
	DnnSmfInfoList []TDnnSmfInfoItem `json:"dnnSmfInfoList"`
}
type TDnnSmfInfoItem struct {
	Dnn string `json:"dnn"`
}
type TIpv6PrefixRange struct {
	End   string `json:"end,omitempty"`
	Start string `json:"start,omitempty"`
}
type TIpEndPoint struct {
	Ipv4Address string      `json:"ipv4Address,omitempty"`
	Ipv6Address string      `json:"ipv6Address,omitempty"`
	Port        *int        `json:"port,omitempty"`
	Transport   interface{} `json:"transport,omitempty"`
}
type TNotificationData struct {
	Event         interface{} `json:"event"`
	NewProfile    *TNFProfile `json:"newProfile,omitempty"`
	NfInstanceUri string      `json:"nfInstanceUri"`
}

type TNFProfileDB struct {
	ExpiredTime       uint64         `json:"expiredTime"`
	LastUpdateTime    uint64         `json:"lastUpdateTime"`
	ProfileUpdateTime uint64         `json:"profileUpdateTime"`
	Provisioned       int32          `json:"provisioned"`
	Md5sum            interface{}    `json:"md5sum"`
	Body              *TNFProfile    `json:"body"`
	OverrideInfo      []OverrideInfo `json:"overrideInfo,omitempty"`
	ProvSupiVersion   int64          `json:"provSupiVersion,omitempty"`
	ProvGpsiVersion   int64          `json:"provGpsiVersion,omitempty"`
}

type TNFProfile struct {
	NfInstanceId         string        `json:"nfInstanceId"`
	NfType               string        `json:"nfType"`
	NfStatus             string        `json:"nfStatus"`
	HeartBeatTimer       *int          `json:"heartBeatTimer,omitempty"`
	PlmnList             []TPlmnId     `json:"plmnList,omitempty"`
	SNssais              []TSnssai     `json:"sNssais,omitempty"`
	NsiList              []string      `json:"nsiList,omitempty"`
	Fqdn                 string        `json:"fqdn,omitempty"`
	InterPlmnFqdn        string        `json:"interPlmnFqdn,omitempty"`
	Ipv4Addresses        []string      `json:"ipv4Addresses,omitempty"`
	Ipv6Addresses        []string      `json:"ipv6Addresses,omitempty"`
	AllowedPlmns         []TPlmnId     `json:"allowedPlmns,omitempty"`
	AllowedNfTypes       []interface{} `json:"allowedNfTypes,omitempty"`
	AllowedNfDomains     []string      `json:"allowedNfDomains,omitempty"`
	AllowedNssais        []TSnssai     `json:"allowedNssais,omitempty"`
	Priority             *int          `json:"priority,omitempty"`
	Capacity             *int          `json:"capacity,omitempty"`
	Load                 int           `json:"load,omitempty"`
	Locality             string        `json:"locality,omitempty"`
	UdrInfo              *TUdrInfo     `json:"udrInfo,omitempty"`
	UdmInfo              *TUdmInfo     `json:"udmInfo,omitempty"`
	AusfInfo             *TAusfInfo    `json:"ausfInfo,omitempty"`
	AmfInfo              *TAmfInfo     `json:"amfInfo,omitempty"`
	SmfInfo              *TSmfInfo     `json:"smfInfo,omitempty"`
	UpfInfo              *TUpfInfo     `json:"upfInfo,omitempty"`
	PcfInfo              *TPcfInfo     `json:"pcfInfo,omitempty"`
	BsfInfo              *TBsfInfo     `json:"bsfInfo,omitempty"`
	ChfInfo              *TChfInfo     `json:"chfInfo,omitempty"`
	NrfInfo              *TNrfInfo     `json:"nrfInfo,omitempty"`
	CustomInfo           interface{}   `json:"customInfo,omitempty"`
	RecoveryTime         string        `json:"recoveryTime,omitempty"`
	NfServicePersistence *bool         `json:"nfServicePersistence,omitempty"`
	NfServices           []TNFService  `json:"nfServices,omitempty"`
	ProvisionInfo        *TProvInfo    `json:"provisionInfo,omitempty"`
}
type TInterfaceUpfInfoItem struct {
	EndpointFqdn          string      `json:"endpointFqdn,omitempty"`
	InterfaceType         interface{} `json:"interfaceType"`
	Ipv4EndpointAddresses []string    `json:"ipv4EndpointAddresses,omitempty"`
	Ipv6EndpointAddresses []string    `json:"ipv6EndpointAddresses,omitempty"`
	NetworkInstance       string      `json:"networkInstance,omitempty"`
}
type TUdrInfo struct {
	ExternalGroupIdentifiersRanges []TIdentityRange `json:"externalGroupIdentifiersRanges,omitempty"`
	GpsiRanges                     []TIdentityRange `json:"gpsiRanges,omitempty"`
	GroupId                        string           `json:"groupId,omitempty"`
	SupiRanges                     []TSupiRange     `json:"supiRanges,omitempty"`
	SupportedDataSets              []interface{}    `json:"supportedDataSets,omitempty"`
}
type TPcfInfo struct {
	DnnList     []string     `json:"dnnList,omitempty"`
	SupiRanges  []TSupiRange `json:"supiRanges,omitempty"`
	GroupId     string       `json:"groupId,omitempty"`
	RxDiamHost  string       `json:"rxDiamHost,omitempty"`
	RxDiamRealm string       `json:"rxDiamRealm,omitempty"`
}
type TIpv4AddressRange struct {
	End   string `json:"end,omitempty"`
	Start string `json:"start,omitempty"`
}
type TAmfInfo struct {
	BackupInfoAmfFailure []TGuami             `json:"backupInfoAmfFailure,omitempty"`
	BackupInfoAmfRemoval []TGuami             `json:"backupInfoAmfRemoval,omitempty"`
	GuamiList            []TGuami             `json:"guamiList"`
	TaiList              []TTai               `json:"taiList,omitempty"`
	TaiRangeList         []TTaiRange          `json:"taiRangeList,omitempty"`
	AmfRegionId          string               `json:"amfRegionId"`
	AmfSetId             string               `json:"amfSetId"`
	N2InterfaceAmfInfo   *TN2InterfaceAmfInfo `json:"n2InterfaceAmfInfo,omitempty"`
}

//TTaiRange is used to define the structure of TaiRange
type TTaiRange struct {
	PlmnID       *TPlmnId    `json:"plmnId"`
	TacRangeList []TTacRange `json:"tacRangeList"`
}

//TTacRange is used to define the structure of TacRange
type TTacRange struct {
	Start   string `json:"start,omitempty"`
	End     string `json:"end,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}

//TN2InterfaceAmfInfo is used to define the structure of N2InterfaceAmfInfo
type TN2InterfaceAmfInfo struct {
	Ipv4EndpointAddress []string `json:"ipv4EndpointAddress,omitempty"`
	Ipv6EndpointAddress []string `json:"ipv6EndpointAddress,omitempty"`
	AmfName             string   `json:"amfName,omitempty"`
}

type TNFServiceVersion struct {
	Expiry          string `json:"expiry,omitempty"`
	ApiFullVersion  string `json:"apiFullVersion"`
	ApiVersionInUri string `json:"apiVersionInUri"`
}
type TPlmnId struct {
	Mnc string `json:"mnc"`
	Mcc string `json:"mcc"`
}
type TNFService struct {
	Load                             int                                `json:"load,omitempty"`
	SupportedFeatures                string                             `json:"supportedFeatures,omitempty"`
	ServiceInstanceId                string                             `json:"serviceInstanceId"`
	ServiceName                      string                             `json:"serviceName"`
	Versions                         []TNFServiceVersion                `json:"versions"`
	Scheme                           string                             `json:"scheme"`
	NfServiceStatus                  string                             `json:"nfServiceStatus"`
	Fqdn                             string                             `json:"fqdn,omitempty"`
	ApiPrefix                        string                             `json:"apiPrefix,omitempty"`
	AllowedNfTypes                   []interface{}                      `json:"allowedNfTypes,omitempty"`
	Capacity                         *int                               `json:"capacity,omitempty"`
	AllowedNssais                    []TSnssai                          `json:"allowedNssais,omitempty"`
	IpEndPoints                      []TIpEndPoint                      `json:"ipEndPoints,omitempty"`
	AllowedNfDomains                 []string                           `json:"allowedNfDomains,omitempty"`
	AllowedPlmns                     []TPlmnId                          `json:"allowedPlmns,omitempty"`
	DefaultNotificationSubscriptions []TDefaultNotificationSubscription `json:"defaultNotificationSubscriptions,omitempty"`
	InterPlmnFqdn                    string                             `json:"interPlmnFqdn,omitempty"`
	RecoveryTime                     string                             `json:"recoveryTime,omitempty"`
	Priority                         *int                               `json:"priority,omitempty"`
	ChfServiceInfo                   *TChfServiceInfo                   `json:"chfServiceInfo,omitempty"`
}
type TDefaultNotificationSubscription struct {
	NotificationType   interface{} `json:"notificationType"`
	CallbackUri        string      `json:"callbackUri"`
	N1MessageClass     interface{} `json:"n1MessageClass,omitempty"`
	N2InformationClass interface{} `json:"n2InformationClass,omitempty"`
}
type TPatchItem struct {
	Value interface{} `json:"value,omitempty"`
	From  string      `json:"from,omitempty"`
	Op    interface{} `json:"op"`
	Path  string      `json:"path"`
}
type TSnssai struct {
	Sst int    `json:"sst"`
	Sd  string `json:"sd,omitempty"`
}
type TUdmInfo struct {
	SupiRanges                     []TSupiRange     `json:"supiRanges,omitempty"`
	ExternalGroupIdentifiersRanges []TIdentityRange `json:"externalGroupIdentifiersRanges,omitempty"`
	GpsiRanges                     []TIdentityRange `json:"gpsiRanges,omitempty"`
	GroupId                        string           `json:"groupId,omitempty"`
	RoutingIndicators              []string         `json:"routingIndicators,omitempty"`
}

//TSubscriptionData is used to define the structure of SubscriptionData
type TSubscriptionData struct {
	NfStatusNotificationUri string           `json:"nfStatusNotificationUri"`
	SubscrCond              *TSubscrCond     `json:"subscrCond,omitempty"`
	SubscriptionID          string           `json:"subscriptionId,omitempty"`
	ValidityTime            string           `json:"validityTime,omitempty"`
	ReqNotifEvents          []string         `json:"reqNotifEvents,omitempty"`
	PlmnId                  *TPlmnId         `json:"plmnId,omitempty"`
	NotifCondition          *TNotifCondition `json:"notifCondition,omitempty"`
	ReqNfType               string           `json:"reqNfType,omitempty"`
	ReqNfFqdn               string           `json:"reqNfFqdn,omitempty"`
}

//TSubscrCond is used to define the structure of SubscrCond
type TSubscrCond struct {
	NfInstanceID string    `json:"nfInstanceId,omitempty"`
	NfType       string    `json:"nfType,omitempty"`
	ServiceName  string    `json:"serviceName,omitempty"`
	AmfSetID     string    `json:"amfSetId,omitempty"`
	AmfRegionID  string    `json:"amfRegionId,omitempty"`
	GuamiList    []TGuami  `json:"guamiList,omitempty"`
	SnssaiList   []TSnssai `json:"snssaiList,omitempty"`
	NsiList      []string  `json:"nsiList,omitempty"`
	NfGroupID    string    `json:"nfGroupId,omitempty"`
}

//TNotifCondition is used to define the structure of NotifCondition
type TNotifCondition struct {
	MonitoredAttributes   []string `json:"monitoredAttributes,omitempty"`
	UnmonitoredAttributes []string `json:"unmonitoredAttributes,omitempty"`
}

type TNrfInfo struct {
	AmfInfoSum  *TAmfInfoSum  `json:"amfInfoSum,omitempty"`
	SmfInfoSum  *TSmfInfoSum  `json:"smfInfoSum,omitempty"`
	UdmInfoSum  *TUdmInfoSum  `json:"udmInfoSum,omitempty"`
	AusfInfoSum *TAusfInfoSum `json:"ausfInfoSum,omitempty"`
	PcfInfoSum  *TPcfInfoSum  `json:"pcfInfoSum,omitempty"`
}

type TAmfInfoSum struct {
	GuamiList       []TGuami    `json:"guamiList"`
	TaiList         []TTai      `json:"taiList,omitempty"`
	TaiRangeList    []TTaiRange `json:"taiRangeList,omitempty"`
	AmfRegionIdList []string    `json:"amfRegionIdList"`
	AmfSetIdList    []string    `json:"amfSetIdList"`
}

type TAusfInfoSum struct {
	GroupIdList          []string     `json:"groupIdList,omitempty"`
	RoutingIndicatorList []string     `json:"routingIndicatorList,omitempty"`
	SupiRanges           []TSupiRange `json:"supiRanges,omitempty"`
}

type TSmfInfoSum struct {
	DnnList      []string    `json:"dnnList"`
	PgwFqdnList  []string    `json:"pgwFqdnList,omitempty"`
	TaiList      []TTai      `json:"taiList,omitempty"`
	TaiRangeList []TTaiRange `json:"taiRangeList,omitempty"`
}

type TUdmInfoSum struct {
	SupiRanges                       []TSupiRange     `json:"supiRanges,omitempty"`
	ExternalGroupIdentityfiersRanges []TIdentityRange `json:"externalGroupIdentityfiersRanges,omitempty"`
	GpsiRanges                       []TIdentityRange `json:"gpsiRanges,omitempty"`
	GroupIdList                      []string         `json:"groupIdList,omitempty"`
	RoutingIndicatorList             []string         `json:"routingIndicatorList,omitempty"`
}

type TPcfInfoSum struct {
	DnnList    []string     `json:"dnnList,omitempty"`
	SupiRanges []TSupiRange `json:"supiRanges,omitempty"`
}

type TChfServiceInfo struct {
	PrimaryChfServiceInstance   string `json:"primaryChfServiceInstance,omitempty"`
	SecondaryChfServiceInstance string `json:"secondaryChfServiceInstance,omitempty"`
}

// NrfInfoPatchData is a patch item to update nrfInfo
type NrfInfoPatchData struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	From  string      `json:"from,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// NfProfilePatchData is a patch item to notify nfProfile changes
type NfProfilePatchData struct {
	Op        string      `json:"op"`
	Path      string      `json:"path"`
	From      string      `json:"from,omitempty"`
	OrigValue interface{} `json:"origValue,omitempty"`
	NewValue  interface{} `json:"newValue,omitempty"`
}

// TLink is to define the link object
type TLink struct {
	Href string `json:"href"`
}

// TLinks is to define the _links object, it has two members whose names are item and self.
type TLinks struct {
	Self TLink   `json:"self"`
	Item []TLink `json:"item"`
}

// TNfInstancesGetResponse is to define the get instances response
type TNfInstancesGetResponse struct {
	Links TLinks `json:"_links"`
}

//TProvInfo is to store the information of ProvisionInfo
type TProvInfo struct {
	CreateMode       string   `json:"createMode"`
	OverrideAttrList []string `json:"overrideAttrList,omitempty"`
}

// OverrideInfo is to store the override infomation from provision
type OverrideInfo struct {
	Path   string `json:"path"`
	Action string `json:"action"`
	Value  string `json:"value,omitempty"`
}

//GroupProfile supi group profile info
type GroupProfile struct {
	GroupProfileID string      `json:"groupProfileId"`
	NfType         []string    `json:"nfType"`
	SupiRanges     []SupiRange `json:"supiRanges"`
	GroupID        string      `json:"groupId"`
}

//GroupProfileReq supi group profile info for request message
type GroupProfileReq struct {
	GroupProfileID string      `json:"groupProfileId,omitempty"`
	NfType         []string    `json:"nfType"`
	SupiRanges     []SupiRange `json:"supiRanges"`
	GroupID        string      `json:"groupId"`
}

//SupiRange identity range info
type SupiRange struct {
	Start   string `json:"start,omitempty"`
	End     string `json:"end,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}

//GpsiProfile gpsi profile info
type GpsiProfile struct {
	GpsiProfileID string      `json:"gpsiProfileId"`
	NfType        []string    `json:"nfType"`
	GpsiRanges    []GpsiRange `json:"gpsiRanges"`
	GroupID       string      `json:"groupId"`
}

//GpsiProfileReq gpsi profile info for request message
type GpsiProfileReq struct {
	GpsiProfileID string      `json:"gpsiProfileId,omitempty"`
	NfType        []string    `json:"nfType"`
	GpsiRanges    []GpsiRange `json:"gpsiRanges"`
	GroupID       string      `json:"groupId"`
}

//GpsiRange identity range info
type GpsiRange struct {
	Start   string `json:"start,omitempty"`
	End     string `json:"end,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}
