package main

import (
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMapToString(t *testing.T) {
	m := make(map[string]string)
	if mapToString(m) != "" {
		t.Fail()
	}
	m["key"] = "value"
	m["key1"] = "value1"

	if mapToString(m) != "key=value,key1=value1" && mapToString(m) != "key1=value1,key=value" {
		t.Log(m)
		t.Fail()
	}
}
func h2cServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprint(writer, "hello")
	})
	h2s := &http2.Server{}
	ts := httptest.NewServer(h2c.NewHandler(mux, h2s))
	return ts
}

func TestH2c(t *testing.T) {
	ts := h2cServer()
	defer ts.Close()
	if resp, err := httpClient.Get(ts.URL); err != nil {
		t.Fail()
	} else if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fail()
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			t.Log(string(body))
		}
	} else {
		t.Fail()
	}
}

const namespace = "nef"
const test3gppSrvName = "nnef-bdt"
const testK8SrvName = "eric-nef"

func TestMonitorAndReport(t *testing.T) {
	// Create the fake client
	pod1 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: "pod1"},
		Status: v1.PodStatus{
			Conditions: []v1.PodCondition{{Type: "Ready", Status: "True"}},
		},
	}
	pod2 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: "pod2", Labels: map[string]string{"name": "eric-nef-pod", "type": "pod"}},
		Status: v1.PodStatus{
			Conditions: []v1.PodCondition{{Type: "Ready", Status: "True"}},
		},
	}
	srv1 := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: "eric-nef"},
		Spec:       v1.ServiceSpec{Selector: map[string]string{"name": "eric-nef-pod", "type": "pod"}},
	}
	client := fake.NewSimpleClientset(pod1, pod2, srv1)

	// Create test server
	ts := newTestServer(t)
	defer ts.Close()

	// fill the config used for the test
	config = &Config{
		namespace:         "nef",
		prefix:            "eric",
		nrfAgentUri:       ts.URL,
		nfCMProfileName:   "ericsson-nef",
		nfProfilePath:     "ericsson-nef:nef",
		heartbeatInterval: 10,
		retryInterval:     2,
		retryTimes:        1,
	}

	if err := monitorAndReportSrvOnce(testK8SrvName, test3gppSrvName, client); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestMonitorAndReport_Neg1(t *testing.T) {
	// Create the fake client
	pod1 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: "pod1"},
		Status: v1.PodStatus{
			Conditions: []v1.PodCondition{{Type: "Ready", Status: "False"}},
		},
	}
	pod2 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: "pod2", Labels: map[string]string{"name": "eric-nef-pod", "type": "pod"}},
		Status: v1.PodStatus{
			Conditions: []v1.PodCondition{{Type: "Ready", Status: "False"}},
		},
	}
	srv1 := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: "eric-nef"},
		Spec:       v1.ServiceSpec{Selector: map[string]string{"name": "eric-nef-pod", "type": "pod"}},
	}
	ts := newTestServer(t)
	defer ts.Close()
	config = &Config{
		namespace:         "nef",
		prefix:            "eric",
		nrfAgentUri:       ts.URL,
		nfCMProfileName:   "ericsson-nef",
		nfProfilePath:     "ericsson-nef:nef",
		heartbeatInterval: 10,
		retryInterval:     2,
		retryTimes:        1,
	}
	client := fake.NewSimpleClientset(pod1, pod2, srv1)
	if err := monitorAndReportSrvOnce(testK8SrvName, test3gppSrvName, client); err != nil && err.Error() != "service "+testK8SrvName+" is not healthy" {
		t.Error(err)
		t.FailNow()
	}
}

