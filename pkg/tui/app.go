package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/darling/mana/pkg/llm"
)

func Run(manager *llm.Manager) error {
	root := NewRootModel(manager)

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
