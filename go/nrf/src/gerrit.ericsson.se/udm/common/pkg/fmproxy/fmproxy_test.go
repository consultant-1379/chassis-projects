package fmproxy

import (
	"fmt"
	"testing"

	//"gerrit.ericsson.se/udm/common/pkg/fmproxy/mock_sarama"
	//"github.com/golang/mock/gomock"
)
/*
func TestStructureAlarmMsg(t *testing.T) {

	jsonAddtion := &AddtionMultiKeyValue{Key: "connection1", Value: "failure"}
	alarm := &AlarmInfo{IsAutoResend: true,
		FaultName:             "ausftest1",
		FaultyResource:        "eric-nrf-management",
		Expiration:            40,
		AdditionalInformation: jsonAddtion,
	}
	alarmMsg := structureAlarmMsg(alarm, true)
	if alarmMsg != `{ "version": "0.2", "faultName": "ausftest1", "serviceName": "", "faultyResource": "eric-nrf-management", "expiration": 40, "additionalInformation": { "connection1": "failure" } }` {
		t.Fatalf("StructureAlarmMsg fail!")
	}

	// test: no expiratime and without addtionalinfo
	alarm = &AlarmInfo{IsAutoResend: true,
		FaultName:      "ausftest1",
		FaultyResource: "eric-nrf-management",
		Expiration:     -1,
	}
	alarmMsg = structureAlarmMsg(alarm, true)

	if alarmMsg != `{ "version": "0.2", "faultName": "ausftest1", "serviceName": "", "faultyResource": "eric-nrf-management", "expiration": 0 }` {
		t.Fatalf("StructureAlarmMsg fail!")
	}

	// test: expiratime uses default value and clean alarm
	alarm = &AlarmInfo{IsAutoResend: true,
		FaultName:      "ausftest1",
		FaultyResource: "eric-nrf-management",
		Expiration:     9,
	}
	alarmMsg = structureAlarmMsg(alarm, false)

	if alarmMsg != `{ "version": "0.2", "faultName": "ausftest1", "serviceName": "", "faultyResource": "eric-nrf-management", "severity": "Clear", "expiration": 300 }` {
		t.Fatalf("StructureAlarmMsg fail!")
	}

}
*/
func TestInit(t *testing.T) {

	err := Init("eric-data-message-bus-kf:9092", "")
	if err == nil {
		t.Fatalf("init alarm fail")
	}

	err = Init("kafka:9092", "nrfservice")
	if kafkaConnection != "kafka:9092" {

		t.Fatalf("init alarm fail")
	}

	err = Init("", "nrfservice")
	if kafkaConnection != DefaultKafka {
		fmt.Printf("ddd")
		t.Fatalf("init alarm fail")
	}

}
/*
func TestSendAlarm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSyncProducer := mock_sarama.NewMockSyncProducer(ctrl)
	syncProducer = mockSyncProducer

	// faultName is empty
	alarm := &AlarmInfo{IsAutoResend: false,
		FaultName:      "",
		FaultyResource: "eric-nrf-management-6745577c8d-t2vq8",
	}
	err := SendAlarm(alarm, true)
	if err == nil {
		t.Fatalf("Send Alarm fail")
	}

	// FaultyResource is empty
	alarm = &AlarmInfo{IsAutoResend: false,
		FaultName:      "test1",
		FaultyResource: "",
	}
	err = SendAlarm(alarm, true)
	if err == nil {
		t.Fatalf("Send Alarm fail")
	}

	jsonAddtion := &AddtionMultiKeyValue{Key: "connection1", Value: "failure"}
	alarm = &AlarmInfo{IsAutoResend: true,
		FaultName:             "test1",
		FaultyResource:        "eric-nrf-management-6745577c8d-t2vq8",
		Expiration:            600,
		AdditionalInformation: jsonAddtion,
	}

	// send alarm
	mockSyncProducer.EXPECT().SendMessage(gomock.Any()).Return(int32(1), int64(1), nil)
	err = SendAlarm(alarm, true)
	_, found := resendFaultMap["test1"]
	if !found {
		t.Fatalf("Send Alarm fail")
	}

	// clean alarm
	mockSyncProducer.EXPECT().SendMessage(gomock.Any()).Return(int32(1), int64(1), nil)
	err = SendAlarm(alarm, false)
	_, found = resendFaultMap["test1"]
	if found {
		t.Fatalf("Send Alarm fail")
	}

	// not resend automatically
	alarm.IsAutoResend = false
	mockSyncProducer.EXPECT().SendMessage(gomock.Any()).Return(int32(1), int64(1), nil)
	err = SendAlarm(alarm, true)
	_, found = resendFaultMap["test1"]
	if found {
		t.Fatalf("Send Alarm fail")
	}
}*/
