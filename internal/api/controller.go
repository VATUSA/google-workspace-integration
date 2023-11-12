package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VATUSA/google-workspace-integration/internal/config"
	"io/ioutil"
)

type StaffListDataWrapper struct {
	Data []ControllerData `json:"data"`
}

type ControllerDataWrapper struct {
	Data ControllerData `json:"data"`
}

type ControllerData struct {
	CID                uint    `json:"cid"`
	FirstName          string  `json:"fname"`
	LastName           string  `json:"lname"`
	Email              *string `json:"email"`
	Facility           string  `json:"facility"`
	Rating             int     `json:"rating"`
	RatingShort        string  `json:"rating_short"`
	FlagHomeController bool    `json:"flag_homecontroller"`
	Roles              []ControllerRoleData
	VisitingFacilities []ControllerVisitingFacilityData `json:"visiting_facilities"`
}

type ControllerRoleData struct {
	Facility string `json:"facility"`
	Role     string `json:"role"`
}

type ControllerVisitingFacilityData struct {
	Facility string `json:"facility"`
}

func GetControllerData(cid uint) (*ControllerData, error) {
	response, err := Get(fmt.Sprintf("/user/%d?apikey=%s", cid, config.VATUSA_API2_KEY))
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 404 {
		return nil, nil
	}
	if response.StatusCode != 200 {
		return nil, errors.New("HTTP Error when fetching controller data")
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var wrapper ControllerDataWrapper
	err = json.Unmarshal(responseData, &wrapper)
	if err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}

func GetStaffMembers() ([]ControllerData, error) {
	response, err := Get("/integration/staff")
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 404 {
		return nil, nil
	}
	if response.StatusCode != 200 {
		return nil, errors.New("HTTP Error when fetching controller data")
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var wrapper StaffListDataWrapper
	err = json.Unmarshal(responseData, &wrapper)
	if err != nil {
		return nil, err
	}
	return wrapper.Data, nil
}
