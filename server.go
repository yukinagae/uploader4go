package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
	"text/template"
)

func main() {
	fmt.Println("helo")

	// Routing
	r := mux.NewRouter()

	// for "/"
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("view.html")
		t.Execute(w, &Page{Title: "hoge"})
	}).Methods("GET")

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":3000")
}
