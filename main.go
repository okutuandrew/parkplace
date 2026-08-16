package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "parkplace/workerupdates"
     "parkplace/logs"
     "encoding/json"
     "strconv"
     "parkplace/Middlewares"
     "parkplace/WbSocks"
     "parkplace/static"
    "parkplace/bots"
    "time"

)



func checkCookie(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("session")
        if err != nil || cookie.Value == "" {
            http.Redirect(w, r, "/", http.StatusSeeOther)
            return
        }
        next(w, r)
    }
}




type PageData struct {
    Street string
}



type FormData struct {  
    Username string
    Password string
    TotalParking int
    Spaces   int
}



func main() {



    go func() {
        for {
            // Call your bot function to generate a random location
            bots.NewParkingSpace()

            

            // Pause the execution of this specific loop for exactly 10 seconds
            time.Sleep(3 * time.Second)
        }
    }() 


       go func() {
        for {
            // Call your bot function to generate a random location
            bots.DeleteRandomEntry()
            

            // Pause the execution of this specific loop for exactly 10 seconds
            time.Sleep(9 * time.Second)
        }
    }() 




       go func() {
        for {
            // Call your bot function to generate a random location
            bots.ResetParkingData()
            

            // Pause the execution of this specific loop for exactly 10 seconds
            time.Sleep(180 * time.Second)
        }
    }() 


    


     



    bots.NewParkingSpace()
    static.ParkingData(); 
    logs.SysLogs()
    log.Println(workerupdates.Workerdata())
    // Static files handler - MUST come BEFORE the specific routes

    http.HandleFunc("/",LandingPage)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
    http.Handle("/F-MAP",Middlewares.SessionTracker(http.HandlerFunc(   checkCookie(FrontMapHandler ) )))
    // Your page route
    //http.HandleFunc("/F-MAP", FrontMapHandler)
	http.HandleFunc("/W-MAP", checkCookie(WorkerMapHandler))
	http.HandleFunc("/PARKINGUPDATES", checkCookie(ParkingUpdates))
	http.HandleFunc("/DRIVERLOGGIN",DriverLoggin)
	http.HandleFunc("/WORKERLOGGIN",WorkerLoggin)
   http.HandleFunc("/ws", checkCookie(WbSocks.DriverMapWbScock))


    http.HandleFunc("/DASHBOARD", checkCookie(Dashboard))
    // Optional route
    http.HandleFunc("/BOOKPARKING", checkCookie(Bookparking))

    fmt.Println("🚀 Server running on http://localhost:8080")
    fmt.Println("   Open → http://localhost:8080/F-MAP")
    log.Fatal(http.ListenAndServe(":8080", nil))


  
}


func LandingPage(w http.ResponseWriter, r *http.Request) {
    

       // add  cookies 
     http.SetCookie(w, &http.Cookie{
        Name:  "session",
        Value: "driver-logged-in",
        Path:  "/",
    })

    tmpl, err := template.ParseFiles("forms/landingpage.html")
    if err != nil {
        http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // 2. Fix your log to accurately match the action
    log.Println("✅ User accessed the main Landing Login Page Portal")

    // 3. CRUCIAL: Execute the template to send the HTML down to the browser window
    err = tmpl.Execute(w, nil) // passing nil since this form doesn't need dynamic Go struct data
    if err != nil {
        http.Error(w, "Execute error: "+err.Error(), http.StatusInternalServerError)
    }
}



func Dashboard(w http.ResponseWriter, r *http.Request) {

     data := FormData{}

    

	result := static.ParkingData()
	data.TotalParking = result.Total
    data.Spaces   = result.Spaces


     if r.Method == "POST" {
        data.Username = r.FormValue("username")  
        log.Println("Dashboard login attempt:", data )
     }

    tmpl, err := template.ParseFiles("forms/dashboard.html")
    if err != nil {
        http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // 2. Fix your log to accurately match the action
    log.Println("✅ User accessed the main Landing Login Page Portal")

    // 3. CRUCIAL: Execute the template to send the HTML down to the browser window
    err = tmpl.Execute(w, data) 
    if err != nil {
        http.Error(w, "Execute error: "+err.Error(), http.StatusInternalServerError)
    }
}



func WorkerMapHandler(w http.ResponseWriter, r *http.Request) {
    data := PageData{Street: "KIMATHI STREET"}

    tmpl, err := template.ParseFiles("maps/workermap.html")
    if err != nil {
        http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    log.Println("✅ User accessed W-MAP page")
    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, "Execute error: "+err.Error(), http.StatusInternalServerError)
    }
}



func FrontMapHandler(w http.ResponseWriter, r *http.Request) {
    data := PageData{Street: "KIMATHI STREET"}

    tmpl, err := template.ParseFiles("maps/map.html")
    if err != nil {
        http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
        return
    }

    log.Println("✅ User accessed F-MAP page")
    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, "Execute error: "+err.Error(), http.StatusInternalServerError)
    }
}

func Bookparking(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, this is your Go HTTP server!")
}

func ParkingUpdates(w http.ResponseWriter, r *http.Request) {



PostSpace := workerupdates.Updates{}  

   PostSpace.Lat,_ =  strconv.ParseFloat(r.FormValue("lattitude"), 64)
   PostSpace.Long,_ =  strconv.ParseFloat(r.FormValue("longitude"), 64)
   PostSpace.Color = r.FormValue("color")
   PostSpace.Content = r.FormValue("notes")
   
   PostSpace.Spaces,_ = strconv.Atoi(r.FormValue("spaces") )   

   WorkerupdatesPointer := &PostSpace

	jsonData, _ := json.Marshal(  WorkerupdatesPointer  )

	workerupdates.ScribeUpdates(  string(jsonData),  WorkerupdatesPointer  )
   



	log.Println("✅ Worker Posted updates",PostSpace )

	  data := PageData{Street: "KIMATHI STREET"}

    tmpl, err := template.ParseFiles("maps/workermap.html")
    if err != nil {
        http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, "Execute error: "+err.Error(), http.StatusInternalServerError)
    }
}



func DriverLoggin(w http.ResponseWriter, r *http.Request) {

	log.Println("✅ DRIVER  OPENED APP ")

 

	tmpl, err := template.ParseFiles("forms/landingpage.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil) // 👈 no data passed
	if err != nil {
		http.Error(w, "Execute error: "+err.Error(), http.StatusInternalServerError)
	}
}



func WorkerLoggin(w http.ResponseWriter, r *http.Request) {

	log.Println("✅ Worker   OPENED APP ")

	tmpl, err := template.ParseFiles("forms/workerloggin.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil) // 👈 no data passed
	if err != nil {
		http.Error(w, "Execute error: "+err.Error(), http.StatusInternalServerError)
	}
}

