package email

import (
	"log"
	"sync"

	"github.com/jacobbanks/tmail/auth"
)

// MailProvider defines the interface for interacting with email providers.
// It abstracts the details of connecting to email servers, sending and receiving
// messages, and managing authentication.
type MailProvider interface {
	// Connection management
	Connect() error
	Disconnect() error

	// Mail operations
	SendEmail(message *OutgoingMessage) error
	QuickSend(to, subject, body string) error
	GetEmails(limit int) ([]*IncomingMessage, error)
	GetUserInfo() (auth.Credentials, error)
}

var (
	provider MailProvider
	once     sync.Once
	initErr  error
)

// CreateDefaultMailProvider creates a mail provider with default configuration
func CreateDefaultMailProvider() (MailProvider, error) {
	once.Do(func() {
		userInfo, err := auth.LoadUser()
		if err != nil {
			log.Println("Cannot load user while creating mail provider")
			initErr = err
			return
		}
		provider, err = NewGmailProvider(DefaultConfig, userInfo)
		if err != nil {
			log.Println("Cannot get Default Mail Provider")
			initErr = err
			return
		}
		initErr = nil
	})
	// Create provider with default config
	return provider, initErr
}
