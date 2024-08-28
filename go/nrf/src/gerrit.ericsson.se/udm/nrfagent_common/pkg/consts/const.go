package consts

import (
	"time"
)

//Nrfclient
const (
	AppWorkmodeREG  = "reg"
	AppWorkmodeNTF  = "ntf"
	AppWorkmodeDISC = "disc"
	AppWorkmodeTEST = "test"

	CmdStartREG  = "start_reg"
	CmdStartNTF  = "start_ntf"
	CmdStartDISC = "start_disc"
	CmdTEST      = "test"
	CmdVERSION   = "version"
	CmdHELP      = "help"
	CmdSTATUS    = "status"

	HTTPServerReadTimeout  = 5 * time.Second
	HTTPServerWriteTimeout = 5 * time.Second
	HTTPMessageFormat      = `{"message": "%s"}`

	nfEventCreated = "configCreated"
	nfEventUpdated = "configUpdated"
	nfEventDeleted = "configDeleted"

	ntfMessage     = "N1_MESSAGES"
	ntfInformation = "N2_INFORMATION"
	ntfLocationNTF = "LOCATION_NOTIFICATION"
	ntfRMNTF       = "DATA_REMOVAL_NOTIFICATION"
	ntfCHGNTF      = "DATA_CHANGE_NOTIFICATION"

	dtSetSubs        = "SUBSCRIPTION"
	dtSetPolicy      = "POLICY"
	dtSetExposure    = "EXPOSURE"
	dtSetApplication = "APPLICATION"

	n15GMM = "5GMM"
	n1SM   = "SM"
	n1LPP  = "LPP"
	n1SMS  = "SMS"

	n2InfoSM    = "SM"
	n2InfoNRPPA = "NRPPA"

	nfTypeNRF   = "NRF"
	nfTypeUDM   = "UDM"
	nfTypeAMF   = "AMF"
	nfTypeSMF   = "SMF"
	nfTypeAUSF  = "AUSF"
	nfTypeSMSF  = "SMSF"
	nfTypeNSSF  = "NSSF"
	nfTypeNEF   = "NEF"
	nfTypePCF   = "PCF"
	nfTypeGMLC  = "GMLC"
	nfTypeUDR   = "UDR"
	nfTypeLMF   = "LMF"
	nfType5GEIR = "5G_EIR"
	nfTypeSEPP  = "SEPP"
	nfTypeUPF   = "UPF"
	nfTypeN3IWF = "N3IWF"
	nfTypeAF    = "AF"
	nfTypeUDSF  = "UDSF"
	nfTypeBSF   = "BSF"
	nfTypeCHF   = "CHF"
	nfTypeNWDAF = "NWDAF"

	NfProfileConfFile   = "/nrfclient-nfprofile/..data/nf-profile.json"
	NfDiscConfFile      = "/nrfclient-disc-conf/..data/target-nf-cnf.json"
	NrfMgmConnConfFile  = "/nrf-mgm-connection-conf/..data/nrf-mgm-cnf.json"
	NrfDiscConnConfFile = "/nrf-disc-connection-conf/..data/nrf-disc-cnf.json"

	DiscoveryAgentReadinessProbe = "/nrf-agent-disc/v1/ready-check"

	ServerIsInitializing = "Initialzing"
	ServerIsRunning      = "Running"
	ServerIsClosing      = "Closing"

	StatusReg   = "REGISTERED"
	StatusSupd  = "SUSPENDED"
	StatusDeReg = "DEREGISTERED"

	UpdateLoad = "UPDATELOAD"

	NFRegister         = "NF_REGISTERED"
	NFDeRegister       = "NF_DEREGISTERED"
	NFProfileChg       = "NF_PROFILE_CHANGED"
	NFEventWithoutBody = "AGENT_EVENT_WITHOUT_BODY"
	NFEventDiscResult  = "AGENT_EVENT_DISC_RESULT"

	UPInterfaceTypeN3 = "N3"
	UPInterfaceTypeN6 = "N6"
	UPInterfaceTypeN9 = "N9"

	NotifEventRegister   = "NF_REGISTERED"
	NotifEventDeregister = "NF_DEREGISTERED"
	NotifEventProfileChg = "NF_PROFILE_CHANGED"

	MsgbusTopicNamePrefix = "nrf-agent-"

	EventTypeRegister    = "REGISTER"
	EventTypeDeregister  = "DEREGISTER"
	EventTypeFQDNChanged = "FQDN_CHANGED"

	EventTypeSyncSubscrInfo = "SYNC_SUBSCRINFO"

	NFLoad         = "load"
	HeartBeatTimer = "heartBeatTimer"
)

//consts used in configmap storage function
const (
	//ConfigMapStorage define name of configmap for agent storage
	ConfigMapStorage = "eric-nrfagent-storage"
	//ConfigMapKeyNfInfo define tag for NfInfo
	ConfigMapKeyNfInfo = "nfInfoList"
	//ConfigMapKeySubsInfo define tag for SubscriptionInfo
	ConfigMapKeySubsInfo = "subscriptionInfoList"
)

const (
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

	BoolTrueString  = "true"
	BoolFalseString = "false"

	SearchDataCacheControl        = "Cache-Control"
	SearchDataCacheControlPrivate = "private"
	SearchDataCacheControlNoCache = "no-cache"
	SearchDataCacheControlNoStore = "no-store"
	SearchDataCacheControlMaxAge0 = "max-age=0"

	//ErrorInfo
	MadatoryFieldNotExistFormat = "Madatory field %s doesn't exist in %s"
	FieldEmptyValue             = "Field %s can't be empty"
	FieldMultipleValue          = "Field %s doesn't support multiple value"
	ArrayFileldExistEmptyValue  = "Filed %s is array, when exist, minitem should be 1"
	UnSupportedQueryParameter   = "UNSUPPORTED_QUERY_PARAMETER"

	//enum AccessType
	Access3GPP    = "3GPP_ACCESS"
	NonAccess3GPP = "NON_3GPP_ACCESS"

	RoamSuffix = "-roam"
)

