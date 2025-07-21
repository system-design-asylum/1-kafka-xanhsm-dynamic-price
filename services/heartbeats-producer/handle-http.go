package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// getHeartbeatHttpHandler handles incoming HTTP heartbeat requests
func getHeartbeatHttpHandler(producer *KafkaProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var hb Heartbeat
		err := json.NewDecoder(r.Body).Decode(&hb)
		if err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Set timestamp if not provided (or if 0)
		if hb.Timestamp == 0 {
			hb.Timestamp = time.Now().Unix()
		}

		// Produce the heartbeat to Kafka
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) // Context timeout for Kafka write
		defer cancel()

		err = producer.ProduceKafkaHeartbeatMsg(ctx, hb)
		if err != nil {
			log.Printf("Error producing heartbeat to Kafka: %v", err)
			http.Error(w, "Failed to process heartbeat: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		// if heartbeatsCounter%HEARTBEATS_BEFORE_LOGGING == 0 {
		// 	if !lastHeartbeatRecordTime.IsZero() {
		// 		heartbeatsPerSec := float64(HEARTBEATS_BEFORE_LOGGING) / time.Since(lastHeartbeatRecordTime).Seconds()
		// 		log.Printf("Received %.2f heartbeats per seconds on average\n", heartbeatsPerSec)

		// 		// Reset the counter and last record time every 500 heartbeats
		// 		heartbeatsCounter = 0
		// 	}
		// 	lastHeartbeatRecordTime = time.Now()
		// }
	}
}
