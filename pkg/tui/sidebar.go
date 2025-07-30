package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SidebarModel struct {
	pane *Pane

	selected int
}

func NewSidebarModel() SidebarModel {
	return SidebarModel{
		pane: NewPane("Chats", "Navigation\n\n- Item 1\n- Item 2\n- Item 3"),
	}
}

func (m SidebarModel) Init() tea.Cmd { return nil }

func (m SidebarModel) Update(msg tea.Msg) (SidebarModel, tea.Cmd) {
	return m, nil
}

func (m SidebarModel) View() string {
	return m.pane.Render()
}

func (m *SidebarModel) SetFocus(focused bool) {
	m.pane.SetFocus(focused)
}

func (m *SidebarModel) SetSize(width, height int) {
	m.pane.SetSize(width, height)
}
