package email

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/mail"
)

type Email struct {
	From        string
	To          string
	Subject     string
	Date        time.Time
	Body        string
	IsHTML      bool
	Attachments []string
}

func (email *Email) Parse(msg *imap.Message) error {
	if msg == nil {
		return fmt.Errorf("cannot parse a nil message")
	}

	if len(msg.Body) == 0 {
		// If no body but we have envelope data, create a basic email
		if msg.Envelope != nil {
			createEmailFromEnvelope(email, msg.Envelope)
			return nil
		}
		return fmt.Errorf("message has no body parts")
	}

	reader := findBodyReader(msg)
	if reader == nil {
		return fmt.Errorf("no message body found")
	}

	// Create a message entity
	entity, err := message.Read(reader)
	if err != nil {
		return fmt.Errorf("failed to parse message: %v", err)
	}

	mr := mail.NewReader(entity)
	err = createEmail(mr, email)
	if err != nil {
		log.Printf("")
		return err
	}

	return nil
}

func createEmail(reader *mail.Reader, email *Email) error {
	emailHeader := reader.Header
	err := extractHeaders(emailHeader, email)
	if err != nil {
		log.Printf("Error extracting email headers: %v", err)
		return err
	}

	err = extractBodyAndAttachments(reader, email)
	return nil
}

func findBodyReader(msg *imap.Message) io.Reader {
	emptySection := &imap.BodySectionName{}
	if body, ok := msg.Body[emptySection]; ok {
		return body
	}

	// TEXT section
	textSection := &imap.BodySectionName{
		BodyPartName: imap.BodyPartName{
			Specifier: imap.TextSpecifier,
		},
	}
	if body, ok := msg.Body[textSection]; ok {
		return body
	}

	// Just use any available section
	for _, body := range msg.Body {
		return body
	}

	return nil
}

// createEmailFromEnvelope creates a basic Email from just envelope data
func createEmailFromEnvelope(email *Email, envelope *imap.Envelope) error {
	email.Subject = envelope.Subject
	email.Date = envelope.Date

	if len(envelope.From) > 0 {
		email.From = formatImapAddress(envelope.From[0])
	}

	if len(envelope.To) > 0 {
		email.To = formatImapAddressList(envelope.To)
	}

	email.Body = "(Message body not available)"

	return nil
}

func extractHeaders(header mail.Header, email *Email) error {
	from, err := header.AddressList("From")
	if err != nil {
		log.Printf("Error extracting From header: %v", err)
	}

	to, err := header.AddressList("To")
	if err != nil {
		log.Printf("Error extracting To header: %v", err)
	}

	subject, err := header.Subject()
	if err != nil {
		log.Printf("Error extracting Subject header: %v", err)
		subject = "(No subject)"
	}

	date, err := header.Date()
	if err != nil {
		log.Printf("Error extracting Date header: %v", err)
		date = time.Now() // Fallback to current time
	}

	email.From = formatAddressList(from)
	email.To = formatAddressList(to)
	email.Subject = subject
	email.Date = date

	return nil
}

func extractBodyAndAttachments(reader *mail.Reader, email *Email) error {
	var plainText, htmlText string
	var attachments []string

	// Process each part of the message
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Error reading message part: %v", err)
			continue
		}

		switch header := part.Header.(type) {
		case *mail.InlineHeader:
			// This is message content (either plain text or HTML)
			contentType := "text/plain"

			if ct, _, err := header.ContentType(); err == nil {
				contentType = ct
			} else {
				log.Printf("Error getting content type, using default: %v", err)
				return err
			}

			log.Printf("Processing email part with content type: %s", contentType)

			content := readContent(part.Body)

			if strings.HasPrefix(contentType, "text/plain") {
				plainText = content
			} else if strings.HasPrefix(contentType, "text/html") {
				htmlText = content
			}

		case *mail.AttachmentHeader:
			// This is an attachment
			filename, err := header.Filename()
			if err != nil {
				log.Printf("Error getting attachment filename: %v", err)
				filename = "unknown-attachment"
				return err
			}
			attachments = append(attachments, filename)
		}
	}

	// Prefer plain text over HTML will come back and do HTML if I have time.
	if plainText != "" {
		email.Body, email.IsHTML, email.Attachments = plainText, false, attachments
		log.Printf("plainText email %v", email)
		return nil
	}
	if htmlText != "" {
		log.Printf("html email %v, and email %v", htmlText, email)
		// implement html later
		email.Body, email.IsHTML, email.Attachments = plainText, true, attachments
		return nil
	}

	email.Body, email.IsHTML, email.Attachments = "(No content found)", false, attachments
	return nil
}

func readContent(reader io.Reader) string {
	// Use a limit to avoid any issues with overly large messages
	const maxReadSize = 10 * 1024 * 1024 // 10MB max

	lReader := io.LimitReader(reader, maxReadSize)

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(lReader)
	if err != nil {
		log.Printf("Error reading content: %v", err)
		return "(Error reading content)"
	}

	// Check if we reached the limit
	if buf.Len() >= maxReadSize {
		log.Printf("Warning: Message content exceeded max size limit of %d bytes", maxReadSize)
		return buf.String() + "\n\n[... Message truncated due to size ...]"
	}

	return buf.String()
}

func formatImapAddress(addr *imap.Address) string {
	if addr == nil {
		return ""
	}

	if addr.PersonalName != "" {
		return fmt.Sprintf("%s <%s@%s>", addr.PersonalName, addr.MailboxName, addr.HostName)
	}
	return fmt.Sprintf("%s@%s", addr.MailboxName, addr.HostName)
}

func formatImapAddressList(addresses []*imap.Address) string {
	if len(addresses) == 0 {
		return ""
	}

	var result strings.Builder
	for i, addr := range addresses {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(formatImapAddress(addr))
	}
	return result.String()
}

func formatAddressList(addresses []*mail.Address) string {
	if len(addresses) == 0 {
		return ""
	}

	var result strings.Builder
	for i, addr := range addresses {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(addr.String())
	}
	return result.String()
}
