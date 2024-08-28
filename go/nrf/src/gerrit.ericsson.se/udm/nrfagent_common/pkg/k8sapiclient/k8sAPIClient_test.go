package k8sapiclient

import (
	"os"
	"testing"
)

var configmapContent = []byte(`
{"kind":"ConfigMap","apiVersion":"v1",
"metadata":{
"name":"eric-nrfagent-storage",
"namespace":"default",
"selfLink":"/api/v1/namespaces/default/configmaps/eric-nrfagent-storage",
"uid":"3907fc52-af39-11e8-9eb9-005056076b54",
"resourceVersion":"1227514",
"creationTimestamp":"2018-09-03T05:21:33Z"
},"data":{"nfInfoList":"[
{\"nfInstanceId\":\"nf1\",\"nfType\":\"AUSF\",\"nfFqdn\":\"ausf.se\"},
{\"nfInstanceId\":\"nf2\",\"nfType\":\"AUSF\",\"nfFqdn\":\"ausf.se\"},
{\"nfInstanceId\":\"nf3\",\"nfType\":\"UDM\",\"nfFqdn\":\"udm.se\"},
{\"nfInstanceId\":\"nf4\",\"nfType\":\"AUSF\",\"nfFqdn\":\"ausf.se\"},
{\"nfInstanceId\":\"nf5\",\"nfType\":\"UDR\",\"nfFqdn\":\"udr.se\"}]"}}`)

func TestK8sClientPreCheck(t *testing.T) {
	GetK8sAPIProxyStub()

	t.Run("TestK8sClientPreCheck01", func(t *testing.T) {
		err := k8sClientPreCheck()
		if err == nil {
			t.Errorf("TestGetConfigMapData01 failed")
		}
	})

	t.Run("TestK8sClientPreCheck02", func(t *testing.T) {
		os.Setenv("KUBERNETES_SERVICE_HOST", "127.0.0.1")
		os.Setenv("KUBERNETES_PORT_443_TCP_PORT", "443")
		os.Setenv("POD_NAMESPACE", "default")

		err := k8sClientPreCheck()
		if err != nil {
			t.Errorf("TestK8sClientPreCheck02 failed")
		}
		if k8sServiceURI != "https://127.0.0.1:443/api/v1/namespaces/default" {
			t.Errorf("Check k8sServiceURI failed")
		}
	})
}

func TestGetConfigMapData(t *testing.T) {
	os.Setenv("KUBERNETES_SERVICE_HOST", "127.0.0.1")
	os.Setenv("KUBERNETES_PORT_443_TCP_PORT", "443")
	os.Setenv("POD_NAMESPACE", "default")

	GetConfigMapStub(configmapContent, nil)

	t.Run("TestGetConfigMapData01", func(t *testing.T) {
		_, err := GetConfigMapData("", "nfInfoList")
		if err == nil {
			t.Errorf("TestGetConfigMapData01 failed")
		}
	})
	t.Run("TestGetConfigMapData02", func(t *testing.T) {
		_, err := GetConfigMapData("nrfagent-configmap-storage", "nfInfoList")
		if err != nil {
			t.Errorf("TestGetConfigMapData02 failed")
		}
	})
}

func TestSetConfigMapData(t *testing.T) {
	os.Setenv("KUBERNETES_SERVICE_HOST", "127.0.0.1")
	os.Setenv("KUBERNETES_PORT_443_TCP_PORT", "443")
	os.Setenv("POD_NAMESPACE", "default")

	PatchConfigMapStub(nil)
	var dataTest = []byte(`test message`)

	t.Run("TestSetConfigMapData01", func(t *testing.T) {
		err := SetConfigMapData("", "nfInfoList", dataTest)
		if err == nil {
			t.Errorf("TestGetConfigMapData01 failed")
		}
	})
	t.Run("TestSetConfigMapData02", func(t *testing.T) {
		err := SetConfigMapData("nrfagent-configmap-storage", "nfInfoList", dataTest)
		if err != nil {
			t.Errorf("TestGetConfigMapData02 failed")
		}
	})
}
