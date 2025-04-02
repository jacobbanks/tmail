package auth

import (
	"bufio"
	// "bytes"
	// "context"
	"encoding/json"
	"fmt"
	// "io"
	// "net/http"
	"os"
	// "strings"
	// "time"
	// "bytes"
	"encoding/gob"
	"golang.org/x/oauth2"
)
// bmxr plzg fpgc cjgm
// Google OAuth2 Endpoints for Device Flow


type User struct {
	AppPassword string
	Email string
	Name string
}

// Ask the user if they want to authenticate with Google
func PromptForAuthentication() bool {
	fmt.Print("Hi! Welcome to tmail, what's your name?")

	var userInfo User
	// Read user input
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	userInfo.Name = input

	if userInfo.Name == input {
		fmt.Printf("Thanks %v!, Next, what is your email address? ", userInfo.Name)
	} else {
		fmt.Print("I am sorry. I didn't get that.")
	}

	input, _ = reader.ReadString('\n')
	userInfo.Email = input

	if userInfo.Email == input {
		fmt.Printf("Got it: %v, Next, what is your gmail app password? If you haven't configured one yet, then please see the ReadMe for instructions on how to setup an app password for gmail.", userInfo.Email)
	}

	input, _ = reader.ReadString('\n')
	userInfo.AppPassword = input
	if input == userInfo.AppPassword {
		fmt.Printf("Thanks! Here is your data. %v", userInfo)
	}
	writeUserToFile(userInfo)
	return true
	 // Normalize input
	//
	// // If user says yes, start authentication flow
	// if input == "y" || input == "yes" {
	// 	fmt.Println("Starting Google authentication...")
	// 	token, err := Authenticate()
	// 	if err != nil {
	// 		fmt.Println("Authentication failed:", err)
	// 		return false
	// 	}
	//
	// 	fmt.Println("Successfully authenticated! Access Token:", token)
	// 	return true
	// } else {
	// 	fmt.Println("Skipping authentication.")
	// 	return false
	// }
}



// saveToken stores the OAuth token in a local file
func saveToken(token *oauth2.Token) {
	file, err := os.Create("token.json")
	if err != nil {
		fmt.Println("Unable to save token:", err)
		return
	}
	defer file.Close()
	print("This is your token here %v", token)

	json.NewEncoder(file).Encode(token)
}

func writeUserToFile(userinfo User) {
	file, err := os.Create("bnin/user.bin")
	if err != nil {
		fmt.Println("Unable to create file", err)
		return
	}
	err = gob.NewEncoder(file).Encode(userinfo)
	if err != nil {
		fmt.Errorf("failed to encode data, %v", err)
	}
}

