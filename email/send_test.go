package email

import (
	"testing"
)

func TestSanitizeHeader(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"Normal text", "Normal text"},
		{"Text\nwith\nnewlines", "Text with newlines"},
		{"Text\r\nwith\r\nCRLF", "Text with CRLF"},
		{"Text\rwith\rCR", "Text with CR"},
		{" Trim spaces ", "Trim spaces"},
	}

	for _, tc := range testCases {
		result := sanitizeHeader(tc.input)
		if result != tc.expected {
			t.Errorf("sanitizeHeader(%q) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

func TestSanitizeAddresses(t *testing.T) {
	addresses := []string{
		"user1@example.com",
		"  user2@example.com  ",
		"user3@example.com\r\n",
		"", // Empty address should be filtered out
	}

	expected := []string{
		"user1@example.com",
		"user2@example.com",
		"user3@example.com",
		// Empty string filtered out
	}

	result := sanitizeAddresses(addresses)

	if len(result) != len(expected) {
		t.Errorf("Expected %d addresses, got %d", len(expected), len(result))
		return
	}

	for i, addr := range result {
		if addr != expected[i] {
			t.Errorf("Address at index %d: got %q, expected %q", i, addr, expected[i])
		}
	}
}

func TestValidateEmailMessage(t *testing.T) {
	// Test nil message
	err := validateEmailMessage(nil)
	if err == nil {
		t.Errorf("Expected error for nil message")
	}

	// Test empty recipients
	emptyMsg := &EmailMessage{
		From:    "sender@example.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}
	err = validateEmailMessage(emptyMsg)
	if err == nil {
		t.Errorf("Expected error for message with no recipients")
	}

	// Test valid message with To
	validToMsg := &EmailMessage{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Test Subject",
		Body:    "Test Body",
	}
	err = validateEmailMessage(validToMsg)
	if err != nil {
		t.Errorf("Unexpected error for valid message with To: %v", err)
	}

	// Test valid message with Cc
	validCcMsg := &EmailMessage{
		From:    "sender@example.com",
		Cc:      []string{"recipient@example.com"},
		Subject: "Test Subject",
		Body:    "Test Body",
	}
	err = validateEmailMessage(validCcMsg)
	if err != nil {
		t.Errorf("Unexpected error for valid message with Cc: %v", err)
	}

	// Test valid message with Bcc
	validBccMsg := &EmailMessage{
		From:    "sender@example.com",
		Bcc:     []string{"recipient@example.com"},
		Subject: "Test Subject",
		Body:    "Test Body",
	}
	err = validateEmailMessage(validBccMsg)
	if err != nil {
		t.Errorf("Unexpected error for valid message with Bcc: %v", err)
	}
}

func TestEmailMessageMethods(t *testing.T) {
	// Create a message manually for testing
	msg := &EmailMessage{
		From: "test@example.com",
		To:   []string{},
		Cc:   []string{},
		Bcc:  []string{},
	}

	// Test AddRecipient
	msg.AddRecipient("recipient@example.com")
	if len(msg.To) != 1 || msg.To[0] != "recipient@example.com" {
		t.Errorf("AddRecipient failed")
	}

	// Test AddCC
	msg.AddCC("cc@example.com")
	if len(msg.Cc) != 1 || msg.Cc[0] != "cc@example.com" {
		t.Errorf("AddCC failed")
	}

	// Test AddBCC
	msg.AddBCC("bcc@example.com")
	if len(msg.Bcc) != 1 || msg.Bcc[0] != "bcc@example.com" {
		t.Errorf("AddBCC failed")
	}

	// Test SetTextBody
	msg.SetTextBody("Test plain text body")
	if msg.Body != "Test plain text body" || msg.IsHTML {
		t.Errorf("SetTextBody failed")
	}
}
