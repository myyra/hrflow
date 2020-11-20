package hrflow

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// hrCalndarDay is used for parsing the JSON response, and is not very useful by itself.
type hrCalendarDay struct {
	// TES (työehtosopimus) is always empty for now.
	TES string `json:"tes"`
	// ID seems to correspond with the date.
	ID int `json:"id"`
	// Date in format "2020-04-10T00:00:00"
	Date string `json:"date"`
	// DateType seems empty for now.
	DateType string `json:"dateType"`
	// Workday means if you should work that day.
	Workday bool `json:"workDay"`
	// Holiday weirdly does not mean it's a holiday, but if that day should be used in holiday calculations according to the Finnish annual holiday system.
	Holiday bool `json:"holiday"`
	// Description seems to always be one space.
	Description string `json:"description"`
	// Number of the weekday, 1 being Monday and 7 Sunday.
	Weekday string `json:"weekDay"`
}

// CalendarDay is a simplified and more useful form of the response object.
type CalendarDay struct {
	// Date is the midninght of the day.
	Date time.Time
	// Workday tells if you should work that day.
	Workday bool
	// HolidayCalc tells if the date should be used in Finnish annual holiday system.
	HolidayCalc bool
	// Description explains what holiday it is.
	Description string
	// Weekday tells the day of the week.
	Weekday time.Weekday
}

func (d *hrCalendarDay) calendarDay() (CalendarDay, error) {

	var weekday time.Weekday

	switch d.Weekday {
	case "1":
		weekday = time.Monday
	case "2":
		weekday = time.Tuesday
	case "3":
		weekday = time.Wednesday
	case "4":
		weekday = time.Thursday
	case "5":
		weekday = time.Friday
	case "6":
		weekday = time.Saturday
	case "7":
		weekday = time.Sunday
	}

	date, err := time.Parse("2006-01-02T15:04:05", d.Date)
	if err != nil {
		return CalendarDay{}, errors.Wrap(err, "parsing date")
	}

	return CalendarDay{
		Date:        date,
		Workday:     d.Workday,
		HolidayCalc: d.Holiday,
		Description: d.Description,
		Weekday:     weekday,
	}, nil
}

func (c *Client) Calendar(startDate, endDate time.Time) ([]CalendarDay, error) {

	getCalendarBody := url.Values{}
	getCalendarBody.Add("startDate", startDate.Format(hrFlowDateFormat))
	getCalendarBody.Add("endDate", endDate.Format(hrFlowDateFormat))
	getCalendarRequest, _ := http.NewRequest("GET", "https://hrflow.accountor.fi/KirjaamoWeb/calendar/GetCalendar", strings.NewReader(getCalendarBody.Encode()))
	getCalendarRequest.Header.Add("X-XSRF-TOKEN", c.xsrfToken)
	getCalendarRequest.Header.Add("Accept", "application/json")
	getCalendarRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	resp, err := c.HttpClient.Do(getCalendarRequest)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("http status error %s", resp.Status)
	}

	var response []hrCalendarDay

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, errors.Wrap(err, "decoding calendar response")
	}

	days := []CalendarDay{}
	for _, day := range response {
		calDay, err := day.calendarDay()
		if err != nil {
			return nil, errors.Wrap(err, "converting hrflow day to calendar day")
		}
		days = append(days, calDay)
	}

	return days, nil
}
