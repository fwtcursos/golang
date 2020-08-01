package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	// PORT to access to the aplication
	PORT = ":8080"
	// DBHost .
	DBHost = "localhost"
	// DBPort .
	DBPort = "3306"
	// DBUser .
	DBUser = "root"
	// DBPass .
	DBPass = "Lescopeta1@"
	// DBDBase .
	DBDBase = "go_building_web_app"
)

// Page represents a db row of pages table
type Page struct {
	Title      string
	RawContent string
	Content    template.HTML
	Date       string
	GUID       string
}

//TruncatedText to show ... if text big
func (p Page) TruncatedText() template.HTML {
	chars := 0
	for i := range p.Content {
		chars++
		if chars > 150 {
			return p.Content[:i] + "..."
		}
	}
	return p.Content
}

var database *sql.DB

func main() {
	ConnDB()

	certificates, err := tls.LoadX509KeyPair("c:/cert/server.crt", "c:/cert/server.key")
	tlsConf := tls.Config{Certificates: []tls.Certificate{certificates}}
	tls.Listen("tcp", ":8080", &tlsConf)
	if err != nil {
		log.Println(err.Error())
	}

	routes := mux.NewRouter()

	routes.HandleFunc("/api/pages", APIPage).
		Methods("GET").
		Schemes("https")
	routes.HandleFunc("/api/page/{guid:[0-9a-zA\\-]+}", APIPage).
		Methods("GET").
		Schemes("https")

	routes.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", ServePage)
	routes.HandleFunc("/", RedirIndex)
	routes.HandleFunc("/home", ServeIndex)
	http.Handle("/", routes)

	log.Println("Server ON -> " + time.Now().String())
	log.Fatal(http.ListenAndServe(PORT, nil))
}

// ServePage .
func ServePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]
	thisPage := Page{}
	err := database.
		QueryRow("SELECT page_title, page_content, page_date FROM pages WHERE page_guid = ?", pageGUID).
		Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	thisPage.Content = template.HTML(thisPage.RawContent)

	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println("Can't get page. GUID:" + pageGUID + " " + err.Error())
		return
	}
	t, err := template.ParseFiles("templates/blog.html")
	if err != nil {
		log.Println("Erro ao carregar template", err.Error())
	}
	t.Execute(w, thisPage)
}

// RedirIndex .
func RedirIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusPermanentRedirect)
}

// ServeIndex .
func ServeIndex(w http.ResponseWriter, r *http.Request) {
	var Pages = []Page{}
	pages, err := database.Query("SELECT page_title, page_content, page_date, page_guid FROM pages ORDER BY ? DESC", "page_date")

	if err != nil {
		fmt.Fprintln(w, err.Error())
	}

	defer pages.Close()
	for pages.Next() {
		thisPage := Page{}
		err := pages.Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date, &thisPage.GUID)
		thisPage.Content = template.HTML(thisPage.RawContent)

		if err != nil {
			log.Fatal("Erro: " + err.Error())
		}

		Pages = append(Pages, thisPage)
	}

	t, err := template.ParseFiles("templates/index.html")

	if err != nil {
		log.Fatal("Erro: " + err.Error())
	}

	t.Execute(w, Pages)
}

// APIPage .
func APIPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]
	thisPage := Page{}
	err := database.
		QueryRow("SELECT page_title, page_content, page_date FROM pages WHERE page_guid = ?", pageGUID).
		Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	thisPage.Content = template.HTML(thisPage.RawContent)

	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println(err.Error())
		return
	}

	// APIOutput, err := json.Marshal(thisPage)
	// fmt.Println(APIOutput)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, thisPage)

}

//ConnDB .
func ConnDB() {
	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUser, DBPass, DBHost, DBPort, DBDBase)
	db, err := sql.Open("mysql", dbConn)

	if err != nil {
		log.Println("Could'n connect", err.Error())
	}

	database = db
}
