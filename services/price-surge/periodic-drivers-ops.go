package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

const WORKING_INTERVAL = 3 * time.Second

func buildValueStringAndArgs(heartbeat Heartbeat) (string, []interface{}) {
	valueStr := fmt.Sprintf("($%d,$%d,$%d,$%d)", heartbeat.DriverID, heartbeat.Latitude, heartbeat.Longitude, heartbeat.Timestamp)
	valueArgs := []interface{}{
		heartbeat.DriverID, heartbeat.Latitude, heartbeat.Longitude, heartbeat.Timestamp,
	}
	return valueStr, valueArgs
}

func batchInsertDriverInfo(db *sql.DB, valueStringBatch []string, valueArgsBatch []interface{}) error {
	if len(valueStringBatch) == 0 {
		return nil
	}

	insertStmt := fmt.Sprintf(SQL_INSERT_TEMPLATE, strings.Join(valueStringBatch, ", "))
	_, err := db.Exec(insertStmt, valueArgsBatch...)
	return err
}

func incrementZoneDriverCount(zoneID int, zoneDriverCount *sync.Map) {
	driverCount, _ := zoneDriverCount.LoadOrStore(zoneID, 0)
	zoneDriverCount.Store(zoneID, driverCount.(int)+1)
}

// driverInfos : sync.Map to store DriverInfo globally, thread-safe across multiple goroutines
func periodicDriverOps(db *sql.DB, availableDrivers *sync.Map, zoneDriverCount *sync.Map, iterationDone chan struct{}) {

	for {
		heartbeatCount := 0

		valueStrBatch := make([]string, 0, HEARTBEATS_PER_BATCH)
		valueArgsBatch := make([]interface{}, 0, HEARTBEATS_PER_BATCH*4)

		availableDrivers.Range(func(key, value any) bool {
			driverInfo := value.(*DriverInfo)
			valueStr, valueArgs := buildValueStringAndArgs(driverInfo.Heartbeat)
			valueStrBatch = append(valueStrBatch, valueStr)
			valueArgsBatch = append(valueArgsBatch, valueArgs...)

			// Increment the count of drivers in the zone
			incrementZoneDriverCount(driverInfo.ZoneID, zoneDriverCount)

			heartbeatCount++
			if heartbeatCount%HEARTBEATS_PER_BATCH == 0 {
				if err := batchInsertDriverInfo(db, valueStrBatch, valueArgsBatch); err != nil {
					log.Printf("Error inserting batch: %v", err)
				} else {
					log.Printf("Inserted %d heartbeats into the database", heartbeatCount)
				}
				valueStrBatch = valueStrBatch[:0]
				valueArgsBatch = valueArgsBatch[:0]
			}

			return true
		})

		iterationDone <- struct{}{} // Signal end of one iteration
		time.Sleep(PERIODIC_OPS_INTERVAL_SECS)
	}
}
