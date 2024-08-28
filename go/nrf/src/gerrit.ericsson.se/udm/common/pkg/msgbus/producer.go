package msgbus

import (
	"errors"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"github.com/Shopify/sarama"
)

// Producer package interface sarama.AsyncProducer or sarama.SyncProducer
type syncProducer struct {
	intf              sarama.SyncProducer
	successes, errors int
}

type asyncProducer struct {
	intf                        sarama.AsyncProducer
	enqueued, successes, errors int
}

func (m *MessageBus) createProducer() error {
	var err error

	// New sarama Config, and set return successes to true
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	m.producer.intf, err = sarama.NewSyncProducer(m.kafka, config)
	if err != nil {
		log.Errorf("failed to create NewSyncProducer, %s", err.Error())
		return err
	}

	m.producerAsync.intf, err = sarama.NewAsyncProducer(m.kafka, config)
	if err != nil {
		log.Errorf("failed to create NewAsyncProducer, %s", err.Error())
		return err
	}
	m.goAsyncProducer()

	// Set manual partion for Producer. If the sarama.Config is not initailize,
	// the parameter "partition" in sarama.ProducerMessage can be useless
	config.Producer.Partitioner = sarama.NewManualPartitioner
	m.producerWithPartition.intf, err = sarama.NewSyncProducer(m.kafka, config)
	if err != nil {
		log.Errorf("failed to create NewSyncProducer, %s", err.Error())
		return err
	}

	return nil
}

func (m *MessageBus) deleteProducer() {
	if m.producer.intf != nil {
		_ = m.producer.intf.Close()
	}
	if m.producerWithPartition.intf != nil {
		_ = m.producerWithPartition.intf.Close()
	}
	if m.producerAsync.intf != nil {
		m.producerAsync.intf.AsyncClose()
	}
}

func (m *MessageBus) goAsyncProducer() {
	if m.producerAsync.intf == nil {
		return
	}

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		for range m.producerAsync.intf.Successes() {
			m.producerAsync.successes++
		}
	}()

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		for err := range m.producerAsync.intf.Errors() {
			log.Errorf("failed to deliver message: %s", err)
			m.producerAsync.errors++
		}
	}()
}

// SendMessage function produce message to kafka via sync producer
func (m *MessageBus) SendMessage(topic string, msg string) error {
	if topic == "" {
		return errors.New("topic name is empty")
	}

	if m.producer.intf == nil {
		return errors.New("no SyncProducer in message bus")
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   nil,
		Value: sarama.StringEncoder(msg),
	}

	//	partition, offset, err := m.producer.intf.SendMessage(message)
	partition, offset, err := m.producer.intf.SendMessage(message)
	if err != nil {
		log.Errorf("failed to produce message: %s", err.Error())
		m.producer.errors++
	} else {
		log.Debugf("produced message, partition:%d offset:%d", partition, offset)
		m.producer.successes++
	}

	return err
}

// SendMessageWithPartition advanced usage of message bus
// Not recommend to send message by this function
func (m *MessageBus) SendMessageWithPartition(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	if m.producerWithPartition.intf == nil {
		return 0, 0, errors.New("no SyncProducer (with manaul partion option) in message bus")
	}

	partition, offset, err = m.producerWithPartition.intf.SendMessage(msg)
	if err != nil {
		log.Errorf("failed to produce message: %s", err.Error())
		m.producerWithPartition.errors++
	} else {
		log.Debugf("produced message, partition:%d offset:%d", partition, offset)
		m.producerWithPartition.successes++
	}
	return partition, offset, err
}

// SendMessageAsync function produce message to kafka via async producer
func (m *MessageBus) SendMessageAsync(topic string, msg string) error {
	if m.producerAsync.intf == nil {
		return errors.New("no AsyncProducer in message bus")
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   nil,
		Value: sarama.StringEncoder(msg),
	}

	m.producerAsync.intf.Input() <- message
	m.producerAsync.enqueued++
	log.Debugf("produced message by async producer")

	// These returns should be useless in Async mode. Since they will
	// be assigned only after this message is realy delivered from queue.
	return nil
}
