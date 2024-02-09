package services

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
)

func SendMail(body bytes.Buffer, recipient string) error {
	from := os.Getenv("SMTP_SENDER")
	password := os.Getenv("SMTP_PASSWORD")

	to := []string{
		recipient,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	addr := fmt.Sprintf("%v:%v", smtpHost, smtpPort)
	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(addr, auth, from, to, body.Bytes())
	if err != nil {
		return err
	}

	return nil
}
