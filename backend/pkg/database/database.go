package database

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type Config struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
}

type DB struct {
	*gorm.DB
}

func InitMySQL(cfg Config) (*DB, error) {
	if cfg.Host == "" || cfg.User == "" || cfg.DBName == "" {
		return nil, errors.New("database config is incomplete")
	}
	if cfg.Port == 0 {
		cfg.Port = 3306
	}
	if cfg.Charset == "" {
		cfg.Charset = "utf8mb4"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local&collation=utf8mb4_unicode_ci", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset)
	gdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
		Logger:                                   gormLogger.Default.LogMode(gormLogger.Warn),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	return &DB{DB: gdb}, nil
}
