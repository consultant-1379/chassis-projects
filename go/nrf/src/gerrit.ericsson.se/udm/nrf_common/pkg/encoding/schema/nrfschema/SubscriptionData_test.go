package nrfschema

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"github.com/buger/jsonparser"
)

func TestGenerateValidatyDateTime(t *testing.T) {

	// case 1: ValidityPeriodOfSubscription in CM = 86400, and validityTime is not included in subsctionData
	{
		subscriptionJsonData := []byte(`{
			                         "subscrCond": {
								        "nfType": "UDM"
									},
								    "nfStatusNotificationUri" : "http://seliius04099.seli.gic.ericsson.se:20001"
								}`)
		cm.ValidityPeriodOfSubscription = 86400
		subscriptionData := TSubscriptionData{}

		json.Unmarshal(subscriptionJsonData, &subscriptionData)
		cmValidateTimeInSeconds := int64(cm.ValidityPeriodOfSubscription) + time.Now().Unix()
		expectedValidityDateTime := time.Unix(cmValidateTimeInSeconds, 0).Format(time.RFC3339)
		if subscriptionData.GenerateValidatyDateTime() != expectedValidityDateTime {
			t.Fatalf("validatiDateTime should be %s, but %s", expectedValidityDateTime, subscriptionData.GenerateValidatyDateTime())
		}
	}

	// case 2: ValidityPeriodOfSubscription in CM = 86400, and validityTime is included in subsctionData, but less than CM
	{
		subscriptionJsonData := []byte(`{
			                         "subscrCond": {
								        "nfType": "UDM"
									},
								    "validityTime" : "2015-12-31T23:59:59.999Z",
								    "nfStatusNotificationUri" : "http://seliius04099.seli.gic.ericsson.se:20001"
								}`)
		cm.ValidityPeriodOfSubscription = 86400
		subscriptionData := TSubscriptionData{}

		json.Unmarshal(subscriptionJsonData, &subscriptionData)
		subscriptionData.ValidityTime = time.Unix(int64(3600)+time.Now().Unix(), 0).Format(time.RFC3339)
		expectedValidityDateTime := subscriptionData.ValidityTime
		if subscriptionData.GenerateValidatyDateTime() != expectedValidityDateTime {
			t.Fatalf("validatiDateTime should be %s, but %s", expectedValidityDateTime, subscriptionData.GenerateValidatyDateTime())
		}
	}

	// case 3: ValidityPeriodOfSubscription in CM = 86400, and validityTime is included in subsctionData, but larger than CM
	{
		subscriptionJsonData := []byte(`{
			                         "subscrCond": {
								        "nfType": "UDM"
									},
								    "validityTime" : "3015-12-31T23:59:59.999Z",
								    "nfStatusNotificationUri" : "http://seliius04099.seli.gic.ericsson.se:20001"
								}`)
		cm.ValidityPeriodOfSubscription = 86400
		subscriptionData := TSubscriptionData{}

		json.Unmarshal(subscriptionJsonData, &subscriptionData)
		cmValidateTimeInSeconds := int64(cm.ValidityPeriodOfSubscription) + time.Now().Unix()
		expectedValidityDateTime := time.Unix(cmValidateTimeInSeconds, 0).Format(time.RFC3339)
		if subscriptionData.GenerateValidatyDateTime() != expectedValidityDateTime {
			t.Fatalf("validatiDateTime should be %s, but %s", expectedValidityDateTime, subscriptionData.GenerateValidatyDateTime())
		}
	}
}

