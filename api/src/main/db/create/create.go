package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/net/context"
)

var db *sql.DB;


func createDB() (result sql.Result, err error) {

	conninfo := "user=postgres host=localhost sslmode=disable"
	db, _ = sql.Open("postgres", conninfo)

	stmt, err := db.Prepare("CREATE DATABASE "+ "pqgotest")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	result, err = stmt.Exec()
	fmt.Println("Created database " + "pqgotest")

	return
}

func ConnectToDB(){

	connStr := "host=localhost dbname=pqgotest sslmode=disable"
	driver := "postgres"

	db, _ = sql.Open(driver, connStr)
	// defer db.Close()

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
}

func createTable() (result sql.Result, err error) {

	columns := "id serial PRIMARY KEY, content text, author text "

	stmt, err := db.Prepare("CREATE TABLE " + "posts" + "("+columns+")")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	result, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created table " + "posts")
	return
}

func main() {
	createDB()
	ConnectToDB()
	createTable()

	// defer db.Close()


}
