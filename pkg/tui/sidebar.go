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
	items    []string
}

func NewSidebarModel() SidebarModel {
	return SidebarModel{
		focused:  false,
		selected: 0,
		items:    []string{"Item 1", "Item 2", "Item 3"},
	}
}

func (m SidebarModel) Init() tea.Cmd {
	return nil
}

func (m SidebarModel) Update(msg tea.Msg) (SidebarModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.focused {
			return m, nil
		}

		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.items)-1 {
				m.selected++
			}
		}
	}

	return m, nil
}

func (m SidebarModel) View() string {
	if len(m.items) == 0 {
		return "No items"
	}

	var items []string
	for i, item := range m.items {
		style := lipgloss.NewStyle()

		if m.selected == i {
			// Highlight selected item
			style = style.
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("24")).
				Bold(true)
			items = append(items, style.Render("> "+item))
		} else {
			// Regular item
			if m.focused {
				style = style.Foreground(lipgloss.Color("15"))
			} else {
				style = style.Foreground(lipgloss.Color("245"))
			}
			items = append(items, style.Render("  "+item))
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, items...)

	// Apply border styling based on focus state
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		Width(m.width).
		Height(m.height)

	if m.focused {
		borderStyle = borderStyle.BorderForeground(lipgloss.Color("39"))
	} else {
		borderStyle = borderStyle.BorderForeground(lipgloss.Color("240"))
	}

	return borderStyle.Render(content)
}

func (m *SidebarModel) SetFocus(focused bool) {
	m.focused = focused
}

func (m *SidebarModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}
