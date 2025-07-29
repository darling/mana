package tui

import tea "github.com/charmbracelet/bubbletea"

func IsQuitKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "q", "esc", "ctrl+c":
		return true
	}
	return false
}