func TestGenerateExpiredTimeInMilSec(t *testing.T) {
	// case 1: ValidityPeriodOfSubscription in CM = 86400, and validityTime is included in subsctionData, larger than Now and  less than CM
	{
		subscriptionJsonData := []byte(`{
			                         "subscrCond": {
								        "nfType": "UDM"
									},
								    "validityTime" : "3015-12-31T23:59:59.999Z",
								    "nfStatusNotificationUri" : "http://seliius04099.seli.gic.ericsson.se:20001"
								}`)
		cm.ValidityPeriodOfSubscription = 86400

		offset := 3600
		validateTimeInSeconds := int64(offset) + time.Now().Unix()
		expectedValidityDateTime := time.Unix(validateTimeInSeconds, 0).Format(time.RFC3339)

		value, _ := jsonparser.Set(subscriptionJsonData, []byte(fmt.Sprintf(`"%s"`, expectedValidityDateTime)), "validityTime")

		subscriptionData := TSubscriptionData{}
		json.Unmarshal(value, &subscriptionData)

		if subscriptionData.GenerateExpiredTimeInMilSec() != validateTimeInSeconds*1000 {
			t.Fatalf("the expired time  should be %d, but %d", validateTimeInSeconds*1000, subscriptionData.GenerateExpiredTimeInMilSec())
		}
	}
}

func TestValidateNotificationURI(t *testing.T) {
	// invalid nfStatusNotificationUri
	body := []byte(`{
		"nfStatusNotificationUri": "ftp://10.10.10.10:80/callback",
		"validityTime": "3015-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		},
		"subscrCond": {
			"nfType": "AMF"
		}
	}`)

	subscriptionData := &TSubscriptionData{}
	err := json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails := subscriptionData.ValidateNotificationURI()
	if problemDetails == nil {
		t.Fatalf("SubscriptionData.ValidateNotificationURI should not return nil, but did!")
	}

	// valid nfStatusNotificationUri
	body = []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"validityTime": "3015-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		},
		"subscrCond": {
			"nfType": "AMF"
		}
	}`)

	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails = subscriptionData.ValidateNotificationURI()
	if problemDetails != nil {
		t.Fatalf("SubscriptionData.ValidateNotificationURI should return nil, but not!")
	}
}

func TestValidateSubscrCond(t *testing.T) {
	// subscriptionData with invalid subscrCond is invalid
	body := []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"validityTime": "3015-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		},
		"subscrCond": {
			"nfInstanceId": "amf01",
			"nfType": "AMF"
		}
	}`)

	subscriptionData := &TSubscriptionData{}
	err := json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails := subscriptionData.ValidateSubscrCond()
	if problemDetails == nil {
		t.Fatalf("SubscriptionData.ValidateSubscrCond should not return nil, but did!")
	}

	// subscriptionData with valid subscrCond is valid
	body = []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"validityTime": "3015-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		},
		"subscrCond": {
			"nfType": "AMF"
		}
	}`)

	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails = subscriptionData.ValidateSubscrCond()
	if problemDetails != nil {
		t.Fatalf("SubscriptionData.ValidateSubscrCond should return nil, but not!")
	}

	// subscriptionData without subscrCond and reqNfType is empty
	body = []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"validityTime": "3015-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		}
	}`)
	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails = subscriptionData.ValidateSubscrCond()
	if problemDetails == nil {
		t.Fatalf("SubscriptionData.ValidateSubscrCond should not return nil, but did!")
	}

	// subscriptionData without subscrCond and reqNfType is not empty, and cm.NrfPolicyProfile.SubscriptionPolicy is not configured
	body = []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"validityTime": "3015-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"reqNfType": "NEF",
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		}
	}`)
	subscriptionPolicy := &cm.TSubscriptionPolicy{}
	cm.NrfPolicy = cm.TNrfPolicy{ManagementService: &cm.TNrfManagementServicePolicy{Subscription: subscriptionPolicy}}
	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails = subscriptionData.ValidateSubscrCond()
	if problemDetails == nil {
		t.Fatalf("SubscriptionData.ValidateSubscrCond should not return nil, but did!")
	}

	// subscriptionData without subscrCond and reqNfType is not empty
	body = []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"validityTime": "3015-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"reqNfType": "NEF",
		"reqNfFqdn": "ericsson.se",
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		}
	}`)
	// subscriptionData without subscrCond and reqNfType is not empty, and reqNfType don't match AllowedNfType configured in CM
	allowedSubscriptionAllNF := cm.TAllowedSubscriptionAllNFs{AllowedNfType: "AMF"}
	cm.NrfPolicy.ManagementService.Subscription.AllowedSubscriptionAllNFs = append(cm.NrfPolicy.ManagementService.Subscription.AllowedSubscriptionAllNFs, allowedSubscriptionAllNF)
	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails = subscriptionData.ValidateSubscrCond()
	if problemDetails == nil {
		t.Fatalf("SubscriptionData.ValidateSubscrCond should not return nil, but did!")
	}

	// subscriptionData without subscrCond and reqNfType is not empty, and reqNfType match AllowedNfType configured in CM
	allowedSubscriptionAllNF = cm.TAllowedSubscriptionAllNFs{AllowedNfType: "NEF"}
	cm.NrfPolicy.ManagementService.Subscription.AllowedSubscriptionAllNFs = append(cm.NrfPolicy.ManagementService.Subscription.AllowedSubscriptionAllNFs, allowedSubscriptionAllNF)
	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails = subscriptionData.ValidateSubscrCond()
	if problemDetails != nil {
		t.Fatalf("SubscriptionData.ValidateSubscrCond should return nil, but don't!")
	}

	// subscriptionData without subscrCond,reqNfFqdn and reqNfType is not empty, and reqNfFqdn don't match AllowedNfFqdn configured in CM
	cm.NrfPolicy.ManagementService.Subscription.AllowedSubscriptionAllNFs[1].AllowedNfDomains = "ericsson01.se"
	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails = subscriptionData.ValidateSubscrCond()
	if problemDetails == nil {
		t.Fatalf("SubscriptionData.ValidateSubscrCond should not return nil, but did!")
	}

	// subscriptionData without subscrCond,reqNfFqdn and reqNfType is not empty, and both reqNfType and reqNfFqdn match AllowedNfType and AllowedNfFqdn configured in CM
	cm.NrfPolicy.ManagementService.Subscription.AllowedSubscriptionAllNFs[1].AllowedNfDomains = "ericsson.se"
	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails = subscriptionData.ValidateSubscrCond()
	if problemDetails != nil {
		t.Fatalf("SubscriptionData.ValidateSubscrCond should not return nil, but did!")
	}
}

