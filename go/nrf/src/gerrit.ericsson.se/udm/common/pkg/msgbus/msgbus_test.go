package msgbus

import (
	"fmt"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
)

func TestMessageBusInitialization(t *testing.T) {
	messagebus := &MessageBus{}

	t.Run("TestCreateProducer", func(t *testing.T) {
		messagebus.createProducer()
	})

	t.Run("TestCreateConsumer", func(t *testing.T) {
		messagebus.createConsumer()
	})

	t.Run("TestDeleteProducer", func(t *testing.T) {
		messagebus.deleteProducer()
	})
	t.Run("TestDeleteProducer2", func(t *testing.T) {
		messagebus.deleteProducer()
	})

	t.Run("TestDeleteProducer", func(t *testing.T) {
		messagebus.deleteConsumer()
	})
	t.Run("TestDeleteProducer2", func(t *testing.T) {
		messagebus.deleteConsumer()
	})

	messagebus.Close()
}

func TestMessageBusProducer(t *testing.T) {
	sp := mocks.NewSyncProducer(t, nil)
	defer func() {
		if err := sp.Close(); err != nil {
			t.Error(err)
		}
	}()

	sp.ExpectSendMessageAndSucceed()
	sp.ExpectSendMessageAndFail(sarama.ErrOutOfBrokers)

	messagebus := &MessageBus{}
	messagebus.producer.intf = sp
	err := messagebus.SendMessage("test", "test")
	if err != nil {
		t.Errorf("The first message should have been produced successfully, but got %s", err)
	}

	err = messagebus.SendMessage("test", "test")
	if err != sarama.ErrOutOfBrokers {
		t.Errorf("The third message should not have been produced successfully")
	}

	if err := sp.Close(); err != nil {
		t.Error(err)
	}
}

func TestMessageBusProducerWithPation(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.Return.Successes = true

	sp := mocks.NewSyncProducer(t, config)
	defer func() {
		if err := sp.Close(); err != nil {
			t.Error(err)
		}
	}()

	sp.ExpectSendMessageAndSucceed()
	sp.ExpectSendMessageAndSucceed()
	sp.ExpectSendMessageAndFail(sarama.ErrOutOfBrokers)

	messagebus := &MessageBus{}
	messagebus.producerWithPartition.intf = sp

	msg1 := &sarama.ProducerMessage{Topic: "test", Value: sarama.StringEncoder("test"), Partition: 1}
	_, offset, err := messagebus.SendMessageWithPartition(msg1)
	if err != nil {
		t.Errorf("The first message should have been produced successfully, but got %s", err)
	}
	if offset != msg1.Offset {
		t.Errorf("The first message should have been assigned offset 1, but got %d", offset)
	}
	//	if p != msg1.Partition {
	//		t.Errorf("The first message should have been assigned pattition 1, but got %d", p)
	//	}

	msg2 := &sarama.ProducerMessage{Topic: "test", Value: sarama.StringEncoder("test"), Partition: 2}
	_, offset, err = messagebus.SendMessageWithPartition(msg2)
	if err != nil {
		t.Errorf("The third message should not have been produced successfully")
	}
	if offset != msg2.Offset {
		t.Errorf("The second message should have been assigned offset 2, but got %d", offset)
	}
	//	if p != msg2.Partition {
	//		t.Errorf("The first message should have been assigned pattition 2, but got %d", p)
	//	}

	_, _, err = messagebus.SendMessageWithPartition(msg2)
	if err != sarama.ErrOutOfBrokers {
		t.Errorf("The third message should not have been produced successfully")
	}

	if err := sp.Close(); err != nil {
		t.Error(err)
	}
}

func TestMessageBusAsyncProducer(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	mp := mocks.NewAsyncProducer(t, config)

	mp.ExpectInputAndSucceed()
	mp.ExpectInputAndSucceed()
	mp.ExpectInputAndFail(sarama.ErrOutOfBrokers)

	messagebus := &MessageBus{}
	messagebus.producerAsync.intf = mp
	messagebus.SendMessageAsync("test 1", "send hello by aync producer")
	messagebus.SendMessageAsync("test 2", "send hello by aync producer")
	messagebus.SendMessageAsync("test 3", "send hello by aync producer")

	msg1 := <-mp.Successes()
	msg2 := <-mp.Successes()
	err1 := <-mp.Errors()

	if msg1.Topic != "test 1" {
		t.Error("Expected message 1 to be returned first")
	}

	if msg2.Topic != "test 2" {
		t.Error("Expected message 2 to be returned second")
	}

	if err1.Msg.Topic != "test 3" || err1.Err != sarama.ErrOutOfBrokers {
		t.Error("Expected message 3 to be returned as error")
	}

	if err := mp.Close(); err != nil {
		t.Error(err)
	}
}

