package k8sapiclient

import (
	"errors"
	"os"

	"github.com/buger/jsonparser"

	"gerrit.ericsson.se/udm/common/pkg/k8sapiproxy"
	"gerrit.ericsson.se/udm/common/pkg/log"
)

var (
	k8sServiceURI string
)

func k8sClientPreCheck() error {
	if k8sServiceURI == "" {
		host := os.Getenv("KUBERNETES_SERVICE_HOST")
		port := os.Getenv("KUBERNETES_PORT_443_TCP_PORT")
		namespace := os.Getenv("POD_NAMESPACE")

		if len(host) == 0 ||
			len(port) == 0 ||
			len(namespace) == 0 {
			return errors.New("failed to get KUBERNETES_SERVICE_HOST|KUBERNETES_PORT_443_TCP_PORT|POD_NAMESPACE")
		}
		k8sServiceURI = "https://" + host + ":" + port + "/api/v1/namespaces/" + namespace
		log.Infof("Kubernetes Service URI: %s", k8sServiceURI)
	}

	if getK8sAPIProxy() == nil {
		return errors.New("failed to get k8sAPIProxy instance")
	}
	return nil
}

var getK8sAPIProxy = func() *k8sapiproxy.K8sAPIProxy {
	return k8sapiproxy.GetK8sAPIProxy()
}

//GetConfigMapData get data from configmap "data:{key}"
func GetConfigMapData(name, key string) ([]byte, error) {
	if name == "" ||
		key == "" {
		return nil, errors.New("name of configmap or json key is empty")
	}

	rawData, err := getConfigMap(name)
	if err != nil {
		return nil, err
	}

	data, _, _, err := jsonparser.Get(rawData, "data", key)
	if err != nil {
		log.Errorf("failed to get jsonpath %s from rawData %s", "/data/"+key, string(rawData))
		return nil, err
	}

	//	 String value from K8S REST API includes many quotes
	//	 Ignored this quotes before delivering to APP
	var cooked []byte
	for _, c := range data {
		if c == 0x5C {
			continue
		}
		cooked = append(cooked, c)
	}
	return cooked, nil
}

//SetConfigMapData set data to configmap "data:{key}"
func SetConfigMapData(name, key string, data []byte) error {
	if name == "" ||
		key == "" {
		return errors.New("name of configmap or json key is empty")
	}

	if data == nil {
		return errors.New("invalid data (nil) for configmap")
	}

	value := string(data)
	return patchConfigMap(name, "add", "/data/"+key, &value)
}

var getConfigMap = func(name string) ([]byte, error) {
	err := k8sClientPreCheck()
	if err != nil {
		return nil, err
	}

	return k8sapiproxy.GetK8sAPIProxy().SendK8sAPIRequest("GET", k8sServiceURI+"/configmaps/"+name, nil)
}

var patchConfigMap = func(name, opr, path string, value interface{}) error {
	err := k8sClientPreCheck()
	if err != nil {
		return err
	}

	data := k8sapiproxy.K8sAPIPatchInfo{
		OperationType: opr,
		Path:          path,
		Value:         value,
	}
	return k8sapiproxy.GetK8sAPIProxy().SendK8sAPIPatchRequest(&data, k8sServiceURI+"/configmaps/"+name)
}
