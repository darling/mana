package tui

import (
	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type ConfigDialog struct {
	focused   bool
	apiKey    textinput.Model
	selection int
	choices   []string
}

func NewConfigDialog() ConfigDialog {
	ti := textinput.New()
	ti.Placeholder = "Enter API key"
	ti.CharLimit = 64
	ti.SetWidth(40)

	return ConfigDialog{
		apiKey:    ti,
		selection: 0,
		choices:   []string{"Save & Close", "Cancel"},
	}
}

func (m ConfigDialog) Init() tea.Cmd { return nil }

func (m ConfigDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case FocusDialogMsg:
		m.focused = msg.Focused
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.selection == 0 { // Save
				// Save logic here
				return m, func() tea.Msg { return CloseDialogMsg{} }
			} else { // Cancel
				return m, func() tea.Msg { return CloseDialogMsg{} }
			}
		case "up":
			m.selection = max(0, m.selection-1)
		case "down":
			m.selection = min(len(m.choices)-1, m.selection+1)
		}
	}

	if m.focused {
		m.apiKey, cmd = m.apiKey.Update(msg)
	}
	return m, cmd
}

func (m ConfigDialog) View() string {
	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(BorderFocused()).
		Padding(1, 2).
		Width(60)

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