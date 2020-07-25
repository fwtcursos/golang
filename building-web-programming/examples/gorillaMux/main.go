package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	// PORT to access to the aplication
	PORT = ":8080"
	// DBHost .
	DBHost = "localhost"
	// DBPort .
	DBPort = "3307"
	// DBUser .
	DBUser = "root"
	// DBPass .
	DBPass = "Lescopeta1@"
	// DBDBase .
	DBDBase = "go_building_web_app"
)

// Page represents a db row of pages table
type Page struct {
	Title   string
	Content string
	Date    string
}

var database *sql.DB

func main() {
	connDB()

	router := mux.NewRouter()
	router.HandleFunc("/page/{id:[0-9]+}", staticPageHandler)
	router.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", dynamicPageHandler)
	router.HandleFunc("/notfound", error404PageHandler)
	http.Handle("/", router)
	http.ListenAndServe(PORT, nil)
}

func staticPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageID := vars["id"]
	fileName := "files/" + pageID + ".html"
	fmt.Println(pageID)
	_, err := os.Stat(fileName)

	if err != nil {
		http.Redirect(w, r, "/notfound", http.StatusMovedPermanently)
		return
	}

	http.ServeFile(w, r, fileName)
}

func dynamicPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]
	thisPage := Page{}
	fmt.Println(pageGUID)
	err := database.
		QueryRow("SELECT page_title, page_content, page_date FROM pages WHERE page_guid = ?", pageGUID).
		Scan(&thisPage.Title, &thisPage.Content, &thisPage.Date)

	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println("Can't get page. GUID:" + pageGUID + " " + err.Error())
	}
	html := getPage(thisPage)
	fmt.Fprintln(w, html)
}

func error404PageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Error")
	fileName := "files/404.html"
	http.ServeFile(w, r, fileName)
}

func getPage(thisPage Page) string {
	return `<html><head><title>` + thisPage.Title +
		`</title></head><body><h1>` + thisPage.Title + `</h1><div>` +
		thisPage.Content + `</div></body></html>`
}

func connDB() {
	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUser, DBPass, DBHost, DBPort, DBDBase)
	fmt.Println("connection string: " + dbConn)
	db, err := sql.Open("mysql", dbConn)

	if err != nil {
		log.Println("Could'n connect")
		log.Println(err.Error())
	}

	database = db
}
