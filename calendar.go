package main

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func calendarCommandFactory() *cli.Command {

	return &cli.Command{
		Name:   "calendar",
		Action: calendar,
		Usage:  "print the next n (default: 40) days",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "print all dates instead of weekdays only",
			},
			&cli.IntFlag{
				Name:    "count",
				Aliases: []string{"c"},
				Value:   39,
				Usage:   "today + `COUNT` days to fetch",
			},
		},
	}
}

func calendar(c *cli.Context) error {

	allDays := c.Bool("all")
	count := c.Int("count")

	client, err := clientFromConfig()
	if err != nil {
		return errors.Wrap(err, "creating client from config")
	}
	err = client.Authenticate()
	if err != nil {
		return errors.Wrap(err, "authentication failed")
	}

	today := time.Now()
	end := time.Now().AddDate(0, 0, count)

	days, err := client.Calendar(today, end)
	if err != nil {
		return errors.Wrap(err, "getting calendar")
	}
	for _, day := range days {
		if !allDays {
			if day.Weekday == time.Saturday || day.Weekday == time.Sunday {
				continue
			}
		}
		var dayType string
		if day.Workday {
			dayType = "workday"
		} else {
			dayType = "holiday"
		}
		fmt.Println(day.Weekday, day.Date.Format("01.02.2006"), dayType)
	}

	return err
}
