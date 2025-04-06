package ui

import (
	"testing"
	"time"

	"github.com/jacobbanks/tmail/email"
)

// Test helper to create a fake email for testing UI components
func createTestEmail() *email.Email {
	return &email.Email{
		From:        "test@example.com",
		To:          "recipient@example.com",
		Subject:     "Test Subject",
		Date:        time.Now(),
		Body:        "This is a test email body",
		IsHTML:      false,
		Attachments: []string{},
	}
}

// Test helper to create a fake email with HTML content
func createTestHTMLEmail() *email.Email {
	return &email.Email{
		From:        "test@example.com",
		To:          "recipient@example.com",
		Subject:     "HTML Test Subject",
		Date:        time.Now(),
		Body:        "This is a test email body",
		HTMLBody:    "<h1>HTML Content</h1><p>This is an HTML email</p>",
		IsHTML:      true,
		Attachments: []string{},
	}
}

// Test helper to create a fake email with attachments
func createTestEmailWithAttachments() *email.Email {
	return &email.Email{
		From:        "test@example.com",
		To:          "recipient@example.com",
		Subject:     "Email with Attachments",
		Date:        time.Now(),
		Body:        "This is a test email with attachments",
		IsHTML:      false,
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
