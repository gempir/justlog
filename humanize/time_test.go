package humanize

import (
	"testing"
	"time"
)

var timeSinceTests = []struct {
	in  time.Time
	out string
}{
	{time.Now().AddDate(0, 0, -1), "1 day"},
	{time.Now().AddDate(0, -1, -1), "1 month 1 day"},
	{time.Now().Add(time.Minute * -10), "10 mins"},
	{time.Now().Add(time.Minute * -10 + time.Second * -30), "10 mins 30 secs"},
}

func TestTimeSince(t *testing.T) {
	for _, testCase := range timeSinceTests {
		if since := TimeSince(testCase.in); since != testCase.out {
			t.Errorf("Incorrect time since string. Expected %s, Actual: %s", testCase.out, since)
		}
	}
}