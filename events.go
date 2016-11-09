package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/hawry/events-are-square/event"

	"comail.io/go/colog"
)

//Website is a shorthand for the map[string]interface{}
type Website map[string]interface{}

// (^(http://)?([^\.]*)\.?(helsingkrona.se){1}|^(http://)?([^\.]*)\.?(bproduction.se){1}) and just continue until end of flags :)

//RegexDomain ... contains the regex pattern to discover top domains in a flag
const RegexDomain = "^(http://)?([^\\.]*)\\.?(%s){1}"

var builtRegex string
var buildVersion string

var (
	port      = kingpin.Flag("port", "port to listen for incoming requests on").PlaceHolder("8080").Short('p').Default("8080").Int()
	topdomain = kingpin.Flag("topdomain", "restrict calendar requests to a specific top-domain").PlaceHolder("hawry.net").Short('t').String()
	timezone  = kingpin.Flag("timezone", "add timezoneid to all events").Short('z').String()
	offset    = kingpin.Flag("offset", "add number of hours as offset and (fake) UTC").Short('o').Default("0").Int()
	version   = kingpin.Flag("version", "only show build version number of eas and then quit").Bool()
	usrTZ     = "UTC"
	zoneMap   map[string]string
)

func handler(w http.ResponseWriter, r *http.Request) {
	transactionID := time.Now().Unix()
	colog.SetPrefix(fmt.Sprintf("[%d-%d]", transactionID, rand.Intn(100)))
	colog.ParseFields(false)
	defer colog.ParseFields(true)

	log.Printf("info: proxy request from '%s'", r.RemoteAddr)

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

	//Append ?format=json unless this is already taken care of
	if !strings.HasSuffix(url, "?format=json") {
		url = fmt.Sprintf("%s?format=json", url)
	}

	//Open the URL and fetch the underlying Reader
	rsp, err := http.Get(url)
	if err != nil {
		log.Printf("error: could not open url '%s' (%v)", url, err)
		return
	}
	defer rsp.Body.Close()

	if rsp.StatusCode == 429 {
		log.Printf("error: remote end returned a '%s' error", http.StatusText(429))
		return
	}

	eventList, err := event.Parse(rsp.Body)
	if err != nil {
		log.Printf("error: could not parse response body for url '%s' (%v)", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Add timezone and offset info
	eventList.TimeZone = usrTZ
	eventList.Offset = *offset

	log.Printf("success: fetched calendar information for '%s'", url)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(eventList.VCal()))
}

func main() {
	kingpin.Parse()

	if *version {
		//Only print binary version and quit
		log.Printf("info: build version %s", buildVersion)
		return
	}

	http.HandleFunc("/", handler)

	colog.Register()
	colog.ParseFields(true)
	colog.SetFlags(log.LstdFlags)
	colog.SetDefaultLevel(colog.LDebug)

	log.Printf("system information: %v, %v, %v", runtime.GOOS, runtime.GOARCH, strconv.IntSize)
	log.Printf("running as server\t%t", true)
	log.Printf("listen port\t%d", *port)

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
