// Package api provides all the routes for the webserver to expose. Queries
// are validated processed.
package api

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/scheedule/coursestore/db"
	"net/http"
	"regexp"
	"strconv"
)

var (
	BadRequestError = errors.New("The request was malformed")
	DBError         = errors.New("Query to database failed.")
	UnmarshalError  = errors.New("Error unmarshalling data from the database")
)

// Type Api contains the database to query and functions we use to query it
type Api struct {
	Mydb *db.DB
}

// Interrogate values and produce JSON
func lookupClass(db *db.DB, department, number string) ([]byte, error) {

	if department == "" || number == "" {
		return nil, BadRequestError
	}

	matched, err := regexp.MatchString("^[A-Z]*$", department)
	if !matched || err != nil {
		log.Warn("Failed to match department: ", err)
		return nil, BadRequestError
	}

	if _, err := strconv.Atoi(number); err != nil {
		log.Warn("Failed to match course number: ", err)
		return nil, BadRequestError
	}

	class, err := db.Lookup(department, number)
	if err != nil {
		log.Warn("DB Lookup Failed: ", err)
		return nil, DBError
	}

	js, err := json.Marshal(class)
	if err != nil {
		log.Error("Class Unmarshal Failed: ", err)
		return nil, UnmarshalError
	}

	return js, nil
}

// Pack all classes into JSON
func packClasses(db *db.DB) ([]byte, error) {
	classes, err := db.GetAll()
	if err != nil {
		log.Error("Failed to query all classes: ", err)
		return nil, DBError
	}

	js, err := json.Marshal(classes)
	if err != nil {
		log.Error("Failed to Unmarshal all classes: ", err)
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
func (a *Api) HandleLookup(w http.ResponseWriter, r *http.Request) {
	department := r.FormValue("department")
	number := r.FormValue("number")

	log.Debug("Looking up: ", department, number)

	js, err := lookupClass(a.Mydb, department, number)
	if err != nil {
		log.Debug("Lookup resulted in error:", err)
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// This route handles requests to get all the class data for every class in one
// request. Data is returned as JSON.
func (a *Api) HandleAll(w http.ResponseWriter, r *http.Request) {
	js, err := packClasses(a.Mydb)
	if err != nil {
		log.Debug("HandleAll resulted in error:", err)
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(w)
	defer gz.Close()
	gz.Write(js)
}
