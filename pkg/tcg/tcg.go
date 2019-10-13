package tcg

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"pokemontcg-api-client/pkg/client"
	"pokemontcg-api-client/pkg/config"
	"pokemontcg-api-client/pkg/controller"
	"pokemontcg-api-client/pkg/dto"
	"pokemontcg-api-client/pkg/mongo"
	"strconv"
	"strings"
)

type PokemonTCGController struct {
	Client *http.Client
	Config config.Config
	Mongo  *mongo.MongoBongo
}

func GetAllSets(p controller.Controller) (*dto.Sets, error) {

	var sets dto.Sets

	//fmt.Printf("pokemon tcg api [ %v ]", p.Config.PokemonAPI)
	resp, err := http.Get(p.Config.PokemonAPI + "/sets")
	if err != nil {
		log.Printf("error getting sets from pokemon tcg api [%v]", err)
		return nil, err
	}

	//Read body into byt slice
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body [%v]", err)
		return nil, err
	}

	//Unmarshal byte slice into sets struct
	if err := json.Unmarshal(body, &sets); err != nil {
		log.Printf("error unmarshalling into sets [%v]", err)
		return nil, err
	}

	return &sets, nil
}

// GetCardsBySetCode will fetch a list of card by their setCode and return it
func GetCardsBySetCode(p controller.Controller, setCode string) *dto.Cards {
	setCode = strings.TrimSpace(setCode)
	if setCode != "" {

		cards := &dto.Cards{}
		err := GetPaginatedCards(p, cards, setCode)
		if err != nil {
			log.Printf("error calling pokemon tcg api [%v]", err)
		}

		return cards
	}
	return &dto.Cards{}
}

func GetPaginatedCards(p controller.Controller, cards *dto.Cards, setCode string) error {
	count := 1
	actualCards := make([]dto.Card, 0)

	for {
		//build query with setCode and replace <count>
		query := p.Config.PokemonAPI + "/cards?page=<count>&setCode=" + setCode
		uri := strings.Replace(query, "<count>", strconv.Itoa(count), 1)

		//call to pokemon tcg api
		resp, err := http.Get(uri)
		if err != nil {
			log.Printf("error in call to uri %v [%v]", uri, err)
			return err
		}

		//decode response from call to dynamic tcg uri
		err = client.DecodeInterface(resp.Body, cards)
		if err != nil {
			log.Printf("error decoding interface [%v]", err)
			return err
		}

		//append cards to temporary card array
		actualCards = append(actualCards, cards.Cards...)
		links := resp.Header.Get("link")
		if !strings.Contains(links, `rel="next"`) {
			break
		}
		count++
	}

	//set cards.Cards equal to the temp card array
	cards.Cards = actualCards

	return nil
}

func PopulateDatabase(p controller.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		sets, err := GetAllSets(p)
		if err != nil {
			log.Printf("unable to get sets [%v]", err)
		}

		for _, s := range sets.Sets {
			err := p.Mongo.Upsert(s, p.Config.Mongo.SetsCollection)
			if err != nil {
				log.Printf("failed to insert set %v [%v]", s.Code, err)
			}
			cards := GetCardsBySetCode(p, s.Code)
			if len(cards.Cards) != s.TotalCards {
				log.Printf("[WARNING] did not receive all card in set [ %s ] - actual [ %d ] / expected: [ %d ]", s.Name, len(cards.Cards), s.TotalCards)
			}

			for _, c := range cards.Cards {
				err := p.Mongo.Upsert(c, p.Config.Mongo.CardsCollection)
				if err != nil {
					log.Printf("failed to insert card %v [%v]", c.ID, err)
				}
			}
		}
	})
}
