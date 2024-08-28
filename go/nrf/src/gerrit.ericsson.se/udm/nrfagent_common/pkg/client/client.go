package client

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/cm"
	"gerrit.ericsson.se/udm/nrfagent_common/pkg/structs"
)

var (
	default_http    *httpclient.HttpClient
	default_https   *httpclient.HttpClient
	defaultH2c      *httpclient.HttpClient
	defaultMultiH2c []*httpclient.HttpClient

	multiH2cTokenMutex sync.Mutex
	multiH2cToken      = 0
)

func InitHttpClient() {
	default_https = httpclient.InitHttpClient(
		httpclient.Connections(cm.Opts.ClientOpts.MaxIdleConnects),
		httpclient.Timeout(time.Duration(cm.Opts.ClientOpts.Timeout)*time.Second),
		httpclient.KeepAlive(cm.Opts.ClientOpts.KeepAlive),
		httpclient.TLSConfig(cm.Opts.ClientOpts.TLSConfig),
		httpclient.HTTP2(true),
	)
	// defaultH2c is used by connecting NRF discovery or another HTTP2 server
	// defaultMultiH2c is used by connecting NRF management
	defaultH2c = httpclient.InitHttpClient(
		httpclient.ResponseTimeout(time.Duration(cm.Opts.ClientOpts.Timeout+1)*time.Second),
		httpclient.Timeout(time.Duration(cm.Opts.ClientOpts.Timeout)*time.Second),
		httpclient.KeepAlive(cm.Opts.ClientOpts.KeepAlive),
		httpclient.SupportH2c(),
	)
	for i := 0; i < cm.GetHTTP2Conns(); i++ {
		defaultMultiH2c = append(defaultMultiH2c, httpclient.InitHttpClient(
			httpclient.ResponseTimeout(time.Duration(cm.Opts.ClientOpts.Timeout+1)*time.Second),
			httpclient.Timeout(time.Duration(cm.Opts.ClientOpts.Timeout)*time.Second),
			httpclient.KeepAlive(cm.Opts.ClientOpts.KeepAlive),
			httpclient.SupportH2c(),
		))
	}

	default_http = httpclient.InitHttpClient(
		httpclient.Connections(cm.Opts.ClientOpts.MaxIdleConnects),
		httpclient.Timeout(time.Duration(cm.Opts.ClientOpts.Timeout)*time.Second),
		httpclient.KeepAlive(cm.Opts.ClientOpts.KeepAlive),
	)
}

func HttpDoJsonBody(httpv, method, url string, body io.Reader) (*httpclient.HttpRespData, error) {
	hdr := make(map[string]string)
	hdr["Content-Type"] = "application/json"
	return HTTPDo(httpv, method, url, hdr, body)
}

// HTTPDo send HTTP request
var HTTPDo = func(httpv, method, url string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
	if httpv == "h1" {
		return default_http.HttpDo(method, url, hdr, body)
	} else if strings.HasPrefix(url, "https") {
		return default_https.HttpDo(method, url, hdr, body)
	} else if strings.HasPrefix(url, "http") {
		if strings.Contains(url, structs.NrfMgmtServiceName) {
			return HTTPDoMultiConns(method, url, hdr, body)
		}
		return defaultH2c.HttpDo(method, url, hdr, body)
	}
	return nil, fmt.Errorf("invalid HTTP version (%s) or URL schema (%s)", httpv, url)
}

//HTTPDoMultiConns send HTTP request via multi HTTP Clients
func HTTPDoMultiConns(method, url string, hdr httpclient.NHeader, body io.Reader) (*httpclient.HttpRespData, error) {
	multiH2cTokenMutex.Lock()
	defer multiH2cTokenMutex.Unlock()

	var resp *httpclient.HttpRespData
	var err error
	connNum := len(defaultMultiH2c)
	for i := multiH2cToken; i < (multiH2cToken + connNum); i++ {
		token := i % connNum
		resp, err = defaultMultiH2c[token].HttpDo(method, url, hdr, body)
		if err == nil {
			multiH2cToken = (token + 1) % connNum
			break
		} else {
			log.Errorf("failed to connect %s (connection:%d)", url, token)
		}
	}
	return resp, err
}
