package main

import (
	"log"
	"net/http"
	"pokemontcg-api-client/pkg/card"
	"pokemontcg-api-client/pkg/client"
	"pokemontcg-api-client/pkg/config"
	"pokemontcg-api-client/pkg/controller"
	"pokemontcg-api-client/pkg/mongo"
	"pokemontcg-api-client/pkg/tcg"

	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.GetConfig("cmd/svr/etcg/config.json")
	if err != nil {
		log.Fatalf("error getting configuration [ %v ]", err)
	}

	db := mongo.InitDatabase(cfg)
	cli := client.InitializeClient(cfg)

	//InitiateControllers
	con := controller.Controller{
		Config: cfg,
		Mongo:  db,
		Client: cli,
	}

	log.Println("log: Elite Trainer TCG [ running ]")
	log.Println("Server is running on port: ", cfg.Port)

	//Mux Router handling
	r := mux.NewRouter().StrictSlash(true)

	//Handler functions for endpoints
	r.Handle("/health", client.Health()).Methods(http.MethodGet)
	r.Handle("/populate-database", tcg.PopulateDatabase(con)).Methods(http.MethodGet)
	r.Handle("/cards", card.GetCards(con)).Methods(http.MethodGet)

	//run server on port
	log.Fatal(http.ListenAndServe(":3000", http.Handler(r)))
}
