package services

import (

	models "proovit-/src/models"
	unspent "proovit-/src/services/unspent"

	"gorm.io/gorm"
)

func ListAllTransactions(db *gorm.DB) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := db.Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func MarkUnspentTransactions(transactions []models.Transaction, db *gorm.DB) error {
	return unspent.MarkUnspentTransactionsAsSpent(transactions, db)
}

func CreateUnspentTransaction(leftoverAmount, checkerAmount, conversionRate float64, db *gorm.DB) error {
	if leftoverAmount > checkerAmount {
		leftoverAmount = leftoverAmount * conversionRate
		return unspent.CreateUnspentTransaction(leftoverAmount, db)
	}
	return nil
}
