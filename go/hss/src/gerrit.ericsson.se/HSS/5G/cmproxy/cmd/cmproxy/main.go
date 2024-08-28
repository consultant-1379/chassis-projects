package main

import (
	"flag"
	"fmt"
	"log"

	"gerrit.ericsson.se/HSS/5G/cmproxy/messagebus"
	"gerrit.ericsson.se/HSS/5G/cmproxy/service"
	"gerrit.ericsson.se/HSS/5G/cmproxy/statistics"
)

var (
	listeningPort     = flag.Int("listenport", 9080, "HTTP2 listening port for incoming requests")
	configurationPort = flag.Int("configport", 9088, "HTTP2 listening port for incoming configuration requests")
	notifEndpoint     = flag.String("notifendpoint", "http://localhost:9081", "app notification endpoint")
	subscribe         = flag.String("subscribe", "ericsson-udm", "name of subscription")
	kafkaEndpoint     = flag.String("kafkaendpoint", "kafka:1000", "Kafka endpoint")
	topic             = flag.String("topic", "eric-udm-config", "topic to subscribe")
	cmendpoint        = flag.String("cmendpoint", "http://eric-cm-mediator:5003", "CM service endpoint")
	apiroot           = flag.String("apiroot", "/cm/api/v1/configurations/", "api root")
)

func main() {
	flag.Parse()

	cmUri := *cmendpoint + *apiroot + *subscribe
	log.Println(cmUri)
	s := service.NewServer(cmUri, *notifEndpoint)
	if s == nil {
		log.Fatal("null server")
	}

	//Start http server
	go func() {
		statistics.StartHttpServer(*subscribe, *configurationPort)
	}()

	mbServer := messagebus.NewKafka(*kafkaEndpoint, *topic, s.GetCB())
	mbServer.Start()
	defer mbServer.Stop()

	cmPort := fmt.Sprintf("%d", *listeningPort)
	if err := s.Start(cmPort); err != nil {
		log.Fatal("The server can not be started. " + err.Error())
	}
}
