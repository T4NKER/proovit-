package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"proovit-/pkg/models"
)

func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func IsValidTransferAmount(amountInEUR, conversionRate float64) error {
	if amountInEUR <= 0  {
		return errors.New("cannot transfer less than 0.00001 EUR")
	}

	amountInBTC := amountInEUR * conversionRate

	if amountInBTC < 0.00001 {
		return errors.New("cannot transfer less than 0.00001 BTC")
	}

	return nil
}

func CheckFunds(totalAmount, amountInEUR, checkerAmount float64) error {
	if totalAmount-amountInEUR < -checkerAmount {
		return errors.New("not enough funds")
	}
	return nil
}

func CalculateTotalAmount(unspentTransactions []models.Transaction, conversionRate, amountInEUR float64) ([]models.Transaction, float64, error) {
	totalAmount := 0.0
	selectedTransactions := make([]models.Transaction, 0)

	for _, transaction := range unspentTransactions {
		totalAmount += transaction.Amount / conversionRate
		selectedTransactions = append(selectedTransactions, transaction)
		if totalAmount >= amountInEUR {
			break
		}
	}

	return selectedTransactions, totalAmount, nil
}