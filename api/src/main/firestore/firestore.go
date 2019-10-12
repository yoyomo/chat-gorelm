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

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

func init() {
	gotenv.Load()
}

var ctx context.Context
var client *firestore.Client

type User struct {
	First  string `firestore:"first,omitempty"`
	Middle string `firestore:"middle,omitempty"`
	Last   string `firestore:"last,omitempty"`
	Born   int    `firestore:"born,omitempty"`
}

func resourceType(resource string) interface{} {

	switch resource {
	case "users":
		return *new(User)
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

		newData := resourceType(resource)

		json.Unmarshal(reqBody, &newData)

		_, _, err = client.Collection(resource).Add(ctx, newData)
		logFatal(err)

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(newData)
	}
}

func get(resource string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resourceID := mux.Vars(r)["id"]

		dsnap, err := client.Collection(resource).Doc(resourceID).Get(ctx)
		logFatal(err)

		data := resourceType(resource)
		dsnap.DataTo(&data)
		dsnap.Data()

		json.NewEncoder(w).Encode(data)
	}
}

func index(resource string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var datas []interface{}

		iter := client.Collection(resource).Documents(ctx)
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}

			data := resourceType(resource)
			doc.DataTo(&data)

			datas = append(datas, data)
		}

		json.NewEncoder(w).Encode(datas)
	}
}

func update(resource string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		reqBody, err := ioutil.ReadAll(r.Body)
		logFatal(err)

		resourceID := mux.Vars(r)["id"]

		updatedData := resourceType(resource)

		json.Unmarshal(reqBody, &updatedData)

		_, err = client.Collection(resource).Doc(resourceID).Set(ctx, updatedData, firestore.MergeAll)
		logFatal(err)

		json.NewEncoder(w).Encode(updatedData)
	}
}

func delete(resource string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resourceID := mux.Vars(r)["id"]
		_, err := client.Collection(resource).Doc(resourceID).Delete(ctx)
		logFatal(err)
	}
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func setupFirestore() {
	ctx = context.Background()

	var firestoreURL string
	var conf *firebase.Config

	switch os.Getenv("GO_ENV") {
	case "production":
		firestoreURL = os.Getenv("FIRESTORE_PROJECT_ID")
		conf = &firebase.Config{ProjectID: firestoreURL}
	case "development":
		firestoreURL = os.Getenv("FIRESTORE_EMULATOR_URL")
		conf = &firebase.Config{ProjectID: "(default)", DatabaseURL: firestoreURL}
	}

	app, err := firebase.NewApp(ctx, conf)
	logFatal(err)

	client, err = app.Firestore(ctx)
	logFatal(err)

	defer client.Close()

	fmt.Println("Firestore:", firestoreURL)

}

func setupRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)

	resources("users", router)

	return router
}

func startServer() {
	router := setupRouter()

	port := os.Getenv("GO_SERVER_PORT")

	fmt.Println("Go Server listening at port:", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func main() {

	setupFirestore()

	startServer()

}
