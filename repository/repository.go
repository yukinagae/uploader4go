package main

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"time"
)

type Post struct {
	Id      int64 `db:"post_id"`
	Created int64
	Title   string
	File    []byte
}

type Repository struct {
	dbmap *gorp.DbMap
}

func newRepo() Repository {
	return Repository{
		dbmap: initDB(),
	}
}

func newPost(title string, file []byte) Post {
	return Post{
		Created: time.Now().UnixNano(),
		Title:   title,
		File:    file,
	}
}

func initDB() *gorp.DbMap {
	db, err := sql.Open("sqlite3", "/tmp/post_db.bin")

	checkErr(err, "sql.Open failed")

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")

	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func (repo Repository) Insert(p Post) Post {
	_ = repo.dbmap.Insert(&p)
	return p
}

func (repo Repository) List() []Post {
	var posts []Post
	repo.dbmap.Select(&posts, "select * from posts order by post_id")
	return posts
}

func (repo Repository) Close() {
	repo.dbmap.Db.Close()
}

func main() {

	repo := newRepo()

	defer repo.Close()

	file_name := "hoge.txt"

	// TODO
	file, _ := ioutil.ReadFile(file_name)

	// create two posts
	p1 := newPost(file_name, file)
	repo.Insert(p1)
}
