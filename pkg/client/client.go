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

func InitializeClient(c config.Config) (client http.Client) {
	client = http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       60 * time.Second,
	}
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
