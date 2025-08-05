package core

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"

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
	case layout.ComponentSizeMsg:
		sidebarWidth := msg.Width / 4
		newM.width = msg.Width - sidebarWidth
		newM.height = msg.Height
	case tea.KeyPressMsg:
		if msg.Code == tea.KeyEnter {
			newM.content += "Updated!"
		}
	}

	return newM, nil
}

func (m MainCmp) View() string {
	content := fmt.Sprintf("Main Content: %s", m.content)

	var boxStyle lipgloss.Style
	if m.focused {
		boxStyle = FocusedBox
	} else {
		boxStyle = BlurredBox
	}

	return boxStyle.Width(m.width).Height(m.height).Render(content)
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
