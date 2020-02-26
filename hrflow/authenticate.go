package hrflow

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

// Authenticate runs the authentication process and stores the needed cookies and header values.
func (c *Client) Authenticate() error {

	start := "https://hrflow.accountor.fi/KirjaamoWeb/login/Employee"

	req, err := http.NewRequest("GET", start, nil)
	if err != nil {
		return errors.Wrap(err, "creating initial request")
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "doing initial request")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(body)

	loginURLRegex := regexp.MustCompile("id=\"options\".*action=\"(.*)\"")
	loginURL := loginURLRegex.FindAllStringSubmatch(bodyString, -1)[0][1]

	loginBody := url.Values{}
	loginBody.Add("UserName", c.username)
	loginBody.Add("Password", c.password)
	loginBody.Add("AuthMethod", "FormsAuthentication")
	loginReq, err := http.NewRequest("POST", loginURL, strings.NewReader(loginBody.Encode()))
	if err != nil {
		return errors.Wrap(err, "creating login request")
	}
	loginReq.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	loginResp, err := c.HttpClient.Do(loginReq)
	if err != nil {
		return errors.Wrap(err, "doing login request")
	}
	defer loginResp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(loginResp.Body)
	if err != nil {
		return errors.Wrap(err, "creating goquery document from login response")
	}

	authURL, err := formActionURL(doc)
	if err != nil {
		return errors.Wrap(err, "getting auth form action")
	}

	authBody, err := formValues(doc)
	if err != nil {
		return errors.Wrap(err, "getting auth form values")
	}

	authReq, err := http.NewRequest("POST", authURL, strings.NewReader(authBody.Encode()))
	if err != nil {
		log.Println(err)
	}
	authReq.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	authResp, err := c.HttpClient.Do(authReq)
	if err != nil {
		log.Println(err)
	}
	defer authResp.Body.Close()

	doc, err = goquery.NewDocumentFromReader(authResp.Body)
	if err != nil {
		return errors.Wrap(err, "creating goquery document from auth response")
	}

	oidcURL, err := formActionURL(doc)
	if err != nil {
		return errors.Wrap(err, "getting auth form action")
	}

	oidcBody, err := formValues(doc)
	if err != nil {
		return errors.Wrap(err, "getting auth form values")
	}

	oidcReq, err := http.NewRequest("POST", oidcURL, strings.NewReader(oidcBody.Encode()))
	if err != nil {
		return errors.Wrap(err, "creating oidc request")
	}
	oidcReq.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	oidcResp, err := c.HttpClient.Do(oidcReq)
	if err != nil {
		return errors.Wrap(err, "doing oidc request")
	}
	defer oidcResp.Body.Close()
	loggedInBody, _ := ioutil.ReadAll(oidcResp.Body)
	loggedInContent := string(loggedInBody)

	rvtRegex := regexp.MustCompile("var RequestVerificationToken = '(.*)'")
	rvt := rvtRegex.FindAllStringSubmatch(loggedInContent, -1)[0][1]

	usrRegex := regexp.MustCompile("var SELECTED_USER_AND_ROLEKEY = \"(.*)\";")
	usr := usrRegex.FindAllStringSubmatch(loggedInContent, -1)[0][1]

	employmentsRegex := regexp.MustCompile(`var employments = (\[{.*}\])`)
	employmentsString := employmentsRegex.FindAllStringSubmatch(loggedInContent, -1)[0][1]
	var employments []Employment
	err = json.Unmarshal([]byte(employmentsString), &employments)
	if err != nil {
		return errors.Wrap(err, "parsing employments")
	}

	c.userRoleKey = usr
	c.xsrfToken = rvt
	c.Employments = employments

	return nil
}

// formActionURL returns the action URL from the first form of the document.
func formActionURL(doc *goquery.Document) (string, error) {

	authForm := doc.Find("form").First()
	authURL, exists := authForm.Attr("action")
	if !exists {
		return "", errors.New("form action attribute not found")
	}
	return authURL, nil
}

// formValues returns the name and value of all input tags from the first form of the document.
func formValues(doc *goquery.Document) (url.Values, error) {

	var err error = nil

	values := url.Values{}
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		name, exists := s.Attr("name")
		if !exists {
			err = errors.New("getting form input tag name")
			return
		}
		value, exists := s.Attr("value")
		if !exists {
			err = errors.New("getting form input tag value")
			return
		}
		values.Add(name, value)
	})

	return values, err
}
