# TMail

A simple terminal-based email client for Gmail

## Features

- Authenticate with Gmail using application password
- Send emails with To, CC, and BCC fields
- Read emails with a clean terminal interface
- Reply to emails
- TUI (Terminal User Interface) experience

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