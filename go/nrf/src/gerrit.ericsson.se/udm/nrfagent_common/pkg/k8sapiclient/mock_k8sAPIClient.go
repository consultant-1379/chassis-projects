package k8sapiclient

import (
	"gerrit.ericsson.se/udm/common/pkg/k8sapiproxy"
)

//GetK8sAPIProxyStub is for UT only
func GetK8sAPIProxyStub() {
	getK8sAPIProxy = func() *k8sapiproxy.K8sAPIProxy {
		return &k8sapiproxy.K8sAPIProxy{}
	}
}

//GetConfigMapStub is for UT only
func GetConfigMapStub(data []byte, e error) {
	getConfigMap = func(name string) ([]byte, error) {
		return data, e
	}
}

//PatchConfigMapStub is for UT only
func PatchConfigMapStub(e error) {
	patchConfigMap = func(name, opr, path string, value interface{}) error {
		return e
	}
}
