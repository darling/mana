package tui

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// Dialog messages
type (
	OpenDialogMsg  struct{ Model tea.Model }
	CloseDialogMsg struct{}
	FocusDialogMsg struct{ Focused bool }
)

// DialogModel manages dialog stack
type DialogModel struct {
	dialogs []tea.Model
	width   int
	height  int
}

func NewDialogModel() DialogModel {
	return DialogModel{dialogs: []tea.Model{}}
}

func (m DialogModel) HasDialogs() bool { return len(m.dialogs) > 0 }

func (m *DialogModel) OpenDialog(model tea.Model) tea.Cmd {
	m.dialogs = append(m.dialogs, model)
	ws := tea.WindowSizeMsg{Width: m.width, Height: m.height}
	return func() tea.Msg { return ws }
}

func (m *DialogModel) CloseDialog() tea.Cmd {
	if len(m.dialogs) > 0 {
		m.dialogs = m.dialogs[:len(m.dialogs)-1]
	}
	return nil
}

func (m *DialogModel) Init() tea.Cmd { return nil }

func (m *DialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case OpenDialogMsg:
		if m.HasDialogs() {
			return m, nil
		}
		return m, m.OpenDialog(msg.Model)
	case CloseDialogMsg, tea.KeyMsg: // Handle ESC
		if key, ok := msg.(tea.KeyMsg); ok && key.String() != "esc" {
			break // Forward to dialog only
		}
		return m, m.CloseDialog()
	}

	// Forward messages to top dialog
	if m.HasDialogs() {
		var cmd tea.Cmd
		topModel := m.dialogs[len(m.dialogs)-1]
		m.dialogs[len(m.dialogs)-1], cmd = topModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m DialogModel) View() string {
	if !m.HasDialogs() {
		return ""
	}

	type Viewer interface {
		View() string
	}

	if v, ok := m.dialogs[len(m.dialogs)-1].(Viewer); ok {
		return v.View()
	}
	return ""
}

func (m *DialogModel) SetSize(width, height int) {
	m.width, m.height = width, height
}

// RenderWithBase creates a modal dialog overlay using lipgloss layers
func (m DialogModel) RenderWithBase(baseView string) string {
	if !m.HasDialogs() {
		return baseView
	}

	// Create base layer
	baseLayer := lipgloss.NewLayer(baseView)

	// Create dialog layers with increasing Z indices
	var layers []*lipgloss.Layer
	layers = append(layers, baseLayer)

	for i, dialog := range m.dialogs {
		if v, ok := dialog.(interface{ View() string }); ok {
			dialogView := v.View()

			// Center the dialog
			x := (m.width - lipgloss.Width(dialogView)) / 2
			y := (m.height - lipgloss.Height(dialogView)) / 2

			// Create semi-transparent overlay for the topmost dialog
			if i == len(m.dialogs)-1 {
				overlay := m.createOverlay()
				overlayLayer := lipgloss.NewLayer(overlay).Z(i + 1)
				layers = append(layers, overlayLayer)
			}

			// Add dialog layer with higher Z index
			dialogLayer := lipgloss.NewLayer(dialogView).
				X(x).
				Y(y).
				Z(i + 2) // Ensure dialog is above overlay
			layers = append(layers, dialogLayer)
		}
	}

	// Create canvas with all layers
	canvas := lipgloss.NewCanvas(layers...)
	return canvas.Render()
}

// Create a semi-transparent overlay
func (m DialogModel) createOverlay() string {
	return ""
}
