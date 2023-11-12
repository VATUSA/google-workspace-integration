package email

import (
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/config"
	"gopkg.in/gomail.v2"
	"strconv"
)

func SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "google_workspace_integration@vatusa.net")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	mailPort, err := strconv.Atoi(config.MAIL_PORT)
	d := gomail.NewDialer(config.MAIL_HOST, mailPort, config.MAIL_USERNAME, config.MAIL_PASSWORD)
	err = d.DialAndSend(m)
	if err != nil {
		return err
	}

	fmt.Printf("Email %s sent to %s", subject, to)
	return nil
}
