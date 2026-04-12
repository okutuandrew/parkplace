package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "parkplace/workerupdates"
)


type PageData struct {
    Street string
}

func main() {

    log.Println(workerupdates.Workerdata())
    // Static files handler - MUST come BEFORE the specific routes
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

    // Your page route
    http.HandleFunc("/F-MAP", FrontMapHandler)
	http.HandleFunc("/W-MAP",WorkerMapHandler)
	http.HandleFunc("/PARKINGUPDATES",ParkingUpdates)
	http.HandleFunc("/DRIVERLOGGIN",DriverLoggin)
	http.HandleFunc("/WORKERLOGGIN",WorkerLoggin)

    // Optional route
    http.HandleFunc("/BOOKPARKING", Bookparking)

    fmt.Println("🚀 Server running on http://localhost:8080")
    fmt.Println("   Open → http://localhost:8080/F-MAP")
    log.Fatal(http.ListenAndServe(":8080", nil))


  
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
   

	log.Println("✅ Worker Posted updates")

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

	tmpl, err := template.ParseFiles("forms/driverloggin.html")
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

