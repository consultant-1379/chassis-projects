package app

import (
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"

	"gerrit.ericsson.se/udm/common/pkg/httpserver"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/pm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/client"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_discovery/app/disc"
	"gerrit.ericsson.se/udm/nrfagent_discovery/pkg/worker"
)

var (
	//ServerStatus indicate server status
	ServerStatus = consts.ServerIsInitializing
)

type Server struct {
	sigurs1_handle bool
	sigMutex       sync.Mutex
	Terminate      chan int
}

//NewServer create server to run nrfclient agent
func NewServer() *Server {
	cm.Opts.ConfigureOptions()
	//cm.Opts.ProcessEnvVar()
	cm.Opts.ProcessOptions()
	//cm.Opts.ProcessConfigFile()
	initLog()
	cm.InitCM()

	ch := make(chan int)
	s := &Server{
		sigurs1_handle: false,
		Terminate:      ch,
	}

	pm.Init(cm.Opts.MetricsServicePort, cm.Opts.MetricsPath)

	client.InitHttpClient()

	log.Warningf("Version: %s", cm.Version)
	return s
}

func initLog() {
	log.SetServiceID("nrfagent-" + cm.Opts.WorkMode)
	log.SetNF("nrfagent")

	log.SetPodIP(os.Getenv("POD_IP"))
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)
	log.SetLevel(log.Level(cm.Opts.LogLevel))
}

//Stop server
func (s *Server) Stop() {
	ServerStatus = consts.ServerIsClosing

	log.Warningf("Exiting service ...")
	worker.StopWorkModeMonitor()
	disc.StopAgentRoleMonitor()
	//s.fsWatcher.StopFsWatcher()
}

//Run server
func (s *Server) Run() {
	s.runAsDiscAgent()

	go func() {
		pm.Run()
	}()

	ServerStatus = consts.ServerIsRunning

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGUSR1)
	go func() {
		for sg := range ch {
			switch sg {
			case syscall.SIGUSR1:
				log.Errorf("goroutine number: %v", runtime.NumGoroutine())
				const size = 64 << 13 //524288
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, true)]
				log.Errorf("%s", buf)
			default:
				log.Error("other signal received:", sg)
			}
		}
	}()
}

func (s *Server) runAsDiscAgent() {
	disc.Setup()
	s.runHTTP2WithoutTLS()
	s.runHTTPWithoutTLS()
	s.handleSignals()

	registerNRFDiscAgentMetrics()
	disc.InitiateRun()
}

func (s *Server) setCommonRoute(h *httpserver.HttpServer) {
	httpserver.PathFunc("/{agentType}/v1/health", "GET", CheckHealth)(h)
	httpserver.PathFunc("/{agentType}/v1/ready-check", "GET", CheckReadiness)(h)
	//httpserver.PathFunc("/{agentType}/v1/configurations", "GET", GetConfs)(h)
	httpserver.PathFunc("/{agentType}/v1/opts", "GET", GetOpts)(h)
	httpserver.PathFunc("/{agentType}/v1/envs", "GET", GetEnvs)(h)
}

func (s *Server) runHTTPWithoutTLS() {
	h := httpserver.InitHTTPServer(
		httpserver.Trace(cm.Opts.HTTPTrace),
		httpserver.HostPort(cm.Opts.Host, strconv.Itoa(cm.Opts.PortHTTPWithoutTLS)),
		httpserver.ReadTimeout(consts.HTTPServerReadTimeout),
		httpserver.WriteTimeout(consts.HTTPServerWriteTimeout),
		httpserver.SetRoute(),
	)
	s.setCommonRoute(h)
	h.Run()
}

func (s *Server) runHTTP2WithoutTLS() {
	h := httpserver.InitHTTPServer(
		httpserver.Trace(cm.Opts.HTTPTrace),
		httpserver.HostPort(cm.Opts.Host, strconv.Itoa(cm.Opts.PortHTTP2WithoutTLS)),
		httpserver.HTTP2(true),
		httpserver.MaxConcurrentStreams(cm.Opts.HTTP2MaxStreamNum),
		httpserver.ReadTimeout(consts.HTTPServerReadTimeout),
		httpserver.WriteTimeout(consts.HTTPServerWriteTimeout),
		httpserver.SetRoute(),
	)
	s.setCommonRoute(h)
	h.Run()
}

func (s *Server) runHTTPWithTLS() {
	h := httpserver.InitHTTPServer(
		httpserver.Trace(cm.Opts.HTTPTrace),
		httpserver.HostPort(cm.Opts.Host, strconv.Itoa(cm.Opts.PortHTTPWithTLS)),
		httpserver.HTTP2(true),
		httpserver.TLSConfig(cm.Opts.TLSConfig),
		httpserver.MaxConcurrentStreams(cm.Opts.HTTP2MaxStreamNum),
		httpserver.ReadTimeout(consts.HTTPServerReadTimeout),
		httpserver.WriteTimeout(consts.HTTPServerWriteTimeout),
		httpserver.SetRoute(),
	)
	s.setCommonRoute(h)
	h.Run()
}

func registerNRFRegAgentMetrics() {
}

func registerNRFNtfAgentMetrics() {
}

func registerNRFDiscAgentMetrics() {
	pm.RegisterMetric(consts.NfDiscoveryRequestsTotal, "Counter", "Total number of Discovery Requests NF service(s) sent to NRF Agent.", nil)
	pm.RegisterMetric(consts.NfDiscoveryResponsesTotal, "Counter", "Total number of Discovery Responses sent from NRF Agent to NF service(s).", nil)

	pm.RegisterMetric(consts.NrfDiscoveryRequestsTotal, "Counter", "Total number of Discovery Requests sent from NRF Agent to NF service(s).", nil)
	pm.RegisterMetric(consts.NrfDiscoveryResponsesTotal, "Counter", "Total number of Discovery Responses sent from NRF to NRF Agent.", nil)
	pm.RegisterMetric(consts.NrfDiscoveryResponses2xx, "Counter", "Total number of Discovery Responses with Status Code 2xx sent from NRF to NRF Agent.", nil)
	pm.RegisterMetric(consts.NrfDiscoveryResponses3xx, "Counter", "Total number of Discovery Responses with Status Code 3xx sent from NRF to NRF Agent.", nil)
	pm.RegisterMetric(consts.NrfDiscoveryResponses4xx, "Counter", "Total number of Discovery Responses with Status Code 4xx sent from NRF to NRF Agent.", nil)
	pm.RegisterMetric(consts.NrfDiscoveryResponses5xx, "Counter", "Total number of Discovery Responses with Status Code 5xx sent from NRF to NRF Agent.", nil)

	pm.RegisterHistogramVecMetric(consts.NfRequestDuration, "HistogramVec", "The duration to handle NF request, partitioned by NF serivce operation.", []float64{0.005, 0.01, 0.05, 0.1, 1}, []string{"operation"})
}
