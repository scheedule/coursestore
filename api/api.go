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
	// BadReqeustError returned when request did not arrive in an expected format
	BadRequestError = errors.New("The request was malformed")

	// DBError returned anytime the database fails to retrieve/store/update a record
	DBError = errors.New("Query to database failed")

	// DecodeError returned whenever we fail to decode data from the database or
	// an incoming request.
	DecodeError = errors.New("Error unmarshalling data from the database")

	// Mapping of errors to respective HTTP status codes
	errorMap = map[error]int{
		BadRequestError: http.StatusBadRequest,
		DBError:         http.StatusNotFound,
		DecodeError:     http.StatusInternalServerError,
	}
)

// Type API contains the database to query and functions we use to query it.
type API struct {
	db *db.DB
}

// Construct a new API object with a pointer to a database to query.
func New(db *db.DB) *API {
	return &API{db}
}

// This route handles all requests to lookup individual class data. Requests
// will have a department and number and class data will be returned as JSON.
func (a *API) HandleLookup(w http.ResponseWriter, r *http.Request) {
	department := r.FormValue("department")
	number := r.FormValue("number")

	if !isValidDepartment(department) || !isValidCourseNumber(number) {
		log.Debug("query does not contain properly formatted dept/num combination")
		handleError(w, BadRequestError)
		return
	}

	class, err := a.db.Lookup(department, number)
	if err != nil {
		log.Warn("DB lookup failed: ", err)
		handleError(w, DBError)
		return
	}

	js, err := json.Marshal(class)
	if err != nil {
		log.Error("class marshal failed: ", err)
		handleError(w, DecodeError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// This route handles requests to get all the class data for every class in one
// request. Data is returned as JSON.
func (a *API) HandleAll(w http.ResponseWriter, r *http.Request) {
	var detailLevel = db.DetailBasic

	detail := r.FormValue("detail")
	if detail == "complete" {
		detailLevel = db.DetailComplete
	}

	classes, err := a.db.GetAll(detailLevel)
	if err != nil {
		log.Error("failed to query all classes: ", err)
		handleError(w, DBError)
		return
	}

	js, err := json.Marshal(classes)
	if err != nil {
		log.Error("failed to marshal all classes: ", err)
		handleError(w, DecodeError)
		return
	}

	/*
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	*/

	w.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(w)
	defer gz.Close()
	gz.Write(js)
}

// Write the appropriate message to the client.
func handleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), errorMap[err])
}

// Return true if and only if the department is formatted correctly. This
// function does not check the database for department existence.
func isValidDepartment(department string) bool {
	// Check empty
	if department == "" {
		return false
	}

	// Check capitalized alphabetic.
	if matched, err := regexp.MatchString("^[A-Z]*$", department); !matched || err != nil {
		return false
	}

	return true
}

// Return true if and only if the course number is formatted correctly.
func isValidCourseNumber(number string) bool {
	// Check empty
	if number == "" {
		return false
	}

	// Check numeric
	if _, err := strconv.Atoi(number); err != nil {
		return false
	}

	return true
}
