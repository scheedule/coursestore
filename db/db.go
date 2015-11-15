// Package db handles all course storing and retrieval from the database.
// This package provides an abstraction to allow users to interact with the
// database with the Class struct type and restrict usage to looking up,
// putting, and purging.
package db

import (
	"errors"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/scheedule/coursestore/types"
)

var (
	// ClassNotFound is returned when a class can't be resolved.
	ClassNotFound error = errors.New("Class Not Found")

	// InternalError is returned when we fail to communicate with the
	// database without error
	InternalError error = errors.New("Internal Database Error")

	// Detail levels
	DetailBasic = bson.M{
		"department":              "1",
		"course_number":           "1",
		"name":                    "1",
		"sections":                "1",
		"sections.crn":            "1",
		"sections.code":           "1",
		"sections.meetings":       "1",
		"sections.meetings.type":  "1",
		"sections.meetings.start": "1",
		"sections.meetings.end":   "1",
		"sections.meetings.days":  "1",
	}
)

// Main primitive to hold db connection and attributes. Users will obtain
// and make requests with the DB type.
type DB struct {
	session        *mgo.Session
	collection     *mgo.Collection
	server         string
	dbName         string
	collectionName string
}

// Construct a new DB type
func New(ip, port, dbName, collectionName string) *DB {
	return &DB{
		server:         ip + ":" + port,
		dbName:         dbName,
		collectionName: collectionName,
	}
}

// Initialize connection to database. An error will be returned if a database
// can't be connected to within a minute.
func (db *DB) Init() error {
	// Initiate DB connection
	session, err := mgo.DialWithTimeout(db.server, 5*time.Second)
	if err != nil {
		log.Error("failed to dial database")
		return InternalError
	}

	db.session = session

	// Establish Session
	db.collection = db.session.DB(db.dbName).C(db.collectionName)
	if db.collection == nil {
		log.Error("failed to establish session and connect to collection")
		return InternalError
	}

	return nil
}

// Drop the specified collection from the database.
func (db *DB) Purge() error {
	err := db.collection.DropCollection()
	if err != nil {
		log.Error("failed to purge database")
		return InternalError
	}

	return nil
}

// Close the session with the database.
func (db *DB) Close() {
	db.session.Close()
}

// Put Class into the database.
func (db *DB) Put(entry types.Class) error {
	err := db.collection.Insert(entry)
	if err != nil {
		log.Error("failed to insert class")
		return InternalError
	}

	return nil
}

// Lookup Class in the database.
func (db *DB) Lookup(department, number string) (*types.Class, error) {
	temp := &types.Class{}
	err := db.collection.Find(bson.M{
		"department":    department,
		"course_number": number,
	}).One(temp)

	if err != nil {
		log.Warn("failed to find class in database")
		return nil, ClassNotFound
	}

	return temp, nil
}

// Get All Class Names from the database.
func (db *DB) GetAll(detail map[string]interface{}) ([]types.Class, error) {
	var result []types.Class
	err := db.collection.Find(nil).Select(detail).All(&result)
	if err != nil {
		log.Error("failed to collect all entries in the collection")
		return nil, InternalError
	}
	return result, nil
}
