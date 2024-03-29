package database

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
const DbPath = "./pkg/database/database.db"

func init() {
	DB = databaseInit(DbPath)
}

func databaseInit(path string) *gorm.DB {
	DB, err := gorm.Open(sqlite.Open(DbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	}

	return DB
}