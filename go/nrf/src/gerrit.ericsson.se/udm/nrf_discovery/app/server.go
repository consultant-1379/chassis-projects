package nrf

import (
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/httpserver"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/pm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/client"
	"gerrit.ericsson.se/udm/nrf_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap"
	"gerrit.ericsson.se/udm/nrf_common/pkg/configmap/internalconf"
	"gerrit.ericsson.se/udm/nrf_common/pkg/constvalue"
	"gerrit.ericsson.se/udm/nrf_common/pkg/dbmgmt"
	"gerrit.ericsson.se/udm/nrf_common/pkg/fm"
	"gerrit.ericsson.se/udm/nrf_common/pkg/fsnotify"
	"gerrit.ericsson.se/udm/nrf_common/pkg/multisite"
	"gerrit.ericsson.se/udm/nrf_common/pkg/options"
	"gerrit.ericsson.se/udm/nrf_common/pkg/probe"
	"gerrit.ericsson.se/udm/nrf_common/pkg/utils"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdisccache"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscservice"
	"gerrit.ericsson.se/udm/nrf_discovery/app/nfdiscovery/nfdiscutil"
	"golang.org/x/net/http2"
)

//Server struct
type Server struct {
	opts           *options.Options
	sigurs1_handle bool
	sig_mutex      sync.Mutex
	Terminate      chan int
}

//NewServer function for server init
func NewServer(opts *options.Options) *Server {
	ch := make(chan int)
	s := &Server{
		opts:           opts,
		sigurs1_handle: false,
		Terminate:      ch,
	}

	s.Init()

	return s
}

//Init initialize server
func (s *Server) Init() {

	//initialize log
	log.SetLevel(log.Level(s.opts.LogLevel))
	log.SetOutput(os.Stdout)

	log.SetServiceID(constvalue.APP_WORKMODE_NRF_DISC)
	log.SetNF("nrf")
	log.SetPodIP(s.opts.PodIP)
	log.SetFormatter(&log.JSONFormatter{})

	//initialize CM
	cm.SetPodIP(s.opts.PodIP)
	cm.SetServiceName(s.opts.WorkMode)
	cm.RegisterDiscLocalCacheHandler(nfdiscservice.LocalCacheCMUpdate)
	//initialize PM
	pm.Init(s.opts.MetricsServicePort, s.opts.MetricsPath)
	registerNRFDiscMetrics()

	//initialize FM
	fm.Init()

	cpuRequest := os.Getenv("CPU_REQUEST")
	goMaxProcs, err := strconv.Atoi(cpuRequest)
	if err == nil {
		log.Warnf("disc GOMAXPROCS = CPU_REQUEST * 2 = %v", goMaxProcs*2)
		runtime.GOMAXPROCS(goMaxProcs * 2)
	}

	//load configuration from configmap
	err = configmap.InitConfigMap()
	if err != nil {
		log.Errorf("Initialize configmap error")
		os.Exit(1)
	}

	//initialize fsnotify
	err = fsnotify.Init()
	if err != nil {
		log.Errorf("Initialize fsnotify error")
		os.Exit(1)
	}

	fsnotify.Run()

	//initialize MarkDiscLocalCacheCapacity, align with cm.DiscLocalCacheCapacity,
	// when DiscLocalCacheCapacity change, modify cache region count
	nfdiscservice.InitMarkDiscLocalCacheCapacity()
	//initilaize http client
	s.opts.ClientOpts.Timeout = 1
	client.InitHttpClient(s.opts.ClientOpts)

	//initialize dbmgmt
	dbmgmt.InitDB(s.opts.GrpcServer)

	//set prefix for sequenceID
	utils.SetPrefix(s.opts.HostName)

	//initialize nfprofile cache
	nfdisccache.InitNfProfileCache()

	//discovery do time statistics
	dbmgmt.StartCalculate()
}

//Stop is function to stop the server
func (s *Server) Stop() {
	log.Warningf("Exiting service ...")
	//defer DestroyDb()
	defer dbmgmt.Close()
	fsnotify.Stop()
}

