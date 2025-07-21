package main

import (
	"database/sql"
	"log"
	"sync"
)

func main() {
	kafkaBroker := "localhost:9092"
	kafkaHeartbeatTopic := "driver-heartbeats"
	kafkaRideRequestsTopic := "ride-requests"
	kafkaGroupID := "heartbound-group"

	// Init db conn
	db, err := sql.Open("postgres", "port=5432 user=postgres password=postgres dbname=grid_db sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize global sync.Map to store heartbeats
	var driversInfo sync.Map

	// Initialize global sync.Map to store ride requests
	var rideRequestsToProcess sync.Map

	// Initialize 2 global sync.Map to store supply, demand of each zone
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
