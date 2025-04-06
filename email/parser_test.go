package email

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/emersion/go-imap"
)

func TestReadContent(t *testing.T) {
	// Test normal reading
	normalContent := "This is normal content"
	reader := bytes.NewBufferString(normalContent)
	result := readContent(reader)
	if result != normalContent {
		t.Errorf("readContent failed for normal content: got %q, expected %q", result, normalContent)
	}

	// Test error handling
	errorReader := errorReaderMock{err: io.ErrUnexpectedEOF}
	result = readContent(errorReader)
	if result != "(Error reading content)" {
		t.Errorf("readContent failed for error: got %q, expected error message", result)
	}

	// Test size limit with a reader that pretends to be infinite
	limitTestReader := &limitTestReaderMock{}
	result = readContent(limitTestReader)

	// Verify truncation message is present
	if !bytes.Contains([]byte(result), []byte("[... Message truncated due to size ...]")) {
		t.Errorf("readContent should indicate truncation in large messages")
	}
}

func TestFormatAddressList(t *testing.T) {
	// Empty address list
	if formatAddressList(nil) != "" {
		t.Errorf("formatAddressList failed for nil")
	}

	// Single address
	addresses := []*imap.Address{
		{PersonalName: "Test User", MailboxName: "test", HostName: "example.com"},
	}
	expected := "Test User <test@example.com>"
	result := formatImapAddressList(addresses)
	if result != expected {
		t.Errorf("formatImapAddressList failed: got %q, expected %q", result, expected)
	}

	// Multiple addresses
	addresses = append(addresses, &imap.Address{
		MailboxName: "another", HostName: "example.com",
	})
	expected = "Test User <test@example.com>, another@example.com"
	result = formatImapAddressList(addresses)
	if result != expected {
		t.Errorf("formatImapAddressList failed: got %q, expected %q", result, expected)
	}
}

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

	// Create an email struct to populate
	email := &Email{}
	err := createEmailFromEnvelope(email, envelope)

	if err != nil {
		t.Fatalf("createEmailFromEnvelope returned error: %v", err)
	}

	if email.Subject != "Test Subject" {
		t.Errorf("Subject mismatch: got %q, expected %q", email.Subject, "Test Subject")
	}

	if email.From != "Sender <sender@example.com>" {
		t.Errorf("From mismatch: got %q, expected %q", email.From, "Sender <sender@example.com>")
	}

	if email.To != "Recipient <recipient@example.com>" {
		t.Errorf("To mismatch: got %q, expected %q", email.To, "Recipient <recipient@example.com>")
	}

	if !email.Date.Equal(now) {
		t.Errorf("Date mismatch: got %v, expected %v", email.Date, now)
	}
}

// Mock for testing error cases
type errorReaderMock struct {
	err error
}

func (e errorReaderMock) Read(p []byte) (n int, err error) {
	return 0, e.err
}

// Mock reader that pretends to be larger than the max size limit
// but actually just returns 'A's indefinitely until hitting the limit
type limitTestReaderMock struct {
	readCount int
}

func (r *limitTestReaderMock) Read(p []byte) (n int, err error) {
	// Fill the buffer with 'A's
	for i := range p {
		p[i] = 'A'
	}

	r.readCount += len(p)

	// Never return EOF - pretend the file is infinite
	// The readContent function's LimitReader will stop reading
	return len(p), nil
}
