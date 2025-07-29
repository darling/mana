package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/darling/mana/cmd"
)

func main() {
	app := &cli.Command{
		Name:        "mana",
		Usage:       "An llm client for the terminal",
		Action:      cmd.DefaultAction,
		Commands: []*cli.Command{
			{
				Name:   "version",
				Usage:  "Show the current version",
				Action: cmd.VersionAction,
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
