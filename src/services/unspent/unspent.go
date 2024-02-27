package unspent

import (
	"errors"
	"log"

	helpers "proovit-/src/helpers"
	models "proovit-/src/models"

	"gorm.io/gorm"
)

func GetUnspentTransactions(db *gorm.DB) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := db.Model(&models.Transaction{}).Where("spent = ?", false).Order("amount ASC").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func MarkUnspentTransactionsAsSpent(transactions []models.Transaction, db *gorm.DB) error {
	for _, transaction := range transactions {
		transaction.Spent = true
		if err := db.Model(&models.Transaction{}).Where("transactionID = ?", transaction.TransactionID).Updates(&transaction).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateUnspentTransaction(leftoverAmount float64, db *gorm.DB) error {
	if leftoverAmount < 0 {
		return errors.New("leftover amount must be non-negative")
	}

	hexadecimal, err := helpers.RandomHex(16)
	if err != nil {
		log.Println("Error creating random hexadecimal:", err)
		return err
	}

	transaction := models.Transaction{
		TransactionID: hexadecimal,
		Amount:        leftoverAmount,
		Spent:         false,
	}

	if err := db.Create(&transaction).Error; err != nil {
		log.Println("Error creating unspent transaction:", err)
		return err
	}

	return nil
}
