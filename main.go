package main

import (

	services "proovit-/services"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)


func main() {

	router := gin.Default()

	router.GET("/transactions", services.listAllTransactions())
	router.GET("/currentBalance", services.currentBalance())
	router.POST("/newTransfer", services.newTransfer())

	router.Run("localhost:8080")

}
