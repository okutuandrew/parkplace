package main

import  (  "parkplace/workerupdates"
            "os"
            "bufio"
            "fmt"
            "strings"
            "strconv"
            "encoding/json"
          
 )


 type Updates struct {
    Lat     float64 `json:"lat"`
    Long    float64 `json:"lng"`
    Title   string  `json:"title"`
    Color   string  `json:"color"`
    Content string  `json:"content"`
} 


func main() {
    workerupdates.Workerdata()
    ParkingDetails()
}

func ParkingDetails() {
	var update Updates

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("ENTER LAT: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	lat, _ := strconv.ParseFloat(input, 64)
	update.Lat = lat

	fmt.Print("ENTER LONG: ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)
	long, _ := strconv.ParseFloat(input, 64)
	update.Long = long

	update.Title = "Entrance A"
	update.Color = "color"
	update.Content = "🚗 Available: 5 spots\n🅿️ Near main gate"

	jsonData, _ := json.Marshal(update)

	fmt.Println(string(jsonData))

	ScribeUpdatesCli(&update)
}


func ScribeUpdatesCli(newData *Updates) {
	filePath := "static/parkingData.json"

	oldData, _ := os.ReadFile(filePath)

	var parking []Updates
	json.Unmarshal(oldData, &parking)

	parking = append(parking, *newData)

	finalData, _ := json.MarshalIndent(parking, "", "  ")

	os.WriteFile(filePath, finalData, 0666)
}