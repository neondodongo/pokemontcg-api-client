package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"pokemontcg-api-client/pkg/client"
	"pokemontcg-api-client/pkg/config"
	"pokemontcg-api-client/pkg/mongo"
)


func main() {

	cfg, err := config.GetConfig("cmd/svr/etcg/config.json")
	if err != nil{
		log.Fatalf("error getting configuration [ %v ]", err)
	}
	cli := client.InitializeClient()
	db := mongo.InitDatabase(cfg)

	log.Println("log: Elite Trainer TCG [ running ]")
	log.Println("Server is running on port: ", cfg.Port)

	//Mux Router handling
	r := mux.NewRouter().StrictSlash(true)

	//Handler functions for endpoints
	r.Handle("/health", Health()).Methods(http.MethodGet)
	r.Handle("/populate-database", db.PopulateDatabase()).Methods(http.MethodGet)

	//run server on port
	log.Fatal(http.ListenAndServe( ":3000", http.Handler(r)))
}

func Health() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var body client.HealthResponse
		body.Status = http.StatusOK
		body.Message = "Application is up and running"

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(body)
	})
}