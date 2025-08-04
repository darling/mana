package core

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/darling/mana/pkg/tui/core/layout"
)

type SidebarItemCmp struct {
	focused bool
	label   string
}

func NewSidebarItemCmp(label string) SidebarItemCmp {
	return SidebarItemCmp{label: label}
}

func (i SidebarItemCmp) Init() tea.Cmd { return nil }

func (i SidebarItemCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return i, nil
}

func (i SidebarItemCmp) View() string {
	if i.focused {
		return fmt.Sprintf("[ %s ]", i.label)
	}
	return fmt.Sprintf("  %s  ", i.label)
}

func (i SidebarItemCmp) SetFocused(focused bool) (layout.Focusable, tea.Cmd) {
	newI := i
	newI.focused = focused
	return newI, nil
}

func (i SidebarItemCmp) IsFocused() bool {
	return i.focused
}

func (i SidebarItemCmp) Clone() layout.Focusable {
	return SidebarItemCmp{
		focused: i.focused,
		label:   i.label,
	}
}
