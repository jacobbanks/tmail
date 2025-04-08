package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)


// TODO: Create an interface for auth that allows for multiple provider authentication.

// Credentials stores user authentication information
type Credentials struct {
	Email       string `json:"email"`
	AppPassword string `json:"app_password"`
	Name        string `json:"name"`
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

	return nil
}

// LoadCredentials reads user credentials from the config file
func LoadCredentials() (Credentials, error) {
	path, err := getCredentialsPath()
	if err != nil {
		return Credentials{}, err
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return Credentials{}, fmt.Errorf("unable to read credentials file: %v", err)
	}

	// Unmarshal JSON
	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return Credentials{}, fmt.Errorf("unable to parse credentials file: %v", err)
	}

	return creds, nil
}

// RemoveCredentials deletes the credentials file
func RemoveCredentials() error {
	path, err := getCredentialsPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to remove credentials file: %v", err)
	}

	return nil
}
