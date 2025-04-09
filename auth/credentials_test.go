package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCredentialsSaveLoad(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "tmail-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testCredentialsPath := filepath.Join(tmpDir, "credentials.json")

	testCreds := Credentials{
		Email:       "test@example.com",
		AppPassword: "test-password-123",
		Name:        "Test User",
	}

	data, err := json.MarshalIndent(testCreds, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal credentials: %v", err)
	}

	err = os.WriteFile(testCredentialsPath, data, 0600)
	if err != nil {
		t.Fatalf("Failed to write credentials file: %v", err)
	}

	if _, err := os.Stat(testCredentialsPath); os.IsNotExist(err) {
		t.Fatalf("Credentials file was not created")
	}

	info, err := os.Stat(testCredentialsPath)
	if err != nil {
		t.Fatalf("Failed to stat credentials file: %v", err)
	}

	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", info.Mode().Perm())
	}

	readData, err := os.ReadFile(testCredentialsPath)
	if err != nil {
		t.Fatalf("Failed to read credentials file: %v", err)
	}

	var loadedCreds Credentials
	err = json.Unmarshal(readData, &loadedCreds)
	if err != nil {
		t.Fatalf("Failed to unmarshal credentials: %v", err)
	}

	if loadedCreds.Email != testCreds.Email {
		t.Errorf("Email mismatch: expected %s, got %s", testCreds.Email, loadedCreds.Email)
	}

	if loadedCreds.AppPassword != testCreds.AppPassword {
		t.Errorf("AppPassword mismatch: expected %s, got %s", testCreds.AppPassword, loadedCreds.AppPassword)
	}

	if loadedCreds.Name != testCreds.Name {
		t.Errorf("Name mismatch: expected %s, got %s", testCreds.Name, loadedCreds.Name)
	}

	err = os.Remove(testCredentialsPath)
	if err != nil {
		t.Fatalf("Failed to remove credentials file: %v", err)
	}

	if _, err := os.Stat(testCredentialsPath); !os.IsNotExist(err) {
		t.Errorf("Credentials file still exists after removal")
	}
}
