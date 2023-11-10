package integration

import (
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/config"
	"github.com/VATUSA/google-workspace-integration/internal/workspace_helper"
)

func SyncAccounts() error {
	svc, err := workspace_helper.GetService()
	if err != nil {
		return err
	}
	users, err := svc.Users.List().Customer(config.GOOGLE_CUSTOMER_ID).Do()
	var existingUserEmails []string

	for _, user := range users.Users {
		existingUserEmails = append(existingUserEmails, user.PrimaryEmail)
		fmt.Printf("%s\n", user.PrimaryEmail)
	}

	// TODO

	return nil
}
