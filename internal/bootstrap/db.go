package bootstrap

import (
	"errors"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg DBConfig, logLevel logger.LogLevel) (*gorm.DB, error) {
	driver := strings.ToLower(cfg.Driver)
	var dialector gorm.Dialector
	if driver == "mysql" {
		dialector = mysql.Open(cfg.DSN)
	} else if driver == "postgres" || driver == "postgresql" {
		dialector = postgres.Open(cfg.DSN)
	} else {
		return nil, errors.New("unsupported db driver")
	}

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.ConnMaxLifetimeSecond > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeSecond) * time.Second)
	}

	return db, nil
}
