package routes

import (
	"log"
	"database/sql"

	r "main/db/resources"


	_ "github.com/lib/pq"
	"golang.org/x/net/context"
)

type DB struct {
	*sql.DB
}

func Open(driver string, cnnStr string) (*DB, error) {
	sqlDB, err := sql.Open(driver, cnnStr)
	return &DB{sqlDB}, err
}

func ConnectToDB() (db *DB){

	connStr := "host=localhost dbname=pqgotest sslmode=disable"
	driver := "postgres"

	db, _ = Open(driver, connStr)
	// defer db.Close()

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	return
}

// (2)投稿1件の取得
func (db *DB) GetPost(id int) (post r.Post, err error) {
	post = r.Post{}
	err = db.QueryRow("select id,content,author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)
	return
}

// (3)新規投稿の生成
func (db *DB) Create(post *r.Post) (err error) {
	statement := "insert into posts (content,author) values ($1,$2) returning id"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
	return
}

func (db *DB)  Update(post *r.Post) (err error) {
	_, err = db.Exec("update posts set content = $2,author=$3 where id = $1", post.Id, post.Content, post.Author)
	return
}

func (db *DB) Delete(post *r.Post) (err error) {
	_, err = db.Exec("delete from posts where id=$1", post.Id)
	return
}

func (db *DB) Posts(limit int) (posts []r.Post, err error) {
	rows, err := db.Query("select id, content,author from posts limit $1", limit)
	if err != nil {
		return
	}
	for rows.Next() {
		post := r.Post{}
		err = rows.Scan(&post.Id, &post.Content, &post.Author)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}
