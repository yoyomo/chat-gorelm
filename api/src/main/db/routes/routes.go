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
func (db *DB) GetPost(id string) (post r.Post, err error) {
	post = r.Post{}
	err = db.QueryRow("select id,content,author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)
	return
}

// (3)新規投稿の生成
func (db *DB) Create(resourceType string, data map[string]interface{}) (err error) {

	statement := "insert into posts (content,author) values ($1,$2) returning id"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	var id string

	err = stmt.QueryRow(data["Content"], data["Author"]).Scan(&id)
	data["Id"] = id
	return
}

func (db *DB)  Update(resourceType string, id string, data map[string]interface{}) (err error) {
	_, err = db.Exec("update posts set content = $2,author=$3 where id = $1", id, data["Content"], data["Author"])
	return
}

func (db *DB) Delete(id string) (err error) {
	_, err = db.Exec("delete from posts where id=$1", id)
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
