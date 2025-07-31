package app

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/darling/mana/cmd"
	"github.com/darling/mana/pkg/version"
)

func New(buildInfo version.BuildInfo) *cli.Command {
	return &cli.Command{
		Name:   "mana",
		Usage:  "The cutest LLM interface for your terminal",
		Action: cmd.DefaultAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Show the current version",
			},
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			if c.Bool("version") {
				if err := cmd.NewVersionAction(buildInfo)(ctx, c); err != nil {
					return ctx, err
				}
				return ctx, cli.Exit("", 0)
			}
			return ctx, nil
		},
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Show the current version",
				Action:  cmd.NewVersionAction(buildInfo),
			},
		},
	}
}
