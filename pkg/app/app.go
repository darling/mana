package app

import (
	"github.com/urfave/cli/v3"

	"github.com/darling/mana/cmd"
)

func New() *cli.Command {
	return &cli.Command{
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
}