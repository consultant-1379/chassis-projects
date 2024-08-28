package cmproxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"gerrit.ericsson.se/udm/common/pkg/httpclient"
	"gerrit.ericsson.se/udm/common/pkg/log"
)

type cmSubscription struct {
	ID                       string   `json:"id,omitempty"`
	ConfigName               string   `json:"configName"`
	Event                    []string `json:"event"`
	UpdateNotificationFormat string   `json:"updateNotificationFormat"`
	LeaseSeconds             int      `json:"leaseSeconds"`
	Callback                 string   `json:"callback"`
}

const (
	// EventConfigCreated is event "configCreated" in notifications
	EventConfigCreated = "configCreated"
	// EventConfigUpdated is event "configUpdated" in notifications
	EventConfigUpdated = "configUpdated"
	// EventConfigDeleted is event "configDeleted" in notifications
	EventConfigDeleted = "configDeleted"

	// NtfFormatFull is updateNotificationFormat "full" in notifications
	NtfFormatFull = "full"
	// NtfFormatPatch is updateNotificationFormat "patch" in notifications
	NtfFormatPatch = "patch"
)

func newSubscription(id, configName string) *cmSubscription {
	var (
		defaultEvent     = []string{EventConfigCreated, EventConfigUpdated, EventConfigDeleted}
		defaultCallback  = "kafka:cmproxy"
		defaultFormat    = NtfFormatPatch
		defaultLeaseSecs = 3600
	)

	s := &cmSubscription{}

	s.ID = id
	s.ConfigName = configName
	s.Event = defaultEvent
	s.Callback = defaultCallback
	s.UpdateNotificationFormat = defaultFormat
	s.LeaseSeconds = defaultLeaseSecs

	return s
}

func cmSubscriptionCreated(id string, s *cmSubscription) bool {
	if id == "" ||
		s == nil {
		return false
	}

	cmSubscriptionsLock.Lock()
	defer cmSubscriptionsLock.Unlock()

	if cmSubscriptions == nil {
		return false
	}

	// Ignored ID in subscriptions
	s.ID = ""
	cmSubscriptions[id] = s

	if status == runningWithMessageBus {
		err := putSubscription(id, s)
		if err != nil {
			return false
		}
	}
	return true
}

func cmSubscriptionDeleted(id string) {
	cmSubscriptionsLock.Lock()
	defer cmSubscriptionsLock.Unlock()

	delete(cmSubscriptions, id)
}

func doSubscriptions(method, url string, s *cmSubscription) (*httpclient.HttpRespData, error) {
	if s != nil {
		jsonBuf, err := json.Marshal(*s)
		if err != nil {
			log.Errorf("failed to Marshal subscription %v, %s", *s, err.Error())
			return nil, err
		}
		return cmHTTPClient.HttpDoJsonBody(method, url, bytes.NewBuffer(jsonBuf))
	}

	return cmHTTPClient.HttpDoJsonBody(method, url, nil)
}

func putSubscription(subscriptionID string, s *cmSubscription) error {
	resp, err := doSubscriptions("PUT", getCmmURL("subscriptions", subscriptionID), s)
	if err != nil {
		log.Errorf("failed to PUT subscriptions %s, %s", subscriptionID, err.Error())
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		s.ID = subscriptionID
		resp, err = doSubscriptions("POST", getCmmURL("subscriptions", ""), s)
		s.ID = ""
		if err != nil {
			log.Errorf("failed to POST subscriptions %s, %s", subscriptionID, err.Error())
			return err
		}
		if resp.StatusCode != http.StatusCreated {
			e := errors.New(string(resp.Body))
			log.Errorf("failed to POST subscriptions, %s", e.Error())
			return e
		}
	default:
		e := errors.New(string(resp.Body))
		log.Errorf("failed to PUT subscriptions %s", e.Error())
		return e
	}
	return nil
}

func deleteSubscription(subscriptionID string) {

	resp, err := doSubscriptions("DELETE", getCmmURL("subscriptions", subscriptionID), nil)
	if err != nil {
		log.Errorf("failed to DELETE subscriptions %s, %s", subscriptionID, err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		log.Errorf("failed to DELETE subscriptions, %s", string(resp.Body))
	}
}
