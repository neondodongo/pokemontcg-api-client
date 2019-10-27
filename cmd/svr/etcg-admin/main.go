package main

import (
	"fmt"
	"log"
	"net/http"
	"pokemontcg-api-client/internal/etcg"
	"pokemontcg-api-client/pkg/client"
	"pokemontcg-api-client/pkg/config"
	"pokemontcg-api-client/pkg/mongo"

	"github.com/gorilla/mux"
)

var con etcg.Controller

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg, err := config.GetConfig("cmd/svr/etcg-admin/config.json")
	if err != nil {
		log.Fatalf("error getting configuration [ %v ]", err)
	}
	db := mongo.InitDatabase(cfg)
	cli := client.InitializeClient(cfg)

	//InitiateControllers
	con := etcg.Controller{
		Client: cli,
		Mongo:  db,
		Config: cfg,
	}
	log.Println("log: Etcg-admin [ running ]")

	////Mux Router handling
	r := mux.NewRouter().StrictSlash(true)
	r.Handle("/cards/get", con.GetCards()).Methods(http.MethodGet)
	r.Handle("/sets/get", con.GetSets()).Methods(http.MethodGet)

	log.Println("Server is running on port: ", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, http.Handler(r)))
}

func PopulateDatabase() {

	fmt.Println("Getting sets")
	sets, err := con.GetAllSets()
	if err != nil {
		log.Printf("unable to get sets [%v]", err)
	}

	for _, set := range sets.Sets {
		err := con.Mongo.Upsert(set)
		if err != nil {
			log.Printf("failed to insert set %v [%v]", set.Code, err)
		}
		cards := con.GetCardsBySetCode(set.Code)
		if len(cards.Cards) != set.TotalCards {
			log.Printf("[WARNING] did not receive all card in set [ %s ] - actual [ %d ] / expected: [ %d ]", set.Name, len(cards.Cards), set.TotalCards)
		}

		for _, card := range cards.Cards {
			err := con.Mongo.Upsert(card)
			if err != nil {
				log.Printf("failed to insert card %v [%v]", card.ID, err)
			}
		}
	}
}