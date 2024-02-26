package main

import (
	handlers "proovit-/src/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("./src/templates/*")

	router.GET("/", handlers.RootHandler)
	router.GET("/transactions", handlers.ListAllTransactionsHandler)
	router.GET("/currentBalance", handlers.CurrentBalanceHandler)
	router.POST("/newTransfer", handlers.NewTransferHandler)

	router.Run("localhost:8080")

}
