package main

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func report(c *cli.Context) error {

	now := time.Now()

	start := c.Timestamp("start")
	end := c.Timestamp("end")
	durationFlag := c.String("duration")
	duration, err := time.ParseDuration(durationFlag)
	if err != nil {
		return errors.Wrap(err, "unable to parse duration, check help for formatting")
	}
	date := c.Timestamp("date")

	if date == nil {
		date = &now
	} else {
		d := time.Date(now.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
		date = &d
	}

	if end == nil {
		t := time.Date(date.Year(), date.Month(), date.Day(), now.Hour(), now.Minute(), 0, 0, time.Local)
		end = &t
	} else {
		t := time.Date(date.Year(), date.Month(), date.Day(), end.Hour(), end.Minute(), 0, 0, time.Local)
		end = &t
	}

	if start == nil {
		t := end.Add(-duration)
		start = &t
	} else {
		t := time.Date(date.Year(), date.Month(), date.Day(), start.Hour(), start.Minute(), 0, 0, time.Local)
		start = &t
	}

	comment := c.String("comment")
	p := c.String("project")
	var project *string
	if len(p) != 0 {
		project = &p
	}

	client, err := clientFromConfig()
	if err != nil {
		return errors.Wrap(err, "creating client from config")
	}
	err = client.Authenticate()
	if err != nil {
		return errors.Wrap(err, "authentication failed")
	}

	err = client.NewWorkLog(*start, *end, comment, project)
	if err != nil {
		return errors.Wrap(err, "creating work log")
	}

	fmt.Println("report logged")

	return nil
}
