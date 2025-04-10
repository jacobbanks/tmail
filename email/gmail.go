package email

import (
	"fmt"
	"net/smtp"

	"github.com/emersion/go-imap"
	imapClient "github.com/emersion/go-imap/client"
	"github.com/jacobbanks/tmail/auth"
)

// GmailProvider implements the MailProvider interface for Gmail
type GmailProvider struct {
	client    *imapClient.Client
	config    Config
	userInfo  auth.Credentials
	connected bool
}

// NewGmailProvider creates a new Gmail provider
func NewGmailProvider(config Config, userInfo auth.Credentials) (*GmailProvider, error) {
	provider := &GmailProvider{
		config:    config,
		userInfo:  userInfo,
		connected: false,
	}

	// Connect immediately
	if err := provider.Connect(); err != nil {
		return nil, err
	}

	return provider, nil
}

// Connect establishes a connection to Gmail's IMAP server using the provider's credentials.
// If already connected, it returns nil without reconnecting.
func (p *GmailProvider) Connect() error {
	if p.connected && p.client != nil {
		return nil // Already connected
	}

	// Validate credentials
	if p.userInfo.Email == "" || p.userInfo.AppPassword == "" {
		return fmt.Errorf("missing email credentials - please set up your account first")
	}

	client, err := imapClient.DialTLS(p.config.GetIMAPAddress(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to IMAP server: %v", err)
	}

	err = client.Login(p.userInfo.Email, p.userInfo.AppPassword)
	if err != nil {
		client.Logout() // Clean up before returning error
		return fmt.Errorf("failed to login: %v", err)
	}

	p.client = client
	p.connected = true
	return nil
}

// Disconnect closes the IMAP connection to the Gmail server.
// If already disconnected, returns nil without any action.
func (p *GmailProvider) Disconnect() error {
	if !p.connected || p.client == nil {
		return nil // Already disconnected
	}

	err := p.client.Logout()
	if err != nil {
		return fmt.Errorf("error during logout: %v", err)
	}

	p.client = nil
	p.connected = false
	return nil
}

// IsConnected checks if the provider is currently connected
func (p *GmailProvider) isConnected() bool {
	return p.connected && p.client != nil
}

// GetEmails retrieves and parses emails
func (p *GmailProvider) GetEmails(limit int) ([]*IncomingMessage, error) {
	if !p.isConnected() {
		if err := p.Connect(); err != nil {
			return nil, err
		}
	}

	// Fetch raw messages first
	messages, err := p.fetchMessages(limit)
	if err != nil {
		return nil, err
	}

	// Parse messages into Email objects
	var emails []*IncomingMessage
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		email := &IncomingMessage{}
		if err := email.Parse(msg); err != nil {
			// Skip emails that fail to parse
			continue
		}
		emails = append(emails, email)
	}

	return emails, nil
}

// fetchMessages gets raw IMAP messages
func (p *GmailProvider) fetchMessages(limit int) ([]*imap.Message, error) {
	// Select inbox
	mailbox, err := p.client.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("failed to select inbox: %v", err)
	}
	seqSet := new(imap.SeqSet)

	if limit <= 0 {
		limit = 10 // Default to 10 emails
	}

	from := uint32(1)
	if mailbox.Messages > uint32(limit) {
		from = mailbox.Messages - uint32(limit) + 1
	}

	seqSet.AddRange(from, mailbox.Messages)

	section := &imap.BodySectionName{}

	items := []imap.FetchItem{
		imap.FetchEnvelope,
		imap.FetchBodyStructure,
		imap.FetchFlags,
		section.FetchItem(),
	}

	messages := make(chan *imap.Message, limit)

	// Start the fetch operation in a goroutine
	done := make(chan error, 1)
	go func() {
		done <- p.client.Fetch(seqSet, items, messages)
	}()

	var emails []*imap.Message
	for msg := range messages {
		emails = append(emails, msg)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("fetch failed: %v", err)
	}

	return emails, nil
}

// SendEmail sends an email message
func (p *GmailProvider) SendEmail(message *OutgoingMessage) error {
	// Validate email message
	if err := validateEmailMessage(message); err != nil {
		return err
	}

	// e := jmail.NewEmail()
	//
	// e.From = message.From
	// e.To = message.To
	// e.Cc = message.Cc
	// e.Bcc = message.Bcc
	// e.Subject = message.Subject
	// e.Text = []byte(message.Text)

	// Add attachments
	for _, path := range message.AttachmentPaths {
		if _, err := message.PrepAttachment(path); err != nil {
			return err
		}
	}

	// Create SMTP auth
	auth := smtp.PlainAuth("", p.userInfo.Email, p.userInfo.AppPassword, p.config.SMTPHost)

	return message.SendMessage(p.config.GetSMTPAddress(), auth)
}

// QuickSend provides a simple way to send a text email
func (p *GmailProvider) QuickSend(to, subject, body string) error {
	message, err := NewOutgoingMessage()
	if err != nil {
		return err
	}

	message.AddRecipient(to)
	message.Subject = subject
	message.SetTextBody(body)

	return p.SendEmail(message)
}

// GetUserInfo returns the user information
func (p *GmailProvider) GetUserInfo() (auth.Credentials, error) {
	return p.userInfo, nil
}

// validateEmailMessage verifies that an email message is valid
func validateEmailMessage(message *OutgoingMessage) error {
	if message == nil {
		return fmt.Errorf("email message is nil")
	}

	if len(message.To) == 0 && len(message.Cc) == 0 && len(message.Bcc) == 0 {
		return fmt.Errorf("email must have at least one recipient")
	}

	return nil
}
