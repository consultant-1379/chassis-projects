package constvalue

import (
	"time"
)

const (
	//Work mode
	APP_WORKMODE_NRF_MGMT = "nrf_mgmt"
	APP_WORKMODE_NRF_DISC = "nrf_disc"
	APP_WORKMODE_NRF_PROV = "nrf_prov"
	AppWorkmodeNrfNotif   = "nrf_notify"
	APP_WORKMODE_TEST     = "test"

	//Event for NT
	NF_EVENT_CREATED = "NF_REGISTERED"      //"configCreated"
	NF_EVENT_UPDATED = "NF_PROFILE_CHANGED" //"configUpdated"
	NF_EVENT_DELETED = "NF_DEREGISTERED"    //"configDeleted"
	// NFEventUpateLoad indicates only load of profileis changed
	NFEventUpateLoad = "NF_PROFILE_LOAD_CHANGE"

	NfTypeNRF   = "NRF"    //NfTypeNRF nrf nfType
	NfTypeUDM   = "UDM"    //NfTypeUDM udm nfType
	NfTypeAMF   = "AMF"    //NfTypeAMF amf nfType
	NfTypeSMF   = "SMF"    //NfTypeSMF smf nfType
	NfTypeAUSF  = "AUSF"   //NfTypeAUSF ausf nfType
	NfTypeNEF   = "NEF"    //NfTypeNEF nef nfType
	NfTypePCF   = "PCF"    //NfTypePCF pcf nfType
	NfTypeSMSF  = "SMSF"   //NfTypeSMSF smsf nfType
	NfTypeNSSF  = "NSSF"   //NfTypeNSSF nssf nfType
	NfTypeUDR   = "UDR"    //NfTypeUDR udr nfType
	NfTypeUPF   = "UPF"    //NfTypeUPF upf nfType
	NfTypeLMF   = "LMF"    //NfTypeLMF lmf nfType
	NfTypeGMLC  = "GMLC"   //NfTypeGMLC gmlc nfType
	NfType5GEIR = "5G-EIR" //NfType5GEIR 5geir nfType
	NfTypeSEPP  = "SEPP"   //NfTypeSEPP sepp nfType
	NfTypeN3IWF = "N3IWF"  //NfTypeN3IWF n3iwf nfType
	NfTypeAF    = "AF"     //NfTypeAF af nfType
	NfTypeUDSF  = "UDSF"   //NfTypeUDSF udsf nfType
	NfTypeBSF   = "BSF"    //NfTypeBSF bsf nfType
	NfTypeCHF   = "CHF"    //NfTypeCHF chf nfType
	NfTypeNWDAF = "NWDAF"  //NfTypeNWDAF nwdaf nfType

	// service name
	NNRFNFM                  = "nnrf-nfm"
	NNRFDISC                 = "nnrf-disc"
	NUDMSDM                  = "nudm-sdm"
	NUDMUECM                 = "nudm-uecm"
	NUDMUEAU                 = "nudm-ueau"
	NUDMEE                   = "nudm-ee"
	NUDMPP                   = "nudm-pp"
	NAMFCOMM                 = "namf-comm"
	NAMFEVTS                 = "namf-evts"
	NAMFMT                   = "namf-mt"
	NAMFLOC                  = "namf-loc"
	NSMFPDUSESSION           = "nsmf-pdusession"
	NSMFEVENTEXPOSURE        = "nsmf-event-exposure"
	NAUSFAUTH                = "nausf-auth"
	NAUSFSORPROTECTION       = "nausf-sorprotection"
	NNEFPFDMANAGEMENT        = "nnef-pfdmanagement"
	NPCFAMPOLICYCONTROL      = "npcf-am-policy-control"
	NPCFSMPOLICYCONTROL      = "npcf-smpolicycontrol"
	NPCFPOLICYAUTHORIZATION  = "npcf-policyauthorization"
	NPCFBDTPOLICYCONTROL     = "npcf-bdtpolicycontrol"
	NPCFEVENTEXPOSURE        = "npcf-eventexposure"
	NPCFUEPOLICYCONTROL      = "npcf-ue-policy-control"
	NSMSFSMS                 = "nsmsf-sms"
	NNSSFNSSELECTION         = "nnssf-nsselection"
	NNSSFNSSAIAVAILABILITY   = "nnssf-nssaiavailability"
	NUDRDR                   = "nudr-dr"
	NLMFLOC                  = "nlmf-loc"
	N5GEIREIC                = "n5g-eir-eic"
	NBSFMANAGEMENT           = "nbsf-management"
	NCHFSPENDINGLIMITCONTROL = "nchf-spendinglimitcontrol"
	NCHFCONVERGEDCHARGING    = "nchf-convergedcharging"
	NNWDAFEVENTSSUBSCRIPTION = "nnwdaf-eventssubscription"
	NNWDAFANALYTICSINFO      = "nnwdaf-analyticsinfo"

	//NF Status
	NFStatusRegistered   = "REGISTERED"
	NFStatusSuspended    = "SUSPENDED"
	NFStatusDeregistered = "DEREGISTERED"

	//enum AccessType
	Access3GPP    = "3GPP_ACCESS"
	NonAccess3GPP = "NON_3GPP_ACCESS"

	//NFAutoRegistered is auto registered
	NFAutoRegistered = "false"
	//NFManualRegistered is manual registered
	NFManualRegistered = "true"

	//Custom Info
	ExpiredTime       = "expiredTime"
	LastUpdateTime    = "lastUpdateTime"
	ProfileUpdateTime = "profileUpdateTime"
	ProvisionedFlag   = "provisioned"
	MD5SUM            = "md5sum"
	BODY              = "body"

	Common = "common"

	//NF Profile
	NfProfile        = "nfProfile"
	NfInstanceId     = "nfInstanceId"
	NfType           = "nfType"
	NfStatus         = "nfStatus"
	Plmn             = "plmn"
	PlmnList         = "plmnList"
	Mcc              = "mcc"
	Mnc              = "mnc"
	Snssais          = "sNssais"
	Sst              = "sst"
	Sd               = "sd"
	Fqdn             = "fqdn"
	InterPlmnFqdn    = "interPlmnFqdn"
	Ipv4Addresses    = "ipv4Addresses"
	Ipv6Addresses    = "ipv6Addresses"
	AllowedPlmns     = "allowedPlmns"
	AllowedNFTypes   = "allowedNfTypes"
	AllowedNfDomains = "allowedNfDomains"
	AllowedNssais    = "allowedNssais"
	Ipv6Prefixes     = "ipv6Prefixes"
	Capacity         = "capacity"
	UdrInfo          = "udrInfo"
	UdmInfo          = "udmInfo"
	AusfInfo         = "ausfInfo"
	AmfInfo          = "amfInfo"
	SmfInfo          = "smfInfo"
	UpfInfo          = "upfInfo"
	PcfInfo          = "pcfInfo"
	BsfInfo          = "bsfInfo"
	ChfInfo          = "chfInfo"
	NrfInfo          = "nrfInfo"
	NfServices       = "nfServices"
	OverrideInfo     = "overrideInfo"
	ProvisionInfo    = "provisionInfo"
	ProvSupiVersion  = "provSupiVersion"
	ProvGpsiVersion  = "provGpsiVersion"

	//HeartBeat
	HeartBeatTimer        = "heartBeatTimer"
	ValidityPeriodOfCache = "ValidityPeriod"
	GuardTime             = 5

	//NF Servcie
	NFServiceInstanceId       = "serviceInstanceId"
	NFServiceName             = "serviceName"
	NFServiceVersions         = "versions"
	NFServiceScheme           = "scheme"
	NFServiceFqdn             = "fqdn"
	NFServicePriority         = "priority"
	NFServiceInterPlmnFqdn    = "interPlmnFqdn"
	NFServiceApiPrefix        = "apiPrefix"
	NFServiceIPEndPoints      = "ipEndPoints"
	NFServiceDefNotifiSub     = "defaultNotificationSubscriptions"
	NFServiceAllowedPlmns     = "allowedPlmns"
	NFServiceAllowedNFTypes   = "allowedNfTypes"
	NFServiceAllowedDomains   = "allowedNfDomains"
	NFServiceAllowedNssais    = "allowedNssais"
	NFServiceCapacity         = "capacity"
	NFServiceStatus           = "nfServiceStatus"
	NFServiceStatusRegistered = "REGISTERED"
	NFServiceChfServiceInfo   = "chfServiceInfo"
	//IpEndPoint
	IPEndPointIpv4Address = "ipv4Address"
	IPEndPointIpv6Address = "ipv6Address"
	IPEndPointIpv6Prefix  = "ipv6Prefix"
	IPEndPointTransport   = "transport"
	IPEndPointPort        = "port"

	ApiVersionInUri = "apiVersionInUri"

	DefNotifiSubNotificationType   = "notificationType"
	DefNotifiSubCallbackURI        = "callbackUri"
	DefNotifiSubN1MessageClass     = "n1MessageClass"
	DefNotifiSubN2InformationClass = "n2InformationClass"

	//UdrInfo
	GroupID                          = "groupId"
	SupiRanges                       = "supiRanges"
	GpsiRanges                       = "gpsiRanges"
	ExternalGroupIdentityfiersRanges = "externalGroupIdentifiersRanges"
	SupportedDataSets                = "supportedDataSets"

	//Amf
	AmfSetID             = "amfSetId"
	AmfRegionID          = "amfRegionId"
	GuamiList            = "guamiList"
	TaiList              = "taiList"
	TaiRangeList         = "taiRangeList"
	BackupInfoAmfFailure = "backupInfoAmfFailure"
	BackupInfoAmfRemoval = "backupInfoAmfRemoval"
	N2InterfaceAmfInfo   = "n2InterfaceAmfInfo"

	// DnnList in PcfInfo
	DnnList = "dnnList"

	//AccessType in SmfInfo
	AccessType = "accessType"

	Locality = "locality"

	//SupiRangeList in ChfInfo
	SupiRangeList = "supiRangeList"
	GpsiRangeList = "gpsiRangeList"
	PlmnRangeList = "plmnRangeList"

	// TacRangeList in TaiRange
	TacRangeList = "tacRangeList"

	ServingArea       = "servingArea"
	PgwFqdn           = "pgwFqdn"
	SNssaiUpfInfoList = "sNssaiUpfInfoList"
	SNssai            = "sNssai"
	Dnn               = "dnn"
	DnnUpfInfoList    = "dnnUpfInfoList"
	DnaiList          = "dnaiList"
	IwkEpsInd         = "iwkEpsInd"
	RoutingIndicator  = "routingIndicator"

	IPDomainList      = "ipDomainList"
	IPv4AddressRanges = "ipv4AddressRanges"
	IPv6PrefixRanges  = "ipv6PrefixRanges"

	// SmfServingArea in upfInfo
	SmfServingArea = "smfServingArea"
	NsiList        = "nsiList"
	// InterfaceUpfInfoList in upfInfo
	InterfaceUpfInfoList = "interfaceUpfInfoList"

	// InterfaceType in InterfaceUpfInfoItem
	InterfaceType = "interfaceType"

	//Ipv4EndpointAddress in InterfaceUpfInfoItem
	Ipv4EndpointAddress = "ipv4EndpointAddress"

	//Ipv6EndpointAddress in InterfaceUpfInfoItem
	Ipv6EndpointAddress = "ipv6EndpointAddress"

	// Ipv6EndpointPrefix in InterfaceUpfInfoItem
	Ipv6EndpointPrefix = "ipv6EndpointPrefix"

	// EndpointFqdn in InterfaceUpfInfoItem
	EndpointFqdn = "endpointFqdn"

	// NetworkInstance in InterfaceUpfInfoItem
	NetworkInstance = "networkInstance"

	//PcfInfoSum in nrfInfo
	PcfInfoSum = "pcfInfoSum"
	//AusfInfoSum in nrfInfo
	AusfInfoSum = "ausfInfoSum"
	//UdmInfoSum in nrfInfo
	UdmInfoSum = "udmInfoSum"
	//SmfInfoSum in nrfInfo
	SmfInfoSum = "smfInfoSum"
	//AmfInfoSum in nrfInfo
	AmfInfoSum = "amfInfoSum"

	//GroupIDList groupIdList info in nrfInfo subItem
	GroupIDList = "groupIdList"
	//RoutingIndicatorList routing info in nrfInfo subItem
	RoutingIndicatorList = "routingIndicatorList"
	//PgwFqdnList fqdn into in nrfInfo subItem
	PgwFqdnList = "pgwFqdnList"
	//AmfSetIDList setIDList info in nrfInfo subItem
	AmfSetIDList = "amfSetIdList"
	//AmfRegionIDList regionIDList info in nrfInfo subItem
	AmfRegionIDList = "amfRegionIdList"

	//SearchData
	SearchDataServiceName         = "service-names"
	SearchDataTargetNfType        = "target-nf-type"
	SearchDataTargetInstID        = "target-nf-instance-id"
	SearchDataRequesterNfType     = "requester-nf-type"
	SearchDataTargetPlmnList      = "target-plmn-list"
	SearchDataRequesterPlmnList   = "requester-plmn-list"
	SearchDataMcc                 = "mcc"
	SearchDataMnc                 = "mnc"
	SearchDataSnssais             = "snssais"
	SearchDataSnssaiSst           = "sst"
	SearchDataSnssaiSd            = "sd"
	SearchDataDnn                 = "dnn"
	SearchDataSmfServingArea      = "smf-serving-area"
	SearchDataTai                 = "tai"
	SearchDataDataEcgi            = "ecgi"
	SearchNcgi                    = "ncgi"
	SearchDataSupi                = "supi"
	SearchDataSupportedFeatures   = "supported-features"
	SearchDataAmfRegionID         = "amf-region-id"
	SearchDataAmfSetID            = "amf-set-id"
	SearchDataGuami               = "guami"
	SearchDataTac                 = "tac"
	SearchDataPlmnID              = "plmnId"
	SearchDataAmfID               = "amfId"
	SearchDataUEIPv4Addr          = "ue-ipv4-address"
	SearchDataIPDoamin            = "ip-domain"
	SearchDataUEIPv6Prefix        = "ue-ipv6-prefix"
	SearchDataRequesterNFInstFQDN = "requester-nf-instance-fqdn"
	SearchDataPGW                 = "pgw"
	SearchDataNsiList             = "nsi-list"
	SearchDataGpsi                = "gpsi"
	SearchDataExterGroupID        = "external-group-identity"
	SearchDataDataSet             = "data-set"
	SearchDataRoutingIndic        = "routing-indicator"
	SearchDataHnrfURI             = "hnrf-uri"
	SearchDataIfNoneMatch         = "If-None-Match"
	SupiFormat                    = "supiformat"
	SearchDataTargetNFFQDN        = "target-nf-fqdn"
	SearchDataForward             = "Forwarded"
	PLMNDiscForwardValue          = "by=_plmn.nrf"
	DiscoveryServiceName          = "nnrf-disc"
	SearchDataGroupIDList         = "group-id-list"
	SearchDataDnaiList            = "dnai-list"
	SearchDataPGWInd              = "pgw-ind"
	SearchDatacomplexQuery        = "complexQuery"
	SearchDataUpfIwkEpsInd        = "upf-iwk-eps-ind"
	SearchDataChfSupportedPlmn    = "chf-supported-plmn"
	SearchDataAccessType          = "access-type"
	SearchDataPreferredLocality   = "preferred-locality"
	BoolTrueString                = "true"
	BoolFalseString               = "false"
	SearchDataCacheControl        = "Cache-Control"
	SearchDataCacheControlPrivate = "private"
	SearchDataCacheControlNoCache = "no-cache"
	SearchDataCacheControlNoStore = "no-store"
	SearchDataCacheControlMaxAge0 = "max-age=0"
	/*
		SearchResult = `
		                  {
		                     "validityPeriod" : %d,
		                     "nfInstances" : %s
		                  }
		               `
	*/
	SearchResult               = "{\n  \"validityPeriod\" : %d,\n  \"nfInstances\" : %s}"
	SearchResultValidityPeriod = "validityPeriod"
	SearchResultNFInstances    = "nfInstances"
	SubsByNotificationURI      = "subs-by-nfStatusNotificationUri"
	SubsToNfInstanceID         = "subs-to-nfInstanceId"

	//SubscriptionData
	SubDataNotificationUri = "nfStatusNotificationUri"
	SubDataSubscrCond      = "subscrCond"
	SubDataValidityTime    = "validityTime"
	SubDatareqNotifEvents  = "reqNotifEvents"
	SubDataPlmnId          = "plmnId"
	SubDataNotifCondition  = "notifCondition"

	SubDataTargetNfType = "nfType"

	SubscriptionId        = "subscriptionId"
	SubscriptionIDPattern = "^([0-9]{5,6}-)?[^-]+$"

	StatusSubscribeRule1 = `subscrCond must be oneOf(NfInstanceIdCond, NfTypeCond, ServiceNameCond, AmfCond, GuamiListCond, NetworkSliceCond, NfGroupCond)`
	StatusSubscribeRule2 = `if nsiList is present, snssaiList shall be present too`
	StatusSubscribeRule3 = `if nfGroupId is present, nfType shall be present too`
	StatusSubscribeRule4 = `for NfGroupCond, nfType shall be oneOf(AUSF, UDM, UDR)`
	StatusSubscribeRule5 = `validityTime can't be before Now(%s)`
	StatusSubscribeRule6 = `attributes monitoredAttributes and unmonitoredAttributes shall not be included simultaneously`
	StatusSubscribeRule7 = `at least one of attributes "subscrCond" and "reqNfType" is present`
	StatusSubscribeRule8 = `the NF Service Consumer of nfType(%s) can not request a subscription to all NFs in the NRF`

	SubscriptionGetRuleRule1 = `at least one of (%s, %s, %s) shall be included in path`
	SubscriptionGetRuleRule2 = `at most one of (%s, %s, %s) shall be included in path`

	//HTTP
	HTTP_SERVER_READ_TIMEOUT              = 5 * time.Second
	HTTP_SERVER_WRITE_TIMEOUT             = 5 * time.Second
	HTTP2_MAX_STREAM_NUM                  = 100
	HTTP_MESSAGE_FORMAT                   = `{"title": "%s"}`
	HTTP_BODY_FORMAT                      = `{"nfInstanceId":"%s"}`
	HTTP_NOTIFICATION_BODY_FORMAT         = `{"subscriptionId":"%s"}`
	HTTP_SUBSCRIPTIONS_BODY_FORMAT        = `{"id":"%s"}`
	HTTP_SUBSCRIPTIONS_NOTIFY_BODY_FORMAT = `{"event": %s, "nfprofile": %s}`
	HTTP_NOTIFICATION_DATA_FORMAT         = `{"timestamp": "%s", "subscriptionID": "%s", "notificationBody": %s}`
	HTTP_HEADER_JSON_PATCH_JSON           = "application/json-patch+json"
	HttpHeaderProblemJson                 = "application/problem+json"
	// HttpHeader3GPPHalJson is the contect-type for 3gppHal+json
	HttpHeader3GPPHalJson = "application/3gppHal+json"
	// HttpHeaderJson is the contect-type for json
	HttpHeaderJson                        = "application/json"
	HTTPHeaderEtag                        = "Etag"
	RegistrationResponseFormat            = `{"heartBeatTimer":%d, "nfProfile":%s}`
	NFNotificationDataRegisterFormat      = `{"event":"%s", "nfInstanceUri":"%s", "nfProfile": %s}`
	NFNotificationDataProfileChangeFormat = `{"event":"%s", "nfInstanceUri":"%s", "profileChanges": %s}`
	NFNotificationDataDeRegFormat         = `{"event":"%s", "nfInstanceUri":"%s"}`
	HTTPServerWriteTimeoutProv            = 60 * time.Second

	//NF Resource URL
	NFInstanceIdInParameter = "nfInstanceID"
	NfInstanceIDResouceURL  = "/nnrf-nfm/v1/nf-instances/{nfInstanceID}"
	NfInstancesResouceURL   = "/nnrf-nfm/v1/nf-instances"

	NRFAddressIDName = "nrfAddressID"

	GroupProfileIDName = "groupProfileID"
	SUPIProfileIDName  = "supiProfileID"
	ImsiProfileIDName  = "imsiProfileID"
	GpsiProfileIDName  = "gpsiProfileID"
	//ErrorInfo
	MadatoryFieldNotExistFormat = "Madatory field %s doesn't exist in %s"
	FieldEmptyValue             = "Field %s can't be empty"
	FieldMultipleValue          = "Field %s doesn't support multiple value"
	ArrayFileldExistEmptyValue  = "Filed %s is array, when exist, minitem should be 1"
	UnSupportedQueryParameter   = "UNSUPPORTED_QUERY_PARAMETER"
	//NfProfileRule1 is a rule for NFProfile
	NfProfileRule1 = "at least one of the addressing parameters (fqdn, ipv4Addresses or ipv6Addresses) shall be included in the NF Profile"

	// NfProfileRule2 is a rule for IpEndPoint
	NfProfileRule2 = "at most one occurrence of either ipv4Address or ipv6Address shall be included in this data structure"

	// NfProfileRule3 is a rule for SupiRange
	NfProfileRule3 = "either the start and end attributes, or the pattern attribute, shall be present"

	// NfProfileRule4 is a rule for IdentityRange
	NfProfileRule4 = "either the start and end attributes, or the pattern attribute, shall be present"

	// NfProfileRule5 is a rule for InterfaceUpfInfoItem
	NfProfileRule5 = "at least one of the addressing parameters (ipv4EndpointAddresses, ipv6EndpointAddresses or endpointFqdn) shall be included in the InterfaceUpfInfoItem"

	// NfProfileRule6 is a rule for TacRange
	NfProfileRule6 = "either the start and end attributes, or the pattern attribute, shall be present"

	// NfProfileRule7 is a rule for N2InterfaceAmfInfo
	NfProfileRule7 = "at least one of the addressing parameters (ipv4EndpointAddress or ipv6EndpointAddress) shall be included"

	// NfProfileRule8 is a rule for ChfServiceInfo
	NfProfileRule8 = "at most one occurrence of either primaryChfServiceInstance or secondaryChfServiceInstance shall be included in this data structure"

	// NfProfileRule9 is a rule for PlmnRange
	NfProfileRule9 = "either the start and end attributes, or the pattern attribute, shall be present"

	// NfProfileRule10 is a rule for service name and nf type
	NfProfileRule10 = `service %s doesn't belog to %s`

	// ForbiddenUnlocalTitle is the title of problemDetail when rejecting an not local Registration
	ForbiddenUnlocalTitle = "not a local nf"

	//Log Message for HTTP Request
	REQUEST_LOG_FORMAT = `{"request":{"sequenceId":"%s", "URL":"%v", "method":"%s", "description":%s}}`

	//Log Message for HTTP Response
	RESPONSE_LOG_FORMAT = `{"response":{"sequenceId":"%s", "statusCode":%d, "description":%s}}`

	//Wild
	Wildcard = "*"
	//NoSubscrCond is a flag indicating a subscription withou subscrCond
	NoSubscrCond = "NOCOND"

	// NF profile expired flag
	NFprofileFlagValid   = 0
	NFprofileFlagExpired = 1
	NFprofileFlagAll     = 2

	// LastUpdateTimeNoneNfInfo is the last update time of newly registered NF profile without nfInfo
	LastUpdateTimeNoneNfInfo = uint64(1)

	//DbproxyGrpcCtxLongTimeout is the db-proxy get contex long timer
	DbproxyGrpcCtxLongTimeout = 10
)

