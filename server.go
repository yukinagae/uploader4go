package main

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/yukinagae/uploader4go/repository"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"text/template"
)

func main() {

	repo := repository.NewRepo()
	defer repo.Close()

	r := mux.NewRouter()

	// routing
	// GET /
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("list.html")
		posts := repo.List()
		t.Execute(w, posts)
	}).Methods("GET")
	// POST / files
	r.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		file, header, _ := r.FormFile("uploadfile")
		defer file.Close()

		file_name := header.Filename
		content, _ := ioutil.ReadAll(file)
		post := repo.Insert(repository.NewPost(file_name, content))

		http.Redirect(w, r, fmt.Sprintf("/files/%d", post.Id), 301)
	}).Methods("POST")
	// GET /files/:id
	r.HandleFunc("/files/{id}", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("detail.html")

		vars := mux.Vars(r)
		fmt.Println(vars)
		id := vars["id"]
		pk, _ := strconv.ParseInt(id, 10, 64)
		post := repo.Select(pk)

		t.Execute(w, post)
	}).Methods("GET")
	// GET /download/:id
	r.HandleFunc("/download/{id}", func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		fmt.Println(vars)
		id := vars["id"]
		pk, _ := strconv.ParseInt(id, 10, 64)
		post := repo.Select(pk)

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", post.Title))
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

		io.Copy(w, bytes.NewReader(post.File))

	}).Methods("GET")
	// DELETE /files/:id
	r.HandleFunc("/delete/{id}", func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		fmt.Println(vars)
		id := vars["id"]
		pk, _ := strconv.ParseInt(id, 10, 64)
		repo.Delete(pk)

		http.Redirect(w, r, "/", 301)
	}).Methods("POST")

	// start
	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":3000")
}
