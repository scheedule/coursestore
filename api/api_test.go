package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/scheedule/coursestore/db"
	"github.com/scheedule/coursestore/types"
)

var testAPI *API

func init() {
	testDB := db.New(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_COLLECTION"))
	testDB.Init()
	testAPI = New(testDB)
}

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
			t.Fatalf("isValidDepartment(%q) => %q, want %q", tt.in, result, tt.out)
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
			t.Fatalf("isValidCourseNumber(%q) => %q, want %q", tt.in, result, tt.out)
		}
	}
}

var sampleClass = types.Class{
	Department:   "CS",
	CourseNumber: "125",
}

var classLookupTests = []struct {
	department string
	number     string
	code       string
}{
	{"CS", "125", "200"},
	{"CS", "225", "404"},
}

func TestLookup(t *testing.T) {
	// Input class into database
	testAPI.db.Purge()
	err := testAPI.db.Put(sampleClass)
	if err != nil {
		t.Fatal("failed to put class in database: ", err)
	}

	for _, tt := range classLookupTests {
		// Make HTTP Request
		urlStr := fmt.Sprintf("/lookup?department=%s&number=%s", tt.department, tt.number)
		req, err := http.NewRequest("GET", urlStr, nil)
		if err != nil {
			t.Fatal("failed to create request object.")
		}

		w := httptest.NewRecorder()
		testAPI.HandleLookup(w, req)

		code := fmt.Sprintf("%d", w.Code)
		if code != tt.code {
			t.Fatalf("response code received %q, want %q", code, tt.code)
		}

		if code != "200" {
			return
		}

		data, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Fatal("failed to read response.")
		}
		proposal := types.Class{}
		err = json.Unmarshal(data, &proposal)
		if err != nil {
			t.Fatal("failed to decode response: ", err)
		}

		if proposal.Department != tt.department {
			t.Fatalf("department contained %q, want %q", proposal.Department, tt.department)
		}
		if proposal.CourseNumber != tt.number {
			t.Fatalf("course number contained %q, want %q", proposal.CourseNumber, tt.department)
		}
	}
}
