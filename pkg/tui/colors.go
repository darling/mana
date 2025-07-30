package tui

import (
	"image/color"

	"github.com/charmbracelet/lipgloss/v2"
)

type BaseColors struct {
	Background color.Color
	Foreground color.Color

	Primary   color.Color
	Secondary color.Color
	Muted     color.Color

	Success color.Color
	Warning color.Color
	Danger  color.Color

	BorderNormal  color.Color
	BorderFocused color.Color
}

type Theme struct {
	BaseColors

	StatusBar struct {
		Background color.Color
		Foreground color.Color
	}

	Pane struct {
		Border struct {
			Normal  color.Color
			Focused color.Color
		}
		Content struct {
			Active   color.Color
			Inactive color.Color
		}
	}
}

var DefaultTheme = Theme{
	BaseColors: BaseColors{
		Background: lipgloss.NoColor{},
		Foreground: lipgloss.ANSIColor(7),

		Primary:   lipgloss.ANSIColor(12),
		Secondary: lipgloss.ANSIColor(8),
		Muted:     lipgloss.ANSIColor(7),

		Success: lipgloss.ANSIColor(10),
		Warning: lipgloss.ANSIColor(11),
		Danger:  lipgloss.ANSIColor(9),

		BorderNormal:  lipgloss.ANSIColor(8),
		BorderFocused: lipgloss.ANSIColor(12),
	},
}

var CurrentTheme Theme

func init() {
	CurrentTheme = DefaultTheme
	CurrentTheme.deriveComponentColors()
}

func (t *Theme) deriveComponentColors() {
	t.StatusBar.Background = t.BaseColors.Background
	t.StatusBar.Foreground = t.BaseColors.Muted

	t.Pane.Border.Normal = t.BaseColors.BorderNormal
	t.Pane.Border.Focused = t.BaseColors.BorderFocused
	t.Pane.Content.Active = t.BaseColors.Foreground
	t.Pane.Content.Inactive = t.BaseColors.Muted
}

func StatusBarBg() color.Color       { return CurrentTheme.StatusBar.Background }
func StatusBarFg() color.Color       { return CurrentTheme.StatusBar.Foreground }
func BorderNormal() color.Color      { return CurrentTheme.Pane.Border.Normal }
func BorderFocused() color.Color     { return CurrentTheme.Pane.Border.Focused }
func ContentFgActive() color.Color   { return CurrentTheme.Pane.Content.Active }
func ContentFgInactive() color.Color { return CurrentTheme.Pane.Content.Inactive }

func SetTheme(theme Theme) {
	CurrentTheme = theme
	CurrentTheme.deriveComponentColors()
}
