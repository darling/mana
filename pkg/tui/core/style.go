package core

import "github.com/charmbracelet/lipgloss/v2"

var (
	// Colors
	subtle    = lipgloss.Color("#383838")
	highlight = lipgloss.Color("#7D56F4") 
	special   = lipgloss.Color("#73F59F")

	// Styles for components
	FocusedBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(highlight).
			Padding(0, 1)

	BlurredBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtle).
			Padding(0, 1)

	FocusedItem = lipgloss.NewStyle().Foreground(special)

	// List header style
	ListHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtle).
			MarginBottom(1)
)