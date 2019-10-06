package main

import (
	"fmt"
	"pokemontcg-api-client/qops"
)

const (
	DB_URI  = "mongodb+srv://neonorca:nakedgyp13@testcluster-i7bzt.gcp.mongodb.net/admin?retryWrites=true&w=majority"
	DB_NAME = "Pokemon-Elite-TCG"
)

var db *qops.MongoBongo

func main() {

	fmt.Println("Program started.")

	db = qops.InitDatabase(DB_URI, DB_NAME)
	db.PopulateDB()
	//db.Count("Cards")

	fmt.Printf("**********PROGRAM END**********")

}
