package client

import (
	//"golang.org/x/net/html"
	"html/template"
	//"strings"

	//"html/template"
	//"log"
	"net/http"
	//"os"
	//"path/filepath"
)
//
//func Health() http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//		reader := strings.NewReader("")
//
//		resp, e := ioutil.ReadFile("template/health.html")
//		if e != nil{
//			client.RespondWithPrettyJSON(w, http.StatusNotFound, e)
//		}
//
//
//		e := html.Render(w, )
//
//	})
//}

type Todo struct {
	Name        string
	Description string
}

func Health(w http.ResponseWriter, r *http.Request){
	title := r.URL.Path

	p, err := LoadPage(title)

	t, err := template.ParseFiles("templates/health.html")
	if err != nil {
		panic(err)
	}
	err = t.Execute(w, p)
	if err != nil {
		panic(err)
	}
}

//func push(w http.ResponseWriter, resource string) {
//	pusher, ok := w.(http.Pusher)
//	if ok {
//		if err := pusher.Push(resource, nil); err == nil {
//			return
//		}
//	}
//}

