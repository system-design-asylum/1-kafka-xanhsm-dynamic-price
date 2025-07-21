package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

// For analysis
// var heartbeatsCounter int = 0
// var lastHeartbeatRecordTime time.Time

const HEARTBEATS_BEFORE_LOGGING = 500

// closes the Kafka writer
func (p *KafkaProducer) Close() error {
	log.Println("Closing Kafka producer...")
	return p.writer.Close()
}

// Sends a heartbeat message to Kafka
func (p *KafkaProducer) ProduceKafkaHeartbeatMsg(ctx context.Context, heartbeat Heartbeat) error {
	// Marshal the heartbeat struct to JSON bytes
	messageValue, err := json.Marshal(heartbeat)
	if err != nil {
		return err
	}

	// Create a Kafka message.
	// Using DriverID as the key ensures all heartbeats for a specific driver
	// go to the same partition, maintaining order for that driver.
	message := kafka.Message{
		Key:   []byte(heartbeat.DriverID),
		Value: messageValue,
	}

	// Send the message to Kafka
	// With BatchTimeout, messages are buffered and sent in batches.
	// The WriteMessages call will block until the batch is sent or timeout occurs.
	err = p.writer.WriteMessages(ctx, message)
	if err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
		return err
	}

	// log.Printf("Produced heartbeat for DriverID: %s", heartbeat.DriverID)
	return nil
}
