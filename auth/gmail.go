package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PromptForAuthentication prompts the user for their credentials
func PromptForAuthentication() (Credentials, error) {
	var creds Credentials
	reader := bufio.NewReader(os.Stdin)

	// Get name
	fmt.Print("Hi! Welcome to tmail, what's your name? ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to read name: %v", err)
	}
	creds.Name = strings.TrimSpace(name)

	// Get email
	fmt.Printf("Thanks %s! What is your email address? ", creds.Name)
	email, err := reader.ReadString('\n')
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to read email: %v", err)
	}
	creds.Email = strings.TrimSpace(email)

	// Get app password
	fmt.Print("What is your Gmail app password? If you haven't set one up, please see the README for instructions. ")
	password, err := reader.ReadString('\n')
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to read password: %v", err)
	}
	creds.AppPassword = strings.TrimSpace(password)

	// Confirm the information
	fmt.Printf("\nName: %s\nEmail: %s\nPassword: (hidden)\n", creds.Name, creds.Email)
	fmt.Print("Is this information correct? (y/n): ")
	confirm, err := reader.ReadString('\n')
	if err != nil {
		return Credentials{}, fmt.Errorf("failed to read confirmation: %v", err)
	}

	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(confirm)), "y") {
		return Credentials{}, fmt.Errorf("user canceled setup")
	}

	// Save credentials
	if err := SaveCredentials(creds); err != nil {
		return Credentials{}, fmt.Errorf("failed to save credentials: %v", err)
	}

	fmt.Println("Authentication setup complete!")
	return creds, nil
}
