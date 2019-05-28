package timeut

import (
	"strconv"
	"time"

	"github.com/peterhellberg/duration"

	ers "github.com/colt3k/nglog/ers/bserr"
	log "github.com/colt3k/nglog/ng"
)

// https://golang.org/pkg/time/#pkg-constants
// https://golang.org/src/time/format.go

const (
	twelveTime    = "03:04:05PM"
	twofourTime   = "15:04:05"
	refTimeLayout = "Mon Jan 2 15:04:05 MST 2006"
)

//ParseRFC3339 go get -u github.com/peterhellberg/duration
// i.e. "P1DT30H4S"  Output: 54h0m4s
func ParseRFC3339(dateTime string) time.Duration {
	d, err := duration.Parse(dateTime)
	if ers.NotErr(err) {
		log.Println(d)
		return d
	}

	return -1
}

func ConvertUnix2Time(unxTime int64) time.Time {
	return  time.Unix(unxTime, 0)
}

func ConvertUnix2TimeStr(unxTime string) time.Time {
	i, err := strconv.ParseInt(unxTime, 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(i, 0)
}

type MyTime struct {
	time.Time
}

func Time(c time.Time) *MyTime {
	t := new(MyTime)
	t.Time = c
	return t
}
func GMTTime() *MyTime {
	t := new(MyTime)
	location, _ := time.LoadLocation("Europe/Rome")

	// this should give you time in location
	tNow := time.Now().In(location)
	t.Time = tNow

	return t
}
func (m *MyTime) GMT() time.Time {
	location, _ := time.LoadLocation("Europe/Rome")

	// this should give you time in location
	t := time.Now().In(location)
	m.Time = t
	return t
}

func (m *MyTime) Add(amt int, duration time.Duration) time.Time {
	tmp := m.Time.Add(time.Duration(amt) * duration)
	m.Time = tmp
	return m.Time
}

func (m *MyTime) Sub(amt int, duration time.Duration) time.Time {
	tmp := m.Time.Add(time.Duration((-1 * amt)) * duration)
	m.Time = tmp
	return m.Time
}
func (m *MyTime) Diff(t time.Time) time.Duration {
	return m.Time.Sub(t)
}
func (m *MyTime) ConvertTo24Hr() string {
	return m.Format(twofourTime)
}
func (m *MyTime) Millis() int64 {
	return m.UnixNano() / 1000000
}

// Utilities to determine the Monday of a week
// y, w := t.ISOWeek()
// i.e. log.Println(StartTime(y, w, time.UTC))
func StartTime(wyear, week int, loc *time.Location) (start time.Time) {
	y, m, d := StartDate(wyear, week)
	return time.Date(y, m, d, 0, 0, 0, 0, loc)
}

func StartDate(wyear, wk int) (year int, month time.Month, day int) {
	return Julian2Date(Date2Julian(wyear, 1, 1)+startOffset(wyear, wk))
}
func startOffset(y, week int) (offset int) {
	// This is optimized version of the following:
	//
	// return week*7 - ISOWeekday(y, 1, 4) - 3
	//
	// Uses Tomohiko Sakamoto's algorithm for calculating the weekday.

	y = y - 1
	return week*7 - (y+y/4-y/100+y/400+3)%7 - 4
}


func Julian2Date(dayNbr int) (y int, m time.Month, d int){
	e := 4*(dayNbr+1401+(4*dayNbr+274277)/146097*3/4-38) + 3
	h := e%1461/4*5 + 2

	d = h%153/5 + 1
	m = time.Month((h/153+2)%12 + 1)
	y = e/1461 - 4716 + (14-int(m))/12

	return y, m, d
}

func Date2Julian(y int, m time.Month, d int) int {

	if m < 3 {
		y = y -1
		m = m + 12
	}
	y = y + 4800

	return d + (153*(int(m)-3)+2)/5 + 365*y +
		y/4 - y/100 + y/400 - 32045
}