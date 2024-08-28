package cm

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"gerrit.ericsson.se/udm/nrfagent_common/pkg/consts"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/utils"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
)

//TLSConfigOpts defination
type TLSConfigOpts struct {
	CertFile string
	KeyFile  string
	CaFile   string
	Verify   bool
}

//ClientOpts defination
type ClientOpts struct {
	MaxIdleConnects int
	Timeout         int
	KeepAlive       bool

	TLSInsecure bool
	TLSConfig   *tls.Config `json:"-"`
}

//Options defination
type Options struct {
	ClientOpts        *ClientOpts
	CPUs              int
	NrfMgmtADDR       string
	NrfDiscPort       string
	WorkMode          string
	HTTP2MaxStreamNum uint32
	NoSigs            bool
	PprofTime         int
	FileNotifyDir     string

	TLSConfig        *tls.Config `json:"-"`
	TLSInsecure      bool
	TLSCertFile      string
	TLSKeyFile       string
	TLSRootCertFiles []string

	LogLevel            uint32
	HTTPTrace           bool
	Host                string
	PortHTTP2WithoutTLS int
	PortHTTPWithoutTLS  int
	PortHTTPWithTLS     int
	NfName              string
	Ingress_cluster     string
	Nrf_mgmt_cluster    string
	Nrf_disc_cluster    string

	configFile string

	//For PM
	MetricsPath        string
	MetricsServicePort string
}

//Version information
var (
	Version = "0.1"
	Opts    = NewOptions()
	PodIp   string
)

func (opts *Options) dumpInfo() {
	data, _ := json.Marshal(opts)
	fmt.Println(string(data))
}

//SetPodIp set pod ip
func SetPodIp(podIp string) {
	PodIp = podIp
}

//NewOptions create new option
func NewOptions() *Options {
	return &Options{
		ClientOpts:       &ClientOpts{},
		TLSRootCertFiles: make([]string, 0),
	}
}

