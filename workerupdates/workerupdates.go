package workerupdates

import ( 

	"encoding/json"
	"os"
	"sync"
)

type Updates struct {
    Lat     float64 `json:"lat"`
    Long    float64 `json:"lng"`
    Title   string  `json:"title"`
    Color   string  `json:"color"`
    Content string  `json:"content"`
} 
func Workerdata()  string {

	var  workerupdates Updates
	workerupdates.Lat =  -0.06575
	workerupdates.Long = 34.77502
	workerupdates.Title =  "Entrance A"
	workerupdates.Color = "color"
	workerupdates.Content = "🚗 Available: 5 spots\n🅿️ Near main gate"

	WorkerupdatesPointer := &workerupdates

	jsonData, _ := json.Marshal( workerupdates  )

	ScribeUpdates(  string(jsonData),  WorkerupdatesPointer  )
    return  string(jsonData)
}



var mu sync.Mutex 


func ScribeUpdates(jsonstring string, Newdata *Updates) {
	

	mu.Lock()

	defer mu.Unlock()

	filePath := "static/parkingData.json"

	// 1. Read existing file
	oldData, err := os.ReadFile(filePath)
	if err != nil {
		oldData = []byte("[]") // if file doesn't exist
	}

	// 2. Convert to slice
	var parking []Updates
	json.Unmarshal(oldData, &parking)

	// 3. Append new data
	parking = append(parking, *Newdata)

	// 4. Convert back to JSON
	finalData, _ := json.MarshalIndent(parking, "", "  ")

	// 5. Write full file (overwrite)
	os.WriteFile(filePath, finalData, 0666)
}