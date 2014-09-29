package repository

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
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

func NewRepo() Repository {
	return Repository{
		dbmap: initDB(),
	}
}

func NewPost(title string, file []byte) Post {
	return Post{
		Created: time.Now().UnixNano(),
		Title:   title,
		File:    file,
	}
}

func initDB() *gorp.DbMap {

	db_path := os.Getenv("HOME") + "\\.up4go\\"

	os.Mkdir(db_path, 0777)

	db_full_path := db_path + ".postdb.bin"

	fmt.Println(db_full_path)

	db, err := sql.Open("sqlite3", db_full_path)

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

func (repo Repository) Select(id int64) Post {
	var post Post
	repo.dbmap.SelectOne(&post, "select * from posts where post_id=?", id)
	return post
}

func (repo Repository) Delete(id int64) {
	post := repo.Select(id)
	repo.dbmap.Delete(&post)
}

func (repo Repository) Close() {
	repo.dbmap.Db.Close()
}
