// Package scrape implements all crawling and parsing associated with getting
// courses from the course store.
package scrape

import (
	"encoding/xml"
	"regexp"
	"strconv"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/scheedule/coursestore/types"
)

var normalizeCreditHoursRE = regexp.MustCompile(`\d+[\.]*[\d]*`)
var normalizeDegreeAttributesRE = regexp.MustCompile(`, and |, | course\.`)

type (
	// Type to unmarshal link XML from UIUC CISAPI
	Link struct {
		Href string `xml:"href,attr"`
	}

	// Type to unmarshal department XML from UIUC CISAPI
	Department struct {
		Courses []Link `xml:"courses>course"`
	}

	// Type to unmarshal term XML from UIUC CISAPI
	Term struct {
		Subjects []Link `xml:"subjects>subject"`
	}

	// Type to unmarshal subject XML from UIUC CISAPI
	Subject struct {
		Department string `xml:"id,attr"`
	}

	// Type to unmarshal course XML from UIUC CISAPI.
	Course struct {
		Number           string          `xml:"id,attr"`
		Name             string          `xml:"label"`
		Subject          Subject         `xml:"parents>subject"`
		Description      string          `xml:"description"`
		CreditHours      string          `xml:"creditHours"`
		DegreeAttributes string          `xml:"sectionDegreeAttributes"`
		Sections         []types.Section `xml:"detailedSections>detailedSection"`
	}
)

// Digest ALL course data from the DB
// Param: XMLData is list of departments
func DigestAll(XMLData []byte, courseChan chan types.Class) {
	log.Debug("starting term digestion")
	term := &Term{}
	err := xml.Unmarshal(XMLData, term)
	if err != nil {
		log.Fatal("failed to unmarshal XML: ", err)
	}

	var wg sync.WaitGroup

	for _, link := range term.Subjects {
		data, err := GetXML(link.Href)
		if err != nil {
			log.Fatal("error retrieving XML: ", err)
		}

		wg.Add(1)
		go digestDepartment(data, courseChan, &wg)
		log.Info("started: ", link.Href)
	}

	wg.Wait()

	close(courseChan)
	log.Info("digestion complete")
}

// Digest all courses from a given department
// Param: XMLData is list of courses for the department
func digestDepartment(XMLData []byte, courseChan chan types.Class, wg *sync.WaitGroup) {
	defer wg.Done()

	department := &Department{}
	err := xml.Unmarshal(XMLData, department)
	if err != nil {
		log.Fatal("failed to unmarshal XML: ", err)
	}

	for _, course := range department.Courses {
		data, err := GetXML(course.Href + "?mode=detail")
		if err != nil {
			log.Fatal("error retrieving XML: ", err)
		}

		c, err := digestClass(data)
		if err != nil {
			log.Fatal("failed to digest class: ", err)
		}

		courseChan <- *c
	}
}

// Extract credit hour numbers from course API string
func normalizeCreditHours(str string) string {
	matches := normalizeCreditHoursRE.FindAllString(str, -1)

	if len(matches) == 1 {
		return matches[0]
	}
	if len(matches) == 2 {
		return matches[0] + "-" + matches[1]
	}

	log.Error("encountered unmatched credit hour string: ", str)
	return ""
}

// Extract Degree Attributes from course API string
func normalizeDegreeAttributes(str string) []string {
	str = normalizeDegreeAttributesRE.ReplaceAllString(str, ",")
	split := strings.Split(str, ",")

	result := make([]string, 0, len(str))

	for _, s := range split {
		if s != "" {
			result = append(result, s)
		}
	}

	return result
}

// Digest all sections from a given class
// Param: XMLData is list of sections
func digestClass(XMLData []byte) (*types.Class, error) {
	course := &Course{}
	err := xml.Unmarshal(XMLData, course)
	if err != nil {
		log.Error("failed to unmarshal XML")
		return nil, err
	}

	// Remove whitespace
	for i, section := range course.Sections {
		course.Sections[i].Code = strings.TrimSpace(section.Code)
		for j, meeting := range section.Meetings {
			course.Sections[i].Meetings[j].Days = strings.TrimSpace(meeting.Days)
		}
	}

	courseNumber, _ := strconv.Atoi(strings.Split(course.Number, " ")[1])

	// Create Class struct
	class := &types.Class{
		Department:       course.Subject.Department,
		CourseNumber:     courseNumber,
		Name:             course.Name,
		Description:      course.Description,
		CreditHours:      normalizeCreditHours(course.CreditHours),
		DegreeAttributes: normalizeDegreeAttributes(course.DegreeAttributes),
		Sections:         course.Sections,
	}

	return class, nil
}
