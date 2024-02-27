package handlers

import (
	database "proovit-/src/database"
	"proovit-/src/helpers"
	services "proovit-/src/services"
	"proovit-/src/services/unspent"

	"github.com/gin-gonic/gin"

)

func RootHandler(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{})
}

func ListAllTransactionsHandler(c *gin.Context) {
	transactions, err := services.ListAllTransactions(database.DB)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve transactions"})
		return
	}

	c.JSON(200, transactions)
}

func CurrentBalanceHandler(c *gin.Context) {
	conversionRate, err := helpers.BTCEURConverter()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve conversionrate"})
		return
	}
	balance, err := services.AccountBalance(database.DB, conversionRate)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve balance"})
		return
	}
	c.JSON(200, balance)
}

func NewTransferHandler(c *gin.Context) {
	var request struct {
		AmountInEUR float64 `json:"amount_in_eur" binding:"required"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	conversionRate, err := helpers.BTCEURConverter()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := helpers.IsValidTransferAmount(request.AmountInEUR, conversionRate); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	unSpentTransactions, err := unspent.GetUnspentTransactions(database.DB)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	selectedTransactions, totalAmount, err := helpers.CalculateTotalAmount(unSpentTransactions, conversionRate, request.AmountInEUR)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	checkerAmount := conversionRate * 0.0001

	if err := helpers.CheckFunds(totalAmount, request.AmountInEUR, checkerAmount); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = services.MarkUnspentTransactions(selectedTransactions, database.DB)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	leftOverAmount := totalAmount - request.AmountInEUR

	if err := services.CreateUnspentTransaction(leftOverAmount, checkerAmount, conversionRate, database.DB); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Transfer completed successfully"})
}

