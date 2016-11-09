[![Build Status](https://travis-ci.org/Hawry/events-are-square.svg?branch=master)](https://travis-ci.org/Hawry/events-are-square)

# SquareSpace event calendar to iCal

SquareSpace currently doesn't support export of entire event calendars to iCal/vCal format to use with i.e. Google Calendar. This little tool acts as a proxy between SquareSpace and an iCalendar provider. The project is still considered as a work in progress, but is currently working as intended - though the number of features are quite low.

Tested with GoLang v1.7 but should work on Go 1.4+, and should work on both 32 and 64-bit architectures.

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
Please use the supplied makefile if possible, since this will add a build version to your binary which will make it a lot easier to troubleshoot in the future. To use the makefile:

**64-bit architectures**: `make`

**32-bit architectures**: `make 32`

If you don't wish to use the makefile, or don't have the possibility due to other reasons, use the go compiler:
`go build -o eas` *Please note that this will not include a build version*

## Usage
The idea is that EaS will act as a proxy between Google (or the provider of your choice) and SquareSpace. Start the EaS-server on a publically available server and then import a webcalendar in Google by using the EaS-as proxy:

`http://your-eas-server.domain.com/?url=http://your.squarespace.com/calendar?format=json`

### Flags and runtime arguments
```
usage: eas [<flags>]

Flags:
      --help                 Show context-sensitive help (also try --help-long
                             and --help-man).
  -p, --port=8080            port to listen for incoming requests on
  -t, --topdomain=hawry.net  restrict calendar requests to a specific top-domain
  -z, --timezone=TIMEZONE    add timezoneid to all events
  -o, --offset=0             add number of hours as offset and (fake) UTC
```

#### Timezones
Unless otherwise specified, the timezone will be in UTC (Zulu-time). To change timezone append the flag `-z` or `--timezone=TIMEZONE` where TIMEZONE is the country code according to the [Zone.tab][1] file.

#### Offset
This flag should only be used when you experience weird time offsets in your Google Calendar vs SquareSpace and any other third-party application you might use. This will add an offset (negative or positive) to the time that SquareSpace are giving you. The reason for this might be that a third-party calendar/application might be poorly

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
