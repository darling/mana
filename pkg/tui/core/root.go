package core

import (
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/darling/mana/pkg/chat"
	"github.com/darling/mana/pkg/tui/core/components"
	"github.com/darling/mana/pkg/tui/core/layout"
)

type RootCmp interface {
	components.Component
}

type rootCmp struct {
	statusbar components.Component

	keys keyMap

	focusManager layout.FocusManager
	layerManager *layout.LayerManager

	width, height int

	chatService chat.Service
}

func NewRootCmp(service chat.Service) RootCmp {
	sidebar := NewSidebarCmp()
	main := NewMainCmp(service)
	statusbar := NewStatusBarCmp("v0.1.0")

	focusables := []layout.Focusable{sidebar.Clone(), main.Clone()}

	fm := layout.NewFocusManager(focusables, true)
	// Focus the main panel (index 1) by default before first render
	fm, _, _ = fm.Focus(1)

	return rootCmp{
		statusbar:    statusbar,
		keys:         DefaultKeyMap,
		focusManager: fm,
		layerManager: layout.NewLayerManager(),
		chatService:  service,
	}
}

func (m rootCmp) Init() tea.Cmd {
	return m.getHelpCmd()
}

func (m rootCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.layerManager.SetSize(msg.Width, msg.Height)
		m.focusManager, cmd = m.focusManager.UpdateAll(layout.ComponentSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - 1, // Account for status bar
		})
		cmds = append(cmds, cmd)
		newStatusBar, cmd := m.statusbar.Update(layout.ComponentSizeMsg{Width: msg.Width})
		m.statusbar = newStatusBar.(components.Component)
		cmds = append(cmds, cmd)

	case layout.FocusChangedMsg:
		return m, m.getHelpCmd()

	case layout.OpenLayerMsg:
		cmd = m.layerManager.Push(msg.Layer)
		cmds = append(cmds, cmd, m.getHelpCmd())

	case layout.DismissLayerMsg:
		cmd = m.layerManager.Pop()
		cmds = append(cmds, cmd, m.getHelpCmd())

	case layout.DismissLayerByIDMsg:
		cmd = m.layerManager.PopByID(msg.ID)
		cmds = append(cmds, cmd, m.getHelpCmd())

	case layout.ShowConfirmDialogMsg:
		dialog := layout.NewConfirmDialog(msg.Text)
		cmd = m.layerManager.Push(dialog)
		cmds = append(cmds, cmd, m.getHelpCmd())

	case layout.ShowPromptDialogMsg:
		dialog := layout.NewPromptDialog("")
		cmd = m.layerManager.Push(dialog)
		cmds = append(cmds, cmd, m.getHelpCmd())

	case layout.ConfirmedMsg:
		cmd = m.layerManager.Pop()
		cmds = append(cmds, cmd, m.getHelpCmd())

	case layout.CancelledMsg:
		cmd = m.layerManager.Pop()
		cmds = append(cmds, cmd, m.getHelpCmd())

	case layout.PromptSubmittedMsg:
		// Dismiss the prompt layer and forward the message to the focused component (main view)
		cmd = m.layerManager.Pop()
		cmds = append(cmds, cmd, m.getHelpCmd())
		m.focusManager, cmd = m.focusManager.UpdateFocused(msg)
		cmds = append(cmds, cmd)

	case tea.KeyPressMsg:
		// First try layer manager
		var handled bool
		m.layerManager, cmd, handled = m.layerManager.Update(msg)
		if handled {
			cmds = append(cmds, cmd, m.getHelpCmd())
		} else {
			// Fall back to base key handling
			m, cmd = m.handleKeyPress(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

	// NOTE: Unhandled types get passed to layers first, then focused component
	default:
		// Try layers first
		var handled bool
		m.layerManager, cmd, handled = m.layerManager.Update(msg)
		if handled {
			cmds = append(cmds, cmd)
		} else {
			// Fall back to focused component
			m.focusManager, cmd = m.focusManager.UpdateFocused(msg)
			cmds = append(cmds, cmd)
		}
	}

	// Messages	always go to the status bar
	var sbCmd tea.Cmd
	newStatusBar, sbCmd := m.statusbar.Update(msg)
	m.statusbar = newStatusBar.(components.Component)
	cmds = append(cmds, sbCmd)

	return m, tea.Batch(cmds...)
}

func (m rootCmp) View() string {
	sidebar, err := m.focusManager.Get(0)
	if err != nil {
		return "Error retrieving sidebar: " + err.Error()
	}
	main, err := m.focusManager.Get(1)
	if err != nil {
		return "Error retrieving main view: " + err.Error()
	}

	// First row: sidebar + main
	top := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebar.View(),
		main.View(),
	)

	// Second row: status bar. Force a single-line status regardless of content above.
	status := lipgloss.NewStyle().Width(m.width).MaxWidth(m.width).Height(1).MaxHeight(1).Render(m.statusbar.View())
	base := lipgloss.JoinVertical(
		lipgloss.Left,
		top,
		status,
	)

	// Let layers render over the base content
	return m.layerManager.RenderOver(base)
}

func (m rootCmp) handleKeyPress(msg tea.KeyPressMsg) (rootCmp, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.FocusNext):
		m.focusManager, cmd = m.focusManager.FocusNext()
		return m, cmd
	default:
		m.focusManager, cmd = m.focusManager.UpdateFocused(msg)
		return m, cmd
	}
}

func (m rootCmp) getHelpCmd() tea.Cmd {
	var bindings []key.Binding

	// First check if there's an active layer with help bindings
	if layerBindings := m.layerManager.HelpBindings(); len(layerBindings) > 0 {
		bindings = append(bindings, layerBindings...)
	} else {
		// Get bindings from the focused component if no layer is active
		focused, err := m.focusManager.GetFocused()
		if err == nil {
			if helpable, ok := focused.(layout.Help); ok {
				bindings = append(bindings, helpable.Bindings()...)
			}
		}
	}

	// Add global key bindings (unless a modal layer is active)
	if top := m.layerManager.Top(); top == nil || !top.LayerMeta().Modal {
		bindings = append(bindings, m.keys.FocusNext, m.keys.Quit)
	}

	return func() tea.Msg {
		return layout.HelpUpdateMsg(bindings)
	}
}
