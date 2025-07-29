package cmd

import (
	"fmt"
)

type VersionCmd struct{}

func (v *VersionCmd) Run() error {
	fmt.Println("mana version 1.0.0")
	return nil
}