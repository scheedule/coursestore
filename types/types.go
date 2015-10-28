// Package types holds types that are used across the coursestore.
// Types are tagged for xml unmarshalling and bson serializing.
package types

import "gopkg.in/mgo.v2/bson"

type (
	// Type to unmarshal course types from the UIUC CISAPI
	CourseType struct {
		Name string `xml:",chardata"`
		Code string `xml:"id,attr"`
	}

	// Type to unmarshal meeting data from the UIUC CISAPI
	Meeting struct {
		Type  CourseType `xml:"type"`
		Start string     `xml:"start" bson:"start"`
		End   string     `xml:"end" bson:"end"`
		Days  string     `xml:"daysOfTheWeek" bson:"days"`
	}

	// Type to unmarshal section data from the UIUC CISAPI
	Section struct {
		CRN      string    `xml:"id,attr" bson:"crn"`
		Code     string    `xml:"sectionNumber" bson:"code"`
		Meetings []Meeting `xml:"meetings>meeting" bson:"meetings"`
	}

	// Type to unmarshal class data from the UIUC CISAPI
	Class struct {
		ID           bson.ObjectId `bson:"_id,omitempty"`
		Department   string        `bson:"department"`
		CourseNumber string        `bson:"course_number"`
		Name         string        `bson:"name"`
		Sections     []Section     `bson:"sections"`
	}
)
