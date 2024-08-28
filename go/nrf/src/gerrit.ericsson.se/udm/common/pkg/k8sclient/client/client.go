/*
Interworking with K8s API server.
*/

package cm

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	//"regexp"
	//"strings"
	//"time"

	log "gerrit.ericsson.se/udm/common/pkg/log"
	"golang.org/x/net/http2"
)

const (
	caFile    string = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	tokenFile string = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	//routeRulePrefix string = "ausfrr-"
	//routeRuleSuffix string = "-routing"
)

// K8sAPIClient is struct of k8sAPI client
type K8sAPIClient struct {
	IstioNS        string
	K8sHost        string
	K8sPort        string
	apiRoot        string
	PodNameSpace   string
	isIstio        bool
	isInit         bool
	isSideInjected bool
	token          []byte
	client         *http.Client
}

var k8sInstance *K8sAPIClient

// GetK8sAPIClient is to get object of k8s client instance
func GetK8sAPIClient() *K8sAPIClient {
	if k8sInstance == nil {
		k8sInstance = &K8sAPIClient{}
		k8sInstance.Init()
	}
	return k8sInstance
}

//Init is to init k8s client
func (apiClient *K8sAPIClient) Init() {
	if apiClient.isInit {
		log.Debug("init had done. do nothing...")
		return
	}
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
	apiClient.client = &http.Client{Transport: tr}
	apiClient.token, err = ioutil.ReadFile(tokenFile)
	if err != nil {
		log.Error("Read token File err:", err)
		return
	}
	apiClient.IstioNS = os.Getenv("ISTIO_NAMESPACE")
	if len(apiClient.IstioNS) == 0 {
		apiClient.isIstio = false
	} else {
		apiClient.isIstio = true
	}
	sideInjected := os.Getenv("SIDECAR_INJECTED")
	if sideInjected == "true" {
		apiClient.isSideInjected = true
	} else {
		apiClient.isSideInjected = false
	}

	apiClient.K8sHost = os.Getenv("KUBERNETES_SERVICE_HOST")
	apiClient.K8sPort = os.Getenv("KUBERNETES_PORT_443_TCP_PORT")
	apiClient.PodNameSpace = os.Getenv("POD_NAMESPACE")

	if len(apiClient.K8sHost) == 0 || len(apiClient.K8sPort) == 0 ||
		len(apiClient.PodNameSpace) == 0 {
		log.Error("env issue:host(", apiClient.K8sHost, "), port(", apiClient.K8sPort, "),  namespace(", apiClient.PodNameSpace, ")")
		apiClient.isInit = false
		return
	}

	apiClient.apiRoot = "https://" + apiClient.K8sHost + ":" + apiClient.K8sPort
	apiClient.isInit = true

}

/*
function: sendK8sAPIRequest
param:  opr       string              http operation "POST", "GET" etc.
        url       string              api url "/api/v1/namespaces/..."
        jsonBuf   *bytes.Buffer       json content of http body
return:           []byte              http server response body
                  error               nil if successful
spec: send configuration request by k8s API.
*/
func (apiClient *K8sAPIClient) SendK8sAPIRequest(opr, url string, jsonBuf *bytes.Buffer) ([]byte, error) {
	fullURL := apiClient.apiRoot + url
	log.Debug("sendK8sAPIRequest enter. opr:", opr, " url:", fullURL)
	var req *http.Request
	var err error
	if jsonBuf != nil {
		log.Debug("jsonbody:", jsonBuf.String())
		req, err = http.NewRequest(opr, fullURL, jsonBuf)
	} else {
		req, err = http.NewRequest(opr, fullURL, nil)
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
	req.Header.Set("Authorization", "Bearer "+string(apiClient.token))
	response, err := apiClient.client.Do(req)
	if err != nil {
		log.Error("http client err:", err)
		return nil, err
	}
	defer func() {
		if cerr := response.Body.Close(); cerr != nil {
			log.Error("Close response body error: ", cerr)
		}
	}()
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("read response body error: ", err)
	}
	log.Debug("Response from kubernets API server:")
	log.Debug(string(buf))
	if response.StatusCode == 409 {
		log.Warnf("The resource to create already exists!")
		return buf, nil
	} else if response.StatusCode < 200 || response.StatusCode > 299 {
		logstr := "k8s API response error code:" + fmt.Sprintf("%d", response.StatusCode)
		log.Error(logstr)
		return buf, errors.New(logstr)
	}
	return buf, nil
}

