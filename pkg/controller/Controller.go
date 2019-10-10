package controller

import (
	"pokemontcg-api-client/pkg/client"
	"pokemontcg-api-client/pkg/config"
	"pokemontcg-api-client/pkg/mongo"
)

type Controller struct {
	Config config.Config
	Mongo  mongo.MongoBongo
	Client client.WebClient
}