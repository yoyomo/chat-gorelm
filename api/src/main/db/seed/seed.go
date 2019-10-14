package main

import (
	"fmt"
	"log"

	r "main/db/resources"
	"main/db/routes"

	_ "github.com/lib/pq"
)

func main() {

	db := routes.ConnectToDB()

	post := r.Post{
		Content: "Hello",
		Author:  "Taro",
	}

	fmt.Println(post)
	err := db.Create(&post)
	if err != nil {
		log.Fatal(err);
	}
	readPost, err := db.GetPost(post.Id)
	if err != nil {
		log.Fatal(err);
	}
	fmt.Println(readPost)
	readPost.Author = "Poerre"
	db.Update(&readPost)

	fmt.Println(db.GetPost(post.Id))


	fmt.Println(db.Posts(5))
}
