package bots 

import (
"fmt"
"encoding/json"
"os"
"sync"
"parkplace/workerupdates" 
"math/rand"
"strconv"

)


var fileMutex sync.Mutex



func randomRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}




func  NewParkingSpace(){

var DummyParking workerupdates.Updates


	latMin := -0.06575 // a
	latMax := -0.08575 // b

	lngMin := 34.77502 // c
	lngMax := 34.7502 // d

colorOptions := []string{"#49d611ff", "#FABB05" , "#34A853" }
randomIndex := rand.Intn(len(colorOptions))

spacesOptions := []int{1,2,3,4,5,6}
randomSpaceIndex := rand.Intn(len( spacesOptions))


DummyParking.Lat = randomRange(latMin, latMax)
DummyParking.Long = randomRange(lngMin, lngMax)
DummyParking.Title = "TESTBOT";
DummyParking.Color = colorOptions[randomIndex] 
SpacesStr:= strconv.Itoa( spacesOptions[randomSpaceIndex]  )

DummyParking.Content = "🚗 Available: " + SpacesStr + " spots\n🅿️ Near main gate"
DummyParking.Spaces = spacesOptions[randomSpaceIndex]; 

WorkerupdatesPointerBot := &DummyParking

jsonData, _ := json.Marshal( WorkerupdatesPointerBot )


workerupdates.ScribeUpdates(  string(jsonData), WorkerupdatesPointerBot )


fmt.Println("WORKER BOT ");

}



func DeleteRandomEntry() {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	filePath := "static/parkingData.json"

	// 1. Read the current JSON data from disk
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Delete error: Can't read file:", err)
		return
	}

	// 2. Unmarshal it into a slice of your workerupdates struct
	var entries []workerupdates.Updates
	err = json.Unmarshal(fileData, &entries)
	if err != nil {
		fmt.Println("Delete error: Can't parse JSON:", err)
		return
	}

	// 3. CRUCIAL: If the file is already empty, stop here so the app doesn't crash!
	if len(entries) == 0 {
		fmt.Println("ℹNo entries left to delete.")
		return
	}

	// 4. Pick a random index position between 0 and the total length - 1
	randomIndex := rand.Intn(len(entries))
	deletedEntry := entries[randomIndex]

	// 5. Remove the item at that index using standard Go slice manipulation
	// This takes everything before the index and appends everything after the index
	entries = append(entries[:randomIndex], entries[randomIndex+1:]...)

	// 6. Encode the updated slice back into formatted text JSON
	updatedJSON, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		fmt.Println("Delete error: Can't pack JSON:", err)
		return
	}

	// 7. Overwrite the file with the remaining entries
	err = os.WriteFile(filePath, updatedJSON, 0666)
	if err != nil {
		fmt.Println("Delete error: Can't write back to file:", err)
		return
	}

	fmt.Printf("🗑️ Randomly deleted entry -> Title: %s | Remaining spots left: %d\n", 
		deletedEntry.Title, len(entries))
}


