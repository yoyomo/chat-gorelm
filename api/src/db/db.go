package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/net/context"
)

type Post struct {
	Id      int
	Content string
	Author  string
}

var DB *sql.DB

// (2)投稿1件の取得
func GetPost(id int) (post Post, err error) {
	post = Post{}
	err = DB.QueryRow("select id,content,author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)
	return
}

// (3)新規投稿の生成
func (post *Post) Create() (err error) {
	statement := "insert into posts (content,author) values ($1,$2) returning id"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
	return
}

func (post *Post) Update() (err error) {
	_, err = DB.Exec("update posts set content = $2,author=$3 where id = $1", post.Id, post.Content, post.Author)
	return
}

func (post *Post) Delete() (err error) {
	_, err = DB.Exec("delete from posts where id=$1", post.Id)
	return
}

func Posts(limit int) (posts []Post, err error) {
	rows, err := DB.Query("select id, content,author from posts limit $1", limit)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Content, &post.Author)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

func createDB() {

	connStr := "user=pqgotest dbname=pqgotest sslmode=disable"
	driver := "postgres"

	DB, _ = sql.Open(driver, connStr)
	defer DB.Close()

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	if err := DB.PingContext(ctx); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
}

func main() {

	createDB()

}
