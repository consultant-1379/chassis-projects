package msgbus

import (
	"os"
	"strings"
	"sync"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
)

var (
	// PendingInitializeFailed is used by UT/BT
	PendingInitializeFailed = true
)

// MessageBus struct define
type MessageBus struct {
	wg    sync.WaitGroup
	kafka []string

	// Interface for Consumers
	consumer consumer
	// Channel to monitor
	consumerChan chan *partitionConsumer

	// Interface for SyncProducer
	producer syncProducer
	// Interface for SyncProducer
	producerWithPartition syncProducer
	// Reserved Interface for AsyncProducer
	producerAsync asyncProducer
}

const (
	// Default Kafka HostName:Port
	defaultKafka = "eric-data-message-bus-kf:9092"
	// Default retry interval of connecting to kafka
	defaultRetryInterval = 5
)

// InitMessageBusLog provide API for Message Bus Logging initialization
func InitMessageBusLog(level log.Level, networkFunc, podIP, serviceID string) {
	// set log level. the value can be ErrorLevel/WarnLevel/InfoLevel/DebugLevel
	log.SetLevel(level)
	// set output
	log.SetOutput(os.Stdout)
	// set network function name, it will be displayed in log output
	log.SetNF(networkFunc)
	// set pod ip, it will be displayed in log output
	log.SetPodIP(podIP)
	// set log format to user-defined json
	// set service ID, it will be displayed in log output
	log.SetServiceID(serviceID)
	// set log format to user-defined json
	log.SetFormatter(&log.JSONFormatter{})
}

// NewMessageBus create Consumers and Producers in MessageBus
func NewMessageBus(kafka string) *MessageBus {
	m := &MessageBus{}

	if kafka == "" {
		log.Infof("kafka connection is not provided, using default value for msgbus")
		m.kafka = strings.Split(defaultKafka, ",")
	} else {
		m.kafka = strings.Split(kafka, ",")
	}
	log.Infof("kafka brokers: %s", m.kafka)

	err := m.createProducer()
	for err != nil && PendingInitializeFailed {
		log.Debugf("retry to create producer")
		m.deleteProducer()

		time.Sleep(defaultRetryInterval * time.Second)
		err = m.createProducer()
	}

	err = m.createConsumer()
	for err != nil && PendingInitializeFailed {
		log.Debugf("retry to create consumer")
		m.deleteConsumer()

		time.Sleep(defaultRetryInterval * time.Second)
		err = m.createConsumer()
	}
	m.consumerChan = make(chan *partitionConsumer)

	go func() {
		m.monitorConsumer()
	}()

	return m
}

// Close function wiil stop goroutines and close connectons between msgbus and kafka
func (m *MessageBus) Close() error {
	m.deleteProducer()
	m.deleteConsumer()
	m.wg.Wait()
	return nil
}
