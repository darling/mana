package tui

import (
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
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
	focusIdx  int
	keys      keyMap
}

func NewRootModel() RootModel {
	sidebar := NewSidebarModel()
	content := NewContentModel()
	statusbar := NewStatusbarModel()

	sidebar.Pane.focused = false
	content.Pane.focused = true

	return RootModel{
		sidebar:   sidebar,
		content:   content,
		statusbar: statusbar,
		focusIdx:  1,
		keys:      keys,
	}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Tab):

			m.focusIdx = (m.focusIdx + 1) % 2

			m.sidebar.Pane.focused = m.focusIdx == 0
			m.content.Pane.focused = m.focusIdx == 1

			return m, nil

		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		statusbarHeight := 1
		mainViewHeight := m.height - statusbarHeight
		sidebarWidth := int(float64(m.width) * 0.3)
		if sidebarWidth < 20 {
			sidebarWidth = 20
		}
		contentWidth := m.width - sidebarWidth

		m.statusbar.SetSize(m.width)

		m.sidebar.Pane.width = sidebarWidth
		m.sidebar.Pane.height = mainViewHeight
		m.sidebar.Pane.handleResize(sidebarWidth, mainViewHeight)

		m.content.Pane.width = contentWidth
		m.content.Pane.height = mainViewHeight
		m.content.Pane.handleResize(contentWidth, mainViewHeight)

		return m, nil
	}

	var sbCmd, ctCmd, stCmd tea.Cmd

	if m.focusIdx == 0 {
		m.sidebar, sbCmd = m.sidebar.Update(msg)
		cmds = append(cmds, sbCmd)
	} else {
		m.content, ctCmd = m.content.Update(msg)
		cmds = append(cmds, ctCmd)
	}

	m.statusbar, stCmd = m.statusbar.Update(msg)
	cmds = append(cmds, stCmd)

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	sidebarView := m.sidebar.View()
	contentView := m.content.View()

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, contentView)

	statusbarView := m.statusbar.View()

	return lipgloss.JoinVertical(lipgloss.Left, mainView, statusbarView)
}
