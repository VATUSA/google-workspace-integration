package integration

import (
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/api"
	"github.com/VATUSA/google-workspace-integration/internal/database"
	"github.com/VATUSA/google-workspace-integration/internal/email"
)

func ProcessPasswordEmails() error {
	accounts, err := database.FetchAccounts()
	if err != nil {
		return err
	}
	for _, account := range accounts {
		if account.TemporaryPassword != "" {
			data, err := api.GetControllerData(account.CID)
			if err != nil {
				return err
			}
			message := fmt.Sprintf("Hello %s %s,\n\n"+
				"Your VATUSA google account password has been reset (or an account has been newly created for you).\n"+
				"Your VATUSA email address is: %s\n Your new password is: \n%s\n "+
				"This password must be changed on your next login."+
				"If your account has been newly created, make sure you set up two factor authentication within one week to avoid getting locked out of your account.",
				data.FirstName, data.LastName, account.PrimaryEmail, account.TemporaryPassword)

			err = email.SendEmail(*data.Email, "VATUSA Google Account - Password Reset", message)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
