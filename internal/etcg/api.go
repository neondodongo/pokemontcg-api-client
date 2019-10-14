package etcg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"pokemontcg-api-client/pkg/config"
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

func (c *PokemonTCGController) GetAllSets() (*dto.Sets, error) {

	var sets dto.Sets

	//fmt.Printf("pokemon tcg api [ %v ]", p.Config.PokemonAPI)
	resp, err := http.Get(c.Config.PokemonAPI + "/sets")
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
func (c *PokemonTCGController) GetCardsBySetCode(setCode string) *dto.Cards {
	setCode = strings.TrimSpace(setCode)
	if setCode != "" {
		var err error
		var cards *dto.Cards
		err, cards = c.GetPaginatedCards(setCode)
		if err != nil {
			log.Printf("error calling pokemon tcg api [%v]", err)
		}

		return cards
	}
	return &dto.Cards{}
}

func (c *PokemonTCGController) GetPaginatedCards(setCode string) (error, *dto.Cards) {
	count := 1
	actualCards := make([]dto.Card, 0)

	for {
		//build query with setCode and replace <count>
		query := c.Config.PokemonAPI + "/cards?page=<count>&setCode=" + setCode
		uri := strings.Replace(query, "<count>", strconv.Itoa(count), 1)
		//call to pokemon tcg api
		resp, err := http.Get(uri)
		if err != nil {
			log.Printf("error in call to uri %v [%v]", uri, err)
			return err, nil
		}

		var cards dto.Cards

		//decode response from call to dynamic tcg uri
		//err = client.DecodeInterface(resp.Body, cards)
		//if err != nil {
		//	log.Fatalf("error decoding interface [%v]", err)
		//	return err, nil
		//}

		b, e := ioutil.ReadAll(resp.Body)
		if e != nil{
			return e, nil
		}

		//log.Println("cards byte slice: ", string(b))

		e = json.Unmarshal(b, &cards)
		if e != nil{
			return e, nil
		}
		log.Println("cards: ", cards)

		//append cards to temporary card array
		actualCards = append(actualCards, cards.Cards...)
		links := resp.Header.Get("link")
		if !strings.Contains(links, `rel="next"`) {
			break
		}
		count++
	}

	//set cards.Cards equal to the temp card array
	cds := &dto.Cards{
		Cards: actualCards,
	}

	return nil, cds
}

func PopulateDatabase(c *PokemonTCGController) {

	fmt.Println("Getting sets")
	sets, err := c.GetAllSets()
	if err != nil {
		log.Printf("unable to get sets [%v]", err)
	}

	for _, s := range sets.Sets {
		err := c.Mongo.Upsert(s, c.Config.Mongo.SetsCollection)
		if err != nil {
			log.Printf("failed to insert set %v [%v]", s.Code, err)
		}
		cards := c.GetCardsBySetCode(s.Code)
		if len(cards.Cards) != s.TotalCards {
			log.Printf("[WARNING] did not receive all card in set [ %s ] - actual [ %d ] / expected: [ %d ]", s.Name, len(cards.Cards), s.TotalCards)
		}

		for _, card := range cards.Cards {
			err := c.Mongo.Upsert(card, c.Config.Mongo.CardsCollection)
			if err != nil {
				log.Printf("failed to insert card %v [%v]", card.ID, err)
			}
		}
	}
}

//func PopulateDatabase
//make call to get all sets
//range over sets
//upsert set information
//get cards paginated list of cards
//range over cards
//upsert each card individually