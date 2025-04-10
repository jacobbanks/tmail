package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jacobbanks/tmail/email"
	"github.com/rivo/tview"
)

// EmailComposer implements a basic TUI for composing emails
type EmailComposer struct {
	app         *tview.Application
	form        *tview.Form
	bodyArea    *tview.TextArea
	statusBar   *tview.TextView
	layout      *tview.Flex
	pages       *tview.Pages
	attachments []string
	debugMode   bool
	sending     bool
	provider    email.MailProvider
}

// Form field indices for the email composer
const (
	// ToField is the index of the To field in the form
	ToField = 0
	// CcField is the index of the CC field in the form
	CcField = 1
	// BccField is the index of the BCC field in the form
	BccField = 2
	// SubjectField is the index of the Subject field in the form
	SubjectField = 3
	// AttachmentField is the index of the Attachment field in the form
	AttachmentField = 4
)

// NewEmailComposer creates a new email composer TUI
func NewEmailComposer(replyTo *email.IncomingMessage, provider email.MailProvider) *EmailComposer {
	composer := &EmailComposer{
		app:         tview.NewApplication(),
		pages:       tview.NewPages(),
		debugMode:   false,
		sending:     false,
		attachments: []string{},
		// eventually refactor to support multiple provider based on config
		provider: provider,
	}

	// Create form and layout
	composer.createLayout(replyTo)

	// Set up keyboard shortcuts
	composer.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if composer.sending {
			// Block all input while sending
			return nil
		}

		switch event.Key() {
		case tcell.KeyCtrlS: // Ctrl+S to send
			composer.sendEmail()
			return nil
		case tcell.KeyCtrlQ, tcell.KeyCtrlC: // Ctrl+Q or Ctrl+C to quit
			composer.app.Stop()
			return nil
		case tcell.KeyCtrlN: // Ctrl+N to focus on body
			composer.app.SetFocus(composer.bodyArea)
			composer.updateStatus("Body Mode")
			return nil
		case tcell.KeyCtrlA: // Ctrl+A to add attachment
			composer.showAttachmentDialog()
			return nil
		case tcell.KeyEscape: // Escape to return to form from body
			if composer.bodyArea.HasFocus() {
				composer.app.SetFocus(composer.form)
				composer.updateStatus("Form Mode")
				return nil
			}
		}
		return event
	})

	return composer
}

// createLayout sets up the UI components
func (c *EmailComposer) createLayout(replyTo *email.IncomingMessage) {
	// Create the form
	c.form = tview.NewForm()
	c.form.SetBorder(true)
	c.form.SetTitle(" Compose Email ")
	c.form.SetTitleAlign(tview.AlignCenter)
	c.form.SetBorderColor(tcell.ColorSteelBlue)

	// Get user configuration for theme
	config, _ := email.LoadUserConfig()
	var primaryColor tcell.Color
	switch config.Theme {
	case "dark":
		primaryColor = tcell.ColorDarkBlue
	case "light":
		primaryColor = tcell.ColorLightBlue
	default:
		primaryColor = tcell.ColorSteelBlue
	}
	c.form.SetBorderColor(primaryColor)

	// Create form fields with help text for multiple addresses
	c.form.AddInputField("To: (separate multiple addresses with commas)", "", 50, nil, nil)
	c.form.AddInputField("Cc: (separate multiple addresses with commas)", "", 50, nil, nil)
	c.form.AddInputField("Bcc: (separate multiple addresses with commas)", "", 50, nil, nil)
	c.form.AddInputField("Subject:", "", 50, nil, nil)

	attachText := "None"
	if len(c.attachments) > 0 {
		var fileNames []string
		for _, path := range c.attachments {
			fileNames = append(fileNames, filepath.Base(path))
		}
		attachText = strings.Join(fileNames, ", ")
	}
	c.form.AddInputField("Attachments:", attachText, 50, nil, nil)
	c.form.GetFormItem(AttachmentField).(*tview.InputField).SetDisabled(true)

	// Create form buttons
	c.form.AddButton("Send", func() {
		c.sendEmail()
	})

	c.form.AddButton("Attach", func() {
		c.showAttachmentDialog()
	})

	c.form.AddButton("Cancel", func() {
		c.app.Stop()
	})

	// Create the body area
	c.bodyArea = tview.NewTextArea()
	c.bodyArea.SetBorder(true)
	c.bodyArea.SetTitle(" Message Body ")
	c.bodyArea.SetTitleAlign(tview.AlignCenter)
	c.bodyArea.SetBorderColor(primaryColor)
	c.bodyArea.SetPlaceholder("Type your message here...")

	// Create the status bar
	c.statusBar = tview.NewTextView()
	c.statusBar.SetDynamicColors(true)
	c.statusBar.SetTextAlign(tview.AlignCenter)
	c.statusBar.SetText("[blue]Tab[white]: Next Field | [blue]Ctrl+N[white]: Edit Body | [blue]Ctrl+A[white]: Add Attachment | [blue]Ctrl+S[white]: Send | [blue]Ctrl+Q[white]: Quit")

	// Create the layout
	c.layout = tview.NewFlex().SetDirection(tview.FlexRow)

	// Create the header
	header := tview.NewTextView().
		SetText("Email Composer").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(primaryColor)

	// Add items to the layout
	c.layout.AddItem(header, 1, 0, false)
	c.layout.AddItem(tview.NewBox(), 1, 0, false) // Spacing

	// Create content area (form + body)
	content := tview.NewFlex().SetDirection(tview.FlexRow)
	content.AddItem(c.form, 14, 0, true) // Increased height for new attachment field
	content.AddItem(c.bodyArea, 0, 1, false)

	// Create centered flex
	centered := tview.NewFlex()
	centered.AddItem(nil, 0, 1, false)
	centered.AddItem(content, 0, 3, true)
	centered.AddItem(nil, 0, 1, false)

	c.layout.AddItem(centered, 0, 1, true)
	c.layout.AddItem(c.statusBar, 1, 0, false)

	// Pre-fill form if replying
	if replyTo != nil {
		// Set To field to original sender
		c.form.GetFormItem(0).(*tview.InputField).SetText(replyTo.From)

		// Set Subject with Re: prefix if needed
		subject := replyTo.Subject
		if !strings.HasPrefix(strings.ToLower(subject), "re:") {
			subject = "Re: " + subject
		}
		c.form.GetFormItem(3).(*tview.InputField).SetText(subject)

		// Add reply content to body
		replyBody := "\n\n-------- Original Message --------\n"
		replyBody += "From: " + replyTo.From + "\n"
		replyBody += "Date: " + replyTo.Date.Format("Mon, 02 Jan 2006 15:04:05 -0700") + "\n"
		replyBody += "Subject: " + replyTo.Subject + "\n\n"
		replyBody += replyTo.Body

		c.bodyArea.SetText(replyBody, true)
	}

	// Add main page to pages
	c.pages.AddPage("main", c.layout, true, true)
}

