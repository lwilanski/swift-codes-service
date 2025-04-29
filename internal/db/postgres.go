package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		env("PGHOST", "db"),
		env("PGUSER", "swift"),
		env("PGPASSWORD", "swift"),
		env("PGDATABASE", "swift"),
		env("PGPORT", "5432"),
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
