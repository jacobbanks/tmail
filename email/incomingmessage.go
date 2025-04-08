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

// IncomingMessage represents an email message received from a mail server.
// It contains the essential fields from the email such as sender, recipient,
// subject, date, and the message body in plain text format.
type IncomingMessage struct {
	From        string
	To          string
	Subject     string
	Date        time.Time
	Body        string
	Attachments []string // Only attachment names, not content
}

// Parse converts an IMAP message into an IncomingMessage structure.
// It extracts headers, body content, and attachment information from the raw message.
func (email *IncomingMessage) Parse(msg *imap.Message) error {
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
		return err
	}

	return nil
}

func createEmail(reader *mail.Reader, email *IncomingMessage) error {
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
func createEmailFromEnvelope(email *IncomingMessage, envelope *imap.Envelope) error {
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

func extractHeaders(header mail.Header, email *IncomingMessage) error {
	from, err := header.AddressList("From")
	if err != nil {
		// Continue with empty From field
	}

	to, err := header.AddressList("To")
	if err != nil {
		// Continue with empty To field
	}

	subject, err := header.Subject()
	if err != nil {
		subject = "(No subject)"
	}

	date, err := header.Date()
	if err != nil {
		date = time.Now() // Fallback to current time
	}

	email.From = formatAddressList(from)
	email.To = formatAddressList(to)
	email.Subject = subject
	email.Date = date

	return nil
}

func extractBodyAndAttachments(reader *mail.Reader, email *IncomingMessage) error {
	var plainText string
	var attachmentNames []string

	// Process each part of the message
	partCount := 0
	maxParts := 20 // Reasonable limit to avoid excessive processing

	for partCount < maxParts {
		partCount++

		part, err := reader.NextPart()
		if err != nil {
			break
		}

		switch header := part.Header.(type) {
		case *mail.InlineHeader:
			// This is message content (either plain text or HTML)
			contentType := "text/plain"

			if ct, _, err := header.ContentType(); err == nil {
				contentType = ct
			} else {
				continue
			}

			// Prioritize plain text for better memory efficiency
			if strings.HasPrefix(contentType, "text/plain") {
				plainText = readContent(part.Body)
				// Break early if we found plain text to avoid processing HTML
				if plainText != "" {
					break
				}
			}
		case *mail.AttachmentHeader:
			// Just store attachment names, not the content
			filename, err := header.Filename()
			if err != nil {
				filename = "unknown-attachment"
			}
			attachmentNames = append(attachmentNames, filename)
		}
	}

	// Store just what we need
	email.Body = plainText
	email.Attachments = attachmentNames

	// If we still have no content
	if email.Body == "" {
		email.Body = "(No content found)"
	}
	return nil
}

func readContent(reader io.Reader) string {
	// Use a much smaller limit for terminal display
	const maxReadSize = 1 * 1024 * 1024 // 1MB max for terminal display

	lReader := io.LimitReader(reader, maxReadSize)

	// Use a fixed-size buffer for better memory management
	buf := bytes.NewBuffer(make([]byte, 0, 32*1024)) // Pre-allocate 32KB
	_, err := buf.ReadFrom(lReader)
	if err != nil {
		return "(Error reading content)"
	}

	// Check if we reached the limit
	if buf.Len() >= maxReadSize {
		return buf.String()[:maxReadSize-256] + "\n\n[... Message truncated due to size ...]"
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
