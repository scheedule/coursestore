package db

import (
	"github.com/scheedule/coursestore/types"
	"testing"
)

func TestNewDB(t *testing.T) {
	mydb := NewDB("mongo", "27017", "test", "test")
	if mydb == nil {
		t.Fail()
	}
}

func TestInit(t *testing.T) {
	mydb := NewDB("mongo", "27017", "test", "test")
	err := mydb.Init()
	if err != nil {
		t.Error("Failed to initialize DB:", err)
	}
	if mydb.session == nil {
		t.Error("Session is nil")
	}
	if mydb.collection == nil {
		t.Error("Collection is nil")
	}
}

func getDB() *DB {
	mydb := NewDB("mongo", "27017", "test", "test")
	_ = mydb.Init()
	return mydb
}

var sampleClass = types.Class{
	Department:   "CS",
	CourseNumber: "125",
}

func TestPurge(t *testing.T) {
	mydb := getDB()
	mydb.Put(sampleClass)
	mydb.Purge()
	classes, err := mydb.GetAll()
	if err != nil {
		t.Error(err)
	}
	if len(classes) > 0 {
		t.Errorf("Database should be empty, yet contains %d classes.", len(classes))
	}
}

func TestClose(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Using a closed session should panic")
		}
	}()

	mydb := getDB()
	mydb.Close()
	_ = mydb.session.Ping()
}

func TestPut(t *testing.T) {
	mydb := getDB()
	mydb.Purge()

	err := mydb.Put(sampleClass)
	if err != nil {
		t.Error("Put returned error: ", err)
	}
}

func TestLookup(t *testing.T) {
	mydb := getDB()
	mydb.Purge()

	err := mydb.Put(sampleClass)
	if err != nil {
		t.Error("Put returned error: ", err)
	}
	class, err := mydb.Lookup("CS", "125")
	if err != nil {
		t.Error("Class lookup returned error: ", err)
	}
	if class.Department != "CS" || class.CourseNumber != "125" {
		t.Error("Lookup result inaccurate: ", class)
	}
}

func TestGetAll(t *testing.T) {
	mydb := getDB()
	mydb.Purge()

	for i := 0; i < 10; i++ {
		err := mydb.Put(types.Class{
			Department: string(i),
		})
		if err != nil {
			t.Error("Put resulted in error: ", err)
		}
	}

	classes, err := mydb.GetAll()
	if err != nil {
		t.Error("GetAll resulted in error: ", err)
	}

	if len(classes) != 10 {
		t.Errorf("GetAll returned %d classes. Expected: %s", len(classes), 10)
	}
}
