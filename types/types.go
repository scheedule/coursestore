// Package types holds types that are used across the coursestore.
// Types are tagged for xml unmarshalling and bson serializing.
package types

import "gopkg.in/mgo.v2/bson"

type (
	// Type to unmarshal instructor types from the UIUC CISAPI
	Instructor struct {
		FirstName string `xml:"firstName,attr" bson:"first" json:"first"`
		LastName  string `xml:"lastName,attr" bson:"last" json:"last"`
	}

	// Type to unmarshal course types from the UIUC CISAPI
	CourseType struct {
		Name string `xml:",chardata" bson:"name" json:"name"`
		Code string `xml:"code,attr" bson:"code" json:"code"`
	}

	// Type to unmarshal meeting data from the UIUC CISAPI
	Meeting struct {
		Type        CourseType   `xml:"type" bson:"type" json:"type"`
		Start       string       `xml:"start" bson:"start" json:"start"`
		End         string       `xml:"end" bson:"end" json:"end"`
		Days        string       `xml:"daysOfTheWeek" bson:"days" json:"days"`
		Building    string       `xml:"buildingName" bson:"building" json:"building,omitempty"`
		Instructors []Instructor `xml:"instructors>instructor" bson:"instructors" json:"instructors,omitempty"`
	}

	// Type to unmarshal section data from the UIUC CISAPI
	Section struct {
		CRN              string    `xml:"id,attr" bson:"crn" json:"crn"`
		Code             string    `xml:"sectionNumber" bson:"code" json:"code"`
		EnrollmentStatus string    `xml:"enrollmentStatus" bson:"enrollment_status" json:"enrollmentStatus,omitempty"`
		Start            string    `xml:"startDate" bson:"start" json:"start,omitempty"`
		End              string    `xml:"endDate" bson:"end" json:"end,omitempty"`
		Meetings         []Meeting `xml:"meetings>meeting" bson:"meetings" json:"meetings"`
	}

	// Type to unmarshal class data from the UIUC CISAPI
	Class struct {
		ID               bson.ObjectId `bson:"_id,omitempty" json:"-"`
		Department       string        `bson:"department" json:"department"`
		CourseNumber     string        `bson:"course_number" json:"courseNumber"`
		Name             string        `bson:"name" json:"name"`
		Description      string        `bson:"description" json:"description,omitempty"`
		CreditHours      string        `bson:"credit_hours" json:"creditHours,omitempty"`
		DegreeAttributes []string      `bson:"degree_attributes" json:"degreeAttributes,omitempty"`
		Sections         []Section     `bson:"sections" json:"sections,omitempty"`
	}
)
