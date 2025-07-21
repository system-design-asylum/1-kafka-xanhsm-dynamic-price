package main

import (
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type AvailableStatus int

const (
	Unavailable AvailableStatus = iota // 0
	Available                          // 1
)

type Heartbeat struct {
	DriverID        string          `json:"driver_id"`
	Latitude        float64         `json:"latitude"`
	Longitude       float64         `json:"longitude"`
	Timestamp       int64           `json:"timestamp"` // Unix timestamp in seconds
	AvailableStatus AvailableStatus `json:"available_status"`
}

type KafkaProducer struct {
	writer *kafka.Writer
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
