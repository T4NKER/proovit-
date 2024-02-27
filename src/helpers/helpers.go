package helpers

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func IsAmountInEur0OrSmaller(amountInEUR float64) bool {
	return amountInEUR <= 0 
}

func IsAmountInBTCLessThanAmount(amountInBTC float64, amount float64) bool {
	return amountInBTC < amount 
}
