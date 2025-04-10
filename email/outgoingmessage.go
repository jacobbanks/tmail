package email

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jacobbanks/tmail/auth"
)

const MaxLineLength = 76 // MaxLineLength is the maximum line length per RFC 2045

// Attachment is a struct representing an email attachment.
// Based on the mime/multipart.FileHeader struct, Attachment contains the name, MIMEHeader, and content of the attachment in question
type Attachment struct {
	Filename    string
	ContentType string
	Header      textproto.MIMEHeader
	Content     []byte
}

// OutgoingMessage represents an email message to be sent.
// It contains all the necessary fields for sending an email, including
// the sender, recipients, subject, body content, and any attachments.
type OutgoingMessage struct {
	Headers         textproto.MIMEHeader
	From            string
	To              []string
	Cc              []string
	Bcc             []string
	Subject         string
	Text            []byte
	Attachments     []*Attachment
	AttachmentPaths []string
	ReplyTo         []string
}

// NewOutgoingMessage creates a new email message with sender information
// populated from the authenticated user's credentials. Returns an error
// if user credentials cannot be loaded or are incomplete.
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
		From:            userInfo.Email,
		To:              []string{},
		Cc:              []string{},
		Bcc:             []string{},
		Attachments:     []*Attachment{},
		AttachmentPaths: []string{},
	}, nil
}

// AddRecipient adds an email address to the To field of the message.
func (msg *OutgoingMessage) AddRecipient(email string) {
	msg.To = append(msg.To, email)
}

// AddCC adds an email address to the CC (carbon copy) field of the message.
func (msg *OutgoingMessage) AddCC(email string) {
	msg.Cc = append(msg.Cc, email)
}

// AddBCC adds an email address to the BCC (blind carbon copy) field of the message.
func (msg *OutgoingMessage) AddBCC(email string) {
	msg.Bcc = append(msg.Bcc, email)
}

// SetTextBody sets the message body as plain text content and marks the message accordingly.
func (msg *OutgoingMessage) SetTextBody(textContent string) {
	msg.Text = []byte(textContent)
}

// AppendAttachmentPath appends the path for the attachment to the emsil.
func (msg *OutgoingMessage) AppendAttachmentPath(path string) error {
	// Check if file exists and is readable
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	msg.AttachmentPaths = append(msg.AttachmentPaths, path)
	return nil
}

// Send an email using the given host and SMTP auth (optional), returns any error thrown by smtp.SendMail
// This function merges the To, Cc, and Bcc fields and calls the smtp.SendMail function using the Email.Bytes() output as the message
func (msg *OutgoingMessage) SendMessage(addr string, a smtp.Auth) error {
	// Merge the To, Cc, and Bcc fields
	to := make([]string, 0, len(msg.To)+len(msg.Cc)+len(msg.Bcc))
	to = append(append(append(to, msg.To...), msg.Cc...), msg.Bcc...)
	for i := 0; i < len(to); i++ {
		addr, err := mail.ParseAddress(to[i])
		if err != nil {
			return err
		}
		to[i] = addr.Address
	}
	// Check to make sure there is at least one recipient and one "From" address
	if msg.From == "" || len(to) == 0 {
		return errors.New("Must specify at least one From address and one To address")
	}
	sender, err := msg.parseSender()
	if err != nil {
		return err
	}
	raw, err := msg.ConvertToBytes()
	if err != nil {
		return err
	}
	return smtp.SendMail(addr, a, sender, to, raw)
}

// The function will return the created Attachment for reference, as well as nil for the error.
func (msg *OutgoingMessage) addAttachment(r io.Reader, filename string, c string) (a *Attachment, err error) {
	var b bytes.Buffer
	if _, err = io.Copy(&b, r); err != nil {
		return
	}
	attach := &Attachment{
		Filename:    filename,
		ContentType: c,
		Header:      textproto.MIMEHeader{},
		Content:     b.Bytes(),
	}
	msg.Attachments = append(msg.Attachments, attach)
	return attach, nil
}

// The function will then return the Attachment for reference, as well as nil for the error.
func (msg *OutgoingMessage) PrepAttachment(filename string) (a *Attachment, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	ct := mime.TypeByExtension(filepath.Ext(filename))
	basename := filepath.Base(filename)
	return msg.addAttachment(f, basename, ct)
}

