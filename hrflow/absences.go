package hrflow

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Absence is a convenience type that is built from HR Flow's own absence info.
type Absence struct {
	// ID is the HR Flow id of the absence.
	ID        string
	StartDate time.Time
	EndDate   time.Time
	// Pretty much the type of the absence.
	Info string
}

type absenceResponse struct {
	AbsenceInfos *[]absenceInfo `json:"absencesInfos,omitempty"`
	Error        *string        `json:"error,omitempty"`
}

type absenceInfo struct {
	ID string `json:"id,omitempty"`
	// StartDate is the start date of the vacation. End date is not directly available.
	StartDate    string `json:"startDate,omitempty"`
	EmploymentID int64  `json:"employmentId,omitempty"`
	PersonID     int64  `json:"personId,omitempty"`
	// AbsenceInfoText contains the information and full date range of the vacation.
	// Format "03.08.2020 - 11.08.2020 Annual leave".
	AbsenceInfoText string `json:"absenceInfoText,omitempty"`
}

type absenceRequestBody struct {
	// Lang is language apparently, request will fail without it. 2 is English.
	Lang           int64           `json:"lang,omitempty"`
	Employments    []Employment    `json:"employments,omitempty"`
	StartDate      Date            `json:"startDate,omitempty"`
	EndDate        Date            `json:"endDate,omitempty"`
	StatusList     []WorkLogStatus `json:"statusList,omitempty"`
	SearchDateType string          `json:"searchDateType,omitempty"`
	GetAll         bool            `json:"getAll,omitempty"`
	GetListsLabels bool            `json:"getListsLabels,omitempty"`
}

func (c *Client) Absences(employments []Employment, startDate, endDate time.Time) ([]Absence, error) {

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
	request, _ := http.NewRequest("POST", "https://hrflow.accountor.fi/KirjaamoWeb/employee/GetAbsences", strings.NewReader(requestBody.Encode()))
	request.Header.Add("X-XSRF-TOKEN", c.xsrfToken)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	// Adding an Accept header makes the API return double encoded JSON, removing it returns proper data.

	resp, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "doing request")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("http status error %s", resp.Status)
	}

	var absenceResponse absenceResponse

	err = json.NewDecoder(resp.Body).Decode(&absenceResponse)
	if err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}
	if absenceResponse.Error != nil {
		return nil, errors.Errorf("hrflow error: %s", *absenceResponse.Error)
	}
	fmt.Println(absenceResponse.AbsenceInfos)

	absences, err := convertAbsences(*absenceResponse.AbsenceInfos)
	if err != nil {
		return nil, errors.Wrap(err, "converting absences")
	}

	return absences, nil
}

func convertAbsences(hrFlowAbsences []absenceInfo) ([]Absence, error) {

	dateRegex := regexp.MustCompile(`(.*) - (\S*) (.*)`)

	absences := make([]Absence, len(hrFlowAbsences))

	for i, hrFlowAbsence := range hrFlowAbsences {

		matches := dateRegex.FindAllStringSubmatch(hrFlowAbsence.AbsenceInfoText, 1)[0]
		startDate, err := time.Parse(hrFlowDateFormat, matches[1])
		if err != nil {
			return nil, errors.Wrapf(err, "parsing start time %s", matches[1])
		}
		endDate, err := time.Parse(hrFlowDateFormat, matches[2])
		if err != nil {
			return nil, errors.Wrapf(err, "parsing end time %s", matches[2])
		}
		info := matches[3]

		absence := Absence{ID: hrFlowAbsence.ID, StartDate: startDate, EndDate: endDate, Info: info}
		absences[i] = absence
	}

	return absences, nil
}