func TestValidateValidityTime(t *testing.T) {
	// subscriptionData without validityTime is valid
	body := []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		},
		"subscrCond": {
			"nfType": "AMF"
		}
	}`)

	subscriptionData := &TSubscriptionData{}
	err := json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails := subscriptionData.ValidateValidityTime()
	if problemDetails != nil {
		t.Fatalf("SubscriptionData.ValidateValidityTime should return nil, but not!")
	}

	// subscriptionData with validityTime which is before now is invalid
	body = []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"validityTime": "2018-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		},
		"subscrCond": {
			"nfType": "AMF"
		}
	}`)

	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails = subscriptionData.ValidateValidityTime()
	if problemDetails == nil {
		t.Fatalf("SubscriptionData.ValidateValidityTime should not return nil, but did!")
	}

	// subscriptionData with validityTime which is after now is valid
	body = []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"validityTime": "3015-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		},
		"subscrCond": {
			"nfType": "AMF"
		}
	}`)

	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails = subscriptionData.ValidateValidityTime()
	if problemDetails != nil {
		t.Fatalf("SubscriptionData.ValidateValidityTime should return nil, but not!")
	}
}

func TestValidateNotifCondition(t *testing.T) {
	// subscriptionData without notifCondition is valid
	body := []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		},
		"subscrCond": {
			"nfType": "AMF"
		}
	}`)

	subscriptionData := &TSubscriptionData{}
	err := json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails := subscriptionData.ValidateNotifCondition()
	if problemDetails != nil {
		t.Fatalf("SubscriptionData.ValidateNotifCondition should return nil, but not!")
	}

	// subscriptionData with valid notifCondition is valid
	body = []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"validityTime": "3015-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		},
		"subscrCond": {
			"nfType": "AMF"
		},
		"notifCondition": {
			"monitoredAttributes": ["load", "plmnList"]
		}
	}`)

	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails = subscriptionData.ValidateNotifCondition()
	if problemDetails != nil {
		t.Fatalf("SubscriptionData.ValidateNotifCondition should return nil, but not!")
	}

	// subscriptionData with invalid notifCondition is invalid
	body = []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"validityTime": "3015-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		},
		"subscrCond": {
			"nfType": "AMF"
		},
		"notifCondition": {
			"monitoredAttributes": ["load", "plmnList"],
			"unmonitoredAttributes": ["nfType", "nfStatus"]
		}
	}`)

	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails = subscriptionData.ValidateNotifCondition()
	if problemDetails == nil {
		t.Fatalf("SubscriptionData.ValidateNotifCondition should not return nil, but did!")
	}
}

