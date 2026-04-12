package workerupdates

import ( 

	"encoding/json"
	 "os"
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

	jsonData, _ := json.Marshal( workerupdates  )
    return  string(jsonData)
}

func ScribeUpdates(jsonstring string){

	

}