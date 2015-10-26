package db

import (
	"errors"
	"github.com/scheedule/coursestore/types"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	ClassNotFound error = errors.New("Class Not Found")
)

type DB struct {
	session         *mgo.Session
	collection      *mgo.Collection
	server          string
	db_name         string
	collection_name string
}

func NewDB(ip, port, db_name, collection_name string) *DB {
	return &DB{
		server:          ip + ":" + port,
		db_name:         db_name,
		collection_name: collection_name,
	}
}

// Init connection
func (db *DB) Init() error {
	// Initiate DB connection
	session, err := mgo.Dial(db.server)
	if err != nil {
		return err
	}

	db.session = session

	// Establish Session
	db.collection = db.session.DB(db.db_name).C(db.collection_name)

	return nil
}

// Purge DB - Mainly for testing
func (db *DB) Purge() {
	db.collection.DropCollection()
}

// Close the session
func (db *DB) Close() {
	db.session.Close()
}

// Put Class
func (db *DB) Put(entry types.Class) error {
	return db.collection.Insert(entry)
}

// Lookup Class
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

// Get All Class Names
func (db *DB) GetAll() ([]types.Class, error) {
	var result []types.Class
	err := db.collection.Find(nil).All(&result)
	return result, err
}
