package ui

import (
	"fmt"
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
}

// NewEmailReader creates a new email reader TUI
func NewEmailReader(emails []*email.Email) *EmailReader {
	reader := &EmailReader{
		app:         tview.NewApplication(),
		pages:       tview.NewPages(),
		emails:      emails,
		currentView: "list",
	}

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
	r.emailList.SetBorderColor(tcell.ColorSteelBlue)

	// Style the list
	r.emailList.SetMainTextColor(tcell.ColorWhite)
	r.emailList.SetSecondaryTextColor(tcell.ColorLightGray)
	r.emailList.SetSelectedTextColor(tcell.ColorBlack)
	r.emailList.SetSelectedBackgroundColor(tcell.ColorSteelBlue)
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

		// Create list item with formatted details
		text := fmt.Sprintf("%s  %s", date, subject)
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
	r.contentView.SetBorderColor(tcell.ColorSteelBlue)

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
	r.statusBar.SetText("[blue]j/k[white]: Navigate | [blue]Enter[white]: View Email | [blue]q[white]: Quit")
}

// setupMainLayout organizes the UI components into a layout
func (r *EmailReader) setupMainLayout() {
	// Create the main layout
	r.mainLayout = tview.NewFlex().SetDirection(tview.FlexRow)

	// Add a title/header
	header := tview.NewTextView()
	header.SetText("Email Reader")
	header.SetTextAlign(tview.AlignCenter)
	header.SetTextColor(tcell.ColorSteelBlue)

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
					// Move down in list
					current := r.emailList.GetCurrentItem()
					if current < r.emailList.GetItemCount()-1 {
						r.emailList.SetCurrentItem(current + 1)
					}
					return nil
				} else if r.currentView == "content" {
					// Scroll down in content view
					row, col := r.contentView.GetScrollOffset()
					r.contentView.ScrollTo(row+1, col)
					return nil
				}
			case 'k':
				if r.currentView == "list" {
					// Move up in list
					current := r.emailList.GetCurrentItem()
					if current > 0 {
						r.emailList.SetCurrentItem(current - 1)
					}
					return nil
				} else if r.currentView == "content" {
					// Scroll up in content view
					row, col := r.contentView.GetScrollOffset()
					if row > 0 {
						r.contentView.ScrollTo(row-1, col)
					}
					return nil
				}
			case 'r':
				if r.currentView == "content" {
					// Reply to email
					index := r.emailList.GetCurrentItem()
					if index >= 0 && index < len(r.emails) {
						r.replyToEmail(index)
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
	content.WriteString(fmt.Sprintf("[yellow]Subject:[white] %s\n\n", email.Subject))

	// Add a separator
	content.WriteString("[blue]" + strings.Repeat("─", 60) + "[white]\n\n")

	// Add the email body
	content.WriteString(email.Body)

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
			"q: Quit\n" +
			"?: Show this help").
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			r.pages.SwitchToPage("main")
		})

	// Style the modal
	modal.SetBorderColor(tcell.ColorSteelBlue)
	modal.SetBackgroundColor(tcell.ColorBlack)
	modal.SetTextColor(tcell.ColorWhite)
	modal.SetButtonBackgroundColor(tcell.ColorSteelBlue)
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
		r.statusBar.SetText("[blue]j/k[white]: Scroll | [blue]Esc[white]: Back to List | [blue]r[white]: Reply | [blue]q[white]: Quit")
	}
}

// Run starts the email reader application
func (r *EmailReader) Run() error {
	// Initially focus on the email list
	r.app.SetFocus(r.emailList)

	// Start the application
	return r.app.SetRoot(r.pages, true).EnableMouse(true).Run()
}
