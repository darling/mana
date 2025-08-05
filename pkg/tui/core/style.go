package core

import "github.com/charmbracelet/lipgloss/v2"

var (
	// Colors using ANSI terminal colors
	subtle    = lipgloss.Color("8") // Bright Black (Dark Gray)
	highlight = lipgloss.Color("5") // Magenta
	special   = lipgloss.Color("2") // Green

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
