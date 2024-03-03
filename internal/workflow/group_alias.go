package workflow

import (
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/config"
	"github.com/VATUSA/google-workspace-integration/internal/database"
	"github.com/VATUSA/google-workspace-integration/internal/google"
	"log"
	"slices"
)

func GroupAliasesMain() error {
	log.Printf("Start GroupAliasesMain")
	groups, err := database.FetchGroups()
	if err != nil {
		return err
	}
	var groupsByEmail = make(map[string]database.Group)
	for _, group := range groups {
		groupsByEmail[group.PrimaryEmail] = group
	}

	createGroupCustomDomainAliases(groupsByEmail)
	createGroupExtraAliases(groupsByEmail)

	log.Printf("End GroupAliasesMain")
	return nil
}

func createGroupCustomDomainAliases(groupsByEmail map[string]database.Group) {
	for _, group := range groupsByEmail {
		var existingAliases []string
		for _, alias := range group.Aliases {
			existingAliases = append(existingAliases, alias.Email)
		}
		for _, domain := range config.FacilityDomains[group.Facility] {
			if group.GroupType == "" {
				// Don't try to create an alias for the general group
				continue
			}
			aliasEmail := fmt.Sprintf("%s@%s", group.GroupType, domain)
			if !slices.Contains(existingAliases, aliasEmail) {
				groupAlias := database.GroupAlias{
					Email:             aliasEmail,
					GroupPrimaryEmail: group.PrimaryEmail,
					Group:             &group,
					Facility:          group.Facility,
					Domain:            domain,
				}
				log.Printf("Creating group alias %s for group %s", aliasEmail, group.PrimaryEmail)
				err := google.AddGroupAlias(group.PrimaryEmail, aliasEmail)
				if err != nil {
					log.Printf("Error creating group alias %s for group %s - %v", aliasEmail, group.PrimaryEmail, err)
					continue
				}
				err = groupAlias.Save()
				if err != nil {
					log.Printf("Error saving groupAlias record %s - %v", aliasEmail, err)
				}
			}
		}
	}
}

func createGroupExtraAliases(groupsByEmail map[string]database.Group) {
	for _, group := range groupsByEmail {
		var existingAliases []string
		for _, alias := range group.Aliases {
			existingAliases = append(existingAliases, alias.Email)
		}
		for _, alias := range config.FacilityGroupCustomDomainAliases[group.GroupType] {
			for _, domain := range config.FacilityDomains[group.Facility] {
				aliasEmail := fmt.Sprintf("%s@%s", alias, domain)
				if !slices.Contains(existingAliases, aliasEmail) {
					groupAlias := database.GroupAlias{
						Email:             aliasEmail,
						GroupPrimaryEmail: group.PrimaryEmail,
						Group:             &group,
						Facility:          group.Facility,
						Domain:            domain,
					}
					log.Printf("Creating group alias %s for group %s", aliasEmail, group.PrimaryEmail)
					err := google.AddGroupAlias(group.PrimaryEmail, aliasEmail)
					if err != nil {
						log.Printf("Error creating group %s alias %s - %v", group.PrimaryEmail, aliasEmail, err)
						continue
					}
					err = groupAlias.Save()
					if err != nil {
						log.Printf("Error saving groupAlias record %s - %v", aliasEmail, err)
					}
				}
			}
		}
	}
}
