package cmd

import (
	"github.com/darling/mana/pkg/tui"
)

type DefaultCmd struct {
	Bar string `cmd:"" help:"A value for the default command."`
}

func (d *DefaultCmd) Run() error {
	return tui.Run()
}


type CLI struct {
	Default DefaultCmd `cmd:"" name:"__defaultRunCmd" hidden:"" help:"The default command." default:"withargs"`
	Version VersionCmd `cmd:"" help:"Show the current version."`
}
