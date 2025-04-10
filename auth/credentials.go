package auth

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"golang.org/x/term"
)

// TODO: Create an interface for auth that allows for multiple provider authentication.

// Credentials stores user authentication information
type Credentials struct {
	Email       string `json:"email"`
	AppPassword string `json:"app_password"`
	Name        string `json:"name"`
}

var isAuthed = false

// PromptForAuthentication prompts the user for their credentials
func PromptForAuthentication() error {
	var creds Credentials
	reader := bufio.NewReader(os.Stdin)

	// Get name
	fmt.Print("Hi! Welcome to tmail, what's your name? ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read name: %v", err)
	}
	creds.Name = strings.TrimSpace(name)

	// Get email
	fmt.Printf("Thanks %s! What is your email address? ", creds.Name)
	email, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read email: %v", err)
	}
	creds.Email = strings.TrimSpace(email)

	// Get app password
	fmt.Print("What is your Gmail app password? If you haven't set one up, please see the README for instructions. ")
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to read password: %v", err)
	}
	creds.AppPassword = strings.TrimSpace(string(bytePassword))
	// Confirm the information
	fmt.Printf("\nName: %s\nEmail: %s\nPassword: (hidden)\n", creds.Name, creds.Email)
	fmt.Print("Is this information correct? (y/n): ")
	confirm, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read confirmation: %v", err)
	}

	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(confirm)), "y") {
		return fmt.Errorf("user canceled setup")
	}

	// Save credentials
	if err := SaveCredentials(creds); err != nil {
		return fmt.Errorf("failed to save credentials: %v", err)
	}

	fmt.Println("Authentication setup complete!")
	return nil
}

// getConfigDir returns the directory where configuration is stored
func getConfigDir() (string, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to find home directory: %v", err)
	}

	// Create path for the config directory using XDG Base Directory spec
	configDir := filepath.Join(homeDir, ".config", "tmail")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", fmt.Errorf("unable to create config directory: %v", err)
	}

	return configDir, nil
}

// getCredentialsPath returns the path to the credentials file
func getCredentialsPath() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "credentials.json"), nil
}

// SaveCredentials stores user credentials to the config file
func SaveCredentials(creds Credentials) error {
	path, err := getCredentialsPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal credentials: %v", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("unable to write credentials file: %v", err)
	}

	isAuthed = true

	return nil
}

// LoadCredentials reads user credentials from the config file
func LoadUser() (Credentials, error) {
	path, err := getCredentialsPath()
	var creds Credentials
	if isAuthed == true {
		return creds, nil
	}
	if err != nil {
		return Credentials{}, err
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return Credentials{}, fmt.Errorf("unable to read credentials file: %v", err)
	}

	// Unmarshal JSON
	if err := json.Unmarshal(data, &creds); err != nil {
		return Credentials{}, fmt.Errorf("unable to parse credentials file: %v", err)
	}

	return creds, nil
}

// RemoveCredentials deletes the credentials file
func RemoveCredentials() error {
	if isAuthed == false {
		return nil
	}
	path, err := getCredentialsPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to remove credentials file: %v", err)
	}
	isAuthed = false
	return nil
}
