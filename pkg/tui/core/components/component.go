package components

import tea "github.com/charmbracelet/bubbletea/v2"

type Component interface {
	tea.Model
	tea.ViewModel
}
