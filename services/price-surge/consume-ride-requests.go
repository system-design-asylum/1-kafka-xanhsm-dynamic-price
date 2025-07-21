package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
)

func (c *KafkaConsumer) StartConsumeRideRequests(rideRequestsToProcess *sync.Map, zoneRideRequests *sync.Map) {
	for {
		m, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var rideRequest RideRequest
		if err := json.Unmarshal(m.Value, &rideRequest); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		rideRequestsToProcess.Store(rideRequest.CustomerID, rideRequest)
	}

}
