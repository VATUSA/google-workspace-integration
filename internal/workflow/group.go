package workflow

import (
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/config"
	"github.com/VATUSA/google-workspace-integration/internal/database"
	"github.com/VATUSA/google-workspace-integration/internal/google"
	"log"
	"strings"
)

func GroupsMain() error {
	log.Printf("Start GroupsMain")
	groups, err := database.FetchGroups()
	if err != nil {
		return err
	}
	var groupsByEmail = make(map[string]database.Group)
	for _, group := range groups {
		groupsByEmail[group.PrimaryEmail] = group
	}

	createFacilityGroups(groupsByEmail)
	createRoleGroups(groupsByEmail)

	log.Printf("End GroupsMain")
	return nil
}

func createFacilityGroups(groupsByEmail map[string]database.Group) {
	for _, facility := range config.AllGroupFacilities {
		for _, groupType := range config.FacilityGroupTypes {
			groupEmail := facilityGroupEmail(facility, groupType)
			_, exists := groupsByEmail[groupEmail]
			if !exists {
				group := database.Group{
					PrimaryEmail: groupEmail,
					DisplayName:  facilityGroupDisplayName(facility, groupType),
					Facility:     facility,
					GroupType:    groupType,
				}
				log.Printf("Creating group %s", groupEmail)
				err := google.CreateGroup(groupEmail, group.DisplayName)
				if err != nil {
					log.Printf("Error creating group %s - %v", groupEmail, err)
					continue
				}
				err = group.Save()
				if err != nil {
					log.Printf("Error saving group record %s - %v", groupEmail, err)
				}
			}
		}
	}
}

func createRoleGroups(groupsByEmail map[string]database.Group) {
	for groupType, groupEmail := range config.RoleGroups {
		_, exists := groupsByEmail[groupEmail]
		if !exists {
			group := database.Group{
				PrimaryEmail: groupEmail,
				DisplayName:  groupEmail,
				Facility:     "",
				GroupType:    groupType,
			}
			log.Printf("Creating group %s", groupEmail)
			err := google.CreateGroup(groupEmail, group.DisplayName)
			if err != nil {
				log.Printf("Error creating group %s - %v", groupEmail, err)
				continue
			}
			err = group.Save()
			if err != nil {
				log.Printf("Error saving group record %s - %v", groupEmail, err)
			}
		}
	}
}

func facilityGroupEmail(facility string, groupType string) string {
	if groupType == "" {
		return fmt.Sprintf("%s@vatusa.net", strings.ToLower(facility))
	}
	return fmt.Sprintf("%s-%s@vatusa.net", strings.ToLower(facility), groupType)
}

func facilityGroupDisplayName(facility string, groupType string) string {
	typeName := config.FacilityGroupTypeNamesMap[groupType]
	return fmt.Sprintf("%s %s", strings.ToUpper(facility), typeName)
}
