package database

func MigrateDB() error {
	err := DB.AutoMigrate(
		&Account{},
		&Alias{},
		&FallbackAlias{},
		&GroupAlias{},
		&Group{},
		&GroupMembership{},
	)
	if err != nil {
		return err
	}
	return nil
}
