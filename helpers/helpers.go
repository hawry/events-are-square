package helpers

import (
	"fmt"
	"strings"
	"time"

	strftime "github.com/jehiah/go-strftime"
)

//To8601 reformats a unix timestamp from json-timestamp to ISO-8601 in UTC (YYYYMMDDTHHmmssZ)
func To8601(t int64, tz string, offset int) string {
	t /= 1000
	t += fixOffset(offset)
	ts := time.Unix(t, 0)
	if strings.Compare(tz, "UTC") != 0 {
		sTime := strftime.Format("%Y%m%dT%H%M%S", ts.Local())
		return fmt.Sprintf(";TZID=%s:%s", tz, sTime)
	}
	sTime := strftime.Format("%Y%m%dT%H%M%SZ", ts.UTC())
	return fmt.Sprintf(":%s", sTime)
}

func fixOffset(offset int) int64 {
	return int64((offset) * 60 * 60)
}
