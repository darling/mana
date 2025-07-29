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

	// Start with content focused
	content.SetFocus(true)

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
			m.sidebar.SetFocus(m.focusIdx == 0)
			m.content.SetFocus(m.focusIdx == 1)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		statusbarHeight := 1
		mainViewHeight := m.height - statusbarHeight

		sidebarWidth := int(float64(m.width) * 0.3)
		contentWidth := m.width - sidebarWidth

		m.sidebar.SetSize(sidebarWidth, mainViewHeight)
		m.content.SetSize(contentWidth, mainViewHeight)
		m.statusbar.SetSize(m.width)
	}

	// Update child models
	m.sidebar, cmd = m.sidebar.Update(msg)
	cmds = append(cmds, cmd)

	m.content, cmd = m.content.Update(msg)
	cmds = append(cmds, cmd)

	m.statusbar, cmd = m.statusbar.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	// Get views from child models (they handle their own styling now)
	sidebarView := m.sidebar.View()
	contentView := m.content.View()

	// Simple horizontal layout - no styling, just positioning
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, contentView)

	// Get statusbar view
	statusbarView := m.statusbar.View()

	// Join main view and statusbar vertically
	return lipgloss.JoinVertical(lipgloss.Left, mainView, statusbarView)
}
