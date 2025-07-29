package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Tab   key.Binding
	Help  key.Binding
	Quit  key.Binding
	Enter key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter},
		{k.Tab, k.Help, k.Quit},
	}
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
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
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
	width    int
	height   int
	sidebar  SidebarModel
	content  ContentModel
	focusIdx int // 0 for sidebar, 1 for content
	help     help.Model
	keys     keyMap
}

func NewRootModel() RootModel {
	sidebar := NewSidebarModel()
	content := NewContentModel()

	// Start with content focused
	content.SetFocus(true)

	h := help.New()
	h.ShowAll = false // Start with short help

	return RootModel{
		sidebar:  sidebar,
		content:  content,
		focusIdx: 1, // 1 for content focus
		help:     h,
		keys:     keys,
	}
}

func (m RootModel) Init() tea.Cmd {
	return tea.Batch(m.sidebar.Init(), m.content.Init())
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Tab):
			m.focusIdx = (m.focusIdx + 1) % 2
			m.sidebar.SetFocus(m.focusIdx == 0)
			m.content.SetFocus(m.focusIdx == 1)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update help model width
		m.help.Width = msg.Width

		// Calculate responsive sidebar width
		sidebarWidth := m.width / 3 // Default to 1/3 of width
		contentWidth := m.width - sidebarWidth

		// Reserve space for help bar at bottom
		availableHeight := m.height - 1

		// Account for border space when setting component sizes
		borderPadding := 2 // 1 character border on each side
		m.sidebar.SetSize(sidebarWidth-borderPadding, availableHeight-borderPadding)
		m.content.SetSize(contentWidth-borderPadding, availableHeight-borderPadding)
	}

	// Update child models
	var cmd tea.Cmd
	m.sidebar, cmd = m.sidebar.Update(msg)
	cmds = append(cmds, cmd)

	m.content, cmd = m.content.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	// Get views from child models (they handle their own styling now)
	sidebarView := m.sidebar.View()
	contentView := m.content.View()

	// Simple horizontal layout - no styling, just positioning
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, contentView)

	// Create help view
	helpView := m.help.View(m.keys)
	helpStyle := lipgloss.NewStyle().
		Width(m.width).
		Padding(0, 1)

	styledHelp := helpStyle.Render(helpView)

	// Join main view and help vertically
	return lipgloss.JoinVertical(lipgloss.Left, mainView, styledHelp)
}
