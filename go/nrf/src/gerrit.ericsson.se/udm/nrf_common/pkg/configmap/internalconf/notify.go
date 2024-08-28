package internalconf

import (
	"fmt"

	"gerrit.ericsson.se/udm/common/pkg/httpserver"
)

var (
	// EnableNotification is a switch to enable/disbale trigger notification
	EnableNotification = true

	// callBackURISupportIPv6 is to control whether callbackURI support ipv6
	CallBackURISupportIPv6 = true

	// EnableLongConnection is a switch to used long connection or short connection to NF
	EnableLongConnection = true

	// NotificationAlwaysWithFullNfProfile is to control whether notification body include either full nfProfile or patch of nfProfile when nfProfile change
	NotificationAlwaysWithFullNfProfile = false

	// NotificationMaxJobWorker is the maximum worker for handling notification job
	NotificationMaxJobWorker = 20

	// NotificationMaxJobQueue is the maximum queue for notification job
	NotificationMaxJobQueue = 3000

	// NotificationMaxNotifyWorker is the maximum worker for sending notification
	NotificationMaxNotifyWorker = 5

	// NotificationTimeout is timeout time for notification job in queue
	NotificationTimeout = 3

	// EnableDebugLogForOverload enable debug log for notificatino overload
	EnableDebugLogForOverload = false
)

// InternalNotifyConf is format of internalConf.json in NRF Notification
type InternalNotifyConf struct {
	HttpServer     HttpServerConf     `json:"httpServer"`
	Dbproxy        DbproxyConf        `json:"dbProxy"`
	NFStatusNotify NFStatusNotifyConf `json:"NFStatusNotify"`
}

// NFStatusNotifyConf is format of "NFStatusNotify" in internalConf.json
type NFStatusNotifyConf struct {
	Enable                              bool `json:"enable"`
	CallBackURISupportIPv6              bool `json:"callBackURISupportIPv6"`
	EnableLongConnection                bool `json:"enableLongConnection"`
	NotificationAlwaysWithFullNfProfile bool `json:"notificationAlwaysWithFullNfProfile"`
	NotificationMaxJobQueue             int  `json:"notificationMaxJobQueue"`
	NotificationMaxJobWorker            int  `json:"notificationMaxJobWorker"`
	NotificationMaxNotifyWorker         int  `json:"notificationMaxNotifyWorker"`
	NotificationTimeout                 int  `json:"notificationTimeout"`
	EnableDebugLogForOverload           bool `json:"enableDebugLogForOverload"`
}

// ParseConf sets the value of global variable from struct
func (conf *InternalNotifyConf) ParseConf() {

	httpServerIdleTimeout := conf.HttpServer.IdleTimeout
	if httpServerIdleTimeout <= 0 {
		httpServerIdleTimeout = defaultHttpServerIdleTimeout
	}
	httpserver.SetIdleTimeout(httpServerIdleTimeout)

	httpServerActiveTimeout := conf.HttpServer.ActiveTimeout
	if httpServerActiveTimeout <= 0 {
		httpServerActiveTimeout = defaultHttpServerActiveTimeout
	}
	httpserver.SetActiveTimeout(httpServerActiveTimeout)

	HTTPWithXVersion = conf.HttpServer.HTTPWithXVersion

	DbproxyConnectionNum = conf.Dbproxy.ConnectionNum
	if DbproxyConnectionNum <= 0 {
		DbproxyConnectionNum = defaultDbproxyConnectionNum
	}

	DbproxyGrpcCtxTimeout = conf.Dbproxy.GrpcCtxTimeout
	if DbproxyGrpcCtxTimeout <= 0 {
		DbproxyGrpcCtxTimeout = defaultDbproxyGrpcCtxTimeout
	}

	EnableNotification = conf.NFStatusNotify.Enable

	CallBackURISupportIPv6 = conf.NFStatusNotify.CallBackURISupportIPv6

	EnableLongConnection = conf.NFStatusNotify.EnableLongConnection

	NotificationAlwaysWithFullNfProfile = conf.NFStatusNotify.NotificationAlwaysWithFullNfProfile

	if conf.NFStatusNotify.NotificationMaxJobWorker > 5 {
		NotificationMaxJobWorker = conf.NFStatusNotify.NotificationMaxJobWorker
	}

	if conf.NFStatusNotify.NotificationMaxJobQueue > 500 {
		NotificationMaxJobQueue = conf.NFStatusNotify.NotificationMaxJobQueue
	}

	if conf.NFStatusNotify.NotificationMaxNotifyWorker > 0 {
		NotificationMaxNotifyWorker = conf.NFStatusNotify.NotificationMaxNotifyWorker
	}

	if conf.NFStatusNotify.NotificationTimeout > 0 {
		NotificationTimeout = conf.NFStatusNotify.NotificationTimeout
	}

	EnableDebugLogForOverload = conf.NFStatusNotify.EnableDebugLogForOverload

}

// ShowInternalConf shows the value of global variable
func (conf *InternalNotifyConf) ShowInternalConf() {
	fmt.Printf("http server idleTimeout : %v\n", httpserver.GetIdleTimeout())
	fmt.Printf("http server activeTimeout : %v\n", httpserver.GetActiveTimeout())
	fmt.Printf("http header with x-version : %v\n", HTTPWithXVersion)
	fmt.Printf("DbproxyConnectionNum  : %d\n", DbproxyConnectionNum)
	fmt.Printf("DbproxyGrpcCtxTimeout : %d\n", DbproxyGrpcCtxTimeout)
	fmt.Printf("EnableNotification : %v\n", EnableNotification)
	fmt.Printf("CallBackURISupportIPv6 : %v\n", CallBackURISupportIPv6)
	fmt.Printf("EnableLongConnection : %v\n", EnableLongConnection)
	fmt.Printf("NotificationAlwaysWithFullNfProfile : %v\n", NotificationAlwaysWithFullNfProfile)
	fmt.Printf("NotificationMaxJobWorker : %d\n", NotificationMaxJobWorker)
	fmt.Printf("NotificationMaxJobQueue : %d\n", NotificationMaxJobQueue)
	fmt.Printf("NotificationMaxNotifyWorker : %d\n", NotificationMaxNotifyWorker)
	fmt.Printf("NotificationTimeout : %d\n", NotificationTimeout)
	fmt.Printf("EnableDebugLogForOverload : %v\n", EnableDebugLogForOverload)
}
