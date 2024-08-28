package cm

import (
	"math"
	"strings"

	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
)

// EricssonNrfFunction is struct of nrf's configuration
type EricssonNrfFunction struct {
	NrfCommon         *TCommon            `json:"common"`
	NfProfile         *TNfProfile         `json:"nf-profile"`
	ManagementService *TManagementService `json:"management-service"`
	NrfPolicy         *TNrfPolicy         `json:"policy"`
	DiscoveryService  *TDiscoveryService  `json:"discovery-service"`
	NfServiceLogs     []TNfServiceLog     `json:"nf-service-log"`
	ProvisionService  *TProvisionService  `json:"provisioning-service"`
}

// --------------------------------------------struct list for common -------------------------------------------------------------

// TCommon is configuration of nrf common
type TCommon struct {
	Role                 string                 `json:"role"`
	PlmnNrf              *TPlmnNrf              `json:"plmn-nrf"`
	RemoteDefaultSetting *TRemoteDefaultSetting `json:"remote-default-setting"`
	GeoRed               *TGeoRed               `json:"geo-red"`
}

// TPlmnNrf is configuration of nrf plmn
type TPlmnNrf struct {
	Mode    string        `json:"mode"`
	Profile []TNrfProfile `json:"profile"`
}

// TNrfProfile is configuration of nrf profile
type TNrfProfile struct {
	ID          string        `json:"id"`
	Fqdn        string        `json:"fqdn"`
	Ipv4Address []string      `json:"ipv4-address"`
	Ipv6Address []string      `json:"ipv6-address"`
	Priority    *int          `json:"priority,omitempty"`
	Capacity    *int          `json:"capacity,omitempty"`
	locality    string        `json:"locality"`
	Service     []TNrfService `json:"service"`
}

// TNrfService is configuration of nrf service
type TNrfService struct {
	ID                int           `json:"id"`
	Name              string        `json:"name"`
	Version           []TVersion    `json:"version"`
	Scheme            string        `json:"scheme"`
	Fqdn              string        `json:"fqdn"`
	IPEndpoint        []TIPEndpoint `json:"ip-endpoint"`
	Priority          *int          `json:"priority,omitempty"`
	Capacity          *int          `json:"capacity,omitempty"`
	APIPrefix         string        `json:"api-prefix"`
	SupportedFeatures string        `json:"supported-features"`
}

// TVersion is configuration of nfService's version
type TVersion struct {
	APIVersionInURI string `json:"api-version-in-uri"`
	APIFullVersion  string `json:"api-full-version"`
	Expiry          string `json:"expiry"`
}

// TIPEndpoint is configuration of end point
type TIPEndpoint struct {
	ID          int    `json:"id"`
	Transport   string `json:"transport"`
	Port        int    `json:"port"`
	Ipv4Address string `json:"ipv4-address"`
	Ipv6Address string `json:"ipv6-address,omitempty"`
}

// TRemoteDefaultSetting is configuration of egress
type TRemoteDefaultSetting struct {
	Scheme string `json:"scheme"`
	Port   int    `json:"port"`
}

// TGeoRed is configuration of geoRed function
type TGeoRed struct {
	KeepManagementService bool        `json:"keep-management-service"`
	KeepDiscoveryService  bool        `json:"keep-discovery-service"`
	WitnessNF             *TWitnessNF `json:"witness-nf"`
}

// TWitnessNF is configuration of geoRed witness
type TWitnessNF struct {
	IdentityType  string `json:"identity-type"`
	IdentityValue string `json:"identity-value"`
}

/*
 * --------------------------------------------struct list for nf-profile -------------------------------------------------------------
 */

