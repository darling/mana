package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Tab   key.Binding
	Quit  key.Binding
	Enter key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch focus"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
}

type RootModel struct {
	width     int
	height    int
	sidebar   SidebarModel
	content   ContentModel
	statusbar StatusbarModel
	focusIdx  int // 0 for sidebar, 1 for content
	keys      keyMap
}

func NewRootModel() RootModel {
	sidebar := NewSidebarModel()
	content := NewContentModel()
	statusbar := NewStatusbarModel()

	// Start with content focused.
	// Since this is initialization, we can set state directly.
	content.pane.focused = true

	return RootModel{
		sidebar:   sidebar,
		content:   content,
		statusbar: statusbar,
		focusIdx:  1, // 1 for content focus
		keys:      keys,
	}
}

func (m RootModel) Init() tea.Cmd {
	return tea.Batch(m.sidebar.Init(), m.content.Init(), m.statusbar.Init())
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Tab):
			m.focusIdx = (m.focusIdx + 1) % 2
			var targetID string
			if m.focusIdx == 0 {
				targetID = "sidebar"
			} else {
				targetID = "content"
			}
			// Send a targeted message to focus the correct pane.
			return m, func() tea.Msg { return FocusPaneMsg(targetID) }
		}

	// This case allows the RootModel to track which pane is focused,
	// even when focus is changed by a mouse click.
	case FocusPaneMsg:
		if string(msg) == "sidebar" {
			m.focusIdx = 0
		} else {
			m.focusIdx = 1
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		statusbarHeight := 1
		mainViewHeight := m.height - statusbarHeight

		sidebarWidth := int(float64(m.width) * 0.3)
		contentWidth := m.width - sidebarWidth

		// Update statusbar directly
		m.statusbar.SetSize(m.width)

		// Send specific size messages to the child panes.
		m.sidebar, cmd = m.sidebar.Update(PaneSizeMsg{Width: sidebarWidth, Height: mainViewHeight})
		cmds = append(cmds, cmd)

		m.content, cmd = m.content.Update(PaneSizeMsg{Width: contentWidth, Height: mainViewHeight})
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
	}

	// For all other messages, propagate them to children.
	m.sidebar, cmd = m.sidebar.Update(msg)
	cmds = append(cmds, cmd)

	m.content, cmd = m.content.Update(msg)
	cmds = append(cmds, cmd)

	m.statusbar, cmd = m.statusbar.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	sidebarView := m.sidebar.View()
	contentView := m.content.View()
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, contentView)
	statusbarView := m.statusbar.View()
	return lipgloss.JoinVertical(lipgloss.Left, mainView, statusbarView)
}
