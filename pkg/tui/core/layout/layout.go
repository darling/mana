package layout

import (
	"github.com/charmbracelet/bubbles/v2/key"
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

type Help interface {
	Bindings() []key.Binding
}
