package pgsql

import (
	postgres "go.elastic.co/apm/module/apmgormv2/v2/driver/postgres"
	"gorm.io/gorm"
)

func NewDialector(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}

