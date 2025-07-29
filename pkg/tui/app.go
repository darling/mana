package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() error {
	p := tea.NewProgram(NewRootModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}
