package google

import (
	"github.com/VATUSA/google-workspace-integration/internal/config"
	"github.com/VATUSA/google-workspace-integration/internal/workspace_helper"
	admin "google.golang.org/api/admin/directory/v1"
	"log"
)

func AddUserAlias(userEmail string, alias string) (err error) {
	svc, err := workspace_helper.GetService()
	if err != nil {
		return
	}
	newAlias := &admin.Alias{
		Alias:        alias,
		PrimaryEmail: userEmail,
	}
	existingAlias, _ := GetAlias(userEmail, alias)
	if existingAlias != nil {
		// Don't try to create the alias if it already exists
		// This ensures failure recovery and backwards compatibility with the existing group memberships
		log.Printf("Prevented alias creation attempt for user: %s - alias: %s. "+
			"This should only happen if the database is purged or if a user is manually created. ",
			userEmail, alias)
		return
	}
	if config.DEBUG {
		return
	}
	_, err = svc.Users.Aliases.Insert(userEmail, newAlias).Do()
	return
}

func DeleteUserAlias(userEmail string, alias string) (err error) {
	svc, err := workspace_helper.GetService()
	if err != nil {
		return
	}
	if config.DEBUG {
		return
	}
	err = svc.Users.Aliases.Delete(userEmail, alias).Do()
	return
}

func GetAlias(userEmail string, alias string) (existingAlias *admin.Alias, err error) {
	svc, err := workspace_helper.GetService()
	aliases, err := svc.Users.Aliases.List(userEmail).Do()
	if err != nil {
		return
	}
	for _, a := range aliases.Aliases {
		if a.(admin.Alias).Alias == alias {
			tmp := a.(admin.Alias)
			return &tmp, nil
		}
	}
	return
}
