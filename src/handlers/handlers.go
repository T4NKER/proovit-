package handlers

import (
	"log"
	"proovit-/src/converter"
	database "proovit-/src/database"
	services "proovit-/src/services"

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
	conversionRate, err := converter.BTCEURConverter()
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
		log.Println("ERROR", err)
		return
	}

	conversionRate, err := converter.BTCEURConverter()
	if err != nil {
		log.Println("ERROR", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = services.NewTransfer(request.AmountInEUR, database.DB, conversionRate)
	if err != nil {
		log.Println("ERROR", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Transfer completed successfully"})
}
