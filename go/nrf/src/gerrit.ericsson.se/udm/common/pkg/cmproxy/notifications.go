package cmproxy

import (
	"encoding/json"
	"os"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"gerrit.ericsson.se/udm/common/pkg/msgbus"
)

type cmNotification struct {
	ConfigName string          `json:"configName"`
	Event      string          `json:"event"`
	BaseETag   string          `json:"baseETag,omitempty"`
	ConfigETag string          `json:"configETag,omitempty"`
	Patch      json.RawMessage `json:"patch,omitempty"`
	Data       json.RawMessage `json:"data,omitempty"`
}

func newNotification(rawData []byte) *cmNotification {
	notification := &cmNotification{}

	err := json.Unmarshal(rawData, notification)
	if err != nil {
		log.Errorf("failed to Unmarshal notification from message bus, %s, %s.", string(rawData), err.Error())
		return nil
	}
	return notification
}

func validateNotification(msg []byte) *cmNotification {
	n := newNotification(msg)
	if n == nil {
		return nil
	}
	if n.ConfigName == "" || n.Event == "" {
		log.Errorf("configName or event is missing in notification from message bus, %v", n)
		return nil
	}
	return n
}

func notificationsHandler(rawMessage []byte) {
	if status == initializing || status == shuttingDown {
		return
	}

	// Patch message from CMM includes many quotes
	// Ignored this quotes before deliver to App Service
	// Quotes removed by cm mediator in latest version - 1.0.0-225

	//	message, e := strconv.Unquote(string(rawMessage))
	//	if e != nil {
	//		log.Errorf("failed to Unquote notification from message bus, %s", e.Error())
	//		return
	//	}

	n := validateNotification(rawMessage)
	if n == nil {
		log.Errorf("failed to validate notification from message bus")
		return
	}

	// Default - updatedFormat:full, rawData:json:data
	// Include conditions:
	// configCreated - key: configName configETag
	// configDeleted - key: configName
	format := NtfFormatFull
	rawData := n.Data
	if n.Patch != nil {
		format = NtfFormatPatch
		rawData = n.Patch
	}
	cmConfigUpdated(n.Event, n.ConfigName, format, n.BaseETag, n.ConfigETag, rawData)
}

func initMessagebus() bool {
	if cmMessageBus == nil {
		cmMessageBus = msgbus.NewMessageBus(os.Getenv("MESSAGE_BUS_KAFKA"))
	}
	if cmMessageBus == nil {
		log.Errorf("failed to initialize message bus for cmproxy")
		return false
	}
	for topic := range notificationsTopics {
		err := cmMessageBus.ConsumeTopic(topic, notificationsHandler)
		if err != nil {
			log.Errorf("failed to ConsumeTopic %s for cmproxy", topic)
			return false
		}
	}
	log.Infof("initialize message bus for cmproxy done")
	return true
}

func closeMessageBus() {
	if cmMessageBus != nil {
		_ = cmMessageBus.Close()
	}
}

// SetMessageBus provide interface for importing extenal messagebus instance
func SetMessageBus(m *msgbus.MessageBus) {
	if cmMessageBus == nil {
		cmMessageBus = m
	}
}
