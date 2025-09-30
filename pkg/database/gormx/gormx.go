package gormx

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// PoolConfig defines the database connection pool settings.
// If not provided, Go's sql.DB will use its own default values.
type PoolConfig struct {
	MaxOpenCon     int           // Maximum number of open connections to the database.
	MaxIdleCon     int           // Maximum number of idle connections in the pool.
	MaxLifetimeCon time.Duration // Maximum amount of time a connection may be reused.
	MaxIdleTimeCon time.Duration // Maximum amount of time a connection may remain idle.
}

// Config defines the configuration for the GORM client.
type Config struct {
	Dialector  gorm.Dialector
	PoolConfig *PoolConfig  // Optional connection pool configuration.
	GormConfig *gorm.Config // Optional GORM configuration (can be nil).
}

// NewClient initializes and returns a new GORM DB client.
//
// By default, it uses the APM instrumented MySQL driver
// so that all queries are automatically traced by Elastic APM.
//
// Steps performed:
// 1. Open the database connection with GORM.
// 2. Apply connection pool configuration if provided.
// 3. Ping the database to ensure the connection is alive.
func NewClient(cfg Config) (*gorm.DB, error) {
	if cfg.GormConfig == nil {
		cfg.GormConfig = &gorm.Config{}
	}
	db, err := gorm.Open(cfg.Dialector, cfg.GormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if cfg.PoolConfig != nil {
		sqlDB.SetMaxOpenConns(cfg.PoolConfig.MaxOpenCon)
		sqlDB.SetMaxIdleConns(cfg.PoolConfig.MaxIdleCon)
		sqlDB.SetConnMaxLifetime(cfg.PoolConfig.MaxLifetimeCon)
		sqlDB.SetConnMaxIdleTime(cfg.PoolConfig.MaxIdleTimeCon)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Disconnect closes the underlying sql.DB connection.
//
// Should be called when the application is shutting down
// to gracefully release all database resources.
func Disconnect(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from gorm.DB: %s", err.Error())
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %s", err.Error())
	}

	return nil
}
