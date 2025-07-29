package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type RootModel struct {
	width  int
	height int
}

func NewRootModel() RootModel {
	return RootModel{}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if IsQuitKey(msg) {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m RootModel) View() string {
	return "Hello, world!"
}