func TestSubscriptionValidate(t *testing.T) {

	// subscriptionData with invalid nfStatusNotificationUri is invalid
	body := []byte(`{
		"nfStatusNotificationUri": "ftp://10.10.10.10:80/callback",
		"validityTime": "3015-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		},
		"subscrCond": {
			"nfType": "AMF"
		}
	}`)

	subscriptionData := &TSubscriptionData{}
	err := json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails := subscriptionData.Validate()
	if problemDetails == nil {
		t.Fatalf("SubscriptionData.Validate should not return nil, but did!")
	}

	// subscriptionData with invalid validityTime is invalid
	body = []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"validityTime": "2018-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		},
		"subscrCond": {
			"nfType": "AMF"
		}
	}`)

	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails = subscriptionData.Validate()
	if problemDetails == nil {
		t.Fatalf("SubscriptionData.Validate should not return nil, but did!")
	}

	// subscriptionData with invalid subscrCond is invalid
	body = []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"validityTime": "3015-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		},
		"subscrCond": {
			"nfInstanceId": "amf01",
			"nfType": "AMF"
		}
	}`)

	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails = subscriptionData.Validate()
	if problemDetails == nil {
		t.Fatalf("SubscriptionData.Validate should not return nil, but did!")
	}

	// subscriptionData with invalid notifCondition is invalid
	body = []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"validityTime": "3015-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		},
		"subscrCond": {
			"nfInstanceId": "amf01"
		},
		"notifCondition": {
			"monitoredAttributes": ["load", "plmnList"],
			"unmonitoredAttributes": ["nfType", "nfStatus"]
		}
	}`)

	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails = subscriptionData.Validate()
	if problemDetails == nil {
		t.Fatalf("SubscriptionData.Validate should not return nil, but did!")
	}

	// subscriptionData is invalid
	body = []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"validityTime": "3015-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"plmnId": {
		    "mcc": "460",
			"mnc": "00"
		},
		"subscrCond": {
			"nfInstanceId": "amf01"
		},
		"notifCondition": {
			"monitoredAttributes": ["load", "plmnList"]
		}
	}`)

	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}
	problemDetails = subscriptionData.Validate()
	if problemDetails != nil {
		t.Fatalf("SubscriptionData.Validate should return nil, but not!")
	}
}