var (
	NF_EVENT_MAP = map[string]bool{
		NF_EVENT_CREATED: true,
		NF_EVENT_UPDATED: true,
		NF_EVENT_DELETED: true,
	}
	NF_NFTYPES_MAP = map[string]bool{
		NfTypeUDM:   true,
		NfTypeAMF:   true,
		NfTypeSMF:   true,
		NfTypeAUSF:  true,
		NfTypeSMSF:  true,
		NfTypeNSSF:  true,
		NfTypeNEF:   true,
		NfTypePCF:   true,
		NfTypeNRF:   true,
		NfTypeUDR:   true,
		NfTypeUPF:   true,
		NfTypeLMF:   true,
		NfTypeGMLC:  true,
		NfType5GEIR: true,
		NfTypeSEPP:  true,
		NfTypeN3IWF: true,
		NfTypeAF:    true,
		NfTypeUDSF:  true,
		NfTypeBSF:   true,
		NfTypeCHF:   true,
		NfTypeNWDAF: true,
	}

	NFStatusMap = map[string]bool{
		NFStatusRegistered: true,
		NFStatusSuspended:  true,
	}

	NFProvisionMap = map[string]bool{
		NFAutoRegistered:   true,
		NFManualRegistered: true,
	}

	// NFInfoMap record the nfType which contained in nrfInfo
	NFInfoMap = map[string]string{
		NfTypeAMF:  AmfInfo,
		NfTypeAUSF: AusfInfo,
		NfTypePCF:  PcfInfo,
		NfTypeSMF:  SmfInfo,
		NfTypeUDM:  UdmInfo,
		NfTypeNRF:  NrfInfo,
	}

	// NRFProfile disc supported parameters
	NRFParaMap = map[string]bool{
		SearchDataTargetNfType:        true,
		SearchDataRequesterNfType:     true,
		SearchDataServiceName:         true,
		SearchDataRequesterNFInstFQDN: true,
		SearchDataTargetPlmnList:      true,
		SearchDataRequesterPlmnList:   true,
		SearchDataTargetInstID:        true,
		SearchDataTargetNFFQDN:        true,
		SearchDataSnssais:             true,
		SearchDataNsiList:             true,
		SearchDataSupportedFeatures:   true,
	}

	//TargetNFInfo is nfType reference to nfInfo
	TargetNFInfo = map[string]string{
		NfTypeAMF:  AmfInfo,
		NfTypeAUSF: AusfInfo,
		NfTypeBSF:  BsfInfo,
		NfTypeCHF:  ChfInfo,
		NfTypeNRF:  NrfInfo,
		NfTypePCF:  PcfInfo,
		NfTypeSMF:  SmfInfo,
		NfTypeUDM:  UdmInfo,
		NfTypeUDR:  UdrInfo,
		NfTypeUPF:  UpfInfo,
	}

	// ServiceNameNFTypeMap is the mapping relation between service name and NF type
	ServiceNameNFTypeMap = map[string]string{
		NNRFNFM:                  NfTypeNRF,
		NNRFDISC:                 NfTypeNRF,
		NUDMSDM:                  NfTypeUDM,
		NUDMUECM:                 NfTypeUDM,
		NUDMUEAU:                 NfTypeUDM,
		NUDMEE:                   NfTypeUDM,
		NUDMPP:                   NfTypeUDM,
		NAMFCOMM:                 NfTypeAMF,
		NAMFEVTS:                 NfTypeAMF,
		NAMFMT:                   NfTypeAMF,
		NAMFLOC:                  NfTypeAMF,
		NSMFPDUSESSION:           NfTypeSMF,
		NSMFEVENTEXPOSURE:        NfTypeSMF,
		NAUSFAUTH:                NfTypeAUSF,
		NAUSFSORPROTECTION:       NfTypeAUSF,
		NNEFPFDMANAGEMENT:        NfTypeNEF,
		NPCFAMPOLICYCONTROL:      NfTypePCF,
		NPCFSMPOLICYCONTROL:      NfTypePCF,
		NPCFPOLICYAUTHORIZATION:  NfTypePCF,
		NPCFBDTPOLICYCONTROL:     NfTypePCF,
		NPCFEVENTEXPOSURE:        NfTypePCF,
		NPCFUEPOLICYCONTROL:      NfTypePCF,
		NSMSFSMS:                 NfTypeSMSF,
		NNSSFNSSELECTION:         NfTypeNSSF,
		NNSSFNSSAIAVAILABILITY:   NfTypeNSSF,
		NUDRDR:                   NfTypeUDR,
		NLMFLOC:                  NfTypeLMF,
		N5GEIREIC:                NfType5GEIR,
		NBSFMANAGEMENT:           NfTypeBSF,
		NCHFSPENDINGLIMITCONTROL: NfTypeCHF,
		NCHFCONVERGEDCHARGING:    NfTypeCHF,
		NNWDAFEVENTSSUBSCRIPTION: NfTypeNWDAF,
		NNWDAFANALYTICSINFO:      NfTypeNWDAF,
	}

	// CountLabelListRequest for request labels
	CountLabelListRequest = []string{"resource", "operation", "remote_endpoint"}

	// CountLabelListSuccResponse for successful response labels
	CountLabelListSuccResponse = []string{"resource", "operation", "remote_endpoint", "status_code"}

	// CountLabelListUnSuccResponse for unsuccessful response labels
	CountLabelListUnSuccResponse = []string{"resource", "operation", "remote_endpoint", "status_code", "detailed_info"}
)

