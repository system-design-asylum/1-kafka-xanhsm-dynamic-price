package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func getRideRequestHttpHandler(producer *KafkaProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rideRequest RideRequest
		if err := json.NewDecoder(r.Body).Decode(&rideRequest); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Set timestamp if not provided (or if 0)
		if rideRequest.RequestedAt == 0 {
			rideRequest.RequestedAt = time.Now().Unix()
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) // Context timeout for Kafka write
		defer cancel()

		if err := producer.ProduceKafkaRideRequestMsg(ctx, rideRequest); err != nil {
			http.Error(w, "Failed to produce ride request", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
