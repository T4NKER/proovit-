package models

import (
	"time"
)

type Transaction struct {
	TransactionID string    `gorm:"column:transactionID" json:"transctionID"`
	Amount        float64   `json:"amount"`
	Spent         bool      `json:"spent"`
	CreatedAt     time.Time `gorm:"column:createdAt" json:"createdAt"`
}

type Balance struct {
	AmountInEUR float64 `json:"amountInEUR"`
	AmountInBTC float64 `json:"amountInBTC"`
}
