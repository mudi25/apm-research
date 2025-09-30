package mysql

import (
	apmmysql "go.elastic.co/apm/module/apmgormv2/v2/driver/mysql"
	"gorm.io/gorm"
)

func NewDialector(dsn string) gorm.Dialector {
	return apmmysql.Open(dsn)
}
