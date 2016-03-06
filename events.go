package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	append = kingpin.Flag("autoappend", "append 'format=pretty-json' to source URL automatically").Short('a').Default("false").Bool()
	// server = kingpin.Flag("srv", "run as server (false=run once and log to file instead of serving web requests)").Short('d').Default("true").Bool()
	port = kingpin.Flag("port", "port to listen for incoming requests on").Short('p').Default("8080").Int()
)

func fetchEvents(url string) (string, error) {
	if *append {
		url = fmt.Sprintf("%s?format=json", url)
	}
	rsp, err := http.Get(url)
	if err != nil {
		rerr := fmt.Errorf("could not open url '%s' (%v)", url, err)
		return "", rerr
	}

	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		rerr := fmt.Errorf("could not read stream content (%v)", err)
		return "", rerr
	}
	w := Upcoming{}
	if err := json.Unmarshal(body, &w); err != nil {
		rerr := fmt.Errorf("could not unmarshal response body (%v)", err)
		return "", rerr
	}
	log.Printf("info: source format OK, unmarshalling")
	var sVal string
	sVal += "BEGIN:VCALENDAR\r\n"
	sVal += "VERSION:2.0\r\n"
	for _, e := range w.Events {
		sVal += "BEGIN:VEVENT\r\n"
		uid := fmt.Sprintf("UID:%s\r\n", e.Id)
		start := fmt.Sprintf("DTSTART:%s\r\n", to8601(e.StartDate))
		end := fmt.Sprintf("DTEND:%s\r\n", to8601(e.EndDate))
		summary := fmt.Sprintf("SUMMARY:%s\r\n", e.Title)
		desc := fmt.Sprintf("DESCRIPTION:%s\r\n", strip.StripTags(e.Body))
		sVal += uid + start + end + summary + desc
		sVal += "END:VEVENT\r\n"
	}
	sVal += "END:VCALENDAR\r\n"
	log.Printf("info: sending response in VCAL format to requester")
	return sVal, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if !(len(url) > 0) || url == "" {
		log.Printf("warning: could not find url in request (%v)", r.Header)
		w.WriteHeader(405)
		return
	}
	log.Printf("info: fetching resources at '%s'", url)
	respBody, err := fetchEvents(url)
	if err != nil {
		log.Printf("err: could not fetch events (%v)", err)
		w.WriteHeader(405)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(respBody))
	log.Printf("success: sending parsed ical format to requester")
}

func main() {
	http.HandleFunc("/", handler)

	kingpin.Parse()
	colog.Register()
	colog.ParseFields(true)
	colog.SetFlags(log.LstdFlags)
	colog.SetDefaultLevel(colog.LInfo)

	log.Printf("running as server\t%t", true)
	log.Printf("listen port\t%d", *port)
	log.Printf("autoappend\t%t", *append)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

//to8601 reformats a unix timestamp from json-timestamp to ISO-8601 in UTC (YYYYMMDDTHHmmssZ)
func to8601(t int) string {
	s := int64(t)
	s /= 1000
	ts := time.Unix(s, 0)
	return strftime.Format("%Y%m%dT%H%M%SZ", ts.UTC())
}
