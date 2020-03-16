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

type WorkLogRow struct {
	// id is always 0.
	Id int64 `json:"id"`
	// customer_id is always 1.
	CustomerID   int64 `json:"customerId"`
	EmploymentID int64 `json:"employmentId"`
	PersonID     int64 `json:"personId"`
	GroupID      int64 `json:"groupId"`
	VersionID    int64 `json:"versionId"`
	// date is the day of the report, in format "2020-02-25 00:00:00.000".
	Date string `json:"date"`
	// mainSalaryFactorID seems to be always 1.
	MainSalaryFactorID int64 `json:"mainSalaryFactorId"`
	// mainAmount is the number of hours that should be reported as a string.
	MainAmount string `json:"mainAmount"`
	// mainUnit is DURATION for monthly workers, and HOURS for hourly workers.
	MainUnit       string          `json:"mainUnit"`
	WorkLogFactors []WorkLogFactor `json:"workLogFactors"`
	StartTime      string          `json:"startTime"`
	EndTime        string          `json:"endTime"`
	// 99002 for DURATION (monthly workers), 11000 for HOURS (hourly workers)
	SalaryGroupValue string `json:"salaryGroupValue"`
	// status is NEW if being created.
	Status string `json:"status"`
	// created date.
	Created                 string `json:"created"`
	LastModifierPersonId    int64  `json:"lastModifierPersonId"`
	LastModifierUserID      int64  `json:"lastModifierUser_Id"`
	LastModifierDisplayName string `json:"lastModifierDisplayName"`
	// modifiedBy is the same as username.
	ModifiedBy string `json:"modifiedBy"`
	// lunchBreak is the duration of the lunch break in minutes.
	LunchBreak int64 `json:"lunchBreak"`
	// cutLunchFromAmount is either Y or N depending on if lunch time should be deducted from total.
	CutLunchFromAmount string `json:"cutLunchFromAmount"`
	// entryText is the same as the comment field in the website.
	EntryText *string `json:"entryText"`
	// entryTextType is probably always TEXTHASHTAG
	EntryTextType     string           `json:"entryTextType"`
	WorkLogRowLinks   []WorkLogRowLink `json:"workLogRowLinks"`
	WorkLogComments   []string         `json:"workLogComments"`
	SourceId          *string          `json:"sourceId"`
	BunchId           *string          `json:"bunchId"`
	ExternalId        *string          `json:"externalId"`
	CreatedFromSource *string          `json:"createdFromSource"`
	StartDate         *string          `json:"startDate"`
	EndDate           *string          `json:"endDate"`
}

type WorkLogFactor struct {
	// id seems to always be 0
	Id int64 `json:"id"`
	// workLogRowID seems to always be 0.
	WorkLogRowID int64 `json:"workLogRowId"`
	// factorID seems to always be 1.
	FactorID int64 `json:"factorId"`
	// amount is the number of hours to be reported.
	Amount float64 `json:"amount"`
	// unit is always DURATION.
	Unit string `json:"unit"`
	// created is the time the entry was created.
	Created string `json:"created"`
	// creator is the same as username
	Creator string `json:"creator"`
	// modifiedBy is the same as username
	ModifiedBy string `json:"modifiedBy"`
	// modified date "2020-02-25 13:40:33.000"
	Modified string `json:"modified"`
}

type WorkLogRowLink struct {
	// always 0
	Id int64 `json:"id"`
	// always 0
	WorkLogRowID int64 `json:"workLogRowId"`
	// 7 for OSASTOT, 8 for KUSTPAIKAT, 9 for PROJEKTIT
	ColNumber int64 `json:"colNumber"`
	// SELECT
	ColControlType string `json:"colControlType"`
	// LIST
	DimensionLinkType string `json:"dimensionLinkType"`
	// PARAM
	DimensionSourceType string `json:"dimensionSourceType"`
	// listID is the list the link is referencing.
	ListID  string  `json:"listId"`
	InputID *string `json:"inputId"`
	Value   *string `json:"value"`
	Label   *string `json:"label"`
}

type WorkLogRequest struct {
	IsFixedProcess                    bool         `json:"isFixedProcess,omitempty"`
	GetListsLabels                    bool         `json:"getListsLabels"`
	GetAll                            bool         `json:"getAll,omitempty"`
	IsSetStartAndEndDateFromCopyDates bool         `json:"isSetStartAndEndDateFromCopyDates"`
	IsUpdateRow                       bool         `json:"isUpdateRow"`
	ViewName                          string       `json:"viewName"`
	Lang                              string       `json:"lang"`
	Employments                       []Employment `json:"employments"`
	StartDate                         string       `json:"startDate"`
	EndDate                           string       `json:"endDate"`
	StatusList                        []string     `json:"statusList"`
	SearchDateType                    string       `json:"searchDateType,omitempty"`
	EmailReceiver                     string       `json:"emailReceiver"`
	// doesn't have to be filled
	EmailReceiverName string `json:"emailReceiverName,omitempty"`
	EmailChangesText  string `json:"emailChangesText,omitempty"`
	EmailComment      string `json:"emailComment,omitempty"`
}