const (
	SupiRanges = "supiRanges"
	GpsiRanges = "gpsiRanges"
)

const (
	//Log Message for HTTP Request
	REQUEST_LOG_FORMAT = `{"request":{"sequenceId":"%s", "URL":"%v", "method":"%s", "description":%s}}`

	//Log Message for HTTP Response
	RESPONSE_LOG_FORMAT = `{"response":{"sequenceId":"%s", "statusCode":%d, "description":%s}}`
)

//enum value map define
var (
	NfEventMap = map[string]bool{
		nfEventCreated: true,
		nfEventUpdated: true,
		nfEventDeleted: true,
	}
	NfNfTypesMap = map[string]bool{
		nfTypeUDM:   true,
		nfTypeAMF:   true,
		nfTypeSMF:   true,
		nfTypeAUSF:  true,
		nfTypeSMSF:  true,
		nfTypeNSSF:  true,
		nfTypeNEF:   true,
		nfTypePCF:   true,
		nfTypeNRF:   true,
		nfTypeUDR:   true,
		nfTypeLMF:   true,
		nfType5GEIR: true,
		nfTypeSEPP:  true,
		nfTypeUPF:   true,
		nfTypeN3IWF: true,
		nfTypeAF:    true,
		nfTypeUDSF:  true,
		nfTypeBSF:   true,
		nfTypeCHF:   true,
		nfTypeNWDAF: true,
	}

	NtfTypeMap = map[string]bool{
		ntfMessage:     true,
		ntfInformation: true,
		ntfLocationNTF: true,
		ntfRMNTF:       true,
		ntfCHGNTF:      true,
	}

	N1MsgClass = map[string]bool{
		n15GMM: true,
		n1LPP:  true,
		n1SM:   true,
		n1SMS:  true,
	}

	N2InfoClass = map[string]bool{
		n2InfoNRPPA: true,
		n2InfoSM:    true,
	}

	DateSetID = map[string]bool{
		dtSetApplication: true,
		dtSetExposure:    true,
		dtSetPolicy:      true,
		dtSetSubs:        true,
	}

	UPInterfaceTypeMap = map[string]bool{
		UPInterfaceTypeN3: true,
		UPInterfaceTypeN6: true,
		UPInterfaceTypeN9: true,
	}
	CMIgnoreParamterMap = map[string]bool{
		NFLoad:         true,
		HeartBeatTimer: true,
	}
)

//consts used in PM function
const (
	//	// NfRegistrationRequestsTotal is used to record the
	//	// total number of registration requests received in the NRF.
	//	NfRegistrationRequestsTotal = "nf_registration_requests_total"

	//	// NfRegistrationSuccessTotal is used to record the
	//	// total number of successful registrations.
	//	NfRegistrationSuccessTotal = "nf_registration_success_total"

	//	// NfRegistrationFailureTotal is used to record the
	//	// total number of failed registrations.
	//	NfRegistrationFailureTotal = "nf_registration_failure_total"

	// NfRequestDuration is used to record nf request duration
	NfRequestDuration = "nf_request_duration_seconds"

	// NfHeartBeat Service Operation
	NfHeartBeat = "NFHeartBeat"

	// NfRegister Service Operation
	NfRegister = "NFRegister"

	// NfUpdate Service Operation
	NfUpdate = "NFUpdate"

	// NfDeregister Service Operation
	NfDeregister = "NFDeregister"

	// NfDiscovery Service Operation
	NfDiscovery = "NFDiscovery"

	// NfDiscoveryRequestsTotal is used to record the
	// total number of discovery requests received from NF.
	NfDiscoveryRequestsTotal = "nf_discovery_requests_total"

	// NfDiscoveryResponsesTotal is used to record the
	// total number of discovery responses send to NF.
	NfDiscoveryResponsesTotal = "nf_discovery_responses_total"

	// NrfDiscoveryRequestsTotal is used to record the
	// total number of discovery requests send to NRF.
	NrfDiscoveryRequestsTotal = "nrf_discovery_requests_total"

	// NrfDiscoveryResponsesTotal is used to record the
	// total number of discovery responses received from NRF.
	NrfDiscoveryResponsesTotal = "nrf_discovery_responses_total"

	// NrfDiscoveryResponses2xx is used to record the
	// total numer of discovery 2xx responses received from NRF.
	NrfDiscoveryResponses2xx = "nrf_discovery_responses_2xx"

	// NrfDiscoveryResponses3xx is used to record the
	// total number of discovery 3xx responses received from NRF.
	NrfDiscoveryResponses3xx = "nrf_discovery_responses_3xx"

	// NrfDiscoveryResponses4xx is used to record the
	// total number of discovery 4xx responses received from NRF.
	NrfDiscoveryResponses4xx = "nrf_discovery_responses_4xx"

	// NrfDiscoveryResponses5xx is used to record the
	// total number of discovery 5xx responses received from NRF.
	NrfDiscoveryResponses5xx = "nrf_discovery_responses_5xx"
)

//consts used in FM function
const (
	//AlarmRegServiceName define FM service name of reg agent
	AlarmRegServiceName = "eric-nrfagent-RegisterAgent"
	//AlarmRegServiceName define FM service name of disc agent
	AlarmDiscServiceName = "eric-nrfagent-DiscoverAgent"
)

//consts used in messagebus topic name postfix
const (
	RegDiscInner  = "regdiscinner"
	NtfDiscInner  = "ntfdiscinner"
	DiscDiscInner = "discdiscinner"
)
