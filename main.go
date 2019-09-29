package main

import (
	"fmt"
	"pokemontcg-api-client/routes"
)

func main() {

	fmt.Println("Program started.")

	// Get all set data from Pokemon TCG API
	routes.GetAllSets()

	fmt.Printf("**********PROGRAM END**********")

}
