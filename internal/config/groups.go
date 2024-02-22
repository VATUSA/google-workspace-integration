package config

// RoleGroups are assigned to everyone with the given role, regardless of facility
var RoleGroups = map[string]string{
	AirTrafficManager:       "atm-all@vatusa.net",
	DeputyAirTrafficManager: "datm-all@vatusa.net",
	TrainingAdministrator:   "ta-all@vatusa.net",
	EventCoordinator:        "ec-all@vatusa.net",
	FacilityEngineer:        "fe-all@vatusa.net",
	WebMaster:               "wm-all@vatusa.net",
	Instructor:              "instructor-all@vatusa.net",
	Mentor:                  "mentor-all@vatusa.net",
	DICETeam:                "dice@vatusa.net",
	USWebTeam:               "techteam@vatusa.net",
	SocialMediaTeam:         "social-media-team@vatusa.net",
}

// Facility Group Suffix
const (
	FacilityGroupGeneral       = ""
	FacilityGroupSeniorStaff   = "sstf"
	FacilityGroupStaff         = "staff"
	FacilityGroupEvents        = "events"
	FacilityGroupFacilities    = "facilities"
	FacilityGroupWeb           = "web"
	FacilityGroupInstructors   = "instructors"
	FacilityGroupTrainingStaff = "training"
)

var FacilityGroupCustomDomainAliases = map[string][]string{
	FacilityGroupSeniorStaff: {"management"},
}

var FacilityGroupTypes = []string{
	FacilityGroupGeneral,
	FacilityGroupSeniorStaff,
	FacilityGroupStaff,
	FacilityGroupEvents,
	FacilityGroupFacilities,
	FacilityGroupWeb,
	FacilityGroupInstructors,
	FacilityGroupTrainingStaff,
}

var FacilityGroupTypeNamesMap = map[string]string{
	FacilityGroupGeneral:       "ARTCC",
	FacilityGroupSeniorStaff:   "Senior Staff",
	FacilityGroupStaff:         "Staff",
	FacilityGroupEvents:        "Events",
	FacilityGroupFacilities:    "Facility Engineers",
	FacilityGroupWeb:           "Web",
	FacilityGroupInstructors:   "Instructors",
	FacilityGroupTrainingStaff: "Training Staff",
}

var FacilityGroupRequiredRolesMap = map[string][]string{
	FacilityGroupGeneral:       {"ATM", "DATM", "TA"},
	FacilityGroupSeniorStaff:   {"ATM", "DATM", "TA"},
	FacilityGroupStaff:         {"ATM", "DATM", "TA", "EC", "FE", "WM"},
	FacilityGroupEvents:        {"EC"},
	FacilityGroupFacilities:    {"FE"},
	FacilityGroupWeb:           {"WM"},
	FacilityGroupInstructors:   {"INS"},
	FacilityGroupTrainingStaff: {"TA", "INS", "MTR"},
}

var FacilityGroupManagerRolesMap = map[string][]string{
	FacilityGroupGeneral: {"ATM", "DATM"},
	FacilityGroupStaff:   {"ATM", "DATM"},
}

var AllGroupFacilities = []string{
	Albuquerque,
	Anchorage,
	Atlanta,
	Boston,
	Chicago,
	Cleveland,
	Denver,
	FortWorth,
	Honolulu,
	Houston,
	Indianapolis,
	Jacksonville,
	KansasCity,
	LosAngeles,
	Memphis,
	Miami,
	Minneapolis,
	NewYork,
	Oakland,
	SaltLake,
	Seattle,
	Washington,
}
