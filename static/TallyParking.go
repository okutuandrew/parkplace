package static

import (
	"encoding/json"
	"log"
	"os"
	"time"
	"strings"
)

type ParkingTally struct {
	Total     int
	Available int
	Full      int
	Warning   int
	Spaces 	  int
	Attendants int 
}

type ParkingEntry struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Title   string  `json:"title"`
	Color   string  `json:"color"`
	Content string  `json:"content"`
	Spaces  int     `json:"spaces"`  // 1. FIX: Add this field so Go can read "spaces" from the JSON
}

type AttendantEntry struct {
	Username  string    `json:"username"`
	DateTime  time.Time `json:"date_time"`
	Longitude float64   `json:"longitude"`
	Latitude  float64   `json:"latitude"`
}

type AttendantTally struct {
	Total      int // Total log submissions in the JSON file
	Attendants int // Total unique usernames active
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
	seenTitles := make(map[string]bool)

	for _, entry := range entries {
		totalSpaces += entry.Spaces
		seenTitles[entry.Title] = true
	}

	result.Spaces = totalSpaces // 2. FIX: Assign the calculated sum to your result struct
	result.Attendants = len(seenTitles)  // FEATURE :Used title as names of the attendant  but will later  use  UserId

	return result
}

func TallyAttendantsData() AttendantTally {
	// 1. Read directly from your new targeted location
	file, err := os.ReadFile("static/TallyAttendants.json")
	if err != nil {
		log.Println("Can't read TallyAttendants.json:", err)
		return AttendantTally{}
	}

	// 2. Parse the JSON array into our matching slice type
	var entries []AttendantEntry
	err = json.Unmarshal(file, &entries)
	if err != nil {
		log.Println("Can't parse TallyAttendants.json:", err)
		return AttendantTally{}
	}

	result := AttendantTally{}
	
	// 3. The raw length represents the total count of logs recorded
	result.Total = len(entries)

	// 4. Track distinct individuals using a map as a unique set
	uniqueUsernames := make(map[string]bool)

	for _, entry := range entries {
		// Clean up spacing variations in usernames to ensure a precise tally
		trimmedName := strings.TrimSpace(entry.Username)
		if trimmedName != "" {
			uniqueUsernames[trimmedName] = true
		}
	}

	// 5. Assign the number of unique active usernames found
	result.Attendants = len(uniqueUsernames)

	return result
}


