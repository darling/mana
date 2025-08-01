package tui

import (
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/textarea"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// PromptSubmittedMsg is sent when the user submits a prompt
type PromptSubmittedMsg struct {
	Prompt string
}

type PromptDialog struct {
	ta     textarea.Model
	keys   keyMap
	submit key.Binding
	cancel key.Binding
}

func NewPromptDialog() PromptDialog {
	ta := textarea.New()
	ta.Placeholder = "Type your message here..."

	ta.SetWidth(60)
	ta.SetHeight(5)
	ta.Focus()

	submit := key.NewBinding(
		key.WithKeys("enter", "shift+enter"),
		key.WithHelp("shift+enter", "send"),
	)
	cancel := key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)

	return PromptDialog{
		ta:     ta,
		keys:   keys,
		submit: submit,
		cancel: cancel,
	}
}

func (m PromptDialog) Init() tea.Cmd {
	return textarea.Blink
}

func (m PromptDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.submit):
			// Submit the prompt
			prompt := m.ta.Value()
			if prompt != "" {
				return m, tea.Batch(
					func() tea.Msg { return PromptSubmittedMsg{Prompt: prompt} },
					func() tea.Msg { return CloseDialogMsg{} },
				)
			}
			return m, func() tea.Msg { return CloseDialogMsg{} }
		case key.Matches(msg, m.cancel):
			// Close without sending
			return m, func() tea.Msg { return CloseDialogMsg{} }
		}
	}

	// Let textarea handle all other messages
	m.ta, cmd = m.ta.Update(msg)
	return m, cmd
}

func (m PromptDialog) View() string {
	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(BorderFocused()).
		Padding(1, 2).
		Width(70)

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(BorderFocused()).
		Render("Compose Message")

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Shift+Enter to send • Esc to cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"\n",
		m.ta.View(),
		"\n",
		help,
	)

	return containerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, title, content),
	)
}
