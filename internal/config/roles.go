package config

const (
	AirTrafficManager       = "ATM"
	DeputyAirTrafficManager = "DATM"
	TrainingAdministrator   = "TA"
	EventCoordinator        = "EC"
	FacilityEngineer        = "FE"
	WebMaster               = "WM"
	Instructor              = "INS"
	Mentor                  = "MTR"
	DICETeam                = "DICE"
	USWebTeam               = "USWT"
	EmailUser               = "EMAIL"
	SocialMediaTeam         = "SMT"
)

// USA Staff roles
// TODO: Fix this so that there is one common role that all USA Staff members have
const (
	US1 = "US1"
	US2 = "US2"
	US3 = "US3"
	US4 = "US4"
	US5 = "US5"
	US6 = "US6"
	US7 = "US7"
	US8 = "US8"
	US9 = "US9"
)

var AccountEntitlementRoles = []string{
	AirTrafficManager,
	DeputyAirTrafficManager,
	TrainingAdministrator,
	EventCoordinator,
	FacilityEngineer,
	WebMaster,
	Instructor,
	Mentor,
	DICETeam,
	USWebTeam,
	EmailUser,
	SocialMediaTeam,
	US1,
	US2,
	US3,
	US4,
	US5,
	US6,
	US7,
	US8,
	US9,
}

var FacilityNameAliasEntitlementRoles = []string{
	AirTrafficManager,
	DeputyAirTrafficManager,
	TrainingAdministrator,
	EventCoordinator,
	FacilityEngineer,
	WebMaster,
	Instructor,
	Mentor,
	EmailUser,
}