// parseSender() parses From address
func (msg *OutgoingMessage) parseSender() (string, error) {
	from, err := mail.ParseAddress(msg.From)
	if err != nil {
		return "", err
	}
	return from.Address, nil
}

// Bytes converts the OutgoingMessage object to a []byte representation, including all needed MIMEHeaders, boundaries, etc.
func (msg *OutgoingMessage) ConvertToBytes() ([]byte, error) {
	buff := bytes.NewBuffer(make([]byte, 0, 4096))

	headers, err := msg.formatOutgoingMsgHeaders()
	if err != nil {
		return nil, err
	}

	var hasAttachments = len(msg.Attachments) > 0

	var w *multipart.Writer
	if hasAttachments {
		w = multipart.NewWriter(buff)
	}
	switch {
	case hasAttachments:
		headers.Set("Content-Type", "multipart/mixed;\r\n boundary="+w.Boundary())
	default:
		headers.Set("Content-Type", "text/plain; charset=UTF-8")
		headers.Set("Content-Transfer-Encoding", "quoted-printable")
	}

	headerToBytes(buff, headers)
	_, err = io.WriteString(buff, "\r\n")
	if err != nil {
		return nil, err
	}

	// Check to see if there is Text
	if len(msg.Text) > 0 {
		var subWriter *multipart.Writer

		if hasAttachments {
			subWriter = multipart.NewWriter(buff)
			header := textproto.MIMEHeader{
				"Content-Type": {"multipart/alternative;\r\n boundary=" + subWriter.Boundary()},
			}
			if _, err := w.CreatePart(header); err != nil {
				return nil, err
			}
		} else {
			subWriter = w
		}
		if len(msg.Text) > 0 {
			// Write the text
			if err := writeMessage(buff, msg.Text, hasAttachments, "text/plain", subWriter); err != nil {
				return nil, err
			}
		}
		if hasAttachments {
			if err := subWriter.Close(); err != nil {
				return nil, err
			}
		}
	}
	// Create attachment part, if necessary
	for _, a := range msg.Attachments {
		a.setDefaultHeaders()
		p, err := w.CreatePart(a.Header)
		if err != nil {
			return nil, err
		}
		// Write the base64Wrapped content to the part
		base64Encode(p, a.Content)
	}
	if hasAttachments {
		if err := w.Close(); err != nil {
			return nil, err
		}
	}
	return buff.Bytes(), nil
}

func (msg *OutgoingMessage) formatOutgoingMsgHeaders() (textproto.MIMEHeader, error) {
	res := make(textproto.MIMEHeader, len(msg.Headers)+6)
	if msg.Headers != nil {
		for _, h := range []string{"Reply-To", "To", "Cc", "From", "Subject", "Date", "Message-Id", "MIME-Version"} {
			if v, ok := msg.Headers[h]; ok {
				res[h] = v
			}
		}
	}
	// Set headers if there are values.
	if _, ok := res["Reply-To"]; !ok && len(msg.ReplyTo) > 0 {
		res.Set("Reply-To", strings.Join(msg.ReplyTo, ", "))
	}
	if _, ok := res["To"]; !ok && len(msg.To) > 0 {
		res.Set("To", strings.Join(msg.To, ", "))
	}
	if _, ok := res["Cc"]; !ok && len(msg.Cc) > 0 {
		res.Set("Cc", strings.Join(msg.Cc, ", "))
	}
	if _, ok := res["Subject"]; !ok && msg.Subject != "" {
		res.Set("Subject", msg.Subject)
	}
	if _, ok := res["Message-Id"]; !ok {
		id, err := createMessageID()
		if err != nil {
			return nil, err
		}
		res.Set("Message-Id", id)
	}
	// Date and From are required headers.
	if _, ok := res["From"]; !ok {
		res.Set("From", msg.From)
	}
	if _, ok := res["Date"]; !ok {
		res.Set("Date", time.Now().Format(time.RFC1123Z))
	}
	if _, ok := res["MIME-Version"]; !ok {
		res.Set("MIME-Version", "1.0")
	}
	for field, vals := range msg.Headers {
		if _, ok := res[field]; !ok {
			res[field] = vals
		}
	}
	return res, nil
}