// updateAttachmentField updates the attachment field display
func (c *EmailComposer) updateAttachmentField() {
	// Update attachment display
	attachText := "None"
	if len(c.attachments) > 0 {
		var fileNames []string
		for _, path := range c.attachments {
			fileNames = append(fileNames, filepath.Base(path))
		}
		attachText = strings.Join(fileNames, ", ")
	}
	c.form.GetFormItem(AttachmentField).(*tview.InputField).SetText(attachText)
}

// SetDebugMode enables or disables debug mode
func (c *EmailComposer) SetDebugMode(debug bool) {
	c.debugMode = debug
}

// Run starts the email composer application
func (c *EmailComposer) Run() error {
	c.app.SetFocus(c.form)
	c.updateStatus("Form Mode")
	c.provider.Connect()
	return c.app.SetRoot(c.pages, true).EnableMouse(true).Run()
}

// updateStatus updates the status bar text
func (c *EmailComposer) updateStatus(status string) {
	text := "[blue]Tab[white]: Next Field | [blue]Ctrl+N[white]: Body | [blue]Ctrl+A[white]: Attach | [blue]Ctrl+S[white]: Send | [blue]Ctrl+C/Q[white]: Quit"
	if status != "" {
		text = status + " | " + text
	}
	c.statusBar.SetText(text)
}

// showAttachmentDialog displays a dialog to input an attachment file path
func (c *EmailComposer) showAttachmentDialog() {
	// Create a form for attachment input
	form := tview.NewForm()
	form.AddInputField("File path:", "", 40, nil, nil)
	form.AddButton("Attach", func() {
		// Get the input field and its text
		field := form.GetFormItem(0).(*tview.InputField)
		filePath := field.GetText()
		if filePath != "" {
			c.addAttachment(filePath)
		}
		c.pages.SwitchToPage("main")
	})
	form.AddButton("Cancel", func() {
		c.pages.SwitchToPage("main")
	})

	// Create a frame for the form
	frame := tview.NewFrame(form).
		SetBorders(1, 1, 1, 1, 2, 2).
		AddText("Add Attachment", true, tview.AlignCenter, tcell.ColorWhite).
		AddText("Enter the path to the file you want to attach", false, tview.AlignCenter, tcell.ColorWhite)

	// Add the frame to pages
	c.pages.AddPage("attachment", frame, true, false)
	c.pages.SwitchToPage("attachment")
	c.app.SetFocus(form)
}

