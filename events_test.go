package main

import "testing"

func TestTimestamp(t *testing.T) {
	correctFormat := "20160406T160000Z"
	var jsTimestamp int64
	jsTimestamp = 1459958400000
	correct := to8601(jsTimestamp)
	if correct != correctFormat {
		t.Fail()
	}
}
