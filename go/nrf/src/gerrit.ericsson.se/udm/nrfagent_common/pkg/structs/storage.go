package structs

import (
	"time"
)

//ConfigMapData define the data structure in configmap
//NfInfoList is actually an array of NfInfo
//Defined as string is easy to replace all data in one patch
type ConfigMapData struct {
	NfInfoList           string `json:"nfInfoList"`
	SubscriptionInfoList string `json:"subscriptionInfoList"`
}

//NfInfo define the data structure of NF info
type NfInfo struct {
	NfInstanceID  string   `json:"nfInstanceId"`
	NfType        string   `json:"nfType"`
	NfFqdn        string   `json:"nfFqdn"`
	NfPlmns       []PlmnID `json:"nfPlmns,omitempty"`
	NRFHBInterval int64    `json:"heartbeatTimer,omitempty"`
}

//SubscriptionInfo define the data structure of Subscription info
type SubscriptionInfo struct {
	RequesterNfType   string          `json:"requesterNfType"`
	TargetNfType      string          `json:"targetNfType"`
	TargetServiceName string          `json:"targetServiceName,omitempty"`
	NfInstanceID      string          `json:"nfInstanceID,omitempty"`
	TargetPlmnID      PlmnID          `json:"targetPlmnID,omitempty"`
	NotifCondition    *NotifCondition `json:"notifCondition,omitempty"`
	SubscriptionID    string          `json:"subscriptionID"`
	ValidityTime      time.Time       `json:"validityTime,omitempty"`
}
