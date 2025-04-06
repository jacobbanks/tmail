package ui

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jacobbanks/tmail/email"
	"github.com/rivo/tview"
)

// EmailReader implements a TUI for reading emails
type EmailReader struct {
	app         *tview.Application
	pages       *tview.Pages
	emails      []*email.Email
	mainLayout  *tview.Flex
	emailList   *tview.List
	contentView *tview.TextView
	statusBar   *tview.TextView
	currentView string // "list" or "content"
	showHTML    bool   // whether to show HTML content
}

// NewEmailReader creates a new email reader TUI
func NewEmailReader(emails []*email.Email) *EmailReader {
	reader := &EmailReader{
		app:         tview.NewApplication(),
		pages:       tview.NewPages(),
		emails:      emails,
		currentView: "list",
	}

	// Apply user configuration
	config, _ := email.LoadUserConfig()
	reader.showHTML = config.ShowHTML

	reader.setupUI()
	return reader
}

// setupUI initializes all UI components
func (r *EmailReader) setupUI() {
	r.setupEmailList()
	r.setupContentView()
	r.setupStatusBar()
	r.setupMainLayout()
	r.setupKeybindings()

	// Add "main" page to the pages component
	r.pages.AddPage("main", r.mainLayout, true, true)
}

// setupEmailList creates and configures the email list
func (r *EmailReader) setupEmailList() {
	r.emailList = tview.NewList()
	r.emailList.SetBorder(true)
	r.emailList.SetTitle(" Inbox ")
	r.emailList.SetTitleAlign(tview.AlignCenter)

	// Apply theme
	config, _ := email.LoadUserConfig()
	var borderColor, selectedBgColor tcell.Color
	switch config.Theme {
	case "dark":
		borderColor = tcell.ColorDarkBlue
		selectedBgColor = tcell.ColorDarkBlue
	case "light":
		borderColor = tcell.ColorLightBlue
		selectedBgColor = tcell.ColorLightBlue
	default:
		borderColor = tcell.ColorSteelBlue
		selectedBgColor = tcell.ColorSteelBlue
	}
	r.emailList.SetBorderColor(borderColor)

	// Style the list
	r.emailList.SetMainTextColor(tcell.ColorWhite)
	r.emailList.SetSecondaryTextColor(tcell.ColorLightGray)
	r.emailList.SetSelectedTextColor(tcell.ColorBlack)
	r.emailList.SetSelectedBackgroundColor(selectedBgColor)
	r.emailList.SetHighlightFullLine(true)
	r.emailList.SetWrapAround(false)

	// Add emails to the list
	for i, email := range r.emails {
		// Format the date for display
		date := email.Date.Format("2006-01-02 15:04")

		// Format the subject (truncate if too long)
		subject := email.Subject
		if len(subject) > 40 {
			subject = subject[:37] + "..."
		}

		// Format the sender (extract just the name or email address)
		sender := email.From
		if idx := strings.LastIndex(sender, "<"); idx > 0 {
			sender = strings.TrimSpace(sender[:idx])
		}
		if len(sender) > 25 {
			sender = sender[:22] + "..."
		}

		// Add attachment indicator if needed
		attachmentIndicator := ""
		if len(email.Attachments) > 0 {
			attachmentIndicator = "📎 "
		}

		// Create list item with formatted details
		text := fmt.Sprintf("%s  %s%s", date, attachmentIndicator, subject)
		secondaryText := fmt.Sprintf("From: %s", sender)

		// Create a fixed value for the closure to capture
		index := i
		r.emailList.AddItem(text, secondaryText, rune('a'+i), func() {
			r.showEmail(index)
		})
	}
}

// setupContentView creates and configures the email content view
func (r *EmailReader) setupContentView() {
	r.contentView = tview.NewTextView()
	r.contentView.SetBorder(true)
	r.contentView.SetTitle(" Email Content ")
	r.contentView.SetTitleAlign(tview.AlignCenter)

	// Apply theme
	config, _ := email.LoadUserConfig()
	var borderColor tcell.Color
	switch config.Theme {
	case "dark":
		borderColor = tcell.ColorDarkBlue
	case "light":
		borderColor = tcell.ColorLightBlue
	default:
		borderColor = tcell.ColorSteelBlue
	}
	r.contentView.SetBorderColor(borderColor)

	r.contentView.SetDynamicColors(true)
	r.contentView.SetRegions(true)
	r.contentView.SetWordWrap(true)
	r.contentView.SetScrollable(true)
}

