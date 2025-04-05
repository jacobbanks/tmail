package email

import (
	"errors"
	"log"
	"net/smtp"
	"strings"

	"github.com/jordan-wright/email"
)

type EmailMessage struct {
	From    string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Body    string
	IsHTML  bool
}

func NewEmailMessage() (*EmailMessage, error) {
	userInfo := getUserInfo()
	if userInfo.Email == "" {
		return nil, errors.New("no user email found - please set up your account first")
	}

	return &EmailMessage{
		From: userInfo.Email,
		To:   []string{},
		Cc:   []string{},
		Bcc:  []string{},
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

	auth := smtp.PlainAuth("", userInfo.Email, userInfo.AppPassword, config.SMTPHost)

	// Log sending info in debug mode
	log.Printf("Sending email from: %s", message.From)
	log.Printf("To: %v", e.To)
	if len(e.Cc) > 0 {
		log.Printf("Cc: %v", e.Cc)
	}
	if len(e.Bcc) > 0 {
		log.Printf("Bcc: %v", e.Bcc)
	}
	log.Printf("Subject: %s", e.Subject)

	return e.Send(config.GetSMTPAddress(), auth)
}

// Used for Testing
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
