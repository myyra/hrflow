package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Version: "v0.3.0",
		Usage:   "A CLI for HR Flow",
		Before:  checkConfig,
		Commands: []*cli.Command{
			reportCommandFactory(),
			calendarCommandFactory(),
		},
		EnableBashCompletion: true,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
