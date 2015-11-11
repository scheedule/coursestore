// Package API provides all the routes for the webserver to expose. Queries
// are validated processed.
package api

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"

	log "github.com/Sirupsen/logrus"

	"github.com/scheedule/coursestore/db"
)

var (
	BadRequestError = errors.New("The request was malformed.")
	DBError         = errors.New("Query to database failed.")
	UnmarshalError  = errors.New("Error unmarshalling data from the database.")
)

// Type API contains the database to query and functions we use to query it.
type API struct {
	DB *db.DB
}

// Interrogate values and produce JSON.
func (a *API) lookupClass(department, number string) ([]byte, error) {

	if department == "" || number == "" {
		return nil, BadRequestError
	}

	// Ensure department is valid string
	matched, err := regexp.MatchString("^[A-Z]*$", department)
	if !matched || err != nil {
		log.Warn("failed to match department: ", err)
		return nil, BadRequestError
	}

	// Ensure number is numeric
	if _, err := strconv.Atoi(number); err != nil {
		log.Warn("failed to match course number: ", err)
		return nil, BadRequestError
	}

	class, err := a.DB.Lookup(department, number)
	if err != nil {
		log.Warn("DB lookup failed: ", err)
		return nil, DBError
	}

	js, err := json.Marshal(class)
	if err != nil {
		log.Error("class unmarshal failed: ", err)
		return nil, UnmarshalError
	}

	return js, nil
}

// Pack all classes into JSON.
func (a *API) packClasses() ([]byte, error) {
	classes, err := a.DB.GetAll()
	if err != nil {
		log.Error("failed to query all classes: ", err)
		return nil, DBError
	}

	js, err := json.Marshal(classes)
	if err != nil {
		log.Error("failed to unmarshal all classes: ", err)
		return nil, UnmarshalError
	}

	return js, nil
}

// Write the appropriate message to the client.
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
func (a *API) HandleLookup(w http.ResponseWriter, r *http.Request) {
	department := r.FormValue("department")
	number := r.FormValue("number")

	log.Debug("looking up: ", department, number)

	js, err := a.lookupClass(department, number)
	if err != nil {
		log.Debug("lookup resulted in error: ", err)
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// This route handles requests to get all the class data for every class in one
// request. Data is returned as JSON.
func (a *API) HandleAll(w http.ResponseWriter, r *http.Request) {
	js, err := a.packClasses()
	if err != nil {
		log.Debug("handleAll resulted in error: ", err)
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(w)
	defer gz.Close()
	gz.Write(js)
}
