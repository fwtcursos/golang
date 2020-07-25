package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	// Port to access de application
	Port = ":8080"
)

func main() {
	http.HandleFunc("/dynamic", serveDynamic)
	http.HandleFunc("/static", serveStatic)
	http.HandleFunc("/", serveHome)

	log.Println("Server ON")
	log.Fatal(http.ListenAndServe(Port, nil))
	// log.Fatal(http.ListenAndServe(Port, http.FileServer(http.Dir("")))) //exibe o diretório para navegar se deixar com ""
}

func serveDynamic(w http.ResponseWriter, r *http.Request) {
	response := "OK - " + time.Now().String()
	fmt.Fprintln(w, response)
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "teste.html")
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	response := r.Method + " " + r.Proto + " " + r.URL.String() + " " + r.UserAgent()
	w.Write([]byte(response))
}
