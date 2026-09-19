package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"parkplace/Middlewares"
	"parkplace/WbSocks"
	"parkplace/bots"
	"parkplace/logs"
	"parkplace/mpesatools"
	"parkplace/static"
	"parkplace/workerupdates"
    "os"
	"github.com/joho/godotenv"
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
	Street      string
	PlateOptions template.HTML
}

type FormData struct {
	Username     string
	Password     string
	TotalParking int
	Spaces       int
	Attendants   int
}

type PaymentDetails struct {
	CarReg      string
	PhoneNumber string
}


type AttendantEntry struct {
	Username  string  `json:"username"`
	DateTime  string  `json:"date_time"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}


func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env vars")
	}

	// ==========================================
	// CREATE NEW PARKING DATA
	// ==========================================
	go func() {
		for {
			bots.NewParkingSpace()
			time.Sleep(3 * time.Second)
		}
	}()

	// ==========================================
	// DELETE RANDOM PARKING ENTRY
	// ==========================================
	go func() {
		for {
			bots.DeleteRandomEntry()
			time.Sleep(9 * time.Second)
		}
	}()

	// ==========================================
	// RESET PARKING DATA
	// ==========================================
	go func() {
		for {
			bots.ResetParkingData()
			time.Sleep(180 * time.Second)
		}
	}()

	// ==========================================
	// INITIAL STARTUP FUNCTIONS
	// ==========================================
	bots.NewParkingSpace()

	static.ParkingData()
	static.GenerateOptionsHTML()

	logs.SysLogs()

	log.Println(workerupdates.Workerdata())

	// ==========================================
	// ROUTES
	// ==========================================

	http.HandleFunc("/", LandingPage)

	// Static files
	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("./static")),
		),
	)

	// Front Map
	http.Handle(
		"/F-MAP",
		Middlewares.SessionTracker(
			http.HandlerFunc(
				checkCookie(FrontMapHandler),
			),
		),
	)

	// Worker Map
	http.HandleFunc(
		"/W-MAP",
		checkCookie(WorkerMapHandler),
	)

	// Parking Updates
	http.HandleFunc(
		"/PARKINGUPDATES",
		checkCookie(ParkingUpdates),
	)

	// Driver Login
	http.HandleFunc(
		"/DRIVERLOGGIN",
		DriverLoggin,
	)

	// Worker Login
	http.HandleFunc(
		"/WORKERLOGGIN",
		WorkerLoggin,
	)

	// WebSockets
	http.HandleFunc(
		"/ws",
		checkCookie(WbSocks.DriverMapWbScock),
	)

	http.HandleFunc(
		"/wss",
		checkCookie(WbSocks.DriverMapWbScock),
	)

	// Plate WebSockets
	http.HandleFunc(
		"/ws/plates",
		checkCookie(static.HandlePlateUpdatesWS),
	)

	http.HandleFunc(
		"/wss/plates",
		checkCookie(static.HandlePlateUpdatesWS),
	)

	// Plate API
	http.HandleFunc(
		"/api/plates",
		checkCookie(static.SearchPlates),
	)

	// M-Pesa
	http.Handle(
		"/MpesaPayment",
		Middlewares.SessionTracker(
			http.HandlerFunc(
				checkCookie(MpesaPaymentHandler),
			),
		),
	)

	// Logout
	http.HandleFunc(
		"/logout/",
		logoutHandler,
	)

	// Dashboard
	http.HandleFunc(
		"/DASHBOARD",
		checkCookie(Dashboard),
	)

	// Book Parking
	http.HandleFunc(
		"/BOOKPARKING",
		checkCookie(Bookparking),
	)

	fmt.Println("🚀 Server running on http://localhost:8080")
	fmt.Println("   Open → http://localhost:8080/F-MAP")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ======================================================
// LOGOUT
// ======================================================

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	username := strings.TrimPrefix(r.URL.Path, "/logout/")

	fmt.Println("Logging out:", username)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ======================================================
// LANDING PAGE
// ======================================================

func LandingPage(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: "driver-logged-in",
		Path:  "/",
	})

	tmpl, err := template.ParseFiles("forms/landingpage.html")
	if err != nil {
		http.Error(
			w,
			"Template error: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	log.Println("User accessed the main Landing Login Page Portal")

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(
			w,
			"Execute error: "+err.Error(),
			http.StatusInternalServerError,
		)
	}
}

// ======================================================
// MPESA PAYMENT
// ======================================================

func MpesaPaymentHandler(w http.ResponseWriter, r *http.Request) {

	var PayReff PaymentDetails

	PayReff.CarReg = r.FormValue("plate")
	PayReff.PhoneNumber = r.FormValue("phone")

	token, err := mpesatools.GetAccessToken(
		PayReff.CarReg,
		PayReff.PhoneNumber,
	)

	if err != nil {
		fmt.Println("Token error:", err)
		return
	}

	fmt.Println("Access Token:", token[:20]+"...")

	mpesatools.StkPush(token)
}

// ======================================================
// DASHBOARD
// ======================================================

func Dashboard(w http.ResponseWriter, r *http.Request) {

	// Create the data that will be sent to dashboard.html
	data := FormData{}

	// ==================================================
	// GET PARKING DATA
	// ==================================================

	result := static.ParkingData()

	// Total parking records
	data.TotalParking = result.Total

	// Total available spaces
	data.Spaces = result.Spaces

	// ==================================================
	// GET ATTENDANTS
	// ==================================================

	data.Attendants = result.Attendants

	// ==================================================
	// HANDLE DASHBOARD LOGIN
	// ==================================================

	if r.Method == "POST" {

		username := r.FormValue("username")
		password := r.FormValue("password")

		log.Println("Dashboard login attempt:", username)

		// Password is captured but not being checked yet
		_ = password

		// ==================================================
		// CREATE ATTENDANT RECORD
		// ==================================================

		attendant := AttendantEntry{
			Username:  username,
			DateTime:  time.Now().UTC().Format(time.RFC3339),
			Longitude: 34.7679,
			Latitude:  -0.0917,
		}

		// ==================================================
		// READ EXISTING TallyAttendants.json
		// ==================================================

		filePath := "static/TallyAttendants.json"

		file, err := os.ReadFile(filePath)
		if err != nil {
			log.Println("Can't read TallyAttendants.json:", err)
			http.Error(w, "Unable to read attendant data", http.StatusInternalServerError)
			return
		}

		var attendants []AttendantEntry

		err = json.Unmarshal(file, &attendants)
		if err != nil {
			log.Println("Can't parse TallyAttendants.json:", err)
			http.Error(w, "Unable to parse attendant data", http.StatusInternalServerError)
			return
		}

		// ==================================================
		// APPEND NEW ATTENDANT
		// ==================================================

		attendants = append(attendants, attendant)

		// ==================================================
		// WRITE UPDATED JSON BACK TO FILE
		// ==================================================

		updatedFile, err := json.MarshalIndent(attendants, "", "  ")
		if err != nil {
			log.Println("Can't encode attendant data:", err)
			http.Error(w, "Unable to encode attendant data", http.StatusInternalServerError)
			return
		}

		err = os.WriteFile(filePath, updatedFile, 0644)
		if err != nil {
			log.Println("Can't write TallyAttendants.json:", err)
			http.Error(w, "Unable to save attendant data", http.StatusInternalServerError)
			return
		}

		log.Println("Attendant added:", attendant)
	}

	// ==================================================
	// LOAD DASHBOARD TEMPLATE
	// ==================================================

	tmpl, err := template.ParseFiles("forms/dashboard.html")
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("User accessed Dashboard")

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Execute error: "+err.Error(), http.StatusInternalServerError)
	}
}

// ======================================================
// WORKER MAP
// ======================================================

func WorkerMapHandler(w http.ResponseWriter, r *http.Request) {

	optionsHTML, err := static.GenerateOptionsHTML()

	if err != nil {
		http.Error(
			w,
			"Plate options error: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	data := PageData{
		Street:       "KIMATHI STREET",
		PlateOptions: template.HTML(optionsHTML),
	}

	tmpl, err := template.ParseFiles("maps/workermap.html")

	if err != nil {
		http.Error(
			w,
			"Template error: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	log.Println("User accessed W-MAP page")

	err = tmpl.Execute(w, data)

	if err != nil {
		http.Error(
			w,
			"Execute error: "+err.Error(),
			http.StatusInternalServerError,
		)
	}
}

// ======================================================
// FRONT MAP
// ======================================================

func FrontMapHandler(w http.ResponseWriter, r *http.Request) {

	data := PageData{
		Street: "KIMATHI STREET",
	}

	tmpl, err := template.ParseFiles("maps/map.html")

	if err != nil {
		http.Error(
			w,
			"Template error: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	log.Println("User accessed F-MAP page")

	err = tmpl.Execute(w, data)

	if err != nil {
		http.Error(
			w,
			"Execute error: "+err.Error(),
			http.StatusInternalServerError,
		)
	}
}

// ======================================================
// BOOK PARKING
// ======================================================

func Bookparking(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(
		w,
		"Hello, this is your Go HTTP server!",
	)
}

// ======================================================
// PARKING UPDATES
// ======================================================

func ParkingUpdates(w http.ResponseWriter, r *http.Request) {

	PostSpace := workerupdates.Updates{}

	PostSpace.Lat, _ = strconv.ParseFloat(
		r.FormValue("lattitude"),
		64,
	)

	PostSpace.Long, _ = strconv.ParseFloat(
		r.FormValue("longitude"),
		64,
	)

	PostSpace.Color = r.FormValue("color")

	PostSpace.Content = r.FormValue("notes")

	PostSpace.Spaces, _ = strconv.Atoi(
		r.FormValue("spaces"),
	)

	WorkerupdatesPointer := &PostSpace

	jsonData, _ := json.Marshal(
		WorkerupdatesPointer,
	)

	workerupdates.ScribeUpdates(
		string(jsonData),
		WorkerupdatesPointer,
	)

	log.Println(
		"Worker Posted updates",
		PostSpace,
	)

	data := PageData{
		Street: "KIMATHI STREET",
	}

	tmpl, err := template.ParseFiles(
		"maps/workermap.html",
	)

	if err != nil {
		http.Error(
			w,
			"Template error: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		http.Error(
			w,
			"Execute error: "+err.Error(),
			http.StatusInternalServerError,
		)
	}
}

// ======================================================
// DRIVER LOGIN
// ======================================================

func DriverLoggin(w http.ResponseWriter, r *http.Request) {

	log.Println("DRIVER OPENED APP")

	tmpl, err := template.ParseFiles(
		"forms/landingpage.html",
	)

	if err != nil {
		http.Error(
			w,
			"Template error: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	err = tmpl.Execute(w, nil)

	if err != nil {
		http.Error(
			w,
			"Execute error: "+err.Error(),
			http.StatusInternalServerError,
		)
	}
}

// ======================================================
// WORKER LOGIN
// ======================================================

func WorkerLoggin(w http.ResponseWriter, r *http.Request) {

	log.Println("Worker OPENED APP")

	tmpl, err := template.ParseFiles(
		"forms/workerloggin.html",
	)

	if err != nil {
		http.Error(
			w,
			"Template error: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	err = tmpl.Execute(w, nil)

	if err != nil {
		http.Error(
			w,
			"Execute error: "+err.Error(),
			http.StatusInternalServerError,
		)
	}
}