// CreateServiceEntry is to create service entry
func (apiClient *K8sAPIClient) CreateServiceEntry(name string, host string, port string) error {
	serviceEntryFormat := `{
							  "apiVersion": "networking.istio.io/v1alpha3",
							  "kind": "ServiceEntry",
							  "metadata": { "name": "%s" },
							  "spec":
							   { "hosts": [ "%s" ],
							     "ports":
							      [ { "number": "%s",
							          "name": "http2",
							          "protocol": "HTTP2"
							        }
							      ]
							
							   }
							 }`
	serviceEntryURI := "/apis/networking.istio.io/v1alpha3/namespaces/" + apiClient.PodNameSpace + "/serviceentries"
	var jsonBuf bytes.Buffer
	jsonBuf.WriteString(fmt.Sprintf(serviceEntryFormat, name, host, port))
	_, err := apiClient.SendK8sAPIRequest("POST", serviceEntryURI, &jsonBuf)
	return err
}

// DeleteServiceEntry is to delete service entry
func (apiClient *K8sAPIClient) DeleteServiceEntry(servienEntryName string) error {
	serviceEntryURI := "/apis/networking.istio.io/v1alpha3/namespaces/" + apiClient.PodNameSpace + "/serviceentries"
	_, err := apiClient.SendK8sAPIRequest("DELETE", serviceEntryURI+"/"+
		servienEntryName, nil)
	return err
}

// CreateEgressRule is to create egress rule
func (apiClient *K8sAPIClient) CreateEgressRule(name string, host string, port string) error {
	if !apiClient.isIstio || !apiClient.isInit || !apiClient.isSideInjected {
		return nil
	}
	tcpEgressRuleFormat := `{
							"apiVersion": "config.istio.io/v1alpha2",
							"kind": "EgressRule",
							"metadata": { 
							    "name": "%s" 
							},
							"spec": {
							  "destination": {
							        "service": "%s"
							  },
							  "ports": [
							    {
							      "port": %s,
							      "protocol": "tcp"
								}
						      ]
							}
						   }`
	egreesRuleURI := "/apis/config.istio.io/v1alpha2/namespaces/" + apiClient.PodNameSpace + "/egressrules"
	hostIP := host
	if net.ParseIP(host) == nil {
		ips, err := net.LookupHost(host)
		if err != nil {
			log.Debugf("Create EgressRule failed as %s:", err.Error())
			return err
		}
		hostIP = ips[0]
	}

	var jsonBuf bytes.Buffer
	jsonBuf.WriteString(fmt.Sprintf(tcpEgressRuleFormat, name, hostIP, port))
	_, err := apiClient.SendK8sAPIRequest("POST", egreesRuleURI, &jsonBuf)
	return err
}

// DeleteEgressRule is to delete service entry
func (apiClient *K8sAPIClient) DeleteEgressRule(egressRuleName string) error {
	if !apiClient.isIstio || !apiClient.isInit || !apiClient.isSideInjected {
		return nil
	}

	egreesRuleURI := "/apis/config.istio.io/v1alpha2/namespaces/" + apiClient.PodNameSpace + "/egressrules"
	_, err := apiClient.SendK8sAPIRequest("DELETE", egreesRuleURI+"/"+
		egressRuleName, nil)
	return err
}

// ModifyConfigMap is to patch Internal Config Map
func (apiClient *K8sAPIClient) ModifyConfigMap(configMapName string, jsonPatchbody []byte) error {
	if !apiClient.isInit {
		return nil
	}
	configMapURI := "/api/v1/namespaces/" + apiClient.PodNameSpace + "/configmaps/" + configMapName
	var jsonBuf bytes.Buffer
	jsonBuf.Write(jsonPatchbody)
	_, err := apiClient.SendK8sAPIRequest("PATCH", configMapURI, &jsonBuf)
	return err
}
