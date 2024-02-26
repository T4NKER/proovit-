package database

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
const DbPath = "./src/database/database.db"

func init() {
	DB = databaseInit(DbPath)
}

func databaseInit(path string) *gorm.DB {
	DB, err := gorm.Open(sqlite.Open(DbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	}

	/* err = db.AutoMigrate(&models.Transaction{})
	if err != nil {
		log.Fatal("Failed to auto migrate model")
	} */

	return DB
}