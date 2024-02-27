package services

import (
	"errors"
	"log"
	"proovit-/src/helpers"
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

func NewTransfer(amountInEUR float64, db *gorm.DB, conversionRate float64) error {
	if helpers.IsAmountInEur0OrSmaller(amountInEUR) {
		return errors.New("cannot transfer less than 0 EUR")
	}

	amountInBTC := amountInEUR * conversionRate

	if helpers.IsAmountInBTCLessThanAmount(amountInBTC, 0.00001) {
		return errors.New("cannot transfer less than 0.00001 BTC")
	}

	unspentTransactions, err := unspent.GetUnspentTransactions(db)
	if err != nil {
		log.Println("Error getting unspent transactions:", err)
		return err
	}

	totalAmount, selectedTransactions, err := calculateTotalAmount(unspentTransactions, conversionRate, amountInEUR)
	if err != nil {
		return err
	}

	checkerAmount := conversionRate * 0.0001
	if err := checkFunds(totalAmount, amountInEUR, checkerAmount); err != nil {
		return err
	}

	err = markUnspentTransactions(selectedTransactions, db)
	if err != nil {
		log.Println("Error marking unspent transactions as spent:", err)
		return err
	}

	leftoverAmount := totalAmount - amountInEUR

	if err := createUnspentTransaction(leftoverAmount, checkerAmount, conversionRate, db); err != nil {
		return err
	}

	return nil
}

func calculateTotalAmount(unspentTransactions []models.Transaction, conversionRate, amountInEUR float64) (float64, []models.Transaction, error) {
	totalAmount := 0.0
	selectedTransactions := make([]models.Transaction, 0)

	for _, transaction := range unspentTransactions {
		totalAmount += transaction.Amount / conversionRate
		selectedTransactions = append(selectedTransactions, transaction)
		if totalAmount >= amountInEUR {
			break
		}
	}

	return totalAmount, selectedTransactions, nil
}

func checkFunds(totalAmount, amountInEUR, checkerAmount float64) error {
	if totalAmount-amountInEUR < -checkerAmount {
		return errors.New("not enough funds")
	}
	return nil
}

func markUnspentTransactions(transactions []models.Transaction, db *gorm.DB) error {
	return unspent.MarkUnspentTransactionsAsSpent(transactions, db)
}

func createUnspentTransaction(leftoverAmount, checkerAmount, conversionRate float64, db *gorm.DB) error {
	if leftoverAmount > checkerAmount {
		leftoverAmount = leftoverAmount * conversionRate
		return unspent.CreateUnspentTransaction(leftoverAmount, db)
	}
	return nil
}