// setupStatusBar creates and configures the status bar
func (r *EmailReader) setupStatusBar() {
	r.statusBar = tview.NewTextView()
	r.statusBar.SetDynamicColors(true)
	r.statusBar.SetTextAlign(tview.AlignCenter)
	r.statusBar.SetText("[blue]j/k[white]: Navigate | [blue]Enter[white]: View Email | [blue]r[white]: Reply | [blue]q[white]: Quit")
}

// setupMainLayout organizes the UI components into a layout
func (r *EmailReader) setupMainLayout() {
	// Create the main layout
	r.mainLayout = tview.NewFlex().SetDirection(tview.FlexRow)

	// Get theme color
	config, _ := email.LoadUserConfig()
	var headerColor tcell.Color
	switch config.Theme {
	case "dark":
		headerColor = tcell.ColorDarkBlue
	case "light":
		headerColor = tcell.ColorLightBlue
	default:
		headerColor = tcell.ColorSteelBlue
	}

	// Add a title/header
	header := tview.NewTextView()
	header.SetText("Email Reader")
	header.SetTextAlign(tview.AlignCenter)
	header.SetTextColor(headerColor)

	// Create a flex for email list and content view
	contentArea := tview.NewFlex()
	contentArea.AddItem(r.emailList, 0, 1, true)
	contentArea.AddItem(r.contentView, 0, 2, false)

	// Create a centered content layout
	centeredFlex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(contentArea, 0, 1, true),
			0, 3, true,
		).
		AddItem(nil, 0, 1, false)

	// Add components to the layout with proper spacing
	r.mainLayout.AddItem(header, 1, 0, false)
	r.mainLayout.AddItem(tview.NewBox(), 1, 0, false) // Spacing
	r.mainLayout.AddItem(centeredFlex, 0, 1, true)    // Content takes remaining space
	r.mainLayout.AddItem(r.statusBar, 1, 0, false)    // Status bar at bottom
}

// setupKeybindings configures global keyboard shortcuts
func (r *EmailReader) setupKeybindings() {
	r.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Global shortcuts
		switch event.Key() {
		case tcell.KeyEscape:
			if r.currentView == "content" {
				r.currentView = "list"
				r.app.SetFocus(r.emailList)
				r.updateStatusBar()
				return nil
			}
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				r.app.Stop()
				return nil
			case '?':
				r.showHelp()
				return nil
			case 'j':
				if r.currentView == "list" {
					current := r.emailList.GetCurrentItem()
					if current < r.emailList.GetItemCount()-1 {
						r.emailList.SetCurrentItem(current + 1)
					}
					return nil
				} else if r.currentView == "content" {
					row, col := r.contentView.GetScrollOffset()
					r.contentView.ScrollTo(row+1, col)
					return nil
				}
			case 'k':
				if r.currentView == "list" {
					current := r.emailList.GetCurrentItem()
					if current > 0 {
						r.emailList.SetCurrentItem(current - 1)
					}
					return nil
				} else if r.currentView == "content" {
					row, col := r.contentView.GetScrollOffset()
					if row > 0 {
						r.contentView.ScrollTo(row-1, col)
					}
					return nil
				}
			case 'r':
				if r.currentView == "content" {
					index := r.emailList.GetCurrentItem()
					if index >= 0 && index < len(r.emails) {
						r.replyToEmail(index)
					}
					return nil
				}
			case 'h':
				// Toggle between HTML and plain text view
				if r.currentView == "content" {
					r.showHTML = !r.showHTML
					index := r.emailList.GetCurrentItem()
					if index >= 0 && index < len(r.emails) {
						r.showEmail(index)
					}
					return nil
				}
			}
		}
		return event
	})
}

