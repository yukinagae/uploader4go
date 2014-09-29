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

var list = `<html>
<head>
	<title>UP4go</title>
</head>
<body>
	<table>
		<thead>
			<tr>
				<th>ID</th>
				<th>File Name</th>
			</tr>
		</thead>
		<tbody>
			{{range $index, $element := . }}
			<tr>
				<td><a href="http://127.0.0.1:3000/files/{{ .Id }}">{{ .Id }}</a></td>
				<td>{{ $element.Title }}</td>
			</tr>
			{{end}}
		</tbody>
	</table>

	<form enctype="multipart/form-data" action="http://127.0.0.1:3000/files" method="post">
		<input type="file" name="uploadfile" />
		<input type="submit" value="UPLOAD" />
	</form>
	<body>
</html>
`

var detail = `
<html>
<head>
	<title>UP4go</title>
</head>
<body>
	<table>
		<thead>
			<tr>
				<th>ID</th>
				<th>File Name</th>
			</tr>
		</thead>
		<tbody>
			<tr>
				<td>{{ .Id }}</td>
				<td>{{ .Title }}</td>
			</tr>
		</tbody>
	</table>
	<form action="http://127.0.0.1:3000/download/{{ .Id }}" method="GET">
		<input type="submit" value="DOWNLOAD" />
	</form>
	<form action="http://127.0.0.1:3000/delete/{{ .Id }}" method="POST">
		<input type="submit" value="DELETE" />
	</form>
	<a href="http://127.0.0.1:3000/">To List</a>
	<body>
</html>
`

func main() {

	repo := repository.NewRepo()
	defer repo.Close()

	r := mux.NewRouter()

	// routing
	// GET /
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.New("LIST").Parse(list)
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
		t, _ := template.New("DETAIL").Parse(detail)

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
