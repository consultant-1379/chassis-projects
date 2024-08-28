package client

import (
	"crypto/tls"
	"io"
	"strings"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
)

var (
	Default_http     *httpclient.HttpClient
	Default_https    *httpclient.HttpClient
	Default_h2c      *httpclient.HttpClient
	NoRedirect_h2c   *httpclient.HttpClient
	NoRedirect_https *httpclient.HttpClient
)

//ClientOpts client config info
type ClientOpts struct {
	MaxIdleConnects int
	Timeout         int
	KeepAlive       bool

	TLSInsecure bool
	TLSConfig   *tls.Config
}

//InitHttpClient init http client
func InitHttpClient(opts *ClientOpts) {
	Default_https = httpclient.InitHttpClient(
		httpclient.Connections(opts.MaxIdleConnects),
		httpclient.Timeout(time.Duration(opts.Timeout)*time.Second),
		httpclient.KeepAlive(opts.KeepAlive),
		httpclient.TLSConfig(opts.TLSConfig),
		httpclient.HTTP2(true),
	)

	Default_http = httpclient.InitHttpClient(
		httpclient.Connections(opts.MaxIdleConnects),
		httpclient.Timeout(time.Duration(opts.Timeout)*time.Second),
		httpclient.KeepAlive(opts.KeepAlive),
	)

	Default_h2c = httpclient.InitHttpClient(
		httpclient.Timeout(time.Duration(opts.Timeout)*time.Second),
		httpclient.KeepAlive(opts.KeepAlive),
		httpclient.SupportH2c(),
		httpclient.ResponseTimeout(time.Duration(opts.Timeout+1)*time.Second),
	)

	NoRedirect_h2c = httpclient.InitHttpClientWithoutRedirect(
		httpclient.Timeout(time.Duration(opts.Timeout)*time.Second),
		httpclient.KeepAlive(opts.KeepAlive),
		httpclient.SupportH2c(),
		httpclient.ResponseTimeout(time.Duration(opts.Timeout+1)*time.Second),
	)

	NoRedirect_https = httpclient.InitHttpClientWithoutRedirect(
		httpclient.Connections(opts.MaxIdleConnects),
		httpclient.Timeout(time.Duration(opts.Timeout)*time.Second),
		httpclient.KeepAlive(opts.KeepAlive),
		httpclient.TLSConfig(opts.TLSConfig),
		httpclient.HTTP2(true),
	)
}

//HttpDoJsonBody http send json body message
func HttpDoJsonBody(method, url string, body io.Reader) {
	//Suitable for TCP egress rule
	if strings.HasPrefix(url, "https") {
		resp, err := Default_https.HttpDoJsonBody(method, url, body)
		if err != nil {
			log.Warnf("err %v", err)
		}
		if resp != nil {
			log.Debugf("resp: %v", resp.SimpleString())
		}
	} else if strings.HasPrefix(url, "http") {
		resp, err := Default_h2c.HttpDoJsonBody(method, url, body)
		if err != nil {
			log.Warnf("err %v", err)
		}
		if resp != nil {
			log.Debugf("resp: %v", resp.SimpleString())
		}
	}

	/* Suitable for Https egress rule
	if strings.HasPrefix(url, "https") {
		url = strings.Replace(url, "https", "http", 1)
	}
	resp, err := default_h2c.HttpDoJsonBody(method, url, body)
	if err != nil {
		log.Debugf("err %v", err)
	}
	if resp != nil {
		log.Debugf("resp: %v", resp.SimpleString())
	}
	*/
}

//HttpDo do http request
func HttpDo(method, url string, hdr httpclient.NHeader, body io.Reader) {
	if strings.HasPrefix(url, "https") {
		resp, err := Default_https.HttpDo(method, url, hdr, body)
		if err != nil {
			log.Warnf("err %v", err)
		}
		if resp != nil {
			log.Debugf("resp: %v", resp.SimpleString())
		}
	} else if strings.HasPrefix(url, "http") {
		resp, err := Default_h2c.HttpDo(method, url, hdr, body)
		if err != nil {
			log.Warnf("err %v", err)
		}
		if resp != nil {
			log.Debugf("resp: %v", resp.SimpleString())
		}
	}
}
