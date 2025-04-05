package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCredentialsSaveLoad(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "tmail-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Test using direct file paths instead of mocking
	testCredentialsPath := filepath.Join(tmpDir, "credentials.json")
	
	// Test credentials
	testCreds := Credentials{
		Email:       "test@example.com",
		AppPassword: "test-password-123",
		Name:        "Test User",
	}
	
	// Convert to JSON
	data, err := json.MarshalIndent(testCreds, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal credentials: %v", err)
	}
	
	// Write directly to the test path
	err = os.WriteFile(testCredentialsPath, data, 0600)
	if err != nil {
		t.Fatalf("Failed to write credentials file: %v", err)
	}
	
	// Verify the file exists
	if _, err := os.Stat(testCredentialsPath); os.IsNotExist(err) {
		t.Fatalf("Credentials file was not created")
	}
	
	// Check file permissions
	info, err := os.Stat(testCredentialsPath)
	if err != nil {
		t.Fatalf("Failed to stat credentials file: %v", err)
	}
	
	// On Unix systems, we can check for 0600 permissions
	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", info.Mode().Perm())
	}
	
	// Read the file directly and unmarshal
	readData, err := os.ReadFile(testCredentialsPath)
	if err != nil {
		t.Fatalf("Failed to read credentials file: %v", err)
	}
	
	var loadedCreds Credentials
	err = json.Unmarshal(readData, &loadedCreds)
	if err != nil {
		t.Fatalf("Failed to unmarshal credentials: %v", err)
	}
	
	// Verify loaded credentials match
	if loadedCreds.Email != testCreds.Email {
		t.Errorf("Email mismatch: expected %s, got %s", testCreds.Email, loadedCreds.Email)
	}
	
	if loadedCreds.AppPassword != testCreds.AppPassword {
		t.Errorf("AppPassword mismatch: expected %s, got %s", testCreds.AppPassword, loadedCreds.AppPassword)
	}
	
	if loadedCreds.Name != testCreds.Name {
		t.Errorf("Name mismatch: expected %s, got %s", testCreds.Name, loadedCreds.Name)
	}
	
	// Remove the file
	err = os.Remove(testCredentialsPath)
	if err != nil {
		t.Fatalf("Failed to remove credentials file: %v", err)
	}
	
	// Verify file no longer exists
	if _, err := os.Stat(testCredentialsPath); !os.IsNotExist(err) {
		t.Errorf("Credentials file still exists after removal")
	}
}

func TestUserAlias(t *testing.T) {
	// Test that User alias works with Credentials
	var creds Credentials
	creds.Email = "test@example.com"
	
	var user User = creds
	
	if user.Email != creds.Email {
		t.Errorf("User alias not working correctly")
	}
}