package core

import (
	"fmt"

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
		NewSidebarItemCmp("Conversations"),
		NewSidebarItemCmp("Models"),
		NewSidebarItemCmp("Settings"),
	}

	fm := layout.NewFocusManager(items, true)
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
	case tea.WindowSizeMsg:
		s.width = msg.Width / 4 // 25% of screen width
		s.height = msg.Height
		s.focusManager, cmd = s.focusManager.UpdateAll(msg)
		cmds = append(cmds, cmd)
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, s.keys.FocusDown):
			s.focusManager, cmd = s.focusManager.FocusNext()
		// TODO:
		// case key.Matches(msg, s.keys.FocusUp):
		// 	s.focusManager, cmd = s.focusManager.FocusPrev()
		default:
			s.focusManager, cmd = s.focusManager.UpdateFocused(msg)
			cmds = append(cmds, cmd)
		}
	default:
		s.focusManager, cmd = s.focusManager.UpdateFocused(msg)
		cmds = append(cmds, cmd)
	}

	return s, tea.Batch(cmds...)
}

func (s SidebarCmp) View() string {
	components := s.focusManager.GetAll()

	viewedComponents := make([]string, len(components))
	for i, component := range components {
		viewedComponents[i] = component.View()
	}

	sidebarHeading := fmt.Sprintf("Mana (%d)", len(viewedComponents))

	list := lipgloss.JoinVertical(lipgloss.Left, viewedComponents...)
	return lipgloss.JoinVertical(lipgloss.Left, sidebarHeading, list)
}

func (s SidebarCmp) SetFocused(focused bool) (layout.Focusable, tea.Cmd) {
	s.focused = focused
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
