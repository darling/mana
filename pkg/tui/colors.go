package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// BaseColors represents semantic colors for UI components
type BaseColors struct {
	Background lipgloss.AdaptiveColor
	Foreground lipgloss.AdaptiveColor

	Primary   lipgloss.AdaptiveColor
	Secondary lipgloss.AdaptiveColor
	Muted     lipgloss.AdaptiveColor

	Success lipgloss.AdaptiveColor
	Warning lipgloss.AdaptiveColor
	Danger  lipgloss.AdaptiveColor

	BorderNormal  lipgloss.AdaptiveColor
	BorderFocused lipgloss.AdaptiveColor
}

// Theme holds all color definitions
type Theme struct {
	BaseColors

	StatusBar struct {
		Background lipgloss.AdaptiveColor
		Foreground lipgloss.AdaptiveColor
	}

	Pane struct {
		Border struct {
			Normal  lipgloss.AdaptiveColor
			Focused lipgloss.AdaptiveColor
		}
		Content struct {
			Active   lipgloss.AdaptiveColor
			Inactive lipgloss.AdaptiveColor
		}
	}
}

// DefaultTheme uses standard terminal color palette
var DefaultTheme = Theme{
	BaseColors: BaseColors{
		Background: lipgloss.AdaptiveColor{Light: "", Dark: ""},
		Foreground: lipgloss.AdaptiveColor{Light: "15", Dark: "7"},

		Primary:   lipgloss.AdaptiveColor{Light: "12", Dark: "12"}, // Blue
		Secondary: lipgloss.AdaptiveColor{Light: "8", Dark: "8"},   // Gray
		Muted:     lipgloss.AdaptiveColor{Light: "7", Dark: "242"}, // Light gray

		Success: lipgloss.AdaptiveColor{Light: "10", Dark: "10"}, // Green
		Warning: lipgloss.AdaptiveColor{Light: "11", Dark: "11"}, // Yellow
		Danger:  lipgloss.AdaptiveColor{Light: "9", Dark: "9"},   // Red

		BorderNormal:  lipgloss.AdaptiveColor{Light: "8", Dark: "8"},
		BorderFocused: lipgloss.AdaptiveColor{Light: "12", Dark: "12"},
	},
}

var CurrentTheme Theme

func init() {
	CurrentTheme = DefaultTheme
	CurrentTheme.deriveComponentColors()
}

func (t *Theme) deriveComponentColors() {
	t.StatusBar.Background = t.BaseColors.Background
	t.StatusBar.Foreground = t.BaseColors.Foreground

	t.Pane.Border.Normal = t.BaseColors.BorderNormal
	t.Pane.Border.Focused = t.BaseColors.BorderFocused
	t.Pane.Content.Active = t.BaseColors.Foreground
	t.Pane.Content.Inactive = t.BaseColors.Muted
}

// Accessor functions for component-specific colors
func StatusBarBg() lipgloss.AdaptiveColor       { return CurrentTheme.StatusBar.Background }
func StatusBarFg() lipgloss.AdaptiveColor       { return CurrentTheme.StatusBar.Foreground }
func BorderNormal() lipgloss.AdaptiveColor      { return CurrentTheme.Pane.Border.Normal }
func BorderFocused() lipgloss.AdaptiveColor     { return CurrentTheme.Pane.Border.Focused }
func ContentFgActive() lipgloss.AdaptiveColor   { return CurrentTheme.Pane.Content.Active }
func ContentFgInactive() lipgloss.AdaptiveColor { return CurrentTheme.Pane.Content.Inactive }

// SetTheme allows changing the active theme
func SetTheme(theme Theme) {
	CurrentTheme = theme
	CurrentTheme.deriveComponentColors()
}
