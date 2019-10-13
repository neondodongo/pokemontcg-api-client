package main

import (
	"log"
	"net/http"
	"pokemontcg-api-client/internal/etcg"
	"pokemontcg-api-client/pkg/client"
	"pokemontcg-api-client/pkg/config"
	"pokemontcg-api-client/pkg/mongo"

	"github.com/gorilla/mux"
)

var con *etcg.PokemonTCGController

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg, err := config.GetConfig("cmd/svr/etcg-admin/config.json")
	if err != nil {
		log.Fatalf("error getting configuration [ %v ]", err)
	}
	db := mongo.InitDatabase(cfg)
	cli := client.InitializeClient(cfg)

	//InitiateControllers
	con = &etcg.PokemonTCGController{
		Client: cli.Client,
		Mongo:  &db,
		Config: cfg,
	}
	log.Println("log: Etcg-admin [ running ]")

	////Mux Router handling
	r := mux.NewRouter().StrictSlash(true)
	//
	////Handler functions for endpoints
	//r.Handle("/health", client.Health()).Methods(http.MethodGet)
	r.HandleFunc("/dashboard", con.GetDashboard).Methods(http.MethodGet)

	log.Println("Server is running on port: ", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, http.Handler(r)))
}
