package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/net/context"
)

var db *sql.DB;

func ConnectToDB(){

	connStr := "host=localhost user=postgres  sslmode=disable"
	driver := "postgres"

	db, _ = sql.Open(driver, connStr)
	// defer db.Close()

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
}

func dropDB() (result sql.Result, err error) {


	stmt, err := db.Prepare("DROP DATABASE pqgotest;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	result, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Dropped database " + "pqgotest")

	return
}

func main() {

	ConnectToDB()
	dropDB()

	// defer db.Close()


}
