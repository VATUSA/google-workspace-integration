package database

func MigrateDB() error {
	err := DB.AutoMigrate(
		&Account{},
		&Alias{},
		&Group{},
		&GroupMembership{},
	)
	if err != nil {
		return err
	}
	return nil
}
