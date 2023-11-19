package integration

import (
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/api"
	"github.com/VATUSA/google-workspace-integration/internal/config"
	"github.com/VATUSA/google-workspace-integration/internal/database"
	"github.com/VATUSA/google-workspace-integration/internal/workspace_helper"
	admin "google.golang.org/api/admin/directory/v1"
	"slices"
	"strings"
)

type GroupInfo struct {
	Name    string
	Email   string
	Aliases []string
}

func makeGroupInfo(facility string, nameFmt string, addressPart string, domains []database.Domain) GroupInfo {
	var aliases []string
	for _, domain := range domains {
		aliases = append(aliases, fmt.Sprintf("%s@%s", addressPart, strings.ToLower(domain.Domain)))
	}
	return GroupInfo{
		Name:    fmt.Sprintf(nameFmt, facility),
		Email:   groupEmail(facility, addressPart),
		Aliases: aliases,
	}
}

func GroupInfoForFacility(facility string) ([]GroupInfo, error) {
	var out []GroupInfo

	domains, err := database.FetchDomainsByFacility(facility)
	if err != nil {
		return nil, err
	}

	out = append(out, makeGroupInfo(facility, "%s ARTCC", "", domains))
	out = append(out, makeGroupInfo(facility, "%s Senior Staff", "sstf", domains))
	out = append(out, makeGroupInfo(facility, "%s Staff", "staff", domains))
	out = append(out, makeGroupInfo(facility, "%s Instructors", "instructors", domains))
	out = append(out, makeGroupInfo(facility, "%s Training Staff", "training", domains))
	out = append(out, makeGroupInfo(facility, "%s Events", "events", domains))
	out = append(out, makeGroupInfo(facility, "%s Facilities", "facilities", domains))
	out = append(out, makeGroupInfo(facility, "%s Web", "web", domains))

	return out, nil
}

type MembershipType = string

var (
	GroupRoleMember  MembershipType = "MEMBER"
	GroupRoleManager MembershipType = "MANAGER"
	GroupRoleOwner   MembershipType = "OWNER"
	facilities                      = []string{
		"ZAB", "ZAN", "ZTL", "ZBW", "ZAU", "ZOB", "ZDV", "ZFW", "HCF", "ZHU", "ZID",
		"ZJX", "ZKC", "ZLA", "ZME", "ZMA", "ZMP", "ZNY", "ZOA", "ZLC", "ZSE", "ZDC"}
)

func groupEmail(facility string, groupType string) string {
	if groupType == "" {
		return fmt.Sprintf("%s@vatusa.net", strings.ToLower(facility))
	}
	return fmt.Sprintf("%s-%s@vatusa.net", strings.ToLower(facility), groupType)
}

func setGroupMembership(out map[string]MembershipType, email string, membership MembershipType) {
	if membership == GroupRoleMember {
		if out[email] != GroupRoleManager {
			out[email] = membership
		}
	} else {
		out[email] = membership
	}
}

