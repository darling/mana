package tui

import (
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type ConfigDialog struct {
	apiKey    textinput.Model
	selection int
	choices   []string
	keys      keyMap
}

func NewConfigDialog() ConfigDialog {
	ti := textinput.New()
	ti.Placeholder = "Enter API key"
	ti.CharLimit = 64
	ti.SetWidth(50)
	ti.Focus()

	return ConfigDialog{
		apiKey:    ti,
		selection: 0,
		choices:   []string{"Save & Close", "Cancel"},
		keys:      keys,
	}
}

func (m ConfigDialog) Init() tea.Cmd { return nil }

func (m ConfigDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			if m.selection == 0 {
				return m, func() tea.Msg { return CloseDialogMsg{} }
			} else {
				return m, func() tea.Msg { return CloseDialogMsg{} }
			}
		case key.Matches(msg, m.keys.Up):
			m.selection = max(0, m.selection-1)
		case key.Matches(msg, m.keys.Down):
			m.selection = min(len(m.choices)-1, m.selection+1)
		}
	}

	m.apiKey, cmd = m.apiKey.Update(msg)
	return m, cmd
}

func (m ConfigDialog) View() string {
	// Create the main dialog content
	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(BorderFocused()).
		Padding(2, 4).
		Width(70)

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(BorderFocused()).
		Render("Settings")

	input := m.apiKey.View()
	buttons := make([]string, len(m.choices))
	for i, choice := range m.choices {
		style := lipgloss.NewStyle()
		if i == m.selection {
			style = style.Bold(true).Underline(true)
		}
		buttons[i] = style.Render("[" + choice + "]")
	}

	buttonRow := lipgloss.JoinHorizontal(lipgloss.Left, buttons...)
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"\n",
		"API Key:",
		input,
		"\n\n",
		buttonRow,
	)

	return containerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, title, content),
	)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

