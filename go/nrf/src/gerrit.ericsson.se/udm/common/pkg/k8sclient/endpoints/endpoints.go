package endpoints

import (
	"encoding/json"
	"fmt"
	"time"

	k8sclient "gerrit.ericsson.se/udm/common/pkg/k8sclient/client"
	"gerrit.ericsson.se/udm/common/pkg/log"
)

// Event is a data type defined for endpoints update
type Event int32

// EndpointsUpdateHandler is an interface, the user must implement it if struct endpoints.Watcher
// is used
type EndpointsUpdateHandler interface {
	HandleIPAddressUpdate(Event, []string)
	GetIPAddresses() []string
	GetInterval() int
}

// Watcher is used to monitor the update of a specified endpoints
type Watcher struct {
	endpointsRepo map[string]EndpointsUpdateHandler
	endpointsName chan string
	quitAll       chan bool
	quitOne       chan bool
}

const (
	// NOCHANGE defines the event when endpoints is not changed
	NOCHANGE Event = iota
	// ADD defines the event when endpoints is created
	ADD
	// UPDATE defines the event when endpoints is changed
	UPDATE
	// REMOVE defines the event when endpoints is deleted
	REMOVE
)

const (
	watchURI        = `/api/v1/namespaces/%s/endpoints/%s`
	defaultInterval = 3
)

var (
	watcher *Watcher
)

// NewWatcher creates a isntance of Watcher
func NewWatcher() *Watcher {
	if watcher == nil {
		watcher = &Watcher{
			endpointsRepo: make(map[string]EndpointsUpdateHandler),
			endpointsName: make(chan string),
			quitAll:       make(chan bool),
			quitOne:       make(chan bool, 10),
		}
	}
	return watcher
}

// AddEndpointsToWatcher adds an specified endpoints to the watcher list
func (w *Watcher) AddEndpointsToWatcher(namespace, endpoints string, handler EndpointsUpdateHandler) {
	if namespace != "" && endpoints != "" && handler != nil {
		URI := fmt.Sprintf(watchURI, namespace, endpoints)
		w.endpointsRepo[URI] = handler
		w.endpointsName <- URI
	}
}

// StartWatcher starts the monitor goroutine
func (w *Watcher) StartWatcher() {
	log.Debugf("Start endpoints watcher")
	go func() {
		for {
			select {
			case URI := <-w.endpointsName:
				w.startWatcherEndpoints(URI)
			case <-w.quitAll:
				log.Debugf("Stop endpoints watcher")
				close(w.quitOne)
				return
			}
		}
	}()
}

// startWatcherEndpoints starts the monitor for a specified endpoints
func (w *Watcher) startWatcherEndpoints(URI string) {
	updateHandler := w.endpointsRepo[URI]
	monitorInterval := updateHandler.GetInterval()
	if monitorInterval <= 0 {
		monitorInterval = defaultInterval
	}
	log.Debugf("Start watcher for endpoints %s with interval %d", URI, monitorInterval)
	ticker := time.NewTicker(time.Second * time.Duration(monitorInterval))

	go func() {
		for {
			select {
			case <-ticker.C:
				// get endpoints from k8s api server
				ipAddresses := getEndpoints(URI)

				// compare to the provious one
				event := w.compare(updateHandler.GetIPAddresses(), ipAddresses)

				// call the updateHandler
				if event != NOCHANGE {
					updateHandler.HandleIPAddressUpdate(event, ipAddresses)
				}
			case <-w.quitOne:
				ticker.Stop()
				log.Debugf("Stop watcher for endpoints %s", URI)
				return
			}
		}
	}()
}

// StopWatcher stop the monitor
func (w *Watcher) StopWatcher() {
	w.quitAll <- true
}

// compare compares the old IP address list to the new one
func (w *Watcher) compare(oldIPAddresses, newIPAddresses []string) Event {
	var event Event

	if oldIPAddresses == nil && newIPAddresses == nil {
		event = NOCHANGE
	} else if oldIPAddresses == nil && newIPAddresses != nil {
		event = ADD
	} else if oldIPAddresses != nil && newIPAddresses == nil {
		event = REMOVE
	} else if oldIPAddresses != nil && newIPAddresses != nil {
		if len(oldIPAddresses) != len(newIPAddresses) {
			event = UPDATE
		} else {
			updated := false
			mapNewIP := make(map[string]bool)
			for _, newIP := range newIPAddresses {
				mapNewIP[newIP] = true
			}

			for _, oldIP := range oldIPAddresses {
				if !mapNewIP[oldIP] {
					updated = true
					break
				}
			}

			if updated {
				event = UPDATE
			} else {
				event = NOCHANGE
			}
		}
	}

	return event
}

// getEndpoints returns IP addresses of a specified Endpoints
// - URI : specifies the URI of the Endpoints
func getEndpoints(URI string) []string {
	log.Debugf("Try to retrieve endpoints %s", URI)
	resp, err := k8sclient.GetK8sAPIClient().SendK8sAPIRequest("GET", URI, nil)
	if err != nil {
		log.Warnf("Retrieve endpoints %s failed. %s", URI, err.Error())
		return nil
	}

	endpoints := &TEndpoints{}
	err = json.Unmarshal(resp, endpoints)
	if err != nil {
		log.Warnf("Retrieve endpoints %s failed. Unmarshal error, %s", URI, err.Error())
		return nil
	}

	var ipAddresses []string
	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			if address.IP != "" {
				ipAddresses = append(ipAddresses, address.IP)
			}
		}
	}

	if ipAddresses == nil {
		log.Warnf("Endpoints %s is unavailable.", URI)
		return nil
	}

	log.Debugf("Retrieve endpoints %s successfully, it's ip addresses is %v.", URI, ipAddresses)
	return ipAddresses
}

// GetEventName returns name of endpoints update event
func GetEventName(event Event) string {
	switch event {
	case NOCHANGE:
		return "NOCHANGE"
	case ADD:
		return "ADD"
	case UPDATE:
		return "UPDATE"
	case REMOVE:
		return "REMOVE"
	}

	return ""
}
