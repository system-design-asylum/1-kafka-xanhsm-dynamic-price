package main

import "github.com/segmentio/kafka-go"

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

type DriverInfo struct {
	Heartbeat
	ZoneID int `json:"zone_id"`
}

type RideRequest struct {
	CustomerID  string  `json:"customer_id"`
	PickupLat   float64 `json:"pickup_lat"`
	PickupLng   float64 `json:"pickup_lng"`
	DropoffLat  float64 `json:"dropoff_lat"`
	DropoffLng  float64 `json:"dropoff_lng"`
	RequestedAt int64   `json:"requested_at"` // Unix timestamp in seconds
}

type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(broker, topic, groupID string) *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{broker},
		Topic:       topic,
		GroupID:     groupID,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.LastOffset,
	})

	return &KafkaConsumer{
		reader: r,
	}
}
