package card

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"pokemontcg-api-client/pkg/client"
	"pokemontcg-api-client/pkg/controller"
	"pokemontcg-api-client/pkg/dto"
)

func GetCards(c controller.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		params, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			log.Println("Unable to read query")
		}
		fmt.Printf("Parameters: %s\n", params)

		var cards []dto.Card
		if len(params) > 0 {
			cards = c.Mongo.GetFilterCards(params)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		client.RespondWithPrettyJSON(w, 200, cards)
	})
}

