package card

import (
	"github.com/gorilla/mux"
	"net/http"
	"pokemontcg-api-client/pkg/client"
	"pokemontcg-api-client/pkg/controller"
)

func GetCard(c controller.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		params := mux.Vars(r)
		cardId := params["cardId"]

		card := c.Mongo.GetCardById(cardId)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		client.RespondWithPrettyJSON(w, 200, card)
	})
}

