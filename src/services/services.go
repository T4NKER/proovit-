package services

import (
	"errors"
	"log"

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

func AccountBalance(db *gorm.DB, conversionRate float64) (models.Balance, error) {
	var balance models.Balance

	transactions, err := unspent.GetUnspentTransactions(db)
	if err != nil {
		return balance, err
	}

	for _, transaction := range transactions {
		balance.AmountInBTC += transaction.Amount
	}

	balance.AmountInEUR = balance.AmountInBTC / conversionRate

	return balance, nil
}

func NewTransfer(amountInEUR float64, db *gorm.DB, conversionRate float64) error {
	log.Printf("Attempting to transfer: %.2f EUR", amountInEUR)
	if amountInEUR <= 0 {
		return errors.New("cannot transfer negative amount or 0")
	}

	amountInBTC := amountInEUR * conversionRate
	if amountInBTC < 0.00001 {
		return errors.New("cannot transfer less than 0.00001 BTC")
	}

	unspentTransactions, err := unspent.GetUnspentTransactions(db)
	if err != nil {
		log.Println("Error getting unspent transactions:", err)
		return err
	}

	totalAmount := 0.0
	selectedTransactions := make([]models.Transaction, 0)

	for _, transaction := range unspentTransactions {
		totalAmount += transaction.Amount / conversionRate
		selectedTransactions = append(selectedTransactions, transaction)
		if totalAmount >= amountInEUR {
			break
		}
	}

	// When the sum is small enough, it should still be possible to
	// conduct the transfer. This is to prevent a situation where the
	// user actually has the money but the transaction fails due to
	// a very small amount.
	// In the current market, when the BTC is about 50 000 EUR, one eur is equal
	// to about 0.000019 BTC meaning that 1/100 of a  cent is about 0.0000000019 BTC. (27/02/2024)

	checkerAmount := conversionRate * 0.0001

	if totalAmount-amountInEUR < -checkerAmount {
		return errors.New("not enough funds")
	}

	err = unspent.MarkUnspentTransactionsAsSpent(selectedTransactions, db)
	if err != nil {
		log.Println("Error marking unspent transactions as spent:", err)
		return err
	}

	leftoverAmount := totalAmount - amountInEUR

	if leftoverAmount > checkerAmount {
		leftoverAmount = leftoverAmount * conversionRate
		err = unspent.CreateUnspentTransaction(leftoverAmount, db)
		if err != nil {
			log.Println("Error creating new unspent transaction:", err)
			return err
		}
	}

	return nil
}
