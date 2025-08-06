package layout

import "github.com/charmbracelet/bubbles/v2/key"

type Help interface {
	Bindings() []key.Binding
}

type HelpUpdateMsg []key.Binding
