package main

import (
	"github.com/alecthomas/kong"

	"github.com/darling/mana/cmd"
)

func main() {
	var cli cmd.CLI

	ctx := kong.Parse(
		&cli,
		kong.Name("mana"),
		kong.Description("An llm client for the terminal"),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
	)

	err := ctx.Run()

	ctx.FatalIfErrorf(err)
}