func TestSubscriptionIsLocalPlmn(t *testing.T) {
	cm.NfProfile.PlmnID = []cm.TPLMN{
		cm.TPLMN{
			Mcc: "460",
			Mnc: "00",
		},
		cm.TPLMN{
			Mcc: "460",
			Mnc: "11",
		},
	}

	// subscriptionData without plmnId is local
	body := []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"validityTime": "3015-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"subscrCond": {
			"nfInstanceId": "amf01"
		}
	}`)

	subscriptionData := &TSubscriptionData{}
	err := json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}

	if !subscriptionData.IsLocalPlmn() {
		t.Fatalf("it is a local subscriptionData, but SubscriptionData.IsLocalPlmn return false!")
	}

	// subscriptionData with plmnId which doesn't belong to local plmn list is not local
	body = []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"validityTime": "3015-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"subscrCond": {
			"nfInstanceId": "amf01"
		},
		"plmnId": {
			"mcc": "460",
			"mnc": "22"
		}
	}`)

	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}

	if subscriptionData.IsLocalPlmn() {
		t.Fatalf("it is not a local subscriptionData, but SubscriptionData.IsLocalPlmn return true!")
	}

	// subscriptionData with plmnId which belongs to local plmn list is local
	body = []byte(`{
		"nfStatusNotificationUri": "http://10.10.10.10:80/callback",
		"validityTime": "3015-12-31T23:59:59.999Z",
		"reqNotifEvents": ["NF_REGISTERED", "NF_DEREGISTERED"],
		"subscrCond": {
			"nfInstanceId": "amf01"
		},
		"plmnId": {
			"mcc": "460",
			"mnc": "00"
		}
	}`)

	subscriptionData = &TSubscriptionData{}
	err = json.Unmarshal(body, subscriptionData)
	if err != nil {
		t.Fatalf("Unmarshal subscriptionData error")
	}

	if !subscriptionData.IsLocalPlmn() {
		t.Fatalf("it is a local subscriptionData, but SubscriptionData.IsLocalPlmn return false!")
	}

}

