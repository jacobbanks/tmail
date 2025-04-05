package ui

import (
	"testing"
	"time"

	"github.com/jacobbanks/tmail/email"
	"github.com/rivo/tview"
)

// TestNewEmailComposer verifies that the composer can be created
func TestNewEmailComposer(t *testing.T) {
	// Create a composer with no reply-to email
	composer := NewEmailComposer(nil)
	if composer == nil {
		t.Fatal("NewEmailComposer returned nil")
	}
	
	// Check that the app was created
	if composer.app == nil {
		t.Error("Composer app was not initialized")
	}
	
	// Verify debug mode is off by default
	if composer.debugMode {
		t.Error("Debug mode should be off by default")
	}
	
	// Test setting debug mode
	composer.SetDebugMode(true)
	if !composer.debugMode {
		t.Error("Debug mode was not set correctly")
	}
	
	// Check that all required form fields are present
	if composer.form == nil {
		t.Error("Form was not initialized")
	} else {
		// Should have 4 fields: To, Cc, Bcc, Subject
		fieldCount := 0
		for i := 0; i < composer.form.GetFormItemCount(); i++ {
			if _, ok := composer.form.GetFormItem(i).(*tview.InputField); ok {
				fieldCount++
			}
		}
		
		if fieldCount != 4 {
			t.Errorf("Expected 4 input fields, got %d", fieldCount)
		}
	}
	
	// Check that body area is initialized
	if composer.bodyArea == nil {
		t.Error("Body text area was not initialized")
	}
}

// TestNewEmailReader verifies that the reader can be created
func TestNewEmailReader(t *testing.T) {
	// Create sample emails for testing
	emails := []*email.Email{
		{
			From:    "sender@example.com",
			To:      "recipient@example.com",
			Subject: "Test Email 1",
			Date:    time.Now(),
			Body:    "This is test email 1",
		},
		{
			From:    "sender2@example.com",
			To:      "recipient2@example.com",
			Subject: "Test Email 2",
			Date:    time.Now().Add(-time.Hour),
			Body:    "This is test email 2",
		},
	}
	
	// Create a reader
	reader := NewEmailReader(emails)
	if reader == nil {
		t.Fatal("NewEmailReader returned nil")
	}
	
	// Check that the app was created
	if reader.app == nil {
		t.Error("Reader app was not initialized")
	}
	
	// Check that emails were stored
	if len(reader.emails) != 2 {
		t.Errorf("Expected 2 emails, got %d", len(reader.emails))
	}
}