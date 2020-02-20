package etcg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"pokemontcg-api-client/pkg/client"
	"pokemontcg-api-client/pkg/config"
	"pokemontcg-api-client/pkg/dto"
	"pokemontcg-api-client/pkg/mongo"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Controller struct {
	Config config.Config
	Mongo  mongo.MongoBongo
	Client http.Client
}

func (c Controller) GetCards() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Access-Control-Allow-Origin", "*") // for local Development only!
		w.Header().Set("Content-Type", "application/json")

		params, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			log.Println("Unable to read query")
		}

		filter := mongo.Filter(params)

		var card dto.Cards
		results, err := c.Mongo.Find(filter, "cards")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			client.RespondWithPrettyJSON(w, 503, fmt.Sprintf("Error finding documents: %v", err))
			return
		}

		b, err := json.Marshal(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			client.RespondWithPrettyJSON(w, 503, fmt.Sprintf("Error finding documents: %v", err))
			return
		}

		// fmt.Println(string(b))
		br := bytes.NewReader(b)

		decoder := json.NewDecoder(br)
		fmt.Println(results)

		if err = decoder.Decode(&card); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			client.RespondWithPrettyJSON(w, 503, fmt.Sprintf("Error finding documents: %v", err))
			return
		}

		// sortCardsByCardNum(cards)
		w.WriteHeader(http.StatusOK)
		client.RespondWithPrettyJSON(w, 200, card)
	})
}

func sortCardsByCardNum(cards []dto.Card) {
	sort.Slice(cards, func(i, j int) bool {

		c1, err := strconv.Atoi(cards[i].Number)
		if err != nil {
			log.Println(err)
		}
		c2, err := strconv.Atoi(cards[j].Number)
		if err != nil {
			log.Println(err)
		}
		return c1 < c2
	})
}

func (c Controller) GetSets() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		_, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			log.Println("Unable to read query")
		}

		var sets []dto.Set
		// sets = c.Mongo.GetFilterSets(params)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		client.RespondWithPrettyJSON(w, 200, sets)
	})
}

//GetAllSets pulls all sets from pokemon tcg api
func (c Controller) GetAllSets() (*dto.Sets, error) {

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
func (c Controller) GetCardsBySetCode(setCode string) *dto.Cards {
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

func (c Controller) GetPaginatedCards(setCode string) (error, *dto.Cards) {
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
		err = client.DecodeBodyToInterface(resp.Body, &cards)
		if err != nil {
			log.Fatalf("error decoding interface [%v]", err)
			return err, nil
		}

		//b, e := ioutil.ReadAll(resp.Body)
		//if e != nil{
		//	return e, nil
		//}
		//
		//e = json.Unmarshal(b, &cards)
		//if e != nil{
		//	return e, nil
		//}
		// log.Println("cards: ", cards)

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

func (c Controller) CreateUser(un, em, pw, pwv string) error {

	if pw != pwv {
		fmt.Println("Password and Password Verify do not match")
		return nil
	}

	// Check if user exists
	// _, err := c.Mongo.FindUserByEmail(em)
	// if err != nil {
	// 	fmt.Printf("Did not find User data [%v]\n", err)
	// } else {
	// 	log.Printf("User already registered with email %s", em)
	// 	return nil
	// }

	// Check if username has been taken
	// _, err = c.Mongo.FindUserByUsername(un)
	// if err != nil {
	// 	fmt.Printf("Failed to find User data [%v]\n", err)
	// } else {
	// 	log.Printf("Username already taken: %s", un)
	// 	return nil
	// }

	// Hash user's password and save to DB
	hash, err := HashPassword(pw)
	if err != nil {
		fmt.Printf("Failed to hash password [%v]\n", err)
		return err
	}

	u := dto.InitUser(un, hash, em)

	r := dto.Role{
		Name:        "AUTH_USER",
		Description: "Basic, authenticated user access to ETCG",
	}

	u.Roles = u.AddRole(r)

	// Save user
	// err = c.Mongo.Upsert(u, bson.M{"id", id}, "user")
	// if err != nil {
	// 	return err
	// }
	return nil
}

func HashPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

//func PopulateDatabase
//make call to get all sets
//range over sets
//upsert set information
//get cards paginated list of cards
//range over cards
//upsert each card individually