//ConfigureOptions config options
func (opts *Options) ConfigureOptions() {
	fs := flag.NewFlagSet("nrfagent", flag.ExitOnError)

	commands := map[string]string{
		consts.CmdStartREG:  " [flags]  -- Start nrfagent as Register agent",
		consts.CmdStartNTF:  " [flags]  -- Start nrfagent as Notify agent",
		consts.CmdStartDISC: " [flags]  -- Start nrfagent as Discovery agent",
		consts.CmdTEST:      " [flags]  -- Start nrfagent in Test mode",
		consts.CmdVERSION:   "  -- Print version and exit",
		consts.CmdHELP:      "  -- Print usage and exit",
		consts.CmdSTATUS:    "  -- Show the status of nrfagent",
	}

	var (
		streams   uint
		loglevel  uint
		rootCerts string
	)

	//fs.BoolVar(&showVersion, "version", false, "")
	//fs.BoolVar(&showHelp, "help", false, "")

	//client
	fs.StringVar(&opts.configFile, "config", "/etc/nrfagent/nrfagent_conf.json", "Configuration file for NRF Agent.")
	fs.BoolVar(&opts.ClientOpts.TLSInsecure, "client_tlsinsecure", true, "Insecure TLS certificates")
	fs.BoolVar(&opts.ClientOpts.KeepAlive, "client_keepalive", true, "Client keep alive")
	fs.IntVar(&opts.ClientOpts.Timeout, "client_timeout", 5, "Client request timeout")
	fs.IntVar(&opts.ClientOpts.MaxIdleConnects, "client_maxide", httpclient.DefaultConnections, "Client max idle connect per host")

	fs.IntVar(&opts.CPUs, "cpus", runtime.NumCPU(), "Number of CPUs to use")
	fs.BoolVar(&opts.TLSInsecure, "tlsinsecure", true, "Insecure server TLS certificates")
	fs.StringVar(&opts.TLSCertFile, "tlscert", "", "certificate file.")
	fs.StringVar(&opts.TLSKeyFile, "tlskey", "", "Private key for certificate.")
	fs.StringVar(&rootCerts, "tlsrootcerts", "", "TLS root certificate files (comma separated list), e.g. a.crt,b.crt,c.crt")
	fs.StringVar(&opts.NrfMgmtADDR, "nrf_mgmt", "", "NRF management service addr, e.g. 146.11.22.222:31145")
	fs.StringVar(&opts.NrfDiscPort, "nrf_disc", "", "NRF discovery service addr, e.g. 146.11.22.222:32345")
	fs.StringVar(&opts.Host, "host", "", "host ip, if empty, use any address")
	fs.IntVar(&opts.PortHTTP2WithoutTLS, "port_h2c", 3000, "port for http/2.0 without TLS")
	fs.IntVar(&opts.PortHTTPWithoutTLS, "port_h1", 3001, "port for http/1.1 without TLS")
	fs.IntVar(&opts.PortHTTPWithTLS, "port_h2", 3002, "port for http/1.1 and http/2.0 with TLS")
	fs.BoolVar(&opts.HTTPTrace, "httptrace", false, "Enable HTTP trace")
	fs.BoolVar(&opts.NoSigs, "nosigs", false, "Disable signal funcation")
	fs.IntVar(&opts.PprofTime, "pproftime", 30, "Pprof time interval")
	fs.UintVar(&streams, "maxstreams", 2000, "MaxStreams number for http2 stream cocurrent")
	fs.UintVar(&loglevel, "loglevel", 3, "loglevel, default warning; 5:debug, 4:info, 3:warning, 2:error, 1:fatal")
	fs.StringVar(&opts.FileNotifyDir, "filenotify", "", "File notify directory")
	fs.StringVar(&opts.NfName, "nfname", "nfname", "The name of NF which the NRF client working for")

	fs.StringVar(&opts.Ingress_cluster, "ingress_cluster", "", "The CallBack of Subscription")
	fs.StringVar(&opts.Nrf_mgmt_cluster, "nrf_mgmt_cluster", "", "The CallBack of Subscription")
	fs.StringVar(&opts.Nrf_disc_cluster, "nrf_disc_cluster", "", "The CallBack of Subscription")
	fs.StringVar(&opts.MetricsPath, "metrics_path", "/metrics", "Endpoint for the metrics of the http server")
	fs.StringVar(&opts.MetricsServicePort, "metrics_port", "3003", "Metric service port of the http server")

	fs.Usage = func() {
		fmt.Println("Usage: nrfClient <command> [command flags]\n\ncommands:")
		for name, option := range commands {
			fmt.Printf("%s%s\n", name, option)
		}
		fmt.Printf("\n\ncommand flags:\n")

		fs.PrintDefaults()
	}

	if len(os.Args) <= 1 {
		fmt.Println("Please input correct command")
		fs.Usage()
		os.Exit(1)
	}
	cmd := os.Args[1]
	switch {
	case cmd == consts.CmdStartREG:
		opts.WorkMode = consts.AppWorkmodeREG
		fmt.Printf("start nrf client as %s agent...\n", consts.AppWorkmodeREG)
	case cmd == consts.CmdStartNTF:
		opts.WorkMode = consts.AppWorkmodeNTF
		fmt.Printf("start nrf client as %s agent...\n", consts.AppWorkmodeNTF)
	case cmd == consts.CmdStartDISC:
		opts.WorkMode = consts.AppWorkmodeDISC
		fmt.Printf("start nrf client as %s agent...\n", consts.AppWorkmodeDISC)
	case cmd == consts.CmdSTATUS:
		fmt.Printf("The status function is under developing.\n")
		os.Exit(0)
	case cmd == consts.CmdTEST:
		fmt.Printf("The test function is under developing.\n")
		os.Exit(0)
	case cmd == consts.CmdHELP:
		fs.Usage()
		os.Exit(0)
	case cmd == consts.CmdVERSION:
		fmt.Println("Version: ", Version)
		os.Exit(0)
	default:
		fmt.Println("Unknow command: ", cmd)
		fs.Usage()
		os.Exit(1)
	}

	if err := fs.Parse(os.Args[2:]); err != nil {
		fmt.Println("Options parse error ", err)
		fs.Usage()
		os.Exit(1)
	}

	opts.HTTP2MaxStreamNum = uint32(streams)
	opts.LogLevel = uint32(loglevel)
	if rootCerts != "" {
		opts.TLSRootCertFiles = append(opts.TLSRootCertFiles, strings.Split(rootCerts, ",")...)
	}
}

//func (opts *Options) ProcessEnvVar() {
//	//redisMasterServiceName := os.Getenv("REDIS_MASTER_SERVICE_NAME")
//	//fmt.Println(redisMasterServiceName)
//}

//func (opts *Options) ProcessConfigFile() {

//}

//ProcessOptions precess option
func (opts *Options) ProcessOptions() {
	opts.dumpInfo()

	runtime.GOMAXPROCS(opts.CPUs)

	var err error
	if opts.TLSConfig, err = utils.GenTlsConfig(opts.TLSInsecure, opts.TLSCertFile, opts.TLSKeyFile, opts.TLSRootCertFiles); err != nil {
		panic(err)
	}

	opts.ClientOpts.TLSConfig = opts.TLSConfig.Clone()
	opts.ClientOpts.TLSConfig.InsecureSkipVerify = opts.ClientOpts.TLSInsecure

	podIp := os.Getenv("POD_IP")
	SetPodIp(podIp)
}
