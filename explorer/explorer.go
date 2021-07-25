package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/lelemita/nomadcoin/blockchain"
)

const (
	port string = ":4000"
	templateDir string = "explorer/templates/"
)
var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks []*blockchain.Block
}

func home (rw http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(rw, "Hello from home!!")
	data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home", data)
}

func add (rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		blockData := r.Form.Get("blockData")
		blockchain.GetBlockchain().AddBlock(blockData)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
	
}

func Start() {
		// page 추가하면서 obj 생성
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	// 생성한 templates obj update
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}