package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/darling/mana/pkg/chat"

	"github.com/darling/mana/pkg/tui/core"
)

func Run(service chat.Service) error {
	root := core.NewRootCmp(service)

	p := tea.NewProgram(
		root,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		tea.WithInputTTY(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	return nil
}
