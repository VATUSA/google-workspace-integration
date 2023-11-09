package integration

import (
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/workspace_helper"
)

func SyncAccounts() error {
	svc, err := workspace_helper.GetService()
	if err != nil {
		return err
	}
	users, err := svc.Users.List().Customer("C038380rr").Do()
	var existingUserEmails []string

	for _, user := range users.Users {
		existingUserEmails = append(existingUserEmails, user.PrimaryEmail)
		fmt.Printf("%s\n", user.PrimaryEmail)
	}

	// TODO

	return nil
}
