[![Build Status](https://travis-ci.org/Hawry/events-are-square.svg?branch=master)](https://travis-ci.org/Hawry/events-are-square)

# SquareSpace event calendar to iCal

SquareSpace currently doesn't support export of entire event calendars to iCal/vCal format to use with i.e. Google Calendar. This little tool acts as a proxy between SquareSpace and an iCalendar provider. The project is still considered as a work in progress, but is currently working as intended - though the number of features are quite low.

Tested with GoLang v1.6 but should work on GoLang 1.4 and 1.5 as well, and should work on both 32-bit architectures as well as 64-bit.

### Current status
As of now, EaS only supports basic event-information and can convert SquareSpace event information to the following iCal tags:
`VEVENT: DTSTART,DTEND,SUMMARY,DESCRIPTION`

## Installation
`go get github.com/hawry/events-are-square`
or clone this repository.
### Dependencies
Make sure you have all the dependencies installed on your system by changing directory to the source-location of EaS and type:
`go get -u`

### Build
`go build -o events`

## Usage
The idea is that EaS will act as a proxy between Google (or the provider of your choice) and SquareSpace. Start the EaS-server on a publically available server and then import a webcalendar in Google by using the EaS-as proxy:

`http://your-eas-server.domain.com/?url=http://your.squarespace.com/calendar?format=json`

### Flags and runtime arguments
```
usage: events [<flags>]

Flags:
      --help                 Show context-sensitive help (also try
                             --help-long and --help-man).
  -a, --autoappend           append 'format=json' to source URL
                             automatically
  -p, --port=8080            port to listen for incoming requests on
  -t, --topdomain=hawry.net  restrict calendar requests to a specific
                             top-domain
  -z, --timezone=TIMEZONE    add timezoneid to all events
```

#### Timezones
Unless otherwise specified, the timezone will be in UTC (Zulu-time). To change timezone append the flag `-z` or `--timezone=TIMEZONE` where TIMEZONE is the country code according to the [Zone.tab][1] file.

## Planned features

*Please note that EaS is a work in progress and is developed during my free time, and therefore might take a while to be updated. You are very welcome to contribute to the project though!*

* Whitelisting/blacklisting of domains
  * Support for multiple domains to deny/allow
* Configuration file
* Adapting release to work with Docker
* Support for entire vCalendar specification
* Code cleanup
* More documentation & use cases
* SSL-support
* Simple systemd configuration

[1]: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
