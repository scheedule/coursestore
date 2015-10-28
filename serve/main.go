// Package main initializes a connection with the database and starts serving
// course data on the specified port. Environment variables must be set to
// override default settings.
//
// Must set the following environment variables:
//
// DB_HOSTNAME: <mongo> The hostname of your db server.
//
// DB_NAME: <test> The name of the db you wish to connect to.
//
// DB_COLLECTION: <classes> The collection on the database with the class data you need.
//
// SERVE_PORT: <7819> The port on which to serve endpoint.
package main

import (
	"fmt"
	"github.com/scheedule/coursestore/db"
	"net/http"
	"os"
)

// Main db connection to query.
var mydb *db.DB

// Hostname of database we intend to connect to.
var HOSTNAME string = func() string {
	if s := os.Getenv("DB_HOSTNAME"); s != "" {
		return s
	}
	return "mongo"
}()

// Database name we intend to connect to.
var DBNAME string = func() string {
	if s := os.Getenv("DB_NAME"); s != "" {
		return s
	}
	return "test"
}()

// Collection with classes we wish to query.
var COLLECTION string = func() string {
	if s := os.Getenv("DB_COLLECTION"); s != "" {
		return s
	}
	return "classes"
}()

// Port we wish to serve on.
var SERVE_PORT string = func() string {
	if s := os.Getenv("SERVE_PORT"); s != "" {
		return s
	}
	return "7819"
}()

// Initialize database connection with. Will panic if connection fails.
func init() {
	mydb = db.NewDB(HOSTNAME, "27017", DBNAME, COLLECTION)
	err := mydb.Init()
	if err != nil {
		panic(err)
	}
}

// Middleware to print requests out to the console.
func printURI(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		fn(w, r)
	}
}

// Entrypoint. Start listening.
func main() {
	http.HandleFunc("/lookup", printURI(HandleLookup(mydb)))
	http.HandleFunc("/all", printURI(HandleAll(mydb)))
	fmt.Println("Listening")
	http.ListenAndServe(":"+SERVE_PORT, nil)
}
