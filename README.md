# TMail

A simple terminal-based email client for Gmail

## Features

- **Authentication**: Connect to Gmail using application password
- **Email Viewing**: 
  - Read emails with a clean terminal interface
  - Toggle between HTML and plaintext views
  - View attachment information
  - Reply to emails
- **Email Composition**:
  - Send emails with To, CC, and BCC fields
  - Support for multiple recipients (comma-separated)
  - Add file attachments
- **User Configuration**:
  - Choose from different color themes
  - Configure number of emails to fetch
  - Toggle HTML display
- **TUI Interface**:
  - Terminal User Interface with intuitive navigation
  - Keyboard shortcuts for all operations
  - Status bars with helpful information

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/jacobbanks/tmail.git
cd tmail

# Build and install
go build -o tmail main.go
go install
```

### With Go

```bash
go install github.com/jacobbanks/tmail@latest
```

### Homebrew

```bash
# Coming soon
brew tap jacobbanks/tap
brew install tmail
```

## Usage

### First-time setup

```bash
# Set up authentication with Gmail
tmail auth
```

You'll need to provide your Gmail address and an App Password (not your regular Gmail password).
[Learn how to create an App Password](https://support.google.com/accounts/answer/185833)

### Reading emails

```bash
# Read recent emails
tmail read
```

### Sending emails

```bash
# Open compose window
tmail send

# Quick send
tmail simple-send --to recipient@example.com --subject "Hello" --body "Hi there!"
```

### Configuration

```bash
# Show current config
tmail config show

# Set theme (blue, dark, light)
tmail config set theme dark

# Set default number of emails to fetch
tmail config set default_mails 100

# Toggle HTML rendering
tmail config set show_html true
```

## Keyboard Shortcuts

### Email List
- `j/k`: Navigate up/down
- `Enter`: View selected email
- `q`: Quit
- `?`: Show help

### Email Viewer
- `j/k`: Scroll up/down
- `Esc`: Return to email list
- `r`: Reply to email
- `h`: Toggle HTML/plaintext view
- `q`: Quit

### Email Composer
- `Tab`: Navigate between fields
- `Ctrl+N`: Focus body area
- `Ctrl+A`: Add attachment
- `Ctrl+S`: Send email
- `Ctrl+Q/C`: Quit

## Development

The project uses standard Go tools and practices:

```bash
# Run tests
go test ./...

# Format code
go fmt ./...

# Build the project
go build -o tmail main.go
```

## Releasing

This project uses GoReleaser for building and publishing releases:

1. Create and push a new tag:
```bash
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
```

2. GitHub Actions will automatically build and publish releases when new tags are pushed.

## License

MIT License