package etcg

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"text/template"

	"pokemontcg-api-client/pkg/client"
)

const (
	templates = "templates/"
	html_ext  = ".html"
)

// UI application for database administration
func (c *PokemonTCGController) GetDashboard(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	title := path[strings.LastIndex(r.URL.Path, "/")+1:]
	if err := renderTemplate(title, w); err != nil {
		log.Printf("Failed to render template [%v]", err)
		e := json.NewEncoder(w).Encode("404 PAGE NOT FOUND")
		if e != nil{
			log.Println("error encoding value: ", err)
		}
	}

}

// Render a template back using the response
func renderTemplate(title string, w http.ResponseWriter) error {

	p, err := client.LoadPage(title)
	if err != nil {
		return err
	}

	t, err := template.ParseFiles(templates + title + html_ext)
	if err != nil {
		return err
	}
	err = t.Execute(w, p)
	if err != nil {
		return err
	}

	return nil
}
