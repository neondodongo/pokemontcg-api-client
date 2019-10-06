package client

import (
	"encoding/json"
	"log"
	"net/http"
)

type HealthResponse struct{
	Status int `json:"status"`
	Message string `json:"message"`
}

func InitializeClient() *http.Client{
	c := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       60,
	}

	return &c
}

func GetRequest(uri string, ){

}

func RespondWithPrettyJSON(w http.ResponseWriter, statusCode int, payload interface{}){
	bSlice, err := json.MarshalIndent(payload, "", "  ")
	if err != nil{
		log.Println("failed to marshal payload")
	}
	w.Write(bSlice)
}