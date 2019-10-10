package client

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"pokemontcg-api-client/pkg/config"
	"time"
)

type HealthResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type WebClient struct {
	Client *http.Client
}

func InitializeClient(c config.Config) (w WebClient) {

	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       60 * time.Second,
	}

	w.Client = &client
	return
}

func (c *WebClient)RespondWithPrettyJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	bSlice, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Println("failed to marshal payload")
	}
	w.Write(bSlice)
}

func (c *WebClient)Get(uri string){

}

func DecodeInterface(io io.Reader, t interface{}) error {
	content, err := ioutil.ReadAll(io)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &t)
	if err != nil {
		return err
	}
	return nil
}
