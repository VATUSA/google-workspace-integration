package google

import (
	"github.com/VATUSA/google-workspace-integration/internal/config"
	"github.com/VATUSA/google-workspace-integration/internal/workspace_helper"
	admin "google.golang.org/api/admin/directory/v1"
	"log"
)

func AddUserToGroup(userEmail string, groupEmail string, isManager bool) (err error) {
	svc, err := workspace_helper.GetService()
	if err != nil {
		return
	}
	role := "MEMBER"
	if isManager {
		role = "MANAGER"
	}
	existingMember, _ := GetGroupMembership(userEmail, groupEmail)
	if existingMember != nil {
		// Don't try to create the group membership if it already exists
		// This ensures failure recovery and backwards compatibility with the existing group memberships
		log.Printf("Prevented membership creation attempt for user: %s - group: %s. "+
			"This should only happen if the database is purged or if a membership is manually created. ",
			userEmail, groupEmail)
		return
	}
	member := admin.Member{
		DeliverySettings: "ALL_MAIL",
		Email:            userEmail,
		Role:             role,
		Type:             "USER",
	}
	if config.DEBUG {
		return
	}
	_, err = svc.Members.Insert(groupEmail, &member).Do()
	return
}

func RemoveUserFromGroup(userEmail string, groupEmail string) (err error) {
	svc, err := workspace_helper.GetService()
	if err != nil {
		return
	}
	_, err = svc.Members.Get(groupEmail, userEmail).Do()
	if err != nil {
		return
	}
	if config.DEBUG {
		return
	}
	err = svc.Members.Delete(groupEmail, userEmail).Do()
	return
}

func ChangeUserGroupRole(userEmail string, groupEmail string, isManager bool) (err error) {
	svc, err := workspace_helper.GetService()
	if err != nil {
		return
	}
	role := "MEMBER"
	if isManager {
		role = "MANAGER"
	}
	member, err := svc.Members.Get(groupEmail, userEmail).Do()
	if err != nil {
		return
	}
	member.Role = role
	if config.DEBUG {
		return
	}
	_, err = svc.Members.Update(groupEmail, userEmail, member).Do()
	return
}

func GetGroupMembership(userEmail string, groupEmail string) (member *admin.Member, err error) {
	svc, err := workspace_helper.GetService()
	if err != nil {
		return
	}
	member, err = svc.Members.Get(groupEmail, userEmail).Do()
	return
}

func GetGroup(groupEmail string) (group *admin.Group, err error) {
	svc, err := workspace_helper.GetService()
	if err != nil {
		return
	}
	group, err = svc.Groups.Get(groupEmail).Do()
	return
}

func CreateGroup(groupEmail string, displayName string) (err error) {
	svc, err := workspace_helper.GetService()
	if err != nil {
		return
	}
	existingGroup, _ := GetGroup(groupEmail)
	if existingGroup != nil {
		// Don't try to create the group if it already exists
		// This ensures failure recovery and backwards compatibility with the existing group memberships
		log.Printf("Prevented group creation attempt for group: %s. "+
			"This should only happen if the database is purged or if a group is manually created. ",
			groupEmail)
		return
	}
	group := &admin.Group{
		Description: "",
		Email:       groupEmail,
		Name:        displayName,
	}
	if config.DEBUG {
		return
	}
	_, err = svc.Groups.Insert(group).Do()
	return
}

func AddGroupAlias(groupEmail string, aliasEmail string) (err error) {
	svc, err := workspace_helper.GetService()
	if err != nil {
		return
	}
	newAlias := &admin.Alias{
		Alias:        aliasEmail,
		PrimaryEmail: groupEmail,
	}
	existingAlias := GroupAliasExists(groupEmail, aliasEmail)
	if !existingAlias {
		// Don't try to create the alias if it already exists
		// This ensures failure recovery and backwards compatibility with the existing group aliases
		log.Printf("Prevented alias creation attempt for group: %s - alias: %s. "+
			"This should only happen if the database is purged or if a user is manually created. ",
			groupEmail, aliasEmail)
		return
	}
	if config.DEBUG {
		return
	}
	_, err = svc.Groups.Aliases.Insert(groupEmail, newAlias).Do()
	return
}

func GroupAliasExists(groupEmail string, alias string) bool {
	svc, err := workspace_helper.GetService()
	aliases, err := svc.Groups.Aliases.List(groupEmail).Do()
	if err != nil {
		return false
	}
	for _, a := range aliases.Aliases {
		if a.(map[string]interface{})["Alias"] == alias {
			return true
		}
	}
	return false
}
