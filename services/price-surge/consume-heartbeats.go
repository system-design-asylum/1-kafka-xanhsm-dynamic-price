package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

const NUM_CLIENTS = 13333
const HCM_MIN_LAT = 10.7
const HCM_MAX_LAT = 10.85
const HCM_MIN_LNG = 106.6
const HCM_MAX_LNG = 106.8
const LAT_STEP = 0.018
const LON_STEP = 0.0183

const HEARTBEATS_PER_BATCH = 100
const SQL_INSERT_TEMPLATE = "INSERT INTO app_data.drivers (driver_id, geom, status, current_zone_id, last_updated_at) VALUES %s"

// driverToHeartbeat: sync.Map to store heartbeats globally, thread-safe across multiple goroutines
func (c *KafkaConsumer) StartConsumeHeartbeats(driverInfos *sync.Map) {

	for {
		m, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var heartbeat Heartbeat
		if err := json.Unmarshal(m.Value, &heartbeat); err != nil {
			log.Printf("Error unmarshalling heartbeat message: %v", err)
			continue
		}
		log.Printf("Received Heartbeat: %s", m.Value)

		// Filter out drivers based on availability status
		if heartbeat.AvailableStatus != Available {
			log.Printf("Driver %s is unavailable, skipping heartbeat", heartbeat.DriverID)
			continue
		}

		driverZoneID := fromCoordToZoneID(heartbeat.Latitude, heartbeat.Longitude)
		newDriverInfo := &DriverInfo{
			Heartbeat: heartbeat,
			ZoneID:    driverZoneID,
		}
		driverInfos.Store(heartbeat.DriverID, newDriverInfo)
	}
}
