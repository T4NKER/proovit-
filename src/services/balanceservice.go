package services

import (
	models "proovit-/src/models"
	queries "proovit-/src/services/queries"

	"gorm.io/gorm"
)

func AccountBalance(db *gorm.DB, conversionRate float64) (models.Balance, error) {
	var balance models.Balance

	transactions, err := queries.GetUnspentTransactions(db)
	if err != nil {
		return balance, err
	}

	for _, transaction := range transactions {
		balance.AmountInBTC += transaction.Amount
	}

	balance.AmountInEUR = balance.AmountInBTC / conversionRate

	return balance, nil
}


