package core

import (
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/darling/mana/pkg/tui/core/components"
	"github.com/darling/mana/pkg/tui/core/layout"
)

type RootCmp interface {
	components.Component
}

type rootCmp struct {
	statusbar components.Component

	keys keyMap

	focusManager layout.FocusManager

	width, height int
}

func NewRootCmp() RootCmp {
	sidebar := NewSidebarCmp()
	main := NewMainCmp()
	statusbar := NewStatusBarCmp("v0.1.0")

	focusables := []layout.Focusable{sidebar.Clone(), main.Clone()}

	fm := layout.NewFocusManager(focusables, true)
	// Ignore the given command as we just want to set state before the first render
	fm, _ = fm.FocusNext()

	return rootCmp{
		statusbar:    statusbar,
		keys:         DefaultKeyMap,
		focusManager: fm,
	}
}

func (m rootCmp) Init() tea.Cmd {
	return m.getHelpCmd()
}

func (m rootCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.focusManager, cmd = m.focusManager.UpdateAll(layout.ComponentSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - 1, // Account for status bar
		})
		cmds = append(cmds, cmd)
		newStatusBar, cmd := m.statusbar.Update(layout.ComponentSizeMsg{Width: msg.Width})
		m.statusbar = newStatusBar.(components.Component)
		cmds = append(cmds, cmd)

	case layout.FocusChangedMsg:
		return m, m.getHelpCmd()

	case tea.KeyPressMsg:
		m, cmd = m.handleKeyPress(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	// NOTE: Unhandled types get passed to the focused component for now
	default:
		m.focusManager, cmd = m.focusManager.UpdateFocused(msg)
		cmds = append(cmds, cmd)
	}

	// Messages	always go to the status bar
	var sbCmd tea.Cmd
	newStatusBar, sbCmd := m.statusbar.Update(msg)
	m.statusbar = newStatusBar.(components.Component)
	cmds = append(cmds, sbCmd)

	return m, tea.Batch(cmds...)
}

func (m rootCmp) View() string {
	sidebar, err := m.focusManager.Get(0)
	if err != nil {
		return "Error retrieving sidebar: " + err.Error()
	}
	main, err := m.focusManager.Get(1)
	if err != nil {
		return "Error retrieving main view: " + err.Error()
	}

	// First row: sidebar + main
	top := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebar.View(),
		main.View(),
	)

	// Second row: status bar
	return lipgloss.JoinVertical(
		lipgloss.Left,
		top,
		m.statusbar.View(),
	)
}

func (m rootCmp) handleKeyPress(msg tea.KeyPressMsg) (rootCmp, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.FocusNext):
		m.focusManager, cmd = m.focusManager.FocusNext()
		return m, cmd
	default:
		m.focusManager, cmd = m.focusManager.UpdateFocused(msg)
		return m, cmd
	}
}

func (m rootCmp) getHelpCmd() tea.Cmd {
	var bindings []key.Binding

	// Get bindings from the focused component
	focused, err := m.focusManager.GetFocused()
	if err == nil {
		if helpable, ok := focused.(layout.Help); ok {
			bindings = append(bindings, helpable.Bindings()...)
		}
	}

	// Add global key bindings
	bindings = append(bindings, m.keys.FocusNext, m.keys.Quit)

	return func() tea.Msg {
		return layout.HelpUpdateMsg(bindings)
	}
}
