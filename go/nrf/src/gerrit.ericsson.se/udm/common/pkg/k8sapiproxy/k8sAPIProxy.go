/*
Package k8sapiproxy for common.
*/
package k8sapiproxy

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/utils"
	"golang.org/x/net/http2"
	"io/ioutil"
	"net/http"
)

var (
	caFile    string = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	tokenFile string = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)

type K8sAPIPatchInfo struct {
	OperationType string      `json:"op"`
	Path          string      `json:"path"`
	Value         interface{} `json:"value,omitempty"`
}

type K8sAPIProxy struct {
	token  []byte
	client *http.Client
}

var k8sInstance *K8sAPIProxy

func GetK8sAPIProxy() *K8sAPIProxy {
	if k8sInstance == nil {
		k8sInstance = &K8sAPIProxy{}
		k8sInstance.Init()
	}
	return k8sInstance
}

func (apiProxy *K8sAPIProxy) Init() {
	//x509.Certificate.
	pool := x509.NewCertPool()
	caCrt, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Error("Read ca cert File err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	tr := &http2.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: pool},
		DisableCompression: true,
	}
	apiProxy.client = &http.Client{Transport: tr}

	apiProxy.token, err = ioutil.ReadFile(tokenFile)
	if err != nil {
		log.Error("Read token File err:", err)
		return
	}
}

func SetFile(ca, token string) {
	caFile = ca
	tokenFile = token
}

/*
function: SendK8sAPIRequest
param:  opr       string              http operation "POST", "GET" etc.
        url       string              http url "https://..."
        jsonBuf   *bytes.Buffer       json content of http body
return:           []byte              http server response body
                  error               nil if successful
spec: send configuration request by k8s API.
*/
func (apiProxy *K8sAPIProxy) SendK8sAPIRequest(opr, url string, jsonBuf *bytes.Buffer) ([]byte, error) {
	log.Debug("sendK8sAPIRequest enter. opr:", opr, " url:", url)

	var req *http.Request
	var err error
	if jsonBuf != nil {
		log.DebugJ(jsonBuf.String())
		req, err = http.NewRequest(opr, url, jsonBuf)
	} else {
		req, err = http.NewRequest(opr, url, nil)
	}
	if err != nil {
		log.Error("new http request err:", err)
		return nil, err
	}

	if opr == "PATCH" {
		req.Header.Set("Content-Type", "application/json-patch+json")
	} else {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+string(apiProxy.token))

	response, err := apiProxy.client.Do(req)
	if err != nil {
		log.Error("http client err:", err)
		return nil, err
	}

	defer func() {
		if cerr := response.Body.Close(); cerr != nil {
			log.Error("Close response body error: ", cerr)
		}
	}()
	buf, _ := ioutil.ReadAll(response.Body)
	log.Debug("Response from kubernets API server:")
	log.DebugJ(utils.ToPrettyJSON(buf))

	if response.StatusCode < 200 || response.StatusCode > 299 {
		logstr := "k8s API response error code:" + fmt.Sprintf("%d", response.StatusCode)
		log.Error(logstr)
		return buf, errors.New(logstr)
	}
	return buf, nil
}

/*
function: SendK8sAPIPatchRequest
param:  data      K8sAPIPatchInfo     patch data
        url       string              http url "https://..."
return:
                  error               nil if successful
spec: send patch request by k8s API.
*/
func (apiProxy *K8sAPIProxy) SendK8sAPIPatchRequest(data *K8sAPIPatchInfo, url string) error {
	var jsonBuf bytes.Buffer
	jsonBody, errJson := json.Marshal(*data)
	if errJson != nil {
		log.Error("json data marshal err:", errJson)
		return errJson
	}
	jsonBuf.WriteByte(0x5b)
	jsonBuf.Write(jsonBody)
	jsonBuf.WriteByte(0x5d)

	_, err := apiProxy.SendK8sAPIRequest("PATCH", url, &jsonBuf)
	return err
}
