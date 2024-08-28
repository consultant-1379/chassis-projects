package internalconf

import (
	"fmt"
	"sync"

	"gerrit.ericsson.se/udm/common/pkg/httpserver"
)

const (
	defaultTokenExpiredTime = 7200
)

// TokenExpiredTime
var (
	TokenExpiredTime    = defaultTokenExpiredTime
	AllowUnknownService = false
	TokenKeyID          = "keyId1"

	NfServiceMapLock       sync.RWMutex
	NfServiceToNfType      = make(map[string]string)
	NfTypeToNfservices     = make(map[string][]string)
	NfServiceAllowedNfType = make(map[string][]string)
)

// InternalAccessTokenConf is format of internalConf.json in NRF Access Token service
type InternalAccessTokenConf struct {
	HTTPServer  HttpServerConf                    `json:"httpServer"`
	Dbproxy     DbproxyConf                       `json:"dbProxy"`
	AccessToken AccessTokenConf                   `json:"accessToken"`
	NfServices  map[string]map[string]AllowedInfo `json:"nfServices"`
}

// AccessTokenConf is used to store provision service info
type AccessTokenConf struct {
	ExpiredTime         int    `json:"expiredTime,omitempty"`
	KeyID               string `json:"keyId,omitempty"`
	AllowUnknownService bool   `json:"allowUnknownService"`
}

// AllowedInfo is used for allowed Info
type AllowedInfo struct {
	AllowedNfType []string `json:"allowed-nf-type"`
}

// ParseConf sets the value of global variable from struct
// Need add RWLock
func (conf *InternalAccessTokenConf) ParseConf() {
	HTTPWithXVersion = conf.HTTPServer.HTTPWithXVersion

	httpServerIdleTimeout := conf.HTTPServer.IdleTimeout
	if httpServerIdleTimeout <= 0 {
		httpServerIdleTimeout = defaultHttpServerIdleTimeout
	}
	httpserver.SetIdleTimeout(httpServerIdleTimeout)

	httpServerActiveTimeout := conf.HTTPServer.ActiveTimeout
	if httpServerActiveTimeout <= 0 {
		httpServerActiveTimeout = defaultHttpServerActiveTimeout
	}
	httpserver.SetActiveTimeout(httpServerActiveTimeout)

	HTTPWithXVersion = conf.HTTPServer.HTTPWithXVersion

	DbproxyConnectionNum = conf.Dbproxy.ConnectionNum
	if DbproxyConnectionNum <= 0 {
		DbproxyConnectionNum = defaultDbproxyConnectionNum
	}

	DbproxyGrpcCtxTimeout = conf.Dbproxy.GrpcCtxTimeout
	if DbproxyGrpcCtxTimeout <= 0 {
		DbproxyGrpcCtxTimeout = defaultDbproxyGrpcCtxTimeout
	}

	TokenExpiredTime = conf.AccessToken.ExpiredTime
	if TokenExpiredTime <= 0 {
		TokenExpiredTime = defaultTokenExpiredTime
	}

	TokenKeyID = conf.AccessToken.KeyID

	AllowUnknownService = conf.AccessToken.AllowUnknownService

	{
		NfServiceMapLock.Lock()
		for nftype, services := range conf.NfServices {
			delete(NfTypeToNfservices, nftype)
			for service, allowedInfo := range services {
				delete(NfServiceToNfType, service)
				delete(NfServiceAllowedNfType, service)
				NfServiceToNfType[service] = nftype
				NfTypeToNfservices[nftype] = append(NfTypeToNfservices[nftype], service)
				NfServiceAllowedNfType[service] = append(NfServiceAllowedNfType[service], allowedInfo.AllowedNfType...)
			}
		}
		NfServiceMapLock.Unlock()
	}
}

// ShowInternalConf shows the value of global variable
func (conf *InternalAccessTokenConf) ShowInternalConf() {
	fmt.Printf("http header with x-version : %v\n", HTTPWithXVersion)
	fmt.Printf("Token Expired Time is: %d\n", TokenExpiredTime)
	fmt.Printf("Token key id is %v", TokenKeyID)
	fmt.Printf("AllowUnknownService is: %v\n", AllowUnknownService)
	fmt.Printf("NfservicesToNfType is: %v\n", NfServiceToNfType)
	fmt.Printf("NfTypeToNfservices is: %v\n", NfTypeToNfservices)
	fmt.Printf("NfServiceAllowedNfType is: %v\n", NfServiceAllowedNfType)
}
