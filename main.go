package main

import (
	"fmt"
	// "log"

	"github.com/jacobbanks/tmail/auth"
	// "github.com/rivo/tview"
)

func main() {
	fmt.Println("Starting authentication process...")
	auth.PromptForAuthentication()

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
