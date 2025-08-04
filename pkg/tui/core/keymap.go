package core

import "github.com/charmbracelet/bubbles/v2/key"

type keyMap struct {
	Quit      key.Binding
	FocusNext key.Binding
}

var DefaultKeyMap = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q", "esc"),
		key.WithHelp("ctrl+c, q, esc", "quit"),
	),
	FocusNext: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "focus next"),
	),
}