// addAttachment adds a file to the list of attachments
func (c *EmailComposer) addAttachment(filePath string) {
	// Expand tilde to home directory if present
	if strings.HasPrefix(filePath, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			filePath = filepath.Join(home, filePath[2:])
		}
	}

	// Check if file exists
	if _, err := os.Stat(filePath); err != nil {
		c.showError(fmt.Sprintf("File not found: %s", filePath))
		return
	}

	// Add to attachments
	c.attachments = append(c.attachments, filePath)

	// Update attachment field
	c.updateAttachmentField()
}

// sendEmail collects form data and sends the email
func (c *EmailComposer) sendEmail() {
	// Prevent multiple sends
	if c.sending {
		return
	}

	// Update status and set sending flag
	c.sending = true

	// Get field values
	toField := c.form.GetFormItem(ToField).(*tview.InputField)
	ccField := c.form.GetFormItem(CcField).(*tview.InputField)
	bccField := c.form.GetFormItem(BccField).(*tview.InputField)
	subjectField := c.form.GetFormItem(SubjectField).(*tview.InputField)

	// Validate required fields
	if toField.GetText() == "" {
		c.showError("Error: Recipient (To) is required")
		return
	}

	if subjectField.GetText() == "" {
		c.showError("Error: Subject is required")
		return
	}

	// Create and populate email message
	message, err := email.NewOutgoingMessage()
	if err != nil {
		c.showError(fmt.Sprintf("Error creating message: %v", err))
		return
	}

	// Add recipients
	toAddresses := strings.Split(toField.GetText(), ",")
	for _, to := range toAddresses {
		to = strings.TrimSpace(to)
		if to != "" {
			message.AddRecipient(to)
		}
	}

	// Add CC recipients
	ccAddresses := strings.Split(ccField.GetText(), ",")
	for _, cc := range ccAddresses {
		cc = strings.TrimSpace(cc)
		if cc != "" {
			message.AddCC(cc)
		}
	}

	// Add BCC recipients
	bccAddresses := strings.Split(bccField.GetText(), ",")
	for _, bcc := range bccAddresses {
		bcc = strings.TrimSpace(bcc)
		if bcc != "" {
			message.AddBCC(bcc)
		}
	}

	// Set subject and body
	message.Subject = strings.TrimSpace(subjectField.GetText())
	message.SetTextBody(c.bodyArea.GetText())

	// Add attachments
	for _, path := range c.attachments {
		if err := message.AppendAttachmentPath(path); err != nil {
			c.showError(fmt.Sprintf("Error adding attachment %s: %v", filepath.Base(path), err))
			return
		}
	}

	if c.debugMode {
		fmt.Println("Debug: Sending email")
		fmt.Printf("To: %v\n", message.To)
		fmt.Printf("CC: %v\n", message.Cc)
		fmt.Printf("BCC: %v\n", message.Bcc)
		fmt.Printf("Subject: %s\n", message.Subject)
		fmt.Printf("Attachments: %d\n", len(message.Attachments))
		fmt.Printf("Body length: %d bytes\n", len(c.bodyArea.GetText()))
	}

	// Create send status label
	statusLabel := tview.NewTextView()
	statusLabel.SetTextAlign(tview.AlignCenter)
	statusLabel.SetTextColor(tcell.ColorWhite)
	statusLabel.SetBackgroundColor(tcell.ColorNavy)
	statusLabel.SetText("Sending email...")

	// Save original layout
	origLayout := c.layout

	// Show sending status
	c.app.SetRoot(statusLabel, true)

	// Actually send the email
	err = c.provider.SendEmail(message)

	// Handle results
	if err != nil {
		// Show error and return to composer
		c.app.SetRoot(origLayout, true)
		c.showError(fmt.Sprintf("Failed to send email: %v", err))
	} else {
		// Show success message
		successText := tview.NewTextView()
		successText.SetTextAlign(tview.AlignCenter)
		successText.SetTextColor(tcell.ColorWhite)
		successText.SetBackgroundColor(tcell.ColorDarkGreen)
		successText.SetText("Email sent successfully!")

		c.app.SetRoot(successText, true)

		// Wait briefly then exit
		time.Sleep(1 * time.Second)
		c.app.Stop()
	}
}

// showError displays an error message and resets the sending state
func (c *EmailComposer) showError(message string) {
	// Create error view
	errorView := tview.NewTextView()
	errorView.SetTextAlign(tview.AlignCenter)
	errorView.SetTextColor(tcell.ColorWhite)
	errorView.SetBackgroundColor(tcell.ColorDarkRed)
	errorView.SetText(message + "\n\nPress any key to continue")

	// Show error
	c.app.SetRoot(errorView, true)

	// Handle key press to dismiss
	errorView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		c.app.SetRoot(c.pages, true)
		c.app.SetFocus(c.form)
		c.sending = false
		c.updateStatus("Form Mode")
		return nil
	})
}
