// Package db handles all course storing and retrieval from the database.
// This package provides an abstraction to allow users to interact with the
// database with the Class struct type and restrict usage to looking up,
// putting, and purging.
package db

import (
	"errors"
	"github.com/scheedule/coursestore/types"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	// ClassNotFound is returned when a class can't be resolved.
	ClassNotFound error = errors.New("Class Not Found")
)

// Main primitive to hold db connection and attributes. Users will obtain
// and make requests with the DB type.
type DB struct {
	session         *mgo.Session
	collection      *mgo.Collection
	server          string
	db_name         string
	collection_name string
}

// Construct a new DB type
func NewDB(ip, port, db_name, collection_name string) *DB {
	return &DB{
		server:          ip + ":" + port,
		db_name:         db_name,
		collection_name: collection_name,
	}
}

// Initialize connection to database. An error will be returned if a database
// can't be connected to within a minute.
func (db *DB) Init() error {
	// Initiate DB connection
	session, err := mgo.DialWithTimeout(db.server, 5*time.Second)
	if err != nil {
		return err
	}

	db.session = session

	// Establish Session
	db.collection = db.session.DB(db.db_name).C(db.collection_name)

	return nil
}

// Drop the specified collection from the database.
func (db *DB) Purge() {
	db.collection.DropCollection()
}

// Close the session with the database.
func (db *DB) Close() {
	db.session.Close()
}

// Put Class into the database.
func (db *DB) Put(entry types.Class) error {
	return db.collection.Insert(entry)
}

// Lookup Class in the database.
func (db *DB) Lookup(department, number string) (*types.Class, error) {
	temp := &types.Class{}
	err := db.collection.Find(bson.M{
		"department":    department,
		"course_number": number,
	}).One(temp)

	if err != nil {
		return nil, ClassNotFound
	}

	return temp, nil
}

// Get All Class Names from the database.
func (db *DB) GetAll() ([]types.Class, error) {
	var result []types.Class
	err := db.collection.Find(nil).All(&result)
	return result, err
}
