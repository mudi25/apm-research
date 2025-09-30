package pgsql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDialector(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}