func testHandler(msg []byte) {
	fmt.Printf("testHandler: %s\n", string(msg))
}
func testP3Handler(msg []byte) {
	fmt.Printf("testP3Handler: %s\n", string(msg))
}
func testP3Handler2(msg *sarama.ConsumerMessage) {
	fmt.Printf("testP3Handler2: %s\n", string(msg.Value))
}

func TestConsumerHandlesExpectations(t *testing.T) {
	consumer := mocks.NewConsumer(t, nil)
	defer func() {
		if err := consumer.Close(); err != nil {
			t.Error(err)
		}
	}()

	consumer.SetTopicMetadata(map[string][]int32{
		"test":    {0},
		"test-p3": {0, 1, 2},
		"other":   {0},
	})

	consumer.ExpectConsumePartition("test", 0, sarama.OffsetNewest).YieldMessage(&sarama.ConsumerMessage{Value: []byte("hello world")})
	consumer.ExpectConsumePartition("test", 0, sarama.OffsetNewest).YieldError(sarama.ErrOutOfBrokers)
	consumer.ExpectConsumePartition("other", 0, sarama.OffsetNewest).YieldMessage(&sarama.ConsumerMessage{Value: []byte("hello world")})
	//	consumer.ExpectConsumePartition("other", 0, sarama.OffsetNewest).YieldError(sarama.ErrOutOfBrokers)
	consumer.ExpectConsumePartition("test-p3", 0, sarama.OffsetNewest).YieldMessage(&sarama.ConsumerMessage{Value: []byte("hello world again")})
	consumer.ExpectConsumePartition("test-p3", 1, sarama.OffsetNewest).YieldMessage(&sarama.ConsumerMessage{Value: []byte("hello world again")})
	consumer.ExpectConsumePartition("test-p3", 2, sarama.OffsetNewest).YieldMessage(&sarama.ConsumerMessage{Value: []byte("hello world again")})

	messagebus := &MessageBus{}

	messagebus.consumerChan = make(chan *partitionConsumer)
	messagebus.consumer.intf = consumer
	messagebus.consumer.partitionConsumers = make(map[string]*partitionConsumer)

	go func() {
		messagebus.monitorConsumer()
	}()

	messagebus.ConsumeTopic("test", testHandler)
	messagebus.ConsumeTopic("test-p3", testP3Handler)

	// abnormal cases
	messagebus.ConsumeTopic("other", nil)
	messagebus.ConsumeTopicWithPartition("test-p3", 0, sarama.OffsetOldest, testP3Handler2)
	messagebus.ConsumeTopicWithPartition("test-p3", 4, sarama.OffsetOldest, testP3Handler2)

	//	pc_test0 := messagebus.consumer.partitionConsumers["test-0"].intf
	//	test0_msg := <-pc_test0.Messages()
	//	if test0_msg.Topic != "test" || test0_msg.Partition != 0 || string(test0_msg.Value) != "hello world" {
	//		t.Error("Message was not as expected:", test0_msg)
	//	}
	//	test0_err := <-pc_test0.Errors()
	//	if test0_err.Err != sarama.ErrOutOfBrokers {
	//		t.Error("Expected sarama.ErrOutOfBrokers, found:", test0_err.Err)
	//	}

	//	pc_test1 := messagebus.consumer.partitionConsumers["test-p3-1"].intf
	//	if pc_test1 == nil {
	//		t.Fatal("nil")
	//	}
	//	test1_msg := <-pc_test1.Messages()
	//	if test1_msg.Topic != "test" || test1_msg.Partition != 1 || string(test1_msg.Value) != "hello world again" {
	//		t.Error("Message was not as expected:", test1_msg)
	//	}

	//	pc_other0, err := consumer.ConsumePartition("other", 0, sarama.OffsetNewest)
	//	if err == nil {
	//		t.Fatal("nil")
	//	}
	//	other0_msg := <-pc_other0.Messages()
	//	if other0_msg.Topic != "other" || other0_msg.Partition != 0 || string(other0_msg.Value) != "hello other" {
	//		t.Error("Message was not as expected:", other0_msg)
	//	}
}
