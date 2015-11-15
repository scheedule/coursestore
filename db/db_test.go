package db

import (
	"os"
	"testing"

	"github.com/scheedule/coursestore/types"
)

func TestNewDB(t *testing.T) {
	myDB := New(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_COLLECTION"))
	if myDB == nil {
		t.Fail()
	}
}

func TestInit(t *testing.T) {
	myDB := New(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_COLLECTION"))
	err := myDB.Init()
	if err != nil {
		t.Error("Failed to initialize DB:", err)
	}
	if myDB.session == nil {
		t.Error("Session is nil")
	}
	if myDB.collection == nil {
		t.Error("Collection is nil")
	}
}

func getDB() *DB {
	myDB := New(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_COLLECTION"))
	_ = myDB.Init()
	return myDB
}

var sampleClass = types.Class{
	Department:   "CS",
	CourseNumber: "125",
}

func TestPurge(t *testing.T) {
	myDB := getDB()
	myDB.Put(sampleClass)
	myDB.Purge()
	classes, err := myDB.GetAll(nil)
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

	myDB := getDB()
	myDB.Close()
	_ = myDB.session.Ping()
}

func TestPut(t *testing.T) {
	myDB := getDB()
	myDB.Purge()

	err := myDB.Put(sampleClass)
	if err != nil {
		t.Error("Put returned error: ", err)
	}
}

func TestLookup(t *testing.T) {
	myDB := getDB()
	myDB.Purge()

	err := myDB.Put(sampleClass)
	if err != nil {
		t.Error("Put returned error: ", err)
	}
	class, err := myDB.Lookup("CS", "125")
	if err != nil {
		t.Error("Class lookup returned error: ", err)
	}
	if class.Department != "CS" || class.CourseNumber != "125" {
		t.Error("Lookup result inaccurate: ", class)
	}
}

func TestGetAll(t *testing.T) {
	myDB := getDB()
	myDB.Purge()

	for i := 0; i < 10; i++ {
		err := myDB.Put(types.Class{
			Department: string(i),
		})
		if err != nil {
			t.Error("Put resulted in error: ", err)
		}
	}

	classes, err := myDB.GetAll(DetailBasic)
	if err != nil {
		t.Error("GetAll resulted in error: ", err)
	}

	if len(classes) != 10 {
		t.Errorf("GetAll returned %d classes. Expected: %s", len(classes), 10)
	}
}