func GetControllerGroups(controller *api.ControllerData) map[string]MembershipType {
	out := make(map[string]MembershipType)

	if controller.Rating == 8 || controller.Rating == 10 {
		if slices.Contains(facilities, controller.Facility) {
			setGroupMembership(out, groupEmail(controller.Facility, "instructors"), GroupRoleMember)
			setGroupMembership(out, groupEmail(controller.Facility, "training"), GroupRoleMember)
		}
		setGroupMembership(out, groupEmail("instructor", "all"), GroupRoleMember)
	}

	for _, role := range controller.Roles {
		if role.Role == "ATM" {
			setGroupMembership(out, groupEmail(controller.Facility, "sstf"), GroupRoleMember)
			setGroupMembership(out, groupEmail(controller.Facility, "staff"), GroupRoleManager)
			setGroupMembership(out, groupEmail(controller.Facility, ""), GroupRoleManager)
			setGroupMembership(out, groupEmail("atm", "all"), GroupRoleMember)
		} else if role.Role == "DATM" {
			setGroupMembership(out, groupEmail(controller.Facility, "sstf"), GroupRoleMember)
			setGroupMembership(out, groupEmail(controller.Facility, "staff"), GroupRoleManager)
			setGroupMembership(out, groupEmail(controller.Facility, ""), GroupRoleManager)
			setGroupMembership(out, groupEmail("datm", "all"), GroupRoleMember)
		} else if role.Role == "TA" {
			setGroupMembership(out, groupEmail(controller.Facility, "sstf"), GroupRoleMember)
			setGroupMembership(out, groupEmail(controller.Facility, "staff"), GroupRoleManager)
			setGroupMembership(out, groupEmail(controller.Facility, ""), GroupRoleManager)
			setGroupMembership(out, groupEmail("ta", "all"), GroupRoleMember)
		} else if role.Role == "EC" {
			setGroupMembership(out, groupEmail(controller.Facility, "staff"), GroupRoleMember)
			setGroupMembership(out, groupEmail("ec", "all"), GroupRoleMember)
		} else if role.Role == "FE" {
			setGroupMembership(out, groupEmail(controller.Facility, "staff"), GroupRoleMember)
			setGroupMembership(out, groupEmail("fe", "all"), GroupRoleMember)
		} else if role.Role == "WM" {
			setGroupMembership(out, groupEmail(controller.Facility, "staff"), GroupRoleMember)
			setGroupMembership(out, groupEmail("wm", "all"), GroupRoleMember)
		} else if role.Role == "INS" {
			setGroupMembership(out, groupEmail(controller.Facility, "instructors"), GroupRoleMember)
			setGroupMembership(out, groupEmail(controller.Facility, "training"), GroupRoleMember)
			setGroupMembership(out, groupEmail("instructor", "all"), GroupRoleMember)
		} else if role.Role == "MTR" {
			setGroupMembership(out, groupEmail(controller.Facility, "training"), GroupRoleMember)
			setGroupMembership(out, groupEmail("mentor", "all"), GroupRoleMember)
		} else if role.Role == "DICE" {
			setGroupMembership(out, groupEmail("dice", ""), GroupRoleMember)
		}
	}

	return out
}

func AllManagedGroupEmails() ([]string, error) {
	var out []string

	allGroups, err := AllManagedGroups()
	if err != nil {
		return nil, err
	}

	for _, gi := range allGroups {
		out = append(out, gi.Email)
	}

	return out, nil
}

func AllManagedGroups() ([]GroupInfo, error) {
	var out []GroupInfo

	out = append(out, GroupInfo{
		Name:  "DICE-Team",
		Email: "dice@vatusa.net",
	})

	roles := []string{"atm", "datm", "ta", "ec", "fe", "wm", "instructor", "mentor"}
	for _, role := range roles {
		out = append(out, GroupInfo{
			Name:  fmt.Sprintf("AlL-%s", strings.ToUpper(role)),
			Email: fmt.Sprintf("%s-all@vatusa.net", role),
		})
	}

	for _, facility := range facilities {
		facilityGroups, err := GroupInfoForFacility(facility)
		if err != nil {
			return nil, err
		}
		for _, gi := range facilityGroups {
			out = append(out, gi)
		}
	}

	return out, nil
}

func CreateFacilityGroups() error {
	svc, err := workspace_helper.GetService()
	if err != nil {
		return err
	}
	groups, err := svc.Groups.List().Customer(config.GOOGLE_CUSTOMER_ID).Do()
	if err != nil {
		return err
	}
	var existingGroups []string
	for _, group := range groups.Groups {
		existingGroups = append(existingGroups, group.Email)
	}

	for _, facility := range facilities {
		facilityGroups, err := GroupInfoForFacility(facility)
		if err != nil {
			return err
		}
		for _, gi := range facilityGroups {
			err := createGroup(svc, existingGroups, gi)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createGroup(svc *admin.Service, existingGroups []string, info GroupInfo) error {
	if !slices.Contains(existingGroups, info.Email) {
		// Don't try to create a group that already exists

		group := admin.Group{
			AdminCreated: true,
			Description:  "",
			Email:        info.Email,
			Name:         info.Name,
		}
		println(fmt.Sprintf("Creating group %s", info.Email))
		_, err := svc.Groups.Insert(&group).Do()
		if err != nil {
			return err
		}
	}
	group, err := svc.Groups.Get(info.Email).Do()
	if err != nil {
		return err
	}
	for _, alias := range info.Aliases {
		if !slices.Contains(group.Aliases, alias) {
			println(fmt.Sprintf("Adding alias %s to group %s", alias, info.Email))
			alias := admin.Alias{
				Alias:        alias,
				PrimaryEmail: info.Email,
			}
			_, err := svc.Groups.Aliases.Insert(info.Email, &alias).Do()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
