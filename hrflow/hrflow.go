package hrflow

import (
	"net/http"
	"net/http/cookiejar"
)

type Client struct {
	username string
	password string

	xsrfToken   string
	userRoleKey string

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
