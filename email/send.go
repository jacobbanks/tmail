package email

import (
	"fmt"
	"net/smtp"
	"log"
)


func SendEmail(to, subject, body string) error {
	fromAddr := "jacobjeffreybanks@gmail.com"
	pw := "bmxr plzg fpgc cjgm"

	host := "smtp.gmail.com"
	smtpPort := "587"

	msg := []byte("To: " + to + "\r\n" + 
			"Subject: " + subject + "\r\n" + "\r\n" + 
			body + "\r\n")

	auth := smtp.PlainAuth("", fromAddr, pw, host)
	err := smtp.SendMail(host+":"+smtpPort, auth,
		fromAddr, []string{to}, msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Email Sent")
	return nil
}


// func main() {
// 	err := SendEmail("bethanyrbanks@gmail.com", "Hello, From Client App", "This is a test email /n I sent it from the Go client I'm building. /n If this worked, that's really cool!")
// 	if err != nil {
// 		log.Fatal(err)
//
// 	}
// }
