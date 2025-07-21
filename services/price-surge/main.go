package main

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Kafka Config
	kafkaBroker := "localhost:9092"
	kafkaHeartbeatTopic := "driver-heartbeats"
	kafkaRideRequestsTopic := "ride-requests"
	kafkaGroupID := "heartbound-group"

	// Init DB Conn
	db, err := GetDbConn()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Init global sync.Map to store heartbeats
	var driversInfo sync.Map

	// Init global sync.Map to store ride requests
	var rideRequestsToProcess sync.Map

	// Init global sync.Map to store ride supply & demand of each zone
	var zoneDriverCount sync.Map
	var zoneRideRequests sync.Map

	go func() {
		heartbeatsConsumer := NewKafkaConsumer(kafkaBroker, kafkaHeartbeatTopic, kafkaGroupID)
		heartbeatsConsumer.StartConsumeHeartbeats(&driversInfo)
	}()

	go func() {
		rideRequestsConsumer := NewKafkaConsumer(kafkaBroker, kafkaRideRequestsTopic, kafkaGroupID)
		rideRequestsConsumer.StartConsumeRideRequests(&rideRequestsToProcess, &zoneRideRequests)
	}()

	driverOpsIterationDone := make(chan struct{})
	go periodicDriverOps(db, &driversInfo, &zoneDriverCount, driverOpsIterationDone)

	rideRequestOpsIterationDone := make(chan struct{})
	go periodicRideRequestOps(&rideRequestsToProcess, &zoneRideRequests, rideRequestOpsIterationDone)

	refreshPriceSurgeOnPremise(
		driverOpsIterationDone,
		rideRequestOpsIterationDone,
		&zoneDriverCount,
		&zoneRideRequests,
	)
}
