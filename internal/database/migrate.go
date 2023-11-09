package database

func MigrateDB() error {
	err := DB.AutoMigrate(
		&Account{},
	)
	if err != nil {
		return err
	}
	return nil
}
