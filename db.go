package main

import (
	"errors"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDB(dbKind string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	if dbKind == "mysql" {
		// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
		dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else if dbKind == "sqlite" {
		db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	} else {
		return nil, errors.New("Unsupported database")
	}
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}
