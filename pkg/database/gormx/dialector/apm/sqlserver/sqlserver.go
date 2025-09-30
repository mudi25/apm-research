package sqlserver

import (
	sqlserver "go.elastic.co/apm/module/apmgormv2/v2/driver/sqlserver"
	"gorm.io/gorm"
)

func NewDialector(dsn string) gorm.Dialector {
	return sqlserver.Open(dsn)
}
