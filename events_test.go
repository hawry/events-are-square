package main

import (
	"strings"
	"testing"
)

func TestTimestamp(t *testing.T) {
	correctFormat := ":20160406T160000Z"
	var jsTimestamp int64
	jsTimestamp = 1459958400000
	correct := to8601(jsTimestamp)
	if correct != correctFormat {
		t.Logf("failed because '%s'!='%s'", correct, correctFormat)
		t.Fail()
	}
}

func TestWithTimezone(t *testing.T) {
	usrTZ = "Europe/Stockholm"
	correctFormat := ";TZID=Europe/Stockholm:20160406T180000"
	var jsTimestamp int64
	jsTimestamp = 1459958400000
	correct := to8601(jsTimestamp)
	if correct != correctFormat {
		t.Logf("failed because '%s'!='%s'", correct, correctFormat)
		t.Fail()
	}
}

func TestDecode(t *testing.T) {
	m := createMap()
	if !(len(m) > 0) {
		t.Fail()
	}

	sval := m["SE"]
	if strings.Compare(sval, "Europe/Stockholm") != 0 {
		t.Fail()
	}
}
