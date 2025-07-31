package cmd

import (
	"context"
	"fmt"

	"github.com/darling/mana/pkg/version"
	"github.com/urfave/cli/v3"
)

func NewVersionAction(buildInfo version.BuildInfo) func(context.Context, *cli.Command) error {
	return func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println(buildInfo.String())
		return nil
	}
}