// TNfProfile is configuration of nf profile
type TNfProfile struct {
	Service                 []TNfService `json:"service,omitempty"`
	InstanceID              string       `json:"instance-id"`
	Type                    string       `json:"type",omitempty`
	Status                  string       `json:"status"`
	RequestedHeartbeatTimer *int         `json:"requested-heartbeat-timer,omitempty"`
	PlmnID                  []TPLMN      `json:"plmn-id,omitempty"`
	Snssai                  []TSnssai    `json:"snssai,omitempty"`
	Nsi                     []string     `json:"nsi,omitempty"`
	Fqdn                    string       `json:"fqdn,omitempty"`
	InterPlmnFqdn           string       `json:"inter-plmn-fqdn,omitempty"`
	Ipv4Address             []string     `json:"ipv4-address,omitempty"`
	Ipv6Address             []string     `json:"ipv6-address,omitempty"`
	AllowedPlmn             []TPLMN      `json:"allowed-plmn,omitempty"`
	AllowedNfType           []string     `json:"allowed-nf-type,omitempty"`
	AllowedNfDomain         []string     `json:"allowed-nf-domain,omitempty"`
	AllowedNssai            []TSnssai    `json:"allowed-nssai,omitempty"`
	Priority                *int         `json:"priority,omitempty"`
	Capacity                *int         `json:"capacity,omitempty"`
	Locality                string       `json:"locality,omitempty"`
	RecoveryTime            string       `json:"recovery-time,omitempty"`
	ServicePersistence      *bool        `json:"service-persistence,omitempty"`
}

// TNfService is configuration of nf service
type TNfService struct {
	InstanceID                      string                             `json:"instance-id"`
	Name                            string                             `json:"name"`
	Version                         []TVersion                         `json:"version"`
	Scheme                          string                             `json:"scheme"`
	Status                          string                             `json:"status"`
	Fqdn                            string                             `json:"fqdn,omitempty"`
	InterPlmnFqdn                   string                             `json:"inter-plmn-fqdn,omitempty"`
	IPEndpoint                      []TIPEndpoint                      `json:"ip-endpoint,omitempty"`
	APIPrefix                       string                             `json:"api-prefix,omitempty"`
	DefaultNotificationSubscription []TDefaultNotificationSubscription `json:"default-notification-subscription,omitempty"`
	AllowedPlmn                     []TPLMN                            `json:"allowed-plmn,omitempty"`
	AllowedNfType                   []string                           `json:"allowed-nf-type,omitempty"`
	AllowedNfDomain                 []string                           `json:"allowed-nf-domain,omitempty"`
	AllowedNssai                    []TSnssai                          `json:"allowed-nssai,omitempty"`
	Priority                        *int                               `json:"priority,omitempty"`
	Capacity                        *int                               `json:"capacity,omitempty"`
	RecoveryTime                    string                             `json:"recovery-time,omitempty"`
	SupportedFeatures               string                             `json:"supported-features,omitempty"`
}

// TDefaultNotificationSubscription is configuration of nfService's default-notification-subscription
type TDefaultNotificationSubscription struct {
	NotificationType   string `json:"notification-type"`
	CallbackURI        string `json:"callback-uri"`
	N1MessageClass     string `json:"n1-message-class,omitempty"`
	N2InformationClass string `json:"n2-information-class,omitempty"`
}

// TPLMN is configuration of PLMN
type TPLMN struct {
	Mcc string `json:"mcc"`
	Mnc string `json:"mnc"`
}

// GetPlmnID returns mcc+mnc
func (p *TPLMN) GetPlmnID() string {
	return p.Mcc + p.Mnc
}

// TSnssai is configuration of Snssai
type TSnssai struct {
	ID  int    `json:"id"`
	Sst int    `json:"sst"`
	Sd  string `json:"sd,omitempty"`
}

/*
 * --------------------------------------------struct list for managment-service --------------------------
 * --------------------------------------------struct list for discovery-service ---------------------------
 */

// TManagementService is configuration of management service profile
type TManagementService struct {
	Heartbeat               *THeartbeat `json:"heartbeat"`
	SubscriptionExpiredTime int         `json:"subscription-expired-time"`
}

// THeartbeat is for heartbeat
type THeartbeat struct {
	Default          int                 `json:"default"`
	DefaultPerNfType []TDefaultPerNfType `json:"default-per-nftype"`
	GracePeriod      int                 `json:"grace-period"`
}

