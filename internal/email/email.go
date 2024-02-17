package email

import (
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/config"
	"gopkg.in/gomail.v2"
	"log"
	"strconv"
)

func SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "google_workspace_integration@vatusa.net")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	if config.DEBUG {
		return nil
	}
	mailPort, err := strconv.Atoi(config.MAIL_PORT)
	d := gomail.NewDialer(config.MAIL_HOST, mailPort, config.MAIL_USERNAME, config.MAIL_PASSWORD)
	err = d.DialAndSend(m)
	if err != nil {
		return err
	}

	log.Printf("Email %s sent to %s\n", subject, to)
	return nil
}

func SendNewAccountEmail(to string, firstName string, lastName string, primaryEmail string, password string) error {
	message := fmt.Sprintf("Hello %s %s,\n\n"+
		"Your VATUSA google account password has been newly created for you.\n"+
		"Your VATUSA email address is: %s\n Your new password is: \n%s\n "+
		"This password must be changed on your next login."+
		"Make sure you set up two factor authentication within one week to avoid getting locked out of your account.",
		firstName, lastName, primaryEmail, password)

	return SendEmail(to, "VATUSA Google Account - Account Created", message)
}

func SendPasswordResetEmail(to string, firstName string, lastName string, primaryEmail string, password string) error {
	message := fmt.Sprintf("Hello %s %s,\n\n"+
		"Your VATUSA google account password has been reset.\n"+
		"Your VATUSA email address is: %s\n Your new password is: \n%s\n "+
		"This password must be changed on your next login.",
		firstName, lastName, primaryEmail, password)

	return SendEmail(to, "VATUSA Google Account - Password Reset", message)
}
