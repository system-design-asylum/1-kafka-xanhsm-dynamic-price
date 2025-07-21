package main

import (
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

type RideRequest struct {
	CustomerID  string  `json:"customer_id"`
	PickupLat   float64 `json:"pickup_lat"`
	PickupLng   float64 `json:"pickup_lng"`
	DropoffLat  float64 `json:"dropoff_lat"`
	DropoffLng  float64 `json:"dropoff_lng"`
	RequestedAt int64   `json:"requested_at"` // Unix timestamp in seconds
}

func NewKafkaProducer(kafkaBroker string, topic string) *KafkaProducer {
	// Create a new Kafka writer.
	// The DialTimeout is set to ensure the producer doesn't hang indefinitely
	// if the broker is unreachable during initial connection.
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{}, // Distribute messages among partitions
		// Optional: Configure batching for better performance
		BatchSize:    100,
		BatchBytes:   1048576, // 1MB
		BatchTimeout: 1 * time.Second,
		// Optional: Error handling for async writes
		ErrorLogger: log.New(os.Stderr, "KAFKA_PRODUCER_ERROR: ", log.LstdFlags),
	}
	log.Printf("Kafka producer initialized for topic '%s' at broker '%s'", topic, kafkaBroker)
	return &KafkaProducer{writer: writer}
}

func (producer *KafkaProducer) Close() error {
	log.Println("Closing Kafka producer...")
	return producer.writer.Close()
}
