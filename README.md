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

## Installation

### Quick Install

```bash
# Using Go (requires Go 1.16+)
go install github.com/jacobbanks/tmail@latest

# Using the Makefile (after cloning the repo)
make install        # Install to $GOPATH/bin
make install-global # Install to /usr/local/bin (requires sudo)
```

### Download Binaries

Pre-built binaries are available on the [Releases page](https://github.com/jacobbanks/tmail/releases).

### Build from Source

```bash
# Clone the repository
git clone https://github.com/jacobbanks/tmail.git
cd tmail

# Build
make build

# Run tests
make test
```

## Usage

### First-time Setup

Before using tmail, you need to set up authentication with Gmail:

```bash
tmail auth
```

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

### Makefile Commands
tmail includes a Makefile with useful commands:

```bash
make build           # Build the binary
make test            # Run tests
make install         # Install to $GOPATH/bin
make install-global  # Install to /usr/local/bin
make clean           # Remove build artifacts
make fmt             # Format code
make lint            # Run linter
make release         # Build release binaries
```

### Release Process

1. Update version in `version.go`
2. Create a tag: `git tag -a v0.1.0 -m "Version 0.1.0"`
3. Push the tag: `git push origin v0.1.0`
4. GitHub Actions will automatically build and create a release

### License

Distributed under the MIT License. See `LICENSE` for more information.
