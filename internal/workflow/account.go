package workflow

import (
	"errors"
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/api"
	"github.com/VATUSA/google-workspace-integration/internal/config"
	"github.com/VATUSA/google-workspace-integration/internal/database"
	"github.com/VATUSA/google-workspace-integration/internal/email"
	"github.com/VATUSA/google-workspace-integration/internal/google"
	"log"
	"regexp"
	"slices"
	"strings"
	"time"
)

func AccountsMain() error {
	log.Print("Start AccountsMain")
	staffMembers, err := api.GetStaffMembers()
	if err != nil {
		return err
	}
	accounts, err := database.FetchAccounts()
	if err != nil {
		return err
	}
	var staffMembersByCID = map[uint64]api.ControllerData{}
	for _, controller := range staffMembers {
		staffMembersByCID[controller.CID] = controller
	}
	var accountsByCID = map[uint64]database.Account{}
	var existingAliases []string
	for _, account := range accounts {
		accountsByCID[account.CID] = account
		existingAliases = append(existingAliases, account.PrimaryEmail)
	}

	createAccounts(staffMembers, accountsByCID, existingAliases)
	updateAccounts(staffMembers, accountsByCID)
	suspendAccounts(accounts, staffMembersByCID)
	deleteAccounts(accounts)
	log.Print("End AccountsMain")
	return nil
}

func createAccounts(
	controllers []api.ControllerData, accountsByCID map[uint64]database.Account, existingAliases []string) {
	for _, controller := range controllers {
		_, ok := accountsByCID[controller.CID]
		if !ok && shouldControllerHaveAccount(controller) {
			alias, err := determineAlias(controller, existingAliases)
			if err != nil {
				log.Printf("Unable to determine alias for CID: %d - FN: %s - LN: %s - %v",
					controller.CID, controller.FirstName, controller.LastName, err)
				continue
			}
			account := database.Account{
				CID:          controller.CID,
				FirstName:    controller.FirstName,
				LastName:     controller.LastName,
				PrimaryAlias: alias,
				PrimaryEmail: fmt.Sprintf("%s@vatusa.net", alias),
				IsManaged:    true,
				IsSuspended:  false,
				SuspendedAt:  nil,
			}
			log.Printf("Creating user %s for CID %d", account.PrimaryEmail, account.CID)
			password, err := google.CreateUser(account.FirstName, account.LastName, account.PrimaryEmail)
			if err != nil {
				log.Printf("Error creating google account for CID: %d - %v", account.CID, err)
				continue
			}
			err = account.Save()
			if err != nil {
				log.Printf("Error saving account object for CID: %d - %v", account.CID, err)
			}
			if password != "" {
				err = email.SendNewAccountEmail(
					*controller.Email, account.FirstName, account.LastName, account.PrimaryEmail, password)
				if err != nil {
					log.Printf("Unable to send new account email for CID: %d - Email: %s - %v",
						account.CID, account.PrimaryEmail, err)
				}
			}
		}
	}
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func clearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func determineAlias(data api.ControllerData, existingAliases []string) (string, error) {
	aliasOverride, hasAliasOverride := config.PrimaryEmailOverrides[data.CID]
	if hasAliasOverride {
		return aliasOverride, nil
	}
	prospect := strings.ToLower(fmt.Sprintf("%s.%s", clearString(data.FirstName), clearString(data.LastName)))
	if slices.Contains(existingAliases, prospect) {
		return "", errors.New("alias already exists")
	}
	return prospect, nil
}

func updateAccounts(
	controllers []api.ControllerData, accountsByCID map[uint64]database.Account) {
	for _, controller := range controllers {
		account, ok := accountsByCID[controller.CID]
		if ok {
			shouldChange := false
			if controller.FirstName != account.FirstName {
				account.FirstName = controller.FirstName
				shouldChange = true
			}
			if controller.LastName != account.LastName {
				account.LastName = controller.LastName
				shouldChange = true
			}
			if shouldChange {
				// TODO: Actually update the user
				// Don't save the account object until we actually update the user
				// account.Save()
			}
		}
	}
}

func suspendAccounts(accounts []database.Account, controllersByCID map[uint64]api.ControllerData) {
	for _, account := range accounts {
		controller, ok := controllersByCID[account.CID]
		if ok {
			shouldHaveAccount := shouldControllerHaveAccount(controller)
			if !account.IsSuspended && !shouldHaveAccount {
				now := time.Now()
				account.IsSuspended = true
				account.SuspendedAt = &now
				err := google.SetUserSuspended(account.PrimaryEmail, true)
				if err != nil {
					log.Printf("Error setting user suspended - CID: %d - %v", account.CID, err)
					continue
				}
				err = account.Save()
				if err != nil {
					log.Printf("Error saving account object for CID: %d - %v", account.CID, err)
				}
			}
			if shouldHaveAccount && account.IsSuspended {
				account.IsSuspended = false
				account.SuspendedAt = nil
				err := google.SetUserSuspended(account.PrimaryEmail, false)
				if err != nil {
					log.Printf("Error setting user un-suspended - CID: %d - %v", account.CID, err)
					continue
				}
				err = account.Save()
				if err != nil {
					log.Printf("Error saving account object for CID: %d - %v", account.CID, err)
				}
			}
		} else {
			controller, err := api.GetControllerData(uint(account.CID))
			if err != nil {
				log.Printf("Error getting controller data for CID: %d - %v", account.CID, err)
				return
			}
			shouldHaveAccount := shouldControllerHaveAccount(*controller)
			if !account.IsSuspended && !shouldHaveAccount {
				now := time.Now()
				account.IsSuspended = true
				account.SuspendedAt = &now
				err := google.SetUserSuspended(account.PrimaryEmail, true)
				if err != nil {
					log.Printf("Error setting user suspended - CID: %d - %v", account.CID, err)
					continue
				}
				err = account.Save()
				if err != nil {
					log.Printf("Error saving account object for CID: %d - %v", account.CID, err)
				}
			}
		}
	}
}

const AccountDeleteDelay = 7 * 24 * time.Hour

func deleteAccounts(accounts []database.Account) {
	for _, account := range accounts {
		if account.IsSuspended {
			if account.SuspendedAt != nil && account.SuspendedAt.Add(AccountDeleteDelay).Before(time.Now()) {
				err := google.DeleteUser(account.PrimaryEmail)
				if err != nil {
					log.Printf("Error deleting user - CID: %d - %v", account.CID, err)
					continue
				}
				err = account.Delete()
				if err != nil {
					log.Printf("Error deleting account object for CID: %d - %v", account.CID, err)
				}
			}
		}
	}
}

func shouldControllerHaveAccount(data api.ControllerData) bool {
	if !data.FlagHomeController {
		return false
	}
	return data.HasAnyRole(config.AccountEntitlementRoles)
}
