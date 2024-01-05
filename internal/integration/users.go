package integration

import (
	"errors"
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/api"
	"github.com/VATUSA/google-workspace-integration/internal/database"
	"regexp"
	"slices"
	"strings"
)

func DoProcessStaff() error {
	staffMembers, err := api.GetStaffMembers()
	if err != nil {
		return err
	}
	return ProcessControllers(staffMembers)
}

func ProcessControllers(data []api.ControllerData) error {
	accounts, err := database.FetchAccounts()
	if err != nil {
		return err
	}

	accountsByCid := make(map[uint]database.Account)
	var existingAliases []string
	var checkedCids []uint

	for _, account := range accounts {
		accountsByCid[account.CID] = account
		existingAliases = append(existingAliases, strings.ToLower(account.Alias))
	}

	for _, c := range data {
		account, ok := accountsByCid[c.CID]
		if ok {
			err := ProcessController(c, &account, existingAliases)
			if err != nil {
				fmt.Printf("Error while processing controller: %d - %v\n", c.CID, err)
			}
		} else {
			err := ProcessController(c, nil, existingAliases)
			if err != nil {
				fmt.Printf("Error while processing controller: %d - %v\n", c.CID, err)
			}
		}
		checkedCids = append(checkedCids, c.CID)
	}

	for cid, account := range accountsByCid {
		if !slices.Contains(checkedCids, cid) {
			data, err := api.GetControllerData(account.CID)
			if err != nil {
				return err
			}
			err = ProcessController(*data, &account, existingAliases)
			if err != nil {
				fmt.Printf("Error while processing controller: %d - %v\n", data.CID, err)
			}
			checkedCids = append(checkedCids, cid)
		}
	}
	return nil
}

func ProcessController(data api.ControllerData, account *database.Account, existingAliases []string) error {
	if !shouldControllerHaveAccount(data) {
		if account != nil {
			account.IsActive = false
			err := account.Save()
			if err != nil {
				return err
			}
		}
		return nil
	}
	if account == nil {
		alias, err := determineAlias(data, existingAliases)
		if err != nil {
			return err
		}
		account := database.NewAccount(data.CID, alias)
		err = account.Save()
		if err != nil {
			return err
		}
	}
	return nil
}

func shouldControllerHaveAccount(data api.ControllerData) bool {
	if !data.FlagHomeController {
		return false
	}
	//if data.Rating == 8 || data.Rating == 10 {
	//	return true
	//}
	allowedRoles := []string{"ATM", "DATM", "TA", "EC", "FE", "WM" /*"INS", "MTR",*/, "DICE", "USWT"}
	for _, role := range data.Roles {
		if slices.Contains(allowedRoles, role.Role) {
			return true
		}
		if strings.HasPrefix(role.Role, "US") {
			return true
		}
	}
	return false
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func clearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func determineAlias(data api.ControllerData, existingAliases []string) (string, error) {
	prospect := strings.ToLower(fmt.Sprintf("%s.%s", clearString(data.FirstName), clearString(data.LastName)))
	if slices.Contains(existingAliases, prospect) {
		return "", errors.New("alias already exists")
	}
	return prospect, nil
}
