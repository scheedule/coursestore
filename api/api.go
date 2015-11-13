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
	errorMap        = map[error]int{
		BadRequestError: http.StatusBadRequest,
		DBError:         http.StatusNotFound,
		UnmarshalError:  http.StatusInternalServerError,
	}
)

// Type API contains the database to query and functions we use to query it.
type API struct {
	db *db.DB
}

func New(db *db.DB) *API {
	return &API{db}
}

// This route handles all requests to lookup individual class data. Requests
// will have a department and number and class data will be returned as JSON.
func (a *API) HandleLookup(w http.ResponseWriter, r *http.Request) {
	department := r.FormValue("department")
	number := r.FormValue("number")

	if !isValidDepartment(department) || !isValidCourseNumber(number) {
		log.Debug("query does not contain properly formatted dept/num combination.", BadRequestError)
		handleError(w, BadRequestError)
		return
	}

	log.Debug("querying: ", department, number)

	class, err := a.db.Lookup(department, number)
	if err != nil {
		log.Warn("DB lookup failed: ", err)
		handleError(w, DBError)
		return
	}

	js, err := json.Marshal(class)
	if err != nil {
		log.Error("class unmarshal failed: ", err)
		handleError(w, UnmarshalError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// This route handles requests to get all the class data for every class in one
// request. Data is returned as JSON.
func (a *API) HandleAll(w http.ResponseWriter, r *http.Request) {
	classes, err := a.db.GetAll()
	if err != nil {
		log.Error("failed to query all classes: ", err)
		handleError(w, DBError)
		return
	}

	js, err := json.Marshal(classes)
	if err != nil {
		log.Error("failed to unmarshal all classes: ", err)
		handleError(w, UnmarshalError)
		return
	}

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
