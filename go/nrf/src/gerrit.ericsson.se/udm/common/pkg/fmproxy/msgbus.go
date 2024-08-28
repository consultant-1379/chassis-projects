package fmproxy

import (
	"errors"
	"strings"
	"time"

	"gerrit.ericsson.se/udm/common/pkg/log"
	"github.com/Shopify/sarama"
)

const (
	// Default Kafka HostName:Port
	DefaultKafka = "eric-data-message-bus-kf:9092"

	// AlarmProducerTopic is topic for alarm
	AlarmProducerTopic = "AdpFaultIndication"
)

var (
	kafkaConnection string
	syncProducer    sarama.SyncProducer
	reConnFlag      = make(chan int)
)

func initMsgbus() error {
	var err error

	// New sarama Config, and set return successes to true
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	go func() {
		for {
			<-reConnFlag
			if syncProducer != nil {
				_ = syncProducer.Close()
				syncProducer = nil
			}
			for {
				syncProducer, err = sarama.NewSyncProducer(strings.Split(kafkaConnection, ","), config)
				if err != nil {
					log.Errorf("re-connecting message bus fail , %s", err.Error())
				} else {
					log.Infof("re-connect message bus succeessfully!")
					break
				}
				time.Sleep(time.Second)
			}
		}
	}()

	syncProducer, err = sarama.NewSyncProducer(strings.Split(kafkaConnection, ","), config)
	if err != nil {
		log.Errorf("connecting message bus fail, %s", err.Error())
		reConnFlag <- int(1)
		return err
	} else {
		log.Infof("connect message bus succeessfully!")
	}

	return nil
}

func sendMsg(alarmPara *AlarmInfo) error {

	alarmMsg := structureAlarmMsg(alarmPara)
	alarmKey := serviceName + "." + alarmPara.FaultName
	msg := &sarama.ProducerMessage{
		Topic:     AlarmProducerTopic,
		Key:       sarama.StringEncoder(alarmKey),
		Value:     sarama.StringEncoder(alarmMsg),
		Timestamp: alarmPara.timestamp,
	}
	if syncProducer != nil {
		log.Warning("Alarm send to Msgbus:")
		log.Warning(msg)
		partition, offset, err := syncProducer.SendMessage(msg)
		if err != nil {
			log.Errorf("failed to send alarm message: %s", err.Error())
			return err
		} else {
			log.Debugf("send alarm successfully, partition:%d offset:%d", partition, offset)
		}
	} else {
		return errors.New("message bus is not available.")
	}

	return nil
}

func Close() error {
	err := syncProducer.Close()
	return err
}
