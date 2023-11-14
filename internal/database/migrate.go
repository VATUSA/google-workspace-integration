package database

func MigrateDB() error {
	err := DB.AutoMigrate(
		&Account{},
		&Domain{},
	)
	if err != nil {
		return err
	}
	return nil
}
