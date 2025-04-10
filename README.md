# TMail

A terminal-based email client for Gmail.

## Features
- **Gmail Integration**
  - Connect securely with application password
  - Read emails with a clean, navigable interface
  - Send emails with attachments
  - Reply to conversations

- **Terminal UI**
  - Intuitive keyboard-driven navigation
  - Email list with sender, subject, and date
  - Message view with formatted content
  - Composer with multiple recipient support

### Quick Install
Run this in your terminal (macOS/Linux):

```bash
curl -sSf https://raw.githubusercontent.com/jacobbanks/tmail/main/install.sh | sh
```

### First-time Setup
Before using tmail, you need a Gmail App Password.
Visit https://myaccount.google.com/security
Enable 2-Step Verification
Click on App Passwords
Under Select App, choose Mail
Under Select Device, choose Other (Custom name) and enter tmail
Click Generate
Copy the 16-character password shown
tmail auth

You'll be prompted to enter:
1. Your Gmail address
2. An App Password (not your regular Gmail password)

> 🔑 **Security Note**: tmail stores credentials locally on your machine. For Gmail, you must create an [App Password](https://support.google.com/accounts/answer/185833).

### Reading Emails

```bash
# Read most recent emails
tmail read
```

### Sending Emails

```bash
# Open the email composer
tmail send

# Quick send from command line
tmail simple-send --to user@example.com --subject "Hello" --body "This is a test email"

# Include attachments
tmail simple-send --to user@example.com --subject "With attachment" --body "See attached file" --attach path/to/file.pdf
```

### Version Information

```bash
tmail version
```

## Keyboard Shortcuts

### Email List View
- `j/k`: Navigate down/up
- `Enter`: Open selected email
- `q`: Quit

### Email Content View
- `j/k`: Scroll down/up
- `Esc`: Return to list view
- `r`: Reply to email
- `q`: Quit

### Email Composer
- `Tab`: Navigate between fields
- `Ctrl+N`: Focus body content
- `Ctrl+A`: Add attachment
- `Ctrl+S`: Send email
- `Ctrl+Q/C`: Quit without sending


### License
Distributed under the MIT License. See `LICENSE` for more information.
