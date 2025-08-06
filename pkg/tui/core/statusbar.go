package core

import (
	"strings"

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/darling/mana/pkg/tui/core/components"
	"github.com/darling/mana/pkg/tui/core/layout"
)

type StatusBarCmp struct {
	width    int
	version  string
	bindings []key.Binding
}

func NewStatusBarCmp(version string) components.Component {
	return &StatusBarCmp{
		version: version,
	}
}

func (s StatusBarCmp) Init() tea.Cmd { return nil }

func (s StatusBarCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case layout.ComponentSizeMsg:
		s.width = msg.Width
	case layout.HelpUpdateMsg:
		s.bindings = msg
	}
	return s, nil
}

func (s StatusBarCmp) View() string {
	if s.width == 0 {
		return ""
	}

	// Render help text from key bindings
	var helpParts []string
	for _, b := range s.bindings {
		// Only display bindings that have help text.
		helpParts = append(helpParts, b.Help().Key+" "+b.Help().Desc)
	}
	helpView := lipgloss.NewStyle().Margin(0, 1).Render(strings.Join(helpParts, " • "))

	// Render version text on the right
	versionView := lipgloss.NewStyle().Margin(0, 1).Render(s.version)

	// Calculate space for the version to align it right
	availableWidth := s.width - lipgloss.Width(helpView)
	versionStyle := lipgloss.NewStyle().Width(availableWidth).Align(lipgloss.Right)

	return lipgloss.JoinHorizontal(lipgloss.Top, helpView, versionStyle.Render(versionView))
}
