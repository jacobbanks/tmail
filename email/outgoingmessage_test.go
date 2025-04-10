package email

import (
	"strings"
	"testing"
)

func TestSanitizeAddresses(t *testing.T) {
	// Test with empty list
	addresses := []string{}
	clean := sanitizeAddresses(addresses)
	if len(clean) != 0 {
		t.Errorf("Expected empty slice, got %v", clean)
	}

	// Test with whitespace
	addresses = []string{" test@example.com ", "  ", "other@example.com  "}
	clean = sanitizeAddresses(addresses)
	if len(clean) != 2 || clean[0] != "test@example.com" || clean[1] != "other@example.com" {
		t.Errorf("Expected cleaned addresses, got %v", clean)
	}
}

func TestValidateEmailMessage_NilMessage(t *testing.T) {
	err := validateEmailMessage(nil)
	if err == nil {
		t.Errorf("Expected error with nil message")
	}
	if !strings.Contains(err.Error(), "nil") {
		t.Errorf("Expected error message to mention 'nil', got %q", err.Error())
	}
}

// This test can be skipped when credentials aren't available
func TestNewOutgoingMessage(t *testing.T) {
	msg, err := NewOutgoingMessage()
	if err != nil {
		t.Skip("Skipping test due to missing credentials: " + err.Error())
	}

	if msg.From == "" {
		t.Errorf("Expected From to be set, got empty string")
	}

	if len(msg.To) != 0 {
		t.Errorf("Expected To to be empty, got %v", msg.To)
	}

	if len(msg.Cc) != 0 {
		t.Errorf("Expected Cc to be empty, got %v", msg.Cc)
	}

	if len(msg.Bcc) != 0 {
		t.Errorf("Expected Bcc to be empty, got %v", msg.Bcc)
	}

	if msg.Subject != "" {
		t.Errorf("Expected Subject to be empty, got %s", msg.Subject)
	}

	if len(msg.Text) > 0 {
		t.Errorf("Expected Body to be empty, got %s", msg.Text)
	}
}
