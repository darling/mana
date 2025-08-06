package core

import (
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/darling/mana/pkg/tui/core/layout"
)

type SidebarCmp struct {
	focusManager layout.FocusManager

	keys sidebarKeyMap

	focused bool

	width  int
	height int
}

func NewSidebarCmp() *SidebarCmp {
	items := []layout.Focusable{
		NewSidebarPaneCmp("Conversations"),
		NewSidebarPaneCmp("Models"),
		NewSidebarPaneCmp("Settings"),
	}

	fm := layout.NewFocusManager(items, false)
	// Focus the first pane by default within the sidebar
	fm, _ = fm.FocusNext()

	return &SidebarCmp{
		focusManager: fm,
		focused:      false,
		keys:         DefaultSidebarKeyMap,
	}
}

func (s SidebarCmp) Init() tea.Cmd {
	return nil
}

func (s SidebarCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case layout.ComponentSizeMsg:
		s.width = msg.Width / 4
		s.height = msg.Height

		// Distribute height among the three panes
		paneHeight := s.height / 3
		remainingHeight := s.height % 3

		panes := s.focusManager.GetAll()
		for i, pane := range panes {
			h := paneHeight
			if i < remainingHeight {
				h++ // Distribute remainder to top panes
			}
			paneSize := layout.ComponentSizeMsg{Width: s.width, Height: h}

			updatedPaneModel, updateCmd := pane.Update(paneSize)
			if updatedPane, ok := updatedPaneModel.(layout.Focusable); ok {
				s.focusManager, _ = s.focusManager.Set(i, updatedPane)
			}
			cmds = append(cmds, updateCmd)
		}

	case tea.KeyPressMsg:
		// Only handle navigation keys if the sidebar itself is the focused component.
		if !s.focused {
			return s, nil
		}

		switch {
		case key.Matches(msg, s.keys.FocusDown):
			s.focusManager, cmd = s.focusManager.FocusNext()
		case key.Matches(msg, s.keys.FocusUp):
			s.focusManager, cmd = s.focusManager.FocusPrev()
		default:
			s.focusManager, cmd = s.focusManager.UpdateFocused(msg)
		}
		cmds = append(cmds, cmd)

	default:
		// Pass other messages to all children, as they might be relevant regardless of focus
		s.focusManager, cmd = s.focusManager.UpdateAll(msg)
		cmds = append(cmds, cmd)
	}

	return s, tea.Batch(cmds...)
}

func (s SidebarCmp) View() string {
	panes := s.focusManager.GetAll()
	viewedPanes := make([]string, len(panes))

	for i, pane := range panes {
		p := pane.Clone()

		// If the sidebar component itself is not focused, then none of its children
		// should appear focused, regardless of their internal state.
		if !s.focused {
			p, _ = p.SetFocused(false)
		}

		viewedPanes[i] = p.View()
	}

	// The sidebar is now just a vertical container for the panes, with no border of its own.
	return lipgloss.JoinVertical(lipgloss.Left, viewedPanes...)
}

func (s SidebarCmp) SetFocused(focused bool) (layout.Focusable, tea.Cmd) {
	s.focused = focused
	// The visual change is handled by the View() method, so no command is needed here.
	return s, nil
}

func (s SidebarCmp) IsFocused() bool {
	return s.focused
}

func (m SidebarCmp) Clone() layout.Focusable {
	return SidebarCmp{
		focusManager: m.focusManager.Clone(),
		keys:         m.keys,

		focused: m.focused,

		width:  m.width,
		height: m.height,
	}
}
