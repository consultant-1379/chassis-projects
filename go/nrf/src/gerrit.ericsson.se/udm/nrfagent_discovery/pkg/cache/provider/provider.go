package provider

var Content = []string{
	`{
		"nfInstanceID": "ausf-5g-01",
		"nfType": "ausf",
		"plmn": "24080",
		"sNssai": {
		  "sst": "0",
		  "sd": "0"
		},
		"fqdn": "seliius03696.seli.gic.ericsson.se",
		"ipAddress": [
		  "172.16.208.1"
		],
		"capacity": 100,
		"nfServiceList": [
		  {
			"serviceInstanceID": "nausf-auth-01",
			"serviceName": "nausf-auth",
			"version": "v1",
			"Schema": "https://",
			"fqdn": "seliius03696.seli.gic.ericsson.se",
			"ipAddress": [
			  "172.16.208.1"
			],
			"port": 30088,
			"callbackUri": [
			  "https://seliius03696.seli.gic.ericsson.se/notification"
			],
			"allowedPlmns": [
			  "46000"
			],
			"allowedNfTypes": [
			  "amf"
			],
			"allowedNssais": [
			  {
				"sst": "0",
				"sd": "0"
			  }
			]
		  }
		]
	   }`,
	`{
		"nfInstanceID": "ausf-5g-02",
		"nfType": "ausf",
		"plmn": "24081",
		"sNssai": {
		  "sst": "0",
		  "sd": "0"
		},
		"fqdn": "seliius03696.seli.gic.ericsson.se",
		"ipAddress": [
		  "172.16.208.1"
		],
		"capacity": 100,
		"nfServiceList": [
		  {
			"serviceInstanceID": "nausf-auth-01",
			"serviceName": "nausf-auth",
			"version": "v1",
			"Schema": "https://",
			"fqdn": "seliius03696.seli.gic.ericsson.se",
			"ipAddress": [
			  "172.16.208.1"
			],
			"port": 30088,
			"callbackUri": [
			  "https://seliius03696.seli.gic.ericsson.se/notification"
			],
			"allowedPlmns": [
			  "46000"
			],
			"allowedNfTypes": [
			  "amf"
			],
			"allowedNssais": [
			  {
				"sst": "0",
				"sd": "0"
			  }
			]
		  }
		]
	   }`,
	`{
		"nfInstanceID": "ausf-5g-03",
		"nfType": "ausf",
		"plmn": "24082",
		"sNssai": {
		  "sst": "0",
		  "sd": "0"
		},
		"fqdn": "seliius03696.seli.gic.ericsson.se",
		"ipAddress": [
		  "172.16.208.1"
		],
		"capacity": 100,
		"nfServiceList": [
		  {
			"serviceInstanceID": "nausf-auth-01",
			"serviceName": "nausf-auth",
			"version": "v1",
			"Schema": "https://",
			"fqdn": "seliius03696.seli.gic.ericsson.se",
			"ipAddress": [
			  "172.16.208.1"
			],
			"port": 30088,
			"callbackUri": [
			  "https://seliius03696.seli.gic.ericsson.se/notification"
			],
			"allowedPlmns": [
			  "46000"
			],
			"allowedNfTypes": [
			  "amf"
			],
			"allowedNssais": [
			  {
				"sst": "0",
				"sd": "0"
			  }
			]
		  }
		]
	   }`,
}

var (
	NfInstanceId       = "nfInstanceId"
	NfType             = "nfType"
	PlmnListOfTargetNf = "plmnList"
	Mcc                = "mcc"
	Mnc                = "mnc"
	Snssais            = "sNssais"
	Sst                = "sst"
	Sd                 = "sd"
	Fqdn               = "fqdn"
	IpAddress          = "ipAddress"
	Capacity           = "capacity"
	NfServices         = "nfServices"
	ServiceName        = "serviceName"
	SmfInfo            = "smfInfo"
	DnnList            = "dnnList"
	NfInstances        = "nfInstances"
	BsfInfo            = "bsfInfo"
	UpfInfo            = "upfInfo"
	AccessType         = "accessType"
	ChfInfo            = "chfInfo"
	PlmnRangeList      = "plmnRangeList"
	Locality           = "locality"
)

var (
	//Db table
	NF_INSTANCE_SET_KEYNAME                   = "nfinstance"
	NF_INSTANCE_ID_KEYVALUE_KEYNAME           = "nfinstance-id:%s"
	NF_INSTANCE_NFTYPE_SET_KEYNAME            = "nfinstance-nftype:%s"
	NF_INSTANCE_REPO_PLMN_SD_SET_KEYNAME      = "nfinstance-repo:%s:%s"
	NF_NOTIFICATON_SET_KEYNAME                = "nfnotification"
	NF_NOTIFICATON_NFTYPE_SET_KEYNAME         = "nfnotification-nftype:%s"
	NF_NOTIFICATION_ID_KEYVALUE_KEYNAME       = "nfnotification-id:%s"
	NF_NOTIFICATION_NFINSTANCE_ID_SET_KEYNAME = "nfnotification-nfinstanceid:%s"
)

const (
	//MatchAllGroupID is for Nfprofile match all groupID
	MatchAllGroupID = "MatchAllGroupID"
)
