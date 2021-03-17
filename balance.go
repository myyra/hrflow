package main

import (
	"fmt"
	"time"

	"github.com/myyra/hrflow/hrflow"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func balanceCommandFactory() *cli.Command {
	return &cli.Command{
		Name:  "balance",
		Usage: "check your working hours balance",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "months",
				Aliases: []string{"m"},
				Value:   30,
				Usage:   "how many `MONTHS` to check backwards",
			},
		},
		Action: balance,
	}
}

func balance(c *cli.Context) error {

	months := c.Int("months")

	client, err := clientFromConfig()
	if err != nil {
		return errors.Wrap(err, "creating client from config")
	}
	err = client.Authenticate()
	if err != nil {
		return errors.Wrap(err, "authentication failed")
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, -months, 0)

	calendarDays, err := client.Calendar(startDate, endDate)
	if err != nil {
		return errors.Wrap(err, "getting calendar")
	}

	absences, err := client.Absences(client.Employments, startDate, endDate)
	if err != nil {
		return errors.Wrap(err, "getting absences")
	}

	absenceMap := make(map[time.Time]bool)
	for _, absence := range absences {
		for _, day := range daysInAbsence(absence) {
			absenceMap[day] = true
		}
	}

	workAmounts, err := client.DailyWorkAmount(client.Employments, startDate, endDate)
	if err != nil {
		return errors.Wrap(err, "getting daily work amount")
	}

	workAmountsMap := make(map[time.Time]float64)
	for _, workAmount := range workAmounts {
		workAmountsMap[workAmount.Date] = workAmount.Amount.Hours()
	}

	neededHours := 0.0
	actualHours := 0.0
	missingDays := make([]time.Time, 0)

	for _, day := range calendarDays {
		if day.Workday {
			neededHours += 7.5
			if _, ok := absenceMap[day.Date]; ok {
				actualHours += 7.5
			} else {
				if hours, ok := workAmountsMap[day.Date]; ok {
					actualHours += hours
				} else {
					actualHours += 7.5
					missingDays = append(missingDays, day.Date)
				}
			}

		}
	}

	fmt.Println("Needed hours:", neededHours)
	fmt.Println("Actual hours:", actualHours)
	fmt.Println("Missing days:", missingDays)

	return nil
}

func daysInAbsence(absence hrflow.Absence) (dates []time.Time) {

	for date := absence.StartDate; date.Before(absence.EndDate.AddDate(0, 0, 1)); date = date.AddDate(0, 0, 1) {
		dates = append(dates, date)
	}

	return dates
}
