package auth

import (
	"bufio"
	"bytes"
	// "context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

// Google OAuth2 Endpoints for Device Flow
const (
	deviceAuthURL  = "https://oauth2.googleapis.com/device/code"
	tokenURL       = "https://oauth2.googleapis.com/token"
	clientID       = "563445405127-d81itr1efp6dg7agrn2de8d3rim3s6m3.apps.googleusercontent.com" // Replace with your actual client ID
	clientSecret   = "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"
)

// DeviceAuthResponse represents the response from Google's device authorization request
type DeviceAuthResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
	Interval        int    `json:"interval"`
	ExpiresIn       int    `json:"expires_in"`
}

// TokenResponse represents the OAuth2 token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// Ask the user if they want to authenticate with Google
func PromptForAuthentication() {
	fmt.Print("Would you like to authenticate with Google? (Y/N): ")

	// Read user input
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input)) // Normalize input

	// If user says yes, start authentication flow
	if input == "y" || input == "yes" {
		fmt.Println("Starting Google authentication...")
		token, err := Authenticate()
		if err != nil {
			fmt.Println("Authentication failed:", err)
			return
		}

		fmt.Println("Successfully authenticated! Access Token:", token)
	} else {
		fmt.Println("Skipping authentication.")
	}
}

// Authenticate performs the Google OAuth2 Device Authorization Flow
func Authenticate() (string, error) {
	// Step 1: Request a Device Code
	deviceResp, err := requestDeviceCode()
	if err != nil {
		return "", err
	}

	// Step 2: Display Code to User
	fmt.Println("Go to this URL and enter the code to authenticate:")
	fmt.Println(deviceResp.VerificationURL)
	fmt.Println("User Code:", deviceResp.UserCode)

	// Step 3: Poll for Authorization
	token, err := pollForToken(deviceResp)
	if err != nil {
		return "", err
	}

	// Step 4: Save token
	saveToken(token)
	fmt.Println("Authentication successful!")
	return "", nil
}

// requestDeviceCode makes a request to Google's device authorization endpoint
func requestDeviceCode() (*DeviceAuthResponse, error) {
	requestBody := map[string]string{
		"client_id": clientID,
		"scope":     "openid email",
	}
	jsonBody, _ := json.Marshal(requestBody)

	resp, err := http.Post(deviceAuthURL, "application/json", bytes.NewBuffer(jsonBody))
	
	if err != nil {
		return nil, fmt.Errorf("failed to request device code: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var deviceResp DeviceAuthResponse
	if err := json.Unmarshal(body, &deviceResp); err != nil {
		return nil, fmt.Errorf("failed to parse device response: %v", err)
	}
	return &deviceResp, nil
}

// pollForToken continuously checks if the user has completed authentication
func pollForToken(deviceResp *DeviceAuthResponse) (*oauth2.Token, error) {
	expirationTime := time.Now().Add(time.Duration(deviceResp.ExpiresIn) * time.Second)
	ticker := time.NewTicker(time.Duration(deviceResp.Interval) * time.Second)
	defer ticker.Stop()

	for {
		// Check if time has expired
		if time.Now().After(expirationTime) {
			return nil, fmt.Errorf("authentication timed out: user did not complete login in time")
		}

		// Request the token
		tokenResp, err := requestToken(deviceResp.DeviceCode)
		if err == nil {
			return tokenResp, nil // Successfully got a token!
		}


		// Handle common OAuth errors
		if err.Error() == "authorization_pending" {
			// Normal case, keep waiting
			fmt.Println("Waiting for user to authenticate...")
		} else if err.Error() == "expired_token" {
			return nil, fmt.Errorf("device code expired, please restart authentication")
		} else if err.Error() == "slow_down" {
			// Google is asking us to slow down requests
			ticker.Reset(10 * time.Second) // Increase interval
			fmt.Println("Google requested slower polling, adjusting...")
		} else {
			// Some other error
			fmt.Errorf("authentication error: %v", err.Error())
		}

		// Wait for the next polling interval
		<-ticker.C
	}
}

// requestToken exchanges a device code for an OAuth2 token
func requestToken(deviceCode string) (*oauth2.Token, error) {
	requestBody := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"device_code":   deviceCode,
		"grant_type":    "urn:ietf:params:oauth:grant-type:device_code",
	}
	jsonBody, _ := json.Marshal(requestBody)

	resp, err := http.Post(tokenURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to request token: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed: %s", body)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %v", err)
	}

	return &oauth2.Token{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		Expiry:       time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		TokenType:    tokenResp.TokenType,
	}, nil
}

// saveToken stores the OAuth token in a local file
func saveToken(token *oauth2.Token) {
	file, err := os.Create("token.json")
	if err != nil {
		fmt.Println("Unable to save token:", err)
		return
	}
	defer file.Close()

	json.NewEncoder(file).Encode(token)
}
