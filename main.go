package main

import (
	"context"
	"log"
	"os"

	"github.com/darling/mana/pkg/app"
	"github.com/darling/mana/pkg/version"
)

var (
	versionStr = "dev"
	commitStr  = "none"
	dateStr    = "unknown"
)

func main() {
	buildInfo := version.BuildInfo{
		Version: versionStr,
		Commit:  commitStr,
		Date:    dateStr,
	}

	if err := app.New(buildInfo).Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
