package main

import (
	"log"
	"net/http"
)

func main() {
	// Configuration
	kafkaBroker := "localhost:9092"
	kafkaTopicName := "driver-heartbeats"
	httpPort := ":8080"

	// Initialize Kafka Producer
	producer := NewKafkaProducer(kafkaBroker, kafkaTopicName)
	defer producer.Close()

	// Set up HTTP server
	http.HandleFunc("/heartbeat", getHeartbeatHttpHandler(producer))

	log.Printf("Starting Heartbeat Service on port %s", httpPort)
	log.Fatal(http.ListenAndServe(httpPort, nil))
}
