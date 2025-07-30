package tui

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	StatusBarBackground lipgloss.Color
	StatusBarForeground lipgloss.Color

	BorderNormal  lipgloss.Color
	BorderFocused lipgloss.Color

	ContentForegroundActive   lipgloss.AdaptiveColor
	ContentForegroundInactive lipgloss.AdaptiveColor
}

var DefaultTheme = Theme{
	StatusBarBackground: lipgloss.Color("235"),
	StatusBarForeground: lipgloss.Color("255"),
	BorderNormal:        lipgloss.Color("240"),
	BorderFocused:       lipgloss.Color("39"),
	ContentForegroundActive: lipgloss.AdaptiveColor{
		Light: "15",
		Dark:  "15",
	},
	ContentForegroundInactive: lipgloss.AdaptiveColor{
		Light: "245",
		Dark:  "245",
	},
}

var CurrentTheme = DefaultTheme

func StatusBarBg() lipgloss.Color               { return CurrentTheme.StatusBarBackground }
func StatusBarFg() lipgloss.Color               { return CurrentTheme.StatusBarForeground }
func BorderNormal() lipgloss.Color              { return CurrentTheme.BorderNormal }
func BorderFocused() lipgloss.Color             { return CurrentTheme.BorderFocused }
func ContentFgActive() lipgloss.AdaptiveColor   { return CurrentTheme.ContentForegroundActive }
func ContentFgInactive() lipgloss.AdaptiveColor { return CurrentTheme.ContentForegroundInactive }

func SetTheme(theme Theme) {
	CurrentTheme = theme
}

