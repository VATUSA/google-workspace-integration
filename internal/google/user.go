package google

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/config"
	"github.com/VATUSA/google-workspace-integration/internal/workspace_helper"
	admin "google.golang.org/api/admin/directory/v1"
	"log"
	"math/rand"
	"strings"
)

func CreateUser(firstName string, lastName string, primaryEmail string) (password string, err error) {
	svc, err := workspace_helper.GetService()
	if err != nil {
		return
	}
	existingUser, _ := svc.Users.Get(primaryEmail).Do()
	if existingUser != nil {
		// Don't try to create the user if it already exists
		// This ensures failure recovery and backwards compatibility with the existing users
		password = "" // Avoid unset variable ambiguity
		log.Printf("Prevented create attempt for existing user %s. "+
			"This should only happen if the database is purged or if a user is manually created. ", primaryEmail)
		return
	}
	password = generatePassword()
	displayName := fmt.Sprintf("%s %s", firstName, lastName)
	newUser := admin.User{
		Aliases:                    nil,
		Archived:                   false,
		ChangePasswordAtNextLogin:  true,
		IncludeInGlobalAddressList: false,
		Name: &admin.UserName{
			DisplayName:     displayName,
			FamilyName:      lastName,
			FullName:        displayName,
			GivenName:       firstName,
			ForceSendFields: nil,
			NullFields:      nil,
		},
		Notes:            nil,
		OrgUnitPath:      "/Managed",
		Password:         fmt.Sprintf("%x", md5.Sum([]byte(password))),
		HashFunction:     "MD5",
		PrimaryEmail:     primaryEmail,
		RecoveryEmail:    "",
		RecoveryPhone:    "",
		Suspended:        false,
		SuspensionReason: "",
	}
	if config.DEBUG {
		return
	}
	_, err = svc.Users.Insert(&newUser).Do()
	return
}

func SetUserSuspended(primaryEmail string, isSuspended bool) (err error) {
	svc, err := workspace_helper.GetService()
	if err != nil {
		return
	}
	user, err := svc.Users.Get(primaryEmail).Do()
	if err != nil {
		return
	}
	user.Suspended = isSuspended
	if config.DEBUG {
		return
	}
	_, err = svc.Users.Update(primaryEmail, user).Do()
	return
}

func DeleteUser(primaryEmail string) (err error) {
	// TODO: Implement DeleteUser
	return errors.New("not implemented")
}

// TODO: Implement something to call this
func ResetPassword(primaryEmail string) (newPassword string, err error) {
	svc, err := workspace_helper.GetService()
	if err != nil {
		return
	}
	user, err := svc.Users.Get(primaryEmail).Do()
	if err != nil {
		return
	}
	newPassword = generatePassword()

	user.HashFunction = "MD5"
	user.Password = newPassword
	user.ChangePasswordAtNextLogin = true

	if config.DEBUG {
		return
	}
	_, err = svc.Users.Update(primaryEmail, user).Do()
	return
}

func generatePassword() string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789" + "!@#$%^&*")
	length := 16
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
