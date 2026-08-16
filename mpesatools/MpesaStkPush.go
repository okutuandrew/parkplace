package mpesatools

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	shortCode   = "174379"
	passkey     = "bfb279f9aa9bdbcf158e97dd71a467cd2e0c893059b10f78e6b72ada1ed2c919"
	testPhone   = "254711487030"
	callbackURL = "https://your-ngrok-url.ngrok-free.app/callback"
)

type MpesaConfig struct {
	ShortCode   string
	Passkey     string
	Phone       string // Customer phone (dynamic)
	Amount      string // Payment amount (dynamic)
	CallbackURL string
	Plate       string // Optional: for reference
}


// GetAccessToken fetches the M-Pesa OAuth token
func GetAccessToken(plate string,Phone string) (string, error) {
	key := os.Getenv("MPESA_CONSUMER_KEY")
	secret := os.Getenv("MPESA_CONSUMER_SECRET")
	if key == "" || secret == "" {
		return "", fmt.Errorf("set MPESA_CONSUMER_KEY and MPESA_CONSUMER_SECRET")
	}

	creds := base64.StdEncoding.EncodeToString([]byte(key + ":" + secret))

	req, _ := http.NewRequest("GET", "https://sandbox.safaricom.co.ke/oauth/v1/generate?grant_type=client_credentials", nil)
	req.Header.Set("Authorization", "Basic "+creds)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

// StkPush initiates the STK Push request
func StkPush(accessToken string) error {
	timestamp := time.Now().Format("20060102150405")
	password := base64.StdEncoding.EncodeToString([]byte(shortCode + passkey + timestamp))

	payload := map[string]string{
		"BusinessShortCode": shortCode,
		"Password":          password,
		"Timestamp":         timestamp,
		"TransactionType":   "CustomerPayBillOnline",
		"Amount":            "1",
		"PartyA":            testPhone,
		"PartyB":            shortCode,
		"PhoneNumber":       testPhone,
		"CallBackURL":       callbackURL,
		"AccountReference":  "DEMO",
		"TransactionDesc":   "Test payment",
	}

	jsonBody, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST",
		"https://sandbox.safaricom.co.ke/mpesa/stkpush/v1/processrequest",
		strings.NewReader(string(jsonBody)))

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response:", string(body))
	return nil
}