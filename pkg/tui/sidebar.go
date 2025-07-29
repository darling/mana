package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SidebarModel struct {
	width    int
	height   int
	focused  bool
	selected int
}

func NewSidebarModel() SidebarModel {
	return SidebarModel{
		focused:  false,
		selected: 0,
	}
}

func (m SidebarModel) Init() tea.Cmd {
	return nil
}

func (m SidebarModel) Update(msg tea.Msg) (SidebarModel, tea.Cmd) {
	return m, nil
}

func (m SidebarModel) View() string {
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true)

	if m.focused {
		borderStyle = borderStyle.BorderForeground(lipgloss.Color("39"))
	} else {
		borderStyle = borderStyle.BorderForeground(lipgloss.Color("240"))
	}

	h, v := borderStyle.GetFrameSize()

	innerWidth := m.width - h
	innerHeight := m.height - v

	content := lipgloss.NewStyle().
		Width(innerWidth).
		Height(innerHeight).
		Render("Content")

	return borderStyle.Render(content)
}

func (m *SidebarModel) SetFocus(focused bool) {
	m.focused = focused
}

func (m *SidebarModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
