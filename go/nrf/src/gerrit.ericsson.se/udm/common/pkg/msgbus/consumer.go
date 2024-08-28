package msgbus

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"github.com/Shopify/sarama"
)

// Consumer package interface sarama.PartitionConsumer
type consumer struct {
	intf sarama.Consumer

	partitionConsumers     map[string]*partitionConsumer
	partitionConsumersLock sync.Mutex
}

type partitionConsumer struct {
	intf sarama.PartitionConsumer

	topic     string
	partition int32
	offset    int64

	handler          func([]byte)
	handlerWithTopic func(string, []byte)
	handlerMsg       func(*sarama.ConsumerMessage)

	consumed, errors int
}

func (m *MessageBus) createConsumer() error {
	var err error

	clientConfig := sarama.NewConfig()
	clientConfig.Consumer.Return.Errors = true
	clientConfig.Metadata.Retry.Backoff = 2 * time.Second

	//clientConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	//m.consumer.intf, err = sarama.NewConsumer(m.kafka, nil)

	m.consumer.intf, err = sarama.NewConsumer(m.kafka, clientConfig)
	if err != nil {
		log.Errorf("failed to create NewConsumer, %s", err.Error())
		return err
	}

	m.consumer.partitionConsumers = make(map[string]*partitionConsumer)
	return nil
}

func (m *MessageBus) createPartitionConsumer() {
	var err error

	m.consumer.partitionConsumersLock.Lock()
	defer m.consumer.partitionConsumersLock.Unlock()

	for _, v := range m.consumer.partitionConsumers {
		if m.consumer.intf == nil {
			// If sarama.consumer is not ready or closed
			// break out this loop and try it 2s later
			break
		}

		if v.intf != nil {
			continue
		}
		v.intf, err = m.consumer.intf.ConsumePartition(v.topic, v.partition, v.offset)
		if err != nil {
			log.Errorf("failed to create ConsumePartition, %s", err.Error())
			continue
		}
		log.Infof("topic/partition (%s/%d) is consumed", v.topic, v.partition)
		m.goPartitionConsumer(v)
	}
}

// GetPartitions function get partition ID list from kafka
func (m *MessageBus) GetPartitions(topic string) ([]int32, error) {
	partitionIDs, err := m.consumer.intf.Partitions(topic)
	for err != nil ||
		len(partitionIDs) == 0 {
		log.Errorf("failed to get partitions on topic %s", topic)
		time.Sleep(defaultRetryInterval * time.Second)
		partitionIDs, err = m.consumer.intf.Partitions(topic)
	}
	log.Infof("topic %s has partitions %v ", topic, partitionIDs)
	return partitionIDs, err
}

// ConsumeTopic function register callback handler to the topic
func (m *MessageBus) ConsumeTopic(topic string, callback func([]byte)) error {
	if m.consumer.intf == nil {
		return errors.New("no Consumer in message bus")
	}

	// Get the list of all partition IDs for the given topic
	partitionIDs, err := m.GetPartitions(topic)
	if err != nil {
		return err
	}

	for _, partition := range partitionIDs {
		m.consumerChan <- &partitionConsumer{
			intf:      nil,
			topic:     topic,
			partition: partition,
			offset:    sarama.OffsetNewest,
			handler:   callback,
		}
	}
	return nil
}

// ConsumeTopic function register callback handler to the topic
func (m *MessageBus) ConsumeTopicPlusTopicName(topic string, callback func(string, []byte)) error {
	if m.consumer.intf == nil {
		return errors.New("no Consumer in message bus")
	}

	// Get the list of all partition IDs for the given topic
	partitionIDs, err := m.GetPartitions(topic)
	if err != nil {
		return err
	}

	for _, partition := range partitionIDs {
		m.consumerChan <- &partitionConsumer{
			intf:             nil,
			topic:            topic,
			partition:        partition,
			offset:           sarama.OffsetNewest,
			handlerWithTopic: callback,
		}
	}
	return nil
}

// ConsumeTopicWithPartition advanced usage of message bus
func (m *MessageBus) ConsumeTopicWithPartition(topic string, partition int32, offset int64, callback func(*sarama.ConsumerMessage)) error {
	var err error

	if m.consumer.intf == nil {
		return errors.New("no Consumer in message bus")
	}

	// Get the list of all partition IDs for the given topic
	partitionIDs, err := m.GetPartitions(topic)
	if err != nil {
		return err
	}

	// Check whether Partition is valid for topic
	existed := false
	for _, p := range partitionIDs {
		if p == partition {
			existed = true
			break
		}
	}
	if !existed {
		log.Infof("partition %d does not exist on topic %s", partition, topic)
		return nil
	}

	m.consumerChan <- &partitionConsumer{
		intf:       nil,
		topic:      topic,
		partition:  partition,
		offset:     offset,
		handlerMsg: callback,
	}

	return nil
}

func (m *MessageBus) monitorConsumer() {
	var pc *partitionConsumer
	ok := true

	for ok {
		m.createPartitionConsumer()

		select {
		case pc, ok = <-m.consumerChan:
			if pc == nil {
				break
			}

			func() {
				m.consumer.partitionConsumersLock.Lock()
				defer m.consumer.partitionConsumersLock.Unlock()

				key := pc.topic + "-" + strconv.Itoa(int(pc.partition))
				// Check whether topic/partions is consumed
				_, existed := m.consumer.partitionConsumers[key]
				if existed {
					log.Infof("topic/partition (%s/%d) was already consumed", pc.topic, pc.partition)
				} else {
					if m.consumer.partitionConsumers != nil {
						m.consumer.partitionConsumers[key] = pc
					}
				}
			}()

		case <-time.After(defaultRetryInterval * time.Second):
		}
	}
}

func (m *MessageBus) goPartitionConsumer(pc *partitionConsumer) {
	if pc == nil ||
		pc.intf == nil {
		log.Errorf("failed to goroutine PartitionConsumer")
		return
	}

	m.wg.Add(1)
	go func(pc *partitionConsumer) {
		defer m.wg.Done()
		for msg := range pc.intf.Messages() {
			pc.consumed++
			log.Debugf("consumed message, topic:%s partition:%d offset:%d key:%s",
				msg.Topic, msg.Partition, msg.Offset, string(msg.Key))

			// deliver message to consumer handler
			if pc.handlerMsg != nil {
				pc.handlerMsg(msg)
			}

			if pc.handlerWithTopic != nil {
				pc.handlerWithTopic(msg.Topic, msg.Value)
			}

			if pc.handler != nil {
				pc.handler(msg.Value)
			}
		}
	}(pc)

	m.wg.Add(1)
	go func(pc *partitionConsumer) {
		defer m.wg.Done()
		for err := range pc.intf.Errors() {
			log.Errorf("failed to consume message, %s", err)
			pc.errors++
		}
	}(pc)
}

func (m *MessageBus) deleteConsumer() {
	// Stop monitorConsumer()
	if m.consumerChan != nil {
		close(m.consumerChan)
	}

	if m.consumer.intf != nil {
		m.consumer.partitionConsumersLock.Lock()
		defer m.consumer.partitionConsumersLock.Unlock()

		// trigger shutdown for consumer goroutines
		for key, pc := range m.consumer.partitionConsumers {
			if pc.intf != nil {
				_ = pc.intf.Close()
			}
			delete(m.consumer.partitionConsumers, key)
			//			m.consumer.partitionConsumers[key] = nil
		}
		_ = m.consumer.intf.Close()
	}
}
