package email

import (
	"errors"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"

	"github.com/jordan-wright/email"
)

type EmailMessage struct {
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Body        string
	IsHTML      bool
	Attachments []string // Paths to attachment files
}

func NewEmailMessage() (*EmailMessage, error) {
	userInfo := getUserInfo()
	if userInfo.Email == "" {
		return nil, errors.New("no user email found - please set up your account first")
	}

	return &EmailMessage{
		From:        userInfo.Email,
		To:          []string{},
		Cc:          []string{},
		Bcc:         []string{},
		Attachments: []string{},
	}, nil
}

func (e *EmailMessage) AddRecipient(email string) {
	e.To = append(e.To, email)
}

func (e *EmailMessage) AddCC(email string) {
	e.Cc = append(e.Cc, email)
}

func (e *EmailMessage) AddBCC(email string) {
	e.Bcc = append(e.Bcc, email)
}

func (e *EmailMessage) SetHTMLBody(htmlContent string) {
	e.Body = htmlContent
	e.IsHTML = true
}

func (e *EmailMessage) SetTextBody(textContent string) {
	e.Body = textContent
	e.IsHTML = false
}

// AddAttachment adds a file as an attachment
func (e *EmailMessage) AddAttachment(filePath string) error {
	// Check if file exists and is readable
	_, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	e.Attachments = append(e.Attachments, filePath)
	return nil
}

func SendEmail(message *EmailMessage) error {
	// Validate email message
	if err := validateEmailMessage(message); err != nil {
		return err
	}

	userInfo := getUserInfo()
	if userInfo.Email == "" || userInfo.AppPassword == "" {
		return errors.New("missing email credentials - please set up your account first")
	}

	// Use default Gmail configuration
	config := DefaultConfig

	e := email.NewEmail()

	e.From = message.From

	e.To = sanitizeAddresses(message.To)
	e.Cc = sanitizeAddresses(message.Cc)
	e.Bcc = sanitizeAddresses(message.Bcc)

	e.Subject = sanitizeHeader(message.Subject)

	if message.IsHTML {
		e.HTML = []byte(message.Body)
	} else {
		e.Text = []byte(message.Body)
	}

	// Add attachments
	for _, attachPath := range message.Attachments {
		if _, err := e.AttachFile(attachPath); err != nil {
			return err
		}
	}

	auth := smtp.PlainAuth("", userInfo.Email, userInfo.AppPassword, config.SMTPHost)

	return e.Send(config.GetSMTPAddress(), auth)
}

func QuickSend(to, subject, body string) error {
	message, err := NewEmailMessage()
	if err != nil {
		return err
	}

	message.AddRecipient(to)
	message.Subject = subject
	message.SetTextBody(body)

	return SendEmail(message)
}

func validateEmailMessage(message *EmailMessage) error {
	if message == nil {
		return errors.New("email message is nil")
	}

	if len(message.To) == 0 && len(message.Cc) == 0 && len(message.Bcc) == 0 {
		return errors.New("email must have at least one recipient")
	}

	// Validate attachments
	for _, path := range message.Attachments {
		if _, err := os.Stat(path); err != nil {
			return errors.New("attachment not found: " + filepath.Base(path))
		}
	}

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

// sanitizeHeader removes CR and LF characters from header fields
func sanitizeHeader(text string) string {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\r\n", " ")
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")

	return text
}
