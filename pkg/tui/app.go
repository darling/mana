package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea/v2"
)

func Run() error {
	p := tea.NewProgram(
		NewRootModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		tea.WithInputTTY(),
	)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}
