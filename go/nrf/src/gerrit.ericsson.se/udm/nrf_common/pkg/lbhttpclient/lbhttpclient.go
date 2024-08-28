package lbhttpclient

import (
	"fmt"
	"io"
	"math/rand"
	"sync"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
	h2cclient "gerrit.ericsson.se/udm/nrf_common/pkg/client"
)

// LBHttpClient is used to create a long-lived http connection pool
type LBHttpClient struct {
	method   string
	uriList  []string
	uriCount int
	header   map[string]string
	mutex    *sync.RWMutex
}

// NewLBHttpClient creates a LBHttpClient instance
func NewLBHttpClient(method string, uriList []string, header map[string]string) *LBHttpClient {
	return &LBHttpClient{
		method:  method,
		uriList: uriList,
		header:  header,
		mutex:   new(sync.RWMutex),
	}
}

// SetURIList sets value to LBHttpClient.uriList
func (client *LBHttpClient) SetURIList(uriList []string) {
	client.mutex.Lock()
	client.uriList = uriList
	client.uriCount = len(client.uriList)
	client.mutex.Unlock()
}

// GetURIList returns LBHttpClient.uriList
func (client *LBHttpClient) GetURIList() []string {
	return client.uriList
}

// Do sends request
func (client *LBHttpClient) Do(body io.Reader) (resp *httpclient.HttpRespData, err error) {
	defer client.mutex.RUnlock()
	client.mutex.RLock()
	if client.uriCount <= 0 {
		return nil, fmt.Errorf("no available server")
	}
	uri := client.uriList[rand.Intn(client.uriCount)]
	log.Debugf("http request URI: %s", uri)
	return h2cclient.Default_h2c.HttpDo(client.method, uri, client.header, body)
}
