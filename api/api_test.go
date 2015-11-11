package api

import (
	"os"
	"testing"

	"github.com/scheedule/coursestore/db"
	"github.com/scheedule/coursestore/types"
)

var myAPI *API

var sampleClass = types.Class{
	Department:   "CS",
	CourseNumber: "125",
}

func init() {
	myDB := db.New(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_COLLECTION"))
	myDB.Init()
	myAPI = &API{myDB}
}

func failWithExpectedError(t *testing.T, expectation, result error) {
	t.Errorf("Expected %q, but got %q", expectation, result)
}

func TestLookupClass(t *testing.T) {
	myAPI.DB.Purge()
	var err error

	// Send empty values
	t.Log("Testing empty values")
	if _, err = myAPI.lookupClass("", ""); err != BadRequestError {
		failWithExpectedError(t, BadRequestError, err)
	}

	// Send one empty value
	t.Log("Testing one empty value")
	if _, err = myAPI.lookupClass("CS", ""); err != BadRequestError {
		failWithExpectedError(t, BadRequestError, err)
	}
	if _, err = myAPI.lookupClass("", "125"); err != BadRequestError {
		failWithExpectedError(t, BadRequestError, err)
	}

	// Send invalid department
	t.Log("Testing invalid department")
	if _, err = myAPI.lookupClass("hello", "125"); err != BadRequestError {
		failWithExpectedError(t, BadRequestError, err)
	}
	if _, err = myAPI.lookupClass("125", "125"); err != BadRequestError {
		failWithExpectedError(t, BadRequestError, err)
	}

	// Send invalid course number
	t.Log("Testing invalid course number")
	if _, err = myAPI.lookupClass("CS", "!!!"); err != BadRequestError {
		failWithExpectedError(t, BadRequestError, err)
	}
	if _, err = myAPI.lookupClass("CS", "\"Hello\""); err != BadRequestError {
		failWithExpectedError(t, BadRequestError, err)
	}

	// Test DB access
	t.Log("Testing DB Access")
	err = myAPI.DB.Put(sampleClass)
	if err != nil {
		t.Error("Putting class in db failed", err)
	}
	if _, err = myAPI.DB.Lookup("CS", "125"); err != nil {
		t.Error("Error upon looking up class", err)
	}
	if class, err := myAPI.DB.Lookup("CS", "225"); err == nil || class != nil {
		t.Errorf("Lookup of unknown course resulted in class: %q and err: %q",
			class, err)
	}
}

func TestPackClasses(t *testing.T) {

	// Clear database
	myAPI.DB.Purge()

	// Put data into the database
	for i := 0; i < 10; i++ {
		err := myAPI.DB.Put(types.Class{
			Department: string(i),
		})
		if err != nil {
			t.Error("Put resulted in error: ", err)
		}
	}

	classes, err := myAPI.DB.GetAll()
	if err != nil {
		t.Error("GetAll classes returned with: ", err)
	}

	if len(classes) != 10 {
		t.Error("GetAll didn't return the correct number of classes")
	}

	for i, class := range classes {
		if class.Department != string(i) {
			t.Error("GetAll returned inaccurate classes")
		}
	}

}
