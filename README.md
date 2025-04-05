# tmail

A simple terminal-based email client for Gmail.

## Overview

`tmail` is a minimalist TUI email client that connects to Gmail using application passwords. It provides a clean interface for reading and sending emails from the terminal.

## Features

- Send plain text emails
- Read emails from your inbox
- Simple, intuitive interface
- Gmail integration
- Secure authentication with app passwords

## Installation

```bash
go install github.com/jacobbanks/tmail@latest
```

Or clone and build manually:

```bash
git clone https://github.com/jacobbanks/tmail.git
cd tmail
go build -o tmail main.go
```

## Gmail App Password Setup

To use `tmail`, you'll need to set up an app password for your Gmail account:

1. Go to your Google Account settings at [myaccount.google.com](https://myaccount.google.com)
2. Enable 2-Step Verification if not already enabled
3. Go to "Security" → "App passwords"
4. Select "Mail" as the app and "Other" as the device
5. Enter "tmail" as the device name
6. Generate and copy the 16-character app password

This app password will be used to authenticate with Gmail.

## Usage

First, set up your credentials:

```bash
tmail auth
```

To send an email:

```bash
tmail send
```

To read your inbox:

```bash
tmail read
```

## License

See the [LICENSE](LICENSE) file for details.