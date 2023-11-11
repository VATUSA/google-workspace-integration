package integration

import (
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/config"
	"github.com/VATUSA/google-workspace-integration/internal/database"
	"github.com/VATUSA/google-workspace-integration/internal/workspace_helper"
	admin "google.golang.org/api/admin/directory/v1"
)

func SyncAccounts() error {
	svc, err := workspace_helper.GetService()
	if err != nil {
		return err
	}
	users, err := svc.Users.List().Customer(config.GOOGLE_CUSTOMER_ID).Do()
	var existingUserEmails []string
	usersByEmail := make(map[string]*admin.User)

	for _, user := range users.Users {
		existingUserEmails = append(existingUserEmails, user.PrimaryEmail)
		usersByEmail[user.PrimaryEmail] = user
		fmt.Printf("%s\n", user.PrimaryEmail)
	}

	accounts, err := database.FetchAccounts()
	if err != nil {
		return nil
	}

	for _, account := range accounts {
		user := usersByEmail[account.PrimaryEmail]
		if user == nil {
			if !account.IsActive {
				continue
			}
			// TODO: Create Account
		} else {

		}
	}

	// TODO

	return nil
}
