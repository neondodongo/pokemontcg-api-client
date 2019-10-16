package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"pokemontcg-api-client/pkg/client"
	"pokemontcg-api-client/pkg/config"
	"pokemontcg-api-client/pkg/dto"
	"pokemontcg-api-client/pkg/mongo"
)

type Controller struct {
	Config config.Config
	Mongo  mongo.MongoBongo
	Client client.WebClient
}

// {{/* a comment */}}	Defines a comment
// {{.}}	Renders the root element
// {{.Title}}	Renders the “Title”-field in a nested element
// {{if .Done}} {{else}} {{end}}	Defines an if-Statement
// {{range .Todos}} {{.}} {{end}}	Loops over all “Todos” and renders each using {{.}}
// {{block "content" .}} {{end}}	Defines a block with the name “content”

func (c *Controller)ViewCard(w http.ResponseWriter, r *http.Request){

	gon := r.URL.Query()
	log.Println("gone: ", gon)

	var cards dto.Cards

	cards.Cards = c.Mongo.GetFilterCards(gon)
	log.Println("cards: ", cards)

	p, e := client.LoadPage("view")
	if e != nil{
		_, e:= fmt.Fprintf(w, fmt.Sprintf("error loading page: %v", e))
		if e != nil{
			fmt.Println("error serving: ", e)
		}
	}

	temp, e := template.ParseFiles("templates/" + p.Title + ".html")
	if e != nil{
		log.Println("error parsing template: ", e)
	}

	e = temp.Execute(w, cards)
	if e != nil{
		log.Println("error executing template: ", e)
	}
}