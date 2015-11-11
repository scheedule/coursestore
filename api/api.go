// Package API provides all the routes for the webserver to expose. Queries
// are validated processed.
package api

import (
	"errors"

	"github.com/scheedule/coursestore/db"
)

var (
	BadRequestError = errors.New("The request was malformed.")
	DBError         = errors.New("Query to database failed.")
	UnmarshalError  = errors.New("Error unmarshalling data from the database.")
)

// Type API contains the database to query and functions we use to query it.
type API struct {
	db *db.DB
}

func New(db *db.DB) *API {
	return &API{db}
}