//Run is function to start the server
func (s *Server) Run() {

	for {
		if configmap.InternalDiscConfInst.PriorityPolicy != nil {
			break
		}

		log.Warningf("Ingress Priority Policy is NOT available, wait for CM update")
		time.Sleep(3 * time.Second)
	}

	if internalconf.OverloadProtection.Enabled {
		http2.EngineManager = http2.NewWorkEngineManager()
		http2.EngineManager.SetOverloadControlLevel(internalconf.OverloadProtection.OverloadControlLevel)
		http2.EngineManager.SetOverloadTriggerLatencyThreshold(internalconf.OverloadProtection.OverloadTriggerLatencyThreshold)
		http2.EngineManager.SetOverloadControlLatencyThreshold(internalconf.OverloadProtection.OverloadControlLatencyThreshold)
		http2.EngineManager.SetOverloadTriggerSampleWindow(internalconf.OverloadProtection.OverloadTriggerSampleWindow)
		http2.EngineManager.SetOverloadControlSampleWindow(internalconf.OverloadProtection.OverloadControlSampleWindow)
		http2.EngineManager.SetIdleInterval(internalconf.OverloadProtection.IdleInterval)
		http2.EngineManager.SetIdleRecoverRatio(internalconf.OverloadProtection.IdleRecoverRatio)
		http2.EngineManager.SetDefaultMessagePriority(configmap.InternalDiscConfInst.PriorityPolicy.DefaultMessagePriority)
		http2.EngineManager.SetCounterReportInterval(internalconf.OverloadProtection.CounterReportInterval)
		http2.EngineManager.SetOverloadAlarmClearWindow(internalconf.OverloadProtection.OverloadAlarmClearWindow)
		http2.EngineManager.SetDeniedRequestWorkerNumber(16)
		http2.EngineManager.SetDeniedRequestQueueCapacity(40960)
		http2.EngineManager.SetStatisticsQueueCapacity(4096)
		for i := range internalconf.OverloadProtection.WorkEngines {
			for _, priorityGroup := range configmap.InternalDiscConfInst.PriorityPolicy.PriorityGroup {
				if priorityGroup.GroupPriority == internalconf.OverloadProtection.WorkEngines[i].GroupPriority {
					engine := http2.NewWorkEngine(priorityGroup.GroupPriority, priorityGroup.PriorityStart, priorityGroup.PriorityEnd, internalconf.OverloadProtection.WorkEngines[i].QueueCapacity, internalconf.OverloadProtection.WorkEngines[i].WorkerNumber, internalconf.OverloadProtection.OverloadControlLevel)
					http2.EngineManager.RegisterWorkEngine(engine)
					break
				}
			}
		}
		http2.EngineManager.Start()
		go http2.EngineManager.ReportCounterForNFDisc(constvalue.NfDiscoveryRequestsTotal, constvalue.NfDiscoveryFailureTotal)
	}
	s.runHTTP2WithoutTLS()
	s.runHTTPWithoutTLS()

	if s.opts.TLSCertFile != "" && s.opts.TLSKeyFile != "" {
		s.runHTTPWithTLS()
	}

	s.handleSignals()

	go pm.Run()

	if configmap.MultisiteEnabled {
		multisite.GetMonitor().Run()
	}
}

func (s *Server) setSpecialRoute(h *httpserver.HttpServer) {
	nfdiscutil.PreComplieRegexp()
	httpserver.PathFunc("/nnrf-disc/v1/nf-instances", "GET", nfdiscovery.NrfDiscSearchGetHandler)(h)
}

func (s *Server) runHTTP2WithoutTLS() {
	h := httpserver.InitHTTPServer(
		httpserver.Trace(s.opts.HttpTrace),
		httpserver.HostPort(s.opts.Host, strconv.Itoa(s.opts.PortHttp2WithoutTLS)),
		httpserver.HTTP2(true),
		httpserver.MaxConcurrentStreams(s.opts.Http2MaxStreamNum),
		httpserver.ReadTimeout(constvalue.HTTP_SERVER_READ_TIMEOUT),
		httpserver.WriteTimeout(constvalue.HTTP_SERVER_WRITE_TIMEOUT),

		httpserver.PathFunc("/health-check", "GET", nfdiscovery.HealthCheckHandler),
		httpserver.PathFunc("/readiness", "GET", probe.ReadinessProbe_Handler),
	)

	s.setSpecialRoute(h)
	h.Run()
}

