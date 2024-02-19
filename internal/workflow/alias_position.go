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

func PositionAliasesMain() error {
	log.Printf("Start PositionAliasesMain")
	facilities, err := api.GetFacilities()
	var facilitiesById = make(map[string]api.FacilityData)
	for _, facility := range facilities {
		facilitiesById[facility.Id] = facility
	}
	accounts, err := database.FetchAccounts()
	if err != nil {
		return err
	}
	var accountsByCID = map[uint64]database.Account{}
	for _, account := range accounts {
		accountsByCID[account.CID] = account
	}

	// Remove needs to be first, otherwise add will fail on replacement
	removePositionAliases(accounts, facilitiesById)
	addPositionAliases(facilities, accountsByCID)

	log.Printf("End PositionAliasesMain")
	return nil
}

func addPositionAliases(facilities []api.FacilityData, accountsByCID map[uint64]database.Account) {
	for _, facility := range facilities {
		for _, role := range config.FacilityAliasRoles {
			holderCID := facilityPositionHolder(facility, role)
			account, ok := accountsByCID[holderCID]
			if ok {
				aliasEmails := positionAliasEmails(facility.Id, role)
				var existingAliases []string
				for _, existingAlias := range account.Aliases {
					existingAliases = append(existingAliases, existingAlias.Email)
				}
				for _, aliasEmail := range aliasEmails {
					if !slices.Contains(existingAliases, aliasEmail) {
						if !slices.Contains(existingAliases, aliasEmail) {
							aliasRecord := database.Alias{
								Email:     aliasEmail,
								AccountId: account.Id,
								Account:   &account,
								AliasType: database.AliasType_FacilityPosition,
								Facility:  facility.Id,
								Role:      role,
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

func removePositionAliases(accounts []database.Account, facilitiesById map[string]api.FacilityData) {
	for _, account := range accounts {
		for _, alias := range account.Aliases {
			if alias.AliasType == database.AliasType_FacilityPosition {
				facility, ok := facilitiesById[alias.Facility]
				if ok {
					if account.CID != facilityPositionHolder(facility, alias.Role) {
						log.Printf("Deleting alias %s for user %s", alias.Email, account.PrimaryEmail)
						err := google.DeleteUserAlias(account.PrimaryEmail, alias.Email)
						if err != nil {
							log.Printf("Error when removing alias %s from user %s - %v",
								alias.Email, account.PrimaryEmail, err)
							continue
						}
						err = alias.Delete()
						if err != nil {
							log.Printf("Error when deleting alias %s record from user %s - %v",
								alias.Email, account.PrimaryEmail, err)
						}
					}
				}
			}
		}
	}
}

func facilityPositionHolder(facility api.FacilityData, role string) uint64 {
	if role == config.AirTrafficManager {
		return facility.AirTrafficManagerCID
	} else if role == config.DeputyAirTrafficManager {
		return facility.DeputyAirTrafficManagerCID
	} else if role == config.TrainingAdministrator {
		return facility.TrainingAdministratorCID
	} else if role == config.EventCoordinator {
		return facility.EventCoordinatorCID
	} else if role == config.FacilityEngineer {
		return facility.FacilityEngineerCID
	} else if role == config.WebMaster {
		return facility.WebMasterCID
	}
	return 0
}

func positionAliasEmails(facility string, position string) []string {
	out := []string{fmt.Sprintf("%s-%s@vatusa.net", facility, position)}
	facilityDomains := config.FacilityDomains[facility]
	for _, domain := range facilityDomains {
		out = append(out, fmt.Sprintf("%s@%s", position, domain))
	}
	return out
}
