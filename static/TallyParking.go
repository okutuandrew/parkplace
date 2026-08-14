package static

import (
	"encoding/json"
	"log"
	"os"
)

type ParkingTally struct {
	Total     int
	Available int
	Full      int
	Warning   int
	Spaces 	  int
}

type ParkingEntry struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Title   string  `json:"title"`
	Color   string  `json:"color"`
	Content string  `json:"content"`
	Spaces  int     `json:"spaces"` // 1. FIX: Add this field so Go can read "spaces" from the JSON
}

func ParkingData() ParkingTally {

	file, err := os.ReadFile("static/parkingData.json")
	if err != nil {
		log.Println("Can't read parkingData.json:", err)
		return ParkingTally{}
	}

	var entries []ParkingEntry
	err = json.Unmarshal(file, &entries)
	if err != nil {
		log.Println("Can't parse parkingData.json:", err)
		return ParkingTally{}
	}

	result := ParkingTally{}
	result.Total = len(entries)

	totalSpaces := 0
	for _, entry := range entries {
		totalSpaces += entry.Spaces
	}

	result.Spaces = totalSpaces // 2. FIX: Assign the calculated sum to your result struct

	return result
}
