package client

import (
	"crypto/tls"
	"io"
	"strings"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"os"
	"fmt"
)

var (
	isServiceMeshHandleTLS = false
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
	
	serviceMeshHandleTLS := os.Getenv("SERVICEMESH_HANDLE_TLS")
	if serviceMeshHandleTLS == "true" {
		isServiceMeshHandleTLS = true
	} else {
		isServiceMeshHandleTLS = false
	}
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

//HTTPDo to handle http request
func HTTPDo(method, URL string, hdr httpclient.NHeader, body io.Reader)(resp *httpclient.HttpRespData, err error){
	if isServiceMeshHandleTLS {
		u, e := httpclient.HTTPSToHTTP(URL)
		if e != nil {
			return nil, fmt.Errorf("invalid url")
		}
		return Default_h2c.HttpDo(method, u, hdr, body)
	}
	if strings.HasPrefix(URL, "https") {
		return Default_https.HttpDo(method, URL, hdr, body)
	} else if strings.HasPrefix(URL, "http") {
		return Default_h2c.HttpDo(method, URL, hdr, body)
	}

	return nil, fmt.Errorf("Unkown URL: %v", URL)
}
//HTTPDoNoRedirect not support redirect, block http package redirect function
func HTTPDoNoRedirect(method, URL string, hdr httpclient.NHeader, body io.Reader)(resp *httpclient.HttpRespData, err error) {
	if isServiceMeshHandleTLS {
		u, e := httpclient.HTTPSToHTTP(URL)
		if e != nil {
			return nil, fmt.Errorf("invalid url")
		}
		return NoRedirect_h2c.HttpDo(method, u, hdr, body)
	}

	if strings.HasPrefix(URL, "https") {
		return NoRedirect_https.HttpDo(method, URL, hdr, body)
	} else if strings.HasPrefix(URL, "http") {
		return NoRedirect_h2c.HttpDo(method, URL, hdr, body)
	}

	return nil, fmt.Errorf("Unkown URL: %v", URL)
}

//HTTPDoRedirect to process request when request has subrequest about redirect
func HTTPDoRedirect(method, URL string, hdr httpclient.NHeader, body []byte, selfURL []string)(resp *httpclient.HttpRespData, realStatusCode int, err error) {
	if isServiceMeshHandleTLS {
		u, e := httpclient.HTTPSToHTTP(URL)
		if e != nil {
			return nil, 0, fmt.Errorf("invalid url")
		}
		return NoRedirect_h2c.HTTPDoProcRedirect(method, u, hdr, body, selfURL, isServiceMeshHandleTLS)
	}

	if strings.HasPrefix(URL, "https") {
		return NoRedirect_https.HTTPDoProcRedirect(method, URL, hdr, body, selfURL, isServiceMeshHandleTLS)
	} else if strings.HasPrefix(URL, "http") {
		return NoRedirect_h2c.HTTPDoProcRedirect(method, URL, hdr, body, selfURL, isServiceMeshHandleTLS)
	}

	return nil, 0, fmt.Errorf("Unkown URL: %v", URL)
}

//HTTPDoOmitOutput do http request
func HTTPDoOmitOutput(method, url string, hdr httpclient.NHeader, body io.Reader) {
	if isServiceMeshHandleTLS {
		u, e := httpclient.HTTPSToHTTP(url)
		if e != nil {
			log.Warnf("invalid url: %v", url)
		}
		resp, err := Default_h2c.HttpDo(method, u, hdr, body)
		if err != nil {
			log.Warnf("err %v", err)
		}
		if resp != nil {
			log.Debugf("resp: %v", resp.SimpleString())
		}
		return
	}

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

//HTTPDoJsonBody to process http request with json body
func HTTPDoJsonBody(method, url string, body io.Reader) (resp *httpclient.HttpRespData, err error) {
	if isServiceMeshHandleTLS {
		u, e := httpclient.HTTPSToHTTP(url)
		if e != nil {
			return nil, fmt.Errorf("invalid url: %v", url)
		}
		return Default_h2c.HttpDoJsonBody(method, u, body)
	}

	if strings.HasPrefix(url, "https") {
		return  Default_https.HttpDoJsonBody(method, url, body)
	} else if strings.HasPrefix(url, "http") {
		return  Default_h2c.HttpDoJsonBody(method, url, body)
	}

	return nil, fmt.Errorf("invalid url: %v", url)

}
