package config

// These overrides exist for backwards compatibility with email addresses that already existed before the new scheme
// Do NOT add new users to this list. Any newly created emails should follow the email addressing scheme
var PrimaryEmailOverrides = map[uint64]string{
	1181029: "b.barrett",
	1394476: "b.wening",
	1505109: "j.owens",
	1350061: "j.west",
	1371112: "s.oneill",
	1293257: "r.patel",
	1241028: "j.kerr",
	1354520: "b.brody",
	1313538: "m.campbell",
	1371395: "m.bruck",
}