func (c *Client) NewWorkLogRow(
	employmentID, personID, groupID int64,
	startTime, endTime time.Time,
	salaryGroupValue string,
	comment string,
	project *string,
) WorkLogRow {

	hours := endTime.Sub(startTime).Hours()
	date := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())

	now := time.Now().Format(hrFlowTimeFormat)
	return WorkLogRow{
		Id:                 0,
		CustomerID:         1,
		EmploymentID:       employmentID,
		PersonID:           personID,
		GroupID:            groupID,
		VersionID:          3,
		Date:               date.Format(hrFlowTimeFormat),
		MainSalaryFactorID: 1,
		MainAmount:         fmt.Sprintf("%.3f", hours),
		MainUnit:           "DURATION",
		WorkLogFactors:     []WorkLogFactor{c.NewWorkLogFactor(hours)},
		StartTime:          startTime.Format(hrFlowTimeFormat),
		EndTime:            endTime.Format(hrFlowTimeFormat),
		SalaryGroupValue:   salaryGroupValue,
		Status:             "NEW",
		Created:            now,
		ModifiedBy:         c.username,
		LunchBreak:         30,
		CutLunchFromAmount: "Y",
		EntryText:          &comment,
		EntryTextType:      "TEXTHASHTAG",
		WorkLogRowLinks: []WorkLogRowLink{
			c.NewWorkLogRowLink(7, "OSASTOT", nil),
			c.NewWorkLogRowLink(8, "KUSTPAIKAT", nil),
			c.NewWorkLogRowLink(9, "PROJEKTIT", project),
		},
		WorkLogComments: []string{},
	}
}

func (c *Client) NewWorkLogFactor(hours float64) WorkLogFactor {
	now := time.Now().Format(hrFlowTimeFormat)
	return WorkLogFactor{
		Id:           0,
		WorkLogRowID: 0,
		FactorID:     1,
		Amount:       hours,
		Unit:         "DURATION",
		Created:      now,
		Creator:      c.username,
		ModifiedBy:   c.username,
		Modified:     now,
	}
}

func (c *Client) NewWorkLogRowLink(colNumber int64, listID string, project *string) WorkLogRowLink {
	link := WorkLogRowLink{
		Id:                  0,
		WorkLogRowID:        0,
		ColNumber:           colNumber,
		ColControlType:      "SELECT",
		DimensionLinkType:   "LIST",
		DimensionSourceType: "PARAM",
		ListID:              listID,
	}
	if project != nil {
		link.Value = &strings.Split(*project, " ")[0]
		link.Label = project
	}

	return link
}

func (c *Client) NewWorkLogRequest() WorkLogRequest {

	now := time.Now().Format(hrFlowDateFormat)
	return WorkLogRequest{
		ViewName:                          "employee",
		Lang:                              "2",
		IsUpdateRow:                       false,
		Employments:                       c.Employments,
		StartDate:                         now,
		EndDate:                           now,
		GetListsLabels:                    false,
		StatusList:                        []string{"NEW"},
		GetAll:                            true,
		SearchDateType:                    "DATE",
		IsSetStartAndEndDateFromCopyDates: false,
		EmailReceiver:                     c.username,
		EmailReceiverName:                 "",
		EmailChangesText:                  "Save",
		EmailComment:                      "",
		IsFixedProcess:                    true,
	}
}

type newWorkLogResponse struct {
	ActionSuccessful bool `json:"actionSuccessful,omitempty"`
}

func (c *Client) NewWorkLog(startTime, endTime time.Time, salaryGroupValue string, comment string, project *string) error {

	if len(c.Employments) == 0 {
		return errors.New("no employment found, cannot log hours")
	}

	employment := c.Employments[0]

	// Salary group is always the same, but probably shouldn't be hardcoded. No good way to get it right now.
	row := c.NewWorkLogRow(employment.EmploymentID, employment.PersonID, employment.GroupID, startTime, endTime, salaryGroupValue, comment, project)
	rowJSON, err := json.Marshal(row)
	if err != nil {
		return errors.Wrap(err, "marshaling work log row")
	}
	workLogRequest := c.NewWorkLogRequest()
	workLogRequestJSON, err := json.Marshal(workLogRequest)
	if err != nil {
		return errors.Wrap(err, "marshaling work log request")
	}
	body := url.Values{}
	body.Add("workLogRow", string(rowJSON))
	body.Add("workLogRequest", string(workLogRequestJSON))
	// This has to be there
	body.Add("action", `{"Action":"T","ActionId":1999,"Receiver":null,"Label":"Save","SelectedReceiver":null,"Comment":""}`)
	body.Add("copyToDates", "[]")

	req, err := http.NewRequest("POST", "https://hrflow.accountor.fi/KirjaamoWeb/employee/NewWorkLogRow", strings.NewReader(body.Encode()))
	if err != nil {
		return errors.Wrap(err, "creating new work log request")
	}
	req.Header.Add("X-XSRF-TOKEN", c.xsrfToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "work log request")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("http status error %d %s", resp.StatusCode, resp.Status)
	}

	var response newWorkLogResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return errors.Wrap(err, "decoding new work log response")
	}

	if !response.ActionSuccessful {
		return errors.New("backend returned unsuccessful status")
	}

	return nil
}
