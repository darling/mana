package tui

import (
	tea "github.com/charmbracelet/bubbletea/v2"
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
	return nil
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

