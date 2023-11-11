package integration

import (
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/api"
	"github.com/VATUSA/google-workspace-integration/internal/config"
	"github.com/VATUSA/google-workspace-integration/internal/workspace_helper"
	admin "google.golang.org/api/admin/directory/v1"
	"slices"
	"strings"
)

type GroupInfo struct {
	Name  string
	Email string
}

func makeGroupInfo(facility string, nameFmt string, emailFmt string) GroupInfo {
	return GroupInfo{
		Name:  fmt.Sprintf(nameFmt, facility),
		Email: fmt.Sprintf(emailFmt, strings.ToLower(facility)),
	}
}

func GroupInfoForFacility(facility string) []GroupInfo {
	var out []GroupInfo

	out = append(out, makeGroupInfo(facility, "%s ARTCC", "%s@vatusa.net"))
	out = append(out, makeGroupInfo(facility, "%s Senior Staff", "%s-sstf@vatusa.net"))
	out = append(out, makeGroupInfo(facility, "%s Staff", "%s-staff@vatusa.net"))
	out = append(out, makeGroupInfo(facility, "%s Instructors", "%s-instructors@vatusa.net"))
	out = append(out, makeGroupInfo(facility, "%s Training Staff", "%s-training@vatusa.net"))

	return out
}

type MembershipType = string

var (
	GROUP_ROLE_MEMBER  MembershipType = "MEMBER"
	GROUP_ROLE_MANAGER MembershipType = "MANAGER"
	GROUP_ROLE_OWNER   MembershipType = "OWNER"
	facilities                        = []string{
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
	if membership == GROUP_ROLE_MEMBER {
		if out[email] != GROUP_ROLE_MANAGER {
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
			setGroupMembership(out, groupEmail(controller.Facility, "instructors"), GROUP_ROLE_MEMBER)
			setGroupMembership(out, groupEmail(controller.Facility, "training"), GROUP_ROLE_MEMBER)
		}
		setGroupMembership(out, groupEmail("instructor", "all"), GROUP_ROLE_MEMBER)
	}

	for _, role := range controller.Roles {
		if role.Role == "ATM" {
			setGroupMembership(out, groupEmail(controller.Facility, "sstf"), GROUP_ROLE_MEMBER)
			setGroupMembership(out, groupEmail(controller.Facility, "staff"), GROUP_ROLE_MANAGER)
			setGroupMembership(out, groupEmail(controller.Facility, ""), GROUP_ROLE_MANAGER)
			setGroupMembership(out, groupEmail("atm", "all"), GROUP_ROLE_MEMBER)
		} else if role.Role == "DATM" {
			setGroupMembership(out, groupEmail(controller.Facility, "sstf"), GROUP_ROLE_MEMBER)
			setGroupMembership(out, groupEmail(controller.Facility, "staff"), GROUP_ROLE_MANAGER)
			setGroupMembership(out, groupEmail(controller.Facility, ""), GROUP_ROLE_MANAGER)
			setGroupMembership(out, groupEmail("datm", "all"), GROUP_ROLE_MEMBER)
		} else if role.Role == "TA" {
			setGroupMembership(out, groupEmail(controller.Facility, "sstf"), GROUP_ROLE_MEMBER)
			setGroupMembership(out, groupEmail(controller.Facility, "staff"), GROUP_ROLE_MANAGER)
			setGroupMembership(out, groupEmail(controller.Facility, ""), GROUP_ROLE_MANAGER)
			setGroupMembership(out, groupEmail("ta", "all"), GROUP_ROLE_MEMBER)
		} else if role.Role == "EC" {
			setGroupMembership(out, groupEmail(controller.Facility, "staff"), GROUP_ROLE_MEMBER)
			setGroupMembership(out, groupEmail("ec", "all"), GROUP_ROLE_MEMBER)
		} else if role.Role == "FE" {
			setGroupMembership(out, groupEmail(controller.Facility, "staff"), GROUP_ROLE_MEMBER)
			setGroupMembership(out, groupEmail("fe", "all"), GROUP_ROLE_MEMBER)
		} else if role.Role == "WM" {
			setGroupMembership(out, groupEmail(controller.Facility, "staff"), GROUP_ROLE_MEMBER)
			setGroupMembership(out, groupEmail("wm", "all"), GROUP_ROLE_MEMBER)
		} else if role.Role == "INS" {
			setGroupMembership(out, groupEmail(controller.Facility, "instructors"), GROUP_ROLE_MEMBER)
			setGroupMembership(out, groupEmail(controller.Facility, "training"), GROUP_ROLE_MEMBER)
			setGroupMembership(out, groupEmail("instructor", "all"), GROUP_ROLE_MEMBER)
		} else if role.Role == "MTR" {
			setGroupMembership(out, groupEmail(controller.Facility, "training"), GROUP_ROLE_MEMBER)
			setGroupMembership(out, groupEmail("mentor", "all"), GROUP_ROLE_MEMBER)
		} else if role.Role == "DICE" {
			setGroupMembership(out, groupEmail("dice", ""), GROUP_ROLE_MEMBER)
		}
	}

	return out
}

func AllManagedGroupEmails() []string {
	var out []string

	for _, gi := range AllManagedGroups() {
		out = append(out, gi.Email)
	}

	return out
}

func AllManagedGroups() []GroupInfo {
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
		for _, gi := range GroupInfoForFacility(facility) {
			out = append(out, gi)
		}
	}

	return out
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
		for _, gi := range GroupInfoForFacility(facility) {
			err := createGroup(svc, existingGroups, gi)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createGroup(svc *admin.Service, existingGroups []string, info GroupInfo) error {
	if slices.Contains(existingGroups, info.Email) {
		// Don't try to create a group that already exists
		return nil
	}
	group := admin.Group{
		AdminCreated: true,
		Aliases:      nil,
		Description:  "",
		Email:        info.Email,
		Kind:         "admin#directory#group",
		Name:         info.Name,
	}
	println(fmt.Sprintf("Creating group %s", info.Email))
	_, err := svc.Groups.Insert(&group).Do()
	if err != nil {
		return err
	}
	return nil
}
