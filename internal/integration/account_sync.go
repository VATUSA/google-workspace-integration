package integration

import (
	"crypto/md5"
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/api"
	"github.com/VATUSA/google-workspace-integration/internal/config"
	"github.com/VATUSA/google-workspace-integration/internal/database"
	"github.com/VATUSA/google-workspace-integration/internal/workspace_helper"
	admin "google.golang.org/api/admin/directory/v1"
	"math/rand"
	"slices"
	"strings"
)

func generatePassword() string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz" + "0123456789" + "!@#$%^&*")
	length := 16
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

var (
	ManagedGroupEmails = AllManagedGroupEmails()
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
	}

	accounts, err := database.FetchAccounts()
	if err != nil {
		return nil
	}

	for _, account := range accounts {
		user := usersByEmail[account.PrimaryEmail]
		data, err := api.GetControllerData(account.CID)
		if err != nil {
			return err
		}
		hasChange := false
		if user == nil {
			if !account.IsActive {
				continue
			}
			password := generatePassword()
			println(fmt.Sprintf("Creating user %s - Password: %s", account.PrimaryEmail, password))
			// TODO: Create Account
			newUser := admin.User{
				Aliases:                    nil,
				Archived:                   false,
				ChangePasswordAtNextLogin:  true,
				IncludeInGlobalAddressList: false,
				Name: &admin.UserName{
					DisplayName:     fmt.Sprintf("%s %s", data.FirstName, data.LastName),
					FamilyName:      data.LastName,
					FullName:        fmt.Sprintf("%s %s", data.FirstName, data.LastName),
					GivenName:       data.FirstName,
					ForceSendFields: nil,
					NullFields:      nil,
				},
				Notes:            nil,
				OrgUnitPath:      "/Managed",
				Password:         fmt.Sprintf("%x", md5.Sum([]byte(password))),
				HashFunction:     "MD5",
				PrimaryEmail:     account.PrimaryEmail,
				RecoveryEmail:    "",
				RecoveryPhone:    "",
				Suspended:        false,
				SuspensionReason: "",
			}
			createdUser, err := svc.Users.Insert(&newUser).Do()
			if err != nil {
				return err
			}
			user = createdUser
			account.IsCreated = true
			account.TemporaryPassword = password
		} else {
			account.IsCreated = true
			if !account.IsActive && !user.Suspended {
				user.Suspended = true
				account.IsDeleted = true
				hasChange = true
			}
		}
		if user.OrgUnitPath != "/Managed" {
			err = account.Save()
			if err != nil {
				return err
			}
			continue
		}
		if account.ShouldResetPassword {
			password := generatePassword()
			user.Password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
			user.HashFunction = "MD5"
			account.ShouldResetPassword = false
			account.TemporaryPassword = password
			fmt.Printf("Reset password for user %s - Password: %s\n", user.PrimaryEmail, password)
		}
		if hasChange && user != nil && user.OrgUnitPath == "/Managed" {
			_, err := svc.Users.Update(user.PrimaryEmail, user).Do()
			if err != nil {
				return err
			}
		}
		existingGroups, err := svc.Groups.List().Customer(config.GOOGLE_CUSTOMER_ID).Query("memberKey=" + user.Id).Do()
		if err != nil {
			return err
		}
		var existingGroupEmails []string
		groups := GetControllerGroups(data)
		for _, eg := range existingGroups.Groups {
			val, ok := groups[eg.Email]
			if !ok {
				if slices.Contains(ManagedGroupEmails, eg.Email) {
					fmt.Printf("Removing user %s from group %s\n", user.PrimaryEmail, eg.Email)
					err := svc.Members.Delete(eg.Email, user.PrimaryEmail).Do()
					if err != nil {
						return err
					}
				}
				continue
			}
			membership, err := svc.Members.Get(eg.Email, user.PrimaryEmail).Do()
			if err != nil {
				return err
			}
			if membership.Role != val {
				fmt.Printf("Changing user %s in group %s from role %s to %s", user.PrimaryEmail, eg.Email, membership.Role, val)
				membership.Role = val
				_, err := svc.Members.Update(eg.Email, user.PrimaryEmail, membership).Do()
				if err != nil {
					return err
				}
			}
			existingGroupEmails = append(existingGroupEmails, eg.Email)
		}

		for group, role := range groups {
			if !slices.Contains(existingGroupEmails, group) {
				member := admin.Member{
					DeliverySettings: "ALL_MAIL",
					Email:            user.PrimaryEmail,
					Role:             role,
					Type:             "USER",
				}
				fmt.Printf("Adding user %s to group %s as %s\n", user.PrimaryEmail, group, role)
				_, err := svc.Members.Insert(group, &member).Do()
				if err != nil {
					return err
				}
			}
		}
		err = account.Save()
		if err != nil {
			return err
		}
	}

	// TODO

	return nil
}