// TDefaultPerNfType is for heartbeat per nf-type
type TDefaultPerNfType struct {
	NfType string `json:"nf-type"`
	Value  int    `json:"value"`
}

// GetDefaultHeartbeatTimer returns the default heartbeat timer
func (heartbeat *THeartbeat) GetDefaultHeartbeatTimer(nfType string) int {
	if len(heartbeat.DefaultPerNfType) > 0 {
		for _, heartbeatTimerPerNfType := range heartbeat.DefaultPerNfType {
			if strings.EqualFold(heartbeatTimerPerNfType.NfType, nfType) {
				return int(math.Max(float64(heartbeatTimerPerNfType.Value), float64(internalconf.HeartBeatTimerMin)))
			}
		}
	}
	return int(math.Max(float64(heartbeat.Default), float64(internalconf.HeartBeatTimerMin)))
}

// TDiscoveryService is configuration of discovery service
type TDiscoveryService struct {
	ResponseCacheTime  int    `json:"response-cache-time"`
	HierarchyMode      string `json:"hierarchy-mode"`
	LocalCacheEnable   bool   `json:"local-cache-enable"`
	LocalCacheTimeout  int    `json:"local-cache-timeout"`
	LocalCacheCapacity int    `json:"local-cache-capacity"`
}

// TProvisionService is configuration of provision service
type TProvisionService struct {
	ProvAddress []TProvAddress `json:"prov-address"`
}

// ServiceInfo is configuration of service info
type ServiceInfo struct {
	Type            string `json:"type"`
	ProductRevision string `json:"product-revision"`
	ProductNumber   string `json:"product-number"`
	Description     string `json:"description"`
	ProductName     string `json:"product-name"`
	ProductionDate  string `json:"production-date"`
}

//TProvAddress is for provision address
type TProvAddress struct {
	ID          int    `json:"id"`
	Scheme      string `json:"scheme"`
	Fqdn        string `json:"fqdn"`
	Transport   string `json:"transport"`
	Port        int    `json:"port"`
	Ipv4Address string `json:"ipv4-address"`
	Ipv6Address string `json:"ipv6-address"`
}

// TNrfPolicy is  the policy of NRF
type TNrfPolicy struct {
	ManagementService *TNrfManagementServicePolicy `json:"management-service"`
}

// TNrfManagementServicePolicy is  the policy of NRF management
type TNrfManagementServicePolicy struct {
	Subscription *TSubscriptionPolicy `json:"subscription"`
}

// TSubscriptionPolicy is the subscription policy
type TSubscriptionPolicy struct {
	AllowedSubscriptionAllNFs []TAllowedSubscriptionAllNFs `json:"allowed-subscription-to-all-nfs"`
}

// TAllowedSubscriptionAllNFs is for the NFs that are allowed to subscribe to all NFs
type TAllowedSubscriptionAllNFs struct {
	AllowedNfType    string `json:"allowed-nf-type"`
	AllowedNfDomains string `json:"allowed-nf-domains"`
}

/*
 * --------------------------------------------struct list for nf-service-log -------------------------------------------------------------
 */

// TNfServiceLog is configuration of service log
type TNfServiceLog struct {
	LogID    string   `json:"log-id,omitempty"`
	Severity string   `json:"severity,omitempty"`
	PodLogs  []PodLog `json:"pod-log"`
}

// PodLog is configuration of pod log
type PodLog struct {
	PodID    string `json:"pod-id,omitempty"`
	Severity string `json:"severity,omitempty"`
}

/*
 * --------------------------------------------define the interface for cm parse ------------------------------
 */

// ConfigCM is interface for nrf configuration
type ConfigCM interface {
	ParseConf()
	Show()
}

var (
	// PodIP is pod Id itself
	PodIP string

	//ServiceName is name of microservice
	ServiceName string
)

// SetPodIP is to set pod Id itself
func SetPodIP(podIP string) {
	PodIP = podIP
}

// SetPodIP is to set pod Id itself
func SetServiceName(name string) {
	ServiceName = name
}