const (
	// NfNotificationEventsTotal is used to record total number of notification event triggered
	NfNotificationEventsTotal = "nf_notification_event_total"

	// NfNotificationEventsDiscardedTotal is used to record total number of notification event discarded
	NfNotificationEventsDiscardedTotal = "nf_notification_event_discarded_total"

	// NfNotificationSentTotal is used to record total number of notification sent to nf instance
	NfNotificationSentTotal = "nf_notification_sent_total"

	// NfRegister Service Operation
	NfRegister = "NFRegister"

	// NfUpdate Service Operation
	NfUpdate = "NFUpdate"

	// NfDeregister Service Operation
	NfDeregister = "NFDeregister"

	// NfStatusSubscribe Service Operation
	NfStatusSubscribe = "NFStatusSubscribe"

	// NfStatusUnSubscribe Service Operation
	NfStatusUnSubscribe = "NFStatusUnSubscribe"

	// NfManagement Service Operation
	NfManagement = "NFManagement"

	// NfDiscovery Service Operation
	NfDiscovery = "NFDiscovery"

	// NfProvision Service Operation
	NfProvision = "NFProvision"

	// NfRequestDuration is used to record nf request duration
	NfRequestDuration = "nf_request_duration_seconds"

	// NfProfiles is used to record nf profiles count
	NfProfiles = "eric_nrf_nf_profiles"

	// NfManagementRequestsTotal
	NfManagementRequestsTotal = "eric_nrf_nnrf_nfm_requests_recv"

	// NfManagementSuccResponseTotal is used to record the
	// total number of successful response sent.
	NfManagementSuccResponseTotal = "eric_nrf_nnrf_nfm_successful_responses_sent"

	// NfManagementUnSuccResponseTotal is used to record the
	// total number of failed response sent.
	NfManagementUnSuccResponseTotal = "eric_nrf_nnrf_nfm_unsuccessful_responses_sent"

	// NfDiscoveryRequestsTotal is used to record the
	// total number of discovery requests received in the NRF.
	NfDiscoveryRequestsTotal = "eric_nrf_nnrf_disc_requests_recv"

	// NfDiscoverySuccessTotal is used to record the
	// total number of successful discovery.
	NfDiscoverySuccessTotal = "eric_nrf_nnrf_disc_successful_responses_sent"

	// NfDiscoveryFailureTotal is used to record the
	// total number of failed discovery.
	NfDiscoveryFailureTotal = "eric_nrf_nnrf_disc_unsuccessful_responses_sent"

	// NfProvisionRequestsTotal is used to record the
	// total number of requests received in the NRF Provision.
	NfProvisionRequestsTotal = "eric_nrf_nnrf_prov_requests_recv"

	// NfProvisionSuccResponsesTotal is used to record the
	// total number of successful responses sent.
	NfProvisionSuccResponsesTotal = "eric_nrf_nnrf_prov_successful_responses_sent"

	// NfProvisionUnSuccResponsesTotal is used to record the
	// total number of failed responses sent.
	NfProvisionUnSuccResponsesTotal = "eric_nrf_nnrf_prov_unsuccessful_responses_sent"
)

