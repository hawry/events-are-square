package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
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

//Upcoming is a parent struct for all events
type Upcoming struct {
	Events []Event `json:"upcoming"`
}

// (^(http://)?([^\.]*)\.?(helsingkrona.se){1}|^(http://)?([^\.]*)\.?(bproduction.se){1}) and just continue until end of flags :)

//RegexDomain ... contains the regex pattern to discover top domains in a flag
const RegexDomain = "^(http://)?([^\\.]*)\\.?(%s){1}"

var builtRegex string

var (
	append    = kingpin.Flag("autoappend", "append 'format=json' to source URL automatically").Short('a').Default("false").Bool()
	port      = kingpin.Flag("port", "port to listen for incoming requests on").PlaceHolder("8080").Short('p').Default("8080").Int()
	topdomain = kingpin.Flag("topdomain", "restrict calendar requests to a specific top-domain").PlaceHolder("hawry.net").Short('t').String()
	timezone  = kingpin.Flag("timezone", "add timezoneid to all events").Short('z').String()
	offset    = kingpin.Flag("offset", "add number of hours as offset and (fake) UTC").Short('o').Default("0").Int()
	usrTZ     = "UTC"
	zoneMap   map[string]string
)

func fetchEvents(url string) (string, error) {
	if *append {
		url = fmt.Sprintf("%s?format=json", url)
		colog.ParseFields(false)
		log.Printf("debug: appended to total '%s'", url)
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
	if rsp.StatusCode == 429 {
		rerr := fmt.Errorf("remote end returned a %s error", http.StatusText(429))
		return "", rerr
	}

	w := Upcoming{}
	if err := json.Unmarshal(body, &w); err != nil {
		rerr := fmt.Errorf("could not unmarshal response body (%v)", err)
		log.Printf("Dump: %v", string(body))
		return "", rerr
	}
	log.Printf("info: source format OK, unmarshalling")
	var sVal string
	sVal += "BEGIN:VCALENDAR\r\n"
	sVal += "VERSION:2.0\r\n"
	for _, e := range w.Events {
		sVal += "BEGIN:VEVENT\r\n"
		uid := fmt.Sprintf("UID:%s\r\n", e.ID)
		start := fmt.Sprintf("DTSTART%s\r\n", to8601(e.StartDate))
		end := fmt.Sprintf("DTEND%s\r\n", to8601(e.EndDate))
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
	colog.ParseFields(false)
	defer colog.ParseFields(true)

	url := r.FormValue("url")
	if !(len(url) > 0) || url == "" {
		log.Printf("warning: could not find url in request (%v)", r.Header)
		w.WriteHeader(405)
		return
	}

	if len(*topdomain) > 0 {
		regex := regexp.MustCompile(builtRegex)
		if s := regex.FindStringSubmatch(url); !(len(s) > 0) {
			log.Printf("%+v", s)
			log.Printf("alert: unauthorized request url found '%s'", url)
			w.WriteHeader(405)
			return
		}
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

	log.Printf("system information: %v, %v, %v", runtime.GOOS, runtime.GOARCH, strconv.IntSize)

	log.Printf("running as server\t%t", true)
	log.Printf("listen port\t%d", *port)
	log.Printf("autoappend\t%t", *append)

	if len(*topdomain) > 0 {
		log.Printf("only serving requests from domain '%s'", *topdomain)
		builtRegex = fmt.Sprintf(RegexDomain, *topdomain)
	}

	if len(*timezone) > 0 {
		zoneMap = createMap()
		usrTZ = "UTC"
		if !(len(zoneMap[*timezone]) > 0) {
			log.Printf("warning: invalid timezone chosen '%s'. All times will be in %s", *timezone, usrTZ)
		} else {
			usrTZ = zoneMap[*timezone]
			log.Printf("info: user defined timezone is '%s'", usrTZ)
		}
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

//to8601 reformats a unix timestamp from json-timestamp to ISO-8601 in UTC (YYYYMMDDTHHmmssZ)
func to8601(t int64) string {
	t /= 1000
	t += fixOffset()
	ts := time.Unix(t, 0)
	if strings.Compare(usrTZ, "UTC") != 0 {
		sTime := strftime.Format("%Y%m%dT%H%M%S", ts.Local())
		return fmt.Sprintf(";TZID=%s:%s", usrTZ, sTime)
	}
	sTime := strftime.Format("%Y%m%dT%H%M%SZ", ts.UTC())
	return fmt.Sprintf(":%s", sTime)
}

func fixOffset() int64 {
	return int64((*offset) * 60 * 60)
}
