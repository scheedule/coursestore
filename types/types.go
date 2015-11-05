// Package types holds types that are used across the coursestore.
// Types are tagged for xml unmarshalling and bson serializing.
package types

import "gopkg.in/mgo.v2/bson"

type (
	// Type to unmarshal instructor types from the UIUC CISAPI
	Instructor struct {
		FirstName string `xml:"firstName,attr" bson:"first"`
		LastName  string `xml:"lastName,attr" bson:"last"`
	}

	// Type to unmarshal course types from the UIUC CISAPI
	CourseType struct {
		Name string `xml:",chardata" bson:"name"`
		Code string `xml:"code,attr" bson:"code"`
	}

	// Type to unmarshal meeting data from the UIUC CISAPI
	Meeting struct {
		Type        CourseType   `xml:"type"`
		Start       string       `xml:"start" bson:"start"`
		End         string       `xml:"end" bson:"end"`
		Days        string       `xml:"daysOfTheWeek" bson:"days"`
		Building    string       `xml:"buildingName" bson:"building"`
		Instructors []Instructor `xml:"instructors>instructor" bson:"instructors"`
	}

	// Type to unmarshal section data from the UIUC CISAPI
	Section struct {
		CRN              string    `xml:"id,attr" bson:"crn"`
		Code             string    `xml:"sectionNumber" bson:"code"`
		EnrollmentStatus string    `xml:"enrollmentStatus" bson:"enrollment_status"`
		Start            string    `xml:"startDate" bson:"start"`
		End              string    `xml:"endDate" bson:"end"`
		Meetings         []Meeting `xml:"meetings>meeting" bson:"meetings"`
	}

	// Type to unmarshal class data from the UIUC CISAPI
	Class struct {
		ID               bson.ObjectId `bson:"_id,omitempty"`
		Department       string        `bson:"department"`
		CourseNumber     string        `bson:"course_number"`
		Name             string        `bson:"name"`
		Description      string        `bson:"description"`
		CreditHours      string        `bson:"credit_hours"`
		DegreeAttributes string        `bson:"degree_attributes"`
		Sections         []Section     `bson:"sections"`
	}
)
