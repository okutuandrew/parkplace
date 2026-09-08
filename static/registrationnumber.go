package static

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
)


type Plate struct {
	RegNumber string `json:"regnumber"`
}


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
	
		result += fmt.Sprintf("  <option value=\"%s\">\n", html.EscapeString(p.RegNumber))
	}
 
	return result, nil
}