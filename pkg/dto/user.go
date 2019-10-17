package dto

import (
	"log"
	"strings"
	"time"
)

type User struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
	Title    string `json:"title" bson:"title"`
	// Rank       string    `json:"rank" bson:"rank"`
	Decks      []Deck    `json:"decks bson:"decks"`
	Verified   bool      `json:"verified" bson:"verified"`
	CreationDt time.Time `json:"creationdt" bson:"creationdt"`
	LastActive time.Time `json:"lastactive" bson:"lastactive"`
	DeckLimit  uint8     `json:"decklimit" bson:"decklimit"`
	Roles      []Role    `json:"roles" bson:"roles"`
}

type Deck struct {
	Trainers []Card `json:"trainers" bson:"trainers"`
	Pokemon  []Card `json:"pokemon" bson:"pokemon"`
	Energy   []Card `json:"energy" bson:"energy"`
}

type Role struct {
	Name        string
	Description string
}

// InitUser creates a new user with only username, email, and password fields
func InitUser(un, pw, em string) (user User) {
	user.Username = un
	user.Email = em
	user.Password = pw
	user.Title = "Newbie"
	user.CreationDt = time.Now()
	user.Decks = []Deck{}
	return
}

func (u User) AddRole(r Role) (roles []Role) {

	if strings.TrimSpace(r.Name) == "" {
		log.Println("Cannot add a Role with no name")
		return nil
	}

	return append(roles, r)

}
