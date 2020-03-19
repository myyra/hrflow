package main

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func reportCommandFactory() *cli.Command {
	return &cli.Command{
		Name:  "report",
		Usage: "add a new hour report",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "duration",
				Aliases: []string{"d"},
				Usage:   "`DURATION` to report, formatted as 8h30m. Will be ignored if both start and end time are defined.",
				Value:   "8h",
			},
			&cli.TimestampFlag{
				Name:        "start",
				Aliases:     []string{"s"},
				Layout:      "15:04",
				Usage:       "Set workday start to `TIME`.",
				DefaultText: "now - DURATION",
			},
			&cli.TimestampFlag{
				Name:        "end",
				Aliases:     []string{"e"},
				Layout:      "15:04",
				Usage:       "Set workday end to `TIME`.",
				DefaultText: "now",
			},
			&cli.StringFlag{
				Name:        "project",
				Aliases:     []string{"p"},
				Usage:       "which `PROJECT` to assign to the report.",
				DefaultText: "none",
			},
			&cli.StringFlag{
				Name:        "comment",
				Aliases:     []string{"c"},
				Usage:       "assign a `COMMENT` to the report.",
				DefaultText: "empty",
			},
			&cli.TimestampFlag{
				Name:        "date",
				Layout:      "2.1.",
				Usage:       "`DATE` for the report, format 'd.M.' (years not supported)",
				DefaultText: "today",
			},
			&cli.BoolFlag{
				Name:  "hourly",
				Value: false,
				Usage: "Report units as an hourly worker, also won't include 30 minute lunch in the duration.",
			},
		},
		Action: report,
	}
}

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

	hourly := c.Bool("hourly")
	var salaryGroupValue string
	if hourly {
		salaryGroupValue = "11000"
	} else {
		salaryGroupValue = "99002"
	}

	comment := c.String("comment")
	p := c.String("project")
	var project *string
	if len(p) != 0 {
		project = &p
	}

	// Lunch is only applicable for monthly workers.
	lunch := !hourly

	client, err := clientFromConfig()
	if err != nil {
		return errors.Wrap(err, "creating client from config")
	}
	err = client.Authenticate()
	if err != nil {
		return errors.Wrap(err, "authentication failed")
	}

	err = client.NewWorkLog(*start, *end, salaryGroupValue, comment, project, lunch)
	if err != nil {
		return errors.Wrap(err, "creating work log")
	}

	fmt.Println("report logged")

	return nil
}
