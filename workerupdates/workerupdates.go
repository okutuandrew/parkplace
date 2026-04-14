package workerupdates

import ( 

	"encoding/json"
)

type Updates struct {
    Lat     float64
    Long    float64
    Title   string
    Color   string
    Content string
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

func ScribeUpdates(jsonstring string, Newdata *Updates){

/*	Newdata.Lat = -0.06575
	Newdata.Long = 34.77502
	Newdata.Title =   "Entrance A"
	Newdata.Color =  "color"
	Newdata.Content = "🚗 Available: 5 spots\n🅿️ Near main gate"

	

	log.Println( Newdata   )

	OldData,_ :=  os.ReadFile("static/parkingData.json") 

	NewdataSlice  := []Updates{ *Newdata}
json.Unmarshal(OldData,&NewdataSlice)
AddData,_ := json.Marshal(NewdataSlice) 
os.WriteFile("static/parkingData.json",  AddData  , 0644)

// WITH this CORRECT code:
var parkingSpots []Updates
if err := json.Unmarshal(OldData, &parkingSpots); err != nil {
    parkingSpots = []Updates{}
}
parkingSpots = append(parkingSpots, *Newdata)
jsonBytes, _ := json.MarshalIndent(parkingSpots, "", "  ")
//os.WriteFile("static/parkingData.json", jsonBytes, 0644)
 */
	

}