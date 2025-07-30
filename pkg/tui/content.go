package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ContentModel struct {
	Pane
	history []string
}

func NewContentModel() ContentModel {
	initialContent := "Welcome to main content!\n\nThis area will display dynamic content."
	return ContentModel{
		Pane:    NewPane("Chat", initialContent),
		history: []string{},
	}
}

func (m ContentModel) Update(msg tea.Msg) (ContentModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" && m.focused {
			m.history = append(m.history, "User input "+time.Now().String())
			content := strings.Join(m.history, "\n")
			m.Pane.content = content
			m.Pane.viewport.SetContent(content)
		}
	}
	var cmd tea.Cmd
	m.Pane, cmd = (&m.Pane).Update(msg)
	return m, cmd
}

func (m ContentModel) View() string {
	return m.Pane.View()
}
