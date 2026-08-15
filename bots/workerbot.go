package bots 

import (
"fmt"
"encoding/json"
_"os"
_"sync"
"parkplace/workerupdates" 
"math/rand"
"strconv"

)



func randomRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}




func  NewParkingSpace(){

var DummyParking workerupdates.Updates


	latMin := -0.06575 // a
	latMax := -0.08575 // b

	lngMin := 34.77502 // c
	lngMax := 34.7502 // d

colorOptions := []string{"#ea4335", "#FABB05" }
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

