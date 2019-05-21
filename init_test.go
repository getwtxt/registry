package registry //import "github.com/getwtxt/registry"

import (
	"log"
	"os"
	"time"
)

func quickErr(err error) {
	if err != nil {
		log.Printf("%v\n", err)
	}
}

// Sets up mock users and statuses
func initTestEnv() *Index {
	hush, err := os.Open("/dev/null")
	quickErr(err)
	log.SetOutput(hush)

	// this is a bit tedious, but set up fake dates
	// for the mock users' join and status timestamps
	timeMonthPrev := time.Now().AddDate(0, -1, 0)
	timeMonthPrevRFC := timeMonthPrev.Format(time.RFC3339)

	timeTwoMonthsPrev := time.Now().AddDate(0, -2, 0)
	timeTwoMonthsPrevRFC := timeTwoMonthsPrev.Format(time.RFC3339)

	timeThreeMonthsPrev := time.Now().AddDate(0, -3, 0)
	timeThreeMonthsPrevRFC := timeThreeMonthsPrev.Format(time.RFC3339)

	timeFourMonthsPrev := time.Now().AddDate(0, -4, 0)
	timeFourMonthsPrevRFC := timeFourMonthsPrev.Format(time.RFC3339)

	var mockusers = []struct {
		url     string
		nick    string
		date    string
		apidate []byte
		status  TimeMap
	}{
		{
			url:  "https://example3.com/twtxt.txt",
			nick: "foo_barrington",
			date: timeTwoMonthsPrevRFC,
			status: TimeMap{
				timeTwoMonthsPrev: timeTwoMonthsPrevRFC + "\tJust got started with #twtxt!",
				timeMonthPrev:     timeMonthPrevRFC + "\tHey <@foo https://example.com/twtxt.txt>, I love programming. Just FYI.",
			},
		},
		{
			url:  "https://example.com/twtxt.txt",
			nick: "foo",
			date: timeFourMonthsPrevRFC,
			status: TimeMap{
				timeFourMonthsPrev:  timeFourMonthsPrevRFC + "\tThis is so much better than #twitter",
				timeThreeMonthsPrev: timeThreeMonthsPrevRFC + "\tI can't wait to start on my next programming #project with <@foo_barrington https://example3.com/twtxt.txt>",
			},
		},
	}
	index := NewIndex()

	// fill the test index with the mock users
	for _, e := range mockusers {
		data := &Data{}
		data.Nick = e.nick
		data.Date = e.date
		data.Status = e.status
		index.Reg[e.url] = data
	}

	return index
}
