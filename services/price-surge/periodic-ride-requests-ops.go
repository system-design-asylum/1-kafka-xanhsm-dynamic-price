package main

import (
	"sync"
	"time"
)

func incrementZoneRideRequestCount(zoneID int, zoneRideRequests *sync.Map) {
	count, _ := zoneRideRequests.LoadOrStore(zoneID, 0)
	zoneRideRequests.Store(zoneID, count.(int)+1)
}

func periodicRideRequestOps(rideRequestsToProcess *sync.Map, zoneRideRequests *sync.Map, iterationDone chan struct{}) {
	for {
		rideRequestsToProcess.Range(func(key, value interface{}) bool {
			rideRequest := value.(RideRequest)
			zoneID := fromCoordToZoneID(rideRequest.PickupLat, rideRequest.PickupLng)
			incrementZoneRideRequestCount(zoneID, zoneRideRequests)
			return true
		})

		iterationDone <- struct{}{} // Signal end of one iteration
		time.Sleep(PERIODIC_OPS_INTERVAL_SECS)
	}
}
