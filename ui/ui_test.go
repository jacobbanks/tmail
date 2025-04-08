package ui

import (
	"testing"
	"time"

	"github.com/jacobbanks/tmail/email"
)

// Test helper to create a fake email for testing UI components
func createTestEmail() *email.IncomingMessage {
	return &email.IncomingMessage{
		From:        "test@example.com",
		To:          "recipient@example.com",
		Subject:     "Test Subject",
		Date:        time.Now(),
		Body:        "This is a test email body",
		Attachments: []string{},
	}
}

// Test helper to create a fake email with HTML-derived content
func createTestHTMLEmail() *email.IncomingMessage {
	return &email.IncomingMessage{
		From:        "test@example.com",
		To:          "recipient@example.com",
		Subject:     "HTML Test Subject",
		Date:        time.Now(),
		Body:        "This is a test email body that was converted from HTML",
		Attachments: []string{},
	}
}

// Test helper to create a fake email with attachments
func createTestEmailWithAttachments() *email.IncomingMessage {
	return &email.IncomingMessage{
		From:        "test@example.com",
		To:          "recipient@example.com",
		Subject:     "Email with Attachments",
		Date:        time.Now(),
		Body:        "This is a test email with attachments",
		Attachments: []string{"document.pdf", "image.jpg", "archive.zip"},
	}
}

// Basic initialization test for EmailComposer
func TestNewEmailComposer(t *testing.T) {
	// Skip UI tests to avoid issues with terminal IO
	t.Skip("Skipping UI tests that require terminal interaction")
}

// Basic initialization test for EmailReader
func TestNewEmailReader(t *testing.T) {
	// Skip UI tests to avoid issues with terminal IO
	t.Skip("Skipping UI tests that require terminal interaction")
}