const (

	//NRFFqdnFormat is used to construct NRF FQDN
	NRFFqdnFormat = "nrf.5gc.mnc%s.mcc%s.3gppnetwork.org"

	//After receiving SIGTERM, how much time need to wait to ensure ongoing traffic is finished.
	TerminateWaitingTime = 6
)

const (
	LT    = -2
	LE    = -1
	EQ    = 0
	GE    = 1
	GT    = 2
	REGEX = 3
)

const (
	ValueString = 1
	ValueNum    = 2
)

const (
	TypeAndExpression    = 1
	TypeORExpression     = 2
	TypeSearchExpression = 3
)

const (
	//TargetNFProfile is used to construct NFProfile string for GRPC interface
	TargetNFProfile = `{"expiredTime": %d, "lastUpdateTime": %d, "profileUpdateTime": %d, "provisioned": %d, "md5sum": {%s}, "helper": {%s}, "body": %s, "provSupiVersion": %d, "provGpsiVersion": %d}`
	//TargetNFProfileWithOverride is used to construct NFProfile string for GRPC interface
	TargetNFProfileWithOverride = `{"expiredTime": %d, "lastUpdateTime": %d, "profileUpdateTime": %d, "provisioned": %d, "md5sum": {%s}, "helper": {%s}, "body": %s, "overrideInfo": %s, "provSupiVersion": %d, "provGpsiVersion": %d}`
	//EmptyAllowedDomain is a default value when allowedDomain is empty
	EmptyAllowedDomain = "RESERVED_EMPTY_DOMAIN"
	//EmptyAllowedNfType is a default value when allowedNfType is empty
	EmptyAllowedNfType = "RESERVED_EMPTY_TYPE"
	//EmptyAllowedPlmnMcc is a default value when mcc is empty
	EmptyAllowedPlmnMcc = "XXX"
	//EmptyAllowedPlmnMnc is a default value when mnc is empty
	EmptyAllowedPlmnMnc = "YYY"
	//EmptyExternalIDPattern is a default value when externalId is empty
	EmptyExternalIDPattern = "RESERVED_EMPTY_EXTERNAL_ID_RANGE_PATTERN"
	//EmptySupiRangePattern is a default value when supirange is empty
	EmptySupiRangePattern = "RESERVED_EMPTY_SUPI_RANGE_PATTERN"
	//EmptyGpsiRangePattern is a default value when gpsirange is empty
	EmptyGpsiRangePattern = "RESERVED_EMPTY_GPSI_RANGE_PATTERN"
	//EmptyPlmnRangePattern is a default value when plmnrange is empty
	EmptyPlmnRangePattern = "RESERVED_EMPTY_PLMN_RANGE_PATTERN"
	//EmptyDnai is a default value when dnai is empty
	EmptyDnai = "RESERVED_EMPTY_DNAI"
	//EmptyMcc is a default value when mcc is empty
	EmptyMcc = "RESERVED_EMPTY_MCC"
	//EmptyMnc is a default value when mnc is empty
	EmptyMnc = "RESERVED_EMPTY_MNC"
	//EmptyTac is a default value when tac is empty
	EmptyTac = "RESERVED_EMPTY_TAC"
	//EmptyTacRangePattern is a default value when tacrange is empty
	EmptyTacRangePattern = "RESERVED_EMPTY_TAC_RANGE_PATTERN"
	//EmptySd is a default value when sd is empty
	EmptySd = "RESERVED_EMPTY_SD"
	//EmptySst is a default value when sst is empty
	EmptySst = uint64(256)
	//EmptyGroupID is a default value when groupId is empty
	EmptyGroupID = "RESERVED_EMPTY_GROUPID"
	//AttrNoAbsence is a default value when attr is empty
	AttrNoAbsence = "ATTRIBUTE_NO_ABSENCE"
	//MatchAll is a flag
	MatchAll = "MATCH_ALL"
	//NoMatchAll is a flag
	NoMatchAll = "NO_MATCH_ALL"
)

