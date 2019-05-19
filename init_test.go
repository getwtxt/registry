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
	timeMonthPrevRFC, err := timeMonthPrev.MarshalText()
	quickErr(err)

	timeTwoMonthsPrev := time.Now().AddDate(0, -2, 0)
	timeTwoMonthsPrevRFC, err := timeTwoMonthsPrev.MarshalText()
	quickErr(err)

	timeThreeMonthsPrev := time.Now().AddDate(0, -3, 0)
	timeThreeMonthsPrevRFC, err := timeThreeMonthsPrev.MarshalText()
	quickErr(err)

	timeFourMonthsPrev := time.Now().AddDate(0, -4, 0)
	timeFourMonthsPrevRFC, err := timeFourMonthsPrev.MarshalText()
	quickErr(err)

	var mockusers = []struct {
		url     string
		nick    string
		date    time.Time
		apidate []byte
		status  TimeMap
	}{
		{
			url:     "https://example3.com/twtxt.txt",
			nick:    "foo_barrington",
			date:    timeTwoMonthsPrev,
			apidate: timeTwoMonthsPrevRFC,
			status: TimeMap{
				timeTwoMonthsPrev: "foo_barrington\thttps://example3.com/twtxt.txt\t" + string(timeTwoMonthsPrevRFC) + "\tJust got started with #twtxt!",
				timeMonthPrev:     "foo_barrington\thttps://example3.com/twtxt.txt\t" + string(timeMonthPrevRFC) + "\tHey <@foo https://example.com/twtxt.txt>, I love programming. Just FYI.",
			},
		},
		{
			url:     "https://example.com/twtxt.txt",
			nick:    "foo",
			date:    timeFourMonthsPrev,
			apidate: timeFourMonthsPrevRFC,
			status: TimeMap{
				timeFourMonthsPrev:  "foo\thttps://example3.com/twtxt.txt\t" + string(timeFourMonthsPrevRFC) + "\tThis is so much better than #twitter",
				timeThreeMonthsPrev: "foo\thttps://example3.com/twtxt.txt\t" + string(timeThreeMonthsPrevRFC) + "\tI can't wait to start on my next programming #project with <@foo_barrington https://example3.com/twtxt.txt>",
			},
		},
	}
	index := NewIndex()

	// fill the test index with the mock users
	for _, e := range mockusers {
		data := &Data{}
		data.Nick = e.nick
		data.APIdate = e.apidate
		data.Status = e.status
		index.Reg[e.url] = data
	}

	return index
}
