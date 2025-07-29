package tui

import (
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
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
	style := lipgloss.NewStyle()
	
	if m.focused {
		style = style.Foreground(lipgloss.Color("15"))
	} else {
		style = style.Foreground(lipgloss.Color("245"))
	}
	
	return style.Render(m.content)
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