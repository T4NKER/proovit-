package services

import (
	"errors"
	"log"

	converter "proovit-/src/converter"
	database "proovit-/src/database"
	models "proovit-/src/models"
	unspent "proovit-/src/services/unspent"
)

func ListAllTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := database.DB.Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func AccountBalance() (models.Balance, error) {
	var balance models.Balance

	transactions, err := unspent.GetUnspentTransactions()
	if err != nil {
		return balance, err
	}

	for _, transaction := range transactions {
		balance.AmountInBTC += transaction.Amount
	}

	balance.AmountInEUR, err = converter.BTCEURConverter(balance.AmountInBTC, "BTCTOEUR")
	if err != nil {
		log.Println(err)
		return balance, err
	}

	return balance, nil
}

func NewTransfer(amountInEUR float64) error {
	if amountInEUR <= 0 {
		return errors.New("cannot transfer negative amount or 0")
	}

	amountInBTC, err := converter.BTCEURConverter(amountInEUR, "EURTOBTC")
	if err != nil {
		log.Println("Error converting EUR to BTC", err)
		return err
	}
	if amountInBTC < 0.00001 {
		return errors.New("cannot transfer less than 0.00001 BTC")
	}

	unspentTransactions, err := unspent.GetUnspentTransactions()
	if err != nil {
		log.Println("Error getting unspent transactions:", err)
		return err
	}

	totalAmount := 0.0
	selectedTransactions := make([]models.Transaction, 0)

	for _, transaction := range unspentTransactions {
		transactionAmountInEUR, _ := converter.BTCEURConverter(transaction.Amount, "BTCTOEUR")
		totalAmount += transactionAmountInEUR
		selectedTransactions = append(selectedTransactions, transaction)
		if totalAmount >= amountInEUR {
			break
		}
	}

	if totalAmount < amountInEUR {
		return errors.New("not enough funds")
	}

	err = unspent.MarkUnspentTransactionsAsSpent(selectedTransactions)
	if err != nil {
		log.Println("Error marking unspent transactions as spent:", err)
		return err
	}

	leftoverAmount := totalAmount - amountInEUR

	if leftoverAmount > 0 {
		leftoverAmount, _ = converter.BTCEURConverter(leftoverAmount, "EURTOBTC")
		err = unspent.CreateUnspentTransaction(leftoverAmount)
		if err != nil {
			log.Println("Error creating new unspent transaction:", err)
			return err
		}
	}

	return nil
}
