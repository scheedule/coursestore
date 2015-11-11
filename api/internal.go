package api

import (
	"encoding/json"
	"regexp"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

// Return true if and only if the department is formatted correctly. This
// function does not check the database for department existence.
func isValidDepartment(department string) bool {
	// Check empty
	if department == "" {
		return false
	}

	// Check capitalized alphabetic.
	if matched, err := regexp.MatchString("^[A-Z]*$", department); !matched || err != nil {
		return false
	}

	return true
}

func isValidCourseNumber(number string) bool {
	// Check empty
	if number == "" {
		return false
	}

	// Check numeric
	if _, err := strconv.Atoi(number); err != nil {
		return false
	}

	return true
}

// Interrogate values and produce JSON.
func (a *API) lookupClass(department, number string) ([]byte, error) {

	class, err := a.db.Lookup(department, number)
	if err != nil {
		log.Warn("DB lookup failed: ", err)
		return nil, DBError
	}

	js, err := json.Marshal(class)
	if err != nil {
		log.Error("class unmarshal failed: ", err)
		return nil, UnmarshalError
	}

	return js, nil
}

// Pack all classes into JSON.
func (a *API) packClasses() ([]byte, error) {
	classes, err := a.db.GetAll()
	if err != nil {
		log.Error("failed to query all classes: ", err)
		return nil, DBError
	}

	js, err := json.Marshal(classes)
	if err != nil {
		log.Error("failed to unmarshal all classes: ", err)
		return nil, UnmarshalError
	}

	return js, nil
}
