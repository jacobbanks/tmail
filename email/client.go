package email

import (
	"fmt"
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/jacobbanks/tmail/auth"
)

func GetEmails(limit int) ([]*imap.Message, error) {
	userInfo, err := auth.LoadUser()
	if err != nil {
		return nil, fmt.Errorf("failed to load user credentials: %v", err)
	}

	if userInfo.Email == "" || userInfo.AppPassword == "" {
		return nil, fmt.Errorf("missing email credentials - please set up your account first")
	}

	// Use default Gmail configuration
	config := DefaultConfig

	c, err := client.DialTLS(config.GetIMAPAddress(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to IMAP server: %v", err)
	}
	defer c.Logout()

	err = c.Login(userInfo.Email, userInfo.AppPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %v", err)
	}

	mailbox, err := c.Select("INBOX", false)
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
		done <- c.Fetch(seqSet, items, messages)
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

func GetUserInfo() (auth.User, error) {
	return auth.LoadUser()
}

// Internal version used by package functions
func getUserInfo() auth.User {
	user, err := auth.LoadUser()
	if err != nil {
		log.Printf("Cannot get user info: %v", err)
		return auth.User{}
	}
	return user
}

func FetchEmails(limit int) []*Email {
	messages, err := GetEmails(limit)
	if err != nil {
		log.Printf("Failed to fetch emails: %v", err)
		return nil
	}

	var parsedEmails []*Email
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		email := Email{}
		err := email.Parse(msg)
		if err != nil {
			log.Printf("Failed to parse email: %v", err)
			continue
		}

		parsedEmails = append(parsedEmails, &email)
	}

	return parsedEmails
}
