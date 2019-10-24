package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Port    string `json:"port"`
	Timeout int    `json:"timeout"`
	Mongo   struct {
		Url             string `json:"url"`
		Database        string `json:"database"`
		CardsCollection string `json:"cardsCollection"`
		SetsCollection  string `json:"setsCollection"`
		UsersCollection string `json:"usersCollection"`
		SetLimit		int64	   `json:"setLimit"`
	}
	PokemonAPI string `json:"pokemonAPI"`
}

func GetConfig(path string) (c Config, err error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	err = json.Unmarshal(content, &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}
