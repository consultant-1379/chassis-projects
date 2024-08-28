package options

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/nrf_common/pkg/client"
	"gerrit.ericsson.se/udm/nrf_common/pkg/utils"
)

//TLSConfigOpts tls config info
type TLSConfigOpts struct {
	CertFile string
	KeyFile  string
	CaFile   string
	Verify   bool
}

//Options config info
type Options struct {
	ClientOpts        *client.ClientOpts
	CPUs              int
	MaxWorker         int
	MaxQueue          int
	GrpcServer        string
	Http2MaxStreamNum uint32
	NoSigs            bool
	PprofTime         int
	PodIP             string
	PodName           string
	ContainerName     string
	HostName          string
	WorkMode          string

	TLSConfig        *tls.Config
	TLSInsecure      bool
	TLSCertFile      string
	TLSKeyFile       string
	TLSRootCertFiles []string

	LogLevel  uint32
	HttpTrace bool

	Host                string
	PortHttp2WithoutTLS int
	PortHttpWithoutTLS  int
	PortHttpWithTLS     int
	PortHTTP2Inner      int

	//For PM
	MetricsPath        string
	MetricsServicePort string

	// For log level control
	LogConfigFile string
}

var (
	Version     = "???" //version info
	optInstance *Options
)

func (opts *Options) dumpInfo() {
	data, _ := json.Marshal(opts)
	fmt.Println(string(data))
}

//Instance get the Options instance
func Instance() *Options {
	if optInstance == nil {
		optInstance = NewOptions()
	}

	return optInstance
}

//NewOptions create config options
func NewOptions() *Options {
	opt := &Options{
		ClientOpts:       &client.ClientOpts{},
		TLSRootCertFiles: make([]string, 0),
	}
	opt.ConfigureOptions()
	opt.ProcessEnvVar()
	opt.ProcessOptions()
	return opt
}

//ConfigureOptions config options
func (opts *Options) ConfigureOptions() {
	fs := flag.NewFlagSet("stone", flag.ExitOnError)

	var (
		showVersion bool
		showHelp    bool
		streams     uint
		loglevel    uint
		rootCerts   string
	)

	fs.BoolVar(&showVersion, "version", false, "Print version and exit")
	fs.BoolVar(&showHelp, "help", false, "Print version and exit")

	//client
	fs.BoolVar(&opts.ClientOpts.TLSInsecure, "client_tlsinsecure", true, "Insecure TLS certificates")
	fs.BoolVar(&opts.ClientOpts.KeepAlive, "client_keepalive", true, "Client keep alive")
	fs.IntVar(&opts.ClientOpts.Timeout, "client_timeout", 10, "Client request timeout")
	fs.IntVar(&opts.ClientOpts.MaxIdleConnects, "client_maxide", httpclient.DefaultConnections, "Client max idle connect per host")

	fs.IntVar(&opts.CPUs, "cpus", runtime.NumCPU(), "Number of CPUs to use")
	fs.IntVar(&opts.MaxWorker, "maxworker", 100, "Maxworker number for subscription/notification")
	fs.IntVar(&opts.MaxQueue, "maxqueueszie", 3000, "MaxQueue number for subscription/notification")
	fs.BoolVar(&opts.TLSInsecure, "tlsinsecure", true, "Insecure server TLS certificates")
	fs.StringVar(&opts.TLSCertFile, "tlscert", "", "certificate file.")
	fs.StringVar(&opts.TLSKeyFile, "tlskey", "", "Private key for certificate.")
	fs.StringVar(&rootCerts, "tlsrootcerts", "", "TLS root certificate files (comma separated list), e.g. a.crt,b.crt,c.crt")
	fs.StringVar(&opts.GrpcServer, "grpcserver", "localhost:50051", "gRpc server addr, e.g. 146.11.22.222:50051")
	fs.StringVar(&opts.Host, "host", "", "host ip, if empty, use any address")
	fs.IntVar(&opts.PortHttp2WithoutTLS, "port_h2c", 3000, "port for http/2.0 without TLS")
	fs.IntVar(&opts.PortHttpWithoutTLS, "port_h1", 3001, "port for http/1.1 without TLS")
	fs.IntVar(&opts.PortHttpWithTLS, "port_h2", 3002, "port for http/1.1 and http/2.0 with TLS")
	fs.IntVar(&opts.PortHTTP2Inner, "port_h2c_inner", 3004, "port for inner service")
	fs.BoolVar(&opts.HttpTrace, "httptrace", false, "Enable HTTP trace")
	fs.BoolVar(&opts.NoSigs, "nosigs", false, "Disable signal funcation")
	fs.IntVar(&opts.PprofTime, "pproftime", 30, "Pprof time interval")
	fs.UintVar(&streams, "maxstreams", 2000, "MaxStreams number for http2 stream cocurrent")
	fs.UintVar(&loglevel, "loglevel", 4, "loglevel, default warning; 5:debug, 4:info, 3:warning, 2:error, 1:fatal")
	fs.StringVar(&opts.MetricsPath, "metrics_path", "/metrics", "Endpoint for the metrics of the http server")
	fs.StringVar(&opts.MetricsServicePort, "metrics_port", "3003", "Metric service port of the http server")

	fs.Usage = func() {
		fmt.Println("Usage: stone <command> [command flags]\n\ncommand:")

		fmt.Printf("\n\ncommand flags:\n")

		fs.PrintDefaults()
	}

	if len(os.Args) <= 1 {
		fmt.Println("Please input correct command")
		fs.Usage()
		os.Exit(1)
	}

	if err := fs.Parse(os.Args[2:]); err != nil {
		fmt.Println("Options parse error ", err)
		fs.Usage()
		os.Exit(1)
	}

	if showVersion {
		fmt.Println("Version: ", Version)
		os.Exit(0)
	}

	if showHelp {
		fs.Usage()
		os.Exit(0)
	}

	opts.Http2MaxStreamNum = uint32(streams)
	opts.LogLevel = uint32(loglevel)
	if rootCerts != "" {
		opts.TLSRootCertFiles = append(opts.TLSRootCertFiles, strings.Split(rootCerts, ",")...)
	}
}

//ProcessEnvVar process env var
func (opts *Options) ProcessEnvVar() {
	opts.PodIP = os.Getenv("POD_IP")
	opts.PodName = os.Getenv("POD_NAME")
	opts.ContainerName = os.Getenv("CONTAINER_NAME")
	opts.HostName = os.Getenv("HOSTNAME")
	opts.WorkMode = os.Getenv("WORK_MODE")
	opts.LogConfigFile = os.Getenv("LOG_CONFIG_FILE")
}

//ProcessOptions process options config
func (opts *Options) ProcessOptions() {
	opts.dumpInfo()

	runtime.GOMAXPROCS(opts.CPUs)

	var err error
	if opts.TLSConfig, err = utils.GenTlsConfig(opts.TLSInsecure, opts.TLSCertFile, opts.TLSKeyFile, opts.TLSRootCertFiles); err != nil {
		panic(err)
	}

	opts.ClientOpts.TLSConfig = opts.TLSConfig.Clone()
	opts.ClientOpts.TLSConfig.InsecureSkipVerify = opts.ClientOpts.TLSInsecure

}
