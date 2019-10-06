package card

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"pokemontcg-api-client/dto"
	"pokemontcg-api-client/pkg/config"
	"pokemontcg-api-client/pkg/mongo"
	"strconv"
	"strings"
)

type CardController struct{
	Cfg config.Config
	Db *mongo.MongoBongo
	Cli *http.Client
}

func GetAllSets(c CardController) (*dto.Sets, error) {

	var sets dto.Sets

	resp, err := http.Get(c.Cfg.PokemonAPI + "/sets")
	if err != nil {
		log.Printf("error getting sets from pokemon-tcg [ %v ]", err)
		return nil, err
	}

	//Read body into byt slice
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading response body [ %v ]", err)
		return nil, err
	}

	//Unmarshal byte slice into sets struct
	if err := json.Unmarshal(body, &sets); err != nil {
		fmt.Println("error unmarshalling into sets")
		return nil, err
	}

	return &sets, nil
}


// GetCardsBySetCode will fetch a list of card by their setCode and return it
func GetCardsBySetCode(c CardController, setCode string) *dto.Cards {

	setCode = strings.TrimSpace(setCode)
	if setCode != "" {

		cards := &dto.Cards{}
		err := getPaginatedCards(c, cards, setCode)
		if err != nil {
			log.Printf("error calling pokemontcg api [ %v ]", err)
		}

		return cards
	}

	return &dto.Cards{}
}

func getPaginatedCards(c CardController, cards *dto.Cards, setCode string) error {
	count := 1
	actualCards := make([]dto.Card, 0)

	for {
		query := c.Cfg.PokemonAPI + "/card?page=<count>&setCode=" + setCode
		uri := strings.Replace(query, "<count>", strconv.Itoa(count), 1)
		fmt.Println("Current URI: " + uri)
		resp, err := c.Cli.Get(uri)
		if err != nil {
			log.Println("initial call")
			return err
		}
		err = decodeInterface(resp.Body, &cards)
		if err != nil {
			log.Println("Error decoding interface", err)
			return err
		}
		actualCards = append(actualCards, cards.Cards...)
		links := resp.Header.Get("link")
		if !strings.Contains(links, `rel="next"`) {
			break
		}
		count++
	}
	cards.Cards = actualCards

	return nil
}

func decodeInterface(io io.Reader, t interface{}) error {
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