func (attach *Attachment) setDefaultHeaders() {
	contentType := "application/octet-stream"
	if len(attach.ContentType) > 0 {
		contentType = attach.ContentType
	}
	attach.Header.Set("Content-Type", contentType)

	if len(attach.Header.Get("Content-Disposition")) == 0 {
		disposition := "attachment"
		attach.Header.Set("Content-Disposition", fmt.Sprintf("%s;\r\n filename=\"%s\"", disposition, attach.Filename))
	}
	if len(attach.Header.Get("Content-ID")) == 0 {
		attach.Header.Set("Content-ID", fmt.Sprintf("<%s>", attach.Filename))
	}
	if len(attach.Header.Get("Content-Transfer-Encoding")) == 0 {
		attach.Header.Set("Content-Transfer-Encoding", "base64")
	}
}

// base64Encode encodes the attachment content, and wraps it according to RFC 2045 standards (every 76 chars)
func base64Encode(w io.Writer, b []byte) {
	// 57 raw bytes per 76-byte base64 line.
	const maxRaw = 57
	// Buffer for each line, including trailing CRLF.
	buffer := make([]byte, MaxLineLength+len("\r\n"))
	copy(buffer[MaxLineLength:], "\r\n")
	// Process raw chunks until there's no longer enough to fill a line.
	for len(b) >= maxRaw {
		base64.StdEncoding.Encode(buffer, b[:maxRaw])
		w.Write(buffer)
		b = b[maxRaw:]
	}
	// Handle the last chunk of bytes.
	if len(b) > 0 {
		out := buffer[:base64.StdEncoding.EncodedLen(len(b))]
		base64.StdEncoding.Encode(out, b)
		out = append(out, "\r\n"...)
		w.Write(out)
	}
}

// field, multiple "Field: value\r\n" lines will be emitted.
func headerToBytes(buff io.Writer, header textproto.MIMEHeader) {
	for field, vals := range header {
		for _, subval := range vals {
			// bytes.Buffer.Write() never returns an error.
			io.WriteString(buff, field)
			io.WriteString(buff, ": ")
			// Write the encoded header if needed
			switch {
			case field == "Content-Type" || field == "Content-Disposition":
				buff.Write([]byte(subval))
			case field == "From" || field == "To" || field == "Cc" || field == "Bcc":
				participants := strings.Split(subval, ",")
				for i, v := range participants {
					addr, err := mail.ParseAddress(v)
					if err != nil {
						continue
					}
					participants[i] = addr.String()
				}
				buff.Write([]byte(strings.Join(participants, ", ")))
			default:
				buff.Write([]byte(mime.QEncoding.Encode("UTF-8", subval)))
			}
			io.WriteString(buff, "\r\n")
		}
	}
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

func writeMessage(buff io.Writer, msg []byte, multipart bool, mediaType string, w *multipart.Writer) error {
	if multipart {
		header := textproto.MIMEHeader{
			"Content-Type":              {mediaType + "; charset=UTF-8"},
			"Content-Transfer-Encoding": {"quoted-printable"},
		}
		if _, err := w.CreatePart(header); err != nil {
			return err
		}
	}

	qp := quotedprintable.NewWriter(buff)
	// Write the text
	if _, err := qp.Write(msg); err != nil {
		return err
	}
	return qp.Close()
}

var maxBigInt = big.NewInt(math.MaxInt64)

// The following parameters are used to generate a Message-ID:
// - The nanoseconds since Epoch
// - The calling PID
// - A cryptographically random int64
// - The sending hostname
func createMessageID() (string, error) {
	t := time.Now().UnixNano()
	pid := os.Getpid()
	rint, err := rand.Int(rand.Reader, maxBigInt)
	if err != nil {
		return "", err
	}
	h, err := os.Hostname()
	// If we can't get the hostname, we'll use localhost
	if err != nil {
		h = "localhost.localdomain"
	}
	msgid := fmt.Sprintf("<%d.%d.%d@%s>", t, pid, rint, h)
	return msgid, nil
}
