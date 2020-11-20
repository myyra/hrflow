package hrflow

type Employment struct {
	IsPassive           bool  `json:"isPassive"`
	IsDefaultEmployment bool  `json:"isDefaultEmployment"`
	EmploymentID        int64 `json:"employmentId"`
	PersonID            int64 `json:"personId"`
	// Name of the employment, contains some ID, name of the person etc.
	Name string `json:"name"`
	// ListName is the name of the person to be used when listing employees, Lastname Firstname style.
	ListName      string  `json:"listName"`
	GroupID       int64   `json:"groupId"`
	GroupName     string  `json:"groupName"`
	ParentGroupID int64   `json:"parentGroupId"`
	CustomerID    int64   `json:"customerId"`
	UserName      string  `json:"userName"`
	StartDate     string  `json:"startDate"`
	EndDate       *string `json:"endDate"`
	// EnterpriseName is the official name of the company.
	EnterpriseName         string   `json:"enterpriseName"`
	AllEmploymentIDs       *[]int64 `json:"allEmploymentIds"`
	ValueSettings          string   `json:"valueSettings"`
	OrganizationPositionID int64    `json:"organizationPositionId"`
}
