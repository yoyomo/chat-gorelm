package main

import (
	"fmt"
	"log"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"

	r "main/db/resources"
	"main/db/routes"

	"golang.org/x/net/context"
)

func init() {
	gotenv.Load()
}

var ctx context.Context
var db *routes.DB

func resourceType(resource string) interface{} {

	switch resource {
	case "posts":
		return *new(r.Post)
	}

	return nil
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func resources(resource string, router *mux.Router) {

	router.HandleFunc("/"+resource, create(resource)).Methods("POST")
	router.HandleFunc("/"+resource, index(resource)).Methods("GET")
	router.HandleFunc("/"+resource+"/{id}", get(resource)).Methods("GET")
	router.HandleFunc("/"+resource+"/{id}", update(resource)).Methods("PATCH")
	router.HandleFunc("/"+resource+"/{id}", delete(resource)).Methods("DELETE")

}

func create(resource string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		reqBody, err := ioutil.ReadAll(r.Body)
		logFatal(err)

		// newData := resourceType(resource)
		var newData map[string]interface{}

		json.Unmarshal(reqBody, &newData)

		err = db.Create(resource, newData)
		logFatal(err)

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(newData)
	}
}

func get(resource string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resourceID := mux.Vars(r)["id"]

		data, err := db.GetPost(resourceID)
		logFatal(err)

		json.NewEncoder(w).Encode(data)
	}
}

func index(resource string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		datas, _ := db.Posts(25)

		json.NewEncoder(w).Encode(datas)
	}
}

func update(resource string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		reqBody, err := ioutil.ReadAll(r.Body)
		logFatal(err)

		resourceID := mux.Vars(r)["id"]

		// updatedData := resourceType(resource)
		var updatedData map[string]interface{}

		json.Unmarshal(reqBody, &updatedData)

		err = db.Update(resource, resourceID, updatedData)
		logFatal(err)

		json.NewEncoder(w).Encode(updatedData)
	}
}

func delete(resource string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resourceID := mux.Vars(r)["id"]
		err := db.Delete(resourceID)
		logFatal(err)
	}
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initDB() {

	db = routes.ConnectToDB()

}

func setupRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)

	resources("posts", router)

	return router
}

func startServer() {
	router := setupRouter()

	port := os.Getenv("GO_SERVER_PORT")

	fmt.Println("Go Server listening at port:", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func main() {

	initDB()

	startServer()

}
