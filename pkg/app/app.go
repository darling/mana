package app

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/darling/mana/cmd"
	"github.com/darling/mana/pkg/llm"
	_ "github.com/darling/mana/pkg/llm/providers/openrouter"
	"github.com/darling/mana/pkg/tui"
	"github.com/darling/mana/pkg/version"
)

func New(buildInfo version.BuildInfo) *cli.Command {
	var (
		openRouterAPIKey string
		llmManager       *llm.Manager
	)

	return &cli.Command{
		Name:    "mana",
		Usage:   "The cutest LLM interface for your terminal",
		Version: buildInfo.GetVersion(),
		Action: func(ctx context.Context, c *cli.Command) error {
			return tui.Run(llmManager)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "openrouter-api-key",
				Usage:       "OpenRouter API key",
				Destination: &openRouterAPIKey,
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("OPENROUTER_API_KEY"),
				),
			},
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			// Initialize LLM manager if API key is provided
			if openRouterAPIKey != "" {
				manager, err := llm.NewManager("openrouter", llm.Config{
					APIKey: openRouterAPIKey,
					Model:  "qwen/qwen3-coder:turbo",
				})
				if err != nil {
					return ctx, err
				}
				llmManager = manager
			}

			return ctx, nil
		},
		After: func(ctx context.Context, c *cli.Command) error {
			if llmManager != nil {
				return llmManager.Close()
			}
			return nil
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
