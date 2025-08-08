package core

import (
	"fmt"

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"

	"github.com/darling/mana/pkg/tui/core/layout"
)

type MainCmp struct {
	focused bool
	width   int
	height  int
	content string
	keys    mainKeyMap
}

func NewMainCmp() MainCmp {
	return MainCmp{
		keys: DefaultMainKeyMap,
	}
}

func (m MainCmp) Init() tea.Cmd { return nil }

func (m MainCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newM := m // Copy

	switch msg := msg.(type) {
	case layout.ComponentSizeMsg:
		sidebarWidth := msg.Width / 4
		newM.width = msg.Width - sidebarWidth
		newM.height = msg.Height
	case layout.ConfirmedMsg:
		newM.content = "Dialog confirmed! Action executed."
	case layout.CancelledMsg:
		newM.content = "Dialog cancelled."
	case tea.KeyPressMsg:
		if !m.focused {
			return newM, nil
		}

		switch {
		case key.Matches(msg, m.keys.Redraw):
			newM.content = "Screen redrawn!"
		case key.Matches(msg, m.keys.Create):
			newM.content = "Creating new prompt..."
		case key.Matches(msg, m.keys.ShowDialog):
			return newM, func() tea.Msg {
				return layout.ShowConfirmDialogMsg{Text: "Are you sure you want to delete this conversation?"}
			}
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
		keys:    m.keys,
	}
}

func (m MainCmp) Bindings() []key.Binding {
	return []key.Binding{m.keys.Redraw, m.keys.Create, m.keys.ShowDialog}
}
