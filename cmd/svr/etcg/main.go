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


func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg, err := config.GetConfig("cmd/svr/etcg/config.json")
	if err != nil {
		log.Fatalf("error getting configuration [ %v ]", err)
	}

	db := mongo.InitDatabase(cfg)
	cli := client.InitializeClient(cfg)

	//InitiateControllers
	con := etcg.Controller{
		Config: cfg,
		Mongo:  db,
		Client: cli,
	}
	//
	log.Println("ELITE TRAINER TCG")
	log.Println("Server is running on port: ", cfg.Port)

	////Mux Router handling
	r := mux.NewRouter().StrictSlash(true)
	r.Handle("/cards/get", con.GetCards()).Methods(http.MethodGet)
	r.Handle("/sets/get", con.GetSets()).Methods(http.MethodGet)

	//run server on port
	log.Fatal(http.ListenAndServe(":3000", http.Handler(r)))

}
