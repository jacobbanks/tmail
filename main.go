package main

import (
	"fmt"
	// "log"

	"github.com/jacobbanks/tmail/auth"
	"github.com/jacobbanks/tmail/email"
	// "github.com/rivo/tview"
)

func main() {
	fmt.Println("Starting authentication process...")
	authenticated := auth.PromptForAuthentication()
	if !authenticated {
		fmt.Errorf("Auth Failed")
	}


	email.FetchEmails()
	// email.GetEmailsFromGmail()



	// Example: Get the user's Gmail profile

	// app := tview.NewApplication()

	//
	// // creating a simple box
	//
	// box := tview.NewBox().SetBorder(true).SetTitle("Welcome to tmail")
	//
	// if err := app.SetRoot(box, true).Run(); err != nil {
	// 	panic(err)
	// }
}
