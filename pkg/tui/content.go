package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ContentModel struct {
	pane *PaneModel
}

func NewContentModel() ContentModel {
	initialContent := "Welcome to main content!\n\nThis area will display dynamic content."
	return ContentModel{
		pane: NewPane("content", "Chat", initialContent),
	}
}

func (m ContentModel) Init() tea.Cmd { return nil }

func (m ContentModel) Update(msg tea.Msg) (ContentModel, tea.Cmd) {
	var cmd tea.Cmd
	
	// Forward messages to the pane model
	updatedPane, paneCmd := m.pane.Update(msg)
	m.pane = updatedPane.(*PaneModel)
	
	if !m.pane.focused {
		return m, paneCmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			// Use message-based approach for content updates
			newContent := m.pane.content + "\nNew line added!"
			cmd = func() tea.Msg { return PaneContentMsg(newContent) }
		}
	}

	return m, tea.Batch(paneCmd, cmd)
}

func (m ContentModel) View() string {
	return m.pane.View()
}

