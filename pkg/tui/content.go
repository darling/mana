package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ContentModel struct {
	width   int
	height  int
	focused bool
	content string
}

func NewContentModel() ContentModel {
	return ContentModel{
		focused: true,
		content: "Welcome to the main content area!\n\nThis is where the main content will be displayed.",
	}
}

func (m ContentModel) Init() tea.Cmd {
	return nil
}

func (m ContentModel) Update(msg tea.Msg) (ContentModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.focused {
			return m, nil
		}

		// Handle content-specific key events here
		switch msg.String() {
		case "enter":
			m.content += "\nNew line added!"
		}
	}

	return m, nil
}

func (m ContentModel) View() string {
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

	contentStyle := lipgloss.NewStyle().
		Width(innerWidth).
		Height(innerHeight)

	if m.focused {
		contentStyle = contentStyle.Foreground(lipgloss.Color("15"))
	} else {
		contentStyle = contentStyle.Foreground(lipgloss.Color("245"))
	}

	content := contentStyle.Render(m.content)

	return borderStyle.Render(content)
}

func (m *ContentModel) SetFocus(focused bool) {
	m.focused = focused
}

func (m *ContentModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *ContentModel) SetContent(content string) {
	m.content = content
}
