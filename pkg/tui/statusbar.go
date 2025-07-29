package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatusbarModel struct {
	width   int
	message string
}

func NewStatusbarModel() StatusbarModel {
	return StatusbarModel{
		message: "tab: switch focus • q: quit • ↑/↓: navigate • enter: select",
	}
}

func (m StatusbarModel) Init() tea.Cmd {
	return nil
}

func (m StatusbarModel) Update(msg tea.Msg) (StatusbarModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	return m, nil
}

func (m StatusbarModel) View() string {
	style := lipgloss.NewStyle().
		Width(m.width).
		Foreground(lipgloss.Color("240")).
		Padding(0, 1)

	return style.Render(m.message)
}

func (m *StatusbarModel) SetMessage(message string) {
	m.message = message
}

func (m *StatusbarModel) SetSize(width int) {
	m.width = width
}