# SquareSpace event calendar to iCal-format

EaS was created to handle automatic parsing of SquareSpace event calendars and reformat them to iCalendar (vCalendar Version 2-format) with the simple purpose of being able to import them to a google calendar. The application is written in GoLang, and acts as a "proxy" between a icalendar compatible calendar and SquareSpace.

Tested with GoLang v1.6 but should work on GoLang 1.4 and 1.5 as well.

## Installation
`go get github.com/hawry/events-are-square`
or clone this repository.
### Dependencies
Make sure you have all the dependencies installed on your system by changing directory to the source-location of EaS and type:
`go get -u`

### Build
`go build`

## Usage
The idea is that EaS will act as a proxy between Google (or the provider of your choice) and SquareSpace. Start the EaS-server on a publically available server and then import a webcalendar in Google by using the EaS-as proxy:

`http://your-eas-server.domain.com/?url=http://your.squarespace.com/calendar?format=json`

**EaS should work on both 32- and 64-bit systems now.**

### Flags and runtime arguments
```
usage: events-are-square [<flags>]

Flags:
  --help        Show context-sensitive help
  -a, --autoappend  append 'format=pretty-json' to source URL automatically (default is false)
  -p, --port=8080   port to listen for incoming requests on
```
