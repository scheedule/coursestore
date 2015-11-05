package scrape

import (
	"github.com/scheedule/coursestore/types"
	"testing"
)

// Error if the string is empty
func emptyCheck(fieldname, str string, t *testing.T) {
	if str == "" {
		t.Errorf("Field %s empty when it shouldn't be", fieldname)
	}
}

// Error if the instructor contains any empty strings
func instructorEmptyCheck(instructor types.Instructor, t *testing.T) {
	emptyCheck("FirstName", instructor.FirstName, t)
	emptyCheck("LastName", instructor.LastName, t)
}

// Error if the courseType contains any empty strings
func courseTypeEmptyCheck(meetingType types.CourseType, t *testing.T) {
	emptyCheck("Code", meetingType.Code, t)
	emptyCheck("Name", meetingType.Name, t)
}

// Error if the meeting contains any empty strings
func meetingEmptyCheck(meeting types.Meeting, t *testing.T) {
	emptyCheck("Days", meeting.Days, t)
	emptyCheck("End", meeting.End, t)
	emptyCheck("Start", meeting.Start, t)
	courseTypeEmptyCheck(meeting.Type, t)
	for i := range meeting.Instructors {
		instructorEmptyCheck(meeting.Instructors[i], t)
	}
}

// Error if the section contains any empty strings
func sectionEmptyCheck(section types.Section, t *testing.T) {
	emptyCheck("CRN", section.CRN, t)
	emptyCheck("Code", section.Code, t)
	for i := range section.Meetings {
		meetingEmptyCheck(section.Meetings[i], t)
	}
}

// Error if the class contains any empty strings
func classEmptyCheck(class types.Class, t *testing.T) {
	emptyCheck("CourseNumber", class.CourseNumber, t)
	emptyCheck("Department", class.Department, t)
	emptyCheck("Name", class.Name, t)
	emptyCheck("Description", class.Description, t)
	emptyCheck("CreditHours", class.CreditHours, t)
	emptyCheck("DegreeAttributes", class.DegreeAttributes, t)
	for i := range class.Sections {
		sectionEmptyCheck(class.Sections[i], t)
	}
}

// Error if we can't retrieve an XML document
func TestGetXML(t *testing.T) {
	url := "http://courses.illinois.edu/cisapp/explorer/schedule/2016.xml"
	_, err := GetXML(url)
	if err != nil {
		t.Error(err)
	}
}

// Error if we can't retrieve an XML document. Testing the limiter is dificult,
// so let's just make sure it gets to the function.
func TestGetXMLLimiter(t *testing.T) {
	url := "http://courses.illinois.edu/cisapp/explorer/schedule/2016.xml"
	_, err := GetXML(url)
	if err != nil {
		t.Error(err)
	}
}

// Error if we fail to parse class data
func TestDigestClass(t *testing.T) {
	url := "http://courses.illinois.edu/cisapp/explorer/schedule/2016/spring/" +
		"AAS/100.xml?mode=detail"
	data, err := GetXML(url)
	if err != nil {
		t.Error(err)
	}

	class, err := digestClass(data)
	if err != nil {
		t.Error(err)
	}

	classEmptyCheck(*class, t)
}
