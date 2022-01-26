package explorer

import (
	"fmt"
	"gotutorial/blockchain"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	templateDir string = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(rw, "Hello from home")
	//tmpl := template.Must(template.ParseFiles("templates/pages/home.gohtml"))
	fmt.Println("asdfad")
	data := homeData{"hello", blockchain.Blocks(blockchain.Blockchain())}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
		break
	case "POST":
		r.ParseForm()
		//data := r.Form.Get("blockData")
		blockchain.Blockchain().AddBlock()
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
		break
	}
}

func Start(port int) {
	router := mux.NewRouter()
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))

	router.HandleFunc("/", home)
	router.HandleFunc("/add", add)
	fmt.Printf("Listening on http://localhost%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

func (h homeData) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("fff")
	fmt.Print(h.PageTitle)
}
