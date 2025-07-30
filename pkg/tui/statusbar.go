package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatusbarModel struct {
	width          int
	helpMessage    string
	rightComponent string
}

func NewStatusbarModel() StatusbarModel {
	return StatusbarModel{
		helpMessage: "tab: switch focus • q: quit • ↑/↓: navigate • enter: select",
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
	// Define styles
	baseStyle := lipgloss.NewStyle().
		Background(StatusBarBg()).
		Foreground(StatusBarFg()).
		Height(1)

	leftStyle := baseStyle.
		Align(lipgloss.Left).
		Padding(0, 1)

	rightStyle := baseStyle.
		Align(lipgloss.Right).
		Padding(0, 1)

	// Calculate widths - right component gets what it needs, left gets the rest
	rightWidth := max(len(m.rightComponent)+2, 2)
	leftWidth := m.width - rightWidth

	// Create content
	leftContent := leftStyle.Width(leftWidth).Render(m.helpMessage)
	rightContent := rightStyle.Width(rightWidth).Render(m.rightComponent)

	// Join sections horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, leftContent, rightContent)
}

func (m *StatusbarModel) SetHelpMessage(message string) {
	m.helpMessage = message
}

func (m *StatusbarModel) SetSize(width int) {
	m.width = width
}

func (m *StatusbarModel) SetRightComponent(component string) {
	m.rightComponent = component
}

