package messagebus

import (
	"context"

	"log"

	"gerrit.ericsson.se/HSS/5G/cmproxy/statistics"
	kafka "github.com/segmentio/kafka-go"
)

// Callback .
type Callback func([]byte)

// Kafka .
type Kafka struct {
	reader *kafka.Reader
	cb     Callback
}

// NewKafka .
func NewKafka(endpoint, topic string, cb Callback) *Kafka {
	log.Println(endpoint)
	log.Println(topic)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{endpoint},
		Topic:     topic,
		Partition: 0,
		MinBytes:  1,
		MaxBytes:  1e6,
	})
	return &Kafka{
		reader: r,
		cb:     cb,
	}
}

func (k *Kafka) loop() {

	for {
		m, err := k.reader.ReadMessage(context.Background())
		statistics.Statistics.NumberOfKafkaMessages = statistics.Statistics.NumberOfKafkaMessages + 1
		if err != nil {
			log.Println(err)
		}
		k.cb(m.Value)
	}
}

// Start .
func (k *Kafka) Start() {
	go k.loop()
}

// Stop .
func (k *Kafka) Stop() {
	k.reader.Close()
}
