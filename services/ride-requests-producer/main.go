package main

import (
	"log"
	"net/http"
)

func main() {
	// Kafka Config
	kafkaBroker := "localhost:9092"
	kafkaTopicName := "ride-requests"
	httpPort := ":8081"

	// Kafka Producer
	producer := NewKafkaProducer(kafkaBroker, kafkaTopicName)
	defer producer.Close()

	http.HandleFunc("/ride-request", getRideRequestHttpHandler(producer))

	log.Printf("Starting Ride Request Handler Service on port %s", httpPort)
	log.Fatal(http.ListenAndServe(httpPort, nil))
}
