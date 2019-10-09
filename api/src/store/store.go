package main

import (
	"database/sql"
	"fmt"
	"log"

	"db"

	_ "github.com/lib/pq"
	"golang.org/x/net/context"
)

var DB *sql.DB

func connectToDB() {

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

	connectToDB()

	post := db.Post{
		Content: "Hello",
		Author:  "Taro",
	}

	fmt.Println(post)
	post.Create()

	fmt.Println(post)

	readPost, _ := db.GetPost(post.Id)
	readPost.Author = "Poerre"
	readPost.Update()

}
