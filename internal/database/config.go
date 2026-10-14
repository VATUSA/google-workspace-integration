package database

import _ "github.com/joho/godotenv/autoload"
import "os"

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	// TLS is passed to the driver's tls= DSN option when set. "true" encrypts and
	// verifies the server certificate against the system roots, which is what Azure
	// Database for MySQL needs (require_secure_transport=ON). Leave unset for
	// DigitalOcean, whose CA isn't in the system roots.
	TLS string
}

var PrimaryDatabaseConfig = DatabaseConfig{
	Host:     os.Getenv("DB_PRIMARY_HOST"),
	Port:     os.Getenv("DB_PRIMARY_PORT"),
	User:     os.Getenv("DB_PRIMARY_USER"),
	Password: os.Getenv("DB_PRIMARY_PASSWORD"),
	Database: os.Getenv("DB_PRIMARY_DATABASE"),
	TLS:      os.Getenv("DB_TLS"),
}