func TestMonitorAndReport_Neg2(t *testing.T) {
	// Create the fake client
	pod1 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: "pod1"},
		Status: v1.PodStatus{
			Conditions: []v1.PodCondition{{Type: "Ready", Status: "True"}},
		},
	}
	pod2 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: "pod2", Labels: map[string]string{"name": "eric-nef-pod", "type": "pod"}},
		Status: v1.PodStatus{
			Conditions: []v1.PodCondition{{Type: "Ready", Status: "True"}},
		},
	}
	srv1 := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: "eric-nef"},
		Spec:       v1.ServiceSpec{Selector: map[string]string{}},
	}
	client := fake.NewSimpleClientset(pod1, pod2, srv1)
	ts := newTestServer(t)
	defer ts.Close()
	config = &Config{
		namespace:         "nef",
		prefix:            "eric",
		nrfAgentUri:       ts.URL,
		nfCMProfileName:   "ericsson-nef",
		nfProfilePath:     "ericsson-nef:nef",
		heartbeatInterval: 10,
		retryInterval:     2,
		retryTimes:        1,
	}
	if err := monitorAndReportSrvOnce(testK8SrvName, test3gppSrvName, client); err != nil && err.Error() != "spec selector of the service is empty" {
		t.Error(err)
	}
}

func TestMonitorAndReport_Neg3(t *testing.T) {
	// Create the fake client
	pod1 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: "pod1"},
		Status: v1.PodStatus{
			Conditions: []v1.PodCondition{{Type: "Ready", Status: "True"}},
		},
	}
	pod2 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: "pod2", Labels: map[string]string{"name": "eric-nef-pod", "type": "pod"}},
		Status: v1.PodStatus{
			Conditions: []v1.PodCondition{{Type: "Ready", Status: "True"}},
		},
	}
	srv1 := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: "eric-nef"},
		Spec:       v1.ServiceSpec{Selector: map[string]string{"name": "eric-nef-pod", "type": "pod"}},
	}
	client := fake.NewSimpleClientset(pod1, pod2, srv1)

	// Create test server
	ts := newTestServer(t)
	defer ts.Close()

	// fill the config used for the test
	config = &Config{
		namespace:         "nef",
		prefix:            "eric",
		nrfAgentUri:       "http://127.0.0.1",
		nfCMProfileName:   "ericsson-nef",
		nfProfilePath:     "ericsson-nef:nef",
		heartbeatInterval: 10,
		retryInterval:     2,
		retryTimes:        1,
	}

	if err := monitorAndReportSrvOnce(testK8SrvName, test3gppSrvName, client); err == nil {
		t.FailNow()
	}
}

func TestMonitorAndReport_Neg4(t *testing.T) {
	// Create the fake client
	pod1 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: "pod1"},
		Status: v1.PodStatus{
			Conditions: []v1.PodCondition{{Type: "Ready", Status: "True"}},
		},
	}
	pod2 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: "pod2", Labels: map[string]string{"name": "eric-nef-pod", "type": "pod"}},
		Status: v1.PodStatus{
			Conditions: []v1.PodCondition{{Type: "Ready", Status: "True"}},
		},
	}
	srv1 := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace + "hi", Name: "eric-nef"},
		Spec:       v1.ServiceSpec{Selector: map[string]string{"name": "eric-nef-pod", "type": "pod"}},
	}
	client := fake.NewSimpleClientset(pod1, pod2, srv1)

	// Create test server
	ts := newTestServer(t)
	defer ts.Close()

	// fill the config used for the test
	config = &Config{
		namespace:         "nef",
		prefix:            "eric",
		nrfAgentUri:       ts.URL,
		nfCMProfileName:   "ericsson-nef",
		nfProfilePath:     "ericsson-nef:nef",
		heartbeatInterval: 10,
		retryInterval:     2,
		retryTimes:        1,
	}

	if err := monitorAndReportSrvOnce(testK8SrvName, test3gppSrvName, client); err != nil && err.Error() != "services \"eric-nef\" not found" {
		t.Error(err)
	}
}

func newTestServer(t *testing.T) *httptest.Server {
	srv := http.NewServeMux()
	srv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Error(r.Method)
			t.FailNow()
		}
		if r.URL.EscapedPath() != "/nrf-register-agent/v1/nf-status/NEF/"+test3gppSrvName {
			t.Error(r.URL.EscapedPath())
			t.FailNow()
		}
		w.WriteHeader(http.StatusOK)
	})
	return httptest.NewServer(h2c.NewHandler(srv, &http2.Server{}))
}

func TestNewKubeClient_Neg(t *testing.T) {
	if _, err := NewKubeClient(); err != rest.ErrNotInCluster {
		t.Error(err)
	}
}