func (s *Server) runHTTPWithTLS() {
	h := httpserver.InitHTTPServer(
		httpserver.Trace(s.opts.HttpTrace),
		httpserver.HostPort(s.opts.Host, strconv.Itoa(s.opts.PortHttpWithTLS)),
		httpserver.HTTP2(true),
		httpserver.TLSConfig(s.opts.TLSConfig),
		httpserver.MaxConcurrentStreams(s.opts.Http2MaxStreamNum),
		httpserver.ReadTimeout(constvalue.HTTP_SERVER_READ_TIMEOUT),
		httpserver.WriteTimeout(constvalue.HTTP_SERVER_WRITE_TIMEOUT),
		httpserver.PathFunc("/health-check", "GET", nfdiscovery.HealthCheckHandler),
		httpserver.PathFunc("/readiness", "GET", probe.ReadinessProbe_Handler),
	)

	s.setSpecialRoute(h)
	h.Run()
}

func (s *Server) runHTTPWithoutTLS() {
	h := httpserver.InitHTTPServer(
		httpserver.Trace(s.opts.HttpTrace),
		httpserver.HostPort(s.opts.Host, strconv.Itoa(s.opts.PortHttpWithoutTLS)),
		//httpserver.HTTP2(true),
		//httpserver.MaxConcurrentStreams(s.opts.Http2MaxStreamNum),
		httpserver.ReadTimeout(constvalue.HTTP_SERVER_READ_TIMEOUT),
		httpserver.WriteTimeout(constvalue.HTTP_SERVER_WRITE_TIMEOUT),
		httpserver.PathFunc("/health-check", "GET", nfdiscovery.HealthCheckHandler),
		httpserver.PathFunc("/readiness", "GET", probe.ReadinessProbe_Handler),
	)

	s.setSpecialRoute(h)
	h.Run()
}

// runHttp2Inner is a H2C server for inner service, such as provision
func (s *Server) runHTTP2Inner() {
	h := httpserver.InitHTTPServer(
		httpserver.Trace(s.opts.HttpTrace),
		httpserver.HostPort(s.opts.Host, strconv.Itoa(s.opts.PortHTTP2Inner)),
		httpserver.HTTP2(true),
		httpserver.MaxConcurrentStreams(s.opts.Http2MaxStreamNum),
		httpserver.ReadTimeout(constvalue.HTTP_SERVER_READ_TIMEOUT),
		httpserver.WriteTimeout(constvalue.HTTP_SERVER_WRITE_TIMEOUT),

		httpserver.PathFunc("/health-check", "GET", nfdiscovery.HealthCheckHandler),
		httpserver.PathFunc("/readiness", "GET", probe.ReadinessProbe_Handler),
	)

	s.setSpecialRoute(h)
	h.Run()
}

func registerNRFDiscMetrics() {
	pm.RegisterMetric(constvalue.NfDiscoveryRequestsTotal, "CounterVec", "Total number of discovery requests received.", constvalue.CountLabelListRequest)
	pm.RegisterMetric(constvalue.NfDiscoverySuccessTotal, "CounterVec", "Total number of successful responses sent.", constvalue.CountLabelListSuccResponse)
	pm.RegisterMetric(constvalue.NfDiscoveryFailureTotal, "CounterVec", "Total number of failed responses sent.", constvalue.CountLabelListUnSuccResponse)
	pm.RegisterHistogramVecMetric(constvalue.NfRequestDuration, "HistogramVec", "The duration to handle NF request, partitioned by NF serivce operation.", []float64{0.005, 0.01, 0.05, 0.1, 1}, []string{"operation"})
}
