package hrflow

import (
	"net/http"
	"net/http/cookiejar"
)

const (
	hrFlowTimeFormat = "2006-01-02 15:04:05.000"
	hrFlowDateFormat = "02.01.2006"
)

type Client struct {
	username string
	password string

	xsrfToken   string
	userRoleKey string
	Employments []Employment

	HttpClient *http.Client
}

func NewClient(username, password string) *Client {

	cookieJar, _ := cookiejar.New(nil)
	return &Client{
		username: username,
		password: password,
		HttpClient: &http.Client{
			Jar: cookieJar,
		},
	}
}
