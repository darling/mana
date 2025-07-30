package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SidebarModel struct {
	pane *PaneModel

	selected int
}

func NewSidebarModel() SidebarModel {
	return SidebarModel{
		pane: NewPane("sidebar", "Chats", "Navigation\n\n- Item 1\n- Item 2\n- Item 3"),
	}
}

func (m SidebarModel) Init() tea.Cmd { return nil }

func (m SidebarModel) Update(msg tea.Msg) (SidebarModel, tea.Cmd) {
	// Forward messages to the pane model
	updatedPane, paneCmd := m.pane.Update(msg)
	m.pane = updatedPane.(*PaneModel)
	
	return m, paneCmd
}

func (m SidebarModel) View() string {
	return m.pane.View()
}

