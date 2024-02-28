package main

import (
	api "proovit-/pkg"
)

func main() {
	API := api.API{}
	API.Initalize()
	API.Router.Run("localhost:8080")
}