const (

	//OverrideAttrPath is for overrideAttrList path value in schema
	OverrideAttrPath = "/provisionInfo/overrideAttrList"
	//Cmode_NFRegiestered_Provisioned is for int value of both NF_REGISTERED and PROVISIONED
	Cmode_NFRegiestered_Provisioned = 0
	//Cmode_NFRegistered is for int value of createMode NF_REGISTERED
	Cmode_NFRegistered = 1
	//Cmode_Provisioned is for int value of createMode PROVISIONED
	Cmode_Provisioned = 2
	//CMODE_NF_REGISTERED is for string value of createMode NF_REGISTERED
	CMODE_NF_REGISTERED = "NF_REGISTERED"
	//CMODE_NF_PROVISIONED is for string value of createMode PROVISIONED
	CMODE_PROVISIONED = "PROVISIONED"
	//Start is for SupiRange/GpsiRange start property
	Start = "start"
	//End is for SupiRange/GpsiRange end property
	End = "end"
)

const (
	//PathSst is to generate sst Attributes parameter
	PathSst = "/sst"
	//PathSd is to generate sd Attributes parameter
	PathSd = "/sd"
	//PathStart is to generate start Attributes parameter
	PathStart = "/start"
	//PathEnd is to generate end Attributes parameter
	PathEnd = "/end"
	//PathStartLength is to generate start's length Attributes parameter
	PathStartLength = "/start/length"
	//PathEndLength is to generate end's length Attributes parameter
	PathEndLength = "/end/length"
	//PathPattern is to generate pattern Attributes parameter
	PathPattern = "/pattern"
	//PathPlmnMcc is to generate plmn/mcc Attributes parameter
	PathPlmnMcc = "/plmnid/mcc"
	//PathPlmnMnc is to generate plmn/mnc Attributes parameter
	PathPlmnMnc = "/plmnid/mnc"
	//PathAmfID is to generate amfId Attributes parameter
	PathAmfID = "/amfId"
	//PathList is to generate list Attributes parameter(tailList)
	PathList = "/list"
	//PathRangeList is to generate rangeList Attributes parameter(taiRangeList)
	PathRangeList = "/rangelist"
	//PathTac is to generate tac Attributes parameter
	PathTac = "/tac"
	//PathBackfailure is to generate amf backup failure Attributes parameter
	PathBackfailure = "/backfailure"
	//PathBackremoval is to generate amf backup removal Attributes parameter
	PathBackremoval = "/backremoval"
	//PathAbsencePattern is to generate supi/gpsi Attributes parameter
	PathAbsencePattern = "/absence/pattern"
)
