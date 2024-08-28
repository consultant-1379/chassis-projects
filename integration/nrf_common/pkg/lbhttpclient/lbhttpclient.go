package lbhttpclient

import (
	"fmt"
	"io"
	"sync"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
	h2cclient "gerrit.ericsson.se/udm/nrf_common/pkg/client"
)

// LBHttpClient is used to create a long-lived http connection pool
type LBHttpClient struct {
	method  string
	uriList []string
	pointer int
	header  map[string][]string
	mutex   *sync.Mutex
}

// NewLBHttpClient creates a LBHttpClient instance
func NewLBHttpClient(method string, uriList []string, header map[string][]string) *LBHttpClient {
	return &LBHttpClient{
		method:  method,
		uriList: uriList,
		pointer: -1,
		header:  header,
		mutex:   new(sync.Mutex),
	}
}

// SetURIList sets value to LBHttpClient.uriList
func (client *LBHttpClient) SetURIList(uriList []string) {
	client.mutex.Lock()
	client.uriList = uriList
	client.mutex.Unlock()
}

// GetURIList returns LBHttpClient.uriList
func (client *LBHttpClient) GetURIList() []string {
	return client.uriList
}

// Do sends request
func (client *LBHttpClient) Do(body io.Reader) (resp *httpclient.HttpRespData, err error) {
	client.mutex.Lock()

	length := len(client.uriList)
	if length <= 0 {
		client.mutex.Unlock()
		return nil, fmt.Errorf("no available server")
	}
	client.pointer++
	if client.pointer >= length {
		client.pointer = 0
	}
	uri := client.uriList[client.pointer]

	client.mutex.Unlock()

	log.Debugf("http request URI: %s", uri)

	return h2cclient.HTTPDo(client.method, uri, client.header, body)
}
