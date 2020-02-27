package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Version: "v0.1.0",
		Usage:   "A CLI for HR Flow",
		Before:  checkConfig,
		Commands: []*cli.Command{
			{
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
				},
				Action: report,
			},
		},
		EnableBashCompletion: true,
	}

	_ = app.Run(os.Args)
}
