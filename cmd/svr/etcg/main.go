package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"pokemontcg-api-client/pkg/client"
	"pokemontcg-api-client/pkg/config"
	"pokemontcg-api-client/pkg/controller"
	"pokemontcg-api-client/pkg/mongo"
	"pokemontcg-api-client/pkg/tcg"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg, err := config.GetConfig("cmd/svr/etcg/config.json")
	if err != nil {
		log.Fatalf("error getting configuration [ %v ]", err)
	}

	db := mongo.InitDatabase(cfg)
	cli := client.InitializeClient(cfg)

	//InitiateControllers
	con := controller.Controller{
		Config: cfg,
		Mongo:  db,
		Client: cli,
	}
	//
	log.Println("log: Elite Trainer TCG [ running ]")
	log.Println("Server is running on port: ", cfg.Port)

	////Mux Router handling
	r := mux.NewRouter().StrictSlash(true)
	//
	////Handler functions for endpoints
	//r.Handle("/health", client.Health()).Methods(http.MethodGet)
	r.HandleFunc("/test-page", testPage)
	r.HandleFunc("/test-load-health", testLoadHealth)
	r.HandleFunc("/sayhelloname", sayhelloName)
	r.Handle("/populate-database", tcg.PopulateDatabase(con)).Methods(http.MethodGet)
<<<<<<< HEAD
	r.Handle("/cards", card.GetCards(con)).Methods(http.MethodGet)

	//run server on port
=======
	////r.Handle("/card/{cardId}", card.GetCard(con)).Methods(http.MethodGet)
	//
	////run server on port
>>>>>>> a1a505e731434ea748203b293421a132bc444d75
	log.Fatal(http.ListenAndServe(":3000", http.Handler(r)))

	//p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	//p1.save()
	//p2, _ := loadPage("TestPage")
	//fmt.Println(string(p2.Body))

	//http.HandleFunc("/health", client.Health)
	//log.Fatal(http.ListenAndServe(":8080", nil))
	//fs := http.FileServer(http.Dir("static"))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))
	//
	//http.HandleFunc("/health", client.Health)
	//
	//log.Println("Listening...")
	//http.ListenAndServe(":3000", nil)
}
func testPage(w http.ResponseWriter, r *http.Request){
	p := client.Page{
		Title: "test-page",
		Body:  []byte("This is a test of Pages"),
	}
	e := p.Save()
	if e != nil{
		n, e := fmt.Fprintf(w, fmt.Sprintf("error saving test-page: %v", e))
		if e != nil{
			fmt.Println("n: ", n, "e: ", e)
		}
	}
}

func testLoadHealth(w http.ResponseWriter, r *http.Request){
	p, e := client.LoadPage("health")
	if e != nil{
		_, e:= fmt.Fprintf(w, fmt.Sprintf("error loading page: %v", e))
		if e != nil{
			fmt.Println("error serving: ", e)
		}
	}
	_, e = fmt.Fprintf(w,"%s", p.Body)
	if e != nil{
		fmt.Println("error serving: ", e)
	}

}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	p, _ := r.URL.User.Password()
	fmt.Println("header  : ", r.Header)
	fmt.Println("user    : ", r.URL.User.String())
	fmt.Println("username: ", r.URL.User.Username())
	fmt.Println("password: ", p)
	fmt.Println("forceQry: ", r.URL.ForceQuery)
	fmt.Println("fragment: ", r.URL.Fragment)
	fmt.Println("host    : ", r.URL.Host)
	fmt.Println("opaque  : ", r.URL.Opaque)
	fmt.Println("rawpath : ", r.URL.RawPath)
	fmt.Println("rawquery: ", r.URL.RawQuery)
	fmt.Println("hostname: ", r.URL.Hostname())

	e := r.ParseForm() // parse arguments, you have to call this by yourself
	if e != nil{
		fmt.Println("error parsing form: ", e)
		return
	}
	fmt.Println("form    : ", r.Form) // print form information in server side
	fmt.Println("path    : ", r.URL.Path)
	fmt.Println("scheme  : ", r.URL.Scheme)
	fmt.Println("url-long:", r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key     :", k)
		fmt.Println("val     :", strings.Join(v, ""))
	}
	_, e = fmt.Fprintf(w, "Hello astaxie!") // send data to client side
	if e != nil{
		fmt.Println("error sending data to client: ", e)
		return
	}
}