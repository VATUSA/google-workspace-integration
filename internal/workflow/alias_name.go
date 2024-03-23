package workflow

import (
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/api"
	"github.com/VATUSA/google-workspace-integration/internal/config"
	"github.com/VATUSA/google-workspace-integration/internal/database"
	"github.com/VATUSA/google-workspace-integration/internal/google"
	"log"
	"slices"
)

func NameAliasesMain() error {
	log.Printf("Start NameAliasesMain")
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

	removeNameAliases(accounts, staffMembersByCID)
	addNameAliases(accounts, staffMembersByCID)

	log.Printf("End NameAliasesMain")
	return nil
}

func addNameAliases(accounts []database.Account, controllersByCID map[uint64]api.ControllerData) {
	for _, account := range accounts {
		controller, ok := controllersByCID[account.CID]
		if ok {
			var existingAliases []string
			for _, alias := range account.Aliases {
				existingAliases = append(existingAliases, alias.Email)
			}

			for facility, facilityDomains := range config.FacilityDomains {
				if shouldHaveFacilityNameAlias(controller, facility) {
					for _, domain := range facilityDomains {
						aliasEmail := domainAlias(account.PrimaryAlias, domain)
						if !slices.Contains(existingAliases, aliasEmail) {
							aliasRecord := database.Alias{
								Email:     aliasEmail,
								AccountId: account.Id,
								Account:   &account,
								AliasType: database.AliasType_FacilityName,
								Facility:  facility,
								Role:      "",
							}
							log.Printf("Creating alias %s for user %s", aliasEmail, account.PrimaryEmail)
							err := google.AddUserAlias(account.PrimaryEmail, aliasEmail)
							if err != nil {
								log.Printf("Error creating alias %s for user %s - %v",
									aliasEmail, account.PrimaryEmail, err)
								continue
							}
							err = aliasRecord.Save()
							if err != nil {
								log.Printf("Error saving alias %s record for user %s - %v",
									aliasEmail, account.PrimaryEmail, err)
							}
						}
					}
				}
			}

		}
	}
}

func removeNameAliases(accounts []database.Account, controllersByCID map[uint64]api.ControllerData) {
	for _, account := range accounts {
		controller, ok := controllersByCID[account.CID]
		if ok {
			for _, alias := range account.Aliases {
				if alias.AliasType == database.AliasType_FacilityName &&
					!shouldHaveFacilityNameAlias(controller, alias.Facility) {
					err := google.DeleteUserAlias(account.PrimaryEmail, alias.Email)
					if err != nil {
						log.Printf("Error removing alias %s for user %s - %v",
							alias.Email, account.PrimaryEmail, err)
						continue
					}
					err = alias.Delete()
					if err != nil {
						log.Printf("Error deleting alias %s record for user %s - %v",
							alias.Email, account.PrimaryEmail, err)
					}
				}
			}
		}
	}
}

func shouldHaveFacilityNameAlias(controller api.ControllerData, facility string) bool {
	return controller.HasAnyRoleAtFacility(config.FacilityNameAliasEntitlementRoles, facility)
}

func domainAlias(primaryAlias string, domain string) string {
	return fmt.Sprintf("%s@%s", primaryAlias, domain)
}
