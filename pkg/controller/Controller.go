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

func (c *Controller)ViewPaginated(w http.ResponseWriter, r *http.Request){

	parms := r.URL.Query()
	log.Println("parameters: ", parms)

	var cards dto.Cards
	var pagi dto.PaginatedCards

	cards.Cards = c.Mongo.GetFilterCards(parms)
	log.Println("cards: ", cards)

	pagi = paginate(cards.Cards)

	log.Printf("Paginated cards: %v", len(pagi.Pages))

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

	e = temp.Execute(w, pagi)
	if e != nil{
		log.Println("error executing template: ", e)
	}
}

func paginate(a []dto.Card) (p dto.PaginatedCards){

	var i int
	var pg dto.Page

	if a != nil{
		numPgs := (len(a) / 6) + 1

		if numPgs == 0{
			numPgs++
		}

		for numPgs != 0{
			for len(pg.Cards) <= 5{
				if i > len(a) - 1 {
					break
				}
				pg.Cards = append(pg.Cards, a[i])
				i++
			}
			p.Pages = append(p.Pages, pg)
			for c := range pg.Cards{
				log.Printf("card: %v", c)
			}
			pg.Cards = nil
			numPgs--
		}
	}
	return
}