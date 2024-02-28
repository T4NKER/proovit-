package services

import (

	models "proovit-/pkg/models"
	queries "proovit-/pkg/services/queries"

	"gorm.io/gorm"
)

func ListAllTransactions(db *gorm.DB) ([]models.Transaction, error) {
	var transactions []models.Transaction
	transactions, err := queries.GetAllTransactions(db)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func MarkUnspentTransactions(transactions []models.Transaction, db *gorm.DB) error {
	return queries.MarkUnspentTransactionsAsSpent(transactions, db)
}

func GetUnspentTransactions(db *gorm.DB) ([]models.Transaction, error) {
	return queries.GetUnspentTransactions(db)
}

func CreateUnspentTransaction(leftoverAmount, checkerAmount, conversionRate float64, db *gorm.DB) error {
	if leftoverAmount > checkerAmount {
		leftoverAmount = leftoverAmount * conversionRate
		return queries.CreateUnspentTransaction(leftoverAmount, db)
	}
	return nil
}
