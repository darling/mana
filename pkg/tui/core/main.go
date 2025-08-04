package core

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/darling/mana/pkg/tui/core/layout"
)

type MainCmp struct {
	focused bool
	width   int
	height  int
	content string
}

func NewMainCmp() MainCmp {
	return MainCmp{}
}

func (m MainCmp) Init() tea.Cmd { return nil }

func (m MainCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newM := m // Copy

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		newM.width = (msg.Width * 3) / 4
		newM.height = msg.Height
	case tea.KeyPressMsg:
		if msg.String() == " " {
			newM.content += "Updated!"
			return newM, nil
		}
		return newM, func() tea.Msg { return layout.UnhandledMsg{Original: msg} }
	}

	return newM, nil
}

func (m MainCmp) View() string {
	prefix := "X"
	if m.focused {
		prefix = "O"
	}
	return fmt.Sprintf("%s Main Content: %s", prefix, m.content)
}

func (m MainCmp) SetFocused(focused bool) (layout.Focusable, tea.Cmd) {
	newM := m
	newM.focused = focused
	return newM, nil
}

func (m MainCmp) IsFocused() bool { return m.focused }

func (m MainCmp) Clone() layout.Focusable {
	return MainCmp{
		focused: m.focused,
		width:   m.width,
		height:  m.height,
		content: m.content,
	}
}
