package main

import (
	"context"
	"log"
	"os"

	"github.com/darling/mana/pkg/app"
)

func main() {
	if err := app.New().Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
