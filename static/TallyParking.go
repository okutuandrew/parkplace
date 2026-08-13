package static

 import (
	"fmt"
	"log"
 )



 type ParkingTally struct {
    Total     int
    Available int
    Full      int
    Warning   int
}


 func ParkingData() {

	  result := ParkingTally{}

	  file, err := os.ReadFile("static/parkingData.json")

	   if err != nil {
		log.Println("❌ Can't read parkingData.json:", err)
        
    }

	fmt.Println("TALLY START ")
 }