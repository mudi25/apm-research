package sqlserver

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func NewDialector(dsn string) gorm.Dialector {
	return sqlserver.Open(dsn)
}
