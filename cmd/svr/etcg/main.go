package main

import (
	"fmt"
	"log"
	"net/http"
	"pokemontcg-api-client/pkg/card"
	"pokemontcg-api-client/pkg/client"
	"pokemontcg-api-client/pkg/config"
	"pokemontcg-api-client/pkg/controller"
	"pokemontcg-api-client/pkg/mongo"

	"strings"

	"github.com/gorilla/mux"
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
	log.Println("ELITE TRAINER TCG")
	log.Println("Server is running on port: ", cfg.Port)

	////Mux Router handling
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", Home)
	r.HandleFunc("/receive", receiveAjax)
	r.HandleFunc("/test-page", testPage)
	r.HandleFunc("/view", con.ViewCard)
	r.HandleFunc("/test-load-health", testLoadHealth)
	r.HandleFunc("/sayhelloname", sayhelloName)
	r.Handle("/card/{cardId}", card.GetCards(con)).Methods(http.MethodGet)

	//run server on port
	log.Fatal(http.ListenAndServe(":3000", http.Handler(r)))

}
func testPage(w http.ResponseWriter, r *http.Request) {
	p := client.Page{
		Title: "test-page",
		Body:  []byte("This is a test of Pages"),
	}
	e := p.Save()
	if e != nil {
		n, e := fmt.Fprintf(w, fmt.Sprintf("error saving test-page: %v", e))
		if e != nil {
			fmt.Println("n: ", n, "e: ", e)
		}
	}
}

func testLoadHealth(w http.ResponseWriter, r *http.Request) {
	p, e := client.LoadPage("health")
	if e != nil {
		_, e := fmt.Fprintf(w, fmt.Sprintf("error loading page: %v", e))
		if e != nil {
			fmt.Println("error serving: ", e)
		}
	}
	_, e = fmt.Fprintf(w, "%s", p.Body)
	if e != nil {
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
	if e != nil {
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
	if e != nil {
		fmt.Println("error sending data to client: ", e)
		return
	}
}


func Home(w http.ResponseWriter, r *http.Request) {
	log.Printf("request: %v", r)
	html := `<head>	
<script src='//ajax.googleapis.com/ajax/libs/jquery/1.11.2/jquery.min.js'></script>
              </head>    
                  <html><body>
                  <h1>Golang Jquery AJAX example</h1>

                  <div id='result'><h3>before</h3></div><br><br>

                  <input id='ajax_btn' type='button' value='POST via AJAX to Golang server'>
                  </body></html>

                   <script>
                   $(document).ready(function () { 
                         $('#ajax_btn').click(function () {
                             $.ajax({
                               url: 'receive',
                               type: 'post',
                               dataType: 'html',
                               data : { ajax_post_data: 'hello'},
                               success : function(data) {
                                 alert('ajax data posted');
                                 $('#result').html(data);
                               },
                             });
                          });
                    });
                    </script>`

	n, e := w.Write([]byte(fmt.Sprintf(html)))
	if e != nil{
		client.RespondWithPrettyJSON(w, n, e)
		return
	}

}

func receiveAjax(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ajaxPostData := r.FormValue("ajax_post_data")
		fmt.Println("Receive ajax post data string ", ajaxPostData)
		n, e := w.Write([]byte("<h2>after<h2>"))
		if e != nil{
			client.RespondWithPrettyJSON(w, n, e)
			return
		}
	}
}