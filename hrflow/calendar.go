package hrflow

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) Calendar() {

	getCalendarBody := url.Values{}
	getCalendarBody.Add("startDate", "01.02.2020")
	getCalendarBody.Add("endDate", "29.02.2020")
	getCalendarRequest, _ := http.NewRequest("POST", "https://hrflow.accountor.fi/KirjaamoWeb/calendar/GetCalendar", strings.NewReader(getCalendarBody.Encode()))
	getCalendarRequest.Header.Add("X-XSRF-TOKEN", c.xsrfToken)
	getCalendarRequest.Header.Add("Accept", "application/json")
	getCalendarRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	calResp, err := c.HttpClient.Do(getCalendarRequest)
	if err != nil {
		log.Println(err)
	}
	defer calResp.Body.Close()
	calRespBody, _ := ioutil.ReadAll(calResp.Body)
	calRespString := string(calRespBody)
	fmt.Println(calRespString)
}
