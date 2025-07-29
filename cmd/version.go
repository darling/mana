package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func VersionAction(ctx context.Context, cmd *cli.Command) error {
	fmt.Println("mana version 1.0.0")
	return nil
}