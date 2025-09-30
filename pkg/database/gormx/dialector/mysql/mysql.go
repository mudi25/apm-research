package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDialector(dsn string) gorm.Dialector {
	return mysql.Open(dsn)
}
