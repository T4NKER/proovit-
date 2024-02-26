package unspent

import (
	"errors"
	"log"

	database "proovit-/src/database"
	helpers "proovit-/src/helpers"
	models "proovit-/src/models"
)

func GetUnspentTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := database.DB.Model(&models.Transaction{}).Where("spent = ?", false).Order("amount ASC").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func MarkUnspentTransactionsAsSpent(transactions []models.Transaction) error {
	for _, transaction := range transactions {
		transaction.Spent = true
		if err := database.DB.Model(&models.Transaction{}).Where("transactionID = ?", transaction.TransactionID).Updates(&transaction).Error; err != nil {
			return err
		}
	}
	return nil
}

func CreateUnspentTransaction(leftoverAmount float64) error {
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

	if err := database.DB.Create(&transaction).Error; err != nil {
		log.Println("Error creating unspent transaction:", err)
		return err
	}

	return nil
}
