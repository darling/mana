package tui

import (
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Tab      key.Binding
	Quit     key.Binding
	Enter    key.Binding
	Settings key.Binding
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
	Settings: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "settings"),
	),
}

type RootModel struct {
	width        int
	height       int
	sidebar      SidebarModel
	content      ContentModel
	statusbar    StatusbarModel
	dialogs      DialogModel
	focusIdx     int
	lastFocusIdx int
	unfocused    bool
	keys         keyMap
}

func NewRootModel() RootModel {
	sidebar := NewSidebarModel()
	content := NewContentModel()
	statusbar := NewStatusbarModel()

	sidebar.focused = false
	content.focused = true

	return RootModel{
		sidebar:      sidebar,
		content:      content,
		statusbar:    statusbar,
		dialogs:      NewDialogModel(),
		focusIdx:     1,
		lastFocusIdx: 1,
		unfocused:    false,
		keys:         keys,
	}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m *RootModel) hasModels() bool {
	return len(m.sidebar.items) > 0
}

func (m *RootModel) updateFocusState() {
	hasModels := m.hasModels()
	hasDialogs := m.dialogs.HasDialogs()
	shouldUnfocus := !hasModels || hasDialogs
	
	if shouldUnfocus && !m.unfocused {
		m.lastFocusIdx = m.focusIdx
		m.unfocused = true
		m.sidebar.focused = false
		m.content.focused = false
	} else if !shouldUnfocus && m.unfocused {
		m.unfocused = false
		m.focusIdx = m.lastFocusIdx
		m.sidebar.focused = m.focusIdx == 0
		m.content.focused = m.focusIdx == 1
	}
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case OpenDialogMsg:
		_, cmd := m.dialogs.Update(msg)
		cmds = append(cmds, cmd)
		m.updateFocusState()
		return m, tea.Batch(cmds...)
	}

	if m.dialogs.HasDialogs() {
		_, cmd := m.dialogs.Update(msg)
		m.updateFocusState()
		return m, cmd
	}

	m.updateFocusState()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Tab) && !m.unfocused:

			m.focusIdx = (m.focusIdx + 1) % 2

			m.sidebar.focused = m.focusIdx == 0
			m.content.focused = m.focusIdx == 1

			return m, nil
		case key.Matches(msg, m.keys.Settings):
			return m, func() tea.Msg {
				return OpenDialogMsg{Model: NewConfigDialog()}
			}
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.dialogs.SetSize(msg.Width, msg.Height)

		statusbarHeight := 1
		mainViewHeight := m.height - statusbarHeight
		sidebarWidth := int(float64(m.width) * 0.3)
		if sidebarWidth < 20 {
			sidebarWidth = 20
		}
		contentWidth := m.width - sidebarWidth

		m.statusbar.SetSize(m.width)

		m.sidebar.width = sidebarWidth
		m.sidebar.height = mainViewHeight
		m.sidebar.handleResize(sidebarWidth, mainViewHeight)

		m.content.width = contentWidth
		m.content.height = mainViewHeight
		m.content.handleResize(contentWidth, mainViewHeight)

		if m.dialogs.HasDialogs() {
			if top := m.dialogs.dialogs[len(m.dialogs.dialogs)-1]; top != nil {
				var cmd tea.Cmd
				top, cmd = top.Update(msg)
				m.dialogs.dialogs[len(m.dialogs.dialogs)-1] = top
				cmds = append(cmds, cmd)
			}
		}

		return m, tea.Batch(cmds...)
	}

	var sbCmd, ctCmd, stCmd tea.Cmd

	if !m.unfocused {
		if m.focusIdx == 0 {
			m.sidebar, sbCmd = m.sidebar.Update(msg)
			cmds = append(cmds, sbCmd)
		} else {
			m.content, ctCmd = m.content.Update(msg)
			cmds = append(cmds, ctCmd)
		}
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
	baseView := lipgloss.JoinVertical(lipgloss.Left, mainView, statusbarView)

	// Use the new layer-based rendering
	return m.dialogs.RenderWithBase(baseView)
}
