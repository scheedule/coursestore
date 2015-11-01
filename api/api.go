// Package api provides all the routes for the webserver to expose. Queries
// are validated processed.
//
// Must set the following environment variables:
//
// DB_HOSTNAME: <mongo> The hostname of your db server.
//
// DB_NAME: <test> The name of the db you wish to connect to.
//
// DB_COLLECTION: <classes> The collection on the database with the class data you need.
package api

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"github.com/scheedule/coursestore/db"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

var (
	BadRequestError = errors.New("The request was malformed")
	DBError         = errors.New("Failed to query the database")
	UnmarshalError  = errors.New("Error unmarshalling data from the database")
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

// Initialize database connection. Will panic if connection fails.
func init() {
	mydb = db.NewDB(HOSTNAME, "27017", DBNAME, COLLECTION)
	err := mydb.Init()
	if err != nil {
		panic(err)
	}
}

// Interrogate values and produce JSON
func lookupClass(db *db.DB, department, number string) ([]byte, error) {

	if department == "" || number == "" {
		return nil, BadRequestError
	}

	matched, err := regexp.MatchString("^[A-Z]*$", department)
	if !matched || err != nil {
		return nil, BadRequestError
	}

	if _, err := strconv.Atoi(number); err != nil {
		return nil, BadRequestError
	}

	class, err := db.Lookup(department, number)
	if err != nil {
		return nil, DBError
	}

	js, err := json.Marshal(class)
	if err != nil {
		return nil, UnmarshalError
	}

	return js, nil
}

// Pack all classes into JSON
func packClasses(db *db.DB) ([]byte, error) {
	classes, err := db.GetAll()
	if err != nil {
		return nil, DBError
	}

	js, err := json.Marshal(classes)
	if err != nil {
		return nil, UnmarshalError
	}

	return js, nil
}

// Write the appropriate message to the client
func handleError(w http.ResponseWriter, err error) {
	switch err {
	case BadRequestError:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case DBError:
		http.Error(w, err.Error(), http.StatusNotFound)
	case UnmarshalError:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// This route handles all requests to lookup individual class data. Requests
// will have a department and number and class data will be returned as JSON.
func HandleLookup(w http.ResponseWriter, r *http.Request) {
	department := r.FormValue("department")
	number := r.FormValue("number")

	js, err := lookupClass(mydb, department, number)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// This route handles requests to get all the class data for every class in one
// request. Data is returned as JSON.
func HandleAll(w http.ResponseWriter, r *http.Request) {
	js, err := packClasses(mydb)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(w)
	defer gz.Close()
	gz.Write(js)
}