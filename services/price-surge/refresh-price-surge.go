package main

import (
	"log"
	"sync"
)

func refreshPriceSurgeOnPremise(
	driverOpsIterationDone chan struct{},
	rideRequestsOpsIterationDone chan struct{},
	zoneDriverCount *sync.Map,
	zoneRideRequests *sync.Map,
) {
	for {
		log.Println("Waiting for driver and ride request operations to complete...")
		<-driverOpsIterationDone
		<-rideRequestsOpsIterationDone

		// Refresh the price surge calculations here
		log.Println("Refreshing price surge calculations...")
		zoneDriverCount.Range(func(zoneID, driverCount interface{}) bool {
			zoneDemand, _ := zoneRideRequests.Load(zoneID)
			zoneSupply := driverCount.(int)
			log.Println("We're here")
			log.Printf("Zone %d: Drivers = %d, Ride Requests = %d", zoneID, zoneSupply, zoneDemand)

			if zoneDemand.(int) > 0 && zoneSupply > 0 {
				priceSurge := float64(zoneDemand.(int)) / float64(zoneSupply) // Calc leverage ratio
				log.Printf("Zone %d: Price Leverage = %.2f\n", zoneID, priceSurge)
			} else {
				log.Printf("Zone %d: No price leverage update needed.\n", zoneID)
			}

			return true
		})
	}
}
