package static
import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/gorilla/websocket"
	"strings"
)

// Upgrader config to turn a normal HTTP connection into a WebSocket pipe
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { 
		return true // Allows your browser to talk to the server locally
	},
}

type Plate struct {
	RegNumber string `json:"regnumber"`
}

// Keep your working function exactly as it is!
func GenerateOptionsHTML() (string, error) {
	jsonPath := "static/registrationnumber.json"
 
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return "", fmt.Errorf("reading file %s: %w", jsonPath, err)
	}
 
	var plates []Plate
	if err := json.Unmarshal(data, &plates); err != nil {
		return "", fmt.Errorf("parsing JSON: %w", err)
	}
 
	result := ""
	for _, p := range plates {
		runes := []rune(p.RegNumber)
		length := len(runes)
		if length > 3 {
			length = 3
		}
		prefix := string(runes[:length])
	
		result += fmt.Sprintf("  <option value=\"%s\">\n", html.EscapeString(prefix))
	}
 
	return result, nil
}

func HandlePlateUpdatesWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(" WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	var lastSentData string

	for {
		// Call your function above to read the JSON file
		optionsHTML, err := GenerateOptionsHTML()
		if err != nil {
			log.Println("Error reading plates file inside WS:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Only send data if something actually changed in the file
		if optionsHTML != lastSentData {
			err = conn.WriteMessage(websocket.TextMessage, []byte(optionsHTML))
			if err != nil {
				break // Stops the loop if the user leaves or closes the tab
			}
			lastSentData = optionsHTML
		}

		// Poll the JSON file for changes every 2 seconds
		time.Sleep(2 * time.Second)
	}
}

func SearchPlates(w http.ResponseWriter, r *http.Request) {

    query := r.URL.Query().Get("q")

    data, err := os.ReadFile("static/registrationnumber.json")
    if err != nil {
        http.Error(w, "Could not read registration numbers", http.StatusInternalServerError)
        return
    }

    var plates []Plate

    if err := json.Unmarshal(data, &plates); err != nil {
        http.Error(w, "Could not read registration numbers", http.StatusInternalServerError)
        return
    }

    var results []string

    for _, plate := range plates {

        if strings.Contains(
            strings.ToLower(plate.RegNumber),
            strings.ToLower(query),
        ) {
            results = append(results, plate.RegNumber)
        }
    }

    w.Header().Set("Content-Type", "application/json")

    json.NewEncoder(w).Encode(results)
}
