package event

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/hawry/events-are-square/helpers"
	"github.com/hawry/events-are-square/strip"
)

//Author represents author info in a specific event
type Author struct {
	ID                  string `json:"id"`
	LastLoginOn         int64  `json:"lastLoginOn"`
	LastActiveOn        int64  `json:"lastActiveOn"`
	IsDeactivated       bool   `json:"isDeactivated"`
	Deleted             bool   `json:"deleted"`
	DisplayName         string `json:"displayName"`
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	EmailVerified       bool   `json:"emailVerified"`
	Bio                 string `json:"bio"`
	RevalidateTimestamp int64  `json:"revalidateTimestamp"`
	SystemGenerated     bool   `json:"systemGenerated"`
}

//StructuredContent conains information about the event startdate/enddate
type StructuredContent struct {
	Type      string `json:"_type"`
	StartDate int64  `json:"startDate"`
	EndDate   int64  `json:"endDate"`
}

//Items is a dummy struct as of now
type Items struct{}

//Event represents a single event in the upcoming list
type Event struct {
	ID                string            `json:"id"`
	CollectionID      string            `json:"collectionId"`
	RecordType        int               `json:"recordType"`
	AddedOn           int64             `json:"addedOn"`
	UpdatedOn         int64             `json:"updatedOn"`
	PublishOn         int64             `json:"publishOn"`
	AuthorID          string            `json:"authorId"`
	URLID             string            `json:"urlId"`
	Title             string            `json:"title"`
	SourceURL         string            `json:"sourceUrl"`
	Body              string            `json:"body"`
	Author            Author            `json:"author"`
	FullURL           string            `json:"fullUrl"`
	AssetURL          string            `json:"assetUrl"`
	ContentType       string            `json:"contentType"`
	StructuredContent StructuredContent `json:"structuredContent"`
	StartDate         int64             `json:"startDate"`
	EndDate           int64             `json:"endDate"`
	Items             []Items           `json:"items"`
}

//List is a parent struct for all events (name because lint suggests stuttering otherwise)
type List struct {
	Events   []Event `json:"upcoming"`
	Offset   int
	TimeZone string
}

//Parse will read the raw data, and unmarshal it into a List struct and return it
func Parse(r io.Reader) (*List, error) {
	eventList := List{}
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return &eventList, err
	}
	err = json.Unmarshal(raw, &eventList)
	if err != nil {
		return &eventList, err
	}
	return &eventList, nil
}

//VCal takes a List-struct and will return the VCAL-formatted string that it represents
func (eventList *List) VCal() string {
	var sVal string

	tz := eventList.TimeZone
	offset := eventList.Offset

	sVal += "BEGIN:VCALENDAR\r\n"
	sVal += "VERSION:2.0\r\n"
	for _, e := range eventList.Events {
		sVal += "BEGIN:VEVENT\r\n"
		uid := fmt.Sprintf("UID:%s\r\n", e.ID)
		start := fmt.Sprintf("DTSTART%s\r\n", helpers.To8601(e.StartDate, tz, offset))
		end := fmt.Sprintf("DTEND%s\r\n", helpers.To8601(e.EndDate, tz, offset))
		summary := fmt.Sprintf("SUMMARY:%s\r\n", e.Title)
		desc := fmt.Sprintf("DESCRIPTION:%s\r\n", strip.StripTags(e.Body))
		sVal += uid + start + end + summary + desc
		sVal += "END:VEVENT\r\n"
	}
	sVal += "END:VCALENDAR\r\n"
	return sVal
}