func TestConstructSubscriptionIndex(t *testing.T) {
	acceptedValidityTime := "2018-12-31T23:59:59.999Z"
	unacceptedValidityTime := "3015-12-31T23:59:59.999Z"

	timeInSecond, _ := time.Parse(time.RFC3339, acceptedValidityTime)
	acceptedValidityTimeInMilSec := timeInSecond.Unix() * 1000

	// subscriptionData without subscrCond
	subscriptionData := &TSubscriptionData{
		NfStatusNotificationUri: "http://10.10.10.10:80/callback",
		ValidityTime:            acceptedValidityTime,
	}

	subscriptionPutIndex := subscriptionData.ConstructSubscriptionIndex()

	if subscriptionPutIndex == nil {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex should not return nil, but did!")
	}

	if subscriptionPutIndex.NfStatusNotificationUri != subscriptionData.NfStatusNotificationUri {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.NoCond != constvalue.NoSubscrCond {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.NfInstanceId != constvalue.Wildcard {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.NfType != constvalue.Wildcard {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.ServiceName != constvalue.Wildcard {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.AmfCond == nil {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	amfCond := subscriptionPutIndex.AmfCond

	if amfCond.SubKey1 != constvalue.Wildcard || amfCond.SubKey2 != constvalue.Wildcard {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.GuamiList == nil || len(subscriptionPutIndex.GuamiList) != 1 {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	guamiList := subscriptionPutIndex.GuamiList

	if guamiList[0].SubKey1 != constvalue.Wildcard || guamiList[0].SubKey2 != constvalue.Wildcard {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.SnssaiList == nil || len(subscriptionPutIndex.SnssaiList) != 1 {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	snssaiList := subscriptionPutIndex.SnssaiList

	if snssaiList[0].SubKey1 != constvalue.Wildcard || snssaiList[0].SubKey2 != constvalue.Wildcard {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.NsiList == nil || len(subscriptionPutIndex.NsiList) != 1 {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	nsiList := subscriptionPutIndex.NsiList

	if nsiList[0] != constvalue.Wildcard {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.NfGroupCond == nil {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	nfGroupCond := subscriptionPutIndex.NfGroupCond

	if nfGroupCond.SubKey1 != constvalue.Wildcard || nfGroupCond.SubKey2 != constvalue.Wildcard {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.ValidityTime != uint64(acceptedValidityTimeInMilSec) {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	// subscriptionData with subscrCond which includes multiple conditions
	subscriptionData = &TSubscriptionData{
		NfStatusNotificationUri: "http://10.10.10.10:80/callback",
		ValidityTime:            unacceptedValidityTime,
		SubscrCond: &TSubscrCond{
			NfInstanceID: "amf01",
			NfType:       "AMF",
			ServiceName:  "serv-amf",
			AmfSetID:     "amfSet01",
			AmfRegionID:  "amfRegion01",
			GuamiList: []TGuami{
				TGuami{
					PlmnId: TPlmnId{
						Mcc: "460",
						Mnc: "00",
					},
					AmfId: "123456",
				},
				TGuami{
					PlmnId: TPlmnId{
						Mcc: "460",
						Mnc: "11",
					},
					AmfId: "234567",
				},
			},
			SnssaiList: []TSnssai{
				TSnssai{
					Sst: 1,
					Sd:  "123456",
				},
				TSnssai{
					Sst: 2,
					Sd:  "234567",
				},
			},
			NsiList:   []string{"nsi1", "nsi2"},
			NfGroupID: "group01",
		},
	}

	subscriptionPutIndex = subscriptionData.ConstructSubscriptionIndex()

	if subscriptionPutIndex == nil {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex should not return nil, but did!")
	}

	if subscriptionPutIndex.NfStatusNotificationUri != subscriptionData.NfStatusNotificationUri {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.NoCond != constvalue.Wildcard {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.NfInstanceId != "amf01" {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.NfType != constvalue.Wildcard {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.ServiceName != "serv-amf" {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.AmfCond == nil {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	amfCond = subscriptionPutIndex.AmfCond

	if amfCond.SubKey1 != "amfSet01" || amfCond.SubKey2 != "amfRegion01" {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.GuamiList == nil || len(subscriptionPutIndex.GuamiList) != 2 {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	guamiList = subscriptionPutIndex.GuamiList

	ok := false

	if (guamiList[0].SubKey1 == "46000" && guamiList[0].SubKey2 == "123456" &&
		guamiList[1].SubKey1 == "46011" && guamiList[1].SubKey2 == "234567") ||
		(guamiList[0].SubKey1 == "46011" && guamiList[0].SubKey2 == "234567" &&
			guamiList[1].SubKey1 == "46000" && guamiList[1].SubKey2 == "123456") {
		ok = true
	}

	if !ok {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.SnssaiList == nil || len(subscriptionPutIndex.SnssaiList) != 2 {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	snssaiList = subscriptionPutIndex.SnssaiList

	ok = false

	if (snssaiList[0].SubKey1 == "1" && snssaiList[0].SubKey2 == "123456" &&
		snssaiList[1].SubKey1 == "2" && snssaiList[1].SubKey2 == "234567") ||
		(snssaiList[0].SubKey1 == "2" && snssaiList[0].SubKey2 == "234567" &&
			snssaiList[1].SubKey1 == "1" && snssaiList[1].SubKey2 == "123456") {
		ok = true
	}

	if !ok {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.NsiList == nil || len(subscriptionPutIndex.NsiList) != 2 {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	nsiList = subscriptionPutIndex.NsiList

	ok = false

	if (nsiList[0] == "nsi1" && nsiList[1] == "nsi2") ||
		(nsiList[0] == "nsi2" && nsiList[1] == "nsi1") {
		ok = true
	}

	if !ok {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.NfGroupCond == nil {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	nfGroupCond = subscriptionPutIndex.NfGroupCond

	if nfGroupCond.SubKey1 != "group01" || nfGroupCond.SubKey2 != "AMF" {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}

	if subscriptionPutIndex.ValidityTime == uint64(acceptedValidityTimeInMilSec) {
		t.Fatalf("TSubscriptionData.ConstructSubscriptionIndex didn't return right value!")
	}
}
