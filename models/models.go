package models

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"


)

type Transaction struct {
   gorm.Model
   TransactionID int       `json:"transactionID"`
   Amount        int       `json:"amountInBTC"`
   Spent         bool      `json:"spent"`
   CreatedAt     time.Time `json:"createdAt"`
}