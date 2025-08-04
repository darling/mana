package core

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/darling/mana/pkg/tui/core/layout"
)

type SidebarCmp struct {
	focused bool

	width  int
	height int
}

func NewSidebarCmp() *SidebarCmp {
	return &SidebarCmp{
		focused: false,
	}
}

func (s SidebarCmp) Init() tea.Cmd {
	return nil
}

func (s SidebarCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width / 4 // 25% of screen width
		s.height = msg.Height
	}
	return s, nil
}

func (s SidebarCmp) View() string {
	focusIndicator := "X"
	if s.focused {
		focusIndicator = "O"
	}
	return focusIndicator + "Sidebar (placeholder)"
}

func (s SidebarCmp) SetFocused(focused bool) (layout.Focusable, tea.Cmd) {
	newS := s
	newS.focused = focused
	return newS, nil
}

func (s SidebarCmp) IsFocused() bool {
	return s.focused
}

func (m SidebarCmp) Clone() layout.Focusable {
	return SidebarCmp{
		focused: m.focused,
		width:   m.width,
		height:  m.height,
	}
}
