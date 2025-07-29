package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/darling/mana/pkg/tui"
)

func DefaultAction(ctx context.Context, cmd *cli.Command) error {
	return tui.Run()
}
