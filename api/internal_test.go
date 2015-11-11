package api

import (
	"github.com/scheedule/coursestore/types"
	"testing"
)

var departmentTests = []struct {
	in  string
	out bool
}{
	{"", false},
	{"hello", false},
	{"123", false},
	{"CS", true},
	{"AAS", true},
}

func TestIsValidDepartment(t *testing.T) {
	for _, tt := range departmentTests {
		result := isValidDepartment(tt.in)
		if result != tt.out {
			t.Errorf("isValidDepartment(%q) => %q, want %q", tt.in, result, tt.out)
		}
	}
}

var courseNumberTests = []struct {
	in  string
	out bool
}{
	{"", false},
	{"hello", false},
	{"123", true},
	{"1234", true},
}

func TestIsValidCourseNumber(t *testing.T) {
	for _, tt := range courseNumberTests {
		result := isValidCourseNumber(tt.in)
		if result != tt.out {
			t.Errorf("isValidCourseNumber(%q) => %q, want %q", tt.in, result, tt.out)
		}
	}
}

var sampleClass = types.Class{
	Department:   "CS",
	CourseNumber: "125",
}

func TestLookupClass(t *testing.T) {
	testAPI.db.Purge()
	err := testAPI.db.Put(sampleClass)
	if err != nil {
		t.Error("failed to put class in database: ", err)
	}

	if _, err := testAPI.lookupClass("CS", "125"); err != nil {
		t.Error("lookupClass of sample resulted in error: ", err)
	}
}

func TestPackClasses(t *testing.T) {
	testAPI.db.Purge()

	for i := 0; i < 10; i++ {
		err := testAPI.db.Put(types.Class{
			Department: string(i),
		})
		if err != nil {
			t.Error("put resulted in error: ", err)
		}
	}

	classes, err := testAPI.db.GetAll()
	if err != nil {
		t.Error("function GetAll resulted in error: ", err)
	}

	if len(classes) != 10 {
		t.Error("incorrect number of classes %q, want %q", len(classes), 10)
	}
}
