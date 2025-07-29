package cmd

import (
	"fmt"
)

type DefaultCmd struct {
	Bar string `cmd:"" help:"A value for the default command."`
}

func (d *DefaultCmd) Run() error {
	fmt.Println("Hello, world!")
	fmt.Println("Use 'mana --help' for more information")
	return nil
}

type VersionCmd struct{}

func (v *VersionCmd) Run() error {
	fmt.Println("mana version 1.0.0")
	return nil
}

type CLI struct {
	Default DefaultCmd `cmd:"" name:"__defaultRunCmd" hidden:"" help:"The default command." default:"withargs"`
	Version VersionCmd `cmd:"" help:"Show the current version."`
}
