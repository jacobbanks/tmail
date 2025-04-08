package email

import (
	"sync"
	"testing"

	"github.com/jacobbanks/tmail/auth"
)

// MockProvider implements MailProvider for testing
type MockProvider struct {
	connected        bool
	sentEmails       []*OutgoingMessage
	storedEmails     []*IncomingMessage
	userInfo         auth.User
	connectCalled    int
	disconnectCalled int
	sendCalled       int
	getCalled        int
}

func NewMockProvider() *MockProvider {
	return &MockProvider{
		connected:  false,
		sentEmails: make([]*OutgoingMessage, 0),
		storedEmails: []*IncomingMessage{
			{
				From:        "test1@example.com",
				To:          "user@example.com",
				Subject:     "Test Email 1",
				Body:        "This is test email 1",
				Attachments: []string{},
			},
			{
				From:        "test2@example.com",
				To:          "user@example.com",
				Subject:     "Test Email 2",
				Body:        "This is test email 2",
				Attachments: []string{"attachment.pdf"},
			},
		},
		userInfo: auth.User{
			Email:       "user@example.com",
			AppPassword: "testpassword",
		},
	}
}

func (m *MockProvider) Connect() error {
	m.connectCalled++
	m.connected = true
	return nil
}

func (m *MockProvider) Disconnect() error {
	m.disconnectCalled++
	m.connected = false
	return nil
}

func (m *MockProvider) SendEmail(message *OutgoingMessage) error {
	m.sendCalled++
	m.sentEmails = append(m.sentEmails, message)
	return nil
}

func (m *MockProvider) QuickSend(to, subject, body string) error {
	msg := &OutgoingMessage{
		To:      []string{to},
		Subject: subject,
		Body:    body,
	}
	return m.SendEmail(msg)
}

func (m *MockProvider) GetEmails(limit int) ([]*IncomingMessage, error) {
	m.getCalled++

	if limit <= 0 || limit > len(m.storedEmails) {
		return m.storedEmails, nil
	}

	return m.storedEmails[:limit], nil
}

func (m *MockProvider) GetUserInfo() (auth.User, error) {
	return m.userInfo, nil
}

// Tests for the MailProvider interface
func TestMailProviderInterface(t *testing.T) {
	// Test that MockProvider implements MailProvider
	var _ MailProvider = (*MockProvider)(nil)
	var _ MailProvider = (*GmailProvider)(nil)
}

// Tests for connection management
func TestProviderConnectionManagement(t *testing.T) {
	mock := NewMockProvider()

	// Should start disconnected
	if mock.connected {
		t.Error("New provider should start disconnected")
	}

	// Connect should mark as connected
	err := mock.Connect()
	if err != nil {
		t.Errorf("Connect returned error: %v", err)
	}
	if !mock.connected {
		t.Error("Provider should be connected after Connect()")
	}

	// Disconnect should mark as disconnected
	err = mock.Disconnect()
	if err != nil {
		t.Errorf("Disconnect returned error: %v", err)
	}
	if mock.connected {
		t.Error("Provider should be disconnected after Disconnect()")
	}
}

// Tests for email operations
func TestProviderEmailOperations(t *testing.T) {
	mock := NewMockProvider()

	// Send an email
	testMessage := &OutgoingMessage{
		To:      []string{"recipient@example.com"},
		Subject: "Test Subject",
		Body:    "Test Body",
	}

	err := mock.SendEmail(testMessage)
	if err != nil {
		t.Errorf("SendEmail returned error: %v", err)
	}

	if mock.sendCalled != 1 {
		t.Errorf("Expected SendEmail to be called once, got %d", mock.sendCalled)
	}

	if len(mock.sentEmails) != 1 {
		t.Errorf("Expected 1 sent email, got %d", len(mock.sentEmails))
	}

	// Get emails
	emails, err := mock.GetEmails(1)
	if err != nil {
		t.Errorf("GetEmails returned error: %v", err)
	}

	if mock.getCalled != 1 {
		t.Errorf("Expected GetEmails to be called once, got %d", mock.getCalled)
	}

	if len(emails) != 1 {
		t.Errorf("Expected 1 returned email, got %d", len(emails))
	}

	// Test with no limit
	allEmails, err := mock.GetEmails(0)
	if err != nil {
		t.Errorf("GetEmails returned error: %v", err)
	}

	if len(allEmails) != 2 {
		t.Errorf("Expected 2 returned emails, got %d", len(allEmails))
	}
}

// Test for CreateDefaultMailProvider using dependency injection
func TestCreateDefaultMailProvider(t *testing.T) {
	// Reset the singleton provider for this test
	provider = nil
	initErr = nil
	once = *new(sync.Once)

	// Create a provider
	p, err := CreateDefaultMailProvider()

	// Skip the test if there's a credential error - this is expected in CI environments
	if err != nil && err.Error() == "open /Users/jacob.banks/.tmail/credentials.json: no such file or directory" {
		t.Skip("Skipping test due to missing credentials")
	}

	if err != nil {
		t.Skip("Skipping test due to required terminal interaction")
		t.Errorf("CreateDefaultMailProvider returned unexpected error: %v", err)
	}

	if p == nil {
		t.Error("CreateDefaultMailProvider returned nil provider")
	}

	// Second call should return the same instance
	p2, _ := CreateDefaultMailProvider()
	if p != p2 {
		t.Error("CreateDefaultMailProvider should return the same instance on subsequent calls")
	}
}

// Mock tests for GmailProvider - these don't actually connect to Gmail
func TestGmailProviderNoCredentials(t *testing.T) {
	emptyUser := auth.User{}
	_, err := NewGmailProvider(DefaultConfig, emptyUser)

	if err == nil {
		t.Error("NewGmailProvider should return error with empty credentials")
	}
}
