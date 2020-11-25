package hrflow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type WorkAmount struct {
	Date         time.Time
	Amount       time.Duration
	EntryCount   int64
	LunchBreak   time.Duration
	NextStarTime time.Time
}

type hrFlowWorkAmount struct {
	EmploymentID  int64    `json:"employmentID,omitempty"`
	GroupID       int64    `json:"groupID,omitempty"`
	Date          JSONDate `json:"date,omitempty"`
	Hours         float64  `json:"hours,omitempty"`
	HoursCount    int64    `json:"hoursCount,omitempty"`
	LunchBreak    int64    `json:"lunchBreak,omitempty"`
	NextStartTime JSONDate `json:"nextStartTime,omitempty"`
}

type workAmountResponse struct {
	DailyWorkAmountList *[]hrFlowWorkAmount `json:"dailyWorkAmountList,omitempty"`
	Error               *string             `json:"error,omitempty"`
}

func (c *Client) DailyWorkAmount(employments []Employment, startDate, endDate time.Time) ([]WorkAmount, error) {

	body := absenceRequestBody{
		Lang:           2,
		Employments:    employments,
		StartDate:      Date(startDate),
		EndDate:        Date(endDate),
		GetListsLabels: false,
		StatusList: []WorkLogStatus{
			WorkLogStatusNew,
			WorkLogStatusSent,
			WorkLogStatusRejected,
			WorkLogStatusTransferred,
			WorkLogStatusWaitingApproval,
		},
		GetAll:         true,
		SearchDateType: "DATE",
	}
	workLog, _ := json.Marshal(body)

	requestBody := url.Values{}
	requestBody.Add("workLogRequest", string(workLog))
	request, _ := http.NewRequest("POST", "https://hrflow.accountor.fi/KirjaamoWeb/employee/GetDailyWorkAmount", strings.NewReader(requestBody.Encode()))
	request.Header.Add("X-XSRF-TOKEN", c.xsrfToken)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	resp, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "doing request")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("http status error %s", resp.Status)
	}

	var workAmountResponse workAmountResponse

	err = json.NewDecoder(resp.Body).Decode(&workAmountResponse)
	if err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}
	if workAmountResponse.Error != nil {
		return nil, errors.Errorf("hrflow error: %s", *workAmountResponse.Error)
	}

	workAmounts := convertWorkAmount(*workAmountResponse.DailyWorkAmountList)

	return workAmounts, nil
}

func convertWorkAmount(hrFlowWorkAmounts []hrFlowWorkAmount) []WorkAmount {

	workAmounts := make([]WorkAmount, len(hrFlowWorkAmounts))

	for i, hrFlowWorkAmount := range hrFlowWorkAmounts {

		amount := time.Duration(hrFlowWorkAmount.Hours * float64(time.Hour))
		lunchBreak := time.Duration(hrFlowWorkAmount.LunchBreak) * time.Minute

		workAmount := WorkAmount{
			Date:         time.Time(hrFlowWorkAmount.Date),
			Amount:       amount,
			EntryCount:   hrFlowWorkAmount.HoursCount,
			LunchBreak:   lunchBreak,
			NextStarTime: time.Time(hrFlowWorkAmount.NextStartTime),
		}

		workAmounts[i] = workAmount
	}

	return workAmounts
}
