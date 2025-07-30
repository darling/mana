package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SidebarModel struct {
	Pane
	items []string
}

func NewSidebarModel() SidebarModel {
	return SidebarModel{
		Pane:  NewPane("Chats", "Navigation\n\n- Item 1\n- Item 2\n- Item 3"),
		items: []string{"Chat 1", "Chat 2", "Chat 3"},
	}
}

func (m SidebarModel) Update(msg tea.Msg) (SidebarModel, tea.Cmd) {
	var cmd tea.Cmd
	m.Pane, cmd = (&m.Pane).Update(msg)
	return m, cmd
}

func (m SidebarModel) View() string {
	return m.Pane.View()
}

