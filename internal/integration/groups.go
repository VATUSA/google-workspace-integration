package integration

import (
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/api"
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

func GetGroupsForController(controller *api.ControllerData) []string {
	var out []string

	if controller.Rating == 8 || controller.Rating == 10 {
		out = append(out, fmt.Sprintf("%s-instructors@vatusa.net", controller.Facility))
		out = append(out, fmt.Sprintf("%s-training@vatusa.net", controller.Facility))
	}

	for _, role := range controller.Roles {
		if role.Role == "ATM" {
			out = append(out, fmt.Sprintf("%s@vatusa.net", role.Facility))
			out = append(out, fmt.Sprintf("%s-staff@vatusa.net", role.Facility))
			out = append(out, fmt.Sprintf("%s-sstf@vatusa.net", role.Facility))
			out = append(out, "atm-all@vatusa.net")

		} else if role.Role == "DATM" {
			out = append(out, fmt.Sprintf("%s@vatusa.net", role.Facility))
			out = append(out, fmt.Sprintf("%s-staff@vatusa.net", role.Facility))
			out = append(out, fmt.Sprintf("%s-sstf@vatusa.net", role.Facility))
			out = append(out, "datm-all@vatusa.net")
		} else if role.Role == "TA" {
			out = append(out, fmt.Sprintf("%s@vatusa.net", role.Facility))
			out = append(out, fmt.Sprintf("%s-staff@vatusa.net", role.Facility))
			out = append(out, fmt.Sprintf("%s-sstf@vatusa.net", role.Facility))
			out = append(out, fmt.Sprintf("%s-instructors@vatusa.net", role.Facility))
			out = append(out, fmt.Sprintf("%s-training@vatusa.net", role.Facility))
			out = append(out, "ta-all@vatusa.net")
		} else if role.Role == "EC" {
			out = append(out, fmt.Sprintf("%s-staff@vatusa.net", role.Facility))
			out = append(out, "ec-all@vatusa.net")
		} else if role.Role == "FE" {
			out = append(out, fmt.Sprintf("%s-staff@vatusa.net", role.Facility))
			out = append(out, "fe-all@vatusa.net")
		} else if role.Role == "WM" {
			out = append(out, fmt.Sprintf("%s-staff@vatusa.net", role.Facility))
			out = append(out, "wm-all@vatusa.net")
		} else if role.Role == "INS" {
			out = append(out, fmt.Sprintf("%s-instructors@vatusa.net", role.Facility))
			out = append(out, fmt.Sprintf("%s-training@vatusa.net", role.Facility))
			out = append(out, "instructor-all@vatusa.net")
		} else if role.Role == "MTR" {
			out = append(out, fmt.Sprintf("%s-training@vatusa.net", role.Facility))
			out = append(out, "mentor-all@vatusa.net")
		}
	}

	return out
}

func CreateFacilityGroups() error {

	facilities := []string{"ZAB", "ZAN", "ZTL", "ZBW", "ZAU", "ZOB", "ZDV", "ZFW", "HCF", "ZHU", "ZID", "ZJX", "ZKC",
		"ZLA", "ZME", "ZMA", "ZMP", "ZNY", "ZOA", "ZLC", "ZSE", "ZDC"}

	svc, err := workspace_helper.GetService()
	if err != nil {
		return err
	}
	groups, err := svc.Groups.List().Customer("C038380rr").Do()
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
