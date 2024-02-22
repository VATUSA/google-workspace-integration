package database

func MigrateDB() error {
	err := DB.AutoMigrate(
		&Account{},
		&Alias{},
		&GroupAlias{},
		&Group{},
		&GroupMembership{},
	)
	if err != nil {
		return err
	}
	return nil
}
