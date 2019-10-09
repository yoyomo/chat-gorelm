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

var db *sql.DB

func createDB() {

	connStr := "user=pqgotest dbname=pqgotest sslmode=disable"
	driver := "postgres"

	db, _ = sql.Open(driver, connStr)
	defer db.Close()

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
}

func main() {

	createDB()

}
