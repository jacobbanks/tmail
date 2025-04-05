package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jacobbanks/tmail/email"
	"github.com/rivo/tview"
)

// EmailComposer implements a basic TUI for composing emails
type EmailComposer struct {
	app       *tview.Application
	form      *tview.Form
	bodyArea  *tview.TextArea
	statusBar *tview.TextView
	layout    *tview.Flex
	debugMode bool
	sending   bool
}

// NewEmailComposer creates a new email composer TUI
func NewEmailComposer(replyTo *email.Email) *EmailComposer {
	composer := &EmailComposer{
		app:       tview.NewApplication(),
		debugMode: false,
		sending:   false,
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
func (c *EmailComposer) createLayout(replyTo *email.Email) {
	// Create the form
	c.form = tview.NewForm()
	c.form.SetBorder(true)
	c.form.SetTitle(" Compose Email ")
	c.form.SetTitleAlign(tview.AlignCenter)
	c.form.SetBorderColor(tcell.ColorSteelBlue)

	// Create form fields
	c.form.AddInputField("To:", "", 50, nil, nil)
	c.form.AddInputField("Cc:", "", 50, nil, nil)
	c.form.AddInputField("Bcc:", "", 50, nil, nil)
	c.form.AddInputField("Subject:", "", 50, nil, nil)

	// Create form buttons
	c.form.AddButton("Send", func() {
		c.sendEmail()
	})
	
	c.form.AddButton("Cancel", func() {
		c.app.Stop()
	})
	
	// Create the body area
	c.bodyArea = tview.NewTextArea()
	c.bodyArea.SetBorder(true)
	c.bodyArea.SetTitle(" Message Body ")
	c.bodyArea.SetTitleAlign(tview.AlignCenter)
	c.bodyArea.SetBorderColor(tcell.ColorSteelBlue)
	c.bodyArea.SetPlaceholder("Type your message here...")
	
	// Create the status bar
	c.statusBar = tview.NewTextView()
	c.statusBar.SetDynamicColors(true)
	c.statusBar.SetTextAlign(tview.AlignCenter)
	c.statusBar.SetText("[blue]Tab[white]: Next Field | [blue]Ctrl+N[white]: Edit Body | [blue]Ctrl+S[white]: Send | [blue]Ctrl+Q[white]: Quit")
	
	// Create the layout
	c.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	
	// Create the header
	header := tview.NewTextView().
		SetText("Email Composer").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorSteelBlue)
	
	// Add items to the layout
	c.layout.AddItem(header, 1, 0, false)
	c.layout.AddItem(tview.NewBox(), 1, 0, false) // Spacing
	
	// Create content area (form + body)
	content := tview.NewFlex().SetDirection(tview.FlexRow)
	content.AddItem(c.form, 13, 0, true)
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
}

// SetDebugMode enables or disables debug mode
func (c *EmailComposer) SetDebugMode(debug bool) {
	c.debugMode = debug
}

// Run starts the email composer application
func (c *EmailComposer) Run() error {
	c.app.SetFocus(c.form)
	c.updateStatus("Form Mode")
	return c.app.SetRoot(c.layout, true).EnableMouse(true).Run()
}

// updateStatus updates the status bar text
func (c *EmailComposer) updateStatus(status string) {
	text := "[blue]Tab[white]: Next Field | [blue]Ctrl+N[white]: Body | [blue]Ctrl+S[white]: Send | [blue]Ctrl+C/Q[white]: Quit"
	if status != "" {
		text = status + " | " + text
	}
	c.statusBar.SetText(text)
}

// sendEmail collects form data and sends the email
func (c *EmailComposer) sendEmail() {
	// Prevent multiple sends
	if c.sending {
		return
	}
	
	// Update status and set sending flag
	c.sending = true
	c.updateStatus("Sending Email...")
	
	// Get field values
	toField := c.form.GetFormItem(0).(*tview.InputField)
	ccField := c.form.GetFormItem(1).(*tview.InputField)
	bccField := c.form.GetFormItem(2).(*tview.InputField)
	subjectField := c.form.GetFormItem(3).(*tview.InputField)
	
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
	message, err := email.NewEmailMessage()
	if err != nil {
		c.showError(fmt.Sprintf("Error creating message: %v", err))
		return
	}
	
	// Add recipients
	for _, to := range strings.Split(toField.GetText(), ",") {
		to = strings.TrimSpace(to)
		if to != "" {
			message.AddRecipient(to)
		}
	}
	
	// Add CC recipients
	for _, cc := range strings.Split(ccField.GetText(), ",") {
		cc = strings.TrimSpace(cc)
		if cc != "" {
			message.AddCC(cc)
		}
	}
	
	// Add BCC recipients
	for _, bcc := range strings.Split(bccField.GetText(), ",") {
		bcc = strings.TrimSpace(bcc)
		if bcc != "" {
			message.AddBCC(bcc)
		}
	}
	
	// Set subject and body
	message.Subject = strings.TrimSpace(subjectField.GetText())
	message.SetTextBody(c.bodyArea.GetText())
	
	// Debug output
	if c.debugMode {
		fmt.Println("Debug: Sending email")
		fmt.Printf("To: %v\n", message.To)
		fmt.Printf("CC: %v\n", message.Cc)
		fmt.Printf("BCC: %v\n", message.Bcc)
		fmt.Printf("Subject: %s\n", message.Subject)
		fmt.Printf("Body length: %d bytes\n", len(c.bodyArea.GetText()))
	}
	
	// Create send status label
	statusLabel := tview.NewTextView().
		SetText("Sending email...").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorWhite).
		SetBackgroundColor(tcell.ColorNavy)
	
	// Save original layout
	origLayout := c.layout
	
	// Show sending status
	c.app.SetRoot(statusLabel, true)
	c.app.Draw()
	
	// Actually send the email
	err = email.SendEmail(message)
	
	// Handle results
	if err != nil {
		// Show error and return to composer
		c.app.SetRoot(origLayout, true)
		c.app.Draw()
		c.showError(fmt.Sprintf("Failed to send email: %v", err))
	} else {
		// Show success message
		successText := tview.NewTextView().
			SetText("Email sent successfully!").
			SetTextAlign(tview.AlignCenter).
			SetTextColor(tcell.ColorWhite).
			SetBackgroundColor(tcell.ColorDarkGreen)
		
		c.app.SetRoot(successText, true)
		c.app.Draw()
		
		// Wait briefly then exit
		time.Sleep(1 * time.Second)
		c.app.Stop()
	}
}

// showError displays an error message and resets the sending state
func (c *EmailComposer) showError(message string) {
	// Create error view
	errorView := tview.NewTextView().
		SetText(message + "\n\nPress any key to continue").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorWhite).
		SetBackgroundColor(tcell.ColorDarkRed)
	
	// Show error
	c.app.SetRoot(errorView, true)
	
	// Handle key press to dismiss
	errorView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		c.app.SetRoot(c.layout, true)
		c.app.SetFocus(c.form)
		c.sending = false
		c.updateStatus("Form Mode")
		return nil
	})
}