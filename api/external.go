package api

import (
	"compress/gzip"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

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
