package email

import (
	"testing"
	"time"

	"github.com/emersion/go-imap"
)

func TestFormatImapAddress(t *testing.T) {
	addr := formatImapAddress(nil)
	if addr != "" {
		t.Errorf("Expected empty string for nil address, got %q", addr)
	}

	imapAddr := &imap.Address{
		MailboxName: "test",
		HostName:    "example.com",
	}
	addr = formatImapAddress(imapAddr)
	expected := "test@example.com"
	if addr != expected {
		t.Errorf("Expected %q, got %q", expected, addr)
	}

	imapAddr = &imap.Address{
		PersonalName: "Test User",
		MailboxName:  "test",
		HostName:     "example.com",
	}
	addr = formatImapAddress(imapAddr)
	expected = "Test User <test@example.com>"
	if addr != expected {
		t.Errorf("Expected %q, got %q", expected, addr)
	}
}

func TestFormatImapAddressList(t *testing.T) {
	addrs := []*imap.Address{}
	result := formatImapAddressList(addrs)
	if result != "" {
		t.Errorf("Expected empty string for empty list, got %q", result)
	}

	addrs = []*imap.Address{
		{PersonalName: "Test User", MailboxName: "test", HostName: "example.com"},
	}
	result = formatImapAddressList(addrs)
	expected := "Test User <test@example.com>"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test with multiple addresses
	addrs = []*imap.Address{
		{PersonalName: "User One", MailboxName: "one", HostName: "example.com"},
		{MailboxName: "two", HostName: "example.com"},
		{PersonalName: "User Three", MailboxName: "three", HostName: "example.com"},
	}
	result = formatImapAddressList(addrs)
	if result == "" {
		t.Errorf("Expected non-empty result for multiple addresses, got empty string")
	}
}

// Mock implementation for testing Parse with envelope
func TestCreateEmailFromEnvelope(t *testing.T) {
	now := time.Now()
	envelope := &imap.Envelope{
		Date:    now,
		Subject: "Test Subject",
		From: []*imap.Address{
			{PersonalName: "Sender", MailboxName: "sender", HostName: "example.com"},
		},
		To: []*imap.Address{
			{PersonalName: "Recipient", MailboxName: "recipient", HostName: "example.com"},
		},
	}

	email := &Email{}
	err := createEmailFromEnvelope(email, envelope)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if email.Subject != "Test Subject" {
		t.Errorf("Expected Subject to be 'Test Subject', got %q", email.Subject)
	}

	expectedFrom := "Sender <sender@example.com>"
	if email.From != expectedFrom {
		t.Errorf("Expected From to be %q, got %q", expectedFrom, email.From)
	}

	if !email.Date.Equal(now) {
		t.Errorf("Expected Date to be %v, got %v", now, email.Date)
	}

	expectedBody := "(Message body not available)"
	if email.Body != expectedBody {
		t.Errorf("Expected Body to be %q, got %q", expectedBody, email.Body)
	}
}
