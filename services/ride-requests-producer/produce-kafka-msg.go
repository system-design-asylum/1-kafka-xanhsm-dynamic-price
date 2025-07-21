package main

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

func (producer *KafkaProducer) ProduceKafkaRideRequestMsg(ctx context.Context, rideRequest RideRequest) error {
	// Marshal ride request to JSON
	value, err := json.Marshal(rideRequest)
	if err != nil {
		return err
	}

	// Create a new Kafka message
	msg := kafka.Message{
		Key:   []byte(rideRequest.CustomerID), // Use CustomerID as the key
		Value: value,
	}

	return producer.writer.WriteMessages(ctx, msg)
}
