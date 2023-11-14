package integration

import (
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/api"
	"slices"
	"strings"
)

func GetControllerAliases(controller *api.ControllerData) []string {
	var out []string

	aliasRoles := []string{"ATM", "DATM", "TA", "EC", "FE", "WM"}
	for _, role := range controller.Roles {
		if slices.Contains(aliasRoles, role.Role) {
			out = append(out,
				fmt.Sprintf("%s-%s@vatusa.net", strings.ToLower(role.Facility), strings.ToLower(role.Role)))
		}
	}

	return out
}
