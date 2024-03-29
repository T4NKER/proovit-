package api

import (
	handlers "proovit-/pkg/handlers"

	"github.com/gin-gonic/gin"
)

type API struct {
	Router *gin.Engine
}

func (a *API) Initalize() {
	router := gin.Default()
	a.Router = router
	router.LoadHTMLGlob("./web/index.html")


	router.GET("/", handlers.RootHandler)
	router.GET("/transactions", handlers.ListAllTransactionsHandler)
	router.GET("/currentBalance", handlers.CurrentBalanceHandler)
	router.POST("/newTransfer", handlers.NewTransferHandler)
}
