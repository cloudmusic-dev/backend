package database

import (
	"github.com/cloudmusic-dev/backend/configuration"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var DB *gorm.DB

func migrateDatabase() error {
	driver, err := mysql.WithInstance(DB.DB(), &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql", driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func InitializeDatabase(configuration configuration.Configuration) {
	var err error
	DB, err = gorm.Open("mysql", configuration.Database.Username+":"+configuration.Database.Password+"@tcp("+configuration.Database.Host+")/"+configuration.Database.Database+"?charset=utf8&parseTime=True&loc=Local&multiStatements=true")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB.LogMode(true)

	err = migrateDatabase()
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

func CloseDatabase() {
	err := DB.Close()
	if err != nil {
		log.Fatalf("Failed to close database connection: %v", err)
	}
}
