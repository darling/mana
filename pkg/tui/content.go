package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ContentModel struct {
	pane *Pane
}

func NewContentModel() ContentModel {
	initialContent := "Welcome to main content!\n\nThis area will display dynamic content."
	return ContentModel{
		pane: NewPane("Chat", initialContent),
	}
}

func (m ContentModel) Init() tea.Cmd { return nil }

func (m ContentModel) Update(msg tea.Msg) (ContentModel, tea.Cmd) {
	if !m.pane.focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			m.pane.SetContent(m.pane.content + "\nNew line added!")
		}
	}

	return m, nil
}

func (m ContentModel) View() string {
	return m.pane.Render()
}

func (m *ContentModel) SetFocus(focused bool) {
	m.pane.SetFocus(focused)
}

func (m *ContentModel) SetSize(width, height int) {
	m.pane.SetSize(width, height)
}
