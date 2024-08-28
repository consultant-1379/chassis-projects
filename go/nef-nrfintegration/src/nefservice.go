package main

import (
	"errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"math/rand"
	"time"
)

func NewKubeClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func hasReadyPod(pods []v1.Pod) bool {
	for _, pod := range pods {
		for _, podCondition := range pod.Status.Conditions {
			if podCondition.Type == "Ready" && podCondition.Status == "True" {
				logger.Debug("pod ", pod.Name, " is healthy")
				return true
			}
		}
	}
	return false
}

func mapToString(m map[string]string) string {
	var val string
	for k, v := range m {
		if val == "" {
			val += k + "=" + v
		} else {
			val += "," + k + "=" + v
		}
	}
	return val
}

func MonitorAndReportSrvWithRetry(k8sSvcName, serviceName3gpp string, client kubernetes.Interface) {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) //Adding random numbers to reduce pressure on NRFAgent
	for i := uint64(0); i == 0 && i < config.retryTimes; i++ {
		if err := monitorAndReportSrvOnce(k8sSvcName, serviceName3gpp, client); err != nil {
			logger.Error(err.Error())
			time.Sleep(time.Duration(config.retryTimes) * time.Second)
		}
	}
}

func monitorAndReportSrvOnce(k8sSvcName, serviceName3gpp string, client kubernetes.Interface) error {
	service, err := client.CoreV1().Services(config.namespace).Get(k8sSvcName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if service.Spec.Selector == nil || len(service.Spec.Selector) == 0 {
		return errors.New("spec selector of the service is empty")
	}
	options := metav1.ListOptions{FieldSelector: "status.phase=Running",
		LabelSelector: mapToString(service.Spec.Selector)}
	pods, err := client.CoreV1().Pods(config.namespace).List(options)
	if err != nil {
		return err
	}
	if !hasReadyPod(pods.Items) {
		return errors.New("service " + service.Name + " is not healthy")
	}
	logger.Debug("service ", service.Name, " is healthy")
	if err = sendHeartbeat(serviceName3gpp); err != nil {
		return err
	}
	return nil
}
