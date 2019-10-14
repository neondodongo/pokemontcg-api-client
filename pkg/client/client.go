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

func RespondWithPrettyJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	bSlice, e := json.MarshalIndent(payload, "", "  ")
	if e != nil {
		log.Println("failed to marshal payload: ", e)
	}
	_, e = w.Write(bSlice)
	if e != nil {
		log.Println("error writing byte slice: ", e)
	}
}

func (c *WebClient) Get(uri string) {

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
