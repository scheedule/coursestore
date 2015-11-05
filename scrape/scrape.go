package scrape

import (
	"encoding/xml"
	log "github.com/Sirupsen/logrus"
	"github.com/scheedule/coursestore/types"
	"strings"
	"sync"
)

type (
	// Type to unmarshal link xml from UIUC CISAPI
	Link struct {
		Href string `xml:"href,attr"`
	}

	// Type to unmarshal department xml from UIUC CISAPI
	Department struct {
		Courses []Link `xml:"courses>course"`
	}

	// Type to unmarshal term xml from UIUC CISAPI
	Term struct {
		Subjects []Link `xml:"subjects>subject"`
	}

	// Type to unmarshal subject xml from UIUC CISAPI
	Subject struct {
		Department string `xml:"id,attr"`
	}

	// Type to unmarshal course xml from UIUC CISAPI.
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
// Param: xml_data is list of departments
func DigestAll(xml_data []byte, course_chan chan types.Class) {
	term := &Term{}
	err := xml.Unmarshal(xml_data, term)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for _, link := range term.Subjects {
		data, err := GetXML(link.Href)
		if err != nil {
			panic(err)
		}

		wg.Add(1)
		go digestDepartment(data, course_chan, &wg)
		log.Info("Started: ", link.Href)
	}

	wg.Wait()

	close(course_chan)
}

// Digest all courses from a given department
// Param: xml_data is list of courses for the department
func digestDepartment(xml_data []byte, course_chan chan types.Class, wg *sync.WaitGroup) {
	defer wg.Done()

	department := &Department{}
	err := xml.Unmarshal(xml_data, department)
	if err != nil {
		panic(err)
	}

	for _, course := range department.Courses {
		data, err := GetXML(course.Href + "?mode=detail")
		if err != nil {
			panic(err)
		}

		c, err := digestClass(data)
		if err != nil {
			panic(err)
		}

		course_chan <- *c
	}
}

// Digest all sections from a given class
// Param: xml_data is list of sections
func digestClass(xml_data []byte) (*types.Class, error) {
	course := &Course{}
	err := xml.Unmarshal(xml_data, course)
	if err != nil {
		return nil, err
	}

	// Remove whitespace
	for i, section := range course.Sections {
		course.Sections[i].Code = strings.TrimSpace(section.Code)
		for j, meeting := range section.Meetings {
			course.Sections[i].Meetings[j].Days = strings.TrimSpace(meeting.Days)
		}
	}

	// Create Class struct
	class := &types.Class{
		Department:       course.Subject.Department,
		CourseNumber:     strings.Split(course.Number, " ")[1],
		Name:             course.Name,
		Description:      course.Description,
		CreditHours:      course.CreditHours,
		DegreeAttributes: course.DegreeAttributes,
		Sections:         course.Sections,
	}

	return class, nil
}
