package layout

import (
	tea "github.com/charmbracelet/bubbletea/v2"
)

type Sizeable interface {
	SetSize(width, height int) tea.Cmd
	GetSize() (int, int)
}

type Contentable interface {
	SetContent(content string) tea.Cmd
	GetContent() string
}

// ComponentSizeMsg is sent to components to inform them of their available space.
type ComponentSizeMsg struct {
	Width  int
	Height int
}
