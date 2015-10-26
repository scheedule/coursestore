package types

import "gopkg.in/mgo.v2/bson"

type (
	CourseType struct {
		Name string `xml:",chardata"`
		Code string `xml:"id,attr"`
	}

	Meeting struct {
		Type  CourseType `xml:"type"`
		Start string     `xml:"start" bson:"start"`
		End   string     `xml:"end" bson:"end"`
		Days  string     `xml:"daysOfTheWeek" bson:"days"`
	}

	Section struct {
		CRN      string    `xml:"id,attr" bson:"crn"`
		Code     string    `xml:"sectionNumber" bson:"code"`
		Meetings []Meeting `xml:"meetings>meeting" bson:"meetings"`
	}

	Class struct {
		ID           bson.ObjectId `bson:"_id,omitempty"`
		Department   string        `bson:"department"`
		CourseNumber string        `bson:"course_number"`
		Name         string        `bson:"name"`
		Sections     []Section     `bson:"sections"`
	}
)
