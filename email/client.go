package email

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/jacobbanks/tmail/auth"
)


func getEmails(userInfo auth.User) ([]imap.Message, error) {
	// load the token recieved from gmail servers

	c, err := client.DialTLS("imap.gmail.com:993", nil) 
	if err != nil {
		return nil, err
	}
	defer c.Logout()

	err = c.Login(userInfo.Email, userInfo.AppPassword)
	if err != nil {
		return nil, err
	}
	// Second Parameter = readOnly
	mailbox, err := c.Select("INBOX", false)
	if err != nil {
		return nil, err
	}


	set := new(imap.SeqSet)
	set.AddRange(mailbox.Messages - 9, mailbox.Messages)

	// Fetch the email
	messages := make(chan *imap.Message, 10)

	go func() {
		err := c.Fetch(set, []imap.FetchItem{imap.FetchEnvelope}, messages)
		if err != nil {
			log.Fatal(err)
		}
	}()

	var emails []imap.Message 
	for msg := range messages {
		emails = append(emails, *msg)
	}

	return emails, nil
} 

func getUserInfo() (auth.User) {
	file, err := os.Open("bin/user.bin")
	if err != nil {
		fmt.Print("Cannot open file")
	}

	dec := gob.NewDecoder(file)
	var user auth.User
	dec.Decode(&user)
	return user
}


func FetchEmails() {
	userInfo := getUserInfo()
	messages, err := getEmails(userInfo)
	if err != nil {
		log.Fatal(err)
	}

	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		if msg.Envelope != nil {
			// Extract sender details
			from := "(Unknown Sender)"
			if len(msg.Envelope.From) > 0 {
				fromAddr := msg.Envelope.From[0]
				from = fmt.Sprintf("%s <%s>", fromAddr.PersonalName, fromAddr.Address())
			}

			// Extract date
			emailDate := msg.Envelope.Date // This is a time.Time object
			formattedDate := emailDate.Format("2006-01-02 15:04:05")

			// Print email details
			fmt.Printf("From: %s\nSubject: %s\nDate: %s\n\n", from, msg.Envelope.Subject, formattedDate)
		}
	}
}
