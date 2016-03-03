package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/hawry/events-are-square/strip"
	"github.com/jehiah/go-strftime"

	"comail.io/go/colog"
)

//Website is a shorthand for the map[string]interface{}
type Website map[string]interface{}

//Author represents author info in a specific event
type Author struct {
	Id                  string `json:"id"`
	LastLoginOn         int    `json:"lastLoginOn"`
	LastActiveOn        int    `json:"lastActiveOn"`
	IsDeactivated       bool   `json:"isDeactivated"`
	Deleted             bool   `json:"deleted"`
	DisplayName         string `json:"displayName"`
	FirstName           string `json:"firstName"`
	LastName            string `json:"lastName"`
	EmailVerified       bool   `json:"emailVerified"`
	Bio                 string `json:"bio"`
	RevalidateTimestamp int    `json:"revalidateTimestamp"`
	SystemGenerated     bool   `json:"systemGenerated"`
}

//StructuredContent conains information about the event startdate/enddate
type StructuredContent struct {
	Type      string `json:"_type"`
	StartDate int    `json:"startDate"`
	EndDate   int    `json:"endDate"`
}

//Items is a dummy struct as of now
type Items struct {
}

//Event represents a single event in the upcoming list
type Event struct {
	Id                string            `json:"id"`
	CollectionId      string            `json:"collectionId"`
	RecordType        int               `json:"recordType"`
	AddedOn           int               `json:"addedOn"`
	UpdatedOn         int               `json:"updatedOn"`
	PublishOn         int               `json:"publishOn"`
	AuthorId          string            `json:"authorId"`
	UrlId             string            `json:"urlId"`
	Title             string            `json:"title"`
	SourceUrl         string            `json:"sourceUrl"`
	Body              string            `json:"body"`
	Author            Author            `json:"author"`
	FullUrl           string            `json:"fullUrl"`
	AssetUrl          string            `json:"assetUrl"`
	ContentType       string            `json:"contentType"`
	StructuredContent StructuredContent `json:"structuredContent"`
	StartDate         int               `json:"startDate"`
	EndDate           int               `json:"endDate"`
	Items             []Items           `json:"items"`
}

//Upcoming is a parent struct for all events
type Upcoming struct {
	Events []Event `json:"upcoming"`
}

var (
	src    = kingpin.Flag("src", "source URL to fetch rss-feed from").Short('s').Default("http://localhost/events").String()
	append = kingpin.Flag("autoappend", "append 'format=pretty-json' to source URL automatically").Short('a').Default("false").Bool()
)

func main() {
	kingpin.UsageTemplate(kingpin.SeparateOptionalFlagsUsageTemplate)
	kingpin.Parse()
	colog.Register()
	colog.SetFlags(log.Lshortfile)

	rsp, err := http.Get(*src)
	if err != nil {
		log.Printf("err: could not open URL '%s' (%v)", *src, err)
		return //TODO log the error to file, and do not exit application
	}

	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("err: could not read from stream (%v)", err)
		return //TODO log error, do not exit
	}

	// var dat map[string]interface{}
	// var w map[string]interface{}
	w := Upcoming{}
	if err := json.Unmarshal(body, &w); err != nil {
		log.Printf("err: could not unmarshal file (%v)", err)
		return
	}

	log.Printf("ok: unmarshalled file")

	f, err := os.Create("./results.ical")
	if err != nil {
		log.Printf("err: could not create file")
		return
	}
	defer f.Close()

	log.Printf("BEGIN:VCALENDAR")
	f.WriteString("BEGIN:VCALENDAR\r\n")

	log.Printf("VERSION:2.0")
	f.WriteString("VERSION:2.0\r\n")

	for _, e := range w.Events {
		f.WriteString("BEGIN:VEVENT\r\n")
		uid := fmt.Sprintf("UID:%s\r\n", e.Id)
		start := fmt.Sprintf("DTSTART:%s\r\n", to8601(e.StartDate))
		end := fmt.Sprintf("DTEND:%s\r\n", to8601(e.EndDate))
		summary := fmt.Sprintf("SUMMARY:%s\r\n", e.Title)
		desc := fmt.Sprintf("DESCRIPTION:%s\r\n", strip.StripTags(e.Body))

		f.WriteString(uid + start + end + summary + desc)
		f.WriteString("END:VEVENT\r\n")
	}

	log.Printf("END:VCALENDAR")
	f.WriteString("END:VCALENDAR\r\n")

}

func to8601(t int) string {
	s := int64(t)
	s /= 1000
	ts := time.Unix(s, 0)
	return strftime.Format("%Y%m%dT%H%M%SZ", ts.UTC())
}
