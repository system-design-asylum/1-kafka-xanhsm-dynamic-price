package main

func fromCoordToZoneID(lat, lng float64) int {
	// Simple logic to determine zone ID based on coordinates
	latZone := int((lat - HCM_MIN_LAT) / LAT_STEP)
	lngZone := int((lng - HCM_MIN_LNG) / LON_STEP)
	return latZone*1000 + lngZone
}
