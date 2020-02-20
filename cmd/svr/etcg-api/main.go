package main

import (
	"net/http"
	"pokemontcg-api-client/internal/etcg"
	"pokemontcg-api-client/pkg/client"
	"pokemontcg-api-client/pkg/config"
	"pokemontcg-api-client/pkg/mongo"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func main() {

	cfg, err := config.GetConfig("./config.json")
	if err != nil {
		log.Fatal().Err(err).Msg("error getting configuration")
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
	log.Info().Msg("ELITE TRAINER TCG")
	log.Info().Msgf("Server is running on port: %v", con.Config.Port)

	////Mux Router handling
	router := mux.NewRouter().StrictSlash(true)

	if err := etcg.Handle(router, con); err != nil {
		log.Fatal().Err(err).Msg("failed to map URIs with application controller")
	}

	//run server on port
	log.Fatal().Err(http.ListenAndServe(":"+con.Config.Port, http.Handler(router))).Msg("SERVER ERROR")

}
