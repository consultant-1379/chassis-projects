package httpclient

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"strings"

	"golang.org/x/net/http2"
)

var (
	DefaultTimeout     = 30 * time.Second
	DefaultTLSConfig   = &tls.Config{InsecureSkipVerify: true}
	DefaultConnections = 2
)

type NHeader map[string]string

type HttpClient struct {
	trace  bool
	client *http.Client
	dialer *net.Dialer
}

type HttpRespData struct {
	StatusCode  int
	Protocol    string
	ContentType string
	Location    string
	Etag        string
	Header      *http.Header
	Body        []byte
}

func (h *HttpRespData) String() string {
	data, err := json.Marshal(h)
	if err == nil {
		return string(data)
	}

	return err.Error()
}

func (h *HttpRespData) SimpleString() string {
	return fmt.Sprintf("{StatusCode: %d, Protocol: %s, ContentType: %s, Etag: %s}",
		h.StatusCode, h.Protocol, h.ContentType, h.Etag)
}

func InitHttpClient(opts ...func(*HttpClient)) *HttpClient {
	a := &HttpClient{
		trace: false,
	}
	a.dialer = &net.Dialer{
		KeepAlive: 30 * time.Second,
		Timeout:   DefaultTimeout,
	}
	a.client = &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			Dial:                  a.dialer.Dial,
			ResponseHeaderTimeout: DefaultTimeout,
			TLSClientConfig:       DefaultTLSConfig,
			TLSHandshakeTimeout:   10 * time.Second,
			MaxIdleConnsPerHost:   DefaultConnections,
		},
	}
	for _, opt := range opts {
		opt(a)
	}

	return a
}

func InitHttpClientWithoutRedirect(opts ...func(*HttpClient)) *HttpClient {
	a := &HttpClient{
		trace: false,
	}
	a.dialer = &net.Dialer{
		KeepAlive: 30 * time.Second,
		Timeout:   DefaultTimeout,
	}
	a.client = &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			Dial:                  a.dialer.Dial,
			ResponseHeaderTimeout: DefaultTimeout,
			TLSClientConfig:       DefaultTLSConfig,
			TLSHandshakeTimeout:   10 * time.Second,
			MaxIdleConnsPerHost:   DefaultConnections,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse },
	}
	for _, opt := range opts {
		opt(a)
	}

	return a
}

func Connections(n int) func(*HttpClient) {
	return func(a *HttpClient) {
		tr := a.client.Transport.(*http.Transport)
		tr.MaxIdleConnsPerHost = n
	}
}

func Trace(b bool) func(*HttpClient) {
	return func(a *HttpClient) {
		a.trace = b
	}
}

func Timeout(d time.Duration) func(*HttpClient) {
	return func(a *HttpClient) {
		tr := a.client.Transport.(*http.Transport)
		tr.ResponseHeaderTimeout = d
		a.dialer.Timeout = d
		tr.Dial = a.dialer.Dial
	}
}

func ResponseTimeout(d time.Duration) func(*HttpClient) {
	return func(a *HttpClient) {
		a.client.Timeout = d
	}
}

func KeepAlive(keepalive bool) func(*HttpClient) {
	return func(a *HttpClient) {
		tr := a.client.Transport.(*http.Transport)
		tr.DisableKeepAlives = !keepalive
		if !keepalive {
			a.dialer.KeepAlive = 0
			tr.Dial = a.dialer.Dial
		}
	}
}

func SupportH2c() func(*HttpClient) {
	return func(a *HttpClient) {
		tr := &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				return a.dialer.Dial(netw, addr)
			},
		}
		tr1 := a.client.Transport.(*http.Transport)
		tr.SetT1(tr1)
		a.client.Transport = tr
	}
}

func TLSConfig(c *tls.Config) func(*HttpClient) {
	return func(a *HttpClient) {
		tr := a.client.Transport.(*http.Transport)
		tr.TLSClientConfig = c
	}
}

func HTTP2(enabled bool) func(*HttpClient) {
	return func(a *HttpClient) {
		if tr := a.client.Transport.(*http.Transport); enabled {
			if err := http2.ConfigureTransport(tr); err != nil {
				fmt.Println("Can not configure HTTP2 ", err)
				panic(err)
			}
		} else {
			tr.TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper{}
		}
	}
}

func (c *HttpClient) Do(req *http.Request) (resp *HttpRespData, err error) {
	if resp, err := c.client.Do(req); err == nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Println("Close http response error ", err)
			}
		}()
		respbody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		respContentType := resp.Header.Get("Content-Type")
		respLocation := resp.Header.Get("Location")
		respEtag := resp.Header.Get("Etag")
		if c.trace {
			fmt.Println(resp.Status, " ", resp.Proto)
		}
		return &HttpRespData{
			StatusCode:  resp.StatusCode,
			Protocol:    resp.Proto,
			ContentType: respContentType,
			Location:    respLocation,
			Etag:        respEtag,
			Header:      &resp.Header,
			Body:        respbody,
		}, nil

	} else {
		return nil, err
	}
}

func (c *HttpClient) HttpDo(method, url string, hdr NHeader, body io.Reader) (resp *HttpRespData, err error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for h, v := range hdr {
		req.Header.Set(h, v)
	}

	return c.Do(req)
}

//HttpDoProcRedirect support redirect but not use golang self redirect function
func (c *HttpClient) HttpDoProcRedirect(method, URL string, hdr NHeader, body io.Reader, selfURL []string) (resp *HttpRespData, err error) {
	var res *HttpRespData
	var e error
	apiURL := URL
	for i := 0; i < 3; i++ {
		req, err := http.NewRequest(method, apiURL, body)
		if err != nil {
			return nil, err
		}
		for h, v := range hdr {
			req.Header.Set(h, v)
		}

		res, e = c.Do(req)
		if e == nil && res.StatusCode == http.StatusTemporaryRedirect {
			l := res.Header.Get("Location")
			if l == "" {
				res.StatusCode = http.StatusBadGateway
				return res, e
			}

			for _, v := range selfURL {
				if strings.HasPrefix(l, v) {
					res.StatusCode = http.StatusBadGateway
					return res, e
				}
			}
			apiURL = l
		} else {
			return res, e
		}
	}

	return res, e
}

func (c *HttpClient) HttpDoJsonBody(method, url string, body io.Reader) (resp *HttpRespData, err error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return c.Do(req)
}
