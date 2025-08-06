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

type sidebarKeyMap struct {
	FocusUp   key.Binding
	FocusDown key.Binding
	Enter     key.Binding
	Create    key.Binding
}

var DefaultSidebarKeyMap = sidebarKeyMap{
	FocusUp: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "focus up"),
	),
	FocusDown: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "focus down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "create new"),
	),
}

type mainKeyMap struct {
	Redraw key.Binding
	Create key.Binding
}

var DefaultMainKeyMap = mainKeyMap{
	Redraw: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "redraw"),
	),
	Create: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "create prompt"),
	),
}
