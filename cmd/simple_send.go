package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jacobbanks/tmail/auth"
	"github.com/jacobbanks/tmail/email"
	"github.com/spf13/cobra"
)

var simpleSendCmd = &cobra.Command{
	Use:   "simple-send",
	Short: "Send an email using the command line",
	Long:  "Send an email without using the TUI (fallback mode)",
	Run:   runSimpleSend,
}

func init() {
	rootCmd.AddCommand(simpleSendCmd)
}

func runSimpleSend(cmd *cobra.Command, args []string) {
	// Check if user is authenticated
	_, err := auth.LoadUser()
	if err != nil {
		fmt.Println("You need to set up your email credentials first.")
		fmt.Println("Please run: tmail auth")
		os.Exit(1)
	}

	// Create a scanner to read user input
	scanner := bufio.NewScanner(os.Stdin)

	// Get recipient
	fmt.Print("To: ")
	scanner.Scan()
	to := scanner.Text()

	// Get CC
	fmt.Print("CC (comma separated, press Enter to skip): ")
	scanner.Scan()
	cc := scanner.Text()

	// Get BCC
	fmt.Print("BCC (comma separated, press Enter to skip): ")
	scanner.Scan()
	bcc := scanner.Text()

	// Get subject
	fmt.Print("Subject: ")
	scanner.Scan()
	subject := scanner.Text()

	// Get body
	fmt.Println("Message body (type '.' on a new line to end):")
	var bodyBuilder strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if line == "." {
			break
		}
		bodyBuilder.WriteString(line)
		bodyBuilder.WriteString("\n")
	}
	body := bodyBuilder.String()

	// Check for errors
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(1)
	}

	// Create email message
	message, err := email.NewEmailMessage()
	if err != nil {
		fmt.Printf("Error creating message: %v\n", err)
		os.Exit(1)
	}

	// Add recipients
	for _, address := range strings.Split(to, ",") {
		address = strings.TrimSpace(address)
		if address != "" {
			message.AddRecipient(address)
		}
	}

	// Add CC recipients
	for _, address := range strings.Split(cc, ",") {
		address = strings.TrimSpace(address)
		if address != "" {
			message.AddCC(address)
		}
	}

	// Add BCC recipients
	for _, address := range strings.Split(bcc, ",") {
		address = strings.TrimSpace(address)
		if address != "" {
			message.AddBCC(address)
		}
	}

	// Set subject and body
	message.Subject = subject
	message.SetTextBody(body)

	// Debug output
	if debugMode {
		fmt.Println("Debug: Sending email")
		fmt.Printf("To: %v\n", message.To)
		if len(message.Cc) > 0 {
			fmt.Printf("CC: %v\n", message.Cc)
		}
		if len(message.Bcc) > 0 {
			fmt.Printf("BCC: %v\n", message.Bcc)
		}
		fmt.Printf("Subject: %s\n", message.Subject)
		fmt.Printf("Body length: %d bytes\n", len(body))
	}

	fmt.Println("Sending email...")
	
	// Send email
	err = email.SendEmail(message)
	if err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Email sent successfully!")
}