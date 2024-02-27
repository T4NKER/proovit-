package services

import (
	"io/ioutil"
	"log"
	"testing"
	"fmt"

	converter "proovit-/src/converter"
	models "proovit-/src/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func localConverter() (rate float64) {
	rate, _ = converter.BTCEURConverter()
	return rate
}

func TestNewTransfer(t *testing.T) {
	setupTestDatabase()

	conversionRate := localConverter()
	fmt.Println()
	LocalAccountBalance(conversionRate)
	fmt.Println()
	LocalListAllTransactions()
	fmt.Println()
	

	// NOTE:
	// The total unspent amount of freshly initalized database is 26 BTC.
	// Since it uses the same database for all tests, we need to make sure that
	// the conversion rate is the same for all tests.
	// Otherwise, the second test pass/fail would be inconsistent.

	tests := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{
			name:    "Valid transfer",
			amount:  16.0 / conversionRate,
			wantErr: false,
		},
		{
			name:    "Valid transfer with remaining funds",
			amount:  10.0 / conversionRate,
			wantErr: false,
		},

		{
			name:    "Invalid transfer with insufficient funds",
			amount:  200.0 / conversionRate,
			wantErr: true,
		},
		{
			name:    "Invalid transfer with negative amount",
			amount:  -50.0 / conversionRate,
			wantErr: true,
		},
		{
			name:    "Invalid transfer with zero amount",
			amount:  0.0 / conversionRate,
			wantErr: true,
		},
		{
			name:    "Invalid transfer with less than minimum BTC",
			amount:  0.000009 / localConverter(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewTransfer(tt.amount, db, conversionRate)
			if tt.wantErr {
				assert.Error(t, err, "Expected an error but got nil")
			} else {
				assert.NoError(t, err, "Expected no error but got an error")
			}
		})
	}
	fmt.Println()
	LocalListAllTransactions()
	fmt.Println()
	LocalAccountBalance(conversionRate)
	fmt.Println()
	cleanupTestDatabase()
}

func LocalAccountBalance(conversionRate float64) {
	balance, err := AccountBalance(db, conversionRate)
	if err != nil {
		log.Fatal("Error getting account balance:", err)
	}
	log.Println("Current balance in EUR: ", balance.AmountInEUR ,"and in BTC: ", balance.AmountInBTC)
}

func LocalListAllTransactions() {
	transactions, err := ListAllTransactions(db)
	if err != nil {
		log.Fatal("Error getting transactions:", err)
	}
	for _, transaction := range transactions {
		log.Println("Transaction: ", transaction)
	}
}

func setupTestDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("Error setting up test database:", err)
	}

	err = db.AutoMigrate(&models.Transaction{})
	if err != nil {
		log.Fatal("Error migrating tables:", err)
	}

	err = loadSeedData(db)
	if err != nil {
		log.Fatal("Error loading seed data:", err)
	}
}

func cleanupTestDatabase() {
	err := db.Migrator().DropTable(&models.Transaction{})
	if err != nil {
		log.Println("Error dropping tables:", err)
	}

}

func loadSeedData(db *gorm.DB) error {
	seedData, err := ioutil.ReadFile("../database/seedData.sql")
	if err != nil {
		return err
	}

	err = db.Exec(string(seedData)).Error
	if err != nil {
		return err
	}

	return nil
}
