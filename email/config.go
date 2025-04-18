package email

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds email provider settings
type Config struct {
	SMTPHost string `json:"smtp_host"`
	SMTPPort string `json:"smtp_port"`
	IMAPHost string `json:"imap_host"`
	IMAPPort string `json:"imap_port"`
}

// UserConfig holds user preferences
type UserConfig struct {
	Theme           string `json:"theme"`         // UI color theme
	DefaultNumMails int    `json:"default_mails"` // Number of emails to fetch
}

// DefaultConfig provides standard connection settings for Gmail's SMTP and IMAP servers.
var DefaultConfig = Config{
	SMTPHost: "smtp.gmail.com",
	SMTPPort: "587",
	IMAPHost: "imap.gmail.com",
	IMAPPort: "993",
}

// DefaultUserConfig provides default UI and behavior settings for the application.
var DefaultUserConfig = UserConfig{
	Theme:           "blue",
	DefaultNumMails: 50,
}

// GetSMTPAddress returns the complete SMTP server address with port for email sending.
func (c *Config) GetSMTPAddress() string {
	return c.SMTPHost + ":" + c.SMTPPort
}

// GetIMAPAddress returns the complete IMAP server address with port for email fetching.
func (c *Config) GetIMAPAddress() string {
	return c.IMAPHost + ":" + c.IMAPPort
}

// GetConfigDir returns the configuration directory
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}
	configDir := filepath.Join(homeDir, ".config", "tmail")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create config directory: %v", err)
	}
	return configDir, nil
}

// SaveUserConfig saves user preferences to config file
func SaveUserConfig(config UserConfig) error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// LoadUserConfig loads user preferences from config file
func LoadUserConfig() (UserConfig, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return DefaultUserConfig, err
	}

	configPath := filepath.Join(configDir, "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If file doesn't exist, save and return defaults
			SaveUserConfig(DefaultUserConfig)
			return DefaultUserConfig, nil
		}
		return DefaultUserConfig, fmt.Errorf("failed to read config file: %v", err)
	}

	var config UserConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return DefaultUserConfig, fmt.Errorf("failed to parse config file: %v", err)
	}

	return config, nil
}
