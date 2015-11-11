package api

import (
	"os"

	"github.com/scheedule/coursestore/db"
)

var testAPI *API

func init() {
	testDB := db.New(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_COLLECTION"))
	testDB.Init()
	testAPI = New(testDB)
}
