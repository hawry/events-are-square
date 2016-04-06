package main

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/jehiah/go-strftime"
)

func TestTimestamp(t *testing.T) {
	correctFormat := "20160406T160000Z"
	var jsTimestamp int64
	jsTimestamp = 1459958400000
	correct := to8601(jsTimestamp)
	if correct != correctFormat {
		t.Fail()
	}
}

func TestRun(t *testing.T) {
	var js int64
	js = 1459958400000 / 1000
	ts := time.Unix(js, 0)
	log.Printf("Hello, world: %v", ts.Local())
	conv := strftime.Format("%Y%m%dT%H%M%S", ts.Local())
	log.Printf("converted: %s", conv)
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
