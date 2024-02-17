package workflow

import (
	"github.com/VATUSA/google-workspace-integration/internal/api"
	"github.com/VATUSA/google-workspace-integration/internal/config"
	"github.com/VATUSA/google-workspace-integration/internal/database"
	"github.com/VATUSA/google-workspace-integration/internal/google"
	"log"
	"slices"
)

func GroupMembershipsMain() error {
	log.Print("Start GroupMembershipsMain")
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

	addAccountGroups(accounts, staffMembersByCID)
	updateAccountGroups(accounts, staffMembersByCID)
	removeAccountGroups(accounts, staffMembersByCID)

	log.Print("End GroupMembershipsMain")
	return nil
}

func addAccountGroups(accounts []database.Account, controllersByCID map[uint64]api.ControllerData) {
	for _, account := range accounts {
		controller, ok := controllersByCID[account.CID]
		if ok {
			var memberships []string
			for _, membership := range account.GroupMemberships {
				memberships = append(memberships, membership.Group.PrimaryEmail)
			}
			addAccountFacilityGroups(account, controller, memberships)
			addAccountRoleGroups(account, controller, memberships)
		}
	}
}

func addAccountFacilityGroups(
	account database.Account, controller api.ControllerData, memberships []string) {
	for _, facility := range config.AllGroupFacilities {
		for _, groupType := range config.FacilityGroupTypes {
			email := facilityGroupEmail(facility, groupType)
			if !slices.Contains(memberships, email) && shouldControllerHaveGroup(controller, facility, groupType) {
				addUserToGroup(account, email, isControllerGroupManager(controller, facility, groupType))
			}
		}
	}
}

func addAccountRoleGroups(account database.Account, controller api.ControllerData, memberships []string) {
	for role, groupEmail := range config.RoleGroups {
		if !slices.Contains(memberships, groupEmail) && controller.HasRole(role) {
			addUserToGroup(account, groupEmail, false)
		}
	}
}

func addUserToGroup(account database.Account, groupEmail string, isManager bool) {
	membership := database.GroupMembership{
		GroupEmail: groupEmail,
		AccountID:  account.Id,
		Account:    &account,
		IsManager:  isManager,
	}
	log.Printf("Adding user %s to group %s - manager: %t\n", account.PrimaryEmail, groupEmail, isManager)
	err := google.AddUserToGroup(
		account.PrimaryEmail, groupEmail, isManager)
	if err != nil {
		log.Printf("Error adding user to groupType %s - %s - %v",
			account.PrimaryEmail, groupEmail, err)
		return
	}
	err = membership.Save()
	if err != nil {
		log.Printf("Error when saving group memership record %s - %s - %v",
			account.PrimaryEmail, groupEmail, err)
	}
}

func updateAccountGroups(accounts []database.Account, controllersByCID map[uint64]api.ControllerData) {
	for _, account := range accounts {
		controller, ok := controllersByCID[account.CID]
		if ok {
			for _, membership := range account.GroupMemberships {
				if membership.Group.Facility == "" {
					continue
				}
				isManager := isControllerGroupManager(controller, membership.Group.Facility, membership.Group.GroupType)
				if membership.IsManager != isManager {
					log.Printf("Updating role for user %s in group %s - Manager (%t -> %t)",
						account.PrimaryEmail, membership.Group.PrimaryEmail,
						membership.IsManager, isManager)
					membership.IsManager = isManager
					err := google.ChangeUserGroupRole(account.PrimaryEmail, membership.Group.PrimaryEmail, isManager)
					if err != nil {
						log.Printf("Error when changing user %s role in group %s to %t - %v",
							account.PrimaryEmail, membership.Group.PrimaryEmail, isManager, err)
					}
					err = membership.Save()
					if err != nil {
						log.Printf("Error when saving group memership record %s - %s - %v",
							account.PrimaryEmail, membership.Group.PrimaryEmail, err)
					}
				}
			}
		}
	}
}

func removeAccountGroups(accounts []database.Account, controllersByCID map[uint64]api.ControllerData) {
	for _, account := range accounts {
		controller, ok := controllersByCID[account.CID]
		if ok {
			for _, membership := range account.GroupMemberships {
				if !shouldControllerHaveGroup(controller, membership.Group.Facility, membership.Group.GroupType) {
					log.Printf("Removing user %s from group %s",
						account.PrimaryEmail, membership.Group.PrimaryEmail)
					err := google.RemoveUserFromGroup(account.PrimaryEmail, membership.Group.PrimaryEmail)
					if err != nil {
						log.Printf("Error when removing user from group %s - %s - %v",
							account.PrimaryEmail, membership.Group.PrimaryEmail, err)
						continue
					}
					err = membership.Delete()
					if err != nil {
						log.Printf("Error when deleting group membership record %s - %s - %v",
							account.PrimaryEmail, membership.Group.PrimaryEmail, err)
					}
				}
			}
		} else {
			// If no controller data, remove the user from all groups
			for _, membership := range account.GroupMemberships {
				err := google.RemoveUserFromGroup(account.PrimaryEmail, membership.Group.PrimaryEmail)
				if err != nil {
					log.Printf("Error when removing user from group %s - %s - %v",
						account.PrimaryEmail, membership.Group.PrimaryEmail, err)
					continue
				}
				err = membership.Delete()
				if err != nil {
					log.Printf("Error when deleting group membership record %s - %s - %v",
						account.PrimaryEmail, membership.Group.PrimaryEmail, err)
				}
			}
		}
	}
}

func shouldControllerHaveGroup(controller api.ControllerData, facility string, groupType string) bool {
	if facility != "" {
		requiredRoles, ok := config.FacilityGroupRequiredRolesMap[groupType]
		if ok {
			return controller.HasAnyRoleAtFacility(requiredRoles, facility)
		}
	} else {
		_, ok := config.RoleGroups[groupType]
		if ok {
			return controller.HasRole(groupType)
		}
	}
	return false
}

func isControllerGroupManager(controller api.ControllerData, facility string, groupType string) bool {
	requiredRoles, ok := config.FacilityGroupManagerRolesMap[groupType]
	if ok {
		return controller.HasAnyRoleAtFacility(requiredRoles, facility)
	}
	return false
}
