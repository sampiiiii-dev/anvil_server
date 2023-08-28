package db

import (
	"fmt"
	"github.com/sampiiiii-dev/anvil_server/anvil/config"
	"github.com/sampiiiii-dev/anvil_server/anvil/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

type Database struct {
	db    *gorm.DB
	dbErr error
}

var once sync.Once
var instance *gorm.DB

func GetDBInstance() *gorm.DB {
	once.Do(func() {
		cfg := config.GetConfigInstance(nil)
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
			cfg.DB.Host, cfg.DB.User, cfg.DB.Pass, cfg.DB.Database, cfg.DB.Port)

		var err error
		instance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("Failed to connect to database: " + err.Error())
		}

		if err := MigrateDB(instance); err != nil {
			panic("Failed to migrate database: " + err.Error())
		}
	})
	return instance
}

// PingDB pings the database to ensure the connection is alive.
func PingDB() error {
	db := GetDBInstance()
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.SignInRecord{}, &models.MailingList{})
}

// GetDBMetrics returns some metrics about the database connection.
func GetDBMetrics() (int, int, int, error) {
	db := GetDBInstance()
	sqlDB, err := db.DB()
	if err != nil {
		return 0, 0, 0, err
	}

	return sqlDB.Stats().OpenConnections,
		sqlDB.Stats().InUse,
		sqlDB.Stats().Idle,
		nil
}

// CloseDB closes the database connection. Call this during graceful shutdown.
func CloseDB() error {
	db := GetDBInstance()
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
