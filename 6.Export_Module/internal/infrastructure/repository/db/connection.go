package db

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// Create connection for SQL Server
func NewConnectionSQLServer(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	return db, err
}
