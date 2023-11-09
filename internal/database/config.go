package database

import _ "github.com/joho/godotenv/autoload"
import "os"

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

var PrimaryDatabaseConfig = DatabaseConfig{
	Host:     os.Getenv("DB_PRIMARY_HOST"),
	Port:     os.Getenv("DB_PRIMARY_PORT"),
	User:     os.Getenv("DB_PRIMARY_USER"),
	Password: os.Getenv("DB_PRIMARY_PASSWORD"),
	Database: os.Getenv("DB_PRIMARY_DATABASE"),
}
