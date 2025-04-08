package email

import (
	"errors"
	"os"
	"strings"

	"github.com/jacobbanks/tmail/auth"
)

type OutgoingMessage struct {
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Body        string
	IsHTML      bool
	Attachments []string // Paths to attachment files
}

func NewOutgoingMessage() (*OutgoingMessage, error) {
	// Get user information
	userInfo, err := auth.LoadUser()
	if err != nil {
		return nil, err
	}

	if userInfo.Email == "" {
		return nil, errors.New("no user email found - please set up your account first")
	}

	return &OutgoingMessage{
		From:        userInfo.Email,
		To:          []string{},
		Cc:          []string{},
		Bcc:         []string{},
		Attachments: []string{},
	}, nil
}

func (e *OutgoingMessage) AddRecipient(email string) {
	e.To = append(e.To, email)
}

func (e *OutgoingMessage) AddCC(email string) {
	e.Cc = append(e.Cc, email)
}

func (e *OutgoingMessage) AddBCC(email string) {
	e.Bcc = append(e.Bcc, email)
}

func (e *OutgoingMessage) SetHTMLBody(htmlContent string) {
	e.Body = htmlContent
	e.IsHTML = true
}

func (e *OutgoingMessage) SetTextBody(textContent string) {
	e.Body = textContent
	e.IsHTML = false
}

// AddAttachment adds a file as an attachment
func (e *OutgoingMessage) AddAttachment(filePath string) error {
	// Check if file exists and is readable
	_, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	e.Attachments = append(e.Attachments, filePath)
	return nil
}

func sanitizeAddresses(addresses []string) []string {
	var clean []string
	for _, addr := range addresses {
		addr = strings.TrimSpace(addr)
		if addr != "" {
			clean = append(clean, addr)
		}
	}
	return clean
}

// SendEmail is a helper function that creates a mail provider
// and sends an email message
// func SendEmail(message *OutgoingMessage) error {
// 	// Create mail provider
// 	provider, err := NewGmailProvider()
// 	if err != nil {
// 		return err
// 	}
// 	defer provider.disconnect()
//
// 	// Send the message
// 	return provider.SendEmail(message)
// }

// // QuickSend provides a simple way to send a text email
// func QuickSend(to, subject, body string) error {
// 	message, err := NewEmailMessage()
// 	if err != nil {
// 		return err
// 	}
//
// 	message.AddRecipient(to)
// 	message.Subject = subject
// 	message.SetTextBody(body)
//
// 	return SendEmail(message)
// }

// // sanitizeHeader removes CR and LF characters from header fields
// func sanitizeHeader(text string) string {
// 	text = strings.TrimSpace(text)
// 	text = strings.ReplaceAll(text, "\r\n", " ")
// 	text = strings.ReplaceAll(text, "\n", " ")
// 	text = strings.ReplaceAll(text, "\r", " ")
//
// 	return text
// }
