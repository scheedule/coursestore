// Package main  connects to the database and starts scraping the course api
// at the specified TERM_URL and populates the database with class data.
// Set environment variables to connect to the appropriate server, db,
// collection, and process the correct term.
//
// Must set the following environment variables:
//
// DB_HOSTNAME: <mongo> The hostname of your db server.
//
// DB_NAME: <test> The name of the db you wish to connect to.
//
// DB_COLLECTION: <classes> The collection on the database with the class data you need.
//
// TERM_URL: The url with the term xml of interest.
package main

import (
	"github.com/scheedule/coursestore/db"
	"github.com/scheedule/coursestore/types"
	"os"
)

// Hostname of db we intend to connect to.
var HOSTNAME string = func() string {
	if s := os.Getenv("DB_HOSTNAME"); s != "" {
		return s
	}
	return "mongo"
}()

// DB name we intend to connect to.
var DBNAME string = func() string {
	if s := os.Getenv("DB_NAME"); s != "" {
		return s
	}
	return "test"
}()

// DB collection with the class data we need.
var COLLECTION string = func() string {
	if s := os.Getenv("DB_COLLECTION"); s != "" {
		return s
	}
	return "classes"
}()

// URL of term we wish to scrape from.
var TERM_URL string = func() string {
	if s := os.Getenv("TERM_URL"); s != "" {
		return s
	}
	return "http://courses.illinois.edu/cisapp/explorer/schedule/2016/spring.xml"
}()

func main() {
	err := PopulateDB(
		TERM_URL,
		HOSTNAME,
		"27017",
		DBNAME,
		COLLECTION,
	)

	if err != nil {
		panic(err)
	}
}

// Populate the selected database with the data scraped from the TERM_URL.
func PopulateDB(term_url, ip, port, db_name, collection_name string) error {
	mydb := db.NewDB(ip, port, db_name, collection_name)

	err := mydb.Init()
	if err != nil {
		return err
	}

	mydb.Purge()

	term, err := GetXML(term_url)

	course_chan := make(chan types.Class)

	go DigestAll(term, course_chan)

	for class := range course_chan {
		err = mydb.Put(class)
		if err != nil {
			return err
		}
	}

	return nil
}