// showEmail displays the selected email in the content view
func (r *EmailReader) showEmail(index int) {
	if index < 0 || index >= len(r.emails) {
		return
	}

	email := r.emails[index]

	// Format the email for display
	var content strings.Builder

	// Add header information with colors
	content.WriteString(fmt.Sprintf("[yellow]From:[white] %s\n", email.From))
	content.WriteString(fmt.Sprintf("[yellow]To:[white] %s\n", email.To))
	content.WriteString(fmt.Sprintf("[yellow]Date:[white] %s\n", email.Date.Format(time.RFC1123Z)))
	content.WriteString(fmt.Sprintf("[yellow]Subject:[white] %s\n", email.Subject))

	// Add attachment information if present
	if len(email.Attachments) > 0 {
		content.WriteString(fmt.Sprintf("[yellow]Attachments:[white] %s\n", strings.Join(email.Attachments, ", ")))
	}

	// Add a separator
	content.WriteString("\n[blue]" + strings.Repeat("─", 60) + "[white]\n\n")

	// Add the email body (either HTML-converted or plain text)
	emailText := ""
	if r.showHTML && email.HTMLBody != "" {
		emailText = email.Body
	} else {
		emailText = email.Body
	}

	// Highlight links for better readability
	if r.showHTML && email.IsHTML {
		emailText = highlightLinks(emailText)
	}

	content.WriteString(emailText)

	// Set the content view text
	r.contentView.SetText(content.String())
	r.contentView.ScrollToBeginning()

	// Update the view state and focus
	r.currentView = "content"
	r.app.SetFocus(r.contentView)

	// Update the status bar
	r.updateStatusBar()
}

// replyToEmail opens a composer to reply to the selected email
func (r *EmailReader) replyToEmail(index int) {
	if index < 0 || index >= len(r.emails) {
		return
	}

	// Stop the current application
	r.app.Stop()

	// Create and run a new email composer in reply mode
	composer := NewEmailComposer(r.emails[index])
	composer.Run()
}

// showHelp displays help information
func (r *EmailReader) showHelp() {
	modal := tview.NewModal().
		SetText("Keyboard Shortcuts:\n\n" +
			"j/k: Navigate up/down\n" +
			"Enter: View selected email\n" +
			"Esc: Return to email list\n" +
			"r: Reply to current email\n" +
			"h: Toggle HTML/plain text view\n" +
			"q: Quit\n" +
			"?: Show this help").
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			r.pages.SwitchToPage("main")
		})

	// Get theme color
	config, _ := email.LoadUserConfig()
	var borderColor tcell.Color
	switch config.Theme {
	case "dark":
		borderColor = tcell.ColorDarkBlue
	case "light":
		borderColor = tcell.ColorLightBlue
	default:
		borderColor = tcell.ColorSteelBlue
	}

	// Style the modal
	modal.SetBorderColor(borderColor)
	modal.SetBackgroundColor(tcell.ColorBlack)
	modal.SetTextColor(tcell.ColorWhite)
	modal.SetButtonBackgroundColor(borderColor)
	modal.SetButtonTextColor(tcell.ColorWhite)

	// Add the modal to pages and show it
	r.pages.AddPage("modal", modal, true, true)
	r.app.SetFocus(modal)
}

// updateStatusBar updates the status bar based on the current view
func (r *EmailReader) updateStatusBar() {
	if r.currentView == "list" {
		r.statusBar.SetText("[blue]j/k[white]: Navigate | [blue]Enter[white]: View Email | [blue]q[white]: Quit")
	} else {
		htmlStatus := ""
		if r.emails[r.emailList.GetCurrentItem()].HTMLBody != "" {
			if r.showHTML {
				htmlStatus = "[blue]h[white]: Show Plain Text | "
			} else {
				htmlStatus = "[blue]h[white]: Show HTML | "
			}
		}
		r.statusBar.SetText("[blue]j/k[white]: Scroll | " + htmlStatus + "[blue]Esc[white]: Back to List | [blue]r[white]: Reply | [blue]q[white]: Quit")
	}
}

// highlightLinks applies color formatting to URLs in text
func highlightLinks(text string) string {
	// Regular expression to match URLs in the text, including those in parentheses
	urlRegex := regexp.MustCompile(`\((https?://[^\s)]+)\)`)

	// Replace each URL with a colored version using tview's color tags
	coloredText := urlRegex.ReplaceAllString(text, "([cyan]$1[white])")

	return coloredText
}

// Run starts the email reader application
func (r *EmailReader) Run() error {
	// Initially focus on the email list
	r.app.SetFocus(r.emailList)

	// Start the application
	return r.app.SetRoot(r.pages, true).EnableMouse(true).Run()
}